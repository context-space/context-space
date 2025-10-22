package notion

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	"github.com/context-space/context-space/backend/internal/shared/utils"
)

// SearchDefaults defines default parameters for search operations.
type SearchDefaults struct {
	PageSize      int    `mapstructure:"page_size"`
	SortTimestamp string `mapstructure:"timestamp"` // Field name for timestamp sorting (e.g., "last_edited_time")
	SortDirection string `mapstructure:"direction"` // Sort direction ("ascending" or "descending")
}

// QueryDatabaseDefaults defines default parameters for database query operations.
type QueryDatabaseDefaults struct {
	PageSize int `mapstructure:"page_size"`
}

// ListUsersDefaults defines default parameters for listing users.
type ListUsersDefaults struct {
	PageSize int `mapstructure:"page_size"`
}

// ListDatabasesDefaults defines default parameters for listing databases.
type ListDatabasesDefaults struct {
	PageSize       int    `mapstructure:"page_size"`
	FilterValue    string `mapstructure:"value"`    // Default filter value (e.g., "database")
	FilterProperty string `mapstructure:"property"` // Default filter property (e.g., "object")
}

// CreatePageDefaults defines default parameters for creating pages (currently none).
type CreatePageDefaults struct {
	// No specific defaults defined in notion.json yet
}

// AppendToBlockDefaults defines default parameters for appending blocks (currently none).
type AppendToBlockDefaults struct {
	// No specific defaults defined in notion.json yet
}

// OperationDefaults holds the default settings for various Notion operations.
type OperationDefaults struct {
	Search        SearchDefaults        `mapstructure:"search"`
	QueryDatabase QueryDatabaseDefaults `mapstructure:"query_database"`
	ListUsers     ListUsersDefaults     `mapstructure:"list_users"`
	ListDatabases ListDatabasesDefaults `mapstructure:"list_databases"`
	CreatePage    CreatePageDefaults    `mapstructure:"create_page"`
	AppendToBlock AppendToBlockDefaults `mapstructure:"append_to_block"`
}

// NotionAdapter is an adapter for the Notion API using OAuth2.
type NotionAdapter struct {
	*base.BaseAdapter
	oauthConfig   *domain.OAuthConfig
	notionVersion string // Notion API version (e.g., "2022-06-28")
	restAdapter   domain.Adapter
	operations    Operations
	permissionSet domain.PermissionSet
	defaults      *OperationDefaults
}

// NewNotionAdapter creates a new Notion adapter.
func NewNotionAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	oauthConfig *domain.OAuthConfig,
	notionVersion string,
	defaults *OperationDefaults,
	permissions domain.PermissionSet,
	restAdapter *rest.RESTAdapter,
) *NotionAdapter {
	baseAdapter := base.NewBaseAdapter(providerInfo, config)

	adapter := &NotionAdapter{
		BaseAdapter:   baseAdapter,
		oauthConfig:   oauthConfig,
		notionVersion: notionVersion,
		restAdapter:   restAdapter,
		defaults:      defaults,
		operations:    make(Operations),
		permissionSet: permissions,
	}

	// Operations are registered via the registerOperations function defined in the _operations.go file
	adapter.registerOperations()

	return adapter
}

// createOAuth2Config creates the oauth2.Config object.
func (a *NotionAdapter) createOAuth2Config(redirectURL string, scopes []string) *oauth2.Config {
	endpoint := oauth2.Endpoint{
		AuthURL:  authURL,
		TokenURL: tokenURL,
	}
	//

	return &oauth2.Config{
		ClientID:     a.oauthConfig.ClientID,
		ClientSecret: a.oauthConfig.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     endpoint,
	}
}

// ShouldRefreshToken checks if the token should be refreshed
func (a *NotionAdapter) ShouldRefreshToken(oldToken *oauth2.Token) bool {
	return false
}

// RefreshOAuthToken refreshes an OAuth token.
// Note: Some providers might not support refresh tokens or use different mechanisms.
func (a *NotionAdapter) RefreshOAuthToken(ctx context.Context, oldToken *oauth2.Token) (*oauth2.Token, error) {
	//
	// Provider does not support standard refresh tokens via this config.
	return nil, domain.NewAdapterError(
		a.GetProviderAdapterInfo().Identifier,
		"oauth_refresh",
		domain.ErrOperationNotSupported,
		"token refresh not supported by this provider configuration",
		http.StatusNotImplemented,
	)
	//
}

