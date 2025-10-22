package domain

import (
	"context"

	contractCredential "github.com/context-space/context-space/backend/internal/shared/contract/credentialmanagement"
	contractAdapter "github.com/context-space/context-space/backend/internal/shared/contract/provideradapter"
	contractProvider "github.com/context-space/context-space/backend/internal/shared/contract/providercore"
)

type ProviderProvider interface {
	// GetProviderByIdentifier returns a provider by Identifier
	GetProviderByIdentifier(ctx context.Context, providerIdentifier string) (*contractProvider.ProviderDTO, error)
}

// AdapterProvider defines an interface for getting provider adapters
type AdapterProvider interface {
	// GetAdapter returns an adapter for the given provider ID
	GetAdapterByProviderIdentifier(ctx context.Context, providerIdentifier string) (contractAdapter.AdapterContract, error)
}

// CredentialProvider defines an interface for getting provider credentials
type CredentialProvider interface {
	// GetCredentialByUserAndProvider retrieves a credential by user ID and provider ID
	GetCredentialByUserAndProvider(ctx context.Context, userID, providerIdentifier string) (interface{}, error)

	// CreateNone creates a new no-auth credential
	CreateNone(ctx context.Context, userID, providerIdentifier string) (*contractCredential.CredentialDTO, error)

	// UpdateCredentialLastUsedAt updates the last used at time of a credential
	UpdateCredentialLastUsedAt(ctx context.Context, credential interface{}) error
}

type TokenRefreshProvider interface {
	// RefreshAccessToken refreshes the access token if needed
	RefreshAccessToken(ctx context.Context, providerIdentifier string, credential interface{}) (interface{}, error)
}
