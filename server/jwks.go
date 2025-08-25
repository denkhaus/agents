// Package server provides JWKS functionality extracted from examples.
package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/mo"
)

// JWKSManager defines JWKS functionality for JWT validation and signing.
type JWKSManager interface {
	// GetJWKS returns the JSON Web Key Set
	GetJWKS() (map[string]interface{}, error)
	
	// ValidateToken validates a JWT token
	ValidateToken(tokenString string) (*jwt.Token, error)
	
	// SignPayload signs a payload and returns a JWT token
	SignPayload(payload map[string]interface{}) (string, error)
	
	// LoadOrGenerateKeys loads existing keys or generates new ones
	LoadOrGenerateKeys(keyDir string) error
	
	// CreateJWKSHandler creates an HTTP handler for JWKS endpoint
	CreateJWKSHandler() http.Handler
}

// jwksManager implements JWKSManager.
// Extracted from examples/jwks/server/main.go
type jwksManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	keyID      string
	issuer     string
	audience   string
}

// JWKSConfig defines JWKS configuration.
type JWKSConfig struct {
	KeySize  mo.Option[int]
	Issuer   mo.Option[string]
	Audience mo.Option[string]
	KeyID    mo.Option[string]
}

// NewJWKSManager creates a new JWKS manager.
func NewJWKSManager(config JWKSConfig) JWKSManager {
	return &jwksManager{
		keyID:    config.KeyID.OrElse("default-key"),
		issuer:   config.Issuer.OrElse("a2a-server"),
		audience: config.Audience.OrElse("a2a-client"),
	}
}

// LoadOrGenerateKeys loads existing keys or generates new ones.
func (jm *jwksManager) LoadOrGenerateKeys(keyDir string) error {
	privateKeyPath := filepath.Join(keyDir, "private.pem")
	publicKeyPath := filepath.Join(keyDir, "public.pem")

	// Try to load existing keys
	if err := jm.loadExistingKeys(privateKeyPath, publicKeyPath); err == nil {
		return nil
	}

	// Generate new keys
	return jm.generateAndSaveKeys(keyDir, privateKeyPath, publicKeyPath)
}

// loadExistingKeys loads existing RSA keys from files.
func (jm *jwksManager) loadExistingKeys(privateKeyPath, publicKeyPath string) error {
	// Load private key
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}

	privateKeyBlock, _ := pem.Decode(privateKeyData)
	if privateKeyBlock == nil {
		return fmt.Errorf("failed to decode private key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Load public key
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyData)
	if publicKeyBlock == nil {
		return fmt.Errorf("failed to decode public key PEM")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("public key is not RSA")
	}

	jm.privateKey = privateKey
	jm.publicKey = publicKey

	return nil
}

// generateAndSaveKeys generates new RSA keys and saves them to files.
func (jm *jwksManager) generateAndSaveKeys(keyDir, privateKeyPath, publicKeyPath string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(keyDir, 0700); err != nil {
		return fmt.Errorf("failed to create key directory: %w", err)
	}

	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Save private key
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to save private key: %w", err)
	}

	// Save public key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		return fmt.Errorf("failed to save public key: %w", err)
	}

	jm.privateKey = privateKey
	jm.publicKey = &privateKey.PublicKey

	return nil
}

// GetJWKS returns the JSON Web Key Set.
func (jm *jwksManager) GetJWKS() (map[string]interface{}, error) {
	if jm.publicKey == nil {
		return nil, fmt.Errorf("public key not loaded")
	}

	// Convert RSA public key to JWK format
	jwk := map[string]interface{}{
		"kty": "RSA",
		"use": "sig",
		"kid": jm.keyID,
		"n":   base64.RawURLEncoding.EncodeToString(jm.publicKey.N.Bytes()),
		"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(jm.publicKey.E)).Bytes()),
	}

	jwks := map[string]interface{}{
		"keys": []interface{}{jwk},
	}

	return jwks, nil
}

// ValidateToken validates a JWT token.
func (jm *jwksManager) ValidateToken(tokenString string) (*jwt.Token, error) {
	if jm.publicKey == nil {
		return nil, fmt.Errorf("public key not loaded")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jm.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	// Validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if err := jm.validateClaims(claims); err != nil {
			return nil, fmt.Errorf("invalid token claims: %w", err)
		}
	} else {
		return nil, fmt.Errorf("invalid token or claims")
	}

	return token, nil
}

// SignPayload signs a payload and returns a JWT token.
func (jm *jwksManager) SignPayload(payload map[string]interface{}) (string, error) {
	if jm.privateKey == nil {
		return "", fmt.Errorf("private key not loaded")
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(payload))

	// Set key ID in header
	token.Header["kid"] = jm.keyID

	// Sign token
	tokenString, err := token.SignedString(jm.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// CreateJWKSHandler creates an HTTP handler for JWKS endpoint.
func (jm *jwksManager) CreateJWKSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		jwks, err := jm.GetJWKS()
		if err != nil {
			http.Error(w, "Failed to get JWKS", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=3600")

		if err := json.NewEncoder(w).Encode(jwks); err != nil {
			http.Error(w, "Failed to encode JWKS", http.StatusInternalServerError)
		}
	})
}

// validateClaims validates JWT claims.
func (jm *jwksManager) validateClaims(claims jwt.MapClaims) error {
	// Validate audience
	if aud, ok := claims["aud"].(string); ok {
		if aud != jm.audience {
			return fmt.Errorf("invalid audience: expected %s, got %s", jm.audience, aud)
		}
	}

	// Validate issuer
	if iss, ok := claims["iss"].(string); ok {
		if iss != jm.issuer {
			return fmt.Errorf("invalid issuer: expected %s, got %s", jm.issuer, iss)
		}
	}

	// Validate expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return fmt.Errorf("token has expired")
		}
	}

	// Validate not before
	if nbf, ok := claims["nbf"].(float64); ok {
		if time.Now().Unix() < int64(nbf) {
			return fmt.Errorf("token not yet valid")
		}
	}

	return nil
}