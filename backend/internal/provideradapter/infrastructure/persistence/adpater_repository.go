package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"gorm.io/gorm"
)

type AdapterRepository struct {
	db  database.Database
	obs *observability.ObservabilityProvider
}

// NewProviderRepository creates a new provider repository
func NewAdapterRepository(db database.Database, observabilityProvider *observability.ObservabilityProvider) *AdapterRepository {
	return &AdapterRepository{
		db:  db,
		obs: observabilityProvider,
	}
}

// GetByID retrieves a provider adapter by ID
func (r *AdapterRepository) GetByID(ctx context.Context, id string) (*domain.ProviderAdapterConfig, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "AdapterRepository.GetByID")
	defer span.End()

	var model ProviderAdapterModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.mapToDomain(&model)
}

// GetByIdentifier retrieves a provider adapter by identifier
func (r *AdapterRepository) GetByIdentifier(ctx context.Context, identifier string) (*domain.ProviderAdapterConfig, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "AdapterRepository.GetByIdentifier")
	defer span.End()

	var model ProviderAdapterModel
	result := r.db.WithContext(ctx).First(&model, "identifier = ?", identifier)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.mapToDomain(&model)
}

// List returns all provider adapters
func (r *AdapterRepository) ListAdapterConfigs(ctx context.Context) ([]*domain.ProviderAdapterConfig, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "AdapterRepository.List")
	defer span.End()

	var models []ProviderAdapterModel
	result := r.db.WithContext(ctx).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	adapters := make([]*domain.ProviderAdapterConfig, len(models))
	for i, model := range models {
		adapter, err := r.mapToDomain(&model)
		if err != nil {
			return nil, err
		}
		adapters[i] = adapter
	}

	return adapters, nil
}

// Create creates a new provider adapter
func (r *AdapterRepository) Create(ctx context.Context, adapter *domain.ProviderAdapterConfig) error {
	ctx, span := r.obs.Tracer.Start(ctx, "AdapterRepository.Create")
	defer span.End()

	model, err := r.mapToModel(adapter)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Create(&model)
	if result.Error != nil {
		return result.Error
	}

	// Update the ID in the domain model
	adapter.ID = model.ID
	return nil
}

// Update updates a provider adapter
func (r *AdapterRepository) Update(ctx context.Context, adapter *domain.ProviderAdapterConfig) error {
	ctx, span := r.obs.Tracer.Start(ctx, "AdapterRepository.Update")
	defer span.End()

	model, err := r.mapToModel(adapter)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Save(&model)
	return result.Error
}

// Delete deletes a provider adapter
func (r *AdapterRepository) Delete(ctx context.Context, id string) error {
	ctx, span := r.obs.Tracer.Start(ctx, "AdapterRepository.Delete")
	defer span.End()

	result := r.db.WithContext(ctx).Delete(&ProviderAdapterModel{}, "id = ?", id)
	return result.Error
}

// mapToDomain converts a persistence model to a domain model
func (r *AdapterRepository) mapToDomain(model *ProviderAdapterModel) (*domain.ProviderAdapterConfig, error) {
	var jsonAttributes struct {
		OAuthConfig  *domain.OAuthConfig    `json:"oauth_config"`
		CustomConfig map[string]interface{} `json:"custom_config"`
	}

	if err := sonic.Unmarshal(model.Configs, &jsonAttributes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal provider info metadata: %w", err)
	}

	return &domain.ProviderAdapterConfig{
		ProviderAdapterInfo: domain.ProviderAdapterInfo{
			Identifier: model.Identifier,
		},
		ID:           model.ID,
		OAuthConfig:  jsonAttributes.OAuthConfig,
		CustomConfig: jsonAttributes.CustomConfig,
	}, nil
}

// mapToModel converts a domain model to a persistence model
func (r *AdapterRepository) mapToModel(adapter *domain.ProviderAdapterConfig) (*ProviderAdapterModel, error) {
	authConfig := struct {
		OAuthConfig  *domain.OAuthConfig    `json:"oauth_config"`
		CustomConfig map[string]interface{} `json:"custom_config"`
	}{
		OAuthConfig:  adapter.OAuthConfig,
		CustomConfig: adapter.CustomConfig,
	}

	authConfigJSON, err := sonic.Marshal(authConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal auth config: %w", err)
	}

	return &ProviderAdapterModel{
		ID:         adapter.ID,
		Identifier: adapter.Identifier,
		Configs:    authConfigJSON,
	}, nil
}
