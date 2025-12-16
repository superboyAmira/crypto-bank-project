package models

import (
	"time"

	"github.com/google/uuid"
)

// CurrencyType represents fiat currency types
type CurrencyType string

const (
	CurrencyUSD CurrencyType = "USD"
	CurrencyEUR CurrencyType = "EUR"
	CurrencyRUB CurrencyType = "RUB"
	CurrencyGBP CurrencyType = "GBP"
)

// Account represents a fiat currency account
type Account struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	UserID    uuid.UUID    `json:"user_id" db:"user_id"`
	Currency  CurrencyType `json:"currency" db:"currency" validate:"required"`
	Balance   float64      `json:"balance" db:"balance"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
}

type CreateAccountRequest struct {
	UserID   uuid.UUID    `json:"user_id" validate:"required"`
	Currency CurrencyType `json:"currency" validate:"required,oneof=USD EUR RUB GBP"`
}

type AccountWithUser struct {
	Account
	User User `json:"user"`
}

