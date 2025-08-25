package settings

import (
	"context"

	"github.com/denkhaus/agents/agents"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
)

type ModelSettings struct {
	ChannelBufferSize int    `yaml:"channel_buffer_size"`
	Name              string `yaml:"model_name"`
}

type AgentSettings struct {
	StreamingEnabled  bool             `yaml:"streaming_enabled"`
	Temperature       float64          `yaml:"temperature"`
	ChannelBufferSize int              `yaml:"channel_buffer_size"`
	MaxTokens         int              `yaml:"max_tokens"`
	Role              agents.AgentRole `yaml:"role"`
	Name              string           `yaml:"name"`
	ID                uuid.UUID        `yaml:"id"`
}

type Settings struct {
	Model ModelSettings `yaml:"model"`
	Agent AgentSettings `yaml:"agent"`
}

type AgentConfiguration interface {
	GetAgentName() string
	GetOptions(ctx context.Context) ([]llmagent.Option, error)
}
