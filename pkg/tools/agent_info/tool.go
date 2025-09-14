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
	ToolName = "list_available_agents"
)

// AgentInfoToolConfig holds configuration for the agent info tool
type AgentInfoToolConfig struct {
	AvailableAgents          []*shared.AgentInfo `json:"available_agents" mapstructure:"available_agents"`
	AllowedToCommunicateWith []uuid.UUID         `json:"allowed_to_communicate_with" mapstructure:"allowed_to_communicate_with"`
	CallingAgentID           uuid.UUID           `json:"calling_agent_id" mapstructure:"calling_agent_id"`
}

// listAvailableAgentsArgs defines the arguments for listing available agents (empty struct)
type listAvailableAgentsArgs struct{}

// listAvailableAgentsResult defines the result of listing available agents
type listAvailableAgentsResult struct {
	Agents  []AgentInfo `json:"agents" description:"List of available agents in the system"`
	Count   int         `json:"count" description:"Number of available agents"`
	Message string      `json:"message" description:"A message describing the result"`
}

// AgentInfo represents information about an available agent
type AgentInfo struct {
	ID                   uuid.UUID        `json:"id" description:"Unique identifier of the agent. Use this ID to communicate with the agent using the 'send_message' tool"`
	Name                 string           `json:"name" description:"Display name of the agent"`
	Role                 shared.AgentRole `json:"role" description:"Role/type of the agent (e.g., researcher, coder, project-manager)"`
	Description          string           `json:"description" description:"Detailed description of the agent's capabilities and purpose"`
	IsSelf               bool             `json:"is_self" description:"Your own information will be marked with 'true'."`
	AllowedToCommunicate bool             `json:"allowed_to_communicate" description:"This defines if you are allowed to communicate with the agent."`
}

type ListAvailableAgentsFunc func(context.Context, listAvailableAgentsArgs) (listAvailableAgentsResult, error)

// listAvailableAgents performs the agent information retrieval operation.
// It returns information about all available agents, marking the calling agent's info as "self".
func listAvailableAgents(
	availableAgents []*shared.AgentInfo,
	callingAgentID uuid.UUID,
	allowedToCommunicateWith []uuid.UUID,
) ListAvailableAgentsFunc {

	isAllowed := func(agentID uuid.UUID) bool {
		for _, id := range allowedToCommunicateWith {
			if agentID == id {
				return true
			}
		}

		return false
	}

	return func(ctx context.Context, args listAvailableAgentsArgs) (listAvailableAgentsResult, error) {

		// Always return all available agents
		agents := make([]AgentInfo, len(availableAgents))
		for i, agent := range availableAgents {
			isSelf := agent.ID == callingAgentID
			allowed := isAllowed(agent.ID)
			if isSelf {
				allowed = false
			}

			agents[i] = AgentInfo{
				ID:                   agent.ID,
				Name:                 agent.Name,
				Description:          agent.Description,
				Role:                 agent.Role,
				IsSelf:               isSelf,
				AllowedToCommunicate: allowed,
			}
		}

		return listAvailableAgentsResult{
			Agents:  agents,
			Count:   len(agents),
			Message: fmt.Sprintf("Found %d available agents", len(agents)),
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

		return New(
			agents,
			toolConfig.CallingAgentID,
			toolConfig.AllowedToCommunicateWith,
		)

	}, nil
}

func New(availableAgents []*shared.AgentInfo, callingAgentID uuid.UUID, allowedToCommunicateWith []uuid.UUID) (tool.Tool, error) {
	// Create agent info tool for retrieving information about available agents.
	agentInfoTool := function.NewFunctionTool(
		listAvailableAgents(availableAgents, callingAgentID, allowedToCommunicateWith),
		function.WithName(ToolName),
		function.WithDescription(`List all available agents in the system with their capabilities and descriptions.
		Use this to discover which agents you can communicate with by the 'send_message' tool.
		Each agent has a unique ID, name, role, and detailed description of their capabilities.
		`))

	return agentInfoTool, nil
}
