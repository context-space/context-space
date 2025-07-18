package mcp

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	providercore "github.com/context-space/context-space/backend/internal/providercore/domain"
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
		Permissions: provider.Permissions,
		Operations:  provider.Operations,
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
	case providercore.AuthTypeNone, providercore.AuthTypeAPIKey:
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
}

// init registers all default MCP templates
func init() {
	for _, template := range DefaultMCPTemplates {
		RegisterMCPTemplate(template.Identifier, template)
	}
}
