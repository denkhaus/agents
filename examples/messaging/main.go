package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/runner"
	sessioninmemory "trpc.group/trpc-go/trpc-agent-go/session/inmemory"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/messaging"
	"github.com/denkhaus/agents/shared"
)

type AgentRunner struct {
	Runner  runner.Runner
	Wrapper shared.TheAgent
	ID      uuid.UUID
	Name    string
}

type ChatSystem struct {
	applicationName string
	human           agent.Agent
	agents          map[string]*AgentRunner
	broker          messaging.MessageBroker
}

// NewChatSystem creates a new chat system
func NewChatSystem(applicationName string) *ChatSystem {
	// Create human agent
	humanAgent := shared.NewHumanAgent(shared.AgentInfoHuman)
	broker := messaging.NewMessageBroker()

	system := &ChatSystem{
		broker:          broker,
		applicationName: applicationName,
		agents:          make(map[string]*AgentRunner),
		human:           humanAgent,
	}

	// Set up message listener to intercept agent-to-agent messages
	system.setupMessageListener()
	system.broker.RegisterAgent(shared.AgentIDHuman, system.human)

	return system
}

// setupMessageListener sets up a listener to display agent-to-agent messages
func (cs *ChatSystem) setupMessageListener() {
	// Add a message interceptor to the broker
	cs.broker.SetMessageInterceptor(func(fromID, toID uuid.UUID, content string) {
		fromName := cs.getAgentNameByID(fromID)
		toName := cs.getAgentNameByID(toID)

		if fromName != "" && toName != "" {
			// Format: "FromName (FromID) -> ToName (ToID)"
			header := fmt.Sprintf("%s (%s) -> %s (%s)",
				fromName, shortenID(fromID.String()),
				toName, shortenID(toID.String()))
			printWithBorder(header, content)
		}
	})
}

// getAgentNameByID returns the agent name for a given ID
func (cs *ChatSystem) getAgentNameByID(id uuid.UUID) string {
	// Check human
	if shared.AgentIDHuman == id {
		return cs.human.Info().Name
	}

	// Check AI agents
	for _, agent := range cs.agents {
		if agent.ID == id {
			return agent.Name
		}
	}

	return ""
}

// getAgentInfoByAuthor returns agent name and ID by author string
func (cs *ChatSystem) getAgentInfoByAuthor(author string) (string, string) {
	// Try to parse as UUID first
	if authorID, err := uuid.Parse(author); err == nil {
		name := cs.getAgentNameByID(authorID)
		if name != "" {
			return name, authorID.String()
		}
	}

	// If not UUID or not found, check if it's already a name
	for _, agent := range cs.agents {
		if agent.Name == author {
			return agent.Name, agent.ID.String()
		}
	}

	if cs.human.Info().Name == author {
		return cs.human.Info().Name, shared.AgentIDHuman.String()
	}

	// Fallback: return as-is
	return author, author
}

// shortenID safely shortens a UUID string to first 8 characters
func shortenID(id string) string {
	if len(id) >= 8 {
		return id[:8] + "..."
	}
	return id
}

// startMessageProcessing starts a goroutine to process incoming messages for an agent
func (cs *ChatSystem) startMessageProcessing(agent *AgentRunner) {
	go func() {
		// Get the message channel for this agent
		msgChan, err := cs.broker.GetMessageChannel(agent.ID)
		if err != nil {
			logger.Log.Error("failed to get message channel for agent", zap.String("agent", agent.Name), zap.Error(err))
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
				logger.Log.Error("failed to process message for agent", zap.String("agent", agent.Name), zap.Error(err))
				continue
			}

			// Process events from the agent's response
			go func() {
				for event := range events {
					cs.processEvent(event)
				}
			}()
		}
	}()
}

