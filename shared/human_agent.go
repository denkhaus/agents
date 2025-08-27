package shared

import (
	"context"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// humanAgentImpl implements the agent.Agent interface for the human
type humanAgentImpl struct {
	AgentInfo
}

func NewHumanAgent(info AgentInfo) TheAgent {
	return &humanAgentImpl{
		AgentInfo: info,
	}
}

func (d *humanAgentImpl) ID() uuid.UUID {
	return d.AgentInfo.ID
}

func (d *humanAgentImpl) Run(ctx context.Context, invocation *agent.Invocation) (<-chan *event.Event, error) {
	// Humans don't process messages automatically, just return empty channel
	ch := make(chan *event.Event)
	close(ch)
	return ch, nil
}

func (d *humanAgentImpl) Info() agent.Info {
	return d.AgentInfo.Info
}

func (d *humanAgentImpl) Tools() []tool.Tool {
	return []tool.Tool{}
}

func (d *humanAgentImpl) FindSubAgent(name string) agent.Agent {
	return nil
}

func (d *humanAgentImpl) SubAgents() []agent.Agent {
	return []agent.Agent{}
}
