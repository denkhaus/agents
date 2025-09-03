package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/denkhaus/agents/pkg/tools/project/repository/postgres/ent"
	"github.com/denkhaus/agents/pkg/tools/project/shared"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// postgresRepository implements the Repository interface using ent ORM
type postgresRepository struct {
	client *ent.Client
	config *Config
}

// NewPostgresRepository creates a new PostgreSQL repository using ent ORM
func NewRepository(opts ...Option) (shared.Repository, error) {
	config := DefaultConfig()

	repo := &postgresRepository{
		config: config,
	}

	// Apply options
	for _, opt := range opts {
		opt(repo)
	}

	if err := repo.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	return repo, nil
}

// initialize sets up the ent client and performs migrations
func (r *postgresRepository) initialize() error {
	// Open database connection

	r.config.Logger.Info("initialize database", zap.String("database_url", r.config.DatabaseURL))

	db, err := sql.Open("postgres", r.config.DatabaseURL)
	if err != nil {
		return NewConnectionError("failed to open database connection", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(r.config.MaxOpenConns)
	db.SetMaxIdleConns(r.config.MaxIdleConns)
	db.SetConnMaxLifetime(r.config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(r.config.ConnMaxIdleTime)

	// Test connection
	if err := db.Ping(); err != nil {
		return NewConnectionError("failed to ping database", err)
	}

	// Create ent client with PostgreSQL driver
	drv := entsql.OpenDB(dialect.Postgres, db)
	r.client = ent.NewClient(ent.Driver(drv))

	// Run auto-migration if enabled
	if r.config.AutoMigrate {
		ctx, cancel := context.WithTimeout(context.Background(), r.config.MigrationTimeout)
		defer cancel()

		if err := r.client.Schema.Create(ctx); err != nil {
			return NewMigrationError("auto-migration failed", err)
		}
	}

	return nil
}

// Close closes the ent client and database connection
func (r *postgresRepository) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// mapError converts ent/database errors to repository errors
func (r *postgresRepository) mapError(operation string, err error) error {
	if err == nil {
		return nil
	}

	if ent.IsNotFound(err) {
		return NewNotFoundError("resource", "unknown")
	}

	if ent.IsConstraintError(err) {
		return NewConstraintViolationError("constraint violation", err)
	}

	// TODO: Add more specific error mapping for different ent error types
	return NewConnectionError(fmt.Sprintf("database operation failed: %s", operation), err)
}
