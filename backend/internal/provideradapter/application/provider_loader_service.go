package application

import (
	"context"
	"fmt"
	"sync"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
)

type ProviderLoadMetadata struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	AuthType string `json:"auth_type"`
	Loaded   bool   `json:"loaded"`
	Error    string `json:"error,omitempty"`
}

// ProviderLoaderInterface define the interface for provider loader
type ProviderLoaderInterface interface {
	// LoadProvider load a single provider
	LoadProvider(config *domain.ProviderAdapterConfig) error
	// GetLoadedProviders get all loaded providers
	GetLoadedProviders() []domain.ProviderAdapterInfo
}

// ProviderLoaderService coordinate provider loading operations across domains
type ProviderLoaderService struct {
	coreDataProvider     domain.ProviderCoreDataProvider
	configRepository     domain.ProviderAdapterConfigRepository
	loader               ProviderLoaderInterface
	providersLoadRecords map[string]*ProviderLoadMetadata
	mu                   sync.RWMutex
}

// NewProviderLoaderService create a new provider loader service
func NewProviderLoaderService(
	coreDataProvider domain.ProviderCoreDataProvider,
	configRepository domain.ProviderAdapterConfigRepository,
	loader ProviderLoaderInterface,
) *ProviderLoaderService {
	return &ProviderLoaderService{
		coreDataProvider:     coreDataProvider,
		configRepository:     configRepository,
		loader:               loader,
		providersLoadRecords: make(map[string]*ProviderLoadMetadata),
	}
}

// LoadAllProviders load all available providers
func (s *ProviderLoaderService) LoadAllProviders(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. get all provider infos from ProviderCore
	providerInfos, err := s.coreDataProvider.ListProviders(ctx)
	if err != nil {
		return fmt.Errorf("failed to get provider infos from core: %w", err)
	}

	// 2. get all adapter configs from ProviderAdapter
	configs, err := s.configRepository.ListAdapterConfigs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get adapter configs: %w", err)
	}

	// 3. create config map for quick lookup
	configMap := make(map[string]domain.ProviderAdapterConfig)
	for _, config := range configs {
		configMap[config.Identifier] = *config
	}

	// 4. load each provider
	for _, providerInfo := range providerInfos {
		// Record the provider load status
		s.providersLoadRecords[providerInfo.Identifier] = &ProviderLoadMetadata{
			ID:       providerInfo.Identifier,
			Name:     providerInfo.Name,
			AuthType: string(providerInfo.AuthType),
			Loaded:   true,
			Error:    "",
		}
		config, exists := configMap[providerInfo.Identifier]
		if !exists {
			s.providersLoadRecords[providerInfo.Identifier].Loaded = false
			s.providersLoadRecords[providerInfo.Identifier].Error = fmt.Sprintf("no adapter config found for provider %s", providerInfo.Identifier)
			continue
		}

		config.Identifier = providerInfo.Identifier
		config.Name = providerInfo.Name
		config.Status = providerInfo.Status
		config.Description = providerInfo.Description
		config.AuthType = providerInfo.AuthType
		config.Permissions = providerInfo.Permissions
		config.Operations = providerInfo.Operations
		err = s.loader.LoadProvider(&config)
		if err != nil {
			s.providersLoadRecords[providerInfo.Identifier].Loaded = false
			s.providersLoadRecords[providerInfo.Identifier].Error = err.Error()
		}
	}

	return nil
}

// GetLoadedProviders get all loaded providers
func (s *ProviderLoaderService) GetProvidersLoadRecords() []ProviderLoadMetadata {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]ProviderLoadMetadata, 0, len(s.providersLoadRecords))

	for _, metadata := range s.providersLoadRecords {
		result = append(result, *metadata)
	}
	return result
}

// ReloadProvider reload a specific provider
func (s *ProviderLoaderService) ReloadProvider(ctx context.Context, identifier string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// clear previous loading status
	delete(s.providersLoadRecords, identifier)

	// 1. get provider info from ProviderCore
	providerInfo, err := s.coreDataProvider.GetProviderCoreData(ctx, identifier)
	if err != nil {
		return fmt.Errorf("failed to get provider info for %s: %w", identifier, err)
	}

	// 2. get adapter config from ProviderAdapter
	config, err := s.configRepository.GetByIdentifier(ctx, identifier)
	if err != nil {
		return fmt.Errorf("failed to get adapter config for %s: %w", identifier, err)
	}
	config.Identifier = providerInfo.Identifier
	config.Name = providerInfo.Name
	config.Status = providerInfo.Status
	config.Description = providerInfo.Description
	config.AuthType = providerInfo.AuthType
	config.Permissions = providerInfo.Permissions
	config.Operations = providerInfo.Operations

	// 3. execute loading
	return s.loader.LoadProvider(config)
}
