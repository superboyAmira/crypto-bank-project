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

type ExchangeService struct {
	exchangeRepo *repositories.ExchangeRepository
	accountRepo  *repositories.AccountRepository
	walletRepo   *repositories.CryptoWalletRepository
	txRepo       *repositories.TransactionRepository
	db           *sql.DB
	rabbitMQ     *rabbitmq.Client
}

func NewExchangeService(
	exchangeRepo *repositories.ExchangeRepository,
	accountRepo *repositories.AccountRepository,
	walletRepo *repositories.CryptoWalletRepository,
	txRepo *repositories.TransactionRepository,
	db *sql.DB,
	rabbitMQ *rabbitmq.Client,
) *ExchangeService {
	return &ExchangeService{
		exchangeRepo: exchangeRepo,
		accountRepo:  accountRepo,
		walletRepo:   walletRepo,
		txRepo:       txRepo,
		db:           db,
		rabbitMQ:     rabbitMQ,
	}
}

// ExchangeCryptoToFiat exchanges cryptocurrency to fiat currency
func (s *ExchangeService) ExchangeCryptoToFiat(req *models.ExchangeCryptoToFiatRequest) (*models.Exchange, error) {
	logger.Info("Exchanging crypto to fiat",
		zap.String("user_id", req.UserID.String()),
		zap.String("from_wallet", req.FromWalletID.String()),
		zap.String("to_account", req.ToAccountID.String()),
		zap.Float64("crypto_amount", req.CryptoAmount),
	)

	// Start database transaction
	dbTx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer dbTx.Rollback()

	// Get wallet and account
	wallet, err := s.walletRepo.GetByID(req.FromWalletID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}

	account, err := s.accountRepo.GetByID(req.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	// Verify ownership
	if wallet.UserID != req.UserID || account.UserID != req.UserID {
		return nil, fmt.Errorf("ownership mismatch")
	}

	// Check balance
	if wallet.Balance < req.CryptoAmount {
		return nil, fmt.Errorf("insufficient balance: have %f, need %f", wallet.Balance, req.CryptoAmount)
	}

	// Get exchange rate
	fromCurrency := string(wallet.CryptoType)
	toCurrency := string(account.Currency)
	rate, err := s.exchangeRepo.GetExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	// Calculate fiat amount
	fiatAmount := req.CryptoAmount * rate

	// Create exchange record
	exchange := &models.Exchange{
		UserID:        req.UserID,
		Type:          models.ExchangeCryptoToFiat,
		Status:        models.ExchangeStatusPending,
		FromCurrency:  fromCurrency,
		ToCurrency:    toCurrency,
		FromAmount:    req.CryptoAmount,
		ToAmount:      fiatAmount,
		ExchangeRate:  rate,
		FromWalletID:  &req.FromWalletID,
		ToAccountID:   &req.ToAccountID,
	}

	if err := s.exchangeRepo.Create(exchange); err != nil {
		return nil, fmt.Errorf("failed to create exchange: %w", err)
	}

	// Update balances
	if err := s.walletRepo.UpdateBalance(req.FromWalletID, -req.CryptoAmount); err != nil {
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	if err := s.accountRepo.UpdateBalance(req.ToAccountID, fiatAmount); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	// Create transaction record
	transaction := &models.Transaction{
		UserID:      req.UserID,
		Type:        models.TransactionTypeExchange,
		Status:      models.TransactionStatusCompleted,
		Amount:      fiatAmount,
		Currency:    toCurrency,
		ToAccountID: &req.ToAccountID,
		ExchangeID:  &exchange.ID,
		Description: fmt.Sprintf("Exchange %f %s to %s", req.CryptoAmount, fromCurrency, toCurrency),
	}

	if err := s.txRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update exchange status
	if err := s.exchangeRepo.UpdateStatus(exchange.ID, models.ExchangeStatusCompleted); err != nil {
		return nil, fmt.Errorf("failed to update exchange status: %w", err)
	}

	// Commit database transaction
	if err := dbTx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	exchange.Status = models.ExchangeStatusCompleted
	exchange.TransactionID = &transaction.ID

	// Update metrics
	metrics.ExchangesTotal.WithLabelValues(string(exchange.Type), string(exchange.Status)).Inc()

	// Publish event
	event := rabbitmq.ExchangeEvent{
		ExchangeID:   exchange.ID.String(),
		UserID:       exchange.UserID.String(),
		FromCurrency: exchange.FromCurrency,
		ToCurrency:   exchange.ToCurrency,
		FromAmount:   exchange.FromAmount,
		ToAmount:     exchange.ToAmount,
		Status:       string(exchange.Status),
	}
	s.rabbitMQ.PublishEvent(rabbitmq.ExchangeEvents, rabbitmq.EventExchangeCompleted, event)

	logger.Info("Crypto to fiat exchange completed", zap.String("exchange_id", exchange.ID.String()))
	return exchange, nil
}

// ExchangeFiatToCrypto exchanges fiat currency to cryptocurrency
func (s *ExchangeService) ExchangeFiatToCrypto(req *models.ExchangeFiatToCryptoRequest) (*models.Exchange, error) {
	logger.Info("Exchanging fiat to crypto",
		zap.String("user_id", req.UserID.String()),
		zap.String("from_account", req.FromAccountID.String()),
		zap.String("to_wallet", req.ToWalletID.String()),
		zap.Float64("fiat_amount", req.FiatAmount),
	)

	// Start database transaction
	dbTx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer dbTx.Rollback()

	// Get account and wallet
	account, err := s.accountRepo.GetByID(req.FromAccountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	wallet, err := s.walletRepo.GetByID(req.ToWalletID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}

	// Verify ownership
	if account.UserID != req.UserID || wallet.UserID != req.UserID {
		return nil, fmt.Errorf("ownership mismatch")
	}

	// Check balance
	if account.Balance < req.FiatAmount {
		return nil, fmt.Errorf("insufficient balance: have %f, need %f", account.Balance, req.FiatAmount)
	}

	// Get exchange rate
	fromCurrency := string(account.Currency)
	toCurrency := string(wallet.CryptoType)
	rate, err := s.exchangeRepo.GetExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	// Calculate crypto amount
	cryptoAmount := req.FiatAmount * rate

	// Create exchange record
	exchange := &models.Exchange{
		UserID:        req.UserID,
		Type:          models.ExchangeFiatToCrypto,
		Status:        models.ExchangeStatusPending,
		FromCurrency:  fromCurrency,
		ToCurrency:    toCurrency,
		FromAmount:    req.FiatAmount,
		ToAmount:      cryptoAmount,
		ExchangeRate:  rate,
		FromAccountID: &req.FromAccountID,
		ToWalletID:    &req.ToWalletID,
	}

	if err := s.exchangeRepo.Create(exchange); err != nil {
		return nil, fmt.Errorf("failed to create exchange: %w", err)
	}

	// Update balances
	if err := s.accountRepo.UpdateBalance(req.FromAccountID, -req.FiatAmount); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	if err := s.walletRepo.UpdateBalance(req.ToWalletID, cryptoAmount); err != nil {
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Create transaction record
	transaction := &models.Transaction{
		UserID:        req.UserID,
		Type:          models.TransactionTypeExchange,
		Status:        models.TransactionStatusCompleted,
		Amount:        req.FiatAmount,
		Currency:      fromCurrency,
		FromAccountID: &req.FromAccountID,
		ExchangeID:    &exchange.ID,
		Description:   fmt.Sprintf("Exchange %f %s to %s", req.FiatAmount, fromCurrency, toCurrency),
	}

	if err := s.txRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update exchange status
	if err := s.exchangeRepo.UpdateStatus(exchange.ID, models.ExchangeStatusCompleted); err != nil {
		return nil, fmt.Errorf("failed to update exchange status: %w", err)
	}

	// Commit database transaction
	if err := dbTx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	exchange.Status = models.ExchangeStatusCompleted
	exchange.TransactionID = &transaction.ID

	// Update metrics
	metrics.ExchangesTotal.WithLabelValues(string(exchange.Type), string(exchange.Status)).Inc()

	// Publish event
	event := rabbitmq.ExchangeEvent{
		ExchangeID:   exchange.ID.String(),
		UserID:       exchange.UserID.String(),
		FromCurrency: exchange.FromCurrency,
		ToCurrency:   exchange.ToCurrency,
		FromAmount:   exchange.FromAmount,
		ToAmount:     exchange.ToAmount,
		Status:       string(exchange.Status),
	}
	s.rabbitMQ.PublishEvent(rabbitmq.ExchangeEvents, rabbitmq.EventExchangeCompleted, event)

	logger.Info("Fiat to crypto exchange completed", zap.String("exchange_id", exchange.ID.String()))
	return exchange, nil
}

// GetExchange retrieves an exchange by ID
func (s *ExchangeService) GetExchange(id uuid.UUID) (*models.Exchange, error) {
	return s.exchangeRepo.GetByID(id)
}

// GetUserExchanges retrieves all exchanges for a user
func (s *ExchangeService) GetUserExchanges(userID uuid.UUID) ([]*models.Exchange, error) {
	return s.exchangeRepo.GetByUserID(userID)
}

