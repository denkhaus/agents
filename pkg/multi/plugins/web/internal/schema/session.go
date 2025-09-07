package schema

import (
	"fmt"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

// NewADKSession converts an internal session object to the ADKSession
func NewADKSession(s *session.Session) (*ADKSession, error) {
	adkEvents := make([]*LLMEvent, 0, len(s.Events))
	for _, e := range s.Events {
		ev, err := NewLLMEvent(&e, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create LLMEvent: %w", err)
		}
		if ev != nil {
			adkEvents = append(adkEvents, ev)
		}
	}

	agentID, err := uuid.Parse(s.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse agentID: %w", err)
	}

	sessionID, err := uuid.Parse(s.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sessionID: %w", err)
	}

	return &ADKSession{
		AppName:        s.AppName,
		AgentID:        agentID,
		ID:             sessionID,
		CreateTime:     s.CreatedAt.Unix(),
		LastUpdateTime: s.UpdatedAt.Unix(),
		State:          map[string][]byte(s.State),
		Events:         adkEvents,
	}, nil
}
