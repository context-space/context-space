package openweathermap

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	identifier = "openweathermap"
	baseURL    = "https://api.openweathermap.org"

	apikeyParamName = "appid"
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

	if provider.AuthType != types.AuthTypeAPIKey {
		return fmt.Errorf("invalid auth_type, must be 'apikey'")
	}

	return nil
}
