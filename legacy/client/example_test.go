package client_test

import (
	"context"
	"fmt"
	"time"

	"github.com/denkhaus/agents/auth"
	"github.com/denkhaus/agents/client"
	"github.com/samber/mo"
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
)

// ExampleNewAPIClient demonstrates basic API client usage.
func ExampleNewAPIClient() {
	// Create a basic API client
	apiClient, err := client.NewAPIClient("http://localhost:8080",
		client.WithTimeout(30*time.Second))
	if err != nil {
		panic(err)
	}

	// Send a message
	ctx := context.Background()
	message := protocol.NewMessage(
		protocol.MessageRoleUser,
		[]protocol.Part{protocol.NewTextPart("Hello, world!")},
	)

	params := protocol.SendMessageParams{
		Message: message,
	}

	result, err := apiClient.SendMessage(ctx, params)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Message sent: %v\n", result)
}

// ExampleNewAPIClient_streaming demonstrates streaming using NewAPIClient with WithStreaming option.
func ExampleNewAPIClient_streaming() {
	// Create JWT authentication
	jwtAuth, _ := auth.NewJWTProvider(auth.JWTConfig{
		Secret:   []byte("my-secret-key"),
		Audience: "a2a-server",
		Issuer:   "example-client",
		Expiry:   time.Hour,
	})

	// Create client with streaming enabled
	apiClient, err := client.NewAPIClient("http://localhost:8080",
		client.WithAuth(jwtAuth),
		client.WithStreaming(true),
		client.WithTimeout(60*time.Second))
	if err != nil {
		panic(err)
	}

	// Type assert to StreamingClient interface to access streaming methods
	streamClient, ok := apiClient.(client.StreamingClient)
	if !ok {
		panic("client does not support streaming")
	}

	// Stream a message
	ctx := context.Background()
	message := protocol.NewMessage(
		protocol.MessageRoleUser,
		[]protocol.Part{protocol.NewTextPart("Stream this message")},
	)

	params := protocol.SendMessageParams{
		Message: message,
	}

	eventChan, err := streamClient.StreamMessage(ctx, params)
	if err != nil {
		panic(err)
	}

	// Process streaming events
	for event := range eventChan {
		fmt.Printf("Received event: %v\n", event)
		break // Just show first event for example
	}
}

// ExampleNewUnifiedClient demonstrates unified client usage.
func ExampleNewUnifiedClient() {
	// Create API key authentication
	apiKeyAuth, _ := auth.NewAPIKeyProvider(auth.APIKeyConfig{
		Key:    "my-api-key",
		Header: mo.Some("X-API-Key"),
	})

	// Create unified client configuration
	config := client.ClientConfig{
		BaseURL:    "http://localhost:8080",
		Timeout:    mo.Some(30 * time.Second),
		UserAgent:  mo.Some("unified-client-example/1.0"),
		Auth:       mo.Some(apiKeyAuth),
		Streaming:  mo.Some(true),
		RetryCount: mo.Some(3),
	}

	// Create unified client
	unifiedClient, err := client.NewUnifiedClient(config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// Fetch agent capabilities
	agentCard, err := unifiedClient.FetchAgentCard(ctx, "http://localhost:8080")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Agent: %s - %s\n", agentCard.Name, agentCard.Description)

	// Send a message
	message := protocol.NewMessage(
		protocol.MessageRoleUser,
		[]protocol.Part{protocol.NewTextPart("Unified client message")},
	)

	params := protocol.SendMessageParams{
		Message: message,
	}

	result, err := unifiedClient.SendMessage(ctx, params)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Message result: %v\n", result)

	// If streaming is supported and enabled, use streaming
	if agentCard.Capabilities.Streaming != nil && *agentCard.Capabilities.Streaming {
		if streamingClient, ok := unifiedClient.(client.StreamingClient); ok {
			eventChan, err := streamingClient.StreamMessage(ctx, params)
			if err == nil {
				for event := range eventChan {
					fmt.Printf("Stream event: %v\n", event)
					break // Just show first event for example
				}
			}
		}
	}
}

// ExampleTaskTracker demonstrates task tracking functionality.
func ExampleTaskTracker() {
	// Create task tracker
	tracker := client.NewTaskTracker()

	// Track a task
	taskID := "task-123"
	tracker.TrackTask(taskID)

	// Update task status
	tracker.UpdateTaskStatus(taskID, "running")
	tracker.UpdateTaskStatus(taskID, "completed")

	// Get task status
	status := tracker.GetTaskStatus(taskID)
	fmt.Printf("Task %s status: %s\n", taskID, status)
}

// ExampleJWKSVerifier demonstrates JWKS verification.
func ExampleJWKSVerifier() {
	// Create JWKS verifier
	verifier, _ := client.NewJWKSVerifier("http://localhost:8080/.well-known/jwks.json", 10*time.Second)

	// Refresh JWKS
	ctx := context.Background()
	err := verifier.RefreshJWKS(ctx)
	if err != nil {
		panic(err)
	}

	// Verify a JWT token
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9..." // Example token
	payload := []byte(`{"message": "test"}`)

	err = verifier.VerifyJWT(ctx, token, payload)
	if err != nil {
		fmt.Printf("JWT verification failed: %v\n", err)
	} else {
		fmt.Println("JWT verification successful")
	}
}
