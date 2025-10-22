package eodhd

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	identifier = "eodhd"
	baseURL    = "https://eodhd.com/api/"

	apikeyParamName = "api_token"
)

// init registers EodhdTemplate when the package is imported
func init() {
	// Type assertion to ensure the adapter implements the necessary interfaces
	var _ domain.APIKeyAdapter = (*EodhdAdapter)(nil)

	template := &EodhdTemplate{}
	registry.RegisterAdapterTemplate(identifier, template)
}

// EodhdTemplate implements domain.AdapterTemplate
// It is responsible for creating EodhdAdapter instances based on configuration parsed from eodhd.json.
type EodhdTemplate struct{}

// CreateAdapter builds a new EodhdAdapter instance from the configuration map.
// Implements the domain.AdapterTemplate interface.
func (t *EodhdTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
	providerInfo := &domain.ProviderAdapterInfo{
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    provider.AuthType,
	}

	// Create Base Adapter Configuration (Extract or use defaults)
	adapterConfig := &domain.AdapterConfig{
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RetryBackoff: 1 * time.Second,
		// Add RateLimit, CircuitBreaker from config if available
	}

	restConfig := &rest.RESTConfig{
		BaseURL: baseURL,
	}
	restAdapterInstance := rest.NewRESTAdapter(providerInfo, adapterConfig, restConfig)

	// Create EodhdAdapter instance with extracted configuration
	// Note: NewEodhdAdapter no longer hardcodes baseURL and apiKeyConfig internally
	adapter := NewEodhdAdapter(
		providerInfo,
		adapterConfig,
		restAdapterInstance,
	)

	return adapter, nil
}

// ValidateConfig checks if the provided configuration map contains the necessary fields for creating an EodhdAdapter.
// Implements the domain.AdapterTemplate interface.
func (t *EodhdTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider model cannot be nil")
	}

	if provider.Identifier != identifier {
		return fmt.Errorf("invalid provider identifier, must be '%s'", identifier)
	}

	if provider.AuthType != types.AuthTypeAPIKey {
		return fmt.Errorf("invalid or missing auth_type, must be 'apikey'")
	}

	return nil
}
