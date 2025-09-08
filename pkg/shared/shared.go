package shared

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
)

const (
	ContextKeyToolInfo  = "tool_info"
	ContextKeyAgentInfo = "agent_info"
)

var (
	AgentInfoHuman = NewAgentInfo(
		AgentIDHuman,
		AgentRoleHuman,
		false,
		"denkhaus",
		`Senior team lead and decision authority with comprehensive project oversight.
		Has the highest rank in the team hierarchy and serves as the final decision maker when agents cannot reach consensus or need guidance.
		Possesses deep domain expertise, strategic vision, and complete visibility across all projects. Responsible for critical decisions, conflict resolution,
		resource allocation, and overall project direction. Other agents should escalate complex decisions, blockers, and strategic questions to this human leader.`,
	)
)

type ToolInfo struct {
	Name        string
	Description string
}

type AgentInfo struct {
	agent.Info
	isStreaming bool
	id          uuid.UUID
	role        AgentRole
}

// MarshalJSON implements the json.Marshaler interface for AgentInfo.
// It ensures that both the embedded agent.Info fields and the private fields
// (isStreaming, id, role) are correctly marshaled into JSON, respecting omitempty.
func (p AgentInfo) MarshalJSON() ([]byte, error) {
	// Define an anonymous struct that mirrors the desired JSON output.
	// This allows explicit control over which fields are marshaled,
	// including unexported fields and fields from embedded structs.
	aux := struct {
		// Fields from embedded agent.Info (assuming they are public and desirable in JSON)
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`

		// Unexported fields from AgentInfo, explicitly included with their JSON tags
		IsStreaming bool      `json:"is_streaming,omitempty"`
		ID          uuid.UUID `json:"id,omitempty"` // Renamed to ID for consistency with JSON tag
		Role        AgentRole `json:"role,omitempty"`
	}{
		Name:        p.Info.Name,
		Description: p.Info.Description,
		IsStreaming: p.isStreaming,
		ID:          p.id,
		Role:        p.role,
	}

	// Use json.Marshal to convert the auxiliary struct to JSON.
	return json.Marshal(aux)
}

// UnmarshalJSON implements the json.Unmarshaler interface for AgentInfo.
// It allows for correct deserialization of the AgentInfo struct,
// including its embedded agent.Info and unexported fields.
func (p *AgentInfo) UnmarshalJSON(data []byte) error {
	// Define an anonymous struct for unmarshaling, mirroring the JSON structure
	// and allowing access to unexported fields.
	aux := struct {
		Name        string    `json:"name,omitempty"`
		Description string    `json:"description,omitempty"`
		IsStreaming bool      `json:"is_streaming,omitempty"`
		ID          uuid.UUID `json:"id,omitempty"`
		Role        AgentRole `json:"role,omitempty"`
	}{}

	// Unmarshal the JSON data into the auxiliary struct.
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Assign the unmarshaled values back to the AgentInfo struct.
	p.Info.Name = aux.Name
	p.Info.Description = aux.Description
	p.isStreaming = aux.IsStreaming
	p.id = aux.ID
	p.role = aux.Role

	return nil
}

func (p *AgentInfo) String() string {
	return fmt.Sprintf("%s-[%s]", p.Name, p.id)
}

func (p *AgentInfo) ID() uuid.UUID {
	return p.id
}

func (p *AgentInfo) SetID(agentID uuid.UUID) {
	p.id = agentID
}

func (p *AgentInfo) Role() AgentRole {
	return p.role
}

func (p *AgentInfo) SetRole(role AgentRole) {
	p.role = role
}

func (p *AgentInfo) Equal(info *AgentInfo) bool {
	return p.role == info.role &&
		p.id == info.id &&
		p.Name == info.Name
}

func (p *AgentInfo) IsStreaming() bool {
	return p.isStreaming
}

func NewAgentInfo(
	agentID uuid.UUID,
	role AgentRole,
	isStreaming bool,
	name string,
	description string,
) AgentInfo {
	return AgentInfo{
		id:          agentID,
		role:        role,
		isStreaming: isStreaming,
		Info: agent.Info{
			Name:        name,
			Description: description,
		},
	}
}

type TheAgent interface {
	agent.Agent
	ID() uuid.UUID
	IsStreaming() bool
	GetInfo() *AgentInfo
	GetRole() AgentRole
}

type theAgentImpl struct {
	agent.Agent
	role        AgentRole
	id          uuid.UUID
	isStreaming bool
}

func (p *theAgentImpl) ID() uuid.UUID {
	return p.id
}

func (p *theAgentImpl) GetRole() AgentRole {
	return p.role
}

func (p *theAgentImpl) IsStreaming() bool {
	return p.isStreaming
}

func (p *theAgentImpl) GetInfo() *AgentInfo {
	return &AgentInfo{
		Info: p.Agent.Info(),
		id:   p.id,
	}
}

func NewAgent(agent agent.Agent, agentID uuid.UUID, isStreaming bool) TheAgent {
	return &theAgentImpl{
		Agent:       agent,
		id:          agentID,
		isStreaming: isStreaming,
	}
}

func TheAgentToInfo(agent TheAgent) *AgentInfo {
	return &AgentInfo{
		Info: agent.Info(),
		id:   agent.ID(),
		role: agent.GetRole(),
	}
}
