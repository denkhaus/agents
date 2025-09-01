package config

import (
	"fmt"

	"github.com/denkhaus/agents/pkg/shared"
	"github.com/denkhaus/agents/pkg/tools"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// cueToolFactoryImpl creates tools from CUE-based configuration
type cueToolFactoryImpl struct {
	injector *do.Injector
}

// NewCUEToolFactory creates a new CUE-based tool factory
func NewCUEToolFactory(injector *do.Injector) (ToolFactory, error) {
	return &cueToolFactoryImpl{
		injector: injector,
	}, nil
}

// CreateTools creates tools and toolsets from CUE configuration
func (f *cueToolFactoryImpl) CreateTools(
	toolsConfig ToolsConfig,
	availableAgents []*shared.AgentInfo,
) ([]tool.Tool, []tool.ToolSet, error) {
	var tools []tool.Tool
	var toolsets []tool.ToolSet

	// Create individual tools
	for toolName, toolConfig := range toolsConfig.Tools {
		if !toolConfig.Enabled {
			continue
		}

		// Create the tool using the registered factory
		tool, err := f.createTool(toolName, toolConfig.Config, availableAgents)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create tool %s: %w", toolName, err)
		}

		tools = append(tools, tool)
	}

	// Create toolsets
	for toolsetName, toolsetConfig := range toolsConfig.ToolSets {
		if !toolsetConfig.Enabled {
			continue
		}

		// Create the toolset using the registered factory
		toolset, err := f.createToolSet(toolsetName, toolsetConfig.Config, availableAgents)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create toolset %s: %w", toolsetName, err)
		}

		toolsets = append(toolsets, toolset)
	}

	return tools, toolsets, nil
}

// createTool creates a single tool using the tool provider's factories
func (f *cueToolFactoryImpl) createTool(
	toolName string,
	config tools.ConfigPayload,
	availableAgents []*shared.AgentInfo,
) (tool.Tool, error) {

	factory, err := do.InvokeNamed[tools.ToolFactoryFunc](f.injector, toolName)
	if err != nil {
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}

	return factory(config, availableAgents)
}

// createToolSet creates a single toolset using the tool provider's factories
func (f *cueToolFactoryImpl) createToolSet(
	toolsetName string,
	config tools.ConfigPayload,
	availableAgents []*shared.AgentInfo,
) (tool.ToolSet, error) {

	factory, err := do.InvokeNamed[tools.ToolSetFactoryFunc](f.injector, toolsetName)
	if err != nil {
		return nil, fmt.Errorf("unknown toolset: %s", toolsetName)
	}

	return factory(config, availableAgents)
}
