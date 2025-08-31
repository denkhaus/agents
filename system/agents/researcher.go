package agents

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/pkg/config"
	"github.com/denkhaus/agents/pkg/provider"
	"github.com/denkhaus/agents/pkg/provider/agent"
	providerConfig "github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/shared"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
)

func CreateResearcherAgent(ctx context.Context, injector *do.Injector) (shared.TheAgent, error) {
	agentID := shared.AgentIDResearcher
	agentProvider := do.MustInvoke[provider.AgentProvider](injector)
	configProvider := do.MustInvoke[config.Service](injector)
	cueConfigProvider := do.MustInvoke[providerConfig.ConfigProvider](injector)

	// Load tool profile configuration
	toolsConfig, err := cueConfigProvider.LoadToolProfile("researcher")
	if err != nil {
		return nil, fmt.Errorf("failed to load researcher tool profile: %w", err)
	}

	// Get workspace path for relative path resolution
	workspacePath, err := configProvider.GetWorkspacePath()
	if err != nil {
		return nil, err
	}

	// Resolve workspace paths in tool configuration
	resolvedToolsConfig := *toolsConfig
	for toolSetName, toolSetConfig := range resolvedToolsConfig.ToolSets {
		if toolSetConfig.Config != nil {
			config := make(map[string]interface{})
			for k, v := range toolSetConfig.Config {
				if k == "workspace_path" && v == "./workspace" {
					config[k] = workspacePath
				} else if k == "base_dir" && v == "./workspace" {
					config[k] = workspacePath
				} else {
					config[k] = v
				}
			}
			toolSetConfig.Config = config
			resolvedToolsConfig.ToolSets[toolSetName] = toolSetConfig
		}
	}

	// Create tools and toolsets using the tool factory
	toolFactory := do.MustInvoke[providerConfig.ToolFactory](injector)
	enabledTools, enabledToolSets, err := toolFactory.CreateTools(resolvedToolsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create tools from configuration: %w", err)
	}

	agent, err := agentProvider.GetAgent(ctx, agentID,
		agent.WithLLMAgentOptions(
			llmagent.WithTools(enabledTools),
			llmagent.WithToolSets(enabledToolSets),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return agent, nil
}
