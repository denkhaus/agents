package di

import (
	"github.com/denkhaus/agents/provider/workspace"
	"github.com/samber/do"
)

func NewContainer() *do.Injector {
	injector := do.New()

	do.Provide(injector, workspace.New)

	provideTools(injector)
	return injector
}
