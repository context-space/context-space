package domain

import "context"

// AdapterRepository defines the interface for provider adapter persistence
type ProviderAdapterConfigRepository interface {
	// GetByID returns a provider by ID
	GetByID(ctx context.Context, id string) (*ProviderAdapterConfig, error)

	// GetByIdentifier returns a provider by identifier
	GetByIdentifier(ctx context.Context, identifier string) (*ProviderAdapterConfig, error)

	// GetByIdentifierWithoutCache returns a provider by identifier without cache
	GetByIdentifierWithoutCache(ctx context.Context, identifier string) (*ProviderAdapterConfig, error)

	// List returns all providers
	ListAdapterConfigs(ctx context.Context) ([]*ProviderAdapterConfig, error)

	// Create creates a new provider
	Create(ctx context.Context, provider *ProviderAdapterConfig) error

	// Update updates a provider
	Update(ctx context.Context, provider *ProviderAdapterConfig) error

	// Delete deletes a provider
	Delete(ctx context.Context, id string) error
}
