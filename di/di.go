package di

import (
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/config"
	"github.com/denkhaus/agents/pkg/provider/agent"
	agentcconfig "github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/pkg/tools/calculator"
	"github.com/denkhaus/agents/pkg/tools/duck"
	"github.com/denkhaus/agents/pkg/tools/fetch"
	"github.com/denkhaus/agents/pkg/tools/project"
	"github.com/denkhaus/agents/pkg/tools/shell"
	"github.com/denkhaus/agents/pkg/tools/state"
	"github.com/denkhaus/agents/pkg/tools/tavily"
	"github.com/denkhaus/agents/pkg/tools/time"

	"github.com/denkhaus/agents/pkg/provider/prompt"
	"github.com/denkhaus/agents/pkg/provider/settings"

	"github.com/samber/do"
)

func NewContainer() *do.Injector {
	injector := do.New()

	do.Provide(injector, config.NewWithDI)
	do.Provide(injector, agentcconfig.NewCUEToolFactory)
	do.Provide(injector, agentcconfig.NewCUEConfigProvider)

	do.Provide(injector, prompt.NewWithDI)
	do.Provide(injector, agent.NewWithDI)
	do.Provide(injector, settings.NewWithDI)
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

	return injector
}
