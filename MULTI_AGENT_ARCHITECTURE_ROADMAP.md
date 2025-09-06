# Multi-Agent Architecture Roadmap

## Aktueller Status: SCHRITT 1 ✅ ABGESCHLOSSEN

**Frontend Agent Store Refactoring** wurde erfolgreich implementiert. Das Frontend behandelt jetzt den Human als vollwertigen Agent mit UUID statt String-Referenzen.

### Was bereits implementiert ist:
- ✅ Agent-ID Konstanten (`AGENT_IDS.HUMAN` etc.)
- ✅ Type-System komplett überarbeitet (alle `AgentId` statt `string`)
- ✅ Chat Store refactoriert (keine "user" Strings mehr)
- ✅ UI Components angepasst (MessageInput, MessageItem)
- ✅ Backend-Frontend ID-Synchronisation

---

## SCHRITT 2: Unified Message Flow 🎯 NÄCHSTER SCHRITT

### Ziel:
Alle Messages (Human + Agent) über einheitliche API-Kanäle verarbeiten, statt separater User/Agent Pfade.

### Aktuelle Probleme:
1. **Doppelte Message-Pfade**: Human Messages über Chat Store, Agent Messages über SSE
2. **Inkonsistente Session-Behandlung**: Human Sessions anders als Agent Sessions
3. **Fragmentierte Message-History**: Verschiedene Speicher-Mechanismen

### Zu implementieren:

#### 2.1 Backend: Unified Session Management
```go
// pkg/multi/plugins/web/server.go
// Erweitere handleRunSSE für alle Agent-Typen (inkl. Human)
func (s *Server) handleRunSSE(w http.ResponseWriter, r *http.Request) {
    // Behandle Human-Agent genauso wie AI-Agents
    // Einheitliche Session-Erstellung und -Verwaltung
}
```

#### 2.2 Frontend: Einheitlicher Message-Service
```typescript
// docker/react-chat/src/lib/services/message-service.ts
class UnifiedMessageService {
    // Alle Messages über SSE senden (Human + Agent)
    sendMessage(fromAgent: AgentId, toAgent: AgentId, content: string)
    
    // Einheitliche Message-Verarbeitung
    processIncomingMessage(event: AgentEvent)
}
```

#### 2.3 Chat Store Vereinfachung
```typescript
// Entferne doppelte Message-Pfade
// Alle Messages über einheitlichen Flow
addMessage() // Nur für lokale UI-Updates
sendMessage() // Für alle Backend-Kommunikation
```

### Dateien zu ändern:
- `pkg/multi/plugins/web/server.go`
- `docker/react-chat/src/lib/store/chat-store.ts`
- `docker/react-chat/src/components/messaging/message-streaming.tsx`
- `docker/react-chat/src/hooks/use-streaming-manager.tsx`

---

## SCHRITT 3: Session Management Harmonisierung

### Ziel:
Human-Agent Sessions wie AI-Agent Sessions behandeln - einheitliche Erstellung, Verwaltung und Persistierung.

### Zu implementieren:

#### 3.1 Backend: Session-Service Erweiterung
```go
// Einheitliche Session-Behandlung für alle Agent-Typen
func (s *SessionService) CreateAgentSession(agentId uuid.UUID, userId string) (*Session, error)
func (s *SessionService) GetAgentSessions(agentId uuid.UUID, userId string) ([]*Session, error)
```

#### 3.2 Frontend: Session-Manager Refactoring
```typescript
// Einheitlicher Session-Manager für alle Agents
class SessionManager {
    createSession(agentId: AgentId): Promise<Session>
    loadSession(agentId: AgentId, sessionId: string): Promise<Session>
    switchSession(agentId: AgentId, sessionId: string): Promise<void>
}
```

### Dateien zu ändern:
- `pkg/multi/plugins/web/server.go`
- `docker/react-chat/src/lib/store/chat-store.ts`
- `docker/react-chat/src/components/navigation/agent-navigation-bar.tsx`

---

## SCHRITT 4: Inter-Agent Communication Standardisierung

### Ziel:
Human als vollwertiger Agent in Inter-Agent-Kommunikation integrieren.

### Zu implementieren:

