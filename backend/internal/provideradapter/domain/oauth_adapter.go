package domain

import (
	"context"

	"golang.org/x/oauth2"
)

// OAuthAdapter extends the base Adapter with OAuth-specific functionality
type OAuthAdapter interface {
	Adapter

	// ShouldRefreshToken checks if the token should be refreshed
	ShouldRefreshToken(oldToken *oauth2.Token) bool

	// RefreshOAuthToken refreshes an OAuth token
	RefreshOAuthToken(ctx context.Context, oldToken *oauth2.Token) (*oauth2.Token, error)

	// GenerateOAuthURL generates an OAuth authorization URL
	GenerateOAuthURL(ctx context.Context, redirectURL, state, codeChallenge string, scopes []string) (string, error)

	// ExchangeCodeForTokens exchanges an authorization code for tokens
	// codeVerifier is required for OAuth PKCE (Proof Key for Code Exchange)
	ExchangeCodeForTokens(ctx context.Context, code, redirectURL, codeVerifier string) (*oauth2.Token, error)

	// CheckMissingPermissions checks if the required permissions are present in the authorized scopes
	CheckMissingPermissions(operationIdentifier string, authorizedScopes []string) (bool, []string, error)

	// GetScopesFromPermissions gets the scopes from the permissions
	GetScopesFromPermissions(permissions []string) ([]string, error)

	// GetPermissionIdentifiersFromScopes gets the permission identifiers from the scopes
	GetPermissionIdentifiersFromScopes(scopes []string) ([]string, error)
}
