package shared

import (
	"fmt"

	"github.com/denkhaus/agents/pkg/utils"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
)

const (
	ContextKeyToolInfo  = "tool_info"
	ContextKeyAgentInfo = "agent_info"
)

var (
	AgentInfoHuman = &AgentInfo{
		ID:          AgentIDHuman,
		Role:        AgentRoleHuman,
		IsStreaming: utils.BoolPtr(false),
		Name:        "denkhaus",
		Description: `Senior team lead and decision authority with comprehensive project oversight.
		Has the highest rank in the team hierarchy and serves as the final decision maker when agents cannot reach consensus or need guidance.
		Possesses deep domain expertise, strategic vision, and complete visibility across all projects. Responsible for critical decisions, conflict resolution,
		resource allocation, and overall project direction. Other agents should escalate complex decisions, blockers, and strategic questions to this human leader.`,
	}
)

type AgentInfo struct {
	Name string `json:"name"`

	// Description is the description of the agent.
	Description string `json:"description"`

	// InputSchema is the input schema of the agent.
	InputSchema map[string]any `json:"input_schema,omitempty"`

	// OutputSchema is the output schema of the agent.
	OutputSchema map[string]any `json:"output_schema,omitempty"`

	// IsStreaming defines if the agent has streaming capabilities
	IsStreaming *bool `json:"is_streaming"`

	// ID is the agents uuid
	ID uuid.UUID `json:"id"`

	// Role is the agents role
	Role AgentRole `json:"role"`
}

func (p *AgentInfo) String() string {
	return fmt.Sprintf("%s-[%s]", p.Name, p.ID)
}

func (p *AgentInfo) Equal(info *AgentInfo) bool {
	return p.Role == info.Role &&
		p.ID == info.ID &&
		p.Name == info.Name
}

type TheAgent interface {
	agent.Agent
	GetIsStreaming() bool
	GetInfo() *AgentInfo
	GetRole() AgentRole
	GetID() uuid.UUID
}

type theAgentImpl struct {
	agent.Agent
	info AgentInfo
}

func (p *theAgentImpl) GetID() uuid.UUID {
	return p.info.ID
}

func (p *theAgentImpl) GetRole() AgentRole {
	return p.info.Role
}

func (p *theAgentImpl) GetIsStreaming() bool {
	if p.info.IsStreaming != nil {
		return *p.info.IsStreaming
	}

	return false
}

func (p *theAgentImpl) GetInfo() *AgentInfo {
	return &p.info
}

func NewAgent(agent agent.Agent, agentID uuid.UUID, isStreaming bool, role AgentRole) TheAgent {
	return &theAgentImpl{
		Agent: agent,
		info: AgentInfo{
			Role:        role,
			ID:          agentID,
			IsStreaming: utils.BoolPtr(isStreaming),
		},
	}
}
