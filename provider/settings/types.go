package settings

import (
	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
)

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
	MaxIterations     int              `yaml:"max_iterations"`
	SubAgents         []uuid.UUID      `yaml:"sub_agents"`
	Role              shared.AgentRole `yaml:"role"`
	Type              shared.AgentType `yaml:"type"`
	Name              string           `yaml:"name"`
}

type Settings struct {
	AgentID uuid.UUID     `yaml:"agent_id"`
	Model   ModelSettings `yaml:"model"`
	Agent   AgentSettings `yaml:"agent"`
}
