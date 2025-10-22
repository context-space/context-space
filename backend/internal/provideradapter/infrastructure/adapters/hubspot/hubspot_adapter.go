package hubspot

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

// OperationDefaults holds the default settings for various HubSpot operations.
type OperationDefaults struct {
}

// HubspotAdapter is an adapter for the HubSpot API using OAuth2.
type HubspotAdapter struct {
	*base.BaseAdapter
	oauthConfig   *domain.OAuthConfig
	restAdapter   domain.Adapter
	operations    Operations
	permissionSet domain.PermissionSet
}

// NewHubspotAdapter creates a new HubSpot adapter.
func NewHubspotAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	oauthConfig *domain.OAuthConfig,
	permissions domain.PermissionSet,
	restAdapter domain.Adapter,
) *HubspotAdapter {
	baseAdapter := base.NewBaseAdapter(providerInfo, config)

	adapter := &HubspotAdapter{
		BaseAdapter:   baseAdapter,
		oauthConfig:   oauthConfig,
		restAdapter:   restAdapter,
		operations:    make(Operations),
		permissionSet: permissions,
	}

	adapter.registerOperations()

	return adapter
}

// createOAuth2Config creates the oauth2.Config object.
func (a *HubspotAdapter) createOAuth2Config(redirectURL string, scopes []string) *oauth2.Config {
	endpoint := oauth2.Endpoint{
		AuthURL:  authURL,
		TokenURL: tokenURL,
		// Add AuthStyle if needed (e.g., oauth2.AuthStyleInParams)
	}

	effectiveRedirectURL := redirectURL

	return &oauth2.Config{
		ClientID:     a.oauthConfig.ClientID,
		ClientSecret: a.oauthConfig.ClientSecret,
		RedirectURL:  effectiveRedirectURL,
		Scopes:       scopes,
		Endpoint:     endpoint,
	}
}

// ShouldRefreshToken checks if the token should be refreshed
func (a *HubspotAdapter) ShouldRefreshToken(oldToken *oauth2.Token) bool {
	if oldToken.Expiry.IsZero() {
		return false
	}
	return oldToken.Expiry.Compare(time.Now().Add(30*time.Minute)) < 1
}

// RefreshOAuthToken refreshes an OAuth token.
func (a *HubspotAdapter) RefreshOAuthToken(ctx context.Context, oldToken *oauth2.Token) (*oauth2.Token, error) {
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
func (a *HubspotAdapter) GenerateOAuthURL(ctx context.Context,
	redirectURL, state, codeChallenge string,
	scopes []string,
) (string, error) {
	oauth2Config := a.createOAuth2Config(redirectURL, scopes)
	return oauth2Config.AuthCodeURL(state /*, oauth2.AccessTypeOffline*/), nil
}

// ExchangeCodeForTokens exchanges an authorization code for tokens.
func (a *HubspotAdapter) ExchangeCodeForTokens(ctx context.Context, code, redirectURL, codeVerifier string) (*oauth2.Token, error) {
	oauth2Config := a.createOAuth2Config(redirectURL, nil)
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			"oauth_exchange",
			domain.ErrCredentialError,
			fmt.Sprintf("failed to exchange code for token: %v", err),
			http.StatusUnauthorized,
		)
	}
	return token, nil
}

// CheckMissingPermissions checks if required permissions are present in authorized scopes.
func (a *HubspotAdapter) CheckMissingPermissions(operationIdentifier string, authorizedScopes []string) (bool, []string, error) {
	opDef, exists := a.operations[operationIdentifier]
	if !exists {
		return false, nil, fmt.Errorf("unknown operation ID: %s", operationIdentifier)
	}
	requiredPermIdentifiers := opDef.PermissionIdentifiers
	allScopesPresent, missingIdentifiers := a.permissionSet.CheckMissingPermissionsByIdentifiers(requiredPermIdentifiers, authorizedScopes)
	return allScopesPresent, missingIdentifiers, nil
}

// GetScopesFromPermissions translates internal permission identifiers to required OAuth scopes.
func (a *HubspotAdapter) GetScopesFromPermissions(permissions []string) ([]string, error) {
	return a.permissionSet.RequiredOAuthScopesByIdentifiers(permissions), nil
}

// GetPermissionIdentifiersFromScopes translates OAuth scopes to internal permission identifiers.
func (a *HubspotAdapter) GetPermissionIdentifiersFromScopes(scopes []string) ([]string, error) {
	return a.permissionSet.GetPermissionIdentifiersFromScopes(scopes), nil
}

// Execute handles API calls based on the operationID using the REST adapter.
func (a *HubspotAdapter) Execute(
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
