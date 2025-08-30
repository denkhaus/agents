package di

import (
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/provider/agent"

	"github.com/denkhaus/agents/provider/prompt"
	"github.com/denkhaus/agents/provider/settings"
	"github.com/denkhaus/agents/provider/workspace"

	"github.com/samber/do"
)

func NewContainer() *do.Injector {
	injector := do.New()

	do.Provide(injector, workspace.New)
	do.Provide(injector, prompt.New)
	do.Provide(injector, agent.New)
	do.Provide(injector, settings.New)
	do.Provide(injector, logger.New)

	return injector
}
