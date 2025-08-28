package plugins

import (
	"fmt"

	"github.com/denkhaus/agents/shared"
)

// GetWelcomeMessage returns a properly formatted welcome message as markdown
func GetWelcomeMessage(agents []shared.AgentInfo, displayWidth int) string {
	// Build agent list dynamically
	agentList := ""
	if len(agents) > 0 {
		for _, info := range agents {
			agentList += fmt.Sprintf("- **%s** (ID: %s)\n", info.Name, info.ID())
		}
	} else {
		agentList = "- No agents available\n"
	}

	// Use template with direct markdown (using string concatenation to avoid backtick issues)
	template := "# Welcome to Multi-Agent Chat System!\n\n" +
		"## Available Agents\n\n" +
		"%s\n" +
		"## Available Commands\n\n" +
		"- `/help` - Show help message\n" +
		"- `/list` - List all available agents\n" +
		"- `/clear` - Clear current agent selection\n" +
		"- `/width <number>` - Set display width (min: 40, current: %d)\n" +
		"- `/<agent-name>` - Select an agent to chat with\n" +
		"- `/exit` - Exit the chat\n\n" +
		"## Quick Start\n\n" +
		"1. Select an agent: `/project-manager`\n" +
		"2. Start chatting: Hello, how can you help?\n" +
		"3. Switch agents anytime: `/another-agent`\n" +
		"4. Get help: `/help`\n\n" +
		"## Message Types\n\n" +
		"- **Yellow boxes**: Reasoning/Planning messages\n" +
		"- **Blue boxes**: Tool calls and actions\n" +
		"- **Purple boxes**: Inter-agent communication\n" +
		"- **White boxes**: Normal responses\n" +
		"- **Green boxes**: System messages\n\n" +
		"---\n\n" +
		"**Ready to chat! Select an agent to get started.**"

	return fmt.Sprintf(template, agentList, displayWidth)
}
