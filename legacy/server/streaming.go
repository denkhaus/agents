// Package server provides streaming functionality extracted from examples.
package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/samber/mo"
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/taskmanager"
)

// StreamingHandler defines streaming message processing.
type StreamingHandler interface {
	// ProcessStream processes a message with streaming support
	ProcessStream(
		ctx context.Context,
		message protocol.Message,
		options taskmanager.ProcessOptions,
		handle taskmanager.TaskHandler,
	) (*taskmanager.MessageProcessingResult, error)
	
	// SupportsStreaming returns true if streaming is supported
	SupportsStreaming() bool
}

// streamingHandler implements StreamingHandler.
// Extracted from examples/streaming/server/main.go
type streamingHandler struct {
	chunkSize    int
	processDelay time.Duration
}

// NewStreamingHandler creates a new streaming handler.
func NewStreamingHandler(chunkSize mo.Option[int], processDelay mo.Option[time.Duration]) StreamingHandler {
	return &streamingHandler{
		chunkSize:    chunkSize.OrElse(5),
		processDelay: processDelay.OrElse(100 * time.Millisecond),
	}
}

// ProcessStream processes a message with streaming support.
func (sh *streamingHandler) ProcessStream(
	ctx context.Context,
	message protocol.Message,
	options taskmanager.ProcessOptions,
	handle taskmanager.TaskHandler,
) (*taskmanager.MessageProcessingResult, error) {
	// Extract text from message
	text := extractTextFromMessage(message)
	if text == "" {
		return createErrorResult("input message must contain text"), nil
	}

	// For non-streaming processing, return direct response
	if !options.Streaming {
		response := sh.processTextDirectly(text)
		return createTextResult(response), nil
	}

	// For streaming processing, create a task
	taskID, err := handle.BuildTask(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build task: %w", err)
	}

	// Subscribe to the task for streaming events
	subscriber, err := handle.SubscribeTask(&taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to task: %w", err)
	}

	// Start processing in a goroutine
	go sh.processStreamingTask(ctx, text, taskID, subscriber, handle)

	// Return the task
	task := &protocol.Task{
		ID: taskID,
		Status: protocol.TaskStatus{
			State:     protocol.TaskStateWorking,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}

	return &taskmanager.MessageProcessingResult{
		Result: task,
	}, nil
}

// SupportsStreaming returns true if streaming is supported.
func (sh *streamingHandler) SupportsStreaming() bool {
	return true
}

// processStreamingTask processes text in chunks with streaming updates.
func (sh *streamingHandler) processStreamingTask(
	ctx context.Context,
	text string,
	taskID string,
	subscriber taskmanager.TaskSubscriber,
	handle taskmanager.TaskHandler,
) {
	defer func() {
		if subscriber != nil {
			subscriber.Close()
		}
		handle.CleanTask(&taskID)
	}()

	// Send start message
	startMsg := createTextMessage("Processing your message...")
	err := subscriber.Send(protocol.StreamingMessageEvent{
		Result: startMsg,
	})
	if err != nil {
		return
	}

	// Split text into chunks
	chunks := sh.splitTextIntoChunks(text)
	totalChunks := len(chunks)

	var processedParts []string
	for i, chunk := range chunks {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Process chunk (reverse it as example)
		processedChunk := sh.reverseString(chunk)
		processedParts = append(processedParts, processedChunk)

		// Send progress update
		progressMsg := createTextMessage(fmt.Sprintf("Processed chunk %d/%d: %s", i+1, totalChunks, processedChunk))
		err = subscriber.Send(protocol.StreamingMessageEvent{
			Result: progressMsg,
		})
		if err != nil {
			return
		}

		// Simulate processing delay
		time.Sleep(sh.processDelay)
	}

	// Send final result
	finalResult := fmt.Sprintf("Streaming complete! Processed %d chunks: %s", totalChunks, strings.Join(processedParts, " "))
	completionMsg := createTextMessage(finalResult)
	
	err = handle.UpdateTaskState(&taskID, protocol.TaskStateCompleted, completionMsg)
	if err != nil {
		return
	}
}

// processTextDirectly processes text without streaming.
func (sh *streamingHandler) processTextDirectly(text string) string {
	chunks := sh.splitTextIntoChunks(text)
	var processedParts []string
	
	for _, chunk := range chunks {
		processedChunk := sh.reverseString(chunk)
		processedParts = append(processedParts, processedChunk)
	}
	
	return fmt.Sprintf("Processed %d chunks: %s", len(chunks), strings.Join(processedParts, " "))
}

// splitTextIntoChunks splits text into chunks of specified size.
func (sh *streamingHandler) splitTextIntoChunks(text string) []string {
	if len(text) <= sh.chunkSize {
		return []string{text}
	}

	var chunks []string
	for i := 0; i < len(text); i += sh.chunkSize {
		end := i + sh.chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[i:end])
	}
	return chunks
}

// reverseString reverses a string.
func (sh *streamingHandler) reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Helper functions for message creation
func extractTextFromMessage(message protocol.Message) string {
	for _, part := range message.Parts {
		if textPart, ok := part.(*protocol.TextPart); ok {
			return textPart.Text
		}
	}
	return ""
}

func createTextMessage(text string) *protocol.Message {
	message := protocol.NewMessage(
		protocol.MessageRoleAgent,
		[]protocol.Part{protocol.NewTextPart(text)},
	)
	return &message
}

func createTextResult(text string) *taskmanager.MessageProcessingResult {
	return &taskmanager.MessageProcessingResult{
		Result: createTextMessage(text),
	}
}

func createErrorResult(errorMsg string) *taskmanager.MessageProcessingResult {
	message := protocol.NewMessage(
		protocol.MessageRoleAgent,
		[]protocol.Part{protocol.NewTextPart(fmt.Sprintf("Error: %s", errorMsg))},
	)
	return &taskmanager.MessageProcessingResult{
		Result: &message,
	}
}