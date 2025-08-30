package fetch

import (
	"context"
	"testing"
)

func TestFetchArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    fetchArgs
		wantErr bool
	}{
		{
			name:    "empty URL should error",
			args:    fetchArgs{URL: "", Method: "GET"},
			wantErr: true,
		},
		{
			name:    "valid URL with GET method",
			args:    fetchArgs{URL: "https://example.com", Method: "GET"},
			wantErr: false,
		},
		{
			name:    "valid URL with POST method",
			args:    fetchArgs{URL: "https://example.com", Method: "POST"},
			wantErr: false,
		},
		{
			name:    "unsupported method should error",
			args:    fetchArgs{URL: "https://example.com", Method: "PUT"},
			wantErr: true,
		},
		{
			name:    "invalid method should error",
			args:    fetchArgs{URL: "https://example.com", Method: "INVALID"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fetch(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFetchDefaultValues(t *testing.T) {
	args := fetchArgs{
		URL: "https://example.com",
		// Method should default to GET
	}

	_, err := fetch(context.Background(), args)
	if err != nil {
		t.Errorf("fetch() with default method failed: %v", err)
	}
}

func TestFetchTimeoutConfiguration(t *testing.T) {
	args := fetchArgs{
		URL:     "https://example.com",
		Method:  "GET",
		Timeout: 5, // 5 seconds
	}

	_, err := fetch(context.Background(), args)
	// We can't easily test timeout behavior without mocking, but we can test it doesn't crash
	if err != nil {
		// This is expected to fail due to network, but shouldn't be a validation error
		if err.Error() == "URL is required" || 
		   err.Error() == "unsupported HTTP method" {
			t.Errorf("unexpected validation error: %v", err)
		}
	}
}