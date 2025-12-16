package models

import (
	"time"

	"github.com/google/uuid"
)

// ExchangeType represents the direction of exchange
type ExchangeType string

const (
	ExchangeCryptoToFiat ExchangeType = "CRYPTO_TO_FIAT"
	ExchangeFiatToCrypto ExchangeType = "FIAT_TO_CRYPTO"
	ExchangeCryptoToCrypto ExchangeType = "CRYPTO_TO_CRYPTO"
)

// ExchangeStatus represents the status of exchange
type ExchangeStatus string

const (
	ExchangeStatusPending   ExchangeStatus = "PENDING"
	ExchangeStatusCompleted ExchangeStatus = "COMPLETED"
	ExchangeStatusFailed    ExchangeStatus = "FAILED"
)

// Exchange represents a currency exchange operation
type Exchange struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	UserID          uuid.UUID      `json:"user_id" db:"user_id"`
	Type            ExchangeType   `json:"type" db:"type"`
	Status          ExchangeStatus `json:"status" db:"status"`
	FromCurrency    string         `json:"from_currency" db:"from_currency"`
	ToCurrency      string         `json:"to_currency" db:"to_currency"`
	FromAmount      float64        `json:"from_amount" db:"from_amount"`
	ToAmount        float64        `json:"to_amount" db:"to_amount"`
	ExchangeRate    float64        `json:"exchange_rate" db:"exchange_rate"`
	FromAccountID   *uuid.UUID     `json:"from_account_id,omitempty" db:"from_account_id"`
	ToAccountID     *uuid.UUID     `json:"to_account_id,omitempty" db:"to_account_id"`
	FromWalletID    *uuid.UUID     `json:"from_wallet_id,omitempty" db:"from_wallet_id"`
	ToWalletID      *uuid.UUID     `json:"to_wallet_id,omitempty" db:"to_wallet_id"`
	TransactionID   *uuid.UUID     `json:"transaction_id,omitempty" db:"transaction_id"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// ExchangeCryptoToFiatRequest represents request to exchange crypto to fiat
type ExchangeCryptoToFiatRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required"`
	FromWalletID  uuid.UUID `json:"from_wallet_id" validate:"required"`
	ToAccountID   uuid.UUID `json:"to_account_id" validate:"required"`
	CryptoAmount  float64   `json:"crypto_amount" validate:"required,gt=0"`
}

// ExchangeFiatToCryptoRequest represents request to exchange fiat to crypto
type ExchangeFiatToCryptoRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required"`
	FromAccountID uuid.UUID `json:"from_account_id" validate:"required"`
	ToWalletID    uuid.UUID `json:"to_wallet_id" validate:"required"`
	FiatAmount    float64   `json:"fiat_amount" validate:"required,gt=0"`
}

// ExchangeRate represents current exchange rate
type ExchangeRate struct {
	ID           uuid.UUID `json:"id" db:"id"`
	FromCurrency string    `json:"from_currency" db:"from_currency"`
	ToCurrency   string    `json:"to_currency" db:"to_currency"`
	Rate         float64   `json:"rate" db:"rate"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

