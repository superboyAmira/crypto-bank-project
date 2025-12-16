package repositories

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/google/uuid"
)

type AccountRepository struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create creates a new account
func (r *AccountRepository) Create(account *models.Account) error {
	account.ID = uuid.New()

	query := r.qb.Insert("accounts").
		Columns("id", "user_id", "currency", "balance").
		Values(account.ID, account.UserID, account.Currency, account.Balance).
		Suffix("RETURNING created_at, updated_at")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(&account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

// GetByID retrieves an account by ID
func (r *AccountRepository) GetByID(id uuid.UUID) (*models.Account, error) {
	var account models.Account

	query := r.qb.Select("id", "user_id", "currency", "balance", "created_at", "updated_at").
		From("accounts").
		Where(sq.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(
		&account.ID, &account.UserID, &account.Currency, &account.Balance,
		&account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

// GetByUserID retrieves all accounts for a user
func (r *AccountRepository) GetByUserID(userID uuid.UUID) ([]*models.Account, error) {
	query := r.qb.Select("id", "user_id", "currency", "balance", "created_at", "updated_at").
		From("accounts").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*models.Account
	for rows.Next() {
		var account models.Account
		err := rows.Scan(
			&account.ID, &account.UserID, &account.Currency, &account.Balance,
			&account.CreatedAt, &account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

// UpdateBalance updates account balance
func (r *AccountRepository) UpdateBalance(id uuid.UUID, amount float64) error {
	query := r.qb.Update("accounts").
		Set("balance", sq.Expr("balance + ?", amount)).
		Where(sq.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

// GetBalance retrieves account balance
func (r *AccountRepository) GetBalance(id uuid.UUID) (float64, error) {
	var balance float64

	query := r.qb.Select("balance").
		From("accounts").
		Where(sq.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("account not found")
		}
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}
