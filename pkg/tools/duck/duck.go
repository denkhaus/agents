package duck

import (
	"github.com/denkhaus/agents/pkg/tools"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/duckduckgo"
)

const (
	ToolName = "duckduckgo_search"
)

func NewWithDI(injector *do.Injector) (tools.ToolFactoryFunc, error) {
	return func(config tools.ConfigPayload) (tool.Tool, error) {
		// Extract options from config if needed
		// For now, use default options
		return duckduckgo.NewTool(), nil
	}, nil
}
