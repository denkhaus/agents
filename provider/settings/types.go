package settings

import (
	"context"

	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
)

type Provider interface {
	GetAgentConfiguration(agentID uuid.UUID) (AgentConfiguration, error)
}

type ModelSettings struct {
	ChannelBufferSize int    `yaml:"channel_buffer_size"`
	Name              string `yaml:"model_name"`
}

type AgentSettings struct {
	PlanningEnabled   bool             `yaml:"planning_enabled"`
	StreamingEnabled  bool             `yaml:"streaming_enabled"`
	Temperature       float64          `yaml:"temperature"`
	ChannelBufferSize int              `yaml:"channel_buffer_size"`
	MaxTokens         int              `yaml:"max_tokens"`
	Role              shared.AgentRole `yaml:"role"`
	Name              string           `yaml:"name"`
}

type Settings struct {
	AgentID uuid.UUID     `yaml:"agent_id"`
	Model   ModelSettings `yaml:"model"`
	Agent   AgentSettings `yaml:"agent"`
}

type AgentConfiguration interface {
	GetAgentName() string
	IsStreamingEnabled() bool
	GetOptions(ctx context.Context) ([]llmagent.Option, error)
}
