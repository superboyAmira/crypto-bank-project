package models

import (
	"time"

	"github.com/google/uuid"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeTransfer  TransactionType = "TRANSFER"
	TransactionTypeDeposit   TransactionType = "DEPOSIT"
	TransactionTypeWithdraw  TransactionType = "WITHDRAW"
	TransactionTypeExchange  TransactionType = "EXCHANGE"
)

// TransactionStatus represents the status of transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID                uuid.UUID         `json:"id" db:"id"`
	UserID            uuid.UUID         `json:"user_id" db:"user_id"`
	Type              TransactionType   `json:"type" db:"type"`
	Status            TransactionStatus `json:"status" db:"status"`
	Amount            float64           `json:"amount" db:"amount"`
	Currency          string            `json:"currency" db:"currency"`
	FromAccountID     *uuid.UUID        `json:"from_account_id,omitempty" db:"from_account_id"`
	ToAccountID       *uuid.UUID        `json:"to_account_id,omitempty" db:"to_account_id"`
	Description       string            `json:"description" db:"description"`
	ExchangeID        *uuid.UUID        `json:"exchange_id,omitempty" db:"exchange_id"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`
}

type CreateTransactionRequest struct {
	FromAccountID uuid.UUID `json:"from_account_id" validate:"required"`
	ToAccountID   uuid.UUID `json:"to_account_id" validate:"required"`
	Amount        float64   `json:"amount" validate:"required,gt=0"`
	Description   string    `json:"description"`
}

type DepositRequest struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required,gt=0"`
}

type WithdrawRequest struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required,gt=0"`
}

