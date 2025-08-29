package bubble

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/denkhaus/agents/multi"
	"github.com/denkhaus/agents/multi/plugins"
	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
)

// EnhancedBubbleTeaChatPlugin implements a modern TUI chat interface with real LLM calls
type EnhancedBubbleTeaChatPlugin struct {
	processor multi.ChatProcessor
	plugins.Options
}

// enhancedChatModel represents the Bubble Tea model
type enhancedChatModel struct {
	processor     multi.ChatProcessor
	agents        []shared.AgentInfo
	currentAgent  *shared.AgentInfo
	messages      []chatMessage
	input         string
	inputHistory  []string // Store previous user inputs
	historyIndex  int      // Current position in history (-1 = not navigating)
	scrollOffset  int      // Offset for scrolling through messages
	inputFocused  bool     // Whether input field has focus (for scroll control)
	busyAgents    map[string]bool
	agentSpinners map[string]*spinner.Spinner
	mainSpinner   *spinner.Spinner
	width         int
	height        int
	ready         bool
	ctx           context.Context
}

// chatMessage represents a chat message with metadata
type chatMessage struct {
	Agent     string
	Content   string
	Type      plugins.MessageType
	Timestamp time.Time
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	agentListStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2).
			Width(30)

	chatAreaStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#04B575")).
			Padding(0, 1)
)

// NewEnhancedBubbleTeaChatPlugin creates a new enhanced Bubble Tea chat plugin
func NewEnhancedBubbleTeaChatPlugin(opts ...plugins.MultiAgentChatOption) plugins.ChatPlugin {
	chat := &EnhancedBubbleTeaChatPlugin{}

	// Apply options
	for _, opt := range opts {
		opt(&chat.Options)
	}

	// Create processor if not provided
	if chat.processor == nil {
		chat.processor = multi.NewChatProcessor(chat.ProcessorOptions...)
	}

	return chat
}

// Start begins the Bubble Tea chat interface
func (p *EnhancedBubbleTeaChatPlugin) Start(ctx context.Context) error {
	// Create main spinner
	mainSpinner := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	mainSpinner.Suffix = " Initializing Multi-Agent Chat..."

	model := enhancedChatModel{
		processor:     p.processor,
		agents:        p.processor.GetAllAgentInfos(),
		messages:      []chatMessage{},
		inputHistory:  []string{},
		historyIndex:  -1,
		scrollOffset:  0,
		inputFocused:  true, // Start with input focused
		busyAgents:    make(map[string]bool),
		agentSpinners: make(map[string]*spinner.Spinner),
		mainSpinner:   mainSpinner,
		width:         120,
		height:        30,
		ctx:           ctx,
	}

	// Create individual spinners for each agent
	for _, agent := range model.agents {
		agentSpinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		agentSpinner.Suffix = fmt.Sprintf(" %s is thinking...", agent.Name)
		model.agentSpinners[agent.ID().String()] = agentSpinner
	}

	// Set up message handlers like in cli_multi_chat.go
	p.processor.SetMessageInterceptor(func(fromID, toID uuid.UUID, content string) {
		fromName := p.processor.GetAgentNameByID(fromID)
		toName := p.processor.GetAgentNameByID(toID)
		model.addMessage(fmt.Sprintf("%s -> %s", fromName, toName), content, plugins.MessageTypeIntercept)
	})

	// Set up callbacks for real LLM responses
	processorOptions := []multi.ChatProcessorOption{
		multi.WithOnMessage(func(info *shared.AgentInfo, content string) {
			msgType := model.detectMessageType(content)
			model.addMessage(info.Name, content, msgType)

			// Stop spinner for this agent
			if spinner, exists := model.agentSpinners[info.ID().String()]; exists {
				spinner.Stop()
			}
			model.busyAgents[info.ID().String()] = false
		}),
		multi.WithOnReasoningMessage(func(info *shared.AgentInfo, reasoning string) {
			model.addMessage(info.Name, reasoning, plugins.MessageTypeReasoningMessage)

			// Stop spinner for this agent
			if spinner, exists := model.agentSpinners[info.ID().String()]; exists {
				spinner.Stop()
			}
			model.busyAgents[info.ID().String()] = false
		}),
		multi.WithOnError(func(info *shared.AgentInfo, err error) {
			model.addMessage(info.Name, fmt.Sprintf("Error: %v", err), plugins.MessageTypeError)

			// Stop spinner for this agent
			if spinner, exists := model.agentSpinners[info.ID().String()]; exists {
				spinner.Stop()
			}
			model.busyAgents[info.ID().String()] = false
		}),
		multi.WithOnProgress(func(messageType multi.SystemMessageType, format string, a ...any) {
			progressMsg := fmt.Sprintf(format, a...)
			model.addMessage("SYSTEM", progressMsg, plugins.MessageTypeSystem)
		}),
	}

	// Apply the processor options if not already set
	if p.processor == nil {
		p.processor = multi.NewChatProcessor(append(p.ProcessorOptions, processorOptions...)...)
		model.processor = p.processor
	}

	// Start the Bubble Tea program
	program := tea.NewProgram(&model, tea.WithAltScreen())
	_, err := program.Run()
	return err
}

