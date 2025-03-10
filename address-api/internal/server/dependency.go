package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sefikcan/address-api/internal/address/handlers"
	"github.com/sefikcan/address-api/internal/address/repository"
	"github.com/sefikcan/address-api/internal/address/service"
	idempotency "github.com/sefikcan/address-api/internal/idempotency/service"
	mw "github.com/sefikcan/address-api/internal/middleware"
	"github.com/sefikcan/address-api/internal/ratelimiter"
	"github.com/sefikcan/address-api/pkg/kafka"
	"github.com/sefikcan/address-api/pkg/metric"
	"github.com/sefikcan/address-api/pkg/redis"
)

func (s *Server) MapHandlers(app *fiber.App) error {
	metrics, err := metric.CreateMetrics(s.cfg.Metric.Url, s.cfg.Metric.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics error: %s", err)
	}

	if err = redis.InitializeRedis(s.cfg); err != nil {
		s.logger.Fatalf("Failed to connect to Redis: %v", err)
	}
	redisClient := redis.GetClient()
	_, err = redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		s.logger.Fatalf("Failed to ping Redis: %v", err)
		return err
	}

	kafkaProducer := kafka.NewKafkaProducer(s.cfg.Kafka.Brokers)
	if err != nil {
		s.logger.Errorf("Error setting up Kafka producers: %v", err)
		return err
	}

	// initialize repositories and service
	addressRepository := repository.NewAddressRepository(s.db)
	addressService := service.NewAddressService(s.cfg, addressRepository, s.logger, kafkaProducer)

	middlewareManager := mw.NewMiddlewareManager(s.cfg, s.logger)

	// set up middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, X-Request-ID",
	}))
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(requestid.New())
	//tb := ratelimiter.NewTokenBucket()
	//app.Use(middlewareManager.RateLimitMiddleware(tb))

	idempotencyService := idempotency.NewIdempotencyService(redisClient)
	app.Use(middlewareManager.IdempotencyMiddleware(idempotencyService))

	distributedTb := ratelimiter.NewDistributedTokenBucket(redisClient)
	app.Use(middlewareManager.DistributedRateLimitMiddleware(distributedTb))
	app.Use(middlewareManager.RequestLogger)
	app.Use(middlewareManager.ErrorLogger)
	app.Use(middlewareManager.Metrics(metrics))

	//swaggerMiddleware := swagger.NewSwagger("./address-api/docs/swagger.json", "/", s.logger)
	//swaggerMiddleware.RegisterSwagger(app)

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

	// initialize handler
	addressHandler := handlers.NewAddressHandler(addressService)

	// initialize handler
	handlers.MapAddressRotes(app, addressHandler, s.logger, middlewareManager, idempotencyService)

	return nil
}
