package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sefikcan/address-api/internal/ratelimiter"
)

func (mw *Manager) EPDistributedRateLimitMiddleware(dtb *ratelimiter.EPDistributedTokenBucket, endpoint string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userKey := c.IP() // Kullanıcı tanımlayıcı olarak başka bir şey de kullanılabilir
		if !dtb.EPAllowRequest(c.Context(), userKey, endpoint, 1) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded for " + endpoint,
			})
		}
		return c.Next()
	}
}
