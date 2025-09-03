package client

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// jwksVerifier implements JWKSVerifier.
type jwksVerifier struct {
	jwksURL      string
	keyset       jwk.Set
	keysetMutex  sync.RWMutex
	lastRefresh  time.Time
	refreshMutex sync.Mutex
	timeout      time.Duration
}

// NewJWKSVerifier creates a new JWKS verifier.
// Returns an error if the JWKS URL is invalid or initial fetch fails.
func NewJWKSVerifier(jwksURL string, timeout time.Duration) (JWKSVerifier, error) {
	if jwksURL == "" {
		return nil, fmt.Errorf("JWKS URL cannot be empty")
	}
	
	if timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive")
	}
	
	verifier := &jwksVerifier{
		jwksURL: jwksURL,
		keyset:  jwk.NewSet(),
		timeout: timeout,
	}

	// Initial JWKS fetch
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := verifier.RefreshJWKS(ctx); err != nil {
		return nil, fmt.Errorf("failed to fetch initial JWKS: %w", err)
	}

	return verifier, nil
}

// VerifyJWT implements JWKSVerifier.
func (v *jwksVerifier) VerifyJWT(ctx context.Context, tokenString string, payload []byte) error {
	keyset := v.getKeySet()
	if keyset == nil || keyset.Len() == 0 {
		return fmt.Errorf("no JWKS available")
	}

	// Calculate payload hash for verification
	hash := sha256.Sum256(payload)
	expectedPayloadHash := fmt.Sprintf("%x", hash)

	// Parse and verify token
	token, err := jwt.Parse([]byte(tokenString), jwt.WithKeySet(keyset), jwt.WithValidate(true))
	if err != nil {
		return fmt.Errorf("JWT verification failed: %w", err)
	}

	// Verify payload hash if present in token
	if requestBodyHash, exists := token.Get("request_body_sha256"); exists {
		if hashStr, ok := requestBodyHash.(string); ok {
			if hashStr != expectedPayloadHash {
				return fmt.Errorf("payload hash mismatch")
			}
		}
	}

	// Verify token age
	if iat, exists := token.Get("iat"); exists {
		if iatTime, ok := iat.(time.Time); ok {
			tokenAge := time.Since(iatTime)
			maxAge := 5 * time.Minute
			if tokenAge > maxAge {
				return fmt.Errorf("token has expired (age: %v)", tokenAge)
			}
		}
	}

	return nil
}

// RefreshJWKS implements JWKSVerifier.
func (v *jwksVerifier) RefreshJWKS(ctx context.Context) error {
	v.refreshMutex.Lock()
	defer v.refreshMutex.Unlock()

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: v.timeout,
	}

	// Fetch JWKS
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.jwksURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch JWKS: HTTP %d", resp.StatusCode)
	}

	// Read and parse JWKS
	keysetData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read JWKS: %w", err)
	}

	keyset, err := jwk.Parse(keysetData)
	if err != nil {
		return fmt.Errorf("failed to parse JWKS: %w", err)
	}

	// Update the keyset
	v.keysetMutex.Lock()
	v.keyset = keyset
	v.lastRefresh = time.Now()
	v.keysetMutex.Unlock()

	return nil
}

// getKeySet returns the current JWKS.
func (v *jwksVerifier) getKeySet() jwk.Set {
	v.keysetMutex.RLock()
	defer v.keysetMutex.RUnlock()
	return v.keyset
}
