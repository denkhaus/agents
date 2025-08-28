package plugins

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/acarl005/stripansi"
	"github.com/denkhaus/agents/multi"
	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"github.com/mattn/go-runewidth"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// ANSI color codes for different message types
const (
	ColorReset = "\033[0m"
	ColorBold  = "\033[1m"

	// Message type colors
	ColorNormal    = "\033[37m" // White - normal messages
	ColorReasoning = "\033[33m" // Yellow - reasoning/planning messages
	ColorTool      = "\033[36m" // Cyan - tool call messages
	ColorIntercept = "\033[35m" // Magenta - intercepted messages
	ColorError     = "\033[31m" // Red - error messages
	ColorSystem    = "\033[32m" // Green - system messages

	// Border colors
	ColorBorderNormal    = "\033[90m" // Dark gray
	ColorBorderReasoning = "\033[93m" // Bright yellow
	ColorBorderTool      = "\033[96m" // Bright cyan
	ColorBorderIntercept = "\033[95m" // Bright magenta
)

// MessageType represents different types of messages for styling
type MessageType int

const (
	MessageTypeNormal MessageType = iota
	MessageTypeReasoningMessage
	MessageTypeToolCall
	MessageTypeIntercept
	MessageTypeError
	MessageTypeSystem
	MessageTypeAgentError
)

// ChatPlugin defines the interface for chat plugins that can be started.
type ChatPlugin interface {
	// Start begins the chat plugin operation with the given context.
	Start(ctx context.Context) error
}

// ChatSystem manages the multi-agent chat
type cliMultiAgentChatImpl struct {
	Options
	currentAgent *shared.AgentInfo // Track the currently selected agent
}

// NewCLIMultiAgentChat creates a new CLI-based multi-agent chat plugin.
// It sets up the chat processor with the provided options and configures message handling.
func NewCLIMultiAgentChat(opts ...MultiAgentChatOption) ChatPlugin {
	chat := &cliMultiAgentChatImpl{
		Options: Options{
			displayWidth: 120, // Default width
		},
	}

	for _, opt := range opts {
		opt(&chat.Options)
	}

	processorOptions := []multi.ChatProcessorOption{
		multi.WithOnProgress(chat.handleOnProgress),
		multi.WithOnMessage(chat.handleOnMessage),
		multi.WithOnReasoningMessage(chat.handleOnReasoningMessage),
		multi.WithOnError(chat.handleOnError),
		multi.WithOnToolCall(chat.handleOnToolCall),
	}

	processorOptions = append(processorOptions, chat.processorOptions...)
	chat.processor = multi.NewChatProcessor(processorOptions...)
	chat.setupMessageListener()

	return chat
}

// setupMessageListener configures the message interceptor to display inter-agent communication.
func (p *cliMultiAgentChatImpl) setupMessageListener() {
	// Add a message interceptor to the broker
	p.processor.SetMessageInterceptor(func(fromID, toID uuid.UUID, content string) {
		fromName := p.processor.GetAgentNameByID(fromID)
		toName := p.processor.GetAgentNameByID(toID)

		if fromName != "" && toName != "" {
			// Format: "FromName (FromID) -> ToName (ToID)"
			header := fmt.Sprintf("%s (%s) -> %s (%s)",
				fromName, shortenID(fromID.String()),
				toName, shortenID(toID.String()),
			)
			p.printWithBorderColored(header, content, MessageTypeIntercept)
		}
	})
}

// handleOnProgress handles progress updates by printing them to stdout.
func (p *cliMultiAgentChatImpl) handleOnProgress(messageType multi.SystemMessageType, format string, a ...any) {
	p.printSystemMessage(format, a...)
}

// handleOnMessage handles agent messages by displaying them with a formatted border.
func (p *cliMultiAgentChatImpl) handleOnMessage(info *shared.AgentInfo, content string) {
	p.printWithBorderColored(info.String(), content, MessageTypeNormal)
}

