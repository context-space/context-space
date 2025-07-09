package domain

import (
	"context"

	providercore "github.com/context-space/context-space/backend/internal/providercore/domain"
)

// ProviderCoreDataProvider defines the interface for accessing provider core data
// This follows DDD anti-corruption layer pattern
type ProviderCoreDataProvider interface {
	// Provider metadata access
	GetProviderCoreData(ctx context.Context, identifier string) (*providercore.Provider, error)
	ListProviders(ctx context.Context) ([]*providercore.Provider, error)
}
