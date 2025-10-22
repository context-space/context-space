package search

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/types"

	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
)

const (
	identifier = "search"
	baseURL    = "https://google.serper.dev"

	apikeyParamName = "X-API-KEY"
)

// init registers SearchTemplate when the package is imported
func init() {
	var _ domain.APIKeyAdapter = (*SearchAdapter)(nil)

	// Create template instance
	template := &SearchTemplate{}
	// Register using provider's unique identifier
	registry.RegisterAdapterTemplate(identifier, template)
}

// SearchTemplate implements application.AdapterTemplate
type SearchTemplate struct {
	// If the template itself needs configuration, fields can be added here, but it's usually stateless
}

// CreateAdapter builds a new SearchAdapter instance from the configuration map
// Implements application.AdapterTemplate interface
func (t *SearchTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
	providerAdapterInfo := &domain.ProviderAdapterInfo{
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    provider.AuthType,
	}

	adapterConfig := &domain.AdapterConfig{
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RetryBackoff: 1 * time.Second,
	}

	restConfig := &rest.RESTConfig{
		BaseURL: baseURL,
	}
	restAdapter := rest.NewRESTAdapter(providerAdapterInfo, adapterConfig, restConfig)

	apiKey := provider.CustomConfig["api_key"].(string)
	adapter := NewSearchAdapter(
		providerAdapterInfo,
		adapterConfig,
		restAdapter,
		apiKey,
	)
	return adapter, nil
}

// ValidateConfig checks if the provided configuration map contains fields required to create a SearchAdapter
// Implements application.AdapterTemplate interface
func (t *SearchTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider is required")
	}

	if provider.Identifier != identifier {
		return fmt.Errorf("invalid provider identifier, must be '%s'", identifier)
	}

	// Validate auth_type type and value (must be "none")
	if provider.AuthType != types.AuthTypeNone {
		return fmt.Errorf("invalid auth_type, must be 'none'")
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
