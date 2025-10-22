package amap

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	identifier = "amap"
	baseURL    = "https://restapi.amap.com"

	apikeyParamName = "key"
)

// Register adapter template during package initialization
func init() {
	// Type assertion to ensure adapter implements necessary interfaces
	var _ domain.APIKeyAdapter = (*AmapAdapter)(nil)

	template := &AmapAdapterTemplate{}
	registry.RegisterAdapterTemplate(identifier, template)
}

// AmapAdapterTemplate implements AdapterTemplate interface
type AmapAdapterTemplate struct {
	// Template-specific configuration can be added here if needed
}

// CreateAdapter creates a new adapter instance from the provided configuration
func (t *AmapAdapterTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {

	// Create provider information
	providerInfo := &domain.ProviderAdapterInfo{
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    provider.AuthType,
	}

	// Create adapter configuration
	adapterConfig := &domain.AdapterConfig{
		Timeout:      30 * time.Second, // Amap API default timeout 30 seconds
		MaxRetries:   3,
		RetryBackoff: 1 * time.Second,
	}

	// Create REST configuration
	restConfig := &rest.RESTConfig{
		BaseURL: baseURL,
	}

	// Create REST adapter
	restAdapter := rest.NewRESTAdapter(providerInfo, adapterConfig, restConfig)

	// Create main adapter (API Key version, no OAuth configuration required)
	adapter := NewAmapAdapter(
		providerInfo,
		adapterConfig,
		restAdapter,
	)

	return adapter, nil
}

// ValidateConfig checks if the provided configuration contains necessary fields
// Implement the AdapterTemplate interface
func (t *AmapAdapterTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider is required")
	}

	if provider.Identifier != identifier {
		return fmt.Errorf("identifier must be '%s'", identifier)
	}

	if provider.AuthType != types.AuthTypeAPIKey {
		return fmt.Errorf("invalid auth_type, must be 'apikey'")
	}

	return nil
}
