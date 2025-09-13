//cue:generate cue get go github.com/denkhaus/agents/pkg/tools/agentinfo

package agentinfo

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/pkg/shared"
	"github.com/denkhaus/agents/pkg/tools"
	"github.com/google/uuid"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// -----------------------------------------------------------------------------
// Constants & helpers ----------------------------------------------------------
// -----------------------------------------------------------------------------
const (
	ToolName = "get_agent_info"
)

// AgentInfoToolConfig holds configuration for the agent info tool
type AgentInfoToolConfig struct {
	AvailableAgents []*shared.AgentInfo `json:"available_agents" mapstructure:"available_agents"`
	CallingAgentID  uuid.UUID           `json:"calling_agent_id" mapstructure:"calling_agent_id"`
}

// getAgentInfoArgs holds the input for the agent info tool.
type getAgentInfoArgs struct {
	// No parameters needed - always returns all available agents
}

// AgentInfoResponse represents information about a single agent
type AgentInfoResponse struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Role        shared.AgentRole `json:"role"`
	IsSelf      bool             `json:"is_self"` // Indicates if this is the calling agent's own info
}

// getAgentInfoResult holds the output for the agent info tool.
type getAgentInfoResult struct {
	Agents []AgentInfoResponse `json:"agents"`
	Count  int                 `json:"count"`
}

// getAgentInfo performs the agent information retrieval operation.
// It returns information about all available agents, marking the calling agent's info as "self".
func getAgentInfo(availableAgents []*shared.AgentInfo, callingAgentID uuid.UUID) func(context.Context, getAgentInfoArgs) (getAgentInfoResult, error) {
	return func(ctx context.Context, args getAgentInfoArgs) (getAgentInfoResult, error) {
		// Always return all available agents
		agents := make([]AgentInfoResponse, len(availableAgents))
		for i, agent := range availableAgents {
			isSelf := agent.ID == callingAgentID
			agents[i] = AgentInfoResponse{
				ID:          agent.ID,
				Name:        agent.Name,
				Description: agent.Description,
				Role:        agent.Role,
				IsSelf:      isSelf,
			}
		}

		return getAgentInfoResult{
			Agents: agents,
			Count:  len(agents),
		}, nil
	}
}

func NewWithDI(injector *do.Injector) (tools.ToolFactoryFunc, error) {
	return func(config tools.ConfigPayload, availableAgents []*shared.AgentInfo) (tool.Tool, error) {
		var toolConfig AgentInfoToolConfig
		if err := config.Bind(&toolConfig); err != nil {
			return nil, fmt.Errorf("failed to bind agent info tool config: %w", err)
		}

		// Use available agents from the factory function if not provided in config
		agents := toolConfig.AvailableAgents
		if len(agents) == 0 {
			agents = availableAgents
		}

		return New(agents, toolConfig.CallingAgentID)
	}, nil
}

func New(availableAgents []*shared.AgentInfo, callingAgentID uuid.UUID) (tool.Tool, error) {
	// Create agent info tool for retrieving information about available agents.
	agentInfoTool := function.NewFunctionTool(
		getAgentInfo(availableAgents, callingAgentID),
		function.WithName(ToolName),
		function.WithDescription(
			"Get information about all available agents in the system. Returns details like name, description, and role for each agent. Your own information will be marked with 'is_self: true'.",
		),
	)

	return agentInfoTool, nil
}
