package plugins

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/denkhaus/agents/multi"
	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
)

// EnhancedBubbleTeaChatPlugin implements a modern TUI chat interface with real LLM calls
type EnhancedBubbleTeaChatPlugin struct {
	processor multi.ChatProcessor
	Options
}

// enhancedChatModel represents the Bubble Tea model
type enhancedChatModel struct {
	processor     multi.ChatProcessor
	agents        []shared.AgentInfo
	currentAgent  *shared.AgentInfo
	messages      []chatMessage
	input         string
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
	Type      MessageType
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
func NewEnhancedBubbleTeaChatPlugin(opts ...MultiAgentChatOption) ChatPlugin {
	chat := &EnhancedBubbleTeaChatPlugin{}

	// Apply options
	for _, opt := range opts {
		opt(&chat.Options)
	}

	// Create processor if not provided
	if chat.processor == nil {
		chat.processor = multi.NewChatProcessor(chat.processorOptions...)
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
		model.addMessage(fmt.Sprintf("%s -> %s", fromName, toName), content, MessageTypeIntercept)
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
			model.addMessage(info.Name, reasoning, MessageTypeReasoningMessage)

			// Stop spinner for this agent
			if spinner, exists := model.agentSpinners[info.ID().String()]; exists {
				spinner.Stop()
			}
			model.busyAgents[info.ID().String()] = false
		}),
		multi.WithOnError(func(info *shared.AgentInfo, err error) {
			model.addMessage(info.Name, fmt.Sprintf("Error: %v", err), MessageTypeError)

			// Stop spinner for this agent
			if spinner, exists := model.agentSpinners[info.ID().String()]; exists {
				spinner.Stop()
			}
			model.busyAgents[info.ID().String()] = false
		}),
		multi.WithOnProgress(func(messageType multi.SystemMessageType, format string, a ...any) {
			progressMsg := fmt.Sprintf(format, a...)
			model.addMessage("SYSTEM", progressMsg, MessageTypeSystem)
		}),
	}

	// Apply the processor options if not already set
	if p.processor == nil {
		p.processor = multi.NewChatProcessor(append(p.processorOptions, processorOptions...)...)
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
	msgType MessageType
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
		m.addMessage("SYSTEM", "Welcome to Multi-Agent Chat System!", MessageTypeSystem)
		m.addMessage("SYSTEM", "Available agents:", MessageTypeSystem)
		for _, agent := range m.agents {
			m.addMessage("SYSTEM", fmt.Sprintf("- %s (ID: %s)", agent.Name, agent.ID().String()), MessageTypeSystem)
		}
		m.addMessage("SYSTEM", "Use /<agent-name> to select an agent, /help for commands", MessageTypeSystem)
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
		m.addMessage(agentName, msg.error, MessageTypeError)
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

		case "enter":
			if m.input != "" {
				m.handleInput()
				m.input = ""
			}
			return m, nil

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
			return m, nil

		default:
			m.input += msg.String()
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

	content := strings.Join(items, "\n")
	return agentListStyle.Height(m.height - 8).Render(content)
}

// renderChatArea renders the chat messages
func (m *enhancedChatModel) renderChatArea() string {
	// Show last 20 messages
	start := 0
	if len(m.messages) > 20 {
		start = len(m.messages) - 20
	}

	var formattedMessages []string
	for _, msg := range m.messages[start:] {
		timestamp := msg.Timestamp.Format("15:04:05")

		// Apply color based on message type
		var style lipgloss.Style
		switch msg.Type {
		case MessageTypeReasoningMessage:
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("yellow"))
		case MessageTypeToolCall:
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("blue"))
		case MessageTypeIntercept:
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("magenta"))
		case MessageTypeSystem:
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("green"))
		case MessageTypeError:
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("red"))
		default:
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("white"))
		}

		formattedMsg := fmt.Sprintf("[%s] %s: %s", timestamp, msg.Agent, msg.Content)
		formattedMessages = append(formattedMessages, style.Render(formattedMsg))
	}

	content := strings.Join(formattedMessages, "\n")
	return chatAreaStyle.Width(m.width - 35).Height(m.height - 8).Render(content)
}

