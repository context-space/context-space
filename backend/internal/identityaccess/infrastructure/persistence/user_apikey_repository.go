package persistence

import (
	"context"
	"errors"
	"fmt"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/identityaccess/domain"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"gorm.io/gorm"
)

// UserAPIKeyRepository implements the domain.APIKeyRepository interface
type UserAPIKeyRepository struct {
	db  database.Database
	obs *observability.ObservabilityProvider
}

// NewUserAPIKeyRepository creates a new API key repository
func NewUserAPIKeyRepository(db database.Database, observability *observability.ObservabilityProvider) *UserAPIKeyRepository {
	return &UserAPIKeyRepository{
		db:  db,
		obs: observability,
	}
}

// Get retrieves an API key by ID
func (r *UserAPIKeyRepository) Get(ctx context.Context, id string) (*domain.APIKey, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "UserAPIKeyRepository.Get")
	defer span.End()

	var model UserAPIKeyModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.mapToDomain(&model), nil
}

// GetByKeyValue retrieves an API key by its value
func (r *UserAPIKeyRepository) GetByKeyValue(ctx context.Context, value string) (*domain.APIKey, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "UserAPIKeyRepository.GetByKeyValue")
	defer span.End()

	var model UserAPIKeyModel
	result := r.db.WithContext(ctx).First(&model, "key_value = ?", value)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.mapToDomain(&model), nil
}

// ListByUserID retrieves API keys for a user with pagination
func (r *UserAPIKeyRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.APIKey, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "UserAPIKeyRepository.ListByUserID")
	defer span.End()

	var models []UserAPIKeyModel
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Map to domain
	apiKeys := make([]*domain.APIKey, len(models))
	for i, model := range models {
		apiKeys[i] = r.mapToDomain(&model)
	}

	return apiKeys, nil
}

// Create creates a new API key
func (r *UserAPIKeyRepository) Create(ctx context.Context, apiKey *domain.APIKey) error {
	ctx, span := r.obs.Tracer.Start(ctx, "UserAPIKeyRepository.Create")
	defer span.End()

	model := r.mapToModel(apiKey)

	result := r.db.WithContext(ctx).Create(&model)
	return result.Error
}

// Update updates an existing API key
func (r *UserAPIKeyRepository) Update(ctx context.Context, apiKey *domain.APIKey) error {
	ctx, span := r.obs.Tracer.Start(ctx, "UserAPIKeyRepository.Update")
	defer span.End()

	model := r.mapToModel(apiKey)

	result := r.db.WithContext(ctx).Model(&model).Updates(map[string]interface{}{
		"name":        apiKey.Name,
		"description": apiKey.Description,
		"last_used":   apiKey.LastUsed,
		"updated_at":  apiKey.UpdatedAt,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("API key not found with id: %s", apiKey.ID)
	}

	return nil
}

// Delete soft-deletes an API key
func (r *UserAPIKeyRepository) Delete(ctx context.Context, id string) error {
	ctx, span := r.obs.Tracer.Start(ctx, "UserAPIKeyRepository.Delete")
	defer span.End()

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&UserAPIKeyModel{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("API key not found with id: %s", id)
	}

	return nil
}

// mapToDomain maps an API key model to a domain API key
func (r *UserAPIKeyRepository) mapToDomain(model *UserAPIKeyModel) *domain.APIKey {
	return &domain.APIKey{
		ID:          model.ID,
		UserID:      model.UserID,
		KeyValue:    model.KeyValue,
		Name:        model.Name,
		Description: model.Description,
		LastUsed:    model.LastUsed,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   parseGormDeletedAt(model.DeletedAt),
	}
}

// mapToModel maps a domain API key to an API key model
func (r *UserAPIKeyRepository) mapToModel(apiKey *domain.APIKey) *UserAPIKeyModel {
	return &UserAPIKeyModel{
		ID:          apiKey.ID,
		UserID:      apiKey.UserID,
		KeyValue:    apiKey.KeyValue,
		Name:        apiKey.Name,
		Description: apiKey.Description,
		LastUsed:    apiKey.LastUsed,
		CreatedAt:   apiKey.CreatedAt,
		UpdatedAt:   apiKey.UpdatedAt,
		DeletedAt:   parseDomainDeletedAt(apiKey.DeletedAt),
	}
}
