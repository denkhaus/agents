// Package auth provides authentication implementations for A2A communication.
//
// This package defines authentication interfaces and provides factory functions
// for creating authentication providers. All concrete implementations are unexported
// and accessed through the Provider interface.
//
// Supported authentication methods:
//   - JWT (JSON Web Tokens)
//   - API Key authentication
//   - OAuth2 client credentials flow
//
// Example usage:
//
//	// JWT Provider
//	jwtProvider := auth.NewJWTProvider(auth.JWTConfig{
//		Secret:   []byte("your-secret"),
//		Audience: "your-audience",
//		Issuer:   "your-issuer",
//		Expiry:   time.Hour,
//	})
//
//	// API Key Provider
//	apiProvider := auth.NewAPIKeyProvider(auth.APIKeyConfig{
//		Key:    "your-api-key",
//		Header: mo.Some("X-Custom-Key"),
//	})
//
//	// OAuth2 Provider
//	oauth2Provider := auth.NewOAuth2Provider(auth.OAuth2Config{
//		ClientID:     "your-client-id",
//		ClientSecret: "your-client-secret",
//		TokenURL:     "https://auth.example.com/token",
//		Scopes:       []string{"read", "write"},
//	})
package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/mo"
)

// Provider defines the authentication provider interface.
// All authentication providers must implement this interface.
type Provider interface {
	// Authenticate adds authentication to the request.
	// The req parameter can be *http.Request or other request types.
	// Returns an error if authentication cannot be applied.
	Authenticate(req interface{}) error

	// IsValid checks if the current authentication is still valid.
	// Returns false if authentication has expired or is invalid.
	IsValid() bool

	// Refresh refreshes authentication if needed and possible.
	// Should be called when IsValid() returns false.
	// Returns an error if refresh fails.
	Refresh(ctx context.Context) error

	// Type returns the authentication type identifier.
	// Used for logging and debugging purposes.
	Type() string
}

// JWTConfig defines JWT authentication configuration.
type JWTConfig struct {
	// Secret is the signing key for JWT tokens
	Secret []byte
	// Audience specifies the intended audience for the JWT
	Audience string
	// Issuer specifies the JWT issuer
	Issuer string
	// Expiry defines how long tokens are valid
	Expiry time.Duration
}

// Validate validates the JWT configuration.
func (c JWTConfig) Validate() error {
	if len(c.Secret) == 0 {
		return fmt.Errorf("JWT secret cannot be empty")
	}
	if c.Audience == "" {
		return fmt.Errorf("JWT audience cannot be empty")
	}
	if c.Issuer == "" {
		return fmt.Errorf("JWT issuer cannot be empty")
	}
	if c.Expiry <= 0 {
		return fmt.Errorf("JWT expiry must be positive")
	}
	return nil
}

// APIKeyConfig defines API key authentication configuration.
type APIKeyConfig struct {
	// Key is the API key value
	Key string
	// Header specifies the HTTP header name for the API key.
	// Defaults to "X-API-Key" if not specified.
	Header mo.Option[string]
}

// Validate validates the API key configuration.
func (c APIKeyConfig) Validate() error {
	if c.Key == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	return nil
}

// OAuth2Config defines OAuth2 authentication configuration.
type OAuth2Config struct {
	// ClientID is the OAuth2 client identifier
	ClientID string
	// ClientSecret is the OAuth2 client secret
	ClientSecret string
	// TokenURL is the OAuth2 token endpoint
	TokenURL string
	// Scopes defines the requested OAuth2 scopes
	Scopes []string
}

// Validate validates the OAuth2 configuration.
func (c OAuth2Config) Validate() error {
	if c.ClientID == "" {
		return fmt.Errorf("OAuth2 client ID cannot be empty")
	}
	if c.ClientSecret == "" {
		return fmt.Errorf("OAuth2 client secret cannot be empty")
	}
	if c.TokenURL == "" {
		return fmt.Errorf("OAuth2 token URL cannot be empty")
	}
	return nil
}

// NewJWTProvider creates a JWT authentication provider.
// Returns an error if the configuration is invalid.
func NewJWTProvider(config JWTConfig) (Provider, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid JWT configuration: %w", err)
	}
	
	return &jwtProvider{
		config: config,
	}, nil
}

// NewAPIKeyProvider creates an API key authentication provider.
// Returns an error if the configuration is invalid.
func NewAPIKeyProvider(config APIKeyConfig) (Provider, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid API key configuration: %w", err)
	}
	
	header := config.Header.OrElse("X-API-Key")
	return &apiKeyProvider{
		key:    config.Key,
		header: header,
	}, nil
}

// NewOAuth2Provider creates an OAuth2 authentication provider.
// Returns an error if the configuration is invalid.
func NewOAuth2Provider(config OAuth2Config) (Provider, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid OAuth2 configuration: %w", err)
	}
	
	return &oauth2Provider{
		config: config,
	}, nil
}
