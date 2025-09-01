package tools

import (
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/mitchellh/mapstructure"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

type ConfigPayload map[string]interface{}

// Bind binds the configuration payload to the provided struct using mapstructure
func (c ConfigPayload) Bind(target interface{}) error {
	return mapstructure.Decode(c, target)
}

type ToolFactoryFunc func(config ConfigPayload, availableAgents []*shared.AgentInfo) (tool.Tool, error)
type ToolSetFactoryFunc func(config ConfigPayload, availableAgents []*shared.AgentInfo) (tool.ToolSet, error)
