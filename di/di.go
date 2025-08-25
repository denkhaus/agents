package di

import "github.com/samber/do"

func NewContainer() *do.Injector {
	injector := do.New()

	provideTools(injector)
	return injector
}
