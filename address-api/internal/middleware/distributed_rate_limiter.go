package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sefikcan/address-api/internal/ratelimiter"
)

func (mw *Manager) DistributedRateLimitMiddleware(dtb *ratelimiter.DistributedTokenBucket) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// get client ip address
		clientIP := ctx.IP()
		key := clientIP
		// check if we can request a token from distributed token bucket structure
		if !dtb.AllowRequest(ctx.Context(), key, 1) {
			// if token not found, return HTTP 429(Too many requests)
			return ctx.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		}
		// if token exist, we move on the next middleware
		return ctx.Next()
	}
}
