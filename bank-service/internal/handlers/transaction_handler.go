package handlers

import (
	"github.com/crypto-bank/bank-service/internal/models"
	"github.com/crypto-bank/bank-service/internal/services"
	"github.com/crypto-bank/bank-service/pkg/metrics"
	"github.com/crypto-bank/bank-service/pkg/response"
	"github.com/crypto-bank/bank-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// CreateTransfer godoc
// @Summary Create a transfer between accounts
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.CreateTransactionRequest true "Transaction data"
// @Success 201 {object} response.Response{data=models.Transaction}
// @Router /api/v1/transactions/transfer [post]
func (h *TransactionHandler) CreateTransfer(c *fiber.Ctx) error {
	var req models.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := validator.Validate(&req); err != nil {
		return response.BadRequest(c, "Validation failed", err)
	}

	transaction, err := h.transactionService.CreateTransfer(&req)
	if err != nil {
		metrics.TransactionsTotal.WithLabelValues("transfer", "failed").Inc()
		return response.InternalServerError(c, "Failed to create transfer", err)
	}

	metrics.TransactionsTotal.WithLabelValues("transfer", "success").Inc()
	return response.Created(c, transaction, "Transfer created successfully")
}

// Deposit godoc
// @Summary Deposit money to an account
// @Tags transactions
// @Accept json
// @Produce json
// @Param deposit body models.DepositRequest true "Deposit data"
// @Success 201 {object} response.Response{data=models.Transaction}
// @Router /api/v1/transactions/deposit [post]
func (h *TransactionHandler) Deposit(c *fiber.Ctx) error {
	var req models.DepositRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := validator.Validate(&req); err != nil {
		return response.BadRequest(c, "Validation failed", err)
	}

	transaction, err := h.transactionService.Deposit(&req)
	if err != nil {
		metrics.TransactionsTotal.WithLabelValues("deposit", "failed").Inc()
		return response.InternalServerError(c, "Failed to deposit", err)
	}

	metrics.TransactionsTotal.WithLabelValues("deposit", "success").Inc()
	return response.Created(c, transaction, "Deposit successful")
}

// Withdraw godoc
// @Summary Withdraw money from an account
// @Tags transactions
// @Accept json
// @Produce json
// @Param withdraw body models.WithdrawRequest true "Withdraw data"
// @Success 201 {object} response.Response{data=models.Transaction}
// @Router /api/v1/transactions/withdraw [post]
func (h *TransactionHandler) Withdraw(c *fiber.Ctx) error {
	var req models.WithdrawRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := validator.Validate(&req); err != nil {
		return response.BadRequest(c, "Validation failed", err)
	}

	transaction, err := h.transactionService.Withdraw(&req)
	if err != nil {
		metrics.TransactionsTotal.WithLabelValues("withdraw", "failed").Inc()
		return response.InternalServerError(c, "Failed to withdraw", err)
	}

	metrics.TransactionsTotal.WithLabelValues("withdraw", "success").Inc()
	return response.Created(c, transaction, "Withdrawal successful")
}

// GetTransaction godoc
// @Summary Get transaction by ID
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} response.Response{data=models.Transaction}
// @Router /api/v1/transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid transaction ID", err)
	}

	transaction, err := h.transactionService.GetTransaction(id)
	if err != nil {
		return response.NotFound(c, "Transaction not found")
	}

	return response.Success(c, transaction, "")
}

// GetUserTransactions godoc
// @Summary Get all transactions for a user
// @Tags transactions
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} response.Response{data=[]models.Transaction}
// @Router /api/v1/users/{user_id}/transactions [get]
func (h *TransactionHandler) GetUserTransactions(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	transactions, err := h.transactionService.GetUserTransactions(userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get transactions", err)
	}

	return response.Success(c, transactions, "")
}
