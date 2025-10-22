package mcp

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

// MCPTemplateRegistry holds all registered MCP templates
var MCPTemplateRegistry = make(map[string]*MCPTemplate)

// RegisterMCPTemplate registers a new MCP template with the given identifier
func RegisterMCPTemplate(identifier string, template *MCPTemplate) {
	MCPTemplateRegistry[identifier] = template
	registry.RegisterAdapterTemplate(identifier, template)
}

// MCPTemplate is a template for creating MCP adapters
type MCPTemplate struct {
	Identifier    string
	DefaultConfig MCPAdapterConfig
}

// CreateAdapter creates a new MCP adapter based on the provided configuration
func (t *MCPTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
	providerInfo := &domain.ProviderAdapterInfo{
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    provider.AuthType,
		Status:      provider.Status,
	}

	// Create adapter configuration with defaults
	adapterConfig := &domain.AdapterConfig{
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RetryBackoff: 1 * time.Second,
		CircuitBreaker: domain.CircuitBreakerConfig{
			FailureThreshold: 5,
			ResetTimeout:     60,
			HalfOpenMaxCalls: 2,
		},
	}

	// Parse MCP-specific configuration
	mcpConfig, err := t.parseMCPConfig(provider.CustomConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MCP configuration: %w", err)
	}

	// Apply template defaults
	if mcpConfig.Timeout == 0 {
		mcpConfig.Timeout = t.DefaultConfig.Timeout
	}
	if mcpConfig.Command == "" {
		mcpConfig.Command = t.DefaultConfig.Command
	}
	if len(mcpConfig.Args) == 0 {
		mcpConfig.Args = t.DefaultConfig.Args
	}

	if len(mcpConfig.Envs) == 0 && len(t.DefaultConfig.Envs) > 0 {
		mcpConfig.Envs = t.DefaultConfig.Envs
	}

	// Merge credential mappings with defaults
	for key, value := range t.DefaultConfig.CredentialMappings {
		if _, exists := mcpConfig.CredentialMappings[key]; !exists {
			mcpConfig.CredentialMappings[key] = value
		}
	}

	// Merge dummy values with defaults
	for key, value := range t.DefaultConfig.DummyCredentials {
		if _, exists := mcpConfig.DummyCredentials[key]; !exists {
			mcpConfig.DummyCredentials[key] = value
		}
	}

	// Merge parameter mappings with defaults
	for key, value := range t.DefaultConfig.ParameterMappings {
		if _, exists := mcpConfig.ParameterMappings[key]; !exists {
			mcpConfig.ParameterMappings[key] = value
		}
	}

	// Merge default parameters with defaults
	for key, value := range t.DefaultConfig.DummyParameters {
		if _, exists := mcpConfig.DummyParameters[key]; !exists {
			mcpConfig.DummyParameters[key] = value
		}
	}

	// Set final default timeout if still not specified
	if mcpConfig.Timeout == 0 {
		mcpConfig.Timeout = 60 * time.Second
	}

	// Initialize credential mappings if not provided
	if mcpConfig.CredentialMappings == nil {
		mcpConfig.CredentialMappings = make(map[string]string)
	}
	if mcpConfig.DummyCredentials == nil {
		mcpConfig.DummyCredentials = make(map[string]string)
	}
	// Initialize parameter mappings if not provided
	if mcpConfig.ParameterMappings == nil {
		mcpConfig.ParameterMappings = make(map[string]string)
	}
	if mcpConfig.DummyParameters == nil {
		mcpConfig.DummyParameters = make(map[string]string)
	}

	// Create base adapter
	baseAdapter := base.NewBaseAdapter(providerInfo, adapterConfig)

	// Create the MCP adapter
	adapter := NewMCPAdapter(baseAdapter, mcpConfig)
	return adapter, nil
}

