package services

import (
	"database/sql"
	"fmt"

	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/crypto-bank/bank-service/internal/repositories"
	"github.com/crypto-bank/bank-service/pkg/logger"
	"github.com/crypto-bank/bank-service/pkg/metrics"
	"github.com/crypto-bank/bank-service/pkg/rabbitmq"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TransactionService struct {
	txRepo      *repositories.TransactionRepository
	accountRepo *repositories.AccountRepository
	db          *sql.DB
	rabbitMQ    *rabbitmq.Client
}

func NewTransactionService(
	txRepo *repositories.TransactionRepository,
	accountRepo *repositories.AccountRepository,
	db *sql.DB,
	rabbitMQ *rabbitmq.Client,
) *TransactionService {
	return &TransactionService{
		txRepo:      txRepo,
		accountRepo: accountRepo,
		db:          db,
		rabbitMQ:    rabbitMQ,
	}
}

// CreateTransfer creates a transfer transaction between accounts
func (s *TransactionService) CreateTransfer(req *models.CreateTransactionRequest) (*models.Transaction, error) {
	logger.Info("Creating transfer",
		zap.String("from_account", req.FromAccountID.String()),
		zap.String("to_account", req.ToAccountID.String()),
		zap.Float64("amount", req.Amount),
	)

	// Start database transaction
	dbTx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer dbTx.Rollback()

	// Get accounts
	fromAccount, err := s.accountRepo.GetByID(req.FromAccountID)
	if err != nil {
		return nil, fmt.Errorf("from account not found: %w", err)
	}

	toAccount, err := s.accountRepo.GetByID(req.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("to account not found: %w", err)
	}

	// Validate currency match
	if fromAccount.Currency != toAccount.Currency {
		return nil, fmt.Errorf("currency mismatch: from %s to %s", fromAccount.Currency, toAccount.Currency)
	}

	// Check balance
	if fromAccount.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance: have %f, need %f", fromAccount.Balance, req.Amount)
	}

	// Create transaction record
	transaction := &models.Transaction{
		UserID:        fromAccount.UserID,
		Type:          models.TransactionTypeTransfer,
		Status:        models.TransactionStatusPending,
		Amount:        req.Amount,
		Currency:      string(fromAccount.Currency),
		FromAccountID: &req.FromAccountID,
		ToAccountID:   &req.ToAccountID,
		Description:   req.Description,
	}

	if err := s.txRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update balances
	if err := s.accountRepo.UpdateBalance(req.FromAccountID, -req.Amount); err != nil {
		return nil, fmt.Errorf("failed to update from account balance: %w", err)
	}

	if err := s.accountRepo.UpdateBalance(req.ToAccountID, req.Amount); err != nil {
		return nil, fmt.Errorf("failed to update to account balance: %w", err)
	}

	// Update transaction status
	if err := s.txRepo.UpdateStatus(transaction.ID, models.TransactionStatusCompleted); err != nil {
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Commit database transaction
	if err := dbTx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	transaction.Status = models.TransactionStatusCompleted

	// Update metrics
	metrics.TransactionsTotal.WithLabelValues(string(transaction.Type), string(transaction.Status)).Inc()
	metrics.TransactionAmount.WithLabelValues(transaction.Currency).Observe(transaction.Amount)

	// Publish events
	event := rabbitmq.TransactionEvent{
		TransactionID: transaction.ID.String(),
		UserID:        transaction.UserID.String(),
		Type:          string(transaction.Type),
		Amount:        transaction.Amount,
		Currency:      transaction.Currency,
		Status:        string(transaction.Status),
	}
	s.rabbitMQ.PublishEvent(rabbitmq.ExchangeEvents, rabbitmq.EventTransactionCompleted, event)

	logger.Info("Transfer completed", zap.String("transaction_id", transaction.ID.String()))
	return transaction, nil
}

// Deposit deposits money to an account
func (s *TransactionService) Deposit(req *models.DepositRequest) (*models.Transaction, error) {
	logger.Info("Creating deposit",
		zap.String("account", req.AccountID.String()),
		zap.Float64("amount", req.Amount),
	)

	account, err := s.accountRepo.GetByID(req.AccountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	transaction := &models.Transaction{
		UserID:      account.UserID,
		Type:        models.TransactionTypeDeposit,
		Status:      models.TransactionStatusPending,
		Amount:      req.Amount,
		Currency:    string(account.Currency),
		ToAccountID: &req.AccountID,
		Description: "Deposit",
	}

	if err := s.txRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := s.accountRepo.UpdateBalance(req.AccountID, req.Amount); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	if err := s.txRepo.UpdateStatus(transaction.ID, models.TransactionStatusCompleted); err != nil {
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	transaction.Status = models.TransactionStatusCompleted

	// Update metrics
	metrics.TransactionsTotal.WithLabelValues(string(transaction.Type), string(transaction.Status)).Inc()
	metrics.TransactionAmount.WithLabelValues(transaction.Currency).Observe(transaction.Amount)

	logger.Info("Deposit completed", zap.String("transaction_id", transaction.ID.String()))
	return transaction, nil
}

// Withdraw withdraws money from an account
func (s *TransactionService) Withdraw(req *models.WithdrawRequest) (*models.Transaction, error) {
	logger.Info("Creating withdrawal",
		zap.String("account", req.AccountID.String()),
		zap.Float64("amount", req.Amount),
	)

	account, err := s.accountRepo.GetByID(req.AccountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	if account.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance: have %f, need %f", account.Balance, req.Amount)
	}

	transaction := &models.Transaction{
		UserID:        account.UserID,
		Type:          models.TransactionTypeWithdraw,
		Status:        models.TransactionStatusPending,
		Amount:        req.Amount,
		Currency:      string(account.Currency),
		FromAccountID: &req.AccountID,
		Description:   "Withdrawal",
	}

	if err := s.txRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := s.accountRepo.UpdateBalance(req.AccountID, -req.Amount); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	if err := s.txRepo.UpdateStatus(transaction.ID, models.TransactionStatusCompleted); err != nil {
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	transaction.Status = models.TransactionStatusCompleted

	// Update metrics
	metrics.TransactionsTotal.WithLabelValues(string(transaction.Type), string(transaction.Status)).Inc()
	metrics.TransactionAmount.WithLabelValues(transaction.Currency).Observe(transaction.Amount)

	logger.Info("Withdrawal completed", zap.String("transaction_id", transaction.ID.String()))
	return transaction, nil
}

// GetTransaction retrieves a transaction by ID
func (s *TransactionService) GetTransaction(id uuid.UUID) (*models.Transaction, error) {
	return s.txRepo.GetByID(id)
}

// GetUserTransactions retrieves all transactions for a user
func (s *TransactionService) GetUserTransactions(userID uuid.UUID) ([]*models.Transaction, error) {
	return s.txRepo.GetByUserID(userID)
}