// handleOnError handles agent errors by displaying them with a formatted border.
func (p *cliMultiAgentChatImpl) handleOnError(info *shared.AgentInfo, err error) {
	p.printWithBorderColored(info.String(), fmt.Sprintf("%+v", err), MessageTypeAgentError)
}



// handleOnToolCall handles tool calls made by agents by displaying them with a formatted border.
func (p *cliMultiAgentChatImpl) handleOnToolCall(info *shared.AgentInfo, functionDef model.FunctionDefinitionParam) {
	toolCallInfo := fmt.Sprintf("Tool Call: %s", functionDef.Name)
	if len(functionDef.Arguments) > 0 {
		toolCallInfo += fmt.Sprintf("\nArguments: %s", string(functionDef.Arguments))
	}
	p.printWithBorderColored(info.String()+" [TOOL]", toolCallInfo, MessageTypeToolCall)
}

// handleOnReasoningMessage handles reasoning messages from agents.
func (p *cliMultiAgentChatImpl) handleOnReasoningMessage(info *shared.AgentInfo, reasoning string) {
	p.printWithBorderColored(info.String(), reasoning, MessageTypeReasoningMessage)
}

// Start runs the interactive chat loop, handling user input and agent communication.
// It supports commands like /exit, /list, /agent-name to select agents, and direct messaging.
func (p *cliMultiAgentChatImpl) Start(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Create prompt showing current agent
		prompt := "you"
		if p.currentAgent != nil {
			prompt = fmt.Sprintf("you [%s]", p.currentAgent.Name)
		}
		fmt.Printf("%s >> ", prompt)

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
				p.printSystemMessage("Goodbye!")
				return nil
			case "list":
				var builder strings.Builder
				builder.WriteString("\n=== Available Agents ===\n")
				for _, info := range p.processor.GetAllAgentInfos() {
					marker := ""
					if p.currentAgent != nil && info.ID == p.currentAgent.ID {
						marker = " (current)"
					}
					builder.WriteString(fmt.Sprintf("- %s (ID: %s)%s\n", info.Name, shortenID(info.ID.String()), marker))
				}
				builder.WriteString("========================")
				p.printSystemMessage(builder.String())
			case "clear":
				p.currentAgent = nil
				p.printSystemMessage("Current agent cleared. Use /<agent-name> to select an agent.")
			case "help":
				p.printSystemMessage(p.getHelpMessage())
			default:
				// Check if it's a width command
				if strings.HasPrefix(command, "width ") {
					widthStr := strings.TrimPrefix(command, "width ")
					if width, err := strconv.Atoi(widthStr); err == nil {
						if width < 40 {
							fmt.Println("Minimum width is 40 characters.")
							width = 40
						}
						p.displayWidth = width
						p.printSystemMessage("Display width set to %d characters.", width)
					} else {
						fmt.Println("Invalid width. Usage: /width <number>")
					}
					continue
				}
				// Try to find agent by name
				agentInfo := p.processor.GetAgentInfoByAuthor(command)
				if agentInfo != nil {
					p.currentAgent = agentInfo
					p.printSystemMessage("Selected agent: %s", agentInfo.Name)
				} else {
					p.printSystemMessage("Unknown command or agent: %s. Use /help for available commands.", command)
				}
			}
			continue
		}

		// Send message to current agent or show help
		if p.currentAgent != nil {
			err := p.processor.SendMessageWithProcessing(ctx, shared.AgentIDHuman, p.currentAgent.ID, input)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				continue
			}
		} else {
			fmt.Println("No agent selected. Use /<agent-name> to select an agent.")
		}
	}

	return nil
}

// printSystemMessage displays a system message with a standard border.
func (p *cliMultiAgentChatImpl) printSystemMessage(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	p.printWithBorderColored("SYSTEM", message, MessageTypeSystem)
}