// ValidateConfig validates the configuration for the MCP template
func (t *MCPTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider configuration cannot be nil")
	}

	if provider.Identifier != t.Identifier {
		return fmt.Errorf("invalid provider identifier, expected '%s', got '%s'", t.Identifier, provider.Identifier)
	}

	// Validate MCP-specific configuration
	mcpConfig, err := t.parseMCPConfig(provider.CustomConfig)
	if err != nil {
		return fmt.Errorf("invalid MCP configuration: %w", err)
	}

	// Validate configuration
	if mcpConfig.Command == "" && t.DefaultConfig.Command == "" {
		return fmt.Errorf("command is required for MCP adapter %s", t.Identifier)
	}

	// Validate auth type is supported
	switch provider.AuthType {
	case types.AuthTypeNone, types.AuthTypeAPIKey:
		// Supported auth types
	default:
		return fmt.Errorf("unsupported auth type: %s", provider.AuthType)
	}

	return nil
}

// parseMCPConfig parses MCP-specific configuration from custom config using JSON marshal/unmarshal
func (t *MCPTemplate) parseMCPConfig(customConfig map[string]interface{}) (*MCPAdapterConfig, error) {
	// Preprocess config to handle multiple field names and ensure proper initialization
	processedConfig := make(map[string]interface{})

	// Copy all fields from original config
	for key, value := range customConfig {
		processedConfig[key] = value
	}

	// Handle field name mapping for dummy credentials
	if dummyCredentials, exists := customConfig["dummy_credentials"]; exists {
		processedConfig["dummy_credentials"] = dummyCredentials
	}

	// Handle field name mapping for dummy parameters
	if dummyParameters, exists := customConfig["dummy_parameters"]; exists {
		processedConfig["dummy_parameters"] = dummyParameters
	}

	// Initialize empty maps if not present to avoid nil pointer issues
	if _, exists := processedConfig["credential_mappings"]; !exists {
		processedConfig["credential_mappings"] = make(map[string]string)
	}
	if _, exists := processedConfig["dummy_credentials"]; !exists {
		processedConfig["dummy_credentials"] = make(map[string]string)
	}
	if _, exists := processedConfig["parameter_mappings"]; !exists {
		processedConfig["parameter_mappings"] = make(map[string]string)
	}
	if _, exists := processedConfig["dummy_parameters"]; !exists {
		processedConfig["dummy_parameters"] = make(map[string]string)
	}

	// Marshal to JSON
	jsonData, err := sonic.Marshal(processedConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	// Unmarshal to MCPAdapterConfig
	var mcpConfig MCPAdapterConfig
	if err := sonic.Unmarshal(jsonData, &mcpConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to MCPAdapterConfig: %w", err)
	}

	return &mcpConfig, nil
}

// DefaultMCPTemplates are Common MCP templates that can be registered
var DefaultMCPTemplates = map[string]*MCPTemplate{
	"sequential_thinking_mcp": {
		Identifier: "sequential_thinking_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "@modelcontextprotocol/server-sequential-thinking@2025.7.1"},
		},
	},
	"brave_search_mcp": {
		Identifier: "brave_search_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "@modelcontextprotocol/server-brave-search@0.6.2"},
			CredentialMappings: map[string]string{
				"apikey": "env:BRAVE_API_KEY",
			},
			DummyCredentials: map[string]string{
				"apikey": "dummy_brave_api_key",
			},
		},
	},
	"newsnow_mcp": {
		Identifier: "newsnow_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "newsnow-mcp-server@0.0.8"},
			Envs: map[string]string{
				"BASE_URL": "https://newsnow.busiyi.world",
			},
		},
	},
	"time_mcp": {
		Identifier: "time_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "uvx",
			Args:    []string{"mcp-server-time@2025.7.1", "--local-timezone", "UTC"},
			Envs: map[string]string{
				"TZ": "UTC", // Fix timezone detection issue
			},
		},
	},

	"context7_mcp": {
		Identifier: "context7_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "@upstash/context7-mcp@1.0.14"},
		},
	},
	"howtocook_mcp": {
		Identifier: "howtocook_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "howtocook-mcp@0.1.1"},
		},
	},
	"airbnb_mcp": {
		Identifier: "airbnb_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "@openbnb/mcp-server-airbnb@0.1.3", "--ignore-robots-txt"},
		},
	},
	"youtube_mcp": {
		Identifier: "youtube_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "youtube-data-mcp-server@1.0.16"},
			Envs: map[string]string{
				"YOUTUBE_TRANSCRIPT_LANG": "en",
			},
			CredentialMappings: map[string]string{
				"apikey": "env:YOUTUBE_API_KEY",
			},
			DummyCredentials: map[string]string{
				"apikey": "dummy_youtube_api_key",
			},
		},
	},
	"yfinance_mcp": {
		Identifier: "yfinance_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "uvx",
			Args:    []string{"mcp-yahoo-finance@0.1.3"},
		},
	},
	"googlemaps_mcp": {
		Identifier: "googlemaps_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "@modelcontextprotocol/server-google-maps@0.6.2"},
			CredentialMappings: map[string]string{
				"apikey": "env:GOOGLE_MAPS_API_KEY",
			},
			DummyCredentials: map[string]string{
				"apikey": "dummy_google_maps_api_key",
			},
		},
	},
	"minimax_mcp": {
		Identifier: "minimax_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "uvx",
			Args:    []string{"minimax-mcp@0.0.17", "-y"},
			Envs: map[string]string{
				"MINIMAX_API_HOST": "https://api.minimaxi.chat",
			},
			CredentialMappings: map[string]string{
				"apikey": "env:MINIMAX_API_KEY",
			},
			DummyCredentials: map[string]string{
				"apikey": "dummy_minimax_api_key",
			},
		},
	},
	"aws_pricing_mcp": {
		Identifier: "aws_pricing_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "uvx",
			Args:    []string{"awslabs.aws-pricing-mcp-server@1.0.6"},
			Envs: map[string]string{
				"AWS_REGION":        "us-east-1",
				"FASTMCP_LOG_LEVEL": "ERROR",
			},
			CredentialMappings: map[string]string{
				"apikey": "env:AWS_ACCESS_KEY_ID",
			},
			DummyCredentials: map[string]string{
				"apikey": "dummy_aws_access_key_id",
			},
		},
	},
	"asana_mcp": {
		Identifier: "asana_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "@roychri/mcp-server-asana@1.7.0"},
			CredentialMappings: map[string]string{
				"access_token": "env:ASANA_ACCESS_TOKEN",
			},
			DummyCredentials: map[string]string{
				"access_token": "dummy_asana_access_token",
			},
		},
	},
	"cloudflare_bindings_mcp": {
		Identifier: "cloudflare_bindings_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"mcp-remote@0.1.18", "https://bindings.mcp.cloudflare.com/sse"},
		},
	},
	"firecrawl_mcp": {
		Identifier: "firecrawl_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "firecrawl-mcp@1.12.0"},
			CredentialMappings: map[string]string{
				"apikey": "env:FIRECRAWL_API_KEY",
			},
			DummyCredentials: map[string]string{
				"apikey": "dummy_firecrawl_api_key",
			},
		},
	},
	"todoist_mcp": {
		Identifier: "todoist_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "@abhiz123/todoist-mcp-server@0.1.0"},
			CredentialMappings: map[string]string{
				"api_token": "env:TODOIST_API_TOKEN",
			},
			DummyCredentials: map[string]string{
				"api_token": "dummy_todoist_api_token",
			},
		},
	},
	"postgre_sql_mcp": {
		Identifier: "postgre_sql_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "@modelcontextprotocol/server-postgres@0.6.2", "postgresql://localhost/mydb"},
		},
	},
	"duckduckgo_mcp": {
		Identifier: "duckduckgo_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "uvx",
			Args:    []string{"duckduckgo-mcp-server@0.1.1"},
		},
	},
	"akshare_mcp": {
		Identifier: "akshare_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "uvx",
			Args:    []string{"akshare-one-mcp@0.2.3"},
		},
	},
	"arxiv_mcp": {
		Identifier: "arxiv_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "uv",
			Args:    []string{"tool", "run", "arxiv-paper-mcp@0.1.2"},
		},
	},
	"exa_mcp": {
		Identifier: "exa_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "exa-mcp-server@2.0.3"},
			CredentialMappings: map[string]string{
				"apikey": "env:EXA_API_KEY",
			},
			DummyCredentials: map[string]string{
				"apikey": "dummy_exa_api_key",
			},
		},
	},
	"everart_mcp": {
		Identifier: "everart_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "@modelcontextprotocol/server-everart@0.6.2"},
			CredentialMappings: map[string]string{
				"apikey": "env:EVERART_API_KEY",
			},
			DummyCredentials: map[string]string{
				"apikey": "dummy_everart_api_key",
			},
		},
	},
	"huggingface_mcp": {
		Identifier: "huggingface_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "uvx",
			Args:    []string{"huggingface-mcp-server@0.1.0"},
		},
	},
	"wikipedia_mcp": {
		Identifier: "wikipedia_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "uvx",
			Args:    []string{"wikipedia-mcp@1.5.5"},
		},
	},
	"calculator_mcp": {
		Identifier: "calculator_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "uvx",
			Args:    []string{"mcp-server-calculator@0.1.1"},
		},
	},
	"jina_mcp": {
		Identifier: "jina_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"jina-mcp-tools@1.1.1"},
			CredentialMappings: map[string]string{
				"apikey": "env:JINA_API_KEY",
			},
			DummyCredentials: map[string]string{
				"apikey": "dummy_jina_api_key",
			},
		},
	},
	"baidu_map_mcp": {
		Identifier: "baidu_map_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "@baidumap/mcp-server-baidu-map@1.0.5"},
			CredentialMappings: map[string]string{
				"apikey": "env:BAIDU_MAP_API_KEY",
			},
			DummyCredentials: map[string]string{
				"apikey": "dummy_baidu_map_api_key",
			},
		},
	},
	"google_scholar_mcp": {
		Identifier: "google_scholar_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "uvx",
			Args:    []string{"google-scholar-mcp-server@0.1.3"},
		},
	},
	"search1api_mcp": {
		Identifier: "search1api_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "search1api-mcp@0.1.8"},
			CredentialMappings: map[string]string{
				"apikey": "env:SEARCH1API_KEY",
			},
			DummyCredentials: map[string]string{
				"apikey": "dummy_search1api_api_key",
			},
		},
	},
	"gitlab_mcp": {
		Identifier: "gitlab_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "@modelcontextprotocol/server-gitlab@2025.4.25"},
			Envs: map[string]string{
				"GITLAB_API_URL": "https://gitlab.com/api/v4",
			},
			CredentialMappings: map[string]string{
				"apikey": "env:GITLAB_PERSONAL_ACCESS_TOKEN",
			},
			DummyCredentials: map[string]string{
				"apikey": "dummy_gitlab_personal_access_token",
			},
		},
	},
	"antv_mcp": {
		Identifier: "antv_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "@antv/mcp-server-chart@0.8.0"},
		},
	},
	"tavily_mcp": {
		Identifier: "tavily_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "tavily-mcp@0.2.9"},
		},
	},
	"web3_research_mcp": {
		Identifier: "web3_research_mcp",
		DefaultConfig: MCPAdapterConfig{
			Command: "npx",
			Args:    []string{"-y", "web3-research-mcp@1.0.1"},
		},
	},
}

// init registers all default MCP templates
func init() {
	for _, template := range DefaultMCPTemplates {
		RegisterMCPTemplate(template.Identifier, template)
	}
}
