package domain

import (
	"context"
)

// InvocationRepository defines the interface for invocation persistence
type InvocationRepository interface {
	// Create creates a new invocation
	Create(ctx context.Context, invocation *Invocation) error

	// Update updates an invocation
	Update(ctx context.Context, invocation *Invocation) error

	// GetByID returns an invocation by ID
	GetByID(ctx context.Context, id string) (*Invocation, error)

	// ListByUserID returns invocations by user ID
	ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*Invocation, error)

	// CountByUserID returns the count of invocations by user ID
	CountByUserID(ctx context.Context, userID string) (int64, error)

	// CountByProviderIdentifier returns the count of invocations by provider identifier
	CountByProviderIdentifier(ctx context.Context, providerIdentifier string) (int64, error)

	// CountByOperationIdentifier returns the count of invocations by operation identifier
	CountByOperationIdentifier(ctx context.Context, providerIdentifier, operationIdentifier string) (int64, error)
}
