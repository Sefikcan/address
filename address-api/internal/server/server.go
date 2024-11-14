package server

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sefikcan/address-api/pkg/config"
	"github.com/sefikcan/address-api/pkg/logger"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	app    *fiber.App
	cfg    *config.Config
	db     *gorm.DB
	logger logger.Logger
}

func (s *Server) Run() error {
	// set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// map handlers before starting server
	if err := s.MapHandlers(s.app); err != nil {
		s.logger.Fatalf("Error map handler: %v", err)
		return err
	}

	go func() {
		address := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.app.Listen(address); err != nil {
			s.logger.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for an interrupt signal for graceful shutdown
	<-quit
	s.logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.CtxTimeout*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		s.logger.Errorf("Server shutdown error: %v", err)
		return err
	}

	s.logger.Info("Server exited properly")
	return nil
}

func NewServer(cfg *config.Config, db *gorm.DB, logger logger.Logger) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * cfg.Server.ReadTimeout,
		WriteTimeout: time.Second * cfg.Server.WriteTimeout,
	})
	return &Server{
		app:    app,
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}
