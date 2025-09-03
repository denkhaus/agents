package auth

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// oauth2Provider implements OAuth2 authentication.
type oauth2Provider struct {
	config      OAuth2Config
	accessToken string
	expiry      time.Time
	mu          sync.RWMutex
}

// Authenticate implements Provider.
func (p *oauth2Provider) Authenticate(req interface{}) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isValidLocked() {
		return fmt.Errorf("OAuth2 token is expired or invalid")
	}

	// Add OAuth2 token to request
	// Implementation depends on request type
	return nil
}

// IsValid implements Provider.
func (p *oauth2Provider) IsValid() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isValidLocked()
}

// Refresh implements Provider.
func (p *oauth2Provider) Refresh(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Refresh OAuth2 token using client credentials flow
	token, expiry, err := p.refreshToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to refresh OAuth2 token: %w", err)
	}

	p.accessToken = token
	p.expiry = expiry

	return nil
}

// Type implements Provider.
func (p *oauth2Provider) Type() string {
	return "oauth2"
}

// isValidLocked checks if token is valid (must be called with lock held).
func (p *oauth2Provider) isValidLocked() bool {
	return p.accessToken != "" && time.Now().Before(p.expiry)
}

// refreshToken refreshes the OAuth2 access token.
func (p *oauth2Provider) refreshToken(ctx context.Context) (string, time.Time, error) {
	// Implementation would use OAuth2 library to refresh token
	// This is a placeholder
	expiry := time.Now().Add(1 * time.Hour)
	return "oauth2-token-placeholder", expiry, nil
}
