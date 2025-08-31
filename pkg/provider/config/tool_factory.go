package config

import (
	"fmt"
	"time"

	"github.com/denkhaus/agents/pkg/tools/calculator"
	"github.com/denkhaus/agents/pkg/tools/fetch"
	"github.com/denkhaus/agents/pkg/tools/project"
	"github.com/denkhaus/agents/pkg/tools/shell"
	"github.com/denkhaus/agents/pkg/tools/tavily"
	timetools "github.com/denkhaus/agents/pkg/tools/time"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/duckduckgo"
	"trpc.group/trpc-go/trpc-agent-go/tool/file"
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
func (f *cueToolFactoryImpl) CreateTools(toolsConfig ToolsConfig) ([]tool.Tool, []tool.ToolSet, error) {
	var tools []tool.Tool
	var toolsets []tool.ToolSet

	// Create individual tools
	for toolName, toolConfig := range toolsConfig.Tools {
		if !toolConfig.Enabled {
			continue
		}

		// Create the tool using the registered factory
		tool, err := f.createTool(toolName, toolConfig.Config)
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
		toolset, err := f.createToolSet(toolsetName, toolsetConfig.Config)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create toolset %s: %w", toolsetName, err)
		}

		toolsets = append(toolsets, toolset)
	}

	return tools, toolsets, nil
}

// createTool creates a single tool using the tool provider's factories
func (f *cueToolFactoryImpl) createTool(toolName string, config map[string]interface{}) (tool.Tool, error) {
	// This is a simplified implementation - in a real scenario, you'd access
	// the tool provider's internal factories or expose them through the interface
	switch toolName {
	case calculator.ToolName:
		return f.createCalculatorTool(config)
	case fetch.ToolName:
		return f.createFetchTool(config)
	case timetools.ToolName:
		return f.createTimeTool(config)
	case "duckduckgo":
		return f.createDuckDuckGoTool(config)
	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
}

// createToolSet creates a single toolset using the tool provider's factories
func (f *cueToolFactoryImpl) createToolSet(toolsetName string, config map[string]interface{}) (tool.ToolSet, error) {
	// This is a simplified implementation - in a real scenario, you'd access
	// the tool provider's internal factories or expose them through the interface
	switch toolsetName {
	case tavily.ToolSetName:
		return f.createTavilyToolSet(config)
	case project.ToolSetName:
		return f.createProjectToolSet(config)
	case shell.ToolSetName:
		return f.createShellToolSet(config)
	case "file":
		return f.createFileToolSet(config)
	default:
		return nil, fmt.Errorf("unknown toolset: %s", toolsetName)
	}
}

func (f *cueToolFactoryImpl) createCalculatorTool(config map[string]interface{}) (tool.Tool, error) {
	return calculator.NewTool()
}

func (f *cueToolFactoryImpl) createFetchTool(config map[string]interface{}) (tool.Tool, error) {
	return fetch.NewTool()
}

func (f *cueToolFactoryImpl) createTimeTool(config map[string]interface{}) (tool.Tool, error) {
	return timetools.NewTool()
}

func (f *cueToolFactoryImpl) createDuckDuckGoTool(config map[string]interface{}) (tool.Tool, error) {
	return duckduckgo.NewTool(), nil
}

// Toolset creation methods
func (f *cueToolFactoryImpl) createTavilyToolSet(config map[string]interface{}) (tool.ToolSet, error) {
	var options []tavily.Option

	if apiKey, ok := config["api_key"].(string); ok && apiKey != "" {
		options = append(options, tavily.WithAPIKey(apiKey))
	}

	return tavily.NewToolSet(options...)
}

func (f *cueToolFactoryImpl) createProjectToolSet(config map[string]interface{}) (tool.ToolSet, error) {
	var options []project.Option
	return project.NewToolSet(options...)
}

func (f *cueToolFactoryImpl) createShellToolSet(config map[string]interface{}) (tool.ToolSet, error) {
	var options []shell.Option

	if baseDir, ok := config["base_dir"].(string); ok && baseDir != "" {
		options = append(options, shell.WithBaseDir(baseDir))
	}

	if timeout, ok := config["timeout"].(int); ok && timeout != 0 {
		options = append(options, shell.WithTimeout(time.Duration(timeout)*time.Second))
	}

	if allowedCommands, ok := config["allowed_commands"].([]interface{}); ok && allowedCommands != nil {
		var commands []string
		for _, cmd := range allowedCommands {
			if cmdStr, ok := cmd.(string); ok {
				commands = append(commands, cmdStr)
			}
		}
		if len(commands) > 0 {
			options = append(options, shell.WithAllowedCommands(commands))
		}
	}

	if executeEnabled, ok := config["execute_enabled"].(bool); ok {
		options = append(options, shell.WithExecuteCommandEnabled(executeEnabled))
	}

	return shell.NewToolSet(options...)
}

func (f *cueToolFactoryImpl) createFileToolSet(config map[string]interface{}) (tool.ToolSet, error) {
	var options []file.Option

	if baseDir, ok := config["base_dir"].(string); ok && baseDir != "" {
		options = append(options, file.WithBaseDir(baseDir))
	}

	if readOnly, ok := config["read_only"].(bool); ok && readOnly {
		options = append(options,
			file.WithListFileEnabled(true),
			file.WithReadFileEnabled(true),
			file.WithSearchFileEnabled(true),
			file.WithSearchContentEnabled(true),
			file.WithReplaceContentEnabled(false),
			file.WithSaveFileEnabled(false),
		)
	} else {
		options = append(options,
			file.WithListFileEnabled(true),
			file.WithReadFileEnabled(true),
			file.WithSearchFileEnabled(true),
			file.WithSearchContentEnabled(true),
			file.WithReplaceContentEnabled(true),
			file.WithSaveFileEnabled(true),
		)
	}

	if maxSize, ok := config["max_file_size"].(int); ok && maxSize != 0 {
		options = append(options, file.WithMaxFileSize(int64(maxSize)))
	}

	return file.NewToolSet(options...)
}
