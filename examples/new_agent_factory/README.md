# New Agent Factory Example

This example demonstrates how to use the new agent factory system implemented in the `provider/config` package.

## Overview

The new agent factory system uses CUE configuration files to create agents with different types and configurations. It supports:

- Default (LLM) agents
- Chain agents (sequential execution)
- Cycle agents (cyclic execution with limits)
- Parallel agents (simultaneous execution)

## Running the Example

To run this example:

```bash
cd /path/to/agents/repository
go run examples/new_agent_factory/main.go
```

## Key Features Demonstrated

1. **Factory Creation**: Creating a new agent factory with CUE configuration path
2. **Configuration Validation**: Validating CUE configurations
3. **Agent Creation**: Creating agents by name and by ID
4. **Configuration Access**: Accessing raw agent configurations

## Migration from Old System

To migrate from the old system:

1. Replace imports of `provider/agent` and `provider/settings` with `provider/config`
2. Update agent creation code to use the new factory methods
3. Convert JSON/YAML configurations to CUE format
4. Update agent configurations to include type specifications

## Configuration Structure

The example assumes CUE configurations are in the `./config` directory with the following structure:

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