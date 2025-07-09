package adapter

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"

	provideradapterApp "github.com/context-space/context-space/backend/internal/provideradapter/application"
)

// OAuthProviderAdapter implements the OAuthProvider interface
// This adapter acts as an anti-corruption layer between Credential Management and Provider Adapter
type OAuthProviderAdapter struct {
	adapterFactory *provideradapterApp.AdapterFactory
}

// NewOAuthProviderAdapter creates a new OAuth provider adapter
func NewOAuthProviderAdapter(adapterFactory *provideradapterApp.AdapterFactory) *OAuthProviderAdapter {
	return &OAuthProviderAdapter{
		adapterFactory: adapterFactory,
	}
}

// GenerateOAuthURL generates an OAuth authorization URL
func (a *OAuthProviderAdapter) GenerateOAuthURL(ctx context.Context, providerIdentifier string, redirectURL, state, codeChallenge string, scopes []string) (string, error) {
	// Get the OAuth adapter for this provider
	adapter, err := a.adapterFactory.GetOAuthAdapter(providerIdentifier)
	if err != nil {
		return "", fmt.Errorf("failed to get OAuth adapter: %w", err)
	}

	// Use the adapter to generate the OAuth URL
	return adapter.GenerateOAuthURL(ctx, redirectURL, state, codeChallenge, scopes)
}

// ExchangeCodeForToken exchanges an authorization code for OAuth tokens
func (a *OAuthProviderAdapter) ExchangeCodeForToken(ctx context.Context, providerIdentifier, code, redirectURL, codeVerifier string) (*oauth2.Token, error) {
	// Get the OAuth adapter for this provider
	adapter, err := a.adapterFactory.GetOAuthAdapter(providerIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth adapter: %w", err)
	}

	// Exchange the code for an OAuth token
	return adapter.ExchangeCodeForTokens(ctx, code, redirectURL, codeVerifier)
}

// ShouldRefreshToken checks if the token should be refreshed
func (a *OAuthProviderAdapter) ShouldRefreshToken(providerIdentifier string, token *oauth2.Token) (bool, error) {
	// Get the OAuth adapter for this provider
	adapter, err := a.adapterFactory.GetOAuthAdapter(providerIdentifier)
	if err != nil {
		return false, fmt.Errorf("failed to get OAuth adapter: %w", err)
	}

	return adapter.ShouldRefreshToken(token), nil
}

// RefreshToken refreshes an OAuth token
func (a *OAuthProviderAdapter) RefreshToken(ctx context.Context, providerIdentifier string, oldToken *oauth2.Token) (*oauth2.Token, error) {
	// Get the OAuth adapter for this provider
	adapter, err := a.adapterFactory.GetOAuthAdapter(providerIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth adapter: %w", err)
	}

	newToken, err := adapter.RefreshOAuthToken(ctx, oldToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh OAuth token: %w", err)
	}

	return newToken, nil
}

// GetScopesFromPermissions gets the scopes from the permissions
func (a *OAuthProviderAdapter) GetScopesFromPermissions(ctx context.Context, providerIdentifier string, permissions []string) ([]string, error) {
	// Get the OAuth adapter for this provider
	adapter, err := a.adapterFactory.GetOAuthAdapter(providerIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth adapter: %w", err)
	}

	// Get the scopes from the permissions
	return adapter.GetScopesFromPermissions(permissions)
}

// GetPermissionIdentifiersFromScopes gets the permission identifiers from the scopes
func (a *OAuthProviderAdapter) GetPermissionIdentifiersFromScopes(ctx context.Context, providerIdentifier string, scopes []string) ([]string, error) {
	// Get the OAuth adapter for this provider
	adapter, err := a.adapterFactory.GetOAuthAdapter(providerIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth adapter: %w", err)
	}

	// Get the permissions from the scopes
	return adapter.GetPermissionIdentifiersFromScopes(scopes)
}
