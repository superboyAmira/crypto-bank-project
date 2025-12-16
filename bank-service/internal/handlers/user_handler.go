package handlers

import (
	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/crypto-bank/bank-service/internal/services"
	"github.com/crypto-bank/bank-service/pkg/response"
	"github.com/crypto-bank/bank-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User data"
// @Success 201 {object} response.Response{data=models.User}
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := validator.Validate(&req); err != nil {
		return response.BadRequest(c, "Validation failed", err)
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create user", err)
	}

	return response.Created(c, user, "User created successfully")
}

// GetUser godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=models.User}
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	return response.Success(c, user, "")
}

// GetAllUsers godoc
// @Summary Get all users
// @Tags users
// @Produce json
// @Success 200 {object} response.Response{data=[]models.User}
// @Router /api/v1/users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return response.InternalServerError(c, "Failed to get users", err)
	}

	return response.Success(c, users, "")
}

// UpdateUser godoc
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUserRequest true "User data"
// @Success 200 {object} response.Response
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := h.userService.UpdateUser(id, &req); err != nil {
		return response.InternalServerError(c, "Failed to update user", err)
	}

	return response.Success(c, nil, "User updated successfully")
}

// DeleteUser godoc
// @Summary Delete user
// @Tags users
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	if err := h.userService.DeleteUser(id); err != nil {
		return response.InternalServerError(c, "Failed to delete user", err)
	}

	return response.Success(c, nil, "User deleted successfully")
}

