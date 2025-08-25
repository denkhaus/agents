package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ResponseParser handles HTTP response parsing.
type ResponseParser struct{}

// NewResponseParser creates a new response parser.
func NewResponseParser() *ResponseParser {
	return &ResponseParser{}
}

// ParseJSON parses JSON response into the provided interface.
func (p *ResponseParser) ParseJSON(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return nil
}

// ReadBody reads the response body as bytes.
func (p *ResponseParser) ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}
	}

	return io.ReadAll(resp.Body)
}
