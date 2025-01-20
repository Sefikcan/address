package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sefikcan/address-api/internal/ratelimiter"
)

func (mw *Manager) RateLimitMiddleware(tb *ratelimiter.TokenBucket) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// check 1 token in bucket
		if !tb.AllowRequest(1) {
			return ctx.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		}
		return ctx.Next()
	}
}
