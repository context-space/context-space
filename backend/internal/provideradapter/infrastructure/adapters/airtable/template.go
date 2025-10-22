package airtable

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	identifier = "airtable"
	baseURL    = "https://api.airtable.com/v0"

	authURL  = "https://airtable.com/oauth2/v1/authorize"
	tokenURL = "https://airtable.com/oauth2/v1/token"
)

// Register the Airtable adapter template during package initialization.
func init() {
	// Type assertion to ensure the adapter implements the necessary interfaces
	var _ domain.OAuthAdapter = (*AirtableAdapter)(nil)

	template := &AirtableTemplate{}
	registry.RegisterAdapterTemplate(identifier, template)
}

// AirtableAdapterTemplate implements the AdapterTemplate interface for Airtable.
type AirtableTemplate struct {
	// Configuration specific to this template could be added here if needed.
}

// CreateAdapter creates a new Airtable adapter instance from the provided configuration.
func (t *AirtableTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {

	providerInfo := &domain.ProviderAdapterInfo{
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    provider.AuthType,
		Status:      provider.Status,
	}

	adapterConfig := &domain.AdapterConfig{
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RetryBackoff: 1 * time.Second,
	}

	// 4. Create the underlying REST Adapter instance
	restConfig := &rest.RESTConfig{
		BaseURL: baseURL,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	restAdapterInstance := rest.NewRESTAdapter(providerInfo, adapterConfig, restConfig)

	permissionsData := provider.Permissions

	// Convert permissions to coreDomain.Permission slice
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

	adapter := NewAirtableAdapter(
		providerInfo,
		adapterConfig,
		oauthConfig,
		restAdapterInstance,
		permissions,
	)
	return adapter, nil
}

// ValidateConfig checks if the provided configuration map contains the necessary fields for an OAuth REST adapter.
func (t *AirtableTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
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
