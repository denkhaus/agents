package di

import (
	"github.com/denkhaus/agents/tools/calculator"
	"github.com/denkhaus/agents/tools/time"
	"github.com/samber/do"
)

func provideTools(i *do.Injector) {
	do.ProvideNamed(i, calculator.ToolName, calculator.NewTool)
	do.ProvideNamed(i, time.ToolName, time.NewTool)
}
