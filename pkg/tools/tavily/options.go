package tavily

import "github.com/iamwavecut/go-tavily"

// Option is a configuration function for the Tavily tool.
type Option func(*TavilyToolSet)

// WithAPIKey sets the Tavily API key.
func WithAPIKey(apiKey string) Option {
	return func(t *TavilyToolSet) {
		t.ApiKey = apiKey
	}
}

// WithClient sets the Tavily client.
func WithClient(client *tavily.Client) Option {
	return func(t *TavilyToolSet) {
		t.client = client
	}
}

// WithSearchEnabled enables or disables the search tool.
func WithSearchEnabled(enabled bool) Option {
	return func(t *TavilyToolSet) {
		t.SearchEnabled = enabled
	}
}

// WithCrawlEnabled enables or disables the crawl tool.
func WithCrawlEnabled(enabled bool) Option {
	return func(t *TavilyToolSet) {
		t.CrawlEnabled = enabled
	}
}

// WithExtractEnabled enables or disables the extract tool.
func WithExtractEnabled(enabled bool) Option {
	return func(t *TavilyToolSet) {
		t.ExtractEnabled = enabled
	}
}

// WithMapEnabled enables or disables the map tool.
func WithMapEnabled(enabled bool) Option {
	return func(t *TavilyToolSet) {
		t.MapEnabled = enabled
	}
}
