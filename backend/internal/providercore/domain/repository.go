package domain

import (
	"context"
)

// ProviderRepository defines the interface for provider persistence
type ProviderRepository interface {
	// GetByID returns a provider by ID
	GetByID(ctx context.Context, id string) (*Provider, error)

	// GetByIdentifier returns a provider by identifier
	GetByIdentifier(ctx context.Context, identifier string) (*Provider, error)

	// List returns all providers
	List(ctx context.Context) ([]*Provider, error)

	// Create creates a new provider
	Create(ctx context.Context, provider *Provider) error

	// Update updates a provider
	Update(ctx context.Context, provider *Provider) error

	// Delete deletes a provider
	Delete(ctx context.Context, id string) error

	// SyncTagsToProvider syncs tags to provider's json_attributes
	SyncTagsToProvider(ctx context.Context, providerIdentifier string, tags []string) error
}

// OperationRepository defines the interface for operation persistence
type OperationRepository interface {
	// ListByProviderID returns all operations for a provider
	ListByProviderID(ctx context.Context, providerID string) ([]*Operation, error)

	// GetByID returns an operation by ID
	GetByID(ctx context.Context, id string) (*Operation, error)

	// GetByProviderIDAndIdentifier returns an operation by provider ID and identifier
	GetByProviderIDAndIdentifier(ctx context.Context, providerID, identifier string) (*Operation, error)

	// Create creates a new operation
	Create(ctx context.Context, operation *Operation) error

	// Update updates an operation
	Update(ctx context.Context, operation *Operation) error

	// Delete deletes an operation
	Delete(ctx context.Context, id string) error
}
