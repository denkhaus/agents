package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/runner"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	sessioninmemory "trpc.group/trpc-go/trpc-agent-go/session/inmemory"

	messaging "github.com/denkhaus/agents/provider/agent"
)

// AgentRunner represents an AI agent with messaging capabilities
type AgentRunner struct {
	ID      uuid.UUID
	Runner  runner.Runner
	Wrapper *messaging.MessagingWrapper
	Name    string
}

// ChatSystem manages the multi-agent chat
type ChatSystem struct {
	broker *messaging.MessageBroker
	agents map[string]*AgentRunner
	human  *HumanAgent
}

// HumanAgent represents the human participant
type HumanAgent struct {
	ID   uuid.UUID
	Name string
}

// NewChatSystem creates a new chat system
func NewChatSystem() *ChatSystem {
	// Create human agent
	humanAgent := &HumanAgent{
		ID:   uuid.New(),
		Name: "Human",
	}

	broker := messaging.NewMessageBroker()

	system := &ChatSystem{
		broker: broker,
		agents: make(map[string]*AgentRunner),
		human:  humanAgent,
	}

	// Set up message listener to intercept agent-to-agent messages
	system.setupMessageListener()

	return system
}

// setupMessageListener sets up a listener to display agent-to-agent messages
func (cs *ChatSystem) setupMessageListener() {
	// Add a message interceptor to the broker
	cs.broker.SetMessageInterceptor(func(fromID, toID uuid.UUID, content string) {
		fromName := cs.getAgentNameByID(fromID)
		toName := cs.getAgentNameByID(toID)

		if fromName != "" && toName != "" {
			printWithBorder(fmt.Sprintf("%s -> %s", fromName, toName), content)
		}
	})
}

// getAgentNameByID returns the agent name for a given ID
func (cs *ChatSystem) getAgentNameByID(id uuid.UUID) string {
	// Check human
	if cs.human.ID == id {
		return cs.human.Name
	}

	// Check AI agents
	for _, agent := range cs.agents {
		if agent.ID == id {
			return agent.Name
		}
	}

	return ""
}

// startMessageProcessing starts a goroutine to process incoming messages for an agent
func (cs *ChatSystem) startMessageProcessing(agent *AgentRunner) {
	go func() {
		// Get the message channel for this agent
		msgChan, err := cs.broker.GetMessageChannel(agent.ID)
		if err != nil {
			log.Printf("Failed to get message channel for agent %s: %v", agent.Name, err)
			return
		}

		// Process incoming messages
		for msg := range msgChan {
			// Create a context for message processing
			ctx := context.Background()
			
			// Format the message content
			messageContent := fmt.Sprintf("Message from %s: %s", cs.getAgentNameByID(msg.From), msg.Content)
			
			// Send to the agent's runner
			events, err := agent.Runner.Run(ctx, msg.From.String(), fmt.Sprintf("msg-%s", msg.ID), model.NewUserMessage(messageContent))
			if err != nil {
				log.Printf("Error processing message for agent %s: %v", agent.Name, err)
				continue
			}

			// Process events from the agent's response
			go func() {
				for event := range events {
					processEvent(event)
				}
			}()
		}
	}()
}

// registerHumanWithBroker registers the human agent with the message broker
func (cs *ChatSystem) registerHumanWithBroker() {
	// Create a dummy agent for the human
	humanAgent := &dummyHumanAgent{id: cs.human.ID, name: cs.human.Name}
	cs.broker.RegisterAgentWithID(cs.human.ID, humanAgent)
}

// dummyHumanAgent implements the agent.Agent interface for the human
type dummyHumanAgent struct {
	id   uuid.UUID
	name string
}

func (d *dummyHumanAgent) Run(ctx context.Context, invocation *agent.Invocation) (<-chan *event.Event, error) {
	// Humans don't process messages automatically, just return empty channel
	ch := make(chan *event.Event)
	close(ch)
	return ch, nil
}

func (d *dummyHumanAgent) Info() agent.Info {
	return agent.Info{
		Name:        d.name,
		Description: "Human user",
	}
}

func (d *dummyHumanAgent) Tools() []tool.Tool {
	return []tool.Tool{}
}

func (d *dummyHumanAgent) FindSubAgent(name string) agent.Agent {
	return nil
}

