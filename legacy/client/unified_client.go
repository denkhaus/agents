package client

import (
	"context"
	"sync"
	"time"

	"github.com/samber/mo"
	"trpc.group/trpc-go/trpc-a2a-go/client"
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/server"
)

// unifiedClient implements all client interfaces.
type unifiedClient struct {
	baseURL    string
	timeout    time.Duration
	userAgent  string
	streaming  bool
	retryCount int
	auth       mo.Option[Authenticator]

	// Internal A2A client
	a2aClient *client.A2AClient

	// Thread safety
	mu sync.RWMutex
}

// SendMessage implements APIClient.
func (c *unifiedClient) SendMessage(ctx context.Context, params protocol.SendMessageParams) (*protocol.MessageResult, error) {
	client, err := c.getA2AClient()
	if err != nil {
		return nil, err
	}

	return client.SendMessage(ctx, params)
}

// GetTasks implements APIClient.
func (c *unifiedClient) GetTasks(ctx context.Context, params protocol.TaskQueryParams) (*protocol.Task, error) {
	client, err := c.getA2AClient()
	if err != nil {
		return nil, err
	}

	return client.GetTasks(ctx, params)
}

// CancelTasks implements APIClient.
func (c *unifiedClient) CancelTasks(ctx context.Context, params protocol.TaskIDParams) (*protocol.Task, error) {
	client, err := c.getA2AClient()
	if err != nil {
		return nil, err
	}

	return client.CancelTasks(ctx, params)
}

// SetPushNotification implements APIClient.
func (c *unifiedClient) SetPushNotification(ctx context.Context, config protocol.TaskPushNotificationConfig) (*protocol.TaskPushNotificationConfig, error) {
	client, err := c.getA2AClient()
	if err != nil {
		return nil, err
	}

	return client.SetPushNotification(ctx, config)
}

// GetPushNotification implements APIClient.
func (c *unifiedClient) GetPushNotification(ctx context.Context, params protocol.TaskIDParams) (*protocol.TaskPushNotificationConfig, error) {
	client, err := c.getA2AClient()
	if err != nil {
		return nil, err
	}

	return client.GetPushNotification(ctx, params)
}

// StreamMessage implements StreamingClient.
func (c *unifiedClient) StreamMessage(ctx context.Context, params protocol.SendMessageParams) (<-chan protocol.StreamingMessageEvent, error) {
	client, err := c.getA2AClient()
	if err != nil {
		return nil, err
	}

	return client.StreamMessage(ctx, params)
}

// FetchAgentCard implements AgentCardFetcher.
func (c *unifiedClient) FetchAgentCard(ctx context.Context, baseURL string) (*server.AgentCard, error) {
	fetcher := newAgentCardFetcher(c.timeout)
	return fetcher.FetchAgentCard(ctx, baseURL)
}

// getA2AClient returns the underlying A2A client, creating it if necessary.
func (c *unifiedClient) getA2AClient() (*client.A2AClient, error) {
	c.mu.RLock()
	if c.a2aClient != nil {
		defer c.mu.RUnlock()
		return c.a2aClient, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check pattern
	if c.a2aClient != nil {
		return c.a2aClient, nil
	}

	// Create client options
	var opts []client.Option
	opts = append(opts, client.WithTimeout(c.timeout))

	// Add authentication if configured
	if auth, ok := c.auth.Get(); ok {
		// Convert our Authenticator to A2A client options
		// This would need to be implemented based on the specific auth type
		_ = auth // TODO: Implement auth conversion
	}

	a2aClient, err := client.NewA2AClient(c.baseURL, opts...)
	if err != nil {
		return nil, err
	}

	c.a2aClient = a2aClient
	return c.a2aClient, nil
}
