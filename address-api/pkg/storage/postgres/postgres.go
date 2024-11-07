package postgres

import (
	"fmt"
	"github.com/sefikcan/address/address-api/internal/address/entity"
	"github.com/sefikcan/address/address-api/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// NewPsqlDb function connect postgresql database
// The gorm.DB object is returned according to the connection settings in the config.
func NewPsqlDb(c *config.Config) (*gorm.DB, error) {
	// create postgresql connection string
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.Postgres.UserName,
		c.Postgres.Password,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.DbName)

	// open postgresql connection
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// we reach *sql.DB object for set database configuration settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// database connection pool settings
	sqlDB.SetMaxOpenConns(c.Postgres.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.Postgres.ConnMaxLifeTime) * time.Second)
	sqlDB.SetMaxIdleConns(c.Postgres.MaxIdleConns)
	sqlDB.SetConnMaxIdleTime(time.Duration(c.Postgres.ConnMaxIdleTime) * time.Second)
	// check database connection
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	// Automatically migrate schema for all specified models
	if err := db.AutoMigrate(
		&entity.Address{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
