package main

import (
	"database/sql"
	"github.com/sefikcan/address/address-api/internal/server"
	"github.com/sefikcan/address/address-api/pkg/config"
	"github.com/sefikcan/address/address-api/pkg/logger"
	"github.com/sefikcan/address/address-api/pkg/storage/postgres"
	"log"
)

// @title Address API
// @version 1.0
// @description This is an Address API for Swagger documentation.
// @host localhost:3048
// @BasePath /
func main() {
	log.Println("Starting api server")

	cfg := config.NewConfig()

	zapLogger := logger.NewLogger(cfg)
	zapLogger.InitLogger()
	zapLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, false)

	psqlDb, err := postgres.NewPsqlDb(cfg)
	db, err := psqlDb.DB()
	if err != nil {
		zapLogger.Fatalf("Postgresql init: %s", err)
	} else {
		zapLogger.Infof("Postgres connected, Status: %#v", db.Stats())
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			zapLogger.Fatal(err)
		}
	}(db)

	// start server
	s := server.NewServer(cfg, psqlDb, zapLogger)
	if err = s.Run(); err != nil {
		zapLogger.Fatal(err)
	}
}
