package fetch

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	identifier = "fetch"
	baseURL    = "https://scrape.serper.dev"

	apikeyParamName = "X-API-KEY"
)

// Register the Fetch adapter template during package initialization
func init() {
	// Type assertion to ensure the adapter implements the necessary interfaces
	var _ domain.APIKeyAdapter = (*FetchAdapter)(nil)

	template := &FetchTemplate{}
	registry.RegisterAdapterTemplate(identifier, template)
}

// FetchTemplate implements the AdapterTemplate interface for Fetch
type FetchTemplate struct{}

// CreateAdapter creates a new Fetch adapter instance from the provided configuration
func (t *FetchTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
	providerInfo := &domain.ProviderAdapterInfo{
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    provider.AuthType,
	}

	// Create Base Adapter Configuration
	adapterConfig := &domain.AdapterConfig{
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RetryBackoff: 1 * time.Second,
	}

	restConfig := &rest.RESTConfig{
		BaseURL: baseURL,
	}
	restAdapterInstance := rest.NewRESTAdapter(providerInfo, adapterConfig, restConfig)

	apiKey := provider.CustomConfig["api_key"].(string)
	adapter := NewFetchAdapter(
		providerInfo,
		adapterConfig,
		restAdapterInstance,
		&opDefaults,
		apiKey,
	)

	return adapter, nil
}

// ValidateConfig checks if the provided configuration contains all required fields
func (t *FetchTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider model cannot be nil")
	}

	if provider.Identifier != identifier {
		return fmt.Errorf("invalid provider identifier, must be '%s'", identifier)
	}

	if provider.AuthType != types.AuthTypeNone {
		return fmt.Errorf("invalid or missing auth_type, must be 'apikey'")
	}

	apiKey, ok := provider.CustomConfig["api_key"]
	if !ok {
		return fmt.Errorf("api_key is required")
	}

	if apiKey.(string) == "" {
		return fmt.Errorf("api_key is required")
	}

	return nil
}
