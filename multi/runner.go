package multi

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/runner"
)

// AgentRunner represents an AI agent with messaging capabilities
type AgentRunner struct {
	runner  runner.Runner
	wrapper shared.TheAgent
}

func (p *AgentRunner) ID() uuid.UUID {
	return p.wrapper.ID()
}

func (p *AgentRunner) Name() string {
	return p.wrapper.Info().Name
}

func (p *AgentRunner) SessionID() string {
	return p.wrapper.Info().Name
}

func (p *AgentRunner) Info() *shared.AgentInfo {
	return shared.TheAgentToInfo(p.wrapper)
}

func (p *AgentRunner) String() string {
	return fmt.Sprintf("%s-[%s]", p.wrapper.Info().Name, p.wrapper.ID())
}

func (p *AgentRunner) Run(
	ctx context.Context,
	fromAgentID uuid.UUID,
	userMessage model.Message,
	runOpts ...agent.RunOption,
) (<-chan *event.Event, error) {

	return p.runner.Run(ctx, fromAgentID.String(), p.SessionID(), userMessage, runOpts...)
}
