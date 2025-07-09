package adapter

import (
	"context"
	"fmt"

	"github.com/context-space/context-space/backend/internal/integration/domain"
	providerAdapterApp "github.com/context-space/context-space/backend/internal/provideradapter/application"
	providerAdapterDomain "github.com/context-space/context-space/backend/internal/provideradapter/domain"
)

// ProviderAdapterImpl is an implementation of the domain.AdapterProvider interface
// that connects to the Provider Adapter Context
type ProviderAdapterImpl struct {
	adapterFactory *providerAdapterApp.AdapterFactory
}

// NewProviderAdapter creates a new provider adapter implementation
func NewProviderAdapter(adapterFactory *providerAdapterApp.AdapterFactory) domain.AdapterProvider {
	return &ProviderAdapterImpl{
		adapterFactory: adapterFactory,
	}
}

// GetAdapter returns an adapter for the given provider identifier
func (p *ProviderAdapterImpl) GetAdapter(ctx context.Context, providerIdentifier string) (providerAdapterDomain.Adapter, error) {
	if p.adapterFactory == nil {
		return nil, fmt.Errorf("adapter factory is not initialized")
	}

	return p.adapterFactory.GetAdapter(providerIdentifier)
}
