package figma

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
	"github.com/context-space/context-space/backend/internal/shared/utils"
)

type GetImageRendersDefaults struct {
	Format string `mapstructure:"format"`
	Scale  int    `mapstructure:"scale"`
}
type GetTeamComponentSetsDefaults struct {
	PageSize int `mapstructure:"page_size"`
}
type GetTeamComponentsDefaults struct {
	PageSize int `mapstructure:"page_size"`
}
type GetTeamStylesDefaults struct {
	PageSize int `mapstructure:"page_size"`
}

// OperationDefaults holds the default settings for various Figma operations.
type OperationDefaults struct {
	GetImageRenders      GetImageRendersDefaults      `mapstructure:"get_image_renders"`
	GetTeamComponentSets GetTeamComponentSetsDefaults `mapstructure:"get_team_component_sets"`
	GetTeamComponents    GetTeamComponentsDefaults    `mapstructure:"get_team_components"`
	GetTeamStyles        GetTeamStylesDefaults        `mapstructure:"get_team_styles"`
}

// FigmaAdapter is an adapter for the Figma API using OAuth2.
type FigmaAdapter struct {
	*base.BaseAdapter
	oauthConfig   *domain.OAuthConfig
	restAdapter   domain.Adapter       // The underlying REST adapter instance
	defaults      *OperationDefaults   // Operation defaults struct defined above
	operations    Operations           // Map of operation ID to definition defined in _operations.go.tmpl
	permissionSet domain.PermissionSet // Permission set defined in providercore
}

// NewFigmaAdapter creates a new Figma adapter.
func NewFigmaAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	oauthConfig *domain.OAuthConfig,
	restAdapter domain.Adapter,
	defaults *OperationDefaults,
	permissions domain.PermissionSet,
) *FigmaAdapter {
	baseAdapter := base.NewBaseAdapter(providerInfo, config)

	adapter := &FigmaAdapter{
		BaseAdapter:   baseAdapter,
		oauthConfig:   oauthConfig,
		restAdapter:   restAdapter,
		defaults:      defaults,
		operations:    make(Operations),
		permissionSet: permissions,
	}

	adapter.registerOperations()

	return adapter
}

// createOAuth2Config creates the oauth2.Config object.
func (a *FigmaAdapter) createOAuth2Config(redirectURL string, scopes []string) *oauth2.Config {
	endpoint := oauth2.Endpoint{
		AuthURL:   authURL,
		TokenURL:  tokenURL,
		AuthStyle: oauth2.AuthStyleInHeader,
	}

	return &oauth2.Config{
		ClientID:     a.oauthConfig.ClientID,
		ClientSecret: a.oauthConfig.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     endpoint,
	}
}

// ShouldRefreshToken checks if the token should be refreshed
func (a *FigmaAdapter) ShouldRefreshToken(oldToken *oauth2.Token) bool {
	if oldToken.Expiry.IsZero() {
		return false
	}
	return oldToken.Expiry.Compare(time.Now().Add(30*time.Minute)) < 1
}

// RefreshOAuthToken refreshes an OAuth token.
func (a *FigmaAdapter) RefreshOAuthToken(ctx context.Context, oldToken *oauth2.Token) (*oauth2.Token, error) {
	oauth2Config := a.createOAuth2Config("", nil)
	tokenSource := oauth2Config.TokenSource(ctx, oldToken)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			"oauth_refresh",
			domain.ErrCredentialError,
			fmt.Sprintf("failed to refresh token: %v", err),
			http.StatusInternalServerError,
		)
	}
	return newToken, nil
}

// GenerateOAuthURL generates an OAuth authorization URL.
func (a *FigmaAdapter) GenerateOAuthURL(ctx context.Context,
	redirectURL, state, codeChallenge string,
	scopes []string,
) (string, error) {
	oauth2Config := a.createOAuth2Config(redirectURL, scopes)
	return oauth2Config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		/*, oauth2.AccessTypeOffline*/), nil
}

// ExchangeCodeForTokens exchanges an authorization code for tokens.
func (a *FigmaAdapter) ExchangeCodeForTokens(ctx context.Context, code, redirectURL, codeVerifier string) (*oauth2.Token, error) {
	oauth2Config := a.createOAuth2Config(redirectURL, nil) // Scopes not needed for exchange usually
	token, err := oauth2Config.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			"oauth_exchange",          // Internal operation identifier
			domain.ErrCredentialError, // Assume exchange failure is credential-related
			fmt.Sprintf("failed to exchange code for token: %v", err),
			http.StatusUnauthorized, // Or another appropriate status code based on the error
		)
	}
	return token, nil
}

// CheckMissingPermissions checks if required permissions are present in authorized scopes.
func (a *FigmaAdapter) CheckMissingPermissions(operationIdentifier string, authorizedScopes []string) (bool, []string, error) {
	opDef, exists := a.operations[operationIdentifier]
	if !exists {
		return false, nil, fmt.Errorf("unknown operation ID: %s", operationIdentifier)
	}
	requiredPermIdentifiers := opDef.PermissionIdentifiers
	allScopesPresent, missingIdentifiers := a.permissionSet.CheckMissingPermissionsByIdentifiers(requiredPermIdentifiers, authorizedScopes)
	return allScopesPresent, missingIdentifiers, nil
}

// GetScopesFromPermissions translates internal permission identifiers to required OAuth scopes.
func (a *FigmaAdapter) GetScopesFromPermissions(permissions []string) ([]string, error) {
	return a.permissionSet.RequiredOAuthScopesByIdentifiers(permissions), nil
}

// GetPermissionIdentifiersFromScopes gets the permission identifiers from the scopes
func (a *FigmaAdapter) GetPermissionIdentifiersFromScopes(scopes []string) ([]string, error) {
	return a.permissionSet.GetPermissionIdentifiersFromScopes(scopes), nil
}

// Execute handles API calls based on the operationID using the REST adapter.
func (a *FigmaAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{}, // User-provided parameters
	credential interface{}, // Expected to be *credDomain.OAuthCredential
) (interface{}, error) {
	oauthCred, ok := credential.(*credDomain.OAuthCredential)
	if !ok || oauthCred == nil || oauthCred.Token == nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrCredentialError, "invalid or missing OAuth credential", http.StatusUnauthorized)
	}

	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrOperationNotSupported, fmt.Sprintf("unknown operation ID: %s", operationID), http.StatusNotFound)
	}

	authorizedScopes := oauthCred.Scopes

	allScopesPresent, missingIDs, err := a.CheckMissingPermissions(operationID, authorizedScopes)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInternal, fmt.Sprintf("error checking permissions: %v", err), http.StatusInternalServerError)
	}
	if !allScopesPresent {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrCredentialError, fmt.Sprintf("missing required permissions: %v", missingIDs), http.StatusForbidden)
	}

	processedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInvalidParameters, fmt.Sprintf("parameter validation failed: %v", err), http.StatusBadRequest)
	}

	handler := opDef.Handler
	restParams, err := handler(ctx, processedParams)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInternal, fmt.Sprintf("operation handler failed: %v", err), http.StatusInternalServerError)
	}

	headers, _ := restParams["headers"].(map[string]string)
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Authorization"] = utils.StringsBuilder("Bearer ", oauthCred.Token.AccessToken)
	restParams["headers"] = headers

	rawResult, err := a.restAdapter.Execute(ctx, operationID, restParams, nil)
	if err != nil {
		return nil, err
	}

	return rawResult, nil
}