// Message types
type readyMsg struct{}
type tickMsg time.Time
type agentResponseMsg struct {
	agentID string
	content string
	msgType plugins.MessageType
}
type agentErrorMsg struct {
	agentID string
	error   string
}

// Init initializes the model
func (m *enhancedChatModel) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
		func() tea.Msg {
			return readyMsg{}
		},
	)
}

// Update handles messages
func (m *enhancedChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case readyMsg:
		m.ready = true
		
		// Create a single combined welcome message
		var welcomeMsg strings.Builder
		welcomeMsg.WriteString("Welcome to Multi-Agent Chat System!\n\n")
		welcomeMsg.WriteString("Available agents:\n")
		for _, agent := range m.agents {
			welcomeMsg.WriteString(fmt.Sprintf("- %s (ID: %s)\n", agent.Name, agent.ID().String()))
		}
		welcomeMsg.WriteString("\nUse /<agent-name> to select an agent, /help for commands")
		
		m.addMessage("SYSTEM", welcomeMsg.String(), plugins.MessageTypeSystem)
		return m, nil

	case agentResponseMsg:
		// Stop spinner for this agent
		if spinner, exists := m.agentSpinners[msg.agentID]; exists {
			spinner.Stop()
		}
		m.busyAgents[msg.agentID] = false

		// Get agent name
		agentName := m.processor.GetAgentNameByID(uuid.MustParse(msg.agentID))
		m.addMessage(agentName, msg.content, msg.msgType)
		return m, nil

	case agentErrorMsg:
		// Stop spinner for this agent
		if spinner, exists := m.agentSpinners[msg.agentID]; exists {
			spinner.Stop()
		}
		m.busyAgents[msg.agentID] = false

		// Get agent name
		agentName := m.processor.GetAgentNameByID(uuid.MustParse(msg.agentID))
		m.addMessage(agentName, msg.error, plugins.MessageTypeError)
		return m, nil

	case tickMsg:
		// Continue ticking for UI updates
		return m, tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			// Toggle focus between input and chat area
			m.inputFocused = !m.inputFocused
			return m, nil

		case "enter":
			if m.inputFocused && m.input != "" {
				m.handleInput()
				m.input = ""
				m.historyIndex = -1 // Reset history navigation
				m.scrollOffset = 0  // Reset scroll to show latest messages
			}
			return m, nil

		case "backspace":
			if m.inputFocused && len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.historyIndex = -1 // Reset history navigation when editing
			}
			return m, nil

		case "up":
			if m.inputFocused {
				// Navigate to previous message in input history
				if len(m.inputHistory) > 0 {
					if m.historyIndex == -1 {
						// First time navigating, start from the most recent
						m.historyIndex = len(m.inputHistory) - 1
					} else if m.historyIndex > 0 {
						// Go to older message
						m.historyIndex--
					}
					if m.historyIndex >= 0 && m.historyIndex < len(m.inputHistory) {
						m.input = m.inputHistory[m.historyIndex]
					}
				}
			} else {
				// Scroll up in chat area
				maxScroll := len(m.messages) - m.getVisibleMessageCount()
				if maxScroll > 0 && m.scrollOffset < maxScroll {
					m.scrollOffset++
				}
			}
			return m, nil

		case "down":
			if m.inputFocused {
				// Navigate to next message in input history
				if len(m.inputHistory) > 0 && m.historyIndex != -1 {
					if m.historyIndex < len(m.inputHistory)-1 {
						// Go to newer message
						m.historyIndex++
						m.input = m.inputHistory[m.historyIndex]
					} else {
						// Go back to empty input (newest)
						m.historyIndex = -1
						m.input = ""
					}
				}
			} else {
				// Scroll down in chat area
				if m.scrollOffset > 0 {
					m.scrollOffset--
				}
			}
			return m, nil

		case "pageup":
			// Page up in chat area (regardless of focus)
			pageSize := m.getVisibleMessageCount() / 2
			maxScroll := len(m.messages) - m.getVisibleMessageCount()
			if maxScroll > 0 {
				m.scrollOffset += pageSize
				if m.scrollOffset > maxScroll {
					m.scrollOffset = maxScroll
				}
			}
			return m, nil

		case "pagedown":
			// Page down in chat area (regardless of focus)
			pageSize := m.getVisibleMessageCount() / 2
			m.scrollOffset -= pageSize
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}
			return m, nil

		case "home":
			// Go to top of chat (oldest messages)
			maxScroll := len(m.messages) - m.getVisibleMessageCount()
			if maxScroll > 0 {
				m.scrollOffset = maxScroll
			}
			return m, nil

		case "end":
			// Go to bottom of chat (newest messages)
			m.scrollOffset = 0
			return m, nil

		default:
			// Only accept text input when input is focused
			if m.inputFocused {
				// Filter out unwanted keys (arrow keys, mouse wheel, etc.)
				key := msg.String()
				if len(key) == 1 || key == "space" || key == "tab" {
					if key == "space" {
						m.input += " "
					} else if key == "tab" {
						m.input += "\t"
					} else {
						m.input += key
					}
					m.historyIndex = -1 // Reset history navigation when typing
				}
			}
			return m, nil
		}
	}

	return m, nil
}

