package duck

import (
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/duckduckgo"
)

const (
	ToolName = "duckduckgo_search"
)

type FactoryFunc func(opts ...duckduckgo.Option) (tool.Tool, error)

func NewWithDI(injector *do.Injector) (FactoryFunc, error) {
	return func(opts ...duckduckgo.Option) (tool.Tool, error) {
		return duckduckgo.NewTool(opts...), nil
	}, nil
}
