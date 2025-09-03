# Design Rationale: Warum keine separate NewStreamingClient Funktion?

## Problem
Ursprünglich hatte das Design sowohl `NewAPIClient` als auch `NewStreamingClient`, was zu Redundanz führte:

```go
// Redundant - zwei Wege für dasselbe Ziel
apiClient := NewAPIClient(url, WithStreaming(true))
streamClient := NewStreamingClient(url) // macht dasselbe
```

## Lösung: Einheitliche Factory mit Options Pattern

### Vorteile der aktuellen Lösung:

1. **Konsistenz**: Ein einziger Einstiegspunkt für alle Client-Typen
2. **Flexibilität**: Streaming kann zur Laufzeit aktiviert/deaktiviert werden
3. **Klarheit**: Explizite Konfiguration über Options
4. **Erweiterbarkeit**: Neue Features können einfach als Options hinzugefügt werden

### Verwendung:

```go
// Standard Client
client := NewAPIClient("http://localhost:8080")

// Client mit Streaming
client := NewAPIClient("http://localhost:8080", 
    WithStreaming(true))

// Client mit allen Features
client := NewAPIClient("http://localhost:8080",
    WithAuth(auth),
    WithStreaming(true),
    WithTimeout(30*time.Second),
    WithRetryCount(3))
```

### Type Assertion für Streaming:

```go
// Sichere Type Assertion für Streaming-Features
if streamClient, ok := client.(StreamingClient); ok {
    eventChan, err := streamClient.StreamMessage(ctx, params)
    // ...
}
```

### Alternative: Unified Client für alle Features

```go
// Wenn alle Interfaces benötigt werden
unifiedClient := NewUnifiedClient(ClientConfig{
    BaseURL: "http://localhost:8080",
    Streaming: mo.Some(true),
    // ...
})
```

## Design Principles

1. **Single Responsibility**: Jede Factory-Funktion hat einen klaren Zweck
2. **Open/Closed**: Erweiterbar durch neue Options, geschlossen für Modifikation
3. **Interface Segregation**: Clients implementieren nur benötigte Interfaces
4. **Dependency Inversion**: Konfiguration durch Injection, nicht durch separate Factories

Diese Lösung eliminiert Redundanz und bietet maximale Flexibilität bei klarer API-Struktur.