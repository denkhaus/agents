package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/event"
)

type Invocation struct {
	ID         uuid.UUID
	EventsByID map[uuid.UUID][]*messaging.LLMEvent `json:"events"`
}

// addEvent adds an LLMEvent to the Invocation's map of events,
// keyed by the event's ID.
func (p *Invocation) addEvent(event *messaging.LLMEvent) {
	if p.EventsByID == nil {
		p.EventsByID = make(map[uuid.UUID][]*messaging.LLMEvent)
	}

	if p.EventsByID[event.ID] == nil {
		p.EventsByID[event.ID] = []*messaging.LLMEvent{event}
	} else {
		p.EventsByID[event.ID] = append(p.EventsByID[event.ID], event)
	}
}

// EventCount returns the total number of LLM events associated with this Invocation.
func (p *Invocation) EventCount() int {
	total := 0
	for _, evts := range p.EventsByID {
		total += len(evts)
	}

	return total
}

type EventCollector struct {
	Invocations      map[uuid.UUID]*Invocation
	creationCtx      context.Context
	signalingCtx     context.Context
	cancel           context.CancelFunc
	signalingTimeout time.Duration
}

// NewEventCollector creates and returns a new EventCollector instance.
// It initializes the collector with a given context and signaling timeout.
func NewEventCollector(ctx context.Context, signalingTimeout time.Duration) *EventCollector {
	coll := &EventCollector{
		signalingTimeout: signalingTimeout,
		creationCtx:      ctx,
	}
	return coll.Touch()
}

// Touch resets the signaling context's timeout, effectively extending the collector's lifetime.
func (p *EventCollector) Touch() *EventCollector {
	p.signalingCtx, p.cancel = context.WithTimeout(p.creationCtx, p.signalingTimeout)
	return p
}

// Close cancels the signaling context, releasing resources and signaling completion.
func (p *EventCollector) Close() {
	p.cancel()
}

// Done returns a channel that is closed when the collector's signaling context is cancelled.
func (p *EventCollector) Done() <-chan struct{} {
	return p.signalingCtx.Done()
}

// addEvent adds an LLMEvent to the appropriate Invocation within the collector.
// If the Invocation does not exist, a new one is created.
func (p *EventCollector) addEvent(event *messaging.LLMEvent) {
	if p.Invocations == nil {
		p.Invocations = make(map[uuid.UUID]*Invocation)
	}

	if inv, ok := p.Invocations[event.InvocationID]; ok {
		inv.addEvent(event)
	} else {
		inv := &Invocation{
			ID: event.InvocationID,
		}
		p.Invocations[event.InvocationID] = inv
		inv.addEvent(event)
	}
}

// EventCount returns the total number of LLM events across all Invocations in the collector.
func (p *EventCollector) EventCount() int {
	total := 0
	for _, inv := range p.Invocations {
		total += inv.EventCount()
	}

	return total
}

// Collect processes an incoming event, converts it to an LLMEvent,
// and adds it to the event collector. It also logs the event details.
func (p *EventCollector) Collect(info *messaging.RoutingInfo, ev *event.Event) {
	p.Touch()

	llmEvent, err := messaging.NewLLMEvent(info, ev)
	if err != nil {
		logger.Log.Error("failed to create llm event", zap.Error(err))
		return
	}

	if llmEvent != nil {
		p.addEvent(llmEvent)
		logger.Log.Debug("Added event to collector",
			zap.String("event_id", llmEvent.ID.String()),
			zap.String("invocation_id", llmEvent.InvocationID.String()),
			zap.Int("total_events", p.EventCount()),
		)
		spew.Dump(llmEvent)
	}
}

// Persist serializes the collected LLM events to a JSON file.
// The filename includes a timestamp for uniqueness.
func (p *EventCollector) Persist() error {

	// Serialize output struct to JSON and save to file
	jsonData, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output to JSON: %w", err)
	}

	// Create filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("llm_events_%s.json", timestamp)

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	logger.Log.Info("LLM events saved to file",
		zap.String("filename", filename),
		zap.Int("event_count", p.EventCount()),
	)

	return nil
}
