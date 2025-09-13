package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
)

func main() {
	// Create some sample agent info
	agentID1 := uuid.New()
	agentID2 := uuid.New()
	callingAgentID := agentID1 // The first agent is calling the tool

	// Create agent info instances
	agent1 := shared.NewAgentInfo(agentID1, shared.AgentRoleCoder, false, "CodeBot", "AI agent specialized in writing and reviewing code")
	agent2 := shared.NewAgentInfo(agentID2, shared.AgentRoleProjectManager, true, "ProjectManager", "AI agent for project management and coordination")

	availableAgents := []*shared.AgentInfo{&agent1, &agent2}

	// Create the agent info tool
	tool, err := messaging.NewAgentInfoTool(availableAgents, callingAgentID)
	if err != nil {
		fmt.Printf("Error creating tool: %v\n", err)
		return
	}

	// Test 1: Get all agents
	fmt.Println("=== Test 1: Get all agents ===")
	args1 := `{}`
	result1, err := tool.Call(context.Background(), []byte(args1))
	if err != nil {
		fmt.Printf("Error calling tool: %v\n", err)
		return
	}
	
	jsonResult1, _ := json.MarshalIndent(result1, "", "  ")
	fmt.Printf("Result: %s\n\n", jsonResult1)

	// Test 2: Get specific agent (self)
	fmt.Println("=== Test 2: Get specific agent (self) ===")
	args2 := fmt.Sprintf(`{"agent_id": "%s"}`, agentID1.String())
	result2, err := tool.Call(context.Background(), []byte(args2))
	if err != nil {
		fmt.Printf("Error calling tool: %v\n", err)
		return
	}
	
	jsonResult2, _ := json.MarshalIndent(result2, "", "  ")
	fmt.Printf("Result: %s\n\n", jsonResult2)

	// Test 3: Get specific agent (other)
	fmt.Println("=== Test 3: Get specific agent (other) ===")
	args3 := fmt.Sprintf(`{"agent_id": "%s"}`, agentID2.String())
	result3, err := tool.Call(context.Background(), []byte(args3))
	if err != nil {
		fmt.Printf("Error calling tool: %v\n", err)
		return
	}
	
	jsonResult3, _ := json.MarshalIndent(result3, "", "  ")
	fmt.Printf("Result: %s\n\n", jsonResult3)
}