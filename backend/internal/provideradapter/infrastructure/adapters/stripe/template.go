package stripe

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	identifier = "stripe"
	baseURL    = "https://api.stripe.com/v1"

	apikeyParamName = "Authorization"
)

// Register the Stripe adapter template during package initialization.
func init() {

	// StripeAdapter now uses APIKey, so it should implement APIKeyAdapter if such an interface exists and is intended.
	// The linter error indicates `domain.APIKeyAdapter` exists and GetAPIKeyConfig() *domain.APIKeyConfig is expected.
	var _ domain.APIKeyAdapter = (*StripeAdapter)(nil)

	template := &StripeTemplate{}
	registry.RegisterAdapterTemplate(identifier, template)
}

// StripeAdapterTemplate implements the AdapterTemplate interface for Stripe.
type StripeTemplate struct{}

// CreateAdapter creates a new Stripe adapter instance from the provided configuration.
func (t *StripeTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
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
	restAdapterInstance := rest.NewRESTAdapter(providerInfo, adapterConfig, restConfig)

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

	adapter := NewStripeAdapter(
		providerInfo,
		adapterConfig,
		restAdapterInstance,
		permissions,
	)

	return adapter, nil
}

// ValidateConfig checks if the provided configuration map contains the necessary fields for an APIKey REST adapter.
func (t *StripeTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider is required")
	}

	if provider.Identifier != identifier {
		return fmt.Errorf("invalid provider identifier, must be '%s'", identifier)
	}

	if provider.AuthType != types.AuthTypeAPIKey {
		return fmt.Errorf("invalid auth_type, must be 'apikey'")
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
