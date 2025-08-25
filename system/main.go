package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/denkhaus/agents/di"
	"github.com/denkhaus/agents/provider/chat"
	"github.com/denkhaus/agents/shared"
	"github.com/samber/do"
)

func enterChat(ctx context.Context, chat chat.Chat) error {

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

func startup(ctx context.Context) error {

	injector := di.NewContainer()
	chatProvider := do.MustInvoke[chat.Provider](injector)

	chat, err := chatProvider.GetChat(ctx, shared.AgentIDCoder)
	if err != nil {
		return fmt.Errorf("failed to getchat for agent %s", shared.AgentIDCoder)
	}

	return enterChat(ctx, chat)
}

func main() {
	if err := startup(context.Background()); err != nil {
		log.Fatal(err)
	}
}
