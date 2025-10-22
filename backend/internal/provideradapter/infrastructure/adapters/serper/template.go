package serper

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	identifier = "serper"
	baseURL    = "https://google.serper.dev"

	apikeyParamName = "X-API-KEY"
)

// init registers SerperTemplate when the package is imported.
func init() {

	var _ domain.APIKeyAdapter = (*SerperAdapter)(nil)

	template := &SerperTemplate{}
	registry.RegisterAdapterTemplate(identifier, template) // Use "serper" identifier
}

// SerperTemplate implements domain.AdapterTemplate.
// It's responsible for creating SerperAdapter instances from parsed serper.json configuration.
type SerperTemplate struct{}

// CreateAdapter constructs a new SerperAdapter instance from a configuration map.
// Implements the domain.AdapterTemplate interface.
func (t *SerperTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {

	providerInfo := &domain.ProviderAdapterInfo{
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
	restAdapter := rest.NewRESTAdapter(providerInfo, adapterConfig, restConfig)

	permissionsData := provider.Permissions

	permissions := make(domain.PermissionSet, len(permissionsData))
	for _, permMap := range permissionsData {
		scopes := permMap.OAuthScopes

		permissions[permMap.Identifier] = *domain.NewPermission(
			permMap.Identifier,
			permMap.Name,
			permMap.Description,
			scopes,
		)
	}

	adapter := NewSerperAdapter(
		providerInfo,
		adapterConfig,
		restAdapter,
		permissions,
	)

	return adapter, nil
}

// ValidateConfig checks if the provided configuration map contains the necessary fields
// to create a SerperAdapter.
// Implements the domain.AdapterTemplate interface.
func (t *SerperTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider is required")
	}

	if provider.Identifier != identifier {
		return fmt.Errorf("invalid provider identifier, must be '%s'", identifier)
	}

	if provider.AuthType != types.AuthTypeAPIKey {
		return fmt.Errorf("invalid 'auth_type' in config (must be 'apikey')")
	}

	// Validate permissions structure
	permissionsData := provider.Permissions
	if len(permissionsData) == 0 {
		return fmt.Errorf("missing or invalid 'permissions' field, must be an array")
	}
	// Validate each permission
	for i, p := range permissionsData {
		if p.Identifier == "" {
			return fmt.Errorf("missing or invalid identifier in permission at index %d", i)
		}
		if p.Name == "" {
			return fmt.Errorf("missing or invalid name in permission at index %d", i)
		}
		if p.Description == "" {
			return fmt.Errorf("missing or invalid description in permission at index %d", i)
		}
	}
	return nil
}
