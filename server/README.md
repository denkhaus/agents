# Unified A2A Server Package

Ein konsolidiertes, production-ready Server-Package für A2A (Agent-to-Agent) Kommunikation, das alle wertvollen Funktionalitäten aus den Beispiel-Implementierungen in einer einheitlichen, erweiterbaren API zusammenfasst.

## 🎯 Tier 1 Features (Implementiert)

### ✅ **Unified Authentication System**
- **Multi-Provider Support**: JWT, API Key, OAuth2 in einer Chain
- **Flexible Configuration**: Type-safe Konfiguration mit `samber/mo`
- **Request Context**: Authentication Context für alle Requests
- **Wiederverwendbare Auth Components**: Nutzt gemeinsame `auth/` Package

### ✅ **Message Processing Framework** 
- **Streaming + Non-Streaming**: Automatische Mode-Detection
- **Base Processor**: Wiederverwendbare Basis-Implementierung
- **Capability System**: Deklarative Processor-Fähigkeiten
- **Echo Processor**: Vollständige Referenz-Implementierung

### ✅ **Task Manager Integration**
- **Memory + Redis Backends**: Pluggable Task Manager Provider
- **Agent-specific Task Managers**: Isolierte Task-Verarbeitung pro Agent
- **Graceful Lifecycle**: Proper Resource Management

### ✅ **Agent Card Management**
- **Dynamic Agent Registration**: Runtime Agent hinzufügen/entfernen
- **Agent Manager**: Centralized Agent Lifecycle Management
- **Path-based Routing**: Flexible Agent URL-Strukturen
- **Agent Metadata**: Registration Time, Last Used Tracking

## 🏗 Architecture

```
server/
├── interfaces.go              # Core Server Interfaces
├── factory.go                # Server Factory mit Options Pattern
├── unified_server.go         # Hauptserver-Implementierung
├── auth_manager.go           # Multi-Provider Authentication
├── agent_manager.go          # Agent Lifecycle Management
├── task_manager_provider.go  # Task Manager Abstraktion
├── processor/               # Message Processing Framework
│   ├── base.go             # Base Processor Implementation
│   └── echo.go             # Echo Processor (Referenz)
├── example_test.go          # Vollständige Usage Examples
└── README.md               # Diese Dokumentation
```

## 🚀 Quick Start

### Basic Server Setup

```go
import (
    "github.com/denkhaus/agents/server"
    "github.com/denkhaus/agents/server/processor"
    "github.com/samber/mo"
)

// Server mit Options erstellen
unifiedServer, err := server.NewServer("localhost",
    server.WithPort(8080),
    server.WithTimeout(30*time.Second),
    server.WithCORS(true))

// Echo Processor erstellen
echoProcessor := processor.NewEchoProcessor()

// Agent Card definieren
agentCard := trpcserver.AgentCard{
    Name:        "EchoAgent",
    Description: "A simple echo agent",
    Capabilities: trpcserver.AgentCapabilities{
        Streaming: &[]bool{true}[0],
    },
    DefaultInputModes:  []string{"text"},
    DefaultOutputModes: []string{"text"},
}

// Agent hinzufügen
err = unifiedServer.AddAgent("echo", agentCard, echoProcessor)

// Server starten
err = unifiedServer.Start("localhost:8080")
```

### Advanced Configuration

```go
// Authentication Manager mit Multi-Provider
authManager := server.NewAuthenticationManager()

// JWT Provider
jwtProvider := auth.NewJWTProvider(auth.JWTConfig{
    Secret:   []byte("secret-key"),
    Audience: "a2a-server",
    Issuer:   "my-server",
    Expiry:   time.Hour,
})
authManager.AddProvider(jwtProvider)

// API Key Provider
apiKeyProvider := auth.NewAPIKeyProvider(auth.APIKeyConfig{
    Key:    "my-api-key",
    Header: mo.Some("X-API-Key"),
})
authManager.AddProvider(apiKeyProvider)

// Server mit Auth erstellen
config := server.ServerConfig{
    Host: "localhost",
    Port: mo.Some(8080),
    Auth: mo.Some(authManager),
    TaskManagerProvider: mo.Some(server.NewRedisTaskManagerProvider(
        server.RedisConfig{
            Addr: "localhost:6379",
            DB:   0,
        })),
    EnableMetrics: mo.Some(true),
    EnableCORS:    mo.Some(true),
}

unifiedServer, err := server.NewUnifiedServer(config)
```

## 🔧 Wiederverwendete Components

### **Aus `auth/` Package**
- ✅ **Provider Interface**: Einheitliche Auth-Abstraktion
- ✅ **JWT/API Key/OAuth2**: Vollständige Provider-Implementierungen
- ✅ **Type-safe Config**: `mo.Option` für optionale Parameter

### **Aus `helper/` Package**
- ✅ **Error Types**: `HTTPError`, `AuthError`, `ValidationError`
- ✅ **Request Handling**: HTTP Request/Response Utilities
- ✅ **Utility Functions**: `StringPtr`, `BoolPtr`, etc.

## 🎛 Interfaces

### **UnifiedServer**
```go
type UnifiedServer interface {
    Start(addr string) error
    Stop(ctx context.Context) error
    AddAgent(name string, card server.AgentCard, processor MessageProcessor) error
    RemoveAgent(name string) error
    GetAgents() map[string]AgentInfo
}
```

### **MessageProcessor**
```go
type MessageProcessor interface {
    taskmanager.MessageProcessor
    GetCapabilities() ProcessorCapabilities
    SupportsStreaming() bool
}
```

### **AuthenticationManager**
```go
type AuthenticationManager interface {
    AddProvider(provider auth.Provider) error
    RemoveProvider(providerType string) error
    Authenticate(req *http.Request) (*AuthContext, error)
    GetProviders() []auth.Provider
}
```

## 🧪 Testing

```bash
cd server
go test -v .
```

## 🔄 Nächste Schritte (Tier 2)

- **Push Notification System**: JWKS-signed notifications
- **Multi-Agent Routing**: Chi Router integration
- **Advanced Streaming**: Chunked processing
- **Configuration Management**: File + ENV support

## 🎯 Design Principles

- **Interface-Driven**: Alle Komponenten über Interfaces abstrahiert
- **Modular**: Jede Funktionalität als separates Subpackage
- **Type-Safe**: `samber/mo` für optionale Konfiguration
- **Production-Ready**: Thread-safe, error-resilient, observable
- **Extensible**: Plugin-basierte Erweiterungen

Das Server-Package ist **Tier 1 complete** und bietet eine solide Basis für production-ready A2A Server mit modernem Go-Design! 🚀