#### 4.1 Backend: Human-Agent Integration
```go
// pkg/multi/plugins/web/multi_chat.go
// Human kann Inter-Agent Messages senden/empfangen
func (s *Server) handleHumanToAgentMessage(fromAgent, toAgent uuid.UUID, message string)
func (s *Server) handleAgentToHumanMessage(fromAgent, toAgent uuid.UUID, message string)
```

#### 4.2 Frontend: Inter-Agent UI für Human
```typescript
// Human kann andere Agents direkt ansprechen
// Inter-Agent Events zeigen Human-Beteiligung
// Unified Event-Display für alle Agent-Typen
```

### Dateien zu ändern:
- `pkg/multi/plugins/web/multi_chat.go`
- `docker/react-chat/src/components/inter-agent/`
- `docker/react-chat/src/lib/streaming/processors/`

---

## SCHRITT 5: Agent Registry & Discovery

### Ziel:
Dynamische Agent-Erkennung und -Verwaltung statt hardcodierter Agent-Listen.

### Zu implementieren:

#### 5.1 Backend: Dynamic Agent Registry
```go
// pkg/shared/registry.go
type AgentRegistry interface {
    RegisterAgent(id uuid.UUID, metadata AgentMetadata) error
    GetAvailableAgents() []Agent
    GetAgentCapabilities(id uuid.UUID) []string
}
```

#### 5.2 Frontend: Dynamic Agent Loading
```typescript
// Agents dynamisch vom Backend laden
// Agent-Capabilities anzeigen
// Agent-Status live tracking
```

### Dateien zu ändern:
- `pkg/shared/` (neue Registry)
- `pkg/multi/plugins/web/server.go`
- `docker/react-chat/src/lib/store/agents-store.ts`

---

## SCHRITT 6: Advanced Features

### 6.1 Agent-to-Agent Direct Communication
- Direkte Agent-Kommunikation ohne Human-Intervention
- Agent-Workflows und -Chains
- Automated Task Distribution

### 6.2 Enhanced Session Management
- Multi-Agent Sessions (mehrere Agents in einer Session)
- Session-Sharing zwischen Agents
- Session-Templates und -Workflows

### 6.3 Advanced UI Features
- Multi-Agent Chat Views
- Agent-Performance Monitoring
- Real-time Agent-Status Dashboard

---

## Prioritäten und Zeitschätzung

### Sofort (Diese Woche):
- **SCHRITT 2**: Unified Message Flow (2-3 Tage)

### Kurzfristig (Nächste 2 Wochen):
- **SCHRITT 3**: Session Management Harmonisierung (3-4 Tage)
- **SCHRITT 4**: Inter-Agent Communication Standardisierung (2-3 Tage)

### Mittelfristig (Nächster Monat):
- **SCHRITT 5**: Agent Registry & Discovery (1 Woche)
- **SCHRITT 6**: Advanced Features (2-3 Wochen)

---

## Erfolgskriterien

### Nach SCHRITT 2:
- ✅ Alle Messages über einheitliche API
- ✅ Keine doppelten Message-Pfade
- ✅ Konsistente Message-History

### Nach SCHRITT 3:
- ✅ Einheitliche Session-Verwaltung
- ✅ Human-Sessions wie Agent-Sessions
- ✅ Nahtloser Agent-Wechsel

### Nach SCHRITT 4:
- ✅ Human in Inter-Agent-Kommunikation
- ✅ Vollständige Agent-Gleichberechtigung
- ✅ Unified Event-System

### Finale Vision:
- 🎯 **Vollständig einheitliche Multi-Agent-Architektur**
- 🎯 **Human als gleichberechtigter Agent**
- 🎯 **Saubere, erweiterbare Codebasis**
- 🎯 **Keine Workarounds oder Legacy-Code**

---

## Nächste Aktion

**Beginne mit SCHRITT 2: Unified Message Flow**

1. Analysiere aktuelle Message-Pfade
2. Implementiere einheitlichen Message-Service
3. Refactore Chat Store für einheitlichen Flow
4. Teste Human-Agent Message-Gleichberechtigung

**Ziel: "Fail fast, fail loud" - Saubere Architektur ohne Kompromisse**