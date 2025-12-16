package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/crypto-bank/bank-service/internal/repositories"
	"github.com/crypto-bank/bank-service/pkg/logger"
	"github.com/crypto-bank/bank-service/pkg/rabbitmq"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CryptoWalletService struct {
	walletRepo *repositories.CryptoWalletRepository
	userRepo   *repositories.UserRepository
	rabbitMQ   *rabbitmq.Client
}

func NewCryptoWalletService(
	walletRepo *repositories.CryptoWalletRepository,
	userRepo *repositories.UserRepository,
	rabbitMQ *rabbitmq.Client,
) *CryptoWalletService {
	return &CryptoWalletService{
		walletRepo: walletRepo,
		userRepo:   userRepo,
		rabbitMQ:   rabbitMQ,
	}
}

// CreateWallet creates a new crypto wallet
func (s *CryptoWalletService) CreateWallet(req *models.CreateCryptoWalletRequest) (*models.CryptoWallet, error) {
	logger.Info("Creating crypto wallet",
		zap.String("user_id", req.UserID.String()),
		zap.String("crypto_type", string(req.CryptoType)),
	)

	// Verify user exists
	user, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate wallet address
	address := s.generateWalletAddress(user.ID, req.CryptoType)

	wallet := &models.CryptoWallet{
		UserID:     req.UserID,
		CryptoType: req.CryptoType,
		Balance:    0,
		Address:    address,
	}

	if err := s.walletRepo.Create(wallet); err != nil {
		logger.Error("Failed to create crypto wallet", zap.Error(err))
		return nil, err
	}

	// Publish event
	event := rabbitmq.WalletEvent{
		WalletID:   wallet.ID.String(),
		UserID:     wallet.UserID.String(),
		CryptoType: string(wallet.CryptoType),
	}
	s.rabbitMQ.PublishEvent(rabbitmq.ExchangeEvents, rabbitmq.EventWalletCreated, event)

	logger.Info("Crypto wallet created", zap.String("wallet_id", wallet.ID.String()))
	return wallet, nil
}

// GetWallet retrieves a crypto wallet by ID
func (s *CryptoWalletService) GetWallet(id uuid.UUID) (*models.CryptoWallet, error) {
	return s.walletRepo.GetByID(id)
}

// GetUserWallets retrieves all crypto wallets for a user
func (s *CryptoWalletService) GetUserWallets(userID uuid.UUID) ([]*models.CryptoWallet, error) {
	return s.walletRepo.GetByUserID(userID)
}

// GetWalletBalance retrieves wallet balance
func (s *CryptoWalletService) GetWalletBalance(id uuid.UUID) (float64, error) {
	return s.walletRepo.GetBalance(id)
}

// generateWalletAddress generates a unique wallet address
func (s *CryptoWalletService) generateWalletAddress(userID uuid.UUID, cryptoType models.CryptoType) string {
	data := fmt.Sprintf("%s-%s-%s", userID.String(), cryptoType, uuid.New().String())
	hash := sha256.Sum256([]byte(data))
	
	// Add prefix based on crypto type
	var prefix string
	switch cryptoType {
	case models.CryptoBTC:
		prefix = "1"
	case models.CryptoETH:
		prefix = "0x"
	case models.CryptoUSDT:
		prefix = "0x"
	case models.CryptoBNB:
		prefix = "0x"
	case models.CryptoSOL:
		prefix = ""
	default:
		prefix = ""
	}
	
	return prefix + hex.EncodeToString(hash[:20])
}

