package cli

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync" // Add this import

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/acarl005/stripansi"
	"github.com/chzyer/readline"
	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/denkhaus/agents/pkg/multi"
	"github.com/denkhaus/agents/pkg/multi/plugins"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
	"github.com/mattn/go-runewidth"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// ANSI color codes for different message types
const (
	ColorReset = "\033[0m"
	ColorBold  = "\033[1m"

	// Message type colors
	ColorNormal    = "\033[37m" // White - normal messages
	ColorReasoning = "\033[33m" // Yellow - reasoning/planning messages
	ColorTool      = "\033[34m" // Blue - tool call messages
	ColorIntercept = "\033[35m" // Magenta - intercepted messages
	ColorError     = "\033[31m" // Red - error messages
	ColorSystem    = "\033[92m" // Bright Green - system messages

	// Border colors
	ColorBorderNormal    = "\033[90m" // Dark gray
	ColorBorderReasoning = "\033[93m" // Bright yellow
	ColorBorderTool      = "\033[94m" // Bright Blue
	ColorBorderIntercept = "\033[95m" // Bright magenta
)

// ChatSystem manages the multi-agent chat
type cliMultiAgentChatImpl struct {
	plugins.Options
	currentAgent *shared.AgentInfo  // Track the currently selected agent
	outputMutex  sync.Mutex         // Mutex to protect concurrent writes to stdout
	activeCancel context.CancelFunc // Cancel function for active agent operations
	activeMutex  sync.Mutex         // Mutex to protect activeCancel
	isAgentBusy  bool               // Track if an agent is currently processing
}

