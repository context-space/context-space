package github

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/shared/types"

	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
)

const (
	identifier = "github"
)

// Register the GitHub adapter template
func init() {

	var _ domain.OAuthAdapter = (*GitHubAdapter)(nil)

	template := &GitHubTemplate{}

	// Register with the central registry
	registry.RegisterAdapterTemplate(identifier, template)
}

// GitHubTemplate is a template for creating GitHub adapters
type GitHubTemplate struct {
	// Template configuration could be injected here if needed
	// Currently using hardcoded values for GitHub
}

// CreateAdapter builds a new GitHubAdapter instance from the provider model
// Implements application.AdapterTemplate interface
func (t *GitHubTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
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
		RateLimit:    5000,
		RatePeriod:   time.Hour,
		CircuitBreaker: domain.CircuitBreakerConfig{
			FailureThreshold: 5,
			ResetTimeout:     60,
			HalfOpenMaxCalls: 2,
		},
	}

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

	adapter := NewGitHubAdapter(
		providerAdapterInfo,
		adapterConfig,
		oauthConfig,
		permissions,
	)

	return adapter, nil
}

// ValidateConfig checks if the provided provider model contains the required fields for creating a GitHubAdapter
// Implements application.AdapterTemplate interface
func (t *GitHubTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
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
