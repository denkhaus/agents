package agents

import (
	"fmt"

	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/file"
)

func CreateFileToolset(workspacePath string, readOnly bool) (toolset tool.ToolSet, err error) {
	options := []file.Option{
		file.WithBaseDir(workspacePath),
	}

	if readOnly {
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

	toolset, err = file.NewToolSet(options...)
	if err != nil {
		return nil, fmt.Errorf("create file tool set: %w", err)
	}

	return toolset, err
}
