package file

import (
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/file"
)

const (
	ToolSetName = "file_toolset"
)

type fileToolSetWrapperImpl struct {
	workspacePath string
	readOnly      bool
}

// Option is a configuration function for the file ToolSet.
type Option func(*fileToolSetWrapperImpl)

func WithReadOnly(readOnly bool) Option {
	return func(t *fileToolSetWrapperImpl) {
		t.readOnly = readOnly
	}
}

func WithWorkspacePath(workspacePath string) Option {
	return func(t *fileToolSetWrapperImpl) {
		t.workspacePath = workspacePath
	}
}

type FactoryFunc func(opts ...Option) (tool.ToolSet, error)

func NewWithDI(injector *do.Injector) (FactoryFunc, error) {
	return func(opts ...Option) (tool.ToolSet, error) {
		return New(opts...)
	}, nil
}

func New(opts ...Option) (tool.ToolSet, error) {
	wrapper := fileToolSetWrapperImpl{}

	for _, opt := range opts {
		opt(&wrapper)
	}

	return wrapper.create()
}

func (p *fileToolSetWrapperImpl) create() (toolset tool.ToolSet, err error) {
	options := []file.Option{
		file.WithBaseDir(p.workspacePath),
	}

	if p.readOnly {
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
