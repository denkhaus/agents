// Package server provides authentication functionality extracted from examples.
package server

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/mo"
)

// SimpleAuthProvider defines a simple authentication provider interface.
type SimpleAuthProvider interface {
	// Authenticate validates a request
	Authenticate(req *http.Request) (*AuthContext, error)
	
	// Type returns the provider type
	Type() string
}

// AuthenticationManager manages multiple authentication providers.
type AuthenticationManager interface {
	// Authenticate validates a request using the provider chain
	Authenticate(req *http.Request) (*AuthContext, error)
	
	// AddProvider adds an authentication provider
	AddProvider(provider SimpleAuthProvider) error
	
	// GetProviders returns all registered providers
	GetProviders() []SimpleAuthProvider
}

// authenticationManager implements AuthenticationManager.
type authenticationManager struct {
	providers []SimpleAuthProvider
}

// NewAuthenticationManager creates a new authentication manager.
func NewAuthenticationManager(providers ...SimpleAuthProvider) AuthenticationManager {
	return &authenticationManager{
		providers: providers,
	}
}

// Authenticate validates a request using the provider chain.
func (am *authenticationManager) Authenticate(req *http.Request) (*AuthContext, error) {
	for _, provider := range am.providers {
		// Try each provider in sequence
		if authContext, err := provider.Authenticate(req); err == nil {
			return authContext, nil
		}
	}
	return nil, fmt.Errorf("authentication failed")
}

// AddProvider adds an authentication provider.
func (am *authenticationManager) AddProvider(provider SimpleAuthProvider) error {
	am.providers = append(am.providers, provider)
	return nil
}

// GetProviders returns all registered providers.
func (am *authenticationManager) GetProviders() []SimpleAuthProvider {
	return am.providers
}

// AuthConfig defines authentication configuration extracted from examples/auth.
type AuthConfig struct {
	JWTSecret    mo.Option[[]byte]
	JWTAudience  mo.Option[string]
	JWTIssuer    mo.Option[string]
	APIKeys      mo.Option[map[string]string]
	APIKeyHeader mo.Option[string]
	EnableOAuth  mo.Option[bool]
}

// CreateAuthProviders creates authentication providers from configuration.
// Simplified implementation without external dependencies
func CreateAuthProviders(config AuthConfig) ([]SimpleAuthProvider, error) {
	var providers []SimpleAuthProvider
	
	// API Key Provider
	if apiKeys, hasKeys := config.APIKeys.Get(); hasKeys {
		header := config.APIKeyHeader.OrElse("X-API-Key")
		apiKeyProvider := NewAPIKeyProvider(apiKeys, header)
		providers = append(providers, apiKeyProvider)
	}
	
	if len(providers) == 0 {
		return nil, fmt.Errorf("no authentication providers configured")
	}
	
	return providers, nil
}

// apiKeyProvider implements SimpleAuthProvider for API key authentication.
type apiKeyProvider struct {
	apiKeys map[string]string
	header  string
}

// NewAPIKeyProvider creates a new API key provider.
func NewAPIKeyProvider(apiKeys map[string]string, header string) SimpleAuthProvider {
	return &apiKeyProvider{
		apiKeys: apiKeys,
		header:  header,
	}
}

// Authenticate validates API key authentication.
func (p *apiKeyProvider) Authenticate(req *http.Request) (*AuthContext, error) {
	apiKey := req.Header.Get(p.header)
	if apiKey == "" {
		return nil, fmt.Errorf("missing API key")
	}
	
	if _, valid := p.apiKeys[apiKey]; !valid {
		return nil, fmt.Errorf("invalid API key")
	}
	
	return &AuthContext{
		UserID:       "api-user",
		ProviderType: "apikey",
		Claims:       map[string]interface{}{"api_key": apiKey},
	}, nil
}

// Type returns the provider type.
func (p *apiKeyProvider) Type() string {
	return "apikey"
}

