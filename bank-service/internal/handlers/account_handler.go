package handlers

import (
	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/crypto-bank/bank-service/internal/services"
	"github.com/crypto-bank/bank-service/pkg/response"
	"github.com/crypto-bank/bank-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AccountHandler struct {
	accountService *services.AccountService
}

func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// CreateAccount godoc
// @Summary Create a new fiat account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body models.CreateAccountRequest true "Account data"
// @Success 201 {object} response.Response{data=models.Account}
// @Router /api/v1/accounts [post]
func (h *AccountHandler) CreateAccount(c *fiber.Ctx) error {
	var req models.CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := validator.Validate(&req); err != nil {
		return response.BadRequest(c, "Validation failed", err)
	}

	account, err := h.accountService.CreateAccount(&req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create account", err)
	}

	return response.Created(c, account, "Account created successfully")
}

// GetAccount godoc
// @Summary Get account by ID
// @Tags accounts
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} response.Response{data=models.Account}
// @Router /api/v1/accounts/{id} [get]
func (h *AccountHandler) GetAccount(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid account ID", err)
	}

	account, err := h.accountService.GetAccount(id)
	if err != nil {
		return response.NotFound(c, "Account not found")
	}

	return response.Success(c, account, "")
}

// GetUserAccounts godoc
// @Summary Get all accounts for a user
// @Tags accounts
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} response.Response{data=[]models.Account}
// @Router /api/v1/users/{user_id}/accounts [get]
func (h *AccountHandler) GetUserAccounts(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	accounts, err := h.accountService.GetUserAccounts(userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get accounts", err)
	}

	return response.Success(c, accounts, "")
}

// GetAccountBalance godoc
// @Summary Get account balance
// @Tags accounts
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} response.Response{data=float64}
// @Router /api/v1/accounts/{id}/balance [get]
func (h *AccountHandler) GetAccountBalance(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid account ID", err)
	}

	balance, err := h.accountService.GetAccountBalance(id)
	if err != nil {
		return response.NotFound(c, "Account not found")
	}

	return response.Success(c, map[string]interface{}{"balance": balance}, "")
}

