// slack/template.go.tmpl - Template for generating adapter template file
package slack

import (
	"fmt"
	"time" // Assuming resilience settings might use time

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest" // Needed to create restAdapter
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	identifier = "slack"
	baseURL    = "https://slack.com/api"

	authURL  = "https://slack.com/oauth/v2/authorize"
	tokenURL = "https://slack.com/api/oauth.v2.access"
)

var opDefaults = OperationDefaults{
	OperationIdentifier1: OperationIdentifier1Defaults{
		Filter: "default_filter_value",
	},
	OperationIdentifier2: OperationIdentifier2Defaults{
		Page: "default_page_value",
	},
}

// Register the Slack adapter template during package initialization.
func init() {
	// Type assertion to ensure the adapter implements the necessary interfaces
	var _ domain.OAuthAdapter = (*SlackAdapter)(nil)

	template := &SlackTemplate{}
	registry.RegisterAdapterTemplate("slack", template)
}

// SlackAdapterTemplate implements the AdapterTemplate interface for Slack.
type SlackTemplate struct {
	// Configuration specific to this template could be added here if needed.
}

// CreateAdapter creates a new Slack adapter instance from the provided configuration.
func (t *SlackTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
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
		permissions[permMap.Identifier] = permMap
	}

	oauthConfig := provider.OAuthConfig.Clone()

	adapter := NewSlackAdapter(
		providerInfo,
		adapterConfig,
		oauthConfig,
		restAdapter,
		&opDefaults,
		permissions,
	)

	return adapter, nil
}

// ValidateConfig checks if the provided configuration map contains the necessary fields for an OAuth REST adapter.
func (t *SlackTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider is required")
	}

	if provider.Identifier != identifier {
		return fmt.Errorf("invalid provider identifier, must be '%s'", identifier)
	}

	if provider.AuthType != types.AuthTypeOAuth {
		return fmt.Errorf("invalid auth_type, must be 'oauth'")
	}

	oauthConfig := provider.OAuthConfig
	if oauthConfig == nil {
		return fmt.Errorf("missing or invalid oauth_config")
	}
	if err := oauthConfig.Validate(); err != nil {
		return fmt.Errorf("[%s] ValidateConfig: invalid oauth_config: %w", provider.Identifier, err)
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
