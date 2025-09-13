package shared

import (
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

type ToolInfo struct {
	Name        string
	Description string
}

func GetToolInfo(tools ...tool.Tool) []ToolInfo {
	var toolInfos []ToolInfo

	for _, tool := range tools {
		decl := tool.Declaration()
		toolInfo := ToolInfo{
			Name:        decl.Name,
			Description: decl.Description,
		}
		toolInfos = append(toolInfos, toolInfo)
	}

	return toolInfos
}

func GetAgentInfoForAgent(agentID uuid.UUID, availableAgents ...*AgentInfo) []*AgentInfo {
	var info []*AgentInfo = []*AgentInfo{
		AgentInfoHuman,
	}

	for _, agent := range availableAgents {
		// Add other AI agents (excluding self)
		if agent.ID != agentID {
			info = append(info, agent)
		}
	}

	return info
}
