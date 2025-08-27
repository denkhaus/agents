# Agent Messaging System

This example demonstrates a generic messaging system that allows any `agent.Agent` implementation to communicate with other agents using unique IDs.

## Features

- **Generic Implementation**: Works with any `agent.Agent` implementation
- **ID-based Communication**: Agents are identified by UUIDs rather than names
- **Message Broker**: Centralized routing of messages between agents
- **Non-blocking Communication**: Asynchronous message passing with timeouts
- **Event Integration**: Messages are converted to events and merged with regular agent events
- **Tool Integration**: Messaging functionality is available as a tool that agents can use

## Components

### MessageBroker
The central component that routes messages between agents:
- Registers and unregisters agents
- Routes messages between agents by ID
- Manages message channels for each agent

### MessagingWrapper
A wrapper that adds messaging capabilities to any existing agent:
- Wraps any `agent.Agent` implementation
- Provides methods for sending messages to other agents
- Merges incoming messages with regular agent events
- Exposes messaging functionality as a tool

### Message
Represents a message between agents with sender, recipient, content, and timestamp.

## Usage

```go
// Create a message broker
broker := NewMessageBroker()

// Wrap existing agents with messaging capabilities
messagingAgent1 := NewMessagingWrapper(existingAgent1, broker)
messagingAgent2 := NewMessagingWrapper(existingAgent2, broker)

// Send a message from one agent to another
err := messagingAgent1.SendMessage(messagingAgent2.ID(), "Hello from Agent1!")

// Listen for messages
msgChan, err := messagingAgent2.GetMessageChannel()
```

## Key Benefits

1. **Decoupling**: Agents don't need to know about each other's implementations
2. **Flexibility**: Works with any agent type (LLM, Chain, Parallel, Cycle, etc.)
3. **Scalability**: Supports multiple agents with the same role through unique IDs
4. **Integration**: Seamlessly integrates with existing event-based workflows
5. **Tool Support**: Messaging functionality is available as a tool for agent use

## Example Output

```
Agent Messaging System Example
==============================
Backend Coder ID: df57524c-8a18-4a0c-b5e1-c2010860e0f5
Frontend Coder ID: 7292dfc1-ff83-4d12-bd0c-b0705bf0083d
Debugger ID: 316e3e4b-ba59-4f42-8fa4-aa1a77f1d3a8

Sending message from Backend Coder to Frontend Coder...
Sending message from Frontend Coder to Debugger...

Waiting for messages in Debugger...
Debugger received message from 7292dfc1-ff83-4d12-bd0c-b0705bf0083d: Frontend is integrated with the new API. Ready for testing.

Running Backend Coder agent...
Events from Backend Coder:
  Event from Backend-Coder: Hello, I'm Backend-Coder, a Backend Developer

Example completed successfully!
```