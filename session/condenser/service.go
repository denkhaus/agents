package condenser

import (
	"fmt"

	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

// New creates a new condenser service
func New(
	sessionService session.Service,
	summarizerLLM model.Model,
	config Config,
	logger *zap.Logger,
) (*Service, error) {
	if sessionService == nil {
		return nil, fmt.Errorf("session service cannot be nil")
	}
	if summarizerLLM == nil {
		return nil, fmt.Errorf("summarizer LLM cannot be nil")
	}

	// Create token counter based on config
	tokenCounter, err := createTokenCounter(summarizerLLM, config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create token counter: %w", err)
	}

	return &Service{
		sessionService: sessionService,
		summarizerLLM:  summarizerLLM,
		tokenCounter:   tokenCounter,
		config:         config,
		logger:         logger.Named("condenser"),
		metrics:        &Metrics{},
	}, nil
}

// NewWithOptions creates a new condenser service with functional options
func NewWithOptions(
	sessionService session.Service,
	summarizerLLM model.Model,
	logger *zap.Logger,
	opts ...ConfigOption,
) (*Service, error) {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(&config)
	}

	return New(sessionService, summarizerLLM, config, logger)
}

// GetMetrics returns current condensation metrics
func (s *Service) GetMetrics() Metrics {
	return *s.metrics
}

// GetConfig returns the current configuration
func (s *Service) GetConfig() Config {
	return s.config
}

// GetTokenCounterInfo returns information about the token counter
func (s *Service) GetTokenCounterInfo() TokenCounterInfo {
	return s.tokenCounter.GetInfo()
}

// Close closes the condenser service
func (s *Service) Close() error {
	return s.sessionService.Close()
}