package helper

import (
	"fmt"
)

// HTTPError represents an HTTP error response.
type HTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Status)
}

// IsClientError returns true if the error is a 4xx client error.
func (e *HTTPError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsServerError returns true if the error is a 5xx server error.
func (e *HTTPError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// ClientError represents a client-side error.
type ClientError struct {
	Operation string
	Cause     error
}

// Error implements the error interface.
func (e *ClientError) Error() string {
	return fmt.Sprintf("client error in %s: %v", e.Operation, e.Cause)
}

// Unwrap returns the underlying error.
func (e *ClientError) Unwrap() error {
	return e.Cause
}

// AuthError represents an authentication error.
type AuthError struct {
	Type    string
	Message string
}

// Error implements the error interface.
func (e *AuthError) Error() string {
	return fmt.Sprintf("authentication error (%s): %s", e.Type, e.Message)
}

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field %s: %s", e.Field, e.Message)
}
