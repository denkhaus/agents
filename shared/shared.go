package shared

import (
	"fmt"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
)

const (
	ContextKeyToolInfo  = "tool_info"
	ContextKeyAgentInfo = "agent_info"
)

var (
	AgentInfoHuman = NewAgentInfo(AgentIDHuman, "Denkhaus-<human>", "A human you can chat with")
)

type ToolInfo struct {
	Name        string
	Description string
}

type AgentInfo struct {
	agent.Info
	ID uuid.UUID
}

func (p *AgentInfo) String() string {
	return fmt.Sprintf("%s-[%s]", p.Name, p.ID)
}

func NewAgentInfo(agentID uuid.UUID, name, description string) AgentInfo {
	return AgentInfo{
		ID: agentID,
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
}

type theAgentImpl struct {
	agent.Agent
	agentID     uuid.UUID
	isStreaming bool
}

func (p *theAgentImpl) ID() uuid.UUID {
	return p.agentID
}

func (p *theAgentImpl) IsStreaming() bool {
	return p.isStreaming
}

func (p *theAgentImpl) GetInfo() *AgentInfo {
	return &AgentInfo{
		Info: p.Agent.Info(),
		ID:   p.agentID,
	}
}

func NewAgent(agent agent.Agent, agentID uuid.UUID, isStreaming bool) TheAgent {
	return &theAgentImpl{
		Agent:       agent,
		agentID:     agentID,
		isStreaming: isStreaming,
	}
}

func TheAgentToInfo(agent TheAgent) *AgentInfo {
	return &AgentInfo{
		Info: agent.Info(),
		ID:   agent.ID(),
	}
}
