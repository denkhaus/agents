# Session Condenser

The Session Condenser is a smart wrapper around session services that automatically condenses session state when it grows too large. It uses token-aware calculations and LLM-powered summarization to maintain conversation context while keeping memory usage under control.

## Overview

The condenser monitors session token usage and automatically triggers condensation when a configurable threshold is reached. It replaces old conversation history with an intelligent summary while preserving recent events, ensuring conversations can continue seamlessly without losing important context.

## Key Features

- **Token-Aware**: Uses actual token counts instead of byte estimates
- **Pluggable Token Counting**: Automatic selection of best available counting method
- **LLM Independence**: Separate LLM for summarization (recommended)
- **Smart Preservation**: Keeps recent events while condensing older history
- **Comprehensive Metrics**: Tracks effectiveness and performance
- **Robust Error Handling**: Never breaks ongoing conversations
- **Configurable**: Extensive options for fine-tuning behavior

## Quick Start

```go
import (
    "github.com/denkhaus/agents/session/condenser"
    "trpc.group/trpc-go/trpc-agent-go/session/inmemory"
    "trpc.group/trpc-go/trpc-agent-go/model/openai"
)

// Create base session service
baseSessionService := inmemory.NewSessionService()

// Create LLM for summarization
summarizerLLM := openai.New("gpt-3.5-turbo") // Cheaper model for summaries

// Create condenser with sensible defaults
condenserService, err := condenser.NewWithOptions(
    baseSessionService,
    summarizerLLM,
    logger,
    condenser.WithMaxContextTokens(4000),
    condenser.WithTriggerThreshold(0.75),
    condenser.WithTokenCountingMethod(condenser.TokenCountingAuto),
)
```

## LLM Selection Strategy

### Why Use Different LLMs?

**It's recommended (but not required) to use a different LLM for condensation than your main chat LLM.** Here's why:

#### 1. **Cost Optimization**
```go
// Main chat: Use powerful, expensive model
chatLLM := openai.New("gpt-4o")

// Condensation: Use cheaper, efficient model
summarizerLLM := openai.New("gpt-3.5-turbo") // 10x cheaper

condenser.NewWithOptions(baseService, summarizerLLM, logger, ...)
```

#### 2. **Performance Isolation**
- Condensation runs in background without blocking main conversation
- Separate rate limits and quotas
- Independent failure modes

#### 3. **Specialization**
```go
// Use a model specifically good at summarization
summarizerLLM := openai.New("gpt-3.5-turbo")
// Or even a specialized summarization model
// summarizerLLM := anthropic.New("claude-3-haiku")
```

#### 4. **Resource Management**
- Separate workload from main conversation processing
- Different scaling requirements
- Independent monitoring and alerting

### When to Use the Same LLM

You **can** use the same LLM for both chat and condensation:

```go
// Same LLM for both (like in our example)
mainLLM := openai.New("gpt-4o")

// Chat agent
chatAgent := llmagent.New("assistant", llmagent.WithModel(mainLLM))

// Condenser (same LLM)
condenserService, err := condenser.NewWithOptions(
    baseService,
    mainLLM, // Same LLM instance
    logger,
    // ... options
)
```

**Use the same LLM when:**
- Simplicity is more important than cost optimization
- You have abundant LLM quota/budget
- You want consistent behavior across all operations
- You're prototyping or in development

## Token Counting Methods

The condenser automatically selects the best available token counting method:

### 1. **Auto Selection (Recommended)**
```go
condenser.WithTokenCountingMethod(condenser.TokenCountingAuto)
```
- Tries LLM-native counting first
- Falls back to heuristic if unavailable
- Best balance of accuracy and reliability

### 2. **Heuristic Counting**
```go
condenser.WithTokenCountingMethod(condenser.TokenCountingHeuristic),
condenser.WithCharsPerToken(4.0), // Customize ratio
```
- Character-based estimation
- Always available, fast
- ~4 characters per token for most models

### 3. **LLM-Native Counting**
```go
condenser.WithTokenCountingMethod(condenser.TokenCountingLLMNative)
```
- Uses model's actual tokenizer
- Most accurate when available
- May have performance overhead

### 4. **Custom Implementation**
```go
type MyTokenCounter struct{}
func (c *MyTokenCounter) CountTokens(ctx context.Context, text string) (int, error) {
    // Your custom logic
}

condenser.WithCustomTokenCounter(&MyTokenCounter{})
```

## Configuration Options

### Basic Configuration
```go
condenser.NewWithOptions(
    sessionService,
    summarizerLLM,
    logger,
    
    // Token limits
    condenser.WithMaxContextTokens(8000),        // Max tokens before condensation
    condenser.WithTriggerThreshold(0.75),        // Trigger at 75% capacity
    
    // Behavior
    condenser.WithRecentEventsToKeep(3),         // Keep 3 most recent events
    condenser.WithSummaryPrompt("Custom prompt..."), // Custom summarization prompt
    
    // Token counting
    condenser.WithTokenCountingMethod(condenser.TokenCountingAuto),
    condenser.WithCharsPerToken(4.0),            // For heuristic method
    
    // Performance
    condenser.WithTokenCaching(true),            // Enable token count caching
    condenser.WithCacheSize(1000),               // Cache size limit
)
```

