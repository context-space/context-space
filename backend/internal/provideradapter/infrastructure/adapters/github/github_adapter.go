package github

import (
	"context"
	"fmt"
	"time"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
)

// GitHubAdapter is an adapter for GitHub
type GitHubAdapter struct {
	*base.BaseAdapter
	oauthConfig   *domain.OAuthConfig
	operations    Operations
	permissionSet domain.PermissionSet
}

// NewGitHubAdapter creates a new GitHub adapter
func NewGitHubAdapter(
	providerAdapterInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	oauthConfig *domain.OAuthConfig,
	permissions domain.PermissionSet,
) *GitHubAdapter {
	baseAdapter := base.NewBaseAdapter(providerAdapterInfo, config)

	adapter := &GitHubAdapter{
		BaseAdapter:   baseAdapter,
		oauthConfig:   oauthConfig,
		operations:    make(Operations),
		permissionSet: permissions,
	}

	adapter.registerOperations()

	return adapter
}

// ShouldRefreshToken checks if the token should be refreshed
func (a *GitHubAdapter) ShouldRefreshToken(oldToken *oauth2.Token) bool {
	if oldToken.Expiry.IsZero() {
		return false
	}
	return oldToken.Expiry.Compare(time.Now().Add(30*time.Minute)) < 1
}

// RefreshOAuthToken refreshes an OAuth token
func (a *GitHubAdapter) RefreshOAuthToken(
	ctx context.Context,
	oldToken *oauth2.Token,
) (*oauth2.Token, error) {
	// GitHub doesn't support refresh tokens for OAuth apps, the following is a sample implementation for other third-party providers
	newToken, err := a.createOAuth2Config("", nil).TokenSource(ctx, oldToken).Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	return newToken, nil
}

// GenerateOAuthURL generates an OAuth authorization URL
func (a *GitHubAdapter) GenerateOAuthURL(
	ctx context.Context,
	redirectURL, state, codeChallenge string,
	scopes []string,
) (string, error) {
	oauth2Config := a.createOAuth2Config(redirectURL, scopes)

	// Generate the authorization URL
	return oauth2Config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		/*, oauth2.AccessTypeOffline*/), nil
}

// ExchangeCodeForTokens exchanges an authorization code for tokens
func (a *GitHubAdapter) ExchangeCodeForTokens(ctx context.Context, code, redirectURL, codeVerifier string) (*oauth2.Token, error) {
	oauth2Config := a.createOAuth2Config(redirectURL, nil)

	opts := []oauth2.AuthCodeOption{}
	if codeVerifier != "" {
		opts = append(opts, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	}

	// Exchange the code for a token with PKCE code verifier
	token, err := oauth2Config.Exchange(ctx, code, opts...)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// Execute handles various GitHub API operations based on the operationID
func (a *GitHubAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{},
	credential interface{},
) (interface{}, error) {
	// Validate the credential scopes
	if err := a.ValidateOperationOAuthScopes(ctx, operationID, credential); err != nil {
		return nil, err
	}

	// Process and validate parameters
	parsedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		return nil, err
	}

	// Get the operation handler
	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, fmt.Errorf("unknown operation ID: %s", operationID)
	}

	// Create GitHub client
	client, err := a.createGitHubClient(credential)
	if err != nil {
		return nil, err
	}

	// Execute the operation
	return opDef.Handler(ctx, client, parsedParams)
}

// CheckMissingPermissions checks if the required permissions are present in the authorized scopes
func (a *GitHubAdapter) CheckMissingPermissions(operationIdentifier string, authorizedScopes []string) (bool, []string, error) {
	opDef, exists := a.operations[operationIdentifier]
	if !exists {
		return false, nil, fmt.Errorf("unknown operation ID: %s", operationIdentifier)
	}

	allScopesPresent, missingPermissionsIdentifiers := a.permissionSet.CheckMissingPermissionsByIdentifiers(opDef.PermissionIdentifiers, authorizedScopes)
	return allScopesPresent, missingPermissionsIdentifiers, nil
}

// GetScopesFromPermissions gets the scopes from the permissions
func (a *GitHubAdapter) GetScopesFromPermissions(permissions []string) ([]string, error) {
	return a.permissionSet.RequiredOAuthScopesByIdentifiers(permissions), nil
}

// GetPermissionIdentifiersFromScopes gets the permission identifiers from the scopes
func (a *GitHubAdapter) GetPermissionIdentifiersFromScopes(scopes []string) ([]string, error) {
	return a.permissionSet.GetPermissionIdentifiersFromScopes(scopes), nil
}

// ValidateOperationOAuthScopes validates the scopes of an OAuth token for an operation
func (a *GitHubAdapter) ValidateOperationOAuthScopes(ctx context.Context, operationID string, credential interface{}) error {
	oauthCredential, ok := credential.(*credDomain.OAuthCredential)
	if !ok {
		return fmt.Errorf("credential is not an OAuthCredential")
	}

	allScopesPresent, missingPermissionsIdentifiers, err := a.CheckMissingPermissions(operationID, oauthCredential.Scopes)
	if err != nil {
		return err
	}

	if !allScopesPresent {
		return fmt.Errorf("missing required scopes: %v", missingPermissionsIdentifiers)
	}

	return nil
}

func (a *GitHubAdapter) createOAuth2Config(redirectURL string, scopes []string) *oauth2.Config {
	oauth2Config := &oauth2.Config{
		ClientID:     a.oauthConfig.ClientID,
		ClientSecret: a.oauthConfig.ClientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     githubOAuth.Endpoint,
		Scopes:       scopes,
	}
	if redirectURL != "" {
		oauth2Config.RedirectURL = redirectURL
	}
	return oauth2Config
}
