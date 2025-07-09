package domain

import (
	"context"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	providerAdapterDomain "github.com/context-space/context-space/backend/internal/provideradapter/domain"
	providerDomain "github.com/context-space/context-space/backend/internal/providercore/domain"
)

type ProviderProvider interface {
	// GetProviderByIdentifier returns a provider by Identifier
	GetProviderByIdentifier(ctx context.Context, providerIdentifier string) (*providerDomain.Provider, error)
}

// AdapterProvider defines an interface for getting provider adapters
type AdapterProvider interface {
	// GetAdapter returns an adapter for the given provider ID
	GetAdapter(ctx context.Context, providerIdentifier string) (providerAdapterDomain.Adapter, error)
}

// CredentialProvider defines an interface for getting provider credentials
type CredentialProvider interface {
	// GetCredentialByUserAndProvider retrieves a credential by user ID and provider ID
	GetCredentialByUserAndProvider(ctx context.Context, userID, providerIdentifier string) (interface{}, error)

	// CreateNone creates a new no-auth credential
	CreateNone(ctx context.Context, userID, providerIdentifier string) (*credDomain.NoneCredential, error)

	// UpdateCredentialLastUsedAt updates the last used at time of a credential
	UpdateCredentialLastUsedAt(ctx context.Context, credential interface{}) error
}
