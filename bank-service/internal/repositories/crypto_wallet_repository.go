package repositories

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/google/uuid"
)

type CryptoWalletRepository struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewCryptoWalletRepository(db *sql.DB) *CryptoWalletRepository {
	return &CryptoWalletRepository{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create creates a new crypto wallet
func (r *CryptoWalletRepository) Create(wallet *models.CryptoWallet) error {
	wallet.ID = uuid.New()

	query := r.qb.Insert("crypto_wallets").
		Columns("id", "user_id", "crypto_type", "balance", "address").
		Values(wallet.ID, wallet.UserID, wallet.CryptoType, wallet.Balance, wallet.Address).
		Suffix("RETURNING created_at, updated_at")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(&wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create crypto wallet: %w", err)
	}

	return nil
}

// GetByID retrieves a crypto wallet by ID
func (r *CryptoWalletRepository) GetByID(id uuid.UUID) (*models.CryptoWallet, error) {
	var wallet models.CryptoWallet

	query := r.qb.Select("id", "user_id", "crypto_type", "balance", "address", "created_at", "updated_at").
		From("crypto_wallets").
		Where(sq.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(
		&wallet.ID, &wallet.UserID, &wallet.CryptoType, &wallet.Balance, &wallet.Address,
		&wallet.CreatedAt, &wallet.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("crypto wallet not found")
		}
		return nil, fmt.Errorf("failed to get crypto wallet: %w", err)
	}

	return &wallet, nil
}

// GetByUserID retrieves all crypto wallets for a user
func (r *CryptoWalletRepository) GetByUserID(userID uuid.UUID) ([]*models.CryptoWallet, error) {
	query := r.qb.Select("id", "user_id", "crypto_type", "balance", "address", "created_at", "updated_at").
		From("crypto_wallets").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get crypto wallets: %w", err)
	}
	defer rows.Close()

	var wallets []*models.CryptoWallet
	for rows.Next() {
		var wallet models.CryptoWallet
		err := rows.Scan(
			&wallet.ID, &wallet.UserID, &wallet.CryptoType, &wallet.Balance, &wallet.Address,
			&wallet.CreatedAt, &wallet.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan crypto wallet: %w", err)
		}
		wallets = append(wallets, &wallet)
	}

	return wallets, nil
}

// UpdateBalance updates crypto wallet balance
func (r *CryptoWalletRepository) UpdateBalance(id uuid.UUID, amount float64) error {
	query := r.qb.Update("crypto_wallets").
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
		return fmt.Errorf("crypto wallet not found")
	}

	return nil
}

// GetBalance retrieves crypto wallet balance
func (r *CryptoWalletRepository) GetBalance(id uuid.UUID) (float64, error) {
	var balance float64

	query := r.qb.Select("balance").
		From("crypto_wallets").
		Where(sq.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("crypto wallet not found")
		}
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}
