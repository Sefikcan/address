package middleware

import (
	"github.com/sefikcan/address-api/pkg/config"
	"github.com/sefikcan/address-api/pkg/logger"
)

type Manager struct {
	cfg    *config.Config
	logger logger.Logger
}

func NewMiddlewareManager(cfg *config.Config, logger logger.Logger) *Manager {
	return &Manager{
		cfg:    cfg,
		logger: logger,
	}
}
