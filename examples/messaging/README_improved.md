# 🤖 Advanced Multi-Agent Chat System

Ein verbessertes Multi-Agent-Chat-System, das Menschen als vollwertige Agents behandelt und dynamische Agent-Registrierung ermöglicht.

## 🎯 Hauptverbesserungen

### ✅ Gelöste Probleme
1. **Agent-ID Reihenfolge**: Agents werden erst erstellt, dann mit korrekten IDs konfiguriert
2. **Dynamische Agent-Kenntnisse**: Alle Agents kennen automatisch alle anderen Agents
3. **Mensch als Agent**: Menschen sind vollwertige Agents mit eigenen IDs
4. **Zentrale Registry**: Alle Agents werden zentral verwaltet
5. **Wiederverwendbare Komponenten**: Modulare, wiederverwendbare Funktionen

### 🏗️ Architektur-Komponenten

#### 1. **AgentRegistry** (`agent_registry.go`)
- Zentrale Verwaltung aller Agents (AI + Human)
- Automatische Agent-Kenntnisse-Generierung
- Message-Routing zwischen Agents
- Thread-sichere Operationen

#### 2. **MessageBroker** (`message_broker.go`)
- Einfacher In-Memory Message Broker
- Pub/Sub Pattern für Agent-Kommunikation
- Asynchrone Message-Verarbeitung

#### 3. **Convenience Functions** (`convenience_functions.go`)
- `CreateAIAgentWithMessaging()`: Erstellt AI-Agents mit Messaging
- `CreateHumanAgent()`: Registriert Menschen als Agents
- `StartMultiAgentChat()`: Startet interaktive Chat-Session
- `ListAllAgents()`: Zeigt alle registrierten Agents

#### 4. **Improved Main** (`main_improved.go`)
- Demonstriert die Verwendung aller Komponenten
- Interaktive Chat-Oberfläche
- Beispiele für Agent-zu-Agent Kommunikation

## 🚀 Verwendung

### Basis-Setup
```go
// Message Broker und Registry erstellen
broker := NewSimpleBroker()
registry := NewAgentRegistry(broker)
config := DefaultAgentConfig()

// AI-Agent erstellen
agentID, err := CreateAIAgentWithMessaging(
    registry,
    config,
    "CodeMaster",
    "Expert software engineer",
    "You are a skilled programmer...",
)

// Human-Agent registrieren
humanID := CreateHumanAgent(registry, "Developer", "Human developer")
```

### Chat-Kommandos
```bash
# Nachricht an Agent senden
@coder Please help me write a sorting function

# Alle Agents auflisten
/list

# Eigene Agent-Info anzeigen
/who

# Hilfe anzeigen
/help

# Chat beenden
/exit
```

### Agent-zu-Agent Kommunikation
```go
// Agents können sich direkt Nachrichten senden
err := registry.SendMessage(ctx, coderID, reviewerID, "Please review this code...")
```

## 🔧 Konfiguration

### Agent-Erstellung Konfiguration
```go
config := AgentCreationConfig{
    AppName:     "multi-agent-chat",
    ModelName:   "deepseek-chat",
    MaxTokens:   500,
    Temperature: 0.7,
}
```

### Automatische Agent-Kenntnisse
Jeder Agent erhält automatisch Informationen über alle anderen Agents:
```
Known Agents in the System:
- CodeMaster (ID: 123e4567-e89b-12d3-a456-426614174000): Expert software engineer [Type: ai]
- CodeReviewer (ID: 123e4567-e89b-12d3-a456-426614174001): Expert code reviewer [Type: ai]
- Developer (ID: 123e4567-e89b-12d3-a456-426614174002): Human developer [Type: human]

You can send messages to any of these agents using the send_message tool with their ID.
```

## 🛠️ Erweiterte Features

### 1. **Dynamische Agent-Registrierung**
```go
// Neue Agents können zur Laufzeit hinzugefügt werden
newAgentID := registry.RegisterAIAgent(name, description, runner, agent)

// Agents können entfernt werden
registry.UnregisterAgent(agentID)
```

### 2. **Message-Tool für Agents**
Alle AI-Agents haben automatisch Zugriff auf das `send_message` Tool:
```json
{
  "name": "send_message",
  "description": "Send a message to another agent in the system",
  "parameters": {
    "to_agent_id": "UUID of the target agent",
    "message": "Message content to send"
  }
}
```

### 3. **Human-Agent Integration**
Menschen erhalten Nachrichten über einen dedizierten Channel:
```go
messageChan, exists := registry.GetHumanMessageChannel(humanID)
if exists {
    for message := range messageChan {
        fmt.Println("Received:", message)
    }
}
```

## 📊 Beispiel-Interaktion

```
🤖 Advanced Multi-Agent Chat System
====================================
🔧 Initializing agents...
✅ System initialized successfully!

Registered Agents:
==================
- CodeMaster (ai)
  ID: 123e4567-e89b-12d3-a456-426614174000
  Description: Expert software engineer specializing in clean, efficient code
  Status: online

- CodeReviewer (ai)
  ID: 123e4567-e89b-12d3-a456-426614174001
  Description: Expert code reviewer focused on quality, security, and best practices
  Status: online

- Developer (human)
  ID: 123e4567-e89b-12d3-a456-426614174002
  Description: Human developer and system operator
  Status: online

💬 Starting interactive chat...

you> @coder Please write a function to reverse a string in Go

📤 Message sent to CodeMaster

[CodeMaster]: Here's a simple and efficient function to reverse a string in Go:

```go
func reverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
```

This function properly handles Unicode characters by converting the string to runes first.

you> @reviewer Please review the code that CodeMaster just provided

📤 Message sent to CodeReviewer

[CodeReviewer]: I've reviewed the string reversal function. Here's my analysis:

**Strengths:**
- ✅ Correctly handles Unicode characters using runes
- ✅ Efficient in-place reversal algorithm
- ✅ Clean, readable code structure

**Suggestions:**
- Consider adding input validation for empty strings
- Could benefit from documentation comments
- Might want to add unit tests

Overall, this is a solid implementation that follows Go best practices!
```

## 🔄 Erweiterungsmöglichkeiten

1. **Persistente Message-Historie**
2. **Web-Interface für Chat**
3. **Agent-Gruppen und Channels**
4. **Message-Verschlüsselung**
5. **Agent-Performance-Metriken**
6. **Externe Message-Broker Integration (Redis, RabbitMQ)**

## 🧪 Testing

```bash
# System starten
go run examples/messaging/*.go

# Umgebungsvariable setzen
export OPENAI_API_KEY="your-api-key"
```

Das verbesserte System löst alle identifizierten Probleme und bietet eine solide Basis für erweiterte Multi-Agent-Kommunikation!