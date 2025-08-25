package di

import (
	"github.com/denkhaus/agents/provider/agent"
	"github.com/denkhaus/agents/provider/chat"
	"github.com/denkhaus/agents/provider/prompt"
	"github.com/denkhaus/agents/provider/settings"
	"github.com/denkhaus/agents/provider/workspace"
	"github.com/samber/do"
)

func NewContainer() *do.Injector {
	injector := do.New()

	do.Provide(injector, workspace.New)
	do.Provide(injector, chat.New)
	do.Provide(injector, prompt.New)
	do.Provide(injector, agent.New)
	do.Provide(injector, settings.New)

	provideTools(injector)
	return injector
}
