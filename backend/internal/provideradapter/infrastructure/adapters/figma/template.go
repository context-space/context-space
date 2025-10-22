package figma

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	identifier = "figma"
	baseURL    = "https://api.figma.com/v1"

	authURL  = "https://www.figma.com/oauth"
	tokenURL = "https://api.figma.com/v1/oauth/token"
)

var operationDefaults = OperationDefaults{
	GetTeamStyles: GetTeamStylesDefaults{
		PageSize: 30,
	},
	GetTeamComponents: GetTeamComponentsDefaults{
		PageSize: 30,
	},
	GetTeamComponentSets: GetTeamComponentSetsDefaults{
		PageSize: 30,
	},
	GetImageRenders: GetImageRendersDefaults{
		Scale:  1,
		Format: "png",
	},
}

// Register the Figma adapter template during package initialization.
func init() {
	// Type assertion to ensure the adapter implements the necessary interfaces
	var _ domain.OAuthAdapter = (*FigmaAdapter)(nil)

	template := &FigmaTemplate{}
	registry.RegisterAdapterTemplate(identifier, template)
}

// FigmaAdapterTemplate implements the AdapterTemplate interface for Figma.
type FigmaTemplate struct {
	// Configuration specific to this template could be added here if needed.
}

// CreateAdapter creates a new Figma adapter instance from the provided configuration.
func (t *FigmaTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
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
	}

	// Create REST Adapter Instance
	restConfig := &rest.RESTConfig{
		BaseURL: baseURL,
		Headers: map[string]string{},
	}
	restAdapterInstance := rest.NewRESTAdapter(providerInfo, adapterConfig, restConfig)

	// Extract permissions
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

	oauthConfig := provider.OAuthConfig.Clone()

	// Create the main Adapter Instance using its constructor
	adapter := NewFigmaAdapter(
		providerInfo,
		adapterConfig,
		oauthConfig,
		restAdapterInstance,
		&operationDefaults,
		permissions,
	)

	return adapter, nil
}

// ValidateConfig checks if the provided configuration map contains the necessary fields for an OAuth REST adapter.
func (t *FigmaTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider model cannot be nil")
	}

	if provider.Identifier != identifier {
		return fmt.Errorf("invalid provider identifier, must be '%s'", identifier)
	}

	if provider.AuthType != types.AuthTypeOAuth {
		return fmt.Errorf("invalid or missing auth_type, must be 'oauth'")
	}

	// Validate oauth_config structure
	oauthConfig := provider.OAuthConfig
	if oauthConfig == nil {
		return fmt.Errorf("oauth_config is required")
	}

	// Validate OAuth configuration
	if err := oauthConfig.Validate(); err != nil {
		return err
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