// renderInputArea renders the input field
func (m *enhancedChatModel) renderInputArea() string {
	prompt := "💬 "
	if m.currentAgent != nil {
		prompt = fmt.Sprintf("💬 [%s] ", m.currentAgent.Name)
	}

	input := fmt.Sprintf("%s%s", prompt, m.input)
	return inputStyle.Width(m.width - 4).Render(input)
}

// addMessage adds a message to the chat
func (m *enhancedChatModel) addMessage(agent, content string, msgType MessageType) {
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
}

// handleInput processes user input
func (m *enhancedChatModel) handleInput() {
	input := strings.TrimSpace(m.input)
	if input == "" {
		return
	}

	m.addMessage("YOU", input, MessageTypeNormal)

	// Handle commands
	if strings.HasPrefix(input, "/") {
		m.handleCommand(input[1:])
		return
	}

	// Send to current agent if selected
	if m.currentAgent != nil {
		m.sendToAgent(input)
	} else {
		m.addMessage("ERROR", "Please select an agent first using /<agent-name>", MessageTypeError)
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
			m.addMessage(m.currentAgent.Name, fmt.Sprintf("Error: %v", err), MessageTypeError)
			return
		}

		// Process events manually
		for event := range events {
			if event.Error != nil {
				m.addMessage(m.currentAgent.Name, fmt.Sprintf("Error: %v", event.Error.Message), MessageTypeError)
				continue
			}

			if event.Response != nil && len(event.Response.Choices) > 0 {
				choice := event.Response.Choices[0]

				// Handle reasoning content
				if choice.Message.ReasoningContent != "" {
					m.addMessage(m.currentAgent.Name, choice.Message.ReasoningContent, MessageTypeReasoningMessage)
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
						m.addMessage(m.currentAgent.Name+" [TOOL]", toolMsg, MessageTypeToolCall)
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
func (m *enhancedChatModel) detectMessageType(content string) MessageType {
	// Check for React planner tags that indicate reasoning/planning content
	if strings.Contains(content, "/PLANNING/") ||
		strings.Contains(content, "/REASONING/") ||
		strings.Contains(content, "/REPLANNING/") ||
		strings.Contains(content, "/*PLANNING*/") ||
		strings.Contains(content, "/*REASONING*/") ||
		strings.Contains(content, "/*REPLANNING*/") {
		return MessageTypeReasoningMessage
	}

	// Check for other reasoning indicators
	if strings.Contains(content, "/ACTION/") ||
		strings.Contains(content, "/*ACTION*/") {
		return MessageTypeReasoningMessage
	}

	return MessageTypeNormal
}

// handleCommand processes commands
func (m *enhancedChatModel) handleCommand(command string) {
	switch command {
	case "help":
		m.addMessage("HELP", "Available commands: /help, /list, /clear, /<agent-name>", MessageTypeSystem)
	case "list":
		m.addMessage("AGENTS", "Available agents:", MessageTypeSystem)
		for _, agent := range m.agents {
			m.addMessage("AGENTS", fmt.Sprintf("  - %s", agent.Name), MessageTypeSystem)
		}
	case "clear":
		m.currentAgent = nil
		m.addMessage("SYSTEM", "Agent selection cleared", MessageTypeSystem)
	default:
		// Try to select agent
		for _, agent := range m.agents {
			if agent.Name == command {
				m.currentAgent = &agent
				m.addMessage("SYSTEM", fmt.Sprintf("Selected agent: %s", agent.Name), MessageTypeSystem)
				return
			}
		}
		m.addMessage("ERROR", fmt.Sprintf("Unknown command or agent: %s", command), MessageTypeError)
	}
}
