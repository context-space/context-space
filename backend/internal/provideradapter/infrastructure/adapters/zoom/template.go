package zoom

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	identifier = "zoom"
	baseURL    = "https://api.zoom.us/v2"

	authURL  = "https://zoom.us/oauth/authorize"
	tokenURL = "https://zoom.us/oauth/token"
)

// Register the Zoom adapter template during package initialization.
func init() {
	// Type assertion to ensure the adapter implements the necessary interfaces
	var _ domain.OAuthAdapter = (*ZoomAdapter)(nil)

	template := &ZoomAdapterTemplate{}
	registry.RegisterAdapterTemplate(identifier, template)
}

// ZoomAdapterTemplate implements the AdapterTemplate interface for Zoom.
type ZoomAdapterTemplate struct {
	// Configuration specific to this template could be added here if needed.
}

// CreateAdapter creates a new Zoom adapter instance from the provided configuration.
func (t *ZoomAdapterTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
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

	oauthConfigData := provider.OAuthConfig

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

	restConfig := &rest.RESTConfig{
		BaseURL: baseURL,
	}
	restAdapterInstance := rest.NewRESTAdapter(providerInfo, adapterConfig, restConfig)

	adapter := NewZoomAdapter(
		providerInfo,
		adapterConfig,
		oauthConfigData,
		restAdapterInstance,
		permissions,
	)

	return adapter, nil
}

// ValidateConfig checks if the provided configuration map contains the necessary fields for an OAuth REST adapter.
func (t *ZoomAdapterTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider model cannot be nil")
	}

	if provider.Identifier != identifier {
		return fmt.Errorf("invalid provider identifier, must be '%s'", identifier)
	}

	// Validate auth_type specifically for OAuth
	if provider.AuthType != types.AuthTypeOAuth {
		return fmt.Errorf("invalid or missing auth_type, must be 'oauth'")
	}

	oauthConfig := provider.OAuthConfig
	if oauthConfig == nil {
		return fmt.Errorf("missing or invalid oauth_config section")
	}
	// Use the struct's validation method
	if err := oauthConfig.Validate(); err != nil {
		return fmt.Errorf("invalid oauth_config: %w", err)
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
		if len(p.OAuthScopes) == 0 {
			return fmt.Errorf("missing or invalid oauth_scopes in permission at index %d", i)
		}
	}
	return nil
}
