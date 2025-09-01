package main
import (
	"fmt"
	"log"
	"os"
	"github.com/denkhaus/agents/pkg/config"
	providerConfig "github.com/denkhaus/agents/pkg/provider/config"
	"github.com/samber/do"
)
func main() {
	os.Unsetenv("TAVILY_API_KEY")
	injector := do.New()
	do.Provide(injector, func(i *do.Injector) (config.Service, error) {
		return config.NewWithDI(i)
	})
	provider, err := providerConfig.NewCUEConfigProvider(injector)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}
	fmt.Println("Testing fallback (no TAVILY_API_KEY)...")
	agentConfig, err := provider.LoadAgentComposition("development", "researcher")
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}
	if tavilyConfig, exists := agentConfig.Tool.ToolSets["tavily_toolset"]; exists {
		if apiKey, ok := tavilyConfig.Config["api_key"]; ok {
			fmt.Printf("✅ Fallback works: %v\n", apiKey)
		}
	}
}
