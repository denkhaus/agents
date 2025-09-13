package multi

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/denkhaus/agents/pkg/shared"
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

// ID returns the unique identifier of the agent.
func (p *AgentRunner) ID() uuid.UUID {
	return p.wrapper.GetID()
}

// Name returns the name of the agent.
func (p *AgentRunner) Name() string {
	return p.wrapper.Info().Name
}

// Info returns the agent's information structure.
func (p *AgentRunner) Info() *shared.AgentInfo {
	return p.wrapper.GetInfo()
}

func (p *AgentRunner) IsStreaming() bool {
	return p.wrapper.GetIsStreaming()
}

// String returns a string representation of the agent runner.
func (p *AgentRunner) String() string {
	return fmt.Sprintf("%s-[%s]", p.wrapper.Info().Name, p.wrapper.GetID())
}

// Run executes the agent with a message from another agent and returns a channel of events.
// The fromAgentID identifies the sender, userMessage contains the message content,
// and runOpts provides additional configuration options.
func (p *AgentRunner) Run(
	ctx context.Context,
	routingInfo *messaging.RoutingInfo,
	userMessage model.Message,
	runOpts ...agent.RunOption,
) (<-chan *event.Event, error) {

	state := agent.WithRuntimeState(map[string]interface{}{
		"from_agent_id": routingInfo.FromAgentID,
		"to_agent_id":   routingInfo.ToAgentID,
		"session_id":    routingInfo.SessionID,
		"streaming":     p.IsStreaming(),
	})

	return p.runner.Run(
		ctx,
		routingInfo.FromAgentID.String(),
		routingInfo.SessionID.String(),
		userMessage,
		state,
	)
}
