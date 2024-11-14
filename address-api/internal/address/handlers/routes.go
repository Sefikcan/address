package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sefikcan/address-api/internal/address/dto/request"
	"github.com/sefikcan/address-api/internal/middleware"
	"github.com/sefikcan/address-api/pkg/logger"
	"github.com/sefikcan/address-api/pkg/util"
)

func MapAddressRotes(app *fiber.App, addressHandler AddressHandler, logger logger.Logger, manager *middleware.Manager) {
	v1 := app.Group("/api/v1/addresses")

	v1.Get("/", addressHandler.GetAll)
	v1.Post("/", middleware.Validator(&request.AddressCreateRequest{}), addressHandler.Create)
	v1.Delete("/:id", addressHandler.Delete)
	v1.Get("/:id", addressHandler.GetById)
	v1.Put("/:id", middleware.Validator(&request.AddressUpdateRequest{}), addressHandler.Update)
	v1.Patch("/:id", middleware.Validator(&request.AddressPatchRequest{}), addressHandler.Patch)

	// health endpoint
	health := v1.Group("/health")
	// Health check endpoint
	health.Get("/", func(c *fiber.Ctx) error {
		logger.Infof("Health check RequestID: %s", util.GetRequestId(c))
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "OK"})
	})

	// Version 2 routes sample
	v2 := app.Group("/api/v2")
	v2.Get("/", addressHandler.GetAllV2)
}
