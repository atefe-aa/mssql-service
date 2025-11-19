package database

import (
	"fmt"
	"mssql-api/internal/config"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func NewGormDatabase(cfg *config.Config) (*gorm.DB, error) {
	// Build connection string
	dsn := buildConnectionString(cfg)

	// Open connection using GORM
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Optional: ping the database to verify connectivity
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	return db, nil
}

func buildConnectionString(cfg *config.Config) string {
	if cfg.DBAuthMode == "windows" {
		// Windows Authentication
		return fmt.Sprintf("sqlserver://%s:%s?database=%s&connection+timeout=30",
			cfg.DBServer,
			cfg.DBPort,
			cfg.DBName,
		)
	} else {
		// SQL Authentication
		return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&connection+timeout=30",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBServer,
			cfg.DBPort,
			cfg.DBName,
		)
	}
}
