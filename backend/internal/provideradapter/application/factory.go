package application

import (
	"fmt"
	"sync"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
)

// AdapterFactory creates and manages adapters for providers
type AdapterFactory struct {
	adapters          map[string]domain.Adapter
	oauthAdapters     map[string]domain.OAuthAdapter
	apiKeyAdapters    map[string]domain.APIKeyAdapter
	basicAuthAdapters map[string]domain.BasicAuthAdapter
	publicAdapters    map[string]domain.PublicAdapter
	mutex             sync.RWMutex
}

// NewAdapterFactory creates a new adapter factory
func NewAdapterFactory() *AdapterFactory {
	return &AdapterFactory{
		adapters:          make(map[string]domain.Adapter),
		oauthAdapters:     make(map[string]domain.OAuthAdapter),
		apiKeyAdapters:    make(map[string]domain.APIKeyAdapter),
		basicAuthAdapters: make(map[string]domain.BasicAuthAdapter),
		publicAdapters:    make(map[string]domain.PublicAdapter),
	}
}

// RegisterAdapter registers an adapter for a provider
func (f *AdapterFactory) RegisterAdapter(providerIdentifier string, adapter domain.Adapter) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.adapters[providerIdentifier] = adapter

	// Register in specific adapter maps if applicable
	if oauthAdapter, ok := adapter.(domain.OAuthAdapter); ok {
		f.oauthAdapters[providerIdentifier] = oauthAdapter
	}

	if apiKeyAdapter, ok := adapter.(domain.APIKeyAdapter); ok {
		f.apiKeyAdapters[providerIdentifier] = apiKeyAdapter
	}

	if basicAuthAdapter, ok := adapter.(domain.BasicAuthAdapter); ok {
		f.basicAuthAdapters[providerIdentifier] = basicAuthAdapter
	}

	if publicAdapter, ok := adapter.(domain.PublicAdapter); ok {
		f.publicAdapters[providerIdentifier] = publicAdapter
	}
}

// GetAdapter returns an adapter for a provider
func (f *AdapterFactory) GetAdapter(providerIdentifier string) (domain.Adapter, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	adapter, exists := f.adapters[providerIdentifier]
	if !exists {
		return nil, fmt.Errorf("no adapter registered for provider: %s", providerIdentifier)
	}

	return adapter, nil
}

// GetOAuthAdapter returns an OAuth adapter for a provider
func (f *AdapterFactory) GetOAuthAdapter(providerIdentifier string) (domain.OAuthAdapter, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	adapter, exists := f.oauthAdapters[providerIdentifier]
	if !exists {
		return nil, fmt.Errorf("no OAuth adapter registered for provider: %s", providerIdentifier)
	}

	return adapter, nil
}

// GetAPIKeyAdapter returns an API key adapter for a provider
func (f *AdapterFactory) GetAPIKeyAdapter(providerIdentifier string) (domain.APIKeyAdapter, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	adapter, exists := f.apiKeyAdapters[providerIdentifier]
	if !exists {
		return nil, fmt.Errorf("no API key adapter registered for provider: %s", providerIdentifier)
	}

	return adapter, nil
}

// GetBasicAuthAdapter returns a basic auth adapter for a provider
func (f *AdapterFactory) GetBasicAuthAdapter(providerIdentifier string) (domain.BasicAuthAdapter, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	adapter, exists := f.basicAuthAdapters[providerIdentifier]
	if !exists {
		return nil, fmt.Errorf("no basic auth adapter registered for provider: %s", providerIdentifier)
	}

	return adapter, nil
}

// GetPublicAdapter returns a public adapter for a provider
func (f *AdapterFactory) GetPublicAdapter(providerIdentifier string) (domain.PublicAdapter, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	adapter, exists := f.publicAdapters[providerIdentifier]
	if !exists {
		return nil, fmt.Errorf("no public adapter registered for provider: %s", providerIdentifier)
	}

	return adapter, nil
}

// UnregisterAdapter removes an adapter for a provider
func (f *AdapterFactory) UnregisterAdapter(providerIdentifier string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// Remove from main adapters map
	delete(f.adapters, providerIdentifier)

	// Remove from specific adapter maps
	delete(f.oauthAdapters, providerIdentifier)
	delete(f.apiKeyAdapters, providerIdentifier)
	delete(f.basicAuthAdapters, providerIdentifier)
	delete(f.publicAdapters, providerIdentifier)
}

// ListRegisteredProviders returns all registered provider identifiers
func (f *AdapterFactory) ListRegisteredProviders() []string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	providers := make([]string, 0, len(f.adapters))
	for identifier := range f.adapters {
		providers = append(providers, identifier)
	}

	return providers
}

// IsAdapterRegistered checks if an adapter is registered for a provider
func (f *AdapterFactory) IsAdapterRegistered(providerIdentifier string) bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	_, exists := f.adapters[providerIdentifier]
	return exists
}
