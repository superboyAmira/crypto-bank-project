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

type ExchangeHandler struct {
	exchangeService *services.ExchangeService
}

func NewExchangeHandler(exchangeService *services.ExchangeService) *ExchangeHandler {
	return &ExchangeHandler{
		exchangeService: exchangeService,
	}
}

// ExchangeCryptoToFiat godoc
// @Summary Exchange cryptocurrency to fiat currency
// @Tags exchanges
// @Accept json
// @Produce json
// @Param exchange body models.ExchangeCryptoToFiatRequest true "Exchange data"
// @Success 201 {object} response.Response{data=models.Exchange}
// @Router /api/v1/exchanges/crypto-to-fiat [post]
func (h *ExchangeHandler) ExchangeCryptoToFiat(c *fiber.Ctx) error {
	var req models.ExchangeCryptoToFiatRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := validator.Validate(&req); err != nil {
		return response.BadRequest(c, "Validation failed", err)
	}

	exchange, err := h.exchangeService.ExchangeCryptoToFiat(&req)
	if err != nil {
		metrics.ExchangesTotal.WithLabelValues("crypto_to_fiat", "failed").Inc()
		return response.InternalServerError(c, "Failed to exchange crypto to fiat", err)
	}

	metrics.ExchangesTotal.WithLabelValues("crypto_to_fiat", "success").Inc()
	return response.Created(c, exchange, "Exchange completed successfully")
}

// ExchangeFiatToCrypto godoc
// @Summary Exchange fiat currency to cryptocurrency
// @Tags exchanges
// @Accept json
// @Produce json
// @Param exchange body models.ExchangeFiatToCryptoRequest true "Exchange data"
// @Success 201 {object} response.Response{data=models.Exchange}
// @Router /api/v1/exchanges/fiat-to-crypto [post]
func (h *ExchangeHandler) ExchangeFiatToCrypto(c *fiber.Ctx) error {
	var req models.ExchangeFiatToCryptoRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := validator.Validate(&req); err != nil {
		return response.BadRequest(c, "Validation failed", err)
	}

	exchange, err := h.exchangeService.ExchangeFiatToCrypto(&req)
	if err != nil {
		metrics.ExchangesTotal.WithLabelValues("fiat_to_crypto", "failed").Inc()
		return response.InternalServerError(c, "Failed to exchange fiat to crypto", err)
	}

	metrics.ExchangesTotal.WithLabelValues("fiat_to_crypto", "success").Inc()
	return response.Created(c, exchange, "Exchange completed successfully")
}

// GetExchange godoc
// @Summary Get exchange by ID
// @Tags exchanges
// @Produce json
// @Param id path string true "Exchange ID"
// @Success 200 {object} response.Response{data=models.Exchange}
// @Router /api/v1/exchanges/{id} [get]
func (h *ExchangeHandler) GetExchange(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid exchange ID", err)
	}

	exchange, err := h.exchangeService.GetExchange(id)
	if err != nil {
		return response.NotFound(c, "Exchange not found")
	}

	return response.Success(c, exchange, "")
}

// GetUserExchanges godoc
// @Summary Get all exchanges for a user
// @Tags exchanges
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} response.Response{data=[]models.Exchange}
// @Router /api/v1/users/{user_id}/exchanges [get]
func (h *ExchangeHandler) GetUserExchanges(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	exchanges, err := h.exchangeService.GetUserExchanges(userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get exchanges", err)
	}

	return response.Success(c, exchanges, "")
}
