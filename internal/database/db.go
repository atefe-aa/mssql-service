package database

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/microsoft/go-mssqldb"
	"mssql-api/internal/config"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	connString := buildConnectionString(cfg)
	
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &Database{DB: db}, nil
}

func buildConnectionString(cfg *config.Config) string {
	query := url.Values{}
	query.Add("database", cfg.DBName)
	query.Add("connection timeout", "30")
	
	var connString string
	
	if cfg.DBAuthMode == "windows" {
		// Windows Authentication
		connString = fmt.Sprintf("sqlserver://%s:%s?%s",
			cfg.DBServer,
			cfg.DBPort,
			query.Encode())
	} else {
		// SQL Authentication
		u := &url.URL{
			Scheme: "sqlserver",
			User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
			Host:   fmt.Sprintf("%s:%s", cfg.DBServer, cfg.DBPort),
		}
		u.RawQuery = query.Encode()
		connString = u.String()
	}
	
	return connString
}

func (d *Database) Close() error {
	return d.DB.Close()
}