// View renders the interface
func (m *enhancedChatModel) View() string {
	if !m.ready {
		return fmt.Sprintf("\n%s\n", m.mainSpinner.Suffix)
	}

	// Title
	title := titleStyle.Render("🤖 Multi-Agent Chat System")

	// Agent list
	agentList := m.renderAgentList()

	// Chat area
	chatArea := m.renderChatArea()

	// Input area
	inputArea := m.renderInputArea()

	// Layout
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		agentList,
		chatArea,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		mainContent,
		inputArea,
	)
}

// renderAgentList renders the agent sidebar
func (m *enhancedChatModel) renderAgentList() string {
	var items []string
	items = append(items, "📋 Available Agents:")
	items = append(items, "")

	for _, agent := range m.agents {
		status := "[ ]"
		if m.busyAgents[agent.ID().String()] {
			status = "[⚡]" // Busy indicator
		} else if m.currentAgent != nil && m.currentAgent.ID() == agent.ID() {
			status = "[✓]" // Selected indicator
		}

		item := fmt.Sprintf("%s %s", status, agent.Name)
		items = append(items, item)
	}

	items = append(items, "")
	items = append(items, "Commands:")
	items = append(items, "/help - Show help")
	items = append(items, "/list - List agents")
	items = append(items, "/clear - Clear selection")
	items = append(items, "q - Quit")
	items = append(items, "")
	items = append(items, "Navigation:")
	items = append(items, "ESC - Toggle input/scroll mode")
	items = append(items, "↑↓ - Scroll/History")
	items = append(items, "PgUp/PgDn - Page scroll")
	items = append(items, "Home/End - Top/Bottom")

	content := strings.Join(items, "\n")
	return agentListStyle.Height(m.height - 8).Render(content)
}

