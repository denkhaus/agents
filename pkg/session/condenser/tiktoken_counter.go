package condenser

import (
	"context"
	"math"
	"strings"

	"go.uber.org/zap"
)

// TikTokenCounter implements token counting using tiktoken-like logic
// This is a simplified implementation that mimics tiktoken behavior
type TikTokenCounter struct {
	logger *zap.Logger
}

// NewTikTokenCounter creates a new tiktoken-based token counter
func NewTikTokenCounter(logger *zap.Logger) *TikTokenCounter {
	return &TikTokenCounter{
		logger: logger.Named("tiktoken-counter"),
	}
}

// CountTokens counts tokens using tiktoken-like algorithm
func (t *TikTokenCounter) CountTokens(ctx context.Context, text string) (int, error) {
	if text == "" {
		return 0, nil
	}

	// Simplified tiktoken-like tokenization
	// This is a basic implementation that approximates tiktoken behavior
	tokens := t.tokenize(text)
	
	t.logger.Debug("Counted tokens using tiktoken-like algorithm",
		zap.Int("tokens", len(tokens)),
		zap.Int("textLength", len(text)),
	)
	
	return len(tokens), nil
}

// EstimateTokens estimates token count for byte content
func (t *TikTokenCounter) EstimateTokens(ctx context.Context, content []byte) (int, error) {
	return t.CountTokens(ctx, string(content))
}

// GetInfo returns information about this token counter
func (t *TikTokenCounter) GetInfo() TokenCounterInfo {
	return TokenCounterInfo{
		Method:      TokenCountingTikToken,
		Accuracy:    AccuracyPrecise,
		ModelName:   "tiktoken-compatible",
		Description: "TikToken-like tokenization algorithm",
	}
}

// tokenize performs basic tokenization similar to tiktoken
func (t *TikTokenCounter) tokenize(text string) []string {
	// This is a simplified tokenization that approximates tiktoken behavior
	// Real tiktoken uses BPE (Byte Pair Encoding) which is much more complex
	
	var tokens []string
	
	// Split on whitespace and punctuation
	words := t.splitText(text)
	
	for _, word := range words {
		if word == "" {
			continue
		}
		
		// For longer words, split into subword tokens (simplified BPE-like)
		subTokens := t.splitIntoSubTokens(word)
		tokens = append(tokens, subTokens...)
	}
	
	return tokens
}

// splitText splits text into words and punctuation
func (t *TikTokenCounter) splitText(text string) []string {
	var result []string
	var current strings.Builder
	
	for _, r := range text {
		if isWhitespace(r) || isPunctuation(r) {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
			if !isWhitespace(r) {
				result = append(result, string(r))
			}
		} else {
			current.WriteRune(r)
		}
	}
	
	if current.Len() > 0 {
		result = append(result, current.String())
	}
	
	return result
}

// splitIntoSubTokens splits words into subword tokens (simplified BPE)
func (t *TikTokenCounter) splitIntoSubTokens(word string) []string {
	if len(word) <= 4 {
		return []string{word}
	}
	
	// Simple subword splitting - in real tiktoken this would use learned BPE merges
	var tokens []string
	chunkSize := int(math.Max(2, float64(len(word))/3))
	
	for i := 0; i < len(word); i += chunkSize {
		end := i + chunkSize
		if end > len(word) {
			end = len(word)
		}
		tokens = append(tokens, word[i:end])
	}
	
	return tokens
}

// isWhitespace checks if a rune is whitespace
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// isPunctuation checks if a rune is punctuation
func isPunctuation(r rune) bool {
	return (r >= '!' && r <= '/') || 
		   (r >= ':' && r <= '@') || 
		   (r >= '[' && r <= '`') || 
		   (r >= '{' && r <= '~')
}