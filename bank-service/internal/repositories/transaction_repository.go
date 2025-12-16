package repositories

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/google/uuid"
)

type TransactionRepository struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(tx *models.Transaction) error {
	tx.ID = uuid.New()

	query := r.qb.Insert("transactions").
		Columns("id", "user_id", "type", "status", "amount", "currency",
			"from_account_id", "to_account_id", "description", "exchange_id").
		Values(tx.ID, tx.UserID, tx.Type, tx.Status, tx.Amount, tx.Currency,
			tx.FromAccountID, tx.ToAccountID, tx.Description, tx.ExchangeID).
		Suffix("RETURNING created_at, updated_at")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(&tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	var tx models.Transaction

	query := r.qb.Select("id", "user_id", "type", "status", "amount", "currency",
		"from_account_id", "to_account_id", "description", "exchange_id", "created_at", "updated_at").
		From("transactions").
		Where(sq.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(
		&tx.ID, &tx.UserID, &tx.Type, &tx.Status, &tx.Amount, &tx.Currency,
		&tx.FromAccountID, &tx.ToAccountID, &tx.Description, &tx.ExchangeID,
		&tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &tx, nil
}

// GetByUserID retrieves all transactions for a user
func (r *TransactionRepository) GetByUserID(userID uuid.UUID) ([]*models.Transaction, error) {
	query := r.qb.Select("id", "user_id", "type", "status", "amount", "currency",
		"from_account_id", "to_account_id", "description", "exchange_id", "created_at", "updated_at").
		From("transactions").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		var tx models.Transaction
		err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.Type, &tx.Status, &tx.Amount, &tx.Currency,
			&tx.FromAccountID, &tx.ToAccountID, &tx.Description, &tx.ExchangeID,
			&tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, &tx)
	}

	return transactions, nil
}

// UpdateStatus updates transaction status
func (r *TransactionRepository) UpdateStatus(id uuid.UUID, status models.TransactionStatus) error {
	query := r.qb.Update("transactions").
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
		return fmt.Errorf("transaction not found")
	}

	return nil
}
