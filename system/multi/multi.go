package multi

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/provider/agent"
	"github.com/denkhaus/agents/provider/chat"
	"github.com/denkhaus/agents/shared"
	"github.com/denkhaus/agents/tools/calculator"
	"github.com/denkhaus/agents/tools/project"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

type MultiAgentSystem interface {
	EnterChat(ctx context.Context, chat provider.Chat) error
	CreateCoderAgentChat(ctx context.Context, injector *do.Injector, opts ...provider.ChatProviderOption) (provider.Chat, error)
	CreateProjectManagerChat(ctx context.Context, injector *do.Injector, opts ...provider.ChatProviderOption) (provider.Chat, error)
}

type multiAgentSystemImpl struct {
	chatProvider provider.ChatProvider
}

func New(injector *do.Injector) (MultiAgentSystem, error) {
	chatProvider := do.MustInvoke[provider.ChatProvider](injector)
	sys := &multiAgentSystemImpl{
		chatProvider: chatProvider,
	}

	return sys, nil
}

func (p *multiAgentSystemImpl) EnterChat(ctx context.Context, chat provider.Chat) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("👤 You: ")
		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "" {
			continue
		}

		// Handle exit command.
		if strings.ToLower(userInput) == "exit" {
			fmt.Println("👋 Goodbye!")
			return nil
		}

		// Process the user message.
		if err := chat.ProcessMessage(ctx, userInput); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}

		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("input scanner error: %w", err)
	}

	return nil
}

func (p *multiAgentSystemImpl) CreateCoderAgentChat(
	ctx context.Context,
	injector *do.Injector,
	opts ...provider.ChatProviderOption,
) (provider.Chat, error) {
	calculatorTool := do.MustInvokeNamed[tool.CallableTool](injector, calculator.ToolName)

	opts = append(opts, chat.WithAgentProviderOptions(
		agent.WithLLMAgentOptions(
			llmagent.WithTools([]tool.Tool{calculatorTool}),
		),
	))

	return p.chatProvider.GetChat(ctx, shared.AgentIDCoder, opts...)
}

func (p *multiAgentSystemImpl) CreateProjectManagerChat(
	ctx context.Context,
	injector *do.Injector,
	opts ...provider.ChatProviderOption,
) (provider.Chat, error) {

	projectManagerToolSet, err := project.NewToolSet()
	if err != nil {
		return nil, fmt.Errorf("failed to create project manager toolset: %w", err)
	}

	opts = append(opts, chat.WithAgentProviderOptions(
		agent.WithLLMAgentOptions(
			llmagent.WithToolSets([]tool.ToolSet{projectManagerToolSet}),
		),
	))

	return p.chatProvider.GetChat(ctx, shared.AgentIDProjectManager, opts...)
}

// func (p *agentSettingsImpl) getToolSets() ([]tool.ToolSet, error) {
// 	workspacePath, err := p.workspace.GetWorkspacePath()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get workspacePath for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
// 	}
// 	// Create file operation tools.
// 	fileToolSet, err := file.NewToolSet(
// 		file.WithBaseDir(workspacePath),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("create file tool set: %w", err)
// 	}

// 	shellToolSet, err := shelltoolset.NewToolSet(
// 		shelltoolset.WithBaseDir(workspacePath),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("create shell tool set: %w", err)
// 	}

// 	return []tool.ToolSet{fileToolSet, shellToolSet}, nil
// }
