package application

import (
	"context"
	"fmt"
	"sync"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"go.uber.org/zap"
)

type ProviderLoadMetadata struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	AuthType string `json:"auth_type"`
	Loaded   bool   `json:"loaded"`
	Error    string `json:"error,omitempty"`
}

// ProviderLoaderService coordinate provider loading operations across domains
type ProviderLoaderService struct {
	coreDataProvider domain.ProviderCoreDataProvider
	configRepository domain.ProviderAdapterConfigRepository
	loader           domain.ProviderAdapterLoader
	obs              *observability.ObservabilityProvider
	mu               sync.RWMutex
}

// NewProviderLoaderService create a new provider loader service
func NewProviderLoaderService(
	coreDataProvider domain.ProviderCoreDataProvider,
	configRepository domain.ProviderAdapterConfigRepository,
	loader domain.ProviderAdapterLoader,
	obs *observability.ObservabilityProvider,
) *ProviderLoaderService {
	return &ProviderLoaderService{
		coreDataProvider: coreDataProvider,
		configRepository: configRepository,
		loader:           loader,
		obs:              obs,
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
		config, exists := configMap[providerInfo.Identifier]
		if !exists {
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
			s.obs.Logger.Error(ctx, "failed to load provider", zap.Error(err))
			continue
		}
	}

	return nil
}

func (s *ProviderLoaderService) GetLoadedProviders() []domain.ProviderAdapterInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.loader.GetLoadedProviders()
}

func (s *ProviderLoaderService) ReloadAllProviders(ctx context.Context) error {
	for _, provider := range s.loader.GetLoadedProviders() {
		if err := s.ReloadProvider(ctx, provider.Identifier); err != nil {
			return fmt.Errorf("failed to reload provider %s: %w", provider.Identifier, err)
		}
	}

	return nil
}

// ReloadProvider reload a specific provider
func (s *ProviderLoaderService) ReloadProvider(ctx context.Context, identifier string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.loader.UnloadProvider(identifier); err != nil {
		return fmt.Errorf("failed to unload provider %s: %w", identifier, err)
	}

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
	if err := s.loader.LoadProvider(config); err != nil {
		return fmt.Errorf("failed to reload provider %s: %w", identifier, err)
	}

	return nil
}
