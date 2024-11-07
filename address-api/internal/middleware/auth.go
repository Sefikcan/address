package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"strings"
)

// Authentication to protect routes
func (mw Manager) Authentication() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Get the token from the Authorization header
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			mw.logger.Warn("Missing or malformed token")

			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed token",
			})
		}

		// Check if the token format is `Bearer <token>`
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			mw.logger.Info("Invalid token format")

			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}

		tokenStr := parts[1]
		// Parse the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				mw.logger.Error("unexpected signing method")

				return nil, errors.New("unexpected signing method")
			}
			// Return the signing key
			return []byte(mw.cfg.Auth.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			mw.logger.Info("Invalid or expired token")

			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Token is valid; set user info in context if needed
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			mw.logger.Info("Invalid or expired token")

			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Failed to parse token claims",
			})
		}

		// You can set the user ID in the context to access it in handlers
		ctx.Locals("userID", claims["user_id"])

		return ctx.Next()
	}
}