// NewCLIMultiAgentChat creates a new CLI-based multi-agent chat plugin.
// It sets up the chat processor with the provided options and configures message handling.
func NewCLIMultiAgentChat(opts ...plugins.MultiAgentChatOption) (plugins.ChatPlugin, error) {
	chat := &cliMultiAgentChatImpl{
		Options: plugins.Options{
			SessionID:    uuid.New(),
			DisplayWidth: 120, // Default width
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
		multi.WithOnRawEvent(chat.handleOnRawEvent),
	}

	processorOptions = append(processorOptions, chat.ProcessorOptions...)

	var err error
	chat.Processor, err = multi.NewChatProcessor(chat.SessionID, processorOptions...)
	if err != nil {
		return nil, err
	}

	chat.setupMessageListener()
	return chat, nil
}

// setupMessageListener configures the message interceptor to display inter-agent communication.
func (p *cliMultiAgentChatImpl) setupMessageListener() {
	// Add a message interceptor to the broker
	p.Processor.SetMessageInterceptor(func(routing *messaging.RoutingInfo, content string) {
		fromName := p.Processor.GetAgentNameByID(routing.FromAgentID)
		toName := p.Processor.GetAgentNameByID(routing.ToAgentID)

		if fromName != "" && toName != "" {
			// Format: "FromName (FromID) -> ToName (ToID)"
			header := fmt.Sprintf("%s (%s) -> %s (%s)",
				fromName, routing.FromAgentID, toName, routing.ToAgentID,
			)
			p.printWithBorderColored(header, content, plugins.MessageTypeIntercept)
		}
	})
}

func (p *cliMultiAgentChatImpl) handleOnRawEvent(info *messaging.RoutingInfo, event *event.Event) {

}

// handleOnProgress handles progress updates by printing them to stdout.
func (p *cliMultiAgentChatImpl) handleOnProgress(routing *messaging.RoutingInfo, messageType multi.SystemMessageType, format string, a ...any) {
	p.printSystemMessage(format, a...)
}

// handleOnMessage handles agent messages by displaying them with a formatted border.
func (p *cliMultiAgentChatImpl) handleOnMessage(routing *messaging.RoutingInfo, content string) {
	// Detect if this is a reasoning/planning message based on content
	msgType := p.detectMessageType(content)
	p.printWithBorderColored(routing.String(), content, msgType)
}

// handleOnError handles agent errors by displaying them with a formatted border.
func (p *cliMultiAgentChatImpl) handleOnError(routing *messaging.RoutingInfo, err error) {
	p.printWithBorderColored(routing.String(), fmt.Sprintf("%+v", err), plugins.MessageTypeAgentError)
}

// handleOnToolCall handles tool calls made by agents by displaying them with a formatted border.
func (p *cliMultiAgentChatImpl) handleOnToolCall(routing *messaging.RoutingInfo, functionDef model.FunctionDefinitionParam) {
	toolCallInfo := fmt.Sprintf("Tool Call: %s", functionDef.Name)
	if len(functionDef.Arguments) > 0 {
		toolCallInfo += fmt.Sprintf("\nArguments: %s", string(functionDef.Arguments))
	}
	p.printWithBorderColored(routing.String()+" [TOOL]", toolCallInfo, plugins.MessageTypeToolCall)
}

// handleOnReasoningMessage handles reasoning messages from agents.
func (p *cliMultiAgentChatImpl) handleOnReasoningMessage(routing *messaging.RoutingInfo, reasoning string) {
	p.printWithBorderColored(routing.String(), reasoning, plugins.MessageTypeReasoningMessage)
}

// Start runs the interactive chat loop, handling user input and agent communication.
// It supports commands like /exit, /list, /agent-name to select agents, and direct messaging.
func (p *cliMultiAgentChatImpl) Start(ctx context.Context) error {
	// Show welcome message with all available commands
	p.showWelcomeMessage()

	// Create readline instance with cursor support and tab completion
	completer := readline.NewPrefixCompleter(
		readline.PcItem("/exit"),
		readline.PcItem("/help"),
		readline.PcItem("/list"),
		readline.PcItem("/clear"),
		readline.PcItem("/width"),
	)

	// Add agent names to completer
	for _, agent := range p.Processor.GetAllAgentInfos() {
		if agent.Role != shared.AgentRoleHuman {
			completer.Children = append(completer.Children, readline.PcItem("/"+agent.Name))
		}
	}

	config := &readline.Config{
		Prompt:          "",
		HistoryFile:     "/tmp/.agents_history",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		return fmt.Errorf("failed to create readline: %w", err)
	}
	defer rl.Close()

	for {
		// Create prompt showing current agent
		prompt := "you"
		if p.currentAgent != nil {
			prompt = fmt.Sprintf("you [%s]", p.currentAgent.Name)
		}
		rl.SetPrompt(fmt.Sprintf("%s >> ", prompt))

		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				// Handle Ctrl+C - exit the application
				p.printSystemMessage("Goodbye!")
				return nil
			} else if err == io.EOF {
				break
			}
			return err
		}

		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}

		// Handle ESC key for agent interruption
		if input == "\x1b" || strings.HasPrefix(input, "\x1b") {
			// ESC key pressed - check if we should cancel agent operation
			p.activeMutex.Lock()
			if p.isAgentBusy && p.activeCancel != nil {
				p.printSystemMessage("[ESC] Cancelling agent operation...")
				p.activeCancel()
				p.activeCancel = nil
				p.isAgentBusy = false
				p.activeMutex.Unlock()
			} else {
				p.activeMutex.Unlock()
				p.printSystemMessage("No active agent operation to cancel")
			}
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
				for _, info := range p.Processor.GetAllAgentInfos() {
					marker := ""
					if p.currentAgent != nil && p.currentAgent.Equal(info) {
						marker = " (current)"
					}
					// Mark human agents as non-selectable
					if info.Role == shared.AgentRoleHuman {
						builder.WriteString(fmt.Sprintf("- %s (ID: %s) [HUMAN - not selectable]%s\n", info.Name, info.ID, marker))
					} else {
						builder.WriteString(fmt.Sprintf("- %s (ID: %s)%s\n", info.Name, info.ID, marker))
					}
				}
				builder.WriteString("=========================")
				p.printSystemText(builder.String())
			case "clear":
				p.currentAgent = nil
				p.printSystemMessage("Current agent cleared. Use /<agent-name> to select an agent.")
			case "help":
				p.printSystemText(p.getHelpMessage())
			default:
				// Check if it's a width command
				if strings.HasPrefix(command, "width ") {
					widthStr := strings.TrimPrefix(command, "width ")
					if width, err := strconv.Atoi(widthStr); err == nil {
						if width < 40 {
							fmt.Println("Minimum width is 40 characters.")
							width = 40
						}
						p.DisplayWidth = width
						p.printSystemMessage("Display width set to %d characters.", width)
					} else {
						fmt.Println("Invalid width. Usage: /width <number>")
					}
					continue
				}
				// Try to find agent by name
				agentInfo := p.Processor.GetAgentInfoByAuthor(command)
				if agentInfo != nil {
					// Prevent selecting the human agent as a target
					if agentInfo.Role == shared.AgentRoleHuman {
						p.printSystemMessage("Cannot select human agent '%s' as a target. You can only send messages to AI agents.", agentInfo.Name)
					} else {
						p.currentAgent = agentInfo
						p.printSystemMessage("Selected agent: %s", agentInfo.Name)
					}
				} else {
					p.printSystemMessage("Unknown command or agent: %s. Use /help for available commands.", command)
				}
			}
			continue
		}

		// Send message to current agent or show help
		if p.currentAgent != nil {
			p.sendMessageToAgent(ctx, input)
		} else {
			fmt.Println("No agent selected. Use /<agent-name> to select an agent.")
		}
	}

	return nil
}

