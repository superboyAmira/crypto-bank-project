package repositories

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/google/uuid"
)

type ExchangeRepository struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewExchangeRepository(db *sql.DB) *ExchangeRepository {
	return &ExchangeRepository{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create creates a new exchange
func (r *ExchangeRepository) Create(exchange *models.Exchange) error {
	exchange.ID = uuid.New()

	query := r.qb.Insert("exchanges").
		Columns("id", "user_id", "type", "status", "from_currency", "to_currency",
			"from_amount", "to_amount", "exchange_rate", "from_account_id", "to_account_id",
			"from_wallet_id", "to_wallet_id", "transaction_id").
		Values(exchange.ID, exchange.UserID, exchange.Type, exchange.Status,
			exchange.FromCurrency, exchange.ToCurrency, exchange.FromAmount, exchange.ToAmount,
			exchange.ExchangeRate, exchange.FromAccountID, exchange.ToAccountID,
			exchange.FromWalletID, exchange.ToWalletID, exchange.TransactionID).
		Suffix("RETURNING created_at, updated_at")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(&exchange.CreatedAt, &exchange.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create exchange: %w", err)
	}

	return nil
}

// GetByID retrieves an exchange by ID
func (r *ExchangeRepository) GetByID(id uuid.UUID) (*models.Exchange, error) {
	var exchange models.Exchange

	query := r.qb.Select("id", "user_id", "type", "status", "from_currency", "to_currency",
		"from_amount", "to_amount", "exchange_rate", "from_account_id", "to_account_id",
		"from_wallet_id", "to_wallet_id", "transaction_id", "created_at", "updated_at").
		From("exchanges").
		Where(sq.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(
		&exchange.ID, &exchange.UserID, &exchange.Type, &exchange.Status,
		&exchange.FromCurrency, &exchange.ToCurrency, &exchange.FromAmount, &exchange.ToAmount,
		&exchange.ExchangeRate, &exchange.FromAccountID, &exchange.ToAccountID,
		&exchange.FromWalletID, &exchange.ToWalletID, &exchange.TransactionID,
		&exchange.CreatedAt, &exchange.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("exchange not found")
		}
		return nil, fmt.Errorf("failed to get exchange: %w", err)
	}

	return &exchange, nil
}

// GetByUserID retrieves all exchanges for a user
func (r *ExchangeRepository) GetByUserID(userID uuid.UUID) ([]*models.Exchange, error) {
	query := r.qb.Select("id", "user_id", "type", "status", "from_currency", "to_currency",
		"from_amount", "to_amount", "exchange_rate", "from_account_id", "to_account_id",
		"from_wallet_id", "to_wallet_id", "transaction_id", "created_at", "updated_at").
		From("exchanges").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchanges: %w", err)
	}
	defer rows.Close()

	var exchanges []*models.Exchange
	for rows.Next() {
		var exchange models.Exchange
		err := rows.Scan(
			&exchange.ID, &exchange.UserID, &exchange.Type, &exchange.Status,
			&exchange.FromCurrency, &exchange.ToCurrency, &exchange.FromAmount, &exchange.ToAmount,
			&exchange.ExchangeRate, &exchange.FromAccountID, &exchange.ToAccountID,
			&exchange.FromWalletID, &exchange.ToWalletID, &exchange.TransactionID,
			&exchange.CreatedAt, &exchange.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exchange: %w", err)
		}
		exchanges = append(exchanges, &exchange)
	}

	return exchanges, nil
}

// UpdateStatus updates exchange status
func (r *ExchangeRepository) UpdateStatus(id uuid.UUID, status models.ExchangeStatus) error {
	query := r.qb.Update("exchanges").
		Set("status", status).
		Where(sq.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("exchange not found")
	}

	return nil
}

// GetExchangeRate retrieves exchange rate between two currencies
func (r *ExchangeRepository) GetExchangeRate(fromCurrency, toCurrency string) (float64, error) {
	var rate float64

	query := r.qb.Select("rate").
		From("exchange_rates").
		Where(sq.Eq{"from_currency": fromCurrency, "to_currency": toCurrency})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(&rate)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("exchange rate not found")
		}
		return 0, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	return rate, nil
}
