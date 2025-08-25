package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/server"
)

// agentCardFetcher implements AgentCardFetcher.
type agentCardFetcher struct {
	timeout time.Duration
}

// newAgentCardFetcher creates a new agent card fetcher.
func newAgentCardFetcher(timeout time.Duration) *agentCardFetcher {
	return &agentCardFetcher{
		timeout: timeout,
	}
}

// FetchAgentCard implements AgentCardFetcher.
func (f *agentCardFetcher) FetchAgentCard(ctx context.Context, baseURL string) (*server.AgentCard, error) {
	// Ensure base URL ends with "/"
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	// Construct agent card URL
	cardURL := baseURL + protocol.AgentCardPath[1:] // Remove leading slash

	// Create request with timeout
	reqCtx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, cardURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agent card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response
	var card server.AgentCard
	if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
		return nil, fmt.Errorf("failed to decode agent card: %w", err)
	}

	return &card, nil
}