// printSystemMessage displays a system message with a standard border.
func (p *cliMultiAgentChatImpl) printSystemMessage(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	p.printWithBorderColored("SYSTEM", message, plugins.MessageTypeSystem)
}

// printSystemText displays a pre-formatted system text without additional formatting.
// Use this for already formatted content like markdown or when the text might contain % characters.
func (p *cliMultiAgentChatImpl) printSystemText(text string) {
	p.printWithBorderColored("SYSTEM", text, plugins.MessageTypeSystem)
}

// showWelcomeMessage displays a welcome message with all available commands and agents.
func (p *cliMultiAgentChatImpl) showWelcomeMessage() {
	agents := p.Processor.GetAllAgentInfos()
	welcomeMarkdown := GetWelcomeMessage(agents, p.DisplayWidth)
	p.printSystemText(welcomeMarkdown) // Use printSystemText for pre-formatted content
}

// getHelpMessage returns the help message as a string.
func (p *cliMultiAgentChatImpl) getHelpMessage() string {
	var builder strings.Builder
	builder.WriteString("\n=== Available Commands ===\n")
	builder.WriteString("/help                 - Show this help message\n")
	builder.WriteString("/list                 - List all available agents\n")
	builder.WriteString("/clear                - Clear current agent selection\n")
	builder.WriteString(fmt.Sprintf("/width <number>       - Set display width (min: 40, current: %d)\n", p.DisplayWidth))
	builder.WriteString("/<agent-name>         - Select an agent to chat with\n")
	builder.WriteString("/exit                 - Exit the chat\n")
	builder.WriteString("\n")
	builder.WriteString("=== Navigation & Control ===\n")
	builder.WriteString("Arrow Keys            - Navigate cursor and command history\n")
	builder.WriteString("Tab                   - Auto-complete commands and agent names\n")
	builder.WriteString("ESC                   - Interrupt active agent processing\n")
	builder.WriteString("Ctrl+C                - Exit the chat\n")
	builder.WriteString("Ctrl+D                - Alternative exit\n")
	builder.WriteString("\n")
	builder.WriteString("=== Usage ===\n")
	builder.WriteString("1. Select an agent: /project-manager\n")
	builder.WriteString("2. Chat directly: Hello, how can you help?\n")
	builder.WriteString("3. Switch agents: /another-agent\n")
	builder.WriteString("4. Interrupt processing: ESC (during agent response)\n")
	builder.WriteString("5. Adjust display: /width 80\n")
	builder.WriteString("===========================")
	return builder.String()
}

// printWithBorderColored prints a message with a decorative colored border for better readability.
func (p *cliMultiAgentChatImpl) printWithBorderColored(sender, message string, msgType plugins.MessageType) {
	p.outputMutex.Lock()         // Acquire lock
	defer p.outputMutex.Unlock() // Release lock when function exits

	// Use configurable width
	width := p.DisplayWidth

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
	renderedMessage := markdown.Render(message, p.DisplayWidth-4, 2)
	for _, line := range strings.Split(string(renderedMessage), "\n") {
		// Apply textColor to the line *after* markdown rendering, and ensure it's reset
		coloredLine := fmt.Sprintf("%s%s%s", textColor, line, ColorReset)

		// Calculate padding based on the *clean* line (without ANSI codes)
		cleanLine := stripansi.Strip(line)
		padding := width - runewidth.StringWidth(cleanLine) - 4
		if padding < 0 {
			padding = 0
		}

		// Print the line with correct padding and border colors
		fmt.Printf("%s│%s %s%s %s│%s\n", borderColor, ColorReset, coloredLine, strings.Repeat(" ", padding), borderColor, ColorReset)
	}

	// Bottom border
	fmt.Printf("%s╰%s╯%s\n", borderColor, strings.Repeat("─", width-2), ColorReset)
	fmt.Println() // Extra line for spacing
}

