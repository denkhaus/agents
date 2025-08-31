package file

import (
	"github.com/denkhaus/agents/pkg/tools"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/file"
)

const (
	ToolSetName = "file_toolset"
)

type FileToolSet struct {
	WorkspacePath string
	ReadOnly      bool
}

func NewWithDI(injector *do.Injector) (tools.ToolSetFactoryFunc, error) {
	return func(config tools.ConfigPayload) (tool.ToolSet, error) {
		// Extract configuration and convert to options
		var settings FileToolSet
		if err := config.Bind(&settings); err != nil {
			return nil, err
		}

		// Create options from settings
		var opts []Option
		if settings.WorkspacePath != "" {
			opts = append(opts, WithWorkspacePath(settings.WorkspacePath))
		}
		opts = append(opts, WithReadOnly(settings.ReadOnly))

		return New(opts...)
	}, nil
}

func New(opts ...Option) (tool.ToolSet, error) {
	wrapper := FileToolSet{}

	for _, opt := range opts {
		opt(&wrapper)
	}

	return wrapper.create()
}

func (p *FileToolSet) create() (toolset tool.ToolSet, err error) {
	options := []file.Option{
		file.WithBaseDir(p.WorkspacePath),
	}

	if p.ReadOnly {
		// Create readonly file operation tools.
		options = append(options,
			file.WithListFileEnabled(true),
			file.WithReadFileEnabled(true),
			file.WithReplaceContentEnabled(false),
			file.WithSaveFileEnabled(false),
			file.WithSearchFileEnabled(true),
			file.WithSearchContentEnabled(true),
		)
	}

	return file.NewToolSet(options...)
}
