# New Agent Factory System

This document describes the new agent factory system implemented in the `provider/config` package, which uses CUE configuration files to create agents with different types and configurations.

## Overview

The new agent factory system replaces the previous implementation in `provider/agent` and `provider/settings` with a more flexible and declarative approach using CUE configuration files. It supports all the functionality of the previous system while providing enhanced configurability.

## Features

- **Multiple Agent Types**: Supports default (LLM), chain, cycle, and parallel agent types
- **CUE Configuration**: Uses CUE for type-safe and declarative configuration
- **Environment-based Configurations**: Supports different configurations for different environments (development, production, etc.)
- **Tool Integration**: Creates tools and toolsets from configuration
- **Dependency Injection**: Integrates with the existing dependency injection system

## Agent Types

The system supports four types of agents:

1. **Default (LLM) Agents**: Standard language model agents with prompts and tools
2. **Chain Agents**: Execute sub-agents in sequential order
3. **Cycle Agents**: Execute sub-agents in a cyclic pattern with a maximum iteration limit
4. **Parallel Agents**: Execute sub-agents in parallel

## Configuration Structure

The configuration is organized in the following structure:

```
config/
├── compositions/
│   ├── environments/
│   │   ├── development.cue
│   │   └── production.cue
│   └── stable/
├── prompts/
├── settings/
└── tools/
```

### Agent Configuration

Each agent is configured with the following properties:

- `agent_id`: Unique identifier for the agent
- `name`: Human-readable name
- `description`: Description of the agent's purpose
- `type`: Agent type (default, chain, cycle, parallel)
- `prompt`: Prompt configuration
- `settings`: Runtime settings
- `tools`: Tool configurations

## Usage

### Creating an Agent Factory

```go
factory := config.NewUnifiedAgentFactory("./config")
```

### Creating an Agent

```go
// Create an agent by name
agent, err := factory.CreateAgent(ctx, "production", "coder")

// Create an agent by ID
agent, err := factory.CreateAgentByID(ctx, shared.AgentIDCoder)
```

### Validating Configuration

```go
if err := factory.ValidateConfiguration(); err != nil {
    log.Fatalf("Configuration validation failed: %v", err)
}
```

## Migration from Old System

To migrate from the old system:

1. Replace `provider/agent` and `provider/settings` usage with `provider/config`
2. Convert JSON/YAML configurations to CUE format
3. Update agent creation code to use the new factory methods

## Examples

See the `examples/new_agent_factory` directory for a complete example of how to use the new system.