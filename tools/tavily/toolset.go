package tavily

import (
	"context"
	"fmt"
	"os"

	"github.com/iamwavecut/go-tavily"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const (
	ToolSetName            = "tavily_toolset"
	DefaultSearchToolName  = "tavily_search"
	DefaultCrawlToolName   = "tavily_crawl"
	DefaultExtractToolName = "tavily_extract"
	DefaultMapToolName     = "tavily_map"
)

// tavilyToolSet implements the ToolSet interface for Tavily API.
type tavilyToolSet struct {
	client         *tavily.Client
	apiKey         string
	searchEnabled  bool
	crawlEnabled   bool
	extractEnabled bool
	mapEnabled     bool
	tools          []tool.CallableTool
}

// NewToolSet creates a new Tavily tool set with the provided options.
func NewToolSet(opts ...Option) (tool.ToolSet, error) {
	t := &tavilyToolSet{
		searchEnabled:  true,
		crawlEnabled:   true,
		extractEnabled: true,
		mapEnabled:     true,
	}

	for _, opt := range opts {
		opt(t)
	}

	if t.apiKey == "" {
		t.apiKey = os.Getenv("TAVILY_API_KEY")
	}

	if t.apiKey == "" {
		return nil, fmt.Errorf("Tavily API key not provided. Set TAVILY_API_KEY environment variable or use WithAPIKey option.")
	}

	if t.client == nil {
		t.client = tavily.New(t.apiKey, nil)
	}

	var tools []tool.CallableTool
	if t.searchEnabled {
		tools = append(tools, t.searchTool())
	}
	if t.crawlEnabled {
		tools = append(tools, t.crawlTool())
	}
	if t.extractEnabled {
		tools = append(tools, t.extractTool())
	}
	if t.mapEnabled {
		tools = append(tools, t.mapTool())
	}
	t.tools = tools

	return t, nil
}

// Tools implements the ToolSet interface.
func (t *tavilyToolSet) Tools(ctx context.Context) []tool.CallableTool {
	return t.tools
}

// Close implements the ToolSet interface.
func (t *tavilyToolSet) Close() error {
	return nil
}

// Search tool
type searchArgs struct {
	Query          string   `json:"query" description:"The search query."`
	SearchDepth    string   `json:"search_depth,omitempty" description:"The depth of the search. Can be 'basic' or 'advanced'."`
	IncludeAnswer  bool     `json:"include_answer,omitempty" description:"Whether to include a direct answer to the query."`
	IncludeDomains []string `json:"include_domains,omitempty" description:"A list of domains to include in the search."`
	ExcludeDomains []string `json:"exclude_domains,omitempty" description:"A list of domains to exclude from the search."`
}

func (t *tavilyToolSet) search(ctx context.Context, args searchArgs) (*tavily.SearchResponse, error) {
	opts := &tavily.SearchOptions{
		SearchDepth:    args.SearchDepth,
		IncludeAnswer:  args.IncludeAnswer,
		IncludeDomains: args.IncludeDomains,
		ExcludeDomains: args.ExcludeDomains,
	}
	return t.client.Search(ctx, args.Query, opts)
}

func (t *tavilyToolSet) searchTool() tool.CallableTool {
	return function.NewFunctionTool(
		t.search,
		function.WithName(DefaultSearchToolName),
		function.WithDescription(`Perform a web search using the Tavily API.
When to use: To find general information on the internet, perform a web search, or get a list of relevant web pages for a query.
When not to use: When you need to extract content from specific URLs (use the `+DefaultExtractToolName+` tool) or crawl a website (use the `+DefaultCrawlToolName+` tool).`),
	)
}

// Crawl tool
type crawlArgs struct {
	URL      string `json:"url" description:"The URL to crawl."`
	MaxDepth int    `json:"max_depth,omitempty" description:"The maximum depth to crawl."`
}

func (t *tavilyToolSet) crawl(ctx context.Context, args crawlArgs) (*tavily.CrawlResponse, error) {
	opts := &tavily.CrawlOptions{
		MaxDepth: args.MaxDepth,
	}
	return t.client.Crawl(ctx, args.URL, opts)
}

func (t *tavilyToolSet) crawlTool() tool.CallableTool {
	return function.NewFunctionTool(
		t.crawl,
		function.WithName(DefaultCrawlToolName),
		function.WithDescription(`Crawl a website using the Tavily API.
When to use: To systematically explore and retrieve content from a website.
When not to use: When you only need to search for general information (use the `+DefaultSearchToolName+` tool) or extract content from known URLs (use the `+DefaultExtractToolName+` tool).`),
	)
}

// Extract tool
type extractArgs struct {
	URLs []string `json:"urls" description:"The URLs to extract content from."`
}

func (t *tavilyToolSet) extract(ctx context.Context, args extractArgs) (*tavily.ExtractResponse, error) {
	return t.client.Extract(ctx, args.URLs, nil)
}

func (t *tavilyToolSet) extractTool() tool.CallableTool {
	return function.NewFunctionTool(
		t.extract,
		function.WithName(DefaultExtractToolName),
		function.WithDescription(`Extract content from a list of URLs using the Tavily API.
When to use: To get the full content of specific web pages when you already have the URLs.
When not to use: When you need to discover new URLs (use the `+DefaultSearchToolName+` or `+DefaultCrawlToolName+` tools) or map a website's structure (use the `+DefaultMapToolName+` tool).`),
	)
}

// Map tool
type mapArgs struct {
	URL      string `json:"url" description:"The URL to map."`
	MaxDepth int    `json:"max_depth,omitempty" description:"The maximum depth to map."`
}

func (t *tavilyToolSet) mapTool() tool.CallableTool {
	return function.NewFunctionTool(
		t.mapFunc,
		function.WithName(DefaultMapToolName),
		function.WithDescription(`Map the structure of a website using the Tavily API.
When to use: To understand the hierarchy and relationships between pages on a website.
When not to use: When you need to search for general information (use the `+DefaultSearchToolName+` tool), crawl content (use the `+DefaultCrawlToolName+` tool), or extract specific page content (use the `+DefaultExtractToolName+` tool).`),
	)
}

func (t *tavilyToolSet) mapFunc(ctx context.Context, args mapArgs) (*tavily.MapResponse, error) {
	opts := &tavily.MapOptions{
		MaxDepth: args.MaxDepth,
	}
	return t.client.Map(ctx, args.URL, opts)
}
