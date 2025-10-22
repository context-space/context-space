package provideradapter

import (
	"context"

	"golang.org/x/oauth2"
)

// AdapterDTO defines the contract interface for provider adapters
// This is used for cross-module communication through the contract layer
type AdapterContract interface {
	// Execute an operation call to the provider
	ExecuteContract(
		ctx context.Context,
		operationID string,
		params map[string]interface{},
		credential interface{},
	) (interface{}, error)

	// GetAdapterInfo returns information about this adapter
	GetAdapterInfoContract() *AdapterInfoDTO
}

// ProviderAdapterContract defines the contract interface for provider adapter operations
// This provides a stable interface for cross-module communication
type ProviderAdapterContract interface {
	// GetAdapter returns an adapter for the given provider ID
	GetAdapterContract(ctx context.Context, providerIdentifier string) (AdapterContract, error)

	// GenerateOAuthURLContract generates an OAuth authorization URL
	GenerateOAuthURLContract(ctx context.Context, providerIdentifier string, redirectURL, state, codeChallenge string, scopes []string) (string, error)

	// ExchangeCodeForTokenContract exchanges an authorization code for OAuth tokens
	ExchangeCodeForTokenContract(ctx context.Context, providerIdentifier, code, redirectURL, codeVerifier string) (*oauth2.Token, error)

	// ShouldRefreshTokenContract checks if a token should be refreshed
	ShouldRefreshTokenContract(providerIdentifier string, token *oauth2.Token) (bool, error)

	// RefreshTokenContract refreshes an OAuth token
	RefreshTokenContract(ctx context.Context, providerIdentifier string, oldToken *oauth2.Token) (*oauth2.Token, error)

	// GetScopesFromPermissionsContract gets the scopes from the permissions
	GetScopesFromPermissionsContract(ctx context.Context, providerIdentifier string, permissions []string) ([]string, error)

	// GetPermissionIdentifiersFromScopesContract gets the permission identifiers from the scopes
	GetPermissionIdentifiersFromScopesContract(ctx context.Context, providerIdentifier string, scopes []string) ([]string, error)
}
