package registry

import (
	"fmt"
	"sync"

	"github.com/context-space/context-space/backend/internal/provideradapter/application"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
)

// ProviderLoader focuses on pure provider loading logic (Infrastructure layer)
type ProviderLoader struct {
	adapterFactory  *application.AdapterFactory
	loadedProviders map[string]domain.ProviderAdapterInfo
	mu              sync.RWMutex
}

// NewProviderLoader creates a new provider loader
func NewProviderLoader(adapterFactory *application.AdapterFactory) *ProviderLoader {
	return &ProviderLoader{
		adapterFactory:  adapterFactory,
		loadedProviders: make(map[string]domain.ProviderAdapterInfo),
	}
}

// LoadProvider implements ProviderLoaderInterface.LoadProvider
func (l *ProviderLoader) LoadProvider(config *domain.ProviderAdapterConfig) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Get the adapterIdentifier (same as providerIdentifier)
	adapterIdentifier := config.Identifier

	// Get the template for this adapter type
	template, ok := GetAdapterTemplate(adapterIdentifier)
	if !ok {
		return fmt.Errorf("unknown adapter_id: %s", adapterIdentifier)
	}

	// Validate the configuration
	if err := template.ValidateConfig(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Create the adapter
	adapter, err := template.CreateAdapter(config)
	if err != nil {
		return fmt.Errorf("failed to create adapter: %w", err)
	}

	// Register the adapter in factory
	l.adapterFactory.RegisterAdapter(config.Identifier, adapter)

	// Store metadata
	adapterInfo := adapter.GetProviderAdapterInfo()
	l.loadedProviders[config.Identifier] = *adapterInfo

	return nil
}

// GetLoadedProviders implements ProviderLoaderInterface.GetLoadedProviders
func (l *ProviderLoader) GetLoadedProviders() []domain.ProviderAdapterInfo {
	l.mu.RLock()
	defer l.mu.RUnlock()

	providers := make([]domain.ProviderAdapterInfo, 0, len(l.loadedProviders))
	for _, provider := range l.loadedProviders {
		providers = append(providers, provider)
	}

	return providers
}

// UnloadProvider 卸载指定的 provider
func (l *ProviderLoader) UnloadProvider(identifier string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Remove from adapter factory
	l.adapterFactory.UnregisterAdapter(identifier)

	// Remove from local records
	delete(l.loadedProviders, identifier)

	return nil
}

// IsProviderLoaded checks if provider is already loaded
func (l *ProviderLoader) IsProviderLoaded(identifier string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	_, exists := l.loadedProviders[identifier]
	return exists
}

// GetLoadedProviderCount gets the count of loaded providers
func (l *ProviderLoader) GetLoadedProviderCount() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return len(l.loadedProviders)
}