// getVisibleMessageCount calculates how many messages can fit in the chat area
func (m *enhancedChatModel) getVisibleMessageCount() int {
	// Estimate based on chat area height - each message box takes roughly 6-8 lines
	// Being conservative with 8 lines per message
	chatHeight := m.height - 8 // Account for title and input area
	if chatHeight < 8 {
		return 1
	}
	return chatHeight / 8
}

// renderChatArea renders the chat messages with scrolling support
func (m *enhancedChatModel) renderChatArea() string {
	if len(m.messages) == 0 {
		content := "No messages yet..."
		return chatAreaStyle.Width(m.width - 35).Height(m.height - 8).Render(content)
	}

	// Calculate visible message range based on scroll offset
	visibleCount := m.getVisibleMessageCount()
	totalMessages := len(m.messages)
	
	// Calculate start index (from the end, accounting for scroll offset)
	start := totalMessages - visibleCount - m.scrollOffset
	if start < 0 {
		start = 0
	}
	
	// Calculate end index
	end := start + visibleCount
	if end > totalMessages {
		end = totalMessages
	}

	var formattedMessages []string
	for _, msg := range m.messages[start:end] {
		timestamp := msg.Timestamp.Format("15:04:05")
		
		// Create colored box for each message similar to CLI version
		boxedMessage := m.createColoredMessageBox(msg.Agent, msg.Content, msg.Type, timestamp)
		formattedMessages = append(formattedMessages, boxedMessage)
	}

	// Add scroll indicator if there are more messages
	var scrollInfo string
	if m.scrollOffset > 0 || start > 0 {
		scrollInfo = fmt.Sprintf(" [Showing %d-%d of %d messages | ESC to toggle focus, ↑↓ to scroll]", 
			start+1, end, totalMessages)
	}

	content := strings.Join(formattedMessages, "\n\n") // Extra spacing between boxes
	if scrollInfo != "" {
		content = scrollInfo + "\n\n" + content
	}
	
	return chatAreaStyle.Width(m.width - 35).Height(m.height - 8).Render(content)
}

// renderInputArea renders the input field
func (m *enhancedChatModel) renderInputArea() string {
	prompt := "💬 "
	statusText := ""
	
	// Show focus indicator
	focusIndicator := ""
	if m.inputFocused {
		focusIndicator = " [INPUT MODE]"
	} else {
		focusIndicator = " [SCROLL MODE - ESC to switch]"
	}
	
	if m.currentAgent != nil {
		prompt = fmt.Sprintf("💬 [%s] ", m.currentAgent.Name)
		
		// Check if current agent is busy and show spinner
		if m.busyAgents[m.currentAgent.ID().String()] {
			if spinner, exists := m.agentSpinners[m.currentAgent.ID().String()]; exists {
				statusText = fmt.Sprintf(" %s", spinner.Suffix)
			}
		}
	}

	// Combine input, status, and focus indicator
	inputContent := fmt.Sprintf("%s%s%s%s", prompt, m.input, statusText, focusIndicator)
	
	// Ensure proper width calculation to include bottom border
	inputWidth := m.width - 2 // Account for left and right margins
	if inputWidth < 20 {
		inputWidth = 20
	}
	
	// Use different border color based on focus
	style := inputStyle
	if !m.inputFocused {
		style = inputStyle.BorderForeground(lipgloss.Color("#666666")) // Dimmed when not focused
	}
	
	return style.Width(inputWidth).Render(inputContent)
}