// getColorsForMessageType returns the appropriate text and border colors for a message type.
func (p *cliMultiAgentChatImpl) getColorsForMessageType(msgType plugins.MessageType) (textColor, borderColor string) {
	switch msgType {
	case plugins.MessageTypeReasoningMessage:
		textColor = ColorReasoning
		borderColor = ColorBorderReasoning
	case plugins.MessageTypeToolCall:
		textColor = ColorTool
		borderColor = ColorBorderTool
	case plugins.MessageTypeIntercept:
		textColor = ColorIntercept
		borderColor = ColorBorderIntercept
	case plugins.MessageTypeError:
		textColor = ColorError
		borderColor = ColorBorderNormal // Keep normal border for errors
	case plugins.MessageTypeAgentError:
		textColor = ColorError
		borderColor = ColorBorderNormal // Keep normal border for errors
	case plugins.MessageTypeSystem:
		textColor = ColorSystem
		borderColor = ColorBorderNormal // Keep normal border for system messages
	default: // MessageTypeNormal
		textColor = ColorNormal
		borderColor = ColorBorderNormal
	}
	return textColor, borderColor
}

// detectMessageType analyzes message content to determine the appropriate message type
func (p *cliMultiAgentChatImpl) detectMessageType(content string) plugins.MessageType {
	// Check for React planner tags that indicate reasoning/planning content
	if strings.Contains(content, "/PLANNING/") ||
		strings.Contains(content, "/REASONING/") ||
		strings.Contains(content, "/REPLANNING/") ||
		strings.Contains(content, "/*PLANNING*/") ||
		strings.Contains(content, "/*REASONING*/") ||
		strings.Contains(content, "/*REPLANNING*/") {
		fmt.Printf("[DEBUG] Detected reasoning content based on planning tags\n")
		return plugins.MessageTypeReasoningMessage
	}

	// Check for other reasoning indicators
	if strings.Contains(content, "/ACTION/") ||
		strings.Contains(content, "/*ACTION*/") {
		// ACTION sections are still part of reasoning flow
		fmt.Printf("[DEBUG] Detected reasoning content based on action tags\n")
		return plugins.MessageTypeReasoningMessage
	}

	return plugins.MessageTypeNormal
}

// sendMessageToAgent sends a message to the current agent with cancellation support
func (p *cliMultiAgentChatImpl) sendMessageToAgent(ctx context.Context, input string) {
	if p.currentAgent == nil {
		return
	}

	// Create a cancellable context for this operation
	agentCtx, cancel := context.WithCancel(ctx)

	// Store the cancel function so it can be called by interrupt handler
	p.activeMutex.Lock()
	p.activeCancel = cancel
	p.isAgentBusy = true
	p.activeMutex.Unlock()

	// Show that agent is processing
	p.printSystemMessage("[PROCESSING] %s is processing your message... (Press ESC to interrupt)", p.currentAgent.Name)

	// Send message in a goroutine to allow interruption
	go func() {
		defer func() {
			// Clean up when done
			p.activeMutex.Lock()
			p.activeCancel = nil
			p.isAgentBusy = false
			p.activeMutex.Unlock()
		}()

		routing := &messaging.RoutingInfo{
			FromAgentID: shared.AgentIDHuman,
			ToAgentID:   p.currentAgent.ID,
			SessionID:   p.SessionID,
			Streaming:   p.currentAgent.IsStreaming,
		}
		// Use SendMessageWithProcessing but with cancellable context
		userMessage := model.NewUserMessage(input)
		err := p.Processor.SendMessageWithProcessing(
			agentCtx,
			routing,
			userMessage,
		)

		if err != nil {
			if agentCtx.Err() == context.Canceled {
				p.printSystemMessage("[CANCELLED] Agent operation was cancelled")
			} else {
				p.printSystemMessage("ERROR: %v", err)
			}
			return
		}

		p.printSystemMessage("[COMPLETED] %s finished processing", p.currentAgent.Name)
	}()
}
