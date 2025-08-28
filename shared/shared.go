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
	AgentInfoHuman = NewAgentInfo(
		AgentIDHuman,
		AgentRoleHuman,
		false,
		"denkhaus",
		"A human you can chat with",
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

func (p *AgentInfo) String() string {
	return fmt.Sprintf("%s-[%s]", p.Name, p.id)
}

func (p *AgentInfo) ID() uuid.UUID {
	return p.id
}

func (p *AgentInfo) Role() AgentRole {
	return p.role
}
func (p *AgentInfo) Equal(info AgentInfo) bool {
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