// JWTSecretManager manages JWT secrets.
// Extracted from examples/auth/server/main.go
type JWTSecretManager interface {
	// LoadOrGenerateSecret loads or generates a JWT secret
	LoadOrGenerateSecret(secretFile string) ([]byte, error)
}

// jwtSecretManager implements JWTSecretManager.
type jwtSecretManager struct{}

// NewJWTSecretManager creates a new JWT secret manager.
func NewJWTSecretManager() JWTSecretManager {
	return &jwtSecretManager{}
}

// LoadOrGenerateSecret loads a JWT secret from file or generates a new one.
func (jsm *jwtSecretManager) LoadOrGenerateSecret(secretFile string) ([]byte, error) {
	// Try to load existing secret
	data, err := os.ReadFile(secretFile)
	if err == nil && len(data) >= 32 {
		return data, nil
	}
	
	// Generate new secret
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("failed to generate JWT secret: %w", err)
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(secretFile)
	if dir != "." {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create directory for JWT secret: %w", err)
		}
	}
	
	// Save for future use with tight permissions
	if err := os.WriteFile(secretFile, secret, 0600); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: Could not save JWT secret to %s: %v\n", secretFile, err)
	}
	
	return secret, nil
}

// MockOAuthServer provides OAuth2 mock functionality.
// Extracted from examples/auth/server/main.go
type MockOAuthServer interface {
	// Start initializes OAuth server handlers
	Start(mux *http.ServeMux)
	
	// GetTokenEndpoint returns the token endpoint path
	GetTokenEndpoint() string
}

// mockOAuthServer implements MockOAuthServer.
type mockOAuthServer struct {
	validCredentials map[string]string
	tokenEndpoint    string
}

// NewMockOAuthServer creates a new mock OAuth2 server.
func NewMockOAuthServer() MockOAuthServer {
	return &mockOAuthServer{
		validCredentials: map[string]string{
			"my-client-id": "my-client-secret",
		},
		tokenEndpoint: "/oauth2/token",
	}
}

// Start initializes the OAuth server handlers.
func (m *mockOAuthServer) Start(mux *http.ServeMux) {
	mux.HandleFunc(m.tokenEndpoint, m.handleTokenRequest)
}

// GetTokenEndpoint returns the token endpoint path.
func (m *mockOAuthServer) GetTokenEndpoint() string {
	return m.tokenEndpoint
}

// handleTokenRequest processes OAuth2 token requests.
func (m *mockOAuthServer) handleTokenRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	
	// Validate grant type
	if r.FormValue("grant_type") != "client_credentials" {
		http.Error(w, "Unsupported grant type", http.StatusBadRequest)
		return
	}
	
	// Get client credentials
	clientID, clientSecret := getClientCredentials(r)
	if clientID == "" || clientSecret == "" {
		w.Header().Set("WWW-Authenticate", `Basic realm="OAuth2 Server"`)
		http.Error(w, "Missing client credentials", http.StatusUnauthorized)
		return
	}
	
	// Validate credentials
	validSecret, ok := m.validCredentials[clientID]
	if !ok || validSecret != clientSecret {
		http.Error(w, "Invalid client credentials", http.StatusUnauthorized)
		return
	}
	
	// Generate token response
	scopes := strings.Split(r.FormValue("scope"), " ")
	token := generateTokenResponse(clientID, scopes)
	
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{
		"access_token": "%s",
		"token_type": "Bearer",
		"expires_in": 3600,
		"scope": "%s"
	}`, token, strings.Join(scopes, " "))))
}

// getClientCredentials extracts client credentials from request.
func getClientCredentials(r *http.Request) (string, string) {
	// Try Basic auth first
	clientID, clientSecret, ok := r.BasicAuth()
	if ok && clientID != "" && clientSecret != "" {
		return clientID, clientSecret
	}
	
	// Try form parameters
	return r.FormValue("client_id"), r.FormValue("client_secret")
}

// generateTokenResponse generates a mock access token.
func generateTokenResponse(clientID string, scopes []string) string {
	return "mock-access-token-" + clientID
}