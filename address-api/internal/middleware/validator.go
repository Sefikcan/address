package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// Validator is a middleware to validate the request body against the provided struct
func Validator(model interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(model); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		validate := validator.New()
		if err := validate.Struct(model); err != nil {
			var validationErrors []string
			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors,
					"Field: "+err.Field()+" failed on the '"+err.Tag()+"' tag")
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": strings.Join(validationErrors, ", "),
			})
		}

		return c.Next()
	}
}