// createColoredMessageBox creates a colored message box similar to CLI version
func (m *enhancedChatModel) createColoredMessageBox(agent, content string, msgType plugins.MessageType, timestamp string) string {
	// Get colors for message type
	textColor, borderColor := m.getColorsForMessageType(msgType)
	
	// Calculate available width for the chat area (accounting for agent list and margins)
	availableWidth := m.width - 35 // Same as chatAreaStyle width
	if availableWidth < 50 {
		availableWidth = 50
	}
	
	// Calculate content width (accounting for border, padding, and margins)
	contentWidth := availableWidth - 8 // Account for border (2) + padding (4) + margin (2)
	if contentWidth < 30 {
		contentWidth = 30
	}
	
	// Header section: Agent name and timestamp in bold
	headerText := fmt.Sprintf(" %s [%s] ", agent, timestamp)
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Width(contentWidth).
		Align(lipgloss.Left)
	
	// Separator line to match CLI version
	separatorStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Width(contentWidth)
	separator := separatorStyle.Render(strings.Repeat("─", contentWidth))
	
	// Content section with appropriate text color
	contentStyle := lipgloss.NewStyle().
		Foreground(textColor).
		Width(contentWidth).
		Padding(0, 1) // Small horizontal padding for content
	
	// Render components
	renderedHeader := headerStyle.Render(headerText)
	renderedContent := contentStyle.Render(content)
	
	// Combine all sections
	boxContent := lipgloss.JoinVertical(
		lipgloss.Left,
		renderedHeader,
		separator,
		renderedContent,
	)
	
	// Create the main box style with colored border - let lipgloss handle the width
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Padding(1, 1) // Add some internal padding
	
	return boxStyle.Render(boxContent)
}

// getColorsForMessageType returns appropriate colors for message types (lipgloss compatible)
func (m *enhancedChatModel) getColorsForMessageType(msgType plugins.MessageType) (textColor, borderColor lipgloss.Color) {
	switch msgType {
	case plugins.MessageTypeReasoningMessage:
		textColor = lipgloss.Color("#FFFF00")    // Yellow
		borderColor = lipgloss.Color("#FFFF87")  // Bright yellow
	case plugins.MessageTypeToolCall:
		textColor = lipgloss.Color("#5F87FF")    // Blue
		borderColor = lipgloss.Color("#87AFFF")  // Bright blue
	case plugins.MessageTypeIntercept:
		textColor = lipgloss.Color("#FF5FFF")    // Magenta
		borderColor = lipgloss.Color("#FF87FF")  // Bright magenta
	case plugins.MessageTypeError, plugins.MessageTypeAgentError:
		textColor = lipgloss.Color("#FF5F5F")    // Red
		borderColor = lipgloss.Color("#808080")  // Dark gray border
	case plugins.MessageTypeSystem:
		textColor = lipgloss.Color("#5FFF5F")    // Bright green
		borderColor = lipgloss.Color("#808080")  // Dark gray border
	default: // MessageTypeNormal
		textColor = lipgloss.Color("#FFFFFF")    // White
		borderColor = lipgloss.Color("#808080")  // Dark gray
	}
	return textColor, borderColor
}

// addMessage adds a message to the chat
func (m *enhancedChatModel) addMessage(agent, content string, msgType plugins.MessageType) {
	m.messages = append(m.messages, chatMessage{
		Agent:     agent,
		Content:   content,
		Type:      msgType,
		Timestamp: time.Now(),
	})

	// Keep only last 100 messages
	if len(m.messages) > 100 {
		m.messages = m.messages[1:]
	}
	
	// Auto-scroll to show new messages (reset scroll offset)
	m.scrollOffset = 0
}

// handleInput processes user input
func (m *enhancedChatModel) handleInput() {
	input := strings.TrimSpace(m.input)
	if input == "" {
		return
	}

	// Add to input history (avoid duplicates of the most recent entry)
	if len(m.inputHistory) == 0 || m.inputHistory[len(m.inputHistory)-1] != input {
		m.inputHistory = append(m.inputHistory, input)
		
		// Keep only last 50 entries in history
		if len(m.inputHistory) > 50 {
			m.inputHistory = m.inputHistory[1:]
		}
	}

	m.addMessage("YOU", input, plugins.MessageTypeNormal)

	// Handle commands
	if strings.HasPrefix(input, "/") {
		m.handleCommand(input[1:])
		return
	}

	// Send to current agent if selected
	if m.currentAgent != nil {
		m.sendToAgent(input)
	} else {
		m.addMessage("ERROR", "Please select an agent first using /<agent-name>", plugins.MessageTypeError)
	}
}