// getHelpMessage returns the help message as a string.
func (p *cliMultiAgentChatImpl) getHelpMessage() string {
	var builder strings.Builder
	builder.WriteString("\n=== Available Commands ===\n")
	builder.WriteString("/help                 - Show this help message\n")
	builder.WriteString("/list                 - List all available agents\n")
	builder.WriteString("/clear                - Clear current agent selection\n")
	builder.WriteString(fmt.Sprintf("/width <number>       - Set display width (min: 40, current: %d)\n", p.displayWidth))
	builder.WriteString("/<agent-name>         - Select an agent to chat with\n")
	builder.WriteString("/exit                 - Exit the chat\n")
	builder.WriteString("\n")
	builder.WriteString("=== Usage ===\n")
	builder.WriteString("1. Select an agent: /project-manager\n")
	builder.WriteString("2. Chat directly: Hello, how can you help?\n")
	builder.WriteString("3. Switch agents: /another-agent\n")
	builder.WriteString("4. Adjust display: /width 80\n")
	builder.WriteString("==========================")
	return builder.String()
}

// printWithBorder prints a message with a decorative border for better readability.
func (p *cliMultiAgentChatImpl) printWithBorder(sender, message string) {
	p.printWithBorderColored(sender, message, MessageTypeNormal)
}

// printWithBorderColored prints a message with a decorative colored border for better readability.
func (p *cliMultiAgentChatImpl) printWithBorderColored(sender, message string, msgType MessageType) {
	// Use configurable width
	width := p.displayWidth

	// Get colors for this message type
	textColor, borderColor := p.getColorsForMessageType(msgType)

	// Top border
	fmt.Printf("%s╭%s╮%s\n", borderColor, strings.Repeat("─", width-2), ColorReset)

	// Sender line with bold text
	senderLine := fmt.Sprintf(" %s%s%s ", ColorBold, sender, ColorReset)
	cleanSender := stripansi.Strip(senderLine)
	senderPadding := width - runewidth.StringWidth(cleanSender) - 2
	if senderPadding < 0 {
		senderPadding = 0
	}
	fmt.Printf("%s│%s%s%s│%s\n", borderColor, senderLine, strings.Repeat(" ", senderPadding), borderColor, ColorReset)

	// Separator
	fmt.Printf("%s├%s┤%s\n", borderColor, strings.Repeat("─", width-2), ColorReset)

	// Message content
	renderedMessage := markdown.Render(message, p.displayWidth-4, 2)
	for _, line := range strings.Split(string(renderedMessage), "\n") {
		cleanLine := stripansi.Strip(line)
		padding := width - runewidth.StringWidth(cleanLine) - 4
		if padding < 0 {
			padding = 0
		}
		fmt.Printf("%s│ %s%s%s │%s\n", borderColor, textColor, line, strings.Repeat(" ", padding), ColorReset)
	}

	// Bottom border
	fmt.Printf("%s╰%s╯%s\n", borderColor, strings.Repeat("─", width-2), ColorReset)
	fmt.Println() // Extra line for spacing
}

// getColorsForMessageType returns the appropriate text and border colors for a message type.
func (p *cliMultiAgentChatImpl) getColorsForMessageType(msgType MessageType) (textColor, borderColor string) {
	switch msgType {
	case MessageTypeReasoningMessage:
		return ColorReasoning, ColorBorderReasoning
	case MessageTypeToolCall:
		return ColorTool, ColorBorderTool
	case MessageTypeIntercept:
		return ColorIntercept, ColorBorderIntercept
	case MessageTypeError:
		return ColorError, ColorError
	case MessageTypeAgentError:
		return ColorError, ColorError
	case MessageTypeSystem:
		return ColorSystem, ColorBorderTool
	default: // MessageTypeNormal
		return ColorNormal, ColorBorderNormal
	}
}

// shortenID safely shortens a UUID string to the first 8 characters for display purposes.
func shortenID(id string) string {
	if len(id) >= 8 {
		return id[:8]
	}
	return id
}
