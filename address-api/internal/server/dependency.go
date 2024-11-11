package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sefikcan/address/address-api/internal/address/handlers"
	"github.com/sefikcan/address/address-api/internal/address/repository"
	"github.com/sefikcan/address/address-api/internal/address/service"
	mw "github.com/sefikcan/address/address-api/internal/middleware"
	"github.com/sefikcan/address/address-api/pkg/metric"
	"github.com/sefikcan/address/address-api/pkg/swagger"
)

func (s *Server) MapHandlers(app *fiber.App) error {
	metrics, err := metric.CreateMetrics(s.cfg.Metric.Url, s.cfg.Metric.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics error: %s", err)
	}

	// initialize repositories and service
	addressRepository := repository.NewAddressRepository(s.db)
	addressService := service.NewAddressService(s.cfg, addressRepository, s.logger)

	// initialize handler
	addressHandler := handlers.NewAddressHandler(addressService)

	middlewareManager := mw.NewMiddlewareManager(s.cfg, s.logger)

	// initialize handler
	handlers.MapAddressRotes(app, addressHandler, s.logger, middlewareManager)

	// set up middleware
	app.Use(middlewareManager.RequestLogger)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, X-Request-ID",
	}))
	app.Use(recover.New(recover.Config{
		EnableStackTrace: false,
	}))
	app.Use(requestid.New())
	app.Use(middlewareManager.Metrics(metrics))

	swaggerMiddleware := swagger.NewSwagger("./address-api/docs/swagger.json", "/", s.logger)
	swaggerMiddleware.RegisterSwagger(app)

	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Content-Security-Policy", "default-src 'self'")
		return c.Next()
	})
	app.Use(func(c *fiber.Ctx) error {
		if c.Request().Header.ContentLength() > (2 * 1024 * 1024) { // 2MB limit
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{"error": "Request body too large"})
		}
		return c.Next()
	})

	return nil
}