// sendToAgent sends a message to the selected agent using real LLM calls
func (m *enhancedChatModel) sendToAgent(message string) {
	if m.currentAgent == nil {
		return
	}

	agentID := m.currentAgent.ID().String()

	// Mark agent as busy and start spinner
	m.busyAgents[agentID] = true
	if spinner, exists := m.agentSpinners[agentID]; exists {
		spinner.Start()
	}

	// Send message using real ChatProcessor
	go func() {
		// Use the real SendMessage method instead to avoid callback issues
		events, err := m.processor.SendMessage(
			m.ctx,
			shared.AgentIDHuman, // From human
			m.currentAgent.ID(), // To selected agent
			message,
		)

		if err != nil {
			// Stop spinner and show error
			if spinner, exists := m.agentSpinners[agentID]; exists {
				spinner.Stop()
			}
			m.busyAgents[agentID] = false
			m.addMessage(m.currentAgent.Name, fmt.Sprintf("Error: %v", err), plugins.MessageTypeError)
			return
		}

		// Process events manually
		for event := range events {
			if event.Error != nil {
				m.addMessage(m.currentAgent.Name, fmt.Sprintf("Error: %v", event.Error.Message), plugins.MessageTypeError)
				continue
			}

			if event.Response != nil && len(event.Response.Choices) > 0 {
				choice := event.Response.Choices[0]

				// Handle reasoning content
				if choice.Message.ReasoningContent != "" {
					m.addMessage(m.currentAgent.Name, choice.Message.ReasoningContent, plugins.MessageTypeReasoningMessage)
				}

				// Handle normal content
				if choice.Message.Content != "" {
					msgType := m.detectMessageType(choice.Message.Content)
					m.addMessage(m.currentAgent.Name, choice.Message.Content, msgType)
				}

				// Handle tool calls
				if len(choice.Message.ToolCalls) > 0 {
					for _, toolCall := range choice.Message.ToolCalls {
						toolMsg := fmt.Sprintf("Tool Call: %s", toolCall.Function.Name)
						m.addMessage(m.currentAgent.Name+" [TOOL]", toolMsg, plugins.MessageTypeToolCall)
					}
				}
			}
		}

		// Stop spinner when done
		if spinner, exists := m.agentSpinners[agentID]; exists {
			spinner.Stop()
		}
		m.busyAgents[agentID] = false
	}()
}

// detectMessageType analyzes message content to determine the appropriate message type
func (m *enhancedChatModel) detectMessageType(content string) plugins.MessageType {
	// Check for React planner tags that indicate reasoning/planning content
	if strings.Contains(content, "/PLANNING/") ||
		strings.Contains(content, "/REASONING/") ||
		strings.Contains(content, "/REPLANNING/") ||
		strings.Contains(content, "/*PLANNING*/") ||
		strings.Contains(content, "/*REASONING*/") ||
		strings.Contains(content, "/*REPLANNING*/") {
		return plugins.MessageTypeReasoningMessage
	}

	// Check for other reasoning indicators
	if strings.Contains(content, "/ACTION/") ||
		strings.Contains(content, "/*ACTION*/") {
		return plugins.MessageTypeReasoningMessage
	}

	return plugins.MessageTypeNormal
}

// handleCommand processes commands
func (m *enhancedChatModel) handleCommand(command string) {
	switch command {
	case "help":
		m.addMessage("HELP", "Available commands: /help, /list, /clear, /<agent-name>", plugins.MessageTypeSystem)
	case "list":
		m.addMessage("AGENTS", "Available agents:", plugins.MessageTypeSystem)
		for _, agent := range m.agents {
			m.addMessage("AGENTS", fmt.Sprintf("  - %s", agent.Name), plugins.MessageTypeSystem)
		}
	case "clear":
		m.currentAgent = nil
		m.addMessage("SYSTEM", "Agent selection cleared", plugins.MessageTypeSystem)
	default:
		// Try to select agent
		for _, agent := range m.agents {
			if agent.Name == command {
				m.currentAgent = &agent
				m.addMessage("SYSTEM", fmt.Sprintf("Selected agent: %s", agent.Name), plugins.MessageTypeSystem)
				return
			}
		}
		m.addMessage("ERROR", fmt.Sprintf("Unknown command or agent: %s", command), plugins.MessageTypeError)
	}
}
