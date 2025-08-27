# 🤖 Multi-Agent Chat System (Clean Version)

Eine saubere, minimale Implementation des Multi-Agent-Chat-Systems, die die bestehende `MessagingWrapper` Infrastruktur nutzt.

## 🎯 **Hauptmerkmale**

- ✅ Nutzt bestehende `messaging.MessagingWrapper` und `messaging.MessageBroker`
- ✅ Keine unbenutzten Variablen oder redundanter Code
- ✅ Einfache, klare Architektur
- ✅ Automatische Agent-zu-Agent Kommunikation über `send_message` Tool
- ✅ Interaktive Chat-Oberfläche

## 🏗️ **Architektur**

```go
ChatSystem
├── MessageBroker (bestehend)    # Message-Routing zwischen Agents
├── AgentRunner                  # AI Agent + MessagingWrapper + Runner
└── Agent-Registry (map)         # Name -> AgentRunner Mapping
```

## 🚀 **Verwendung**

### Setup
```bash
export OPENAI_API_KEY="your-api-key"
go run examples/messaging/main_clean.go
```

### Chat-Kommandos
```bash
# Nachricht an Agent senden
you> @coder Please write a function to sort an array

# Agents auflisten
you> /list

# Chat beenden
you> /exit
```

## 🔧 **Code-Struktur**

### **ChatSystem** (Hauptklasse)
```go
type ChatSystem struct {
    broker *messaging.MessageBroker        // Bestehende Message-Infrastruktur
    agents map[string]*AgentRunner         // Agent-Registry
}
```

### **AgentRunner** (Vereinfacht)
```go
type AgentRunner struct {
    Runner  runner.Runner                  // Agent-Runner
    Wrapper *messaging.MessagingWrapper    // Messaging-Fähigkeiten
    Name    string                         // Agent-Name
}
```

### **Agent-Erstellung** (Eine Zeile)
```go
// Basis-Agent erstellen
baseAgent := llmagent.New(name, options...)

// Mit Messaging umhüllen (automatische Registrierung)
wrapper := messaging.NewMessagingWrapper(baseAgent, broker)

// Runner erstellen
runner := runner.NewRunner(appName, wrapper, options...)
```

## 💬 **Beispiel-Interaktion**

```
🤖 Multi-Agent Chat System
==========================
Creating agents...
✅ Agents created successfully!

Available agents: Coder, Reviewer

Commands:
  @<agent> <message>  - Send message to agent
  /list              - List agents
  /exit              - Exit

you> @coder Please write a bubble sort function in Go

[Coder]: Here's a bubble sort implementation in Go:

```go
func bubbleSort(arr []int) []int {
    n := len(arr)
    for i := 0; i < n; i++ {
        for j := 0; j < n-i-1; j++ {
            if arr[j] > arr[j+1] {
                arr[j], arr[j+1] = arr[j+1], arr[j]
            }
        }
    }
    return arr
}
```

you> @reviewer Please review the code above

[Reviewer] using tool: send_message

[Reviewer]: I'll review the bubble sort code. The implementation is correct but could be optimized...
```

## 🔄 **Agent-zu-Agent Kommunikation**

Agents können sich automatisch über das `send_message` Tool unterhalten:

```go
// Automatisch verfügbar für alle Agents:
{
  "name": "send_message",
  "parameters": {
    "to": "agent-uuid",
    "content": "message content"
  }
}
```

## 📊 **Vorteile der sauberen Version**

1. **🧹 Minimal**: Nur 150 Zeilen Code, keine unbenutzten Variablen
2. **🔄 Wiederverwendung**: Nutzt bestehende Infrastruktur vollständig
3. **📖 Lesbar**: Klare, einfache Struktur
4. **🐛 Fehlerfrei**: Keine redundanten Agent-Erstellungen
5. **⚡ Effizient**: Direkte Nutzung der MessagingWrapper-Features

## 🛠️ **Erweiterungsmöglichkeiten**

- **Human-Agent Integration**: Menschen als vollwertige Agents hinzufügen
- **Persistenz**: Message-Historie speichern
- **Web-Interface**: Browser-basierte Chat-Oberfläche
- **Agent-Gruppen**: Thematische Agent-Gruppierungen

Die saubere Version konzentriert sich auf das Wesentliche und nutzt die bestehende, bewährte Infrastruktur optimal!