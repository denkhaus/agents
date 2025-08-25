package utils

import (
	"context"

	"github.com/denkhaus/agents/shared"
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

func GetToolInfoFromSets(ctx context.Context, sets []tool.ToolSet) []shared.ToolInfo {
	var toolInfos []shared.ToolInfo

	for _, toolSet := range sets {
		tools := toolSet.Tools(ctx)

		for _, tool := range tools {
			decl := tool.Declaration()
			toolInfo := shared.ToolInfo{
				Name:        decl.Name,
				Description: decl.Description,
			}
			toolInfos = append(toolInfos, toolInfo)
		}
	}

	return toolInfos
}
