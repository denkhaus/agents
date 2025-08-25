package auth

import (
	"context"
)

// apiKeyProvider implements API key authentication.
type apiKeyProvider struct {
	key    string
	header string
}

// Authenticate implements Provider.
func (p *apiKeyProvider) Authenticate(req interface{}) error {
	// Add API key to request headers
	// Implementation depends on request type
	return nil
}

// IsValid implements Provider.
func (p *apiKeyProvider) IsValid() bool {
	return p.key != ""
}

// Refresh implements Provider.
func (p *apiKeyProvider) Refresh(ctx context.Context) error {
	// API keys don't need refreshing
	return nil
}

// Type implements Provider.
func (p *apiKeyProvider) Type() string {
	return "apikey"
}
