package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"time"
)

func (mw Manager) ErrorLogger(c *fiber.Ctx) error {
	start := time.Now()

	err := c.Next()

	method := c.Method()
	path := c.Path()
	duration := time.Since(start)

	if err != nil {
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			mw.logger.Errorf("Error: %v, Path: %s, Method: %s, Duration: %", fiberErr.Message, path, method, fiberErr.Code, duration)
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"error": fiberErr.Message,
			})
		}

		mw.logger.Errorf("Error: %v, Path: %s, Method: %s, Status: %d, Duration: %s", err.Error(), path, method, fiber.StatusInternalServerError, duration)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	return nil
}
