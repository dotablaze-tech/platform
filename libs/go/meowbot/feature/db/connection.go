package db

import (
	"context"
	"database/sql"
	"fmt"
	"libs/go/meowbot/util"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection and handles connection pooling and context management.
func InitDB(ctx context.Context) error {
	connStr := util.Cfg.DatabaseURL
	if connStr == "" {
		// Log a warning or use a fallback connection string for local dev environments
		connStr = "postgres://default_user:default_password@127.0.0.1:5432/meowbot?sslmode=disable"
	} else {
		user := util.Cfg.DatabaseUser
		pass := util.Cfg.DatabasePassword
		host := util.Cfg.DatabaseHost
		port := util.Cfg.DatabasePort
		name := util.Cfg.DatabaseName

		connStr = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			user, pass, host, port, name,
		)
	}

	// Open the database connection
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open DB connection: %w", err)
	}

	// Set connection pool options
	DB.SetMaxOpenConns(10)                  // Maximum number of open connections to the database
	DB.SetMaxIdleConns(5)                   // Maximum number of idle connections in the pool
	DB.SetConnMaxLifetime(30 * time.Minute) // Maximum amount of time a connection may be reused

	// Attempt to ping the database to ensure the connection is live
	if err := DB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	// Connection is successful
	return nil
}

// CloseDB ensures that the DB connection is properly closed when the application shuts down.
func CloseDB() error {
	if DB != nil {
		if err := DB.Close(); err != nil {
			return fmt.Errorf("failed to close DB connection: %w", err)
		}
	}
	return nil
}