// CreateAgent creates an AI agent and adds it to the system
func (cs *ChatSystem) CreateAgent(agentName, agentDescription, instruction string, agentID uuid.UUID) error {
	// Get the pre-registered agent entry
	agentEntry, exists := cs.agents[strings.ToLower(agentName)]
	if !exists {
		return fmt.Errorf("agent %s not pre-registered", agentName)
	}

	// Create base agent with full instruction including agent list (excluding self)
	baseAgent := llmagent.New(
		agentName,
		llmagent.WithModel(openai.New("deepseek-chat")),
		llmagent.WithDescription(agentDescription),
		llmagent.WithInstruction(instruction+cs.generateAgentListForAgent(agentName)),
		llmagent.WithGenerationConfig(model.GenerationConfig{
			MaxTokens:   intPtr(500),
			Temperature: floatPtr(0.7),
			Stream:      false,
		}),
	)

	// Wrap with messaging using predefined ID
	wrapper := messaging.NewMessagingWrapper(
		shared.NewAgent(baseAgent, agentID, false),
		cs.broker,
	)

	// Create runner
	agentRunner := runner.NewRunner(
		cs.applicationName,
		wrapper,
		runner.WithSessionService(
			sessioninmemory.NewSessionService(),
		),
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
	result += fmt.Sprintf("- %s (ID: %s) - Human user you can communicate with\n", cs.human.Info().Name, shared.AgentIDHuman.String())

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
	return agent.Runner.Run(ctx, agentName, "session-"+agentName, userMessage)
}

// ListAgents returns a list of agent names
func (cs *ChatSystem) ListAgents() []string {
	var names []string
	// Add human first
	names = append(names, cs.human.Info().Name)
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
	system := NewChatSystem("multi-agent-chat")

	// Create agents in two phases: first register IDs, then create with full knowledge
	fmt.Println("Creating agents...")

	// Phase 1: Pre-register agent metadata
	agentMetadata := []struct {
		name, description, instruction string
		agentID                        uuid.UUID
	}{
		{
			"Coder",
			"Expert software engineer",
			"You are a skilled programmer who writes clean, efficient code. Help with coding tasks and collaborate with other agents when needed.",
			uuid.New(),
		},
		{
			"Reviewer",
			"Expert code reviewer",
			"You are an experienced code reviewer. Analyze code for quality, security, and best practices. Collaborate with other agents when needed.",
			uuid.New(),
		},
	}

	// Pre-create agent entries with IDs
	for _, meta := range agentMetadata {
		system.agents[strings.ToLower(meta.name)] = &AgentRunner{
			ID:   meta.agentID,
			Name: meta.name,
		}
	}

	// Phase 2: Create actual agents with full knowledge
	for _, meta := range agentMetadata {
		err := system.CreateAgent(
			meta.name,
			meta.description,
			meta.instruction,
			meta.agentID,
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
				agents := system.ListAgents()
				fmt.Println("\n=== Available Agents ===")
				for _, agentName := range agents {
					if agentName == "Human" {
						fmt.Printf("- %s (ID: %s) [Type: Human]\n",
							system.human.Info().Name,
							shortenID(shared.AgentIDHuman.String()))
					} else {
						for _, agent := range system.agents {
							if agent.Name == agentName {
								fmt.Printf("- %s (ID: %s) [Type: AI]\n",
									agent.Name,
									shortenID(agent.ID.String()))
								break
							}
						}
					}
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

			agentName := parts[0]
			message := parts[1]

			// Send message and process events
			fmt.Printf(">> Sending message to %s...\n", agentName)
			events, err := system.SendMessage(ctx, agentName, message)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				continue
			}

			fmt.Printf(">> Message delivered to %s - Processing...\n", agentName)

			// Process events
			for event := range events {
				system.processEvent(event)
			}

			fmt.Printf(">> %s finished processing\n", agentName)
			continue
		}

		fmt.Println("Use @<agent> <message> to send a message, or /help for commands")
	}
}

// processEvent handles a single event from an agent
func (cs *ChatSystem) processEvent(event *event.Event) {
	if event.Error != nil {
		printWithBorder("ERROR", event.Error.Message)
		return
	}

	if event.Response != nil && len(event.Response.Choices) > 0 {
		choice := event.Response.Choices[0]

		// Show assistant messages
		if choice.Message.Role == model.RoleAssistant && choice.Message.Content != "" {
			// Get agent info for better display
			agentName, agentID := cs.getAgentInfoByAuthor(event.Author)
			header := fmt.Sprintf("%s (%s)", agentName, shortenID(agentID))
			printWithBorder(header, choice.Message.Content)
		}

		// Show tool calls (but suppress the generic "sending message" for cleaner output)
		if len(choice.Message.ToolCalls) > 0 {
			for _, toolCall := range choice.Message.ToolCalls {
				if toolCall.Function.Name != "send_message" {
					agentName, _ := cs.getAgentInfoByAuthor(event.Author)
					printWithBorder(agentName+" (Tool)", fmt.Sprintf("Using tool: %s", toolCall.Function.Name))
				}
			}
		}

		// Suppress tool responses (they're shown via message interceptor)
	}
}

// printWithBorder prints a message with a simple, clean border
func printWithBorder(sender, message string) {
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

// Helper functions
func intPtr(i int) *int           { return &i }
func floatPtr(f float64) *float64 { return &f }
