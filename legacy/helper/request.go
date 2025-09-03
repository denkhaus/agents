// Package helper provides utility functions for client operations.
package helper

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/samber/mo"
)

// RequestConfig defines HTTP request configuration.
type RequestConfig struct {
	Method     string
	URL        string
	Headers    mo.Option[map[string]string]
	Body       mo.Option[[]byte]
	Timeout    mo.Option[time.Duration]
	RetryCount mo.Option[int]
}

// HTTPClient defines an interface for HTTP operations.
type HTTPClient interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// RequestHandler handles HTTP requests with retry logic.
type RequestHandler struct {
	client     *http.Client
	retryCount int
}

// httpClient implements HTTPClient interface.
type httpClient struct {
	client     *http.Client
	retryCount int
}

// NewRequestHandler creates a new request handler.
func NewRequestHandler(timeout time.Duration, retryCount int) *RequestHandler {
	return &RequestHandler{
		client: &http.Client{
			Timeout: timeout,
		},
		retryCount: retryCount,
	}
}

// NewHTTPClient creates a new HTTP client.
func NewHTTPClient(timeout time.Duration, retryCount int) HTTPClient {
	return &httpClient{
		client: &http.Client{
			Timeout: timeout,
		},
		retryCount: retryCount,
	}
}

// Do implements HTTPClient interface.
func (c *httpClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error
	
	for attempt := 0; attempt <= c.retryCount; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}
		
		// Set context on request
		reqWithCtx := req.WithContext(ctx)
		
		resp, err := c.client.Do(reqWithCtx)
		if err == nil {
			return resp, nil
		}
		
		lastErr = err
		
		// Don't retry on context cancellation
		if ctx.Err() != nil {
			break
		}
	}
	
	return nil, fmt.Errorf("request failed after %d attempts: %w", c.retryCount+1, lastErr)
}

// Execute executes an HTTP request with retry logic.
func (h *RequestHandler) Execute(ctx context.Context, config RequestConfig) (*http.Response, error) {
	timeout := config.Timeout.OrElse(30 * time.Second)
	retryCount := config.RetryCount.OrElse(h.retryCount)

	var lastErr error

	for attempt := 0; attempt <= retryCount; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		resp, err := h.executeOnce(ctx, config, timeout)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Don't retry on context cancellation
		if ctx.Err() != nil {
			break
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", retryCount+1, lastErr)
}

// executeOnce executes a single HTTP request.
func (h *RequestHandler) executeOnce(ctx context.Context, config RequestConfig, timeout time.Duration) (*http.Response, error) {
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var req *http.Request
	var err error

	if body, hasBody := config.Body.Get(); hasBody {
		req, err = http.NewRequestWithContext(reqCtx, config.Method, config.URL,
			NewByteReader(body))
	} else {
		req, err = http.NewRequestWithContext(reqCtx, config.Method, config.URL, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	if headers, hasHeaders := config.Headers.Get(); hasHeaders {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	return h.client.Do(req)
}
