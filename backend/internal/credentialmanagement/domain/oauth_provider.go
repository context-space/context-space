package domain

import (
	"context"

	"golang.org/x/oauth2"
)

// OAuthProvider defines the interface for OAuth operations
// This interface acts as an anti-corruption layer between Credential Management and Provider Adapter
type OAuthProvider interface {
	// GenerateOAuthURL generates an OAuth authorization URL
	GenerateOAuthURL(ctx context.Context, providerIdentifier, redirectURL, state, codeChallenge string, scopes []string) (string, error)

	// ExchangeCodeForToken exchanges an authorization code for OAuth tokens
	// codeVerifier is required for OAuth PKCE (Proof Key for Code Exchange)
	ExchangeCodeForToken(ctx context.Context, providerIdentifier, code, redirectURL, codeVerifier string) (*oauth2.Token, error)

	// ShouldRefreshToken checks if a token should be refreshed
	ShouldRefreshToken(providerIdentifier string, oldToken *oauth2.Token) (bool, error)

	// RefreshToken refreshes an OAuth token
	RefreshToken(ctx context.Context, providerIdentifier string, oldToken *oauth2.Token) (*oauth2.Token, error)

	// GetScopesFromPermissions gets the scopes from the permissions
	GetScopesFromPermissions(ctx context.Context, providerIdentifier string, permissions []string) ([]string, error)

	// GetPermissionIdentifiersFromScopes gets the permission identifiers from the scopes
	GetPermissionIdentifiersFromScopes(ctx context.Context, providerIdentifier string, scopes []string) ([]string, error)
}
