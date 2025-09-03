package auth

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// jwtProvider implements JWT authentication.
type jwtProvider struct {
	config JWTConfig
	token  string
	expiry time.Time
	mu     sync.RWMutex
}

// Authenticate implements Provider.
func (p *jwtProvider) Authenticate(req interface{}) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isValidLocked() {
		return fmt.Errorf("JWT token is expired or invalid")
	}

	// Add JWT token to request
	// Implementation depends on request type
	return nil
}

// IsValid implements Provider.
func (p *jwtProvider) IsValid() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isValidLocked()
}

// Refresh implements Provider.
func (p *jwtProvider) Refresh(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Generate new JWT token
	token, err := p.generateToken()
	if err != nil {
		return fmt.Errorf("failed to generate JWT token: %w", err)
	}

	p.token = token
	p.expiry = time.Now().Add(p.config.Expiry)

	return nil
}

// Type implements Provider.
func (p *jwtProvider) Type() string {
	return "jwt"
}

// isValidLocked checks if token is valid (must be called with lock held).
func (p *jwtProvider) isValidLocked() bool {
	return p.token != "" && time.Now().Before(p.expiry)
}

// generateToken generates a new JWT token.
func (p *jwtProvider) generateToken() (string, error) {
	// Implementation would use JWT library to generate token
	// This is a placeholder
	return "jwt-token-placeholder", nil
}
