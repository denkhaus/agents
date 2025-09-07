package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/denkhaus/agents/pkg/multi/plugins/web/internal/schema"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"trpc.group/trpc-go/trpc-agent-go/log"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

const (
	keyEventID      = "gcp.vertex.agent.event_id"
	keySessionID    = "gcp.vertex.agent.session_id"
	keyInvocationID = "gcp.vertex.agent.invocation_id"
	keyLLMRequest   = "gcp.vertex.agent.llm_request"
	keyLLMResponse  = "gcp.vertex.agent.llm_response"
)

// func setTraceInfo() {
// 	itelemetry.KeyEventID = keyEventID
// 	itelemetry.KeySessionID = keySessionID
// 	itelemetry.KeyLLMRequest = keyLLMRequest
// 	itelemetry.KeyLLMResponse = keyLLMResponse
// 	itelemetry.KeyInvocationID = keyInvocationID
// }

type inMemoryExporter struct {
	sessionTraces map[string]map[string]struct{} // key: session_id, value: map[event_id]struct{}
	spans         []sdktrace.ReadOnlySpan
}

func newInMemoryExporter() *inMemoryExporter {
	return &inMemoryExporter{sessionTraces: make(map[string]map[string]struct{})}
}
func (e *inMemoryExporter) ExportSpans(_ context.Context, spans []sdktrace.ReadOnlySpan) error {
	for _, span := range spans {
		// if span.Name() != itelemetry.SpanNameCallLLM {
		// 	continue
		// }
		for _, attr := range span.Attributes() {
			if attr.Key != keySessionID {
				continue
			}
			sessionID := attr.Value.AsString()
			traceID := span.SpanContext().TraceID().String()
			if _, ok := e.sessionTraces[sessionID]; !ok {
				e.sessionTraces[sessionID] = map[string]struct{}{
					traceID: {},
				}
			} else {
				e.sessionTraces[sessionID][traceID] = struct{}{}
			}
			break
		}
	}
	e.spans = append(e.spans, spans...)
	return nil
}

func (e *inMemoryExporter) Shutdown(_ context.Context) error {
	return nil
}

func (e *inMemoryExporter) getFinishedSpans(sessionID string) []sdktrace.ReadOnlySpan {
	traceIDs := e.sessionTraces[sessionID]
	var spans []sdktrace.ReadOnlySpan
	for traceID := range traceIDs {
		for _, s := range e.spans {
			if s.SpanContext().TraceID().String() == traceID {
				spans = append(spans, s)
			}
		}
	}
	return spans
}

type apiServerSpanExporter struct {
	traces map[string]attribute.Set
}

func newApiServerSpanExporter(ts map[string]attribute.Set) *apiServerSpanExporter {
	return &apiServerSpanExporter{traces: ts}
}

func (e *apiServerSpanExporter) ExportSpans(_ context.Context, spans []sdktrace.ReadOnlySpan) error {
	for _, span := range spans {
		// if name := span.Name(); name != itelemetry.SpanNameCallLLM && !strings.HasPrefix(name, itelemetry.SpanNamePrefixExecuteTool) {
		// 	continue
		// }
		baseAttrs := []attribute.KeyValue{
			attribute.String("trace_id", span.SpanContext().TraceID().String()),
			attribute.String("span_id", span.SpanContext().SpanID().String()),
		}
		allAttrs := append(baseAttrs, span.Attributes()...)
		attributes := attribute.NewSet(allAttrs...)

		if eventID, ok := attributes.Value(keyEventID); ok {
			e.traces[eventID.AsString()] = attributes
		}
	}
	return nil
}

func (e *apiServerSpanExporter) Shutdown(_ context.Context) error {
	return nil
}

func (s *Server) handleEventTrace(w http.ResponseWriter, r *http.Request) {
	log.Infof("handleEventTrace called: path=%s", r.URL.Path)
	vars := mux.Vars(r)
	eventID := vars["event_id"]
	trace, ok := s.traces[eventID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Trace not found"))
		return
	}
	s.writeJSON(w, buildTraceAttributes(trace))
}

func (s *Server) handleSessionTrace(w http.ResponseWriter, r *http.Request) {
	log.Infof("handleSessionTrace called: path=%s", r.URL.Path)
	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	var spans []schema.Span
	for _, span := range s.memoryExporter.getFinishedSpans(sessionID) {
		result := buildTraceAttributes(attribute.NewSet(span.Attributes()...))
		spans = append(spans, schema.Span{
			Name:         span.Name(),
			SpanID:       span.SpanContext().SpanID().String(),
			TraceID:      span.SpanContext().TraceID().String(),
			StartTime:    span.StartTime().UnixNano(),
			EndTime:      span.EndTime().UnixNano(),
			Attributes:   result,
			ParentSpanID: span.Parent().SpanID().String(),
		})
	}
	s.writeJSON(w, spans)
}

func buildTraceAttributes(attributes attribute.Set) map[string]any {
	result := make(map[string]any)
	for iter := attributes.Iter(); iter.Next(); {
		attr := iter.Attribute()
		if attr.Key == keyLLMRequest {
			var req model.Request
			if err := json.Unmarshal([]byte(attr.Value.AsString()), &req); err == nil {
				var contents []schema.Content
				for _, c := range req.Messages {
					contents = append(contents, schema.Content{
						Role: c.Role.String(),
						Parts: []schema.PartIncoming{
							{
								Text: c.Content,
							},
						},
					})
				}
				result[string(attr.Key)] = schema.TraceLLMRequest{
					Contents: contents,
				}
			}
		} else {
			result[string(attr.Key)] = attr.Value.AsString()
		}
	}
	return result
}
