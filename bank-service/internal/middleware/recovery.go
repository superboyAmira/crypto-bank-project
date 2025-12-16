package middleware

import (
	"github.com/crypto-bank/bank-service/pkg/logger"
	"github.com/crypto-bank/bank-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Recovery middleware recovers from panics
func Recovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic recovered",
					zap.Any("error", r),
					zap.String("path", c.Path()),
				)
				response.InternalServerError(c, "Internal server error", nil)
			}
		}()
		return c.Next()
	}
}

