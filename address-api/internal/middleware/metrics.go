package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sefikcan/address/address-api/pkg/metric"
	"time"
)

func (mw *Manager) Metrics(metrics metric.Metrics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		status := c.Response().StatusCode()

		// If there’s an error, use Fiber's error handling to get the status code
		if err != nil {
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				status = fiberErr.Code
			}
		}

		metrics.ObserveResponseTime(status, c.Method(), c.Path(), time.Since(start).Seconds())
		metrics.IncreaseHits(status, c.Method(), c.Path())

		return err
	}
}
