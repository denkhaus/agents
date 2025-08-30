package fetch

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const (
	ToolName = "fetch"
)

// fetchArgs holds the input for the fetch tool.
type fetchArgs struct {
	URL        string            `json:"url" description:"The URL to fetch"`
	Method     string            `json:"method" description:"HTTP method: GET or POST"`
	Headers    map[string]string `json:"headers,omitempty" description:"Optional HTTP headers"`
	Body       string            `json:"body,omitempty" description:"Optional request body for POST requests"`
	Timeout    int               `json:"timeout,omitempty" description:"Timeout in seconds (default: 30)"`
	Username   string            `json:"username,omitempty" description:"Username for Basic Authentication"`
	Password   string            `json:"password,omitempty" description:"Password for Basic Authentication"`
	ReturnType string            `json:"return_type,omitempty" description:"Return type: string or buffer (default: string)"`
}

// fetchResult holds the output for the fetch tool.
type fetchResult struct {
	StatusCode  int               `json:"status_code"`
	Status      string            `json:"status"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	ContentType string            `json:"content_type"`
	MIMEType    string            `json:"mime_type"`
	Size        int               `json:"size"`
	Duration    float64           `json:"duration_seconds"`
}

// fetch performs HTTP requests with configurable options.
func fetch(ctx context.Context, args fetchArgs) (fetchResult, error) {
	start := time.Now()

	// Validate URL
	if args.URL == "" {
		return fetchResult{}, fmt.Errorf("URL is required")
	}

	// Validate HTTP method
	method := strings.ToUpper(args.Method)
	if method == "" {
		method = "GET"
	}
	if method != "GET" && method != "POST" {
		return fetchResult{}, fmt.Errorf("unsupported HTTP method: %s. Only GET and POST are supported", method)
	}

	// Set default timeout
	timeout := 30 * time.Second
	if args.Timeout > 0 {
		timeout = time.Duration(args.Timeout) * time.Second
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: timeout,
	}

	// Create request
	var bodyReader io.Reader
	if method == "POST" && args.Body != "" {
		bodyReader = strings.NewReader(args.Body)
	}

	req, err := http.NewRequestWithContext(ctx, method, args.URL, bodyReader)
	if err != nil {
		return fetchResult{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if args.Headers != nil {
		for key, value := range args.Headers {
			req.Header.Set(key, value)
		}
	}

	// Set Basic Authentication if provided
	if args.Username != "" && args.Password != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(args.Username + ":" + args.Password))
		req.Header.Set("Authorization", "Basic "+auth)
	}

	// Set Content-Type for POST requests if not provided
	if method == "POST" && args.Body != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return fetchResult{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fetchResult{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Convert headers to map
	headers := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// Extract content type and MIME type
	contentType := resp.Header.Get("Content-Type")
	mimeType := strings.Split(contentType, ";")[0] // Remove charset etc.

	// Handle return type
	var bodyContent string
	if args.ReturnType == "buffer" {
		// For binary data, return as base64 encoded string
		bodyContent = base64.StdEncoding.EncodeToString(bodyBytes)
	} else {
		// Default: return as string
		bodyContent = string(bodyBytes)
	}

	duration := time.Since(start).Seconds()

	return fetchResult{
		StatusCode:  resp.StatusCode,
		Status:      resp.Status,
		Headers:     headers,
		Body:        bodyContent,
		ContentType: contentType,
		MIMEType:    mimeType,
		Size:        len(bodyBytes),
		Duration:    duration,
	}, nil
}

func NewTool() (tool.Tool, error) {
	// Create fetch tool for HTTP requests
	fetchTool := function.NewFunctionTool(
		fetch,
		function.WithName(ToolName),
		function.WithDescription(
			"Perform HTTP requests to fetch content from URLs. "+
				"Supports GET and POST methods with configurable timeouts, "+
				"headers, basic authentication, and return type options.",
		),
	)

	return fetchTool, nil
}
