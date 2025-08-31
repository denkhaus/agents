package config

import (
	"errors"
	"fmt"

	"github.com/denkhaus/agents/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/samber/do"
)

// Service defines the interface for configuration retrieval.
type Service interface {
	GetConfigPath() (string, error)
	GetWorkspacePath() (string, error)
}

// Specification for environment variables
type Specification struct {
	ConfigPath    string `envconfig:"CONFIG_PATH" default:"./config"`
	WorkspacePath string `envconfig:"WORKSPACE_PATH"`
}

// AppConfig is a concrete implementation of ConfigService.
type AppConfig struct {
	spec Specification
}

// NewWithDI creates a new AppConfig instance.
func NewWithDI(i *do.Injector) (Service, error) {

	// Try to load .env file from current directory or parent directories
	envPaths := []string{".env", "../.env", "../../.env"}
	var envLoaded bool
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			envLoaded = true
			break
		}
	}
	if !envLoaded {
		// Fallback: try to find .env in the project root
		if err := godotenv.Load(); err != nil {
			logger.Log.Warn("No .env file found, using environment variable defaults")
		}
	}

	var spec Specification
	err := envconfig.Process("AGENTS", &spec)
	if err != nil {
		return nil, fmt.Errorf("failed to process envconfig: %w", err)
	}

	return &AppConfig{
		spec: spec,
	}, nil
}

// GetBasePath returns the base path of the application.
func (c *AppConfig) GetConfigPath() (string, error) {
	if c.spec.ConfigPath == "" {
		return "", errors.New("failed to get config path. please define AGENTS_CONFIG_PATH variable")
	}

	return c.spec.ConfigPath, nil
}

// GetWorkspacePath returns the workspace path of the application.
func (c *AppConfig) GetWorkspacePath() (string, error) {
	if c.spec.ConfigPath == "" {
		return "", errors.New("failed to get workspace path from config. please define AGENTS_WORKSPACE_PATH variable")
	}

	return c.spec.WorkspacePath, nil
}
