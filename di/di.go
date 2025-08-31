package di

import (
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/config"
	agentcconfig "github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/pkg/tools/calculator"
	"github.com/denkhaus/agents/pkg/tools/duck"
	"github.com/denkhaus/agents/pkg/tools/fetch"
	"github.com/denkhaus/agents/pkg/tools/file"
	"github.com/denkhaus/agents/pkg/tools/project"
	"github.com/denkhaus/agents/pkg/tools/shell"
	"github.com/denkhaus/agents/pkg/tools/state"
	"github.com/denkhaus/agents/pkg/tools/tavily"
	"github.com/denkhaus/agents/pkg/tools/time"

	"github.com/samber/do"
)

func NewContainer() *do.Injector {
	injector := do.New()

	do.Provide(injector, config.NewWithDI)
	do.Provide(injector, agentcconfig.NewCUEToolFactory)
	do.Provide(injector, agentcconfig.NewCUEConfigProvider)
	do.Provide(injector, agentcconfig.NewUnifiedAgentFactory)

	do.Provide(injector, logger.NewWithDI)

	// provide tools
	do.ProvideNamed(injector, fetch.ToolName, fetch.NewWithDI)
	do.ProvideNamed(injector, calculator.ToolName, calculator.NewWithDI)
	do.ProvideNamed(injector, time.ToolName, time.NewWithDI)
	do.ProvideNamed(injector, tavily.ToolSetName, tavily.NewWithDI)
	do.ProvideNamed(injector, shell.ToolSetName, shell.NewWithDI)
	do.ProvideNamed(injector, project.ToolSetName, project.NewWithDI)
	do.ProvideNamed(injector, state.ToolName, state.NewWithDI)
	do.ProvideNamed(injector, duck.ToolName, duck.NewWithDI)
	do.ProvideNamed(injector, file.ToolSetName, file.NewWithDI)

	return injector
}
