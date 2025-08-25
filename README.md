# Agents - A2A Communication Library

A unified, production-ready Go library for Agent-to-Agent (A2A) communication with comprehensive authentication, streaming, and server capabilities.

## Features

- **Multiple Authentication Methods**: JWT, API Key, OAuth2
- **Streaming Support**: Real-time message streaming
- **Agent Discovery**: Automatic capability discovery
- **Task Management**: Task tracking and status monitoring
- **JWKS Verification**: JWT verification with JSON Web Key Sets
- **Multi-Agent Routing**: Support for multiple agent endpoints
- **Push Notifications**: Real-time notification support
- **Production Ready**: Comprehensive error handling, logging, and metrics

## Architecture

The library follows Go best practices with:
- **Interface-Driven Design**: All functionality exposed through well-defined interfaces
- **Factory Pattern**: Type-safe configuration through factory functions
- **Unexported Implementations**: All concrete types are unexported for better encapsulation
- **Option Pattern**: Flexible configuration using `mo.Option` for type safety
- **Context Support**: Full context.Context support for cancellation and timeouts
- **Thread Safety**: All components are thread-safe

## Quick Start

### Client Usage

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/denkhaus/agents/auth"
    "github.com/denkhaus/agents/client"
    "github.com/samber/mo"
)

func main() {
    // Create authentication provider
    authProvider, err := auth.NewAPIKeyProvider(auth.APIKeyConfig{
        Key:    "your-api-key",
        Header: mo.Some("X-Custom-Key"),
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create unified client
    client, err := client.NewUnifiedClient(client.ClientConfig{
        BaseURL: "https://api.example.com",
        Timeout: mo.Some(30 * time.Second),
        Auth:    mo.Some(authProvider),
    })
    if err != nil {
        log.Fatal(err)
    }

    // Send a message
    ctx := context.Background()
    result, err := client.SendMessage(ctx, params)
    if err != nil {
        log.Printf("Failed to send message: %v", err)
        return
    }

    log.Printf("Message sent successfully: %+v", result)
}
```

### Server Usage

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/denkhaus/agents/server"
    "trpc.group/trpc-go/trpc-a2a-go/server"
    "trpc.group/trpc-go/trpc-a2a-go/taskmanager"
)

func main() {
    // Create agent card
    agentCard := server.AgentCard{
        Name:        "My Agent",
        Description: "A sample A2A agent",
        Version:     "1.0.0",
    }

    // Create task manager
    taskManager := taskmanager.NewTaskManager()

    // Create server
    srv, err := server.NewServer(agentCard, taskManager,
        server.WithPort(8080),
        server.WithTLS("cert.pem", "key.pem"),
        server.WithCORS(true),
        server.WithTimeout(30*time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Start server
    log.Println("Starting server on :8080")
    if err := srv.Start(":8080"); err != nil {
        log.Printf("Server failed: %v", err)
    }
}
```

## Package Structure

```
github.com/denkhaus/agents/
├── doc.go                 # Package documentation
├── README.md             # This file
├── auth/                 # Authentication providers
│   ├── interfaces.go     # Authentication interfaces and configs
│   ├── jwt.go           # JWT authentication implementation
│   ├── apikey.go        # API Key authentication implementation
│   └── oauth2.go        # OAuth2 authentication implementation
├── client/              # Client implementations
│   ├── interfaces.go    # Client interfaces and configs
│   ├── factory.go       # Client factory functions
│   ├── unified_client.go # Main client implementation
│   ├── agentcard.go     # Agent capability discovery
│   ├── jwks.go          # JWKS verification
│   └── tracker.go       # Task tracking
├── server/              # Server implementations
│   ├── interfaces.go    # Server interfaces and configs
│   ├── factory.go       # Server factory functions
│   ├── unified_server.go # Main server implementation
│   ├── auth.go          # Server authentication
│   ├── streaming.go     # Streaming functionality
│   ├── multiagent.go    # Multi-agent routing
│   └── notifications.go # Push notifications
└── helper/              # Utility functions
    ├── errors.go        # Error types
    ├── request.go       # HTTP request utilities
    ├── response.go      # HTTP response utilities
    └── utils.go         # General utilities
```

## Authentication

### JWT Authentication

```go
jwtProvider, err := auth.NewJWTProvider(auth.JWTConfig{
    Secret:   []byte("your-secret"),
    Audience: "your-audience",
    Issuer:   "your-issuer",
    Expiry:   time.Hour,
})
```

### API Key Authentication

```go
apiProvider, err := auth.NewAPIKeyProvider(auth.APIKeyConfig{
    Key:    "your-api-key",
    Header: mo.Some("X-Custom-Key"), // Optional, defaults to "X-API-Key"
})
```

### OAuth2 Authentication

```go
oauth2Provider, err := auth.NewOAuth2Provider(auth.OAuth2Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    TokenURL:     "https://auth.example.com/token",
    Scopes:       []string{"read", "write"},
})
```

## Error Handling

The library provides comprehensive error handling with specific error types:

```go
import "github.com/denkhaus/agents/helper"

// HTTP errors
if httpErr, ok := err.(*helper.HTTPError); ok {
    if httpErr.IsClientError() {
        // Handle 4xx errors
    } else if httpErr.IsServerError() {
        // Handle 5xx errors
    }
}

// Authentication errors
if authErr, ok := err.(*helper.AuthError); ok {
    log.Printf("Auth error (%s): %s", authErr.Type, authErr.Message)
}

// Validation errors
if valErr, ok := err.(*helper.ValidationError); ok {
    log.Printf("Validation error for %s: %s", valErr.Field, valErr.Message)
}
```

## Configuration

All components use the Option pattern for type-safe configuration:

```go
// Client configuration
client, err := client.NewUnifiedClient(client.ClientConfig{
    BaseURL:    "https://api.example.com",           // Required
    Timeout:    mo.Some(30 * time.Second),          // Optional
    UserAgent:  mo.Some("my-app/1.0"),              // Optional
    Auth:       mo.Some(authProvider),               // Optional
    Streaming:  mo.Some(true),                      // Optional
    RetryCount: mo.Some(3),                         // Optional
})

// Server configuration
server, err := server.NewServer(agentCard, taskManager,
    server.WithHost("localhost"),                    // Optional
    server.WithPort(8080),                          // Optional
    server.WithTLS("cert.pem", "key.pem"),         // Optional
    server.WithTimeout(30*time.Second),             // Optional
    server.WithMaxRequestSize(10*1024*1024),       // Optional
    server.WithCORS(true),                          // Optional
    server.WithMetrics(true),                       // Optional
    server.WithShutdownTimeout(30*time.Second),     // Optional
)
```

## Thread Safety

All components are designed to be thread-safe:
- Authentication providers can be used concurrently
- Clients support concurrent requests
- Servers handle concurrent connections
- Task trackers support concurrent access

## Testing

The library includes comprehensive test coverage. Run tests with:

```bash
go test ./...
```

## Contributing

1. Follow Go best practices and idioms
2. Maintain interface-driven design
3. Keep implementations unexported
4. Add comprehensive documentation
5. Include tests for new functionality
6. Use the Option pattern for configuration

## License

[Add your license information here]