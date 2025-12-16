package models

import (
	"time"

	"github.com/google/uuid"
)

// CryptoType represents cryptocurrency types
type CryptoType string

const (
	CryptoBTC  CryptoType = "BTC"
	CryptoETH  CryptoType = "ETH"
	CryptoUSDT CryptoType = "USDT"
	CryptoBNB  CryptoType = "BNB"
	CryptoSOL  CryptoType = "SOL"
)

// CryptoWallet represents a cryptocurrency wallet
type CryptoWallet struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	CryptoType CryptoType `json:"crypto_type" db:"crypto_type" validate:"required"`
	Balance   float64    `json:"balance" db:"balance"`
	Address   string     `json:"address" db:"address"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateCryptoWalletRequest struct {
	UserID     uuid.UUID  `json:"user_id" validate:"required"`
	CryptoType CryptoType `json:"crypto_type" validate:"required,oneof=BTC ETH USDT BNB SOL"`
}

type CryptoWalletWithUser struct {
	CryptoWallet
	User User `json:"user"`
}