func (d *dummyHumanAgent) SubAgents() []agent.Agent {
	return []agent.Agent{}
}

// CreateAgent creates an AI agent and adds it to the system
func (cs *ChatSystem) CreateAgent(appName, name, description, instruction string) error {
	// Get the pre-registered agent entry
	agentEntry, exists := cs.agents[strings.ToLower(name)]
	if !exists {
		return fmt.Errorf("agent %s not pre-registered", name)
	}

	// Create base agent with full instruction including agent list (excluding self)
	baseAgent := llmagent.New(
		name,
		llmagent.WithModel(openai.New("deepseek-chat")),
		llmagent.WithDescription(description),
		llmagent.WithInstruction(instruction+cs.generateAgentListForAgent(name)),
		llmagent.WithGenerationConfig(model.GenerationConfig{
			MaxTokens:   intPtr(500),
			Temperature: floatPtr(0.7),
			Stream:      false,
		}),
	)

	// Wrap with messaging using predefined ID
	wrapper := messaging.NewMessagingWrapperWithID(baseAgent, cs.broker, agentEntry.ID)

	// Create runner
	agentRunner := runner.NewRunner(
		appName,
		wrapper,
		runner.WithSessionService(sessioninmemory.NewSessionService()),
	)

	// Update the agent entry with actual components
	agentEntry.Runner = agentRunner
	agentEntry.Wrapper = wrapper

	// Start message processing for this agent
	cs.startMessageProcessing(agentEntry)

	return nil
}

// generateAgentList creates a list of known agents for instructions (excluding self)
func (cs *ChatSystem) generateAgentListForAgent(excludeAgentName string) string {
	result := "\n\nKnown agents you can message:\n"

	// Add human agent
	result += fmt.Sprintf("- %s (ID: %s) - Human user you can communicate with\n", cs.human.Name, cs.human.ID.String())

	// Add other AI agents (excluding self)
	for _, agent := range cs.agents {
		if agent.Name != excludeAgentName {
			result += fmt.Sprintf("- %s (ID: %s)\n", agent.Name, agent.ID.String())
		}
	}

	result += "Use the send_message tool to communicate with them.\n"
	return result
}

// SendMessage sends a message to an agent and returns events
func (cs *ChatSystem) SendMessage(ctx context.Context, agentName, message string) (<-chan *event.Event, error) {
	agent, exists := cs.agents[strings.ToLower(agentName)]
	if !exists {
		return nil, fmt.Errorf("agent '%s' not found", agentName)
	}

	userMessage := model.NewUserMessage(message)
	return agent.Runner.Run(ctx, "user", "session-"+agentName, userMessage)
}

// ListAgents returns a list of agent names
func (cs *ChatSystem) ListAgents() []string {
	var names []string
	// Add human first
	names = append(names, cs.human.Name)
	// Add AI agents
	for _, agent := range cs.agents {
		names = append(names, agent.Name)
	}
	return names
}

func main() {
	fmt.Println("🤖 Multi-Agent Chat System")
	fmt.Println("==========================")

	// Check API key
	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal("Please set OPENAI_API_KEY environment variable")
	}

	// Create system
	system := NewChatSystem()

	// Create agents in two phases: first register IDs, then create with full knowledge
	fmt.Println("Creating agents...")

	// Phase 1: Pre-register agent metadata
	agentMetadata := []struct {
		name, description, instruction string
	}{
		{
			"Coder",
			"Expert software engineer",
			"You are a skilled programmer who writes clean, efficient code. Help with coding tasks and collaborate with other agents when needed.",
		},
		{
			"Reviewer",
			"Expert code reviewer",
			"You are an experienced code reviewer. Analyze code for quality, security, and best practices. Collaborate with other agents when needed.",
		},
	}

	// Pre-create agent entries with IDs
	for _, meta := range agentMetadata {
		agentID := uuid.New()
		system.agents[strings.ToLower(meta.name)] = &AgentRunner{
			ID:   agentID,
			Name: meta.name,
		}
	}

	// Phase 2: Create actual agents with full knowledge
	for _, meta := range agentMetadata {
		err := system.CreateAgent(
			"multi-agent-chat",
			meta.name,
			meta.description,
			meta.instruction,
		)
		if err != nil {
			log.Fatal("Failed to create", meta.name+":", err)
		}
	}

	fmt.Println("✅ Agents created successfully!")
	fmt.Println("\nAvailable agents:", strings.Join(system.ListAgents(), ", "))
	fmt.Println("\nCommands:")
	fmt.Println("  @<agent> <message>  - Send message to agent")
	fmt.Println("  /list              - List agents")
	fmt.Println("  /exit              - Exit")
	fmt.Println()

	// Start chat loop
	startChat(system)
}

