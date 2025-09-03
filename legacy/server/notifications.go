// Package server provides push notification functionality extracted from examples.
package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/samber/mo"
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/taskmanager"
	"github.com/denkhaus/agents/helper"
)

// NotificationSender defines push notification capabilities.
type NotificationSender interface {
	// SendNotification sends a push notification
	SendNotification(ctx context.Context, config NotificationConfig, payload interface{}) error
	
	// WrapTaskManager wraps a task manager with notification capabilities
	WrapTaskManager(base taskmanager.TaskManager) taskmanager.TaskManager
}

// NotificationConfig defines notification configuration.
type NotificationConfig struct {
	URL   string
	Token mo.Option[string]
}

// notificationSender implements NotificationSender.
// Extracted from examples/basic/server/main.go
type notificationSender struct {
	httpClient helper.HTTPClient
}

// NewNotificationSender creates a new notification sender.
func NewNotificationSender(timeout mo.Option[time.Duration]) NotificationSender {
	timeoutValue := timeout.OrElse(10 * time.Second)
	return &notificationSender{
		httpClient: helper.NewHTTPClient(timeoutValue, 3),
	}
}

// SendNotification sends a push notification.
func (ns *notificationSender) SendNotification(ctx context.Context, config NotificationConfig, payload interface{}) error {
	// Create JSON payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, config.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Add authentication if configured
	if token, hasToken := config.Token.Get(); hasToken {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	// Send the request
	resp, err := ns.httpClient.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("notification failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// WrapTaskManager wraps a task manager with notification capabilities.
func (ns *notificationSender) WrapTaskManager(base taskmanager.TaskManager) taskmanager.TaskManager {
	return &notificationTaskManager{
		TaskManager: base,
		sender:      ns,
	}
}

// notificationTaskManager wraps a TaskManager to add webhook notifications.
type notificationTaskManager struct {
	taskmanager.TaskManager
	sender NotificationSender
}

// OnSendMessage overrides to add webhook notification.
func (ntm *notificationTaskManager) OnSendMessage(
	ctx context.Context,
	params protocol.SendMessageParams,
) (*protocol.MessageResult, error) {
	// Call the underlying implementation
	result, err := ntm.TaskManager.OnSendMessage(ctx, params)
	if err == nil && result != nil {
		// Check if result is a task and send notification if needed
		if task, ok := result.Result.(*protocol.Task); ok {
			go ntm.maybeSendStatusNotification(ctx, task.ID, task.Status.State)
		}
	}
	return result, err
}

// OnSendMessageStream overrides to add webhook notification.
func (ntm *notificationTaskManager) OnSendMessageStream(
	ctx context.Context,
	params protocol.SendMessageParams,
) (<-chan protocol.StreamingMessageEvent, error) {
	// Call the underlying implementation
	eventChan, err := ntm.TaskManager.OnSendMessageStream(ctx, params)
	if err != nil {
		return nil, err
	}

	// Create a wrapper channel to monitor events and send notifications
	wrappedChan := make(chan protocol.StreamingMessageEvent)

	go func() {
		defer close(wrappedChan)
		for event := range eventChan {
			// Forward the event
			wrappedChan <- event

			// Check if it's a task status update and send notification
			if statusEvent, ok := event.Result.(*protocol.TaskStatusUpdateEvent); ok {
				go ntm.maybeSendStatusNotification(ctx, statusEvent.TaskID, statusEvent.Status.State)
			}
		}
	}()

	return wrappedChan, nil
}

// OnCancelTask overrides to add webhook notification.
func (ntm *notificationTaskManager) OnCancelTask(
	ctx context.Context,
	params protocol.TaskIDParams,
) (*protocol.Task, error) {
	// Call the underlying implementation
	task, err := ntm.TaskManager.OnCancelTask(ctx, params)
	if err == nil && task != nil {
		// Send push notification if task was canceled successfully
		go ntm.maybeSendStatusNotification(ctx, task.ID, task.Status.State)
	}
	return task, err
}

// maybeSendStatusNotification sends a status notification if configured.
func (ntm *notificationTaskManager) maybeSendStatusNotification(
	ctx context.Context,
	taskID string,
	status protocol.TaskState,
) {
	// Get the push notification configuration for this task
	config, err := ntm.TaskManager.OnPushNotificationGet(
		ctx, protocol.TaskIDParams{ID: taskID},
	)
	if err != nil {
		// No configuration found or error occurred - no notification to send
		return
	}

	// Create notification payload
	payload := map[string]interface{}{
		"task_id":   taskID,
		"status":    status,
		"timestamp": time.Now().Unix(),
	}

	// Create notification config
	notificationConfig := NotificationConfig{
		URL:   config.PushNotificationConfig.URL,
		Token: mo.Some(config.PushNotificationConfig.Token),
	}

	// Send the notification
	if err := ntm.sender.SendNotification(ctx, notificationConfig, payload); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to send push notification: %v\n", err)
	}
}

// WebhookConfig defines webhook configuration for notifications.
type WebhookConfig struct {
	URL           string
	Token         mo.Option[string]
	RetryCount    mo.Option[int]
	RetryDelay    mo.Option[time.Duration]
	Timeout       mo.Option[time.Duration]
	EnableLogging mo.Option[bool]
}

// WebhookSender provides webhook-specific notification functionality.
type WebhookSender interface {
	// SendWebhook sends a webhook notification
	SendWebhook(ctx context.Context, config WebhookConfig, payload interface{}) error
	
	// SendWebhookWithRetry sends a webhook with retry logic
	SendWebhookWithRetry(ctx context.Context, config WebhookConfig, payload interface{}) error
}

// webhookSender implements WebhookSender.
type webhookSender struct {
	httpClient helper.HTTPClient
}

// NewWebhookSender creates a new webhook sender.
func NewWebhookSender() WebhookSender {
	return &webhookSender{
		httpClient: helper.NewHTTPClient(30*time.Second, 3),
	}
}

// SendWebhook sends a webhook notification.
func (ws *webhookSender) SendWebhook(ctx context.Context, config WebhookConfig, payload interface{}) error {
	// Create JSON payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, config.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "A2A-Server-Webhook/1.0")

	// Add authentication if configured
	if token, hasToken := config.Token.Get(); hasToken {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	// Apply timeout
	timeout := config.Timeout.OrElse(30 * time.Second)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Send the request
	resp, err := ws.httpClient.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Log success if enabled
	if config.EnableLogging.OrElse(false) {
		fmt.Printf("Webhook sent successfully to %s (status: %d)\n", config.URL, resp.StatusCode)
	}

	return nil
}

// SendWebhookWithRetry sends a webhook with retry logic.
func (ws *webhookSender) SendWebhookWithRetry(ctx context.Context, config WebhookConfig, payload interface{}) error {
	retryCount := config.RetryCount.OrElse(3)
	retryDelay := config.RetryDelay.OrElse(1 * time.Second)

	var lastErr error
	for attempt := 0; attempt <= retryCount; attempt++ {
		if attempt > 0 {
			// Apply retry delay with exponential backoff
			delay := time.Duration(attempt) * retryDelay
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		err := ws.SendWebhook(ctx, config, payload)
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on context cancellation
		if ctx.Err() != nil {
			break
		}
	}

	return fmt.Errorf("webhook failed after %d attempts: %w", retryCount+1, lastErr)
}