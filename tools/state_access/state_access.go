package stateaccess

import (
	"context"
	"encoding/json"
	"fmt"

	"trpc.group/trpc-go/trpc-agent-go/session"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// StateAccessTool provides access to session state data.
type StateAccessTool struct {
	sessionService session.Service
	appName        string
	userID         string
	sessionID      string
}

// Declaration returns tool metadata.
func (t *StateAccessTool) Declaration() *tool.Declaration {
	return &tool.Declaration{
		Name:        "get_session_state",
		Description: "Retrieve data from the current session state. Use this to access information stored by previous agents in the chain.",
		InputSchema: &tool.Schema{
			Type: "object",
			Properties: map[string]*tool.Schema{
				"key": {
					Type:        "string",
					Description: "The key of the data to retrieve from session state.",
				},
			},
			Required: []string{"key"},
		},
	}
}

// Call executes the tool to retrieve data from session state.
func (t *StateAccessTool) Call(ctx context.Context, jsonArgs []byte) (any, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(jsonArgs, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	key, ok := params["key"].(string)
	if !ok {
		return nil, fmt.Errorf("key parameter must be a string")
	}

	// Create session key.
	sessionKey := session.Key{
		AppName:   t.appName,
		UserID:    t.userID,
		SessionID: t.sessionID,
	}

	// Get session state.
	sessionData, err := t.sessionService.GetSession(ctx, sessionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Extract data from session state.
	if sessionData.State == nil {
		return map[string]interface{}{
			"result": "No data found in session state",
		}, nil
	}

	// Look for the specific key in the state.
	if data, exists := sessionData.State[key]; exists {
		return map[string]interface{}{
			"result": fmt.Sprintf("Found data for key '%s': %s", key, string(data)),
		}, nil
	}

	// If key not found, return available keys.
	availableKeys := make([]string, 0, len(sessionData.State))
	for k := range sessionData.State {
		availableKeys = append(availableKeys, k)
	}

	return map[string]interface{}{
		"result": fmt.Sprintf("Key '%s' not found. Available keys: %v", key, availableKeys),
	}, nil
}

// NewStateAccessTool creates a new StateAccessTool with the provided options.
func NewStateAccessTool(opts ...StateAccessToolOption) (tool.CallableTool, error) {
	t := &StateAccessTool{}

	// Apply all options
	for _, opt := range opts {
		opt(t)
	}

	// Validate required fields
	if t.sessionService == nil {
		return nil, fmt.Errorf("session service is required")
	}
	if t.appName == "" {
		return nil, fmt.Errorf("app name is required")
	}
	if t.userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	if t.sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	return t, nil
}