// startChat runs the interactive chat loop
func startChat(system *ChatSystem) {
	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()

	for {
		fmt.Print("you> ")
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
				fmt.Println("Available agents:", strings.Join(system.ListAgents(), ", "))
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

			agentName := parts[0]
			message := parts[1]

			// Send message and process events
			events, err := system.SendMessage(ctx, agentName, message)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			// Process events
			for event := range events {
				processEvent(event)
			}
			continue
		}

		fmt.Println("Use @<agent> <message> to send a message, or /help for commands")
	}
}

// processEvent handles a single event from an agent
func processEvent(event *event.Event) {
	if event.Error != nil {
		printWithBorder("ERROR", event.Error.Message)
		return
	}

	if event.Response != nil && len(event.Response.Choices) > 0 {
		choice := event.Response.Choices[0]

		// Show assistant messages
		if choice.Message.Role == model.RoleAssistant && choice.Message.Content != "" {
			printWithBorder(event.Author, choice.Message.Content)
		}

		// Show tool calls
		if len(choice.Message.ToolCalls) > 0 {
			for _, toolCall := range choice.Message.ToolCalls {
				if toolCall.Function.Name == "send_message" {
					printWithBorder(event.Author+" (Sending)", "Sending message to another agent...")
				} else {
					printWithBorder(event.Author+" (Tool)", fmt.Sprintf("Using tool: %s", toolCall.Function.Name))
				}
			}
		}

		// Show tool responses (messages from other agents)
		if choice.Message.Role == model.RoleTool && choice.Message.Content != "" {
			printWithBorder("Agent Message", choice.Message.Content)
		}
	}
}

// printWithBorder prints a message with a decorative border
func printWithBorder(sender, message string) {
	lines := strings.Split(message, "\n")
	maxLen := len(sender) + 4

	// Find the longest line for border width
	for _, line := range lines {
		if len(line)+4 > maxLen {
			maxLen = len(line) + 4
		}
	}

	// Limit maximum width to prevent overly wide displays
	if maxLen > 80 {
		maxLen = 80
	} else if maxLen < 50 {
		maxLen = 50
	}

	// Top border
	fmt.Printf("+")
	for i := 0; i < maxLen-2; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("+\n")

	// Sender line
	fmt.Printf("| %s", sender)
	for i := len(sender) + 2; i < maxLen-1; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("|\n")

	// Separator
	fmt.Printf("+")
	for i := 0; i < maxLen-2; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("+\n")

	// Message lines with word wrapping
	for _, line := range lines {
		// Word wrap long lines
		if len(line) > maxLen-4 {
			words := strings.Fields(line)
			currentLine := ""
			for _, word := range words {
				if len(currentLine)+len(word)+1 <= maxLen-4 {
					if currentLine != "" {
						currentLine += " "
					}
					currentLine += word
				} else {
					// Print current line and start new one
					if currentLine != "" {
						fmt.Printf("| %s", currentLine)
						for i := len(currentLine) + 2; i < maxLen-1; i++ {
							fmt.Printf(" ")
						}
						fmt.Printf("|\n")
					}
					currentLine = word
				}
			}
			// Print remaining line
			if currentLine != "" {
				fmt.Printf("| %s", currentLine)
				for i := len(currentLine) + 2; i < maxLen-1; i++ {
					fmt.Printf(" ")
				}
				fmt.Printf("|\n")
			}
		} else {
			fmt.Printf("| %s", line)
			for i := len(line) + 2; i < maxLen-1; i++ {
				fmt.Printf(" ")
			}
			fmt.Printf("|\n")
		}
	}

	// Bottom border
	fmt.Printf("+")
	for i := 0; i < maxLen-2; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("+\n\n")
}

// Helper functions
func intPtr(i int) *int           { return &i }
func floatPtr(f float64) *float64 { return &f }
