package postgres

import (
	"time"
)

// Config holds configuration for the PostgreSQL repository
type Config struct {
	// Database connection settings
	DatabaseURL     string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration

	// Migration settings
	AutoMigrate      bool
	MigrationTimeout time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		DatabaseURL:      "postgres://localhost/project_management?sslmode=disable",
		MaxOpenConns:     25,
		MaxIdleConns:     5,
		ConnMaxLifetime:  time.Hour,
		ConnMaxIdleTime:  time.Minute * 15,
		AutoMigrate:      true,
		MigrationTimeout: time.Minute * 5,
	}
}

// Option is a function that configures a PostgreSQL repository
type Option func(*postgresRepository)

// WithConfig sets the entire configuration
func WithConfig(config *Config) Option {
	return func(r *postgresRepository) {
		r.config = config
	}
}

// WithDatabaseURL sets the database URL
func WithDatabaseURL(url string) Option {
	return func(r *postgresRepository) {
		r.config.DatabaseURL = url
	}
}

// WithAutoMigrate enables or disables auto-migration
func WithAutoMigrate(enable bool) Option {
	return func(r *postgresRepository) {
		r.config.AutoMigrate = enable
	}
}

// WithConnectionPool configures the connection pool
func WithConnectionPool(maxOpen, maxIdle int) Option {
	return func(r *postgresRepository) {
		r.config.MaxOpenConns = maxOpen
		r.config.MaxIdleConns = maxIdle
	}
}

// WithConnectionLifetime configures connection lifetimes
func WithConnectionLifetime(maxLifetime, maxIdleTime time.Duration) Option {
	return func(r *postgresRepository) {
		r.config.ConnMaxLifetime = maxLifetime
		r.config.ConnMaxIdleTime = maxIdleTime
	}
}

// WithMigrationTimeout sets the migration timeout
func WithMigrationTimeout(timeout time.Duration) Option {
	return func(r *postgresRepository) {
		r.config.MigrationTimeout = timeout
	}
}
