# Unified A2A Client Package

Ein konsolidiertes, thread-sicheres Client-Paket für A2A (Agent-to-Agent) Kommunikation, das alle Funktionalitäten aus den Beispiel-Implementierungen in einer einheitlichen API zusammenfasst.

## Features

### 🔌 Core Interfaces
- **APIClient**: Grundlegende A2A-Operationen (Nachrichten senden, Tasks verwalten, Push-Benachrichtigungen)
- **StreamingClient**: Echtzeit-Streaming für Nachrichten
- **UnifiedClient**: Kombiniert alle Client-Funktionalitäten

### 🔐 Authentication
- **JWT**: Konfigurierbare Secrets, Audiences und Ablaufzeiten
- **API Key**: Benutzerdefinierte Header-basierte Authentifizierung
- **OAuth2**: Client Credentials Flow

### 🛠 Utility Features
- **AgentCardFetcher**: Automatische Erkennung von Agent-Fähigkeiten
- **JWKSVerifier**: JWT-Token-Verifizierung mit JWKS
- **TaskTracker**: Task-Status-Verfolgung und -Management
- **RequestHandler**: HTTP-Requests mit Retry-Logik
- **ResponseParser**: Strukturierte Response-Verarbeitung

## Installation

```bash
go get github.com/denkhaus/agents/client
```

## Quick Start

### Basic API Client

```go
import (
    "github.com/denkhaus/agents/client"
    "trpc.group/trpc-go/trpc-a2a-go/protocol"
)

// Einfacher API Client
apiClient, err := client.NewAPIClient("http://localhost:8080",
    client.WithTimeout(30*time.Second))

// Nachricht senden
message := protocol.NewMessage(
    protocol.MessageRoleUser,
    []protocol.Part{protocol.NewTextPart("Hello, world!")},
)

result, err := apiClient.SendMessage(ctx, protocol.SendMessageParams{
    Message: message,
})
```

### Streaming-Enabled Client mit Authentication

```go
import (
    "github.com/denkhaus/agents/client"
    "github.com/denkhaus/agents/client/auth"
)

// JWT Authentication
jwtAuth := auth.NewJWTProvider(auth.JWTConfig{
    Secret:   []byte("my-secret-key"),
    Audience: "a2a-server",
    Issuer:   "example-client",
    Expiry:   time.Hour,
})

// Client mit Streaming aktiviert
streamClient, err := client.NewAPIClient("http://localhost:8080",
    client.WithAuth(jwtAuth),
    client.WithStreaming(true),
    client.WithTimeout(60*time.Second))

// Stream Events verarbeiten (nur wenn Streaming aktiviert)
if streamingClient, ok := streamClient.(client.StreamingClient); ok {
    eventChan, err := streamingClient.StreamMessage(ctx, params)
    for event := range eventChan {
        // Event verarbeiten
    }
}
```

### Unified Client (Alle Features)

```go
// Konfiguration mit mo.Option für type-safe defaults
config := client.ClientConfig{
    BaseURL:    "http://localhost:8080",
    Timeout:    mo.Some(30 * time.Second),
    UserAgent:  mo.Some("my-app/1.0"),
    Auth:       mo.Some(auth),
    Streaming:  mo.Some(true),
    RetryCount: mo.Some(3),
}

unifiedClient, err := client.NewUnifiedClient(config)

// Agent Capabilities abrufen
agentCard, err := unifiedClient.FetchAgentCard(ctx, baseURL)

// Standard oder Streaming basierend auf Agent-Fähigkeiten
if agentCard.Capabilities.Streaming != nil && *agentCard.Capabilities.Streaming {
    eventChan, err := unifiedClient.StreamMessage(ctx, params)
} else {
    result, err := unifiedClient.SendMessage(ctx, params)
}
```

## Authentication Providers

### JWT Authentication

```go
jwtAuth := auth.NewJWTProvider(auth.JWTConfig{
    Secret:   []byte("your-secret-key"),
    Audience: "a2a-server",
    Issuer:   "your-client",
    Expiry:   time.Hour,
})
```

### API Key Authentication

```go
apiKeyAuth := auth.NewAPIKeyProvider(auth.APIKeyConfig{
    Key:    "your-api-key",
    Header: mo.Some("X-API-Key"), // Optional, default: "X-API-Key"
})
```

### OAuth2 Authentication

```go
oauth2Auth := auth.NewOAuth2Provider(auth.OAuth2Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    TokenURL:     "https://auth.example.com/oauth2/token",
    Scopes:       []string{"a2a.read", "a2a.write"},
})
```

## Task Management

```go
// Task Tracker erstellen
tracker := client.NewTaskTracker()

// Task verfolgen
tracker.TrackTask("task-123")
tracker.UpdateTaskStatus("task-123", "running")
tracker.UpdateTaskStatus("task-123", "completed")

// Status abrufen
status := tracker.GetTaskStatus("task-123")
```

## JWKS Verification

```go
// JWKS Verifier für Push Notifications
verifier := client.NewJWKSVerifier(
    "http://localhost:8080/.well-known/jwks.json",
    10*time.Second,
)

// JWT Token verifizieren
err := verifier.VerifyJWT(ctx, jwtToken, requestPayload)
```

## Error Handling

Das Paket bietet strukturierte Error-Typen:

```go
import "github.com/denkhaus/agents/client/helper"

// HTTP Errors
if httpErr, ok := err.(*helper.HTTPError); ok {
    if httpErr.IsClientError() {
        // 4xx Fehler behandeln
    }
    if httpErr.IsServerError() {
        // 5xx Fehler behandeln
    }
}

// Authentication Errors
if authErr, ok := err.(*helper.AuthError); ok {
    // Auth-spezifische Behandlung
}
```

## Thread Safety

Alle Client-Implementierungen sind thread-safe und können sicher von mehreren Goroutines verwendet werden.

## Architecture

```
client/
├── interfaces.go          # Core interfaces
├── factory.go            # Factory functions
├── unified_client.go     # Unified implementation
├── agentcard.go         # Agent capability discovery
├── jwks.go              # JWKS verification
├── tracker.go           # Task tracking
├── auth/                # Authentication providers
│   ├── interfaces.go
│   ├── jwt.go
│   ├── apikey.go
│   └── oauth2.go
└── helper/              # Utility functions
    ├── request.go       # HTTP request handling
    ├── response.go      # Response parsing
    ├── errors.go        # Error types
    └── utils.go         # General utilities
```

## Design Principles

- **Modular**: Jede Komponente kann unabhängig verwendet werden
- **Type-Safe**: Verwendung von `samber/mo` für optionale Konfiguration
- **Thread-Safe**: Alle Implementierungen sind concurrent-safe
- **Interface-Driven**: Klare Trennung zwischen Interface und Implementierung
- **Dependency Injection**: Flexible Konfiguration über Factory-Pattern
- **Error-Transparent**: Strukturierte Error-Typen für bessere Fehlerbehandlung

## Examples

Siehe `example_test.go` für vollständige Verwendungsbeispiele aller Features.