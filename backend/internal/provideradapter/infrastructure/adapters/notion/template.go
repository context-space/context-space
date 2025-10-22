package notion

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	identifier    = "notion"
	baseURL       = "https://api.notion.com/v1"
	notionVersion = "2022-06-28"

	authURL  = "https://api.notion.com/v1/oauth/authorize"
	tokenURL = "https://api.notion.com/v1/oauth/token"
)

var opDefaults = OperationDefaults{
	Search: SearchDefaults{
		PageSize:      50,
		SortTimestamp: "last_edited_time",
		SortDirection: "descending",
	},
	QueryDatabase: QueryDatabaseDefaults{
		PageSize: 50,
	},
	ListUsers: ListUsersDefaults{
		PageSize: 100,
	},
	ListDatabases: ListDatabasesDefaults{
		PageSize:       50,
		FilterValue:    "database",
		FilterProperty: "object",
	},
}

// Register the Notion adapter template during package initialization.
func init() {

	var _ domain.OAuthAdapter = (*NotionAdapter)(nil)

	template := &NotionAdapterTemplate{}
	registry.RegisterAdapterTemplate("notion", template)
}

// NotionAdapterTemplate implements the AdapterTemplate interface for Notion.
type NotionAdapterTemplate struct {
	// Configuration specific to this template could be added here if needed.
}

// CreateAdapter creates a new Notion adapter instance from the provided configuration.
func (t *NotionAdapterTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
	providerInfo := &domain.ProviderAdapterInfo{
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    provider.AuthType,
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

	oauthConfigData := provider.OAuthConfig.Clone()

	adapterConfig := &domain.AdapterConfig{
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RetryBackoff: 1 * time.Second,
		RateLimit:    5000, // Example rate limit
		RatePeriod:   time.Hour,
		CircuitBreaker: domain.CircuitBreakerConfig{
			FailureThreshold: 5,
			ResetTimeout:     60,
			HalfOpenMaxCalls: 2,
		},
	}

	restConfig := &rest.RESTConfig{
		BaseURL: baseURL,
	}
	restAdapter := rest.NewRESTAdapter(providerInfo, adapterConfig, restConfig)

	adapter := NewNotionAdapter(
		providerInfo,
		adapterConfig,
		oauthConfigData,
		notionVersion,
		&opDefaults,
		permissions,
		restAdapter,
	)

	return adapter, nil
}

// ValidateConfig checks if the provided configuration map contains the necessary fields for Notion.
func (t *NotionAdapterTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider is required")
	}

	if provider.Identifier != identifier {
		return fmt.Errorf("invalid provider identifier, must be '%s'", identifier)
	}

	// Validate auth_type
	if provider.AuthType != types.AuthTypeOAuth {
		return fmt.Errorf("invalid auth_type, must be 'oauth'")
	}

	// Validate oauth_config structure
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
