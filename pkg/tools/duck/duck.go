//cue:generate cue get go github.com/denkhaus/agents/pkg/tools/duck

package duck

import (
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/denkhaus/agents/pkg/tools"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/duckduckgo"
)

const (
	ToolName = "duckduckgo_search"
)

// ToolConfig holds configuration for the DuckDuckGo search tool
type ToolConfig struct {
	// Currently no configurable options, but struct exists for consistency
}

func NewWithDI(injector *do.Injector) (tools.ToolFactoryFunc, error) {
	return func(config tools.ConfigPayload, _ []*shared.AgentInfo) (tool.Tool, error) {
		// Extract options from config if needed
		// For now, use default options
		return duckduckgo.NewTool(), nil
	}, nil
}
