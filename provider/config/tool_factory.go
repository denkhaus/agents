package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/denkhaus/agents/tools/calculator"
	"github.com/denkhaus/agents/tools/fetch"
	"github.com/denkhaus/agents/tools/project"
	"github.com/denkhaus/agents/tools/shell"
	"github.com/denkhaus/agents/tools/tavily"
	timetools "github.com/denkhaus/agents/tools/time"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/duckduckgo"
	"trpc.group/trpc-go/trpc-agent-go/tool/file"
)

// cueToolFactoryImpl creates tools from CUE-based configuration
type cueToolFactoryImpl struct{}

// NewCUEToolFactory creates a new CUE-based tool factory
func NewCUEToolFactory() ToolFactory {
	return &cueToolFactoryImpl{}
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

		// Resolve environment variables in config
		resolvedConfig, err := f.resolveConfig(toolConfig.Config)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to resolve config for tool %s: %w", toolName, err)
		}

		// Create the tool using the registered factory
		tool, err := f.createTool(toolName, resolvedConfig)
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

		// Resolve environment variables in config
		resolvedConfig, err := f.resolveConfig(toolsetConfig.Config)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to resolve config for toolset %s: %w", toolsetName, err)
		}

		// Create the toolset using the registered factory
		toolset, err := f.createToolSet(toolsetName, resolvedConfig)
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

// resolveConfig resolves environment variables in configuration
func (f *cueToolFactoryImpl) resolveConfig(config map[string]interface{}) (map[string]interface{}, error) {
	if config == nil {
		return make(map[string]interface{}), nil
	}

	resolved := make(map[string]interface{})
	for key, value := range config {
		resolvedValue, err := f.resolveValue(value)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve config key %s: %w", key, err)
		}
		resolved[key] = resolvedValue
	}
	return resolved, nil
}

// resolveValue resolves environment variables in a single value
func (f *cueToolFactoryImpl) resolveValue(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return f.resolveStringValue(v)
	case map[string]interface{}:
		return f.resolveConfig(v)
	case []interface{}:
		resolved := make([]interface{}, len(v))
		for i, item := range v {
			resolvedItem, err := f.resolveValue(item)
			if err != nil {
				return nil, err
			}
			resolved[i] = resolvedItem
		}
		return resolved, nil
	default:
		return value, nil
	}
}

// resolveStringValue resolves environment variables in string values
func (f *cueToolFactoryImpl) resolveStringValue(value string) (interface{}, error) {
	if !strings.HasPrefix(value, "env:") {
		return value, nil
	}

	// Parse env:VAR_NAME or env:VAR_NAME:default
	parts := strings.SplitN(value, ":", 3)
	if len(parts) < 2 {
		return value, nil // Return original if malformed
	}

	envVarName := parts[1]
	if envVarName == "" {
		return value, nil // Return original if empty var name
	}

	envValue := os.Getenv(envVarName)

	// If no value and we have a default
	if envValue == "" && len(parts) == 3 {
		return parts[2], nil // Return default value
	}

	// If no value and no default, return empty string (don't error)
	if envValue == "" {
		return "", nil
	}

	return envValue, nil
}

// getConfigValue safely extracts a typed value from a config map
func (f *cueToolFactoryImpl) getConfigValue(config map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if config == nil {
		return defaultValue
	}
	if value, exists := config[key]; exists {
		return value
	}
	return defaultValue
}

// Individual tool creation methods
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

	if apiKey := f.getConfigValue(config, "api_key", ""); apiKey != "" {
		if keyStr, ok := apiKey.(string); ok {
			options = append(options, tavily.WithAPIKey(keyStr))
		}
	}

	return tavily.NewToolSet(options...)
}

func (f *cueToolFactoryImpl) createProjectToolSet(config map[string]interface{}) (tool.ToolSet, error) {
	var options []project.Option
	return project.NewToolSet(options...)
}

func (f *cueToolFactoryImpl) createShellToolSet(config map[string]interface{}) (tool.ToolSet, error) {
	var options []shell.Option

	if baseDir := f.getConfigValue(config, "base_dir", ""); baseDir != "" {
		if dirStr, ok := baseDir.(string); ok {
			options = append(options, shell.WithBaseDir(dirStr))
		}
	}

	if timeout := f.getConfigValue(config, "timeout", 0); timeout != 0 {
		if timeoutInt, ok := timeout.(int); ok {
			options = append(options, shell.WithTimeout(time.Duration(timeoutInt)*time.Second))
		}
	}

	if allowedCommands := f.getConfigValue(config, "allowed_commands", []string{}); allowedCommands != nil {
		if cmdList, ok := allowedCommands.([]interface{}); ok {
			var commands []string
			for _, cmd := range cmdList {
				if cmdStr, ok := cmd.(string); ok {
					commands = append(commands, cmdStr)
				}
			}
			if len(commands) > 0 {
				options = append(options, shell.WithAllowedCommands(commands))
			}
		}
	}

	executeEnabled := f.getConfigValue(config, "execute_enabled", true)
	if enabled, ok := executeEnabled.(bool); ok {
		options = append(options, shell.WithExecuteCommandEnabled(enabled))
	}

	return shell.NewToolSet(options...)
}

func (f *cueToolFactoryImpl) createFileToolSet(config map[string]interface{}) (tool.ToolSet, error) {
	var options []file.Option

	if baseDir := f.getConfigValue(config, "base_dir", ""); baseDir != "" {
		if dirStr, ok := baseDir.(string); ok {
			options = append(options, file.WithBaseDir(dirStr))
		}
	}

	readOnly := f.getConfigValue(config, "read_only", false)
	if ro, ok := readOnly.(bool); ok && ro {
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

	if maxSize := f.getConfigValue(config, "max_file_size", 0); maxSize != 0 {
		if sizeInt, ok := maxSize.(int); ok {
			options = append(options, file.WithMaxFileSize(int64(sizeInt)))
		}
	}

	return file.NewToolSet(options...)
}
