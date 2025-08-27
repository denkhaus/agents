package plugins

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/denkhaus/agents/multi"
	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
)

type ChatPlugin interface {
	Start(ctx context.Context)
}

// ChatSystem manages the multi-agent chat
type cliMultiAgentChatImpl struct {
	Options
}

func NewCLIMultiAgentChat(opts ...MultiAgentChatOption) ChatPlugin {
	chat := &cliMultiAgentChatImpl{}

	for _, opt := range opts {
		opt(&chat.Options)
	}

	chat.processor = multi.NewChatProcessor(
		multi.WithHumanAgent(shared.NewHumanAgent(shared.AgentInfoHuman)),
		multi.WithOnProgress(chat.handleOnProgress),
		multi.WithOnMessage(chat.handleOnMessage),
		multi.WithOnError(chat.handleOnError),
	)

	chat.setupMessageListener()
	return chat
}

func (p *cliMultiAgentChatImpl) setupMessageListener() {
	// Add a message interceptor to the broker
	p.processor.SetMessageInterceptor(func(fromID, toID uuid.UUID, content string) {
		fromName := p.processor.GetAgentNameByID(fromID)
		toName := p.processor.GetAgentNameByID(toID)

		if fromName != "" && toName != "" {
			// Format: "FromName (FromID) -> ToName (ToID)"
			header := fmt.Sprintf("%s (%s) -> %s (%s)",
				fromName, shortenID(fromID.String()),
				toName, shortenID(toID.String()))
			p.printWithBorder(header, content)
		}
	})
}

func (p *cliMultiAgentChatImpl) handleOnProgress(format string, a ...any) {
	fmt.Printf(format, a...)
}

func (p *cliMultiAgentChatImpl) handleOnMessage(info *shared.AgentInfo, content string) {
	p.printWithBorder(info.String(), content)
}

func (p *cliMultiAgentChatImpl) handleOnError(info *shared.AgentInfo, err error) {
	p.printWithBorder(info.String(), err.Error())
}

// startChat runs the interactive chat loop
func (p *cliMultiAgentChatImpl) Start(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("you >> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Handle commands
		if strings.HasPrefix(input, "/") {
			command := strings.TrimPrefix(input, "/")
			switch command {
			case "exit":
				fmt.Println("Goodbye!")
				return
			case "list":
				fmt.Println("\n=== Available Agents ===")
				for _, info := range p.processor.GetAllAgentInfos() {
					fmt.Printf("- %s (ID: %s)\n", info.Name, shortenID(info.ID.String()))
				}
				fmt.Println("========================")
			default:
				fmt.Println("Unknown command. Use /list or /exit")
			}
			continue
		}

		// Handle agent messages (@agent message)
		if strings.HasPrefix(input, "@") {
			parts := strings.SplitN(input[1:], " ", 2)
			if len(parts) < 2 {
				fmt.Println("Usage: @<agent> <message>")
				continue
			}

			toAgentInfo := p.processor.GetAgentInfoByAuthor(parts[0])
			if toAgentInfo == nil {
				fmt.Printf("ERROR: could not find agent info for %s\n", parts[0])
				continue
			}

			err := p.processor.SendMessageWithProcessing(ctx, shared.AgentIDHuman, toAgentInfo.ID, parts[1])
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				continue
			}

			continue
		}

		fmt.Println("Use @<agent> <message> to send a message, or /help for commands")
	}
}

// printWithBorder prints a message with a simple, clean border
func (p *cliMultiAgentChatImpl) printWithBorder(sender, message string) {
	// Fixed width for consistency
	const width = 120

	// Print top border
	fmt.Println(strings.Repeat("=", width))

	// Print sender with padding
	senderLine := fmt.Sprintf("[ %s ]", sender)
	if len(senderLine) > width-2 {
		senderLine = senderLine[:width-5] + "..."
	}
	fmt.Printf("%s%s\n", senderLine, strings.Repeat(" ", width-len(senderLine)))

	// Print separator
	fmt.Println(strings.Repeat("-", width))

	// Print message lines
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		if len(line) > width-2 {
			// Split long lines
			for len(line) > width-2 {
				fmt.Printf("%s\n", line[:width-2])
				line = line[width-2:]
			}
			if len(line) > 0 {
				fmt.Printf("%s\n", line)
			}
		} else {
			fmt.Printf("%s\n", line)
		}
	}

	// Print bottom border
	fmt.Println(strings.Repeat("=", width))
	fmt.Println() // Extra line for spacing
}

// shortenID safely shortens a UUID string to first 8 characters
func shortenID(id string) string {
	if len(id) >= 8 {
		return id[:8] + "..."
	}
	return id
}
