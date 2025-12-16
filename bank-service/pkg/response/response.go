package response

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a successful response
func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a created response
func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *fiber.Ctx, statusCode int, message string, err error) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

// BadRequest sends a bad request error
func BadRequest(c *fiber.Ctx, message string, err error) error {
	return Error(c, fiber.StatusBadRequest, message, err)
}

// NotFound sends a not found error
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, nil)
}

// InternalServerError sends an internal server error
func InternalServerError(c *fiber.Ctx, message string, err error) error {
	return Error(c, fiber.StatusInternalServerError, message, err)
}

// Unauthorized sends an unauthorized error
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message, nil)
}

