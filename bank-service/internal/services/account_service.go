package services

import (
	"fmt"

	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/crypto-bank/bank-service/internal/repositories"
	"github.com/crypto-bank/bank-service/pkg/logger"
	"github.com/crypto-bank/bank-service/pkg/rabbitmq"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AccountService struct {
	accountRepo *repositories.AccountRepository
	userRepo    *repositories.UserRepository
	rabbitMQ    *rabbitmq.Client
}

func NewAccountService(
	accountRepo *repositories.AccountRepository,
	userRepo *repositories.UserRepository,
	rabbitMQ *rabbitmq.Client,
) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		userRepo:    userRepo,
		rabbitMQ:    rabbitMQ,
	}
}

// CreateAccount creates a new fiat account
func (s *AccountService) CreateAccount(req *models.CreateAccountRequest) (*models.Account, error) {
	logger.Info("Creating account",
		zap.String("user_id", req.UserID.String()),
		zap.String("currency", string(req.Currency)),
	)

	// Verify user exists
	_, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	account := &models.Account{
		UserID:   req.UserID,
		Currency: req.Currency,
		Balance:  0,
	}

	if err := s.accountRepo.Create(account); err != nil {
		logger.Error("Failed to create account", zap.Error(err))
		return nil, err
	}

	// Publish event
	event := rabbitmq.AccountEvent{
		AccountID: account.ID.String(),
		UserID:    account.UserID.String(),
		Currency:  string(account.Currency),
	}
	s.rabbitMQ.PublishEvent(rabbitmq.ExchangeEvents, rabbitmq.EventAccountCreated, event)

	logger.Info("Account created", zap.String("account_id", account.ID.String()))
	return account, nil
}

// GetAccount retrieves an account by ID
func (s *AccountService) GetAccount(id uuid.UUID) (*models.Account, error) {
	return s.accountRepo.GetByID(id)
}

// GetUserAccounts retrieves all accounts for a user
func (s *AccountService) GetUserAccounts(userID uuid.UUID) ([]*models.Account, error) {
	return s.accountRepo.GetByUserID(userID)
}

// GetAccountBalance retrieves account balance
func (s *AccountService) GetAccountBalance(id uuid.UUID) (float64, error) {
	return s.accountRepo.GetBalance(id)
}

