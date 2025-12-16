package repositories

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	user.ID = uuid.New()

	query := r.qb.Insert("users").
		Columns("id", "email", "first_name", "last_name", "phone").
		Values(user.ID, user.Email, user.FirstName, user.LastName, user.Phone).
		Suffix("RETURNING created_at, updated_at")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User

	query := r.qb.Select("id", "email", "first_name", "last_name", "phone", "created_at", "updated_at", "deleted_at").
		From("users").
		Where(sq.Eq{"id": id, "deleted_at": nil})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	query := r.qb.Select("id", "email", "first_name", "last_name", "phone", "created_at", "updated_at", "deleted_at").
		From("users").
		Where(sq.Eq{"email": email, "deleted_at": nil})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetAll retrieves all users
func (r *UserRepository) GetAll() ([]*models.User, error) {
	query := r.qb.Select("id", "email", "first_name", "last_name", "phone", "created_at", "updated_at", "deleted_at").
		From("users").
		Where(sq.Eq{"deleted_at": nil}).
		OrderBy("created_at DESC")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
			&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	return users, nil
}

// Update updates a user
func (r *UserRepository) Update(id uuid.UUID, req *models.UpdateUserRequest) error {
	updateMap := make(map[string]interface{})

	if req.FirstName != "" {
		updateMap["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updateMap["last_name"] = req.LastName
	}
	if req.Phone != "" {
		updateMap["phone"] = req.Phone
	}

	if len(updateMap) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := r.qb.Update("users").
		SetMap(updateMap).
		Where(sq.Eq{"id": id, "deleted_at": nil})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	query := r.qb.Update("users").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id, "deleted_at": nil})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
