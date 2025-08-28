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
	sessionID uuid.UUID
	runner    runner.Runner
	wrapper   shared.TheAgent
}

// ID returns the unique identifier of the agent.
func (p *AgentRunner) ID() uuid.UUID {
	return p.wrapper.ID()
}

// Name returns the name of the agent.
func (p *AgentRunner) Name() string {
	return p.wrapper.Info().Name
}

// SessionID returns the session ID as a string.
func (p *AgentRunner) SessionID() string {
	return p.sessionID.String()
}

// Info returns the agent's information structure.
func (p *AgentRunner) Info() *shared.AgentInfo {
	return p.wrapper.GetInfo()
}

// String returns a string representation of the agent runner.
func (p *AgentRunner) String() string {
	return fmt.Sprintf("%s-[%s]", p.wrapper.Info().Name, p.wrapper.ID())
}

// Run executes the agent with a message from another agent and returns a channel of events.
// The fromAgentID identifies the sender, userMessage contains the message content,
// and runOpts provides additional configuration options.
func (p *AgentRunner) Run(
	ctx context.Context,
	fromAgentID uuid.UUID,
	userMessage model.Message,
	runOpts ...agent.RunOption,
) (<-chan *event.Event, error) {

	return p.runner.Run(ctx, fromAgentID.String(), p.SessionID(), userMessage, runOpts...)
}
