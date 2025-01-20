package middleware

import (
	"github.com/gofiber/fiber/v2"
	idempotency "github.com/sefikcan/address-api/internal/idempotency/service"
)

func (mw *Manager) IdempotencyMiddleware(idempotencyService idempotency.IdempotencyService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		idempotencyKey := ctx.Get("idempotent-Key")
		if idempotencyKey == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "idempotent-Key header is required",
			})
		}

		_, err := idempotencyService.GetByKey(ctx.Context(), idempotencyKey)
		err = ctx.Next()
		if err == nil {
			response := ctx.Response().Body()
			if err := idempotencyService.Save(ctx.Context(), idempotencyKey, response); err != nil {
				mw.logger.Info("Error saving idempotency result: ", err)
			}
		}

		return err
	}
}