// GenerateOAuthURL generates an OAuth authorization URL.
func (a *NotionAdapter) GenerateOAuthURL(ctx context.Context,
	redirectURL, state, codeChallenge string,
	scopes []string,
) (string, error) {
	oauth2Config := a.createOAuth2Config(redirectURL, scopes)
	// Consider AccessTypeOffline if refresh tokens are needed and supported
	// TODO: Make AccessType configurable?
	return oauth2Config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		/*, oauth2.AccessTypeOffline*/), nil
}

// ExchangeCodeForTokens exchanges an authorization code for tokens.
func (a *NotionAdapter) ExchangeCodeForTokens(ctx context.Context, code, redirectURL, codeVerifier string) (*oauth2.Token, error) {
	oauth2Config := a.createOAuth2Config(redirectURL, nil) // Scopes not needed for exchange usually
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		// Wrap error into AdapterError
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			"oauth_exchange",          // Internal operation identifier for code exchange
			domain.ErrCredentialError, // Assuming exchange failure is a credential error
			fmt.Sprintf("failed to exchange code for token: %v", err),
			http.StatusUnauthorized, // Or another appropriate status
		)
	}
	return token, nil
}

// CheckMissingPermissions checks if required permissions are present in authorized scopes.
func (a *NotionAdapter) CheckMissingPermissions(operationIdentifier string, authorizedScopes []string) (bool, []string, error) {
	opDef, exists := a.operations[operationIdentifier]
	if !exists {
		return false, nil, fmt.Errorf("unknown operation ID: %s", operationIdentifier)
	}

	// Get the required permission identifiers for the operation
	requiredPermIdentifiers := opDef.PermissionIdentifiers // Assuming this field exists in OperationDefinition

	// Check missing permissions using the PermissionSet
	allScopesPresent, missingIdentifiers := a.permissionSet.CheckMissingPermissionsByIdentifiers(requiredPermIdentifiers, authorizedScopes)
	return allScopesPresent, missingIdentifiers, nil
}

func (a *NotionAdapter) GetScopesFromPermissions(permissions []string) ([]string, error) {
	return a.permissionSet.RequiredOAuthScopesByIdentifiers(permissions), nil
}

func (a *NotionAdapter) GetPermissionIdentifiersFromScopes(scopes []string) ([]string, error) {
	return []string{"access_notion"}, nil
}

// Execute handles API calls based on the operationID using the REST adapter.
func (a *NotionAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{}, // User-provided parameters
	credential interface{}, // Expected to be *credDomain.OAuthCredential
) (interface{}, error) {
	// 1. Validate Credential Type and Scopes
	oauthCred, ok := credential.(*credDomain.OAuthCredential)
	if !ok || oauthCred == nil || oauthCred.Token == nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrCredentialError, "invalid or missing OAuth credential", http.StatusUnauthorized)
	}

	// 2. Find Operation Definition
	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrOperationNotSupported, fmt.Sprintf("unknown operation ID: %s", operationID), http.StatusNotFound)
	}

	// 3. Check Permissions based on stored credential scopes
	authorizedScopes := oauthCred.Scopes
	allScopesPresent, missingIDs, err := a.CheckMissingPermissions(operationID, authorizedScopes)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInternal, fmt.Sprintf("error checking permissions: %v", err), http.StatusInternalServerError)
	}
	if !allScopesPresent {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrCredentialError, fmt.Sprintf("missing required permissions: %v", missingIDs), http.StatusForbidden)
	}

	// 4. Process User Parameters (Validation based on registered schema)
	processedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInvalidParameters, fmt.Sprintf("parameter validation failed: %v", err), http.StatusBadRequest)
	}

	// 5. Call the Operation Handler to get REST parameters
	handler := opDef.Handler
	restParams, err := handler(ctx, processedParams)
	if err != nil {
		// Wrap handler errors
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInternal, fmt.Sprintf("operation handler failed: %v", err), http.StatusInternalServerError)
	}

	// 6. Inject Authentication and Notion Headers
	headers, _ := restParams["headers"].(map[string]string)
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Authorization"] = utils.StringsBuilder("Bearer ", oauthCred.Token.AccessToken)
	headers["Notion-Version"] = a.notionVersion // Use stored version
	// Content-Type will be handled by restAdapter if body exists
	restParams["headers"] = headers

	// 7. Execute via REST Adapter
	// Pass nil credential to REST adapter as auth is handled here
	rawResult, err := a.restAdapter.Execute(ctx, operationID, restParams, nil)
	if err != nil {
		// Error should already be wrapped by restAdapter
		return nil, err
	}

	// The operationResponseTypes map is now primarily for documentation,
	// testing, or potential future use by callers, not for automatic decoding here.

	// Return the raw result directly from the restAdapter
	return rawResult, nil
}
