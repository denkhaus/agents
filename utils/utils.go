package utils

import (
	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
}

func FloatPtr(f float64) *float64 {
	return &f
}

type ToolInfo struct {
	Name        string
	Description string
}

func GetToolInfo(tools ...tool.Tool) []shared.ToolInfo {
	var toolInfos []shared.ToolInfo

	for _, tool := range tools {
		decl := tool.Declaration()
		toolInfo := shared.ToolInfo{
			Name:        decl.Name,
			Description: decl.Description,
		}
		toolInfos = append(toolInfos, toolInfo)
	}

	return toolInfos
}

func GetAgentInfoForAgent(agentID uuid.UUID, availableAgents ...shared.AgentInfo) []shared.AgentInfo {
	var info []shared.AgentInfo = []shared.AgentInfo{
		shared.AgentInfoHuman,
	}

	for _, agent := range availableAgents {
		// Add other AI agents (excluding self)
		if agent.ID() != agentID {
			info = append(info, agent)
		}
	}

	return info
}
