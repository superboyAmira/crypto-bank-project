package handlers

import (
	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/crypto-bank/bank-service/internal/services"
	"github.com/crypto-bank/bank-service/pkg/response"
	"github.com/crypto-bank/bank-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CryptoWalletHandler struct {
	walletService *services.CryptoWalletService
}

func NewCryptoWalletHandler(walletService *services.CryptoWalletService) *CryptoWalletHandler {
	return &CryptoWalletHandler{
		walletService: walletService,
	}
}

// CreateWallet godoc
// @Summary Create a new crypto wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param wallet body models.CreateCryptoWalletRequest true "Wallet data"
// @Success 201 {object} response.Response{data=models.CryptoWallet}
// @Router /api/v1/wallets [post]
func (h *CryptoWalletHandler) CreateWallet(c *fiber.Ctx) error {
	var req models.CreateCryptoWalletRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := validator.Validate(&req); err != nil {
		return response.BadRequest(c, "Validation failed", err)
	}

	wallet, err := h.walletService.CreateWallet(&req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create wallet", err)
	}

	return response.Created(c, wallet, "Wallet created successfully")
}

// GetWallet godoc
// @Summary Get wallet by ID
// @Tags wallets
// @Produce json
// @Param id path string true "Wallet ID"
// @Success 200 {object} response.Response{data=models.CryptoWallet}
// @Router /api/v1/wallets/{id} [get]
func (h *CryptoWalletHandler) GetWallet(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid wallet ID", err)
	}

	wallet, err := h.walletService.GetWallet(id)
	if err != nil {
		return response.NotFound(c, "Wallet not found")
	}

	return response.Success(c, wallet, "")
}

// GetUserWallets godoc
// @Summary Get all wallets for a user
// @Tags wallets
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} response.Response{data=[]models.CryptoWallet}
// @Router /api/v1/users/{user_id}/wallets [get]
func (h *CryptoWalletHandler) GetUserWallets(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	wallets, err := h.walletService.GetUserWallets(userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get wallets", err)
	}

	return response.Success(c, wallets, "")
}

// GetWalletBalance godoc
// @Summary Get wallet balance
// @Tags wallets
// @Produce json
// @Param id path string true "Wallet ID"
// @Success 200 {object} response.Response{data=float64}
// @Router /api/v1/wallets/{id}/balance [get]
func (h *CryptoWalletHandler) GetWalletBalance(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid wallet ID", err)
	}

	balance, err := h.walletService.GetWalletBalance(id)
	if err != nil {
		return response.NotFound(c, "Wallet not found")
	}

	return response.Success(c, map[string]interface{}{"balance": balance}, "")
}

