package services

import (
	"fmt"

	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/crypto-bank/bank-service/internal/repositories"
	"github.com/crypto-bank/bank-service/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	logger.Info("Creating user", zap.String("email", req.Email))

	// Check if user with email already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	logger.Info("User created", zap.String("user_id", user.ID.String()))
	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.GetAll()
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id uuid.UUID, req *models.UpdateUserRequest) error {
	logger.Info("Updating user", zap.String("user_id", id.String()))

	if err := s.userRepo.Update(id, req); err != nil {
		logger.Error("Failed to update user", zap.Error(err))
		return err
	}

	logger.Info("User updated", zap.String("user_id", id.String()))
	return nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id uuid.UUID) error {
	logger.Info("Deleting user", zap.String("user_id", id.String()))

	if err := s.userRepo.Delete(id); err != nil {
		logger.Error("Failed to delete user", zap.Error(err))
		return err
	}

	logger.Info("User deleted", zap.String("user_id", id.String()))
	return nil
}

