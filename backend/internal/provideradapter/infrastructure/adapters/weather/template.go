package openweathermap

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	providercore "github.com/context-space/context-space/backend/internal/providercore/domain"
)

const (
	identifier = "openweathermap"
	baseURL    = "https://api.openweathermap.org"
)

// Register adapter template during package initialization
func init() {
	// Type assertion to ensure adapter implements necessary interfaces
	var _ domain.APIKeyAdapter = (*OpenWeatherMapAdapter)(nil)

	template := &OpenWeatherMapAdapterTemplate{}
	registry.RegisterAdapterTemplate(identifier, template)
}

// OpenWeatherMapAdapterTemplate implements AdapterTemplate interface
type OpenWeatherMapAdapterTemplate struct {
	// Add template-specific configuration here if needed
}

// CreateAdapter creates a new adapter instance from provided configuration
func (t *OpenWeatherMapAdapterTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
	// Validate configuration
	if err := t.ValidateConfig(provider); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create provider info
	providerInfo := &domain.ProviderAdapterInfo{
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    provider.AuthType,
		Permissions: provider.Permissions,
		Operations:  provider.Operations,
		Status:      provider.Status,
	}

	// Create adapter configuration
	adapterConfig := &domain.AdapterConfig{
		Timeout:      10 * time.Second,
		MaxRetries:   3,
		RetryBackoff: 1 * time.Second,
	}

	// Create REST configuration
	restConfig := &rest.RESTConfig{
		BaseURL: baseURL,
	}

	// Create REST adapter
	restAdapter := rest.NewRESTAdapter(providerInfo, adapterConfig, restConfig)

	// Create OpenWeatherMap adapter
	adapter := NewOpenWeatherMapAdapter(providerInfo, adapterConfig, restAdapter)

	return adapter, nil
}

// ValidateConfig validates provider configuration
func (t *OpenWeatherMapAdapterTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider config cannot be nil")
	}

	// Validate identifier
	if provider.Identifier != identifier {
		return fmt.Errorf("invalid identifier: expected %s, got %s", identifier, provider.Identifier)
	}

	// Validate auth type
	if provider.AuthType == "" {
		return fmt.Errorf("auth type is required")
	}

	if provider.AuthType != providercore.AuthTypeAPIKey {
		return fmt.Errorf("unsupported auth type: %s", provider.AuthType)
	}

	// Validate operations
	if len(provider.Operations) == 0 {
		return fmt.Errorf("at least one operation must be defined")
	}

	// Validate custom configuration
	if provider.CustomConfig != nil {
		// API key validation
		if apiKey, exists := provider.CustomConfig["api_key"]; exists {
			if apiKeyStr, ok := apiKey.(string); !ok || apiKeyStr == "" {
				return fmt.Errorf("api_key must be a non-empty string")
			}
		}
	}

	return nil
}

// GetSupportedOperations returns operations supported by this adapter template
func (t *OpenWeatherMapAdapterTemplate) GetSupportedOperations() []string {
	return []string{
		"get_current_weather",
		"get_weather_forecast",
		"get_one_call_weather",
		"get_air_pollution",
		"get_geocoding",
	}
}