### Default Values
```go
config := condenser.DefaultConfig()
// MaxContextTokens:    8000
// TriggerThreshold:    0.75 (75%)
// RecentEventsToKeep:  3
// TokenCountingMethod: TokenCountingAuto
// CharsPerToken:       4.0
// EnableTokenCaching:  true
// CacheSize:          1000
```

## Monitoring and Metrics

### Getting Metrics
```go
metrics := condenserService.GetMetrics()
fmt.Printf("Condensations: %d\n", metrics.CondensationCount)
fmt.Printf("Tokens Saved: %d\n", metrics.TotalTokensSaved)
fmt.Printf("Avg Reduction: %.2f%%\n", metrics.AverageReductionRatio*100)
fmt.Printf("Failures: %d\n", metrics.FailureCount)
```

### Token Counter Information
```go
info := condenserService.GetTokenCounterInfo()
fmt.Printf("Method: %s\n", info.Method.String())
fmt.Printf("Accuracy: %s\n", info.Accuracy.String())
fmt.Printf("Description: %s\n", info.Description)
```

## How Condensation Works

1. **Monitoring**: Every event append triggers token calculation
2. **Threshold Check**: Compares current tokens vs. configured limit
3. **Summary Generation**: LLM creates concise conversation summary
4. **Event Preservation**: Keeps N most recent events for context
5. **Session Replacement**: Atomically replaces old session with condensed version
6. **Metrics Update**: Tracks effectiveness and performance

### Example Condensation Flow
```
Before Condensation (5000 tokens):
├── Event 1: "Hello" (100 tokens)
├── Event 2: "How are you?" (200 tokens)
├── Event 3: "Tell me about AI" (300 tokens)
├── Event 4: "That's interesting..." (400 tokens)
└── Event 5: "What about ML?" (4000 tokens)

After Condensation (2500 tokens):
├── Summary: "User asked about AI and ML..." (500 tokens)
├── Event 4: "That's interesting..." (400 tokens)  # Recent events kept
└── Event 5: "What about ML?" (1600 tokens)        # Recent events kept

Token Savings: 2500 tokens (50% reduction)
```

## Error Handling

The condenser is designed to never break ongoing conversations:

- **Condensation failures**: Logged but don't stop the session
- **Token counting errors**: Fall back to heuristic methods
- **LLM errors**: Session continues without condensation
- **Session creation failures**: Detailed logging for debugging

## Best Practices

### 1. **Choose Appropriate Thresholds**
```go
// Conservative: Condense early, smaller summaries
condenser.WithTriggerThreshold(0.5)  // 50%

// Aggressive: Condense late, larger summaries  
condenser.WithTriggerThreshold(0.9)  // 90%

// Balanced: Good default
condenser.WithTriggerThreshold(0.75) // 75%
```

### 2. **Optimize Recent Events**
```go
// Keep more events for better context continuity
condenser.WithRecentEventsToKeep(5)

// Keep fewer events for maximum compression
condenser.WithRecentEventsToKeep(1)

// Balanced approach
condenser.WithRecentEventsToKeep(3)
```

### 3. **Custom Summarization Prompts**
```go
condenser.WithSummaryPrompt(`
Summarize this conversation focusing on:
1. Key decisions made
2. Important facts shared
3. Current context needed for continuation
Keep it concise but preserve essential information.
`)
```

### 4. **Monitor Effectiveness**
```go
// Regularly check metrics
metrics := service.GetMetrics()
if metrics.AverageReductionRatio < 0.3 {
    // Condensation not effective, adjust settings
    log.Warn("Low condensation effectiveness")
}
```

## Testing

Run the test suite:
```bash
cd session/condenser
go test -v
```

The tests cover:
- Service creation and configuration
- Token counting accuracy
- Error handling
- Type safety

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Your App      │───▶│ Condenser Service │───▶│ Base Session    │
│                 │    │                  │    │ Service         │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌──────────────────┐
                       │ Summarizer LLM   │
                       │ (Separate/Same)   │
                       └──────────────────┘
                                │
                                ▼
                       ┌──────────────────┐
                       │ Token Counter    │
                       │ (Auto-selected)  │
                       └──────────────────┘
```

## Examples

See the complete example in `examples/agents/tokentracker/main.go` which demonstrates:
- Setting up the condenser service
- Token usage tracking
- Interactive conversation with automatic condensation
- Metrics monitoring

## Contributing

When adding new token counting methods:
1. Implement the `TokenCounter` interface
2. Add the new method to `TokenCountingMethod` enum
3. Update the factory in `factory.go`
4. Add comprehensive tests
5. Update this documentation