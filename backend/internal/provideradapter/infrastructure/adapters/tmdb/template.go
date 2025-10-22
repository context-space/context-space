package tmdb

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	identifier = "tmdb"
	baseURL    = "https://api.themoviedb.org/3"

	apikeyParamName = "api_key"
)

// Register the TMDB adapter template during package initialization
func init() {
	// Type assertion to ensure the adapter implements the necessary interfaces
	var _ domain.APIKeyAdapter = (*TmdbAdapter)(nil)

	template := &TmdbTemplate{}
	registry.RegisterAdapterTemplate(identifier, template)
}

// TmdbTemplate implements the AdapterTemplate interface for TMDB
type TmdbTemplate struct{}

// CreateAdapter creates a new TMDB adapter instance from the provided configuration
func (t *TmdbTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
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

	adapter := NewTmdbAdapter(
		providerInfo,
		adapterConfig,
		restAdapterInstance,
	)

	return adapter, nil
}

// ValidateConfig checks if the provided configuration contains all required fields
// Implements application.AdapterTemplate interface
func (t *TmdbTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider is required")
	}

	if provider.Identifier != identifier {
		return fmt.Errorf("invalid provider identifier, must be '%s'", identifier)
	}

	if provider.AuthType != types.AuthTypeAPIKey {
		return fmt.Errorf("invalid auth_type, must be 'apikey'")
	}

	return nil
}
