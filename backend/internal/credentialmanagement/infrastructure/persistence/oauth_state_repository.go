package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"

	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/cache"
	"github.com/context-space/context-space/backend/internal/shared/utils"
	"gorm.io/gorm"
)

const (
	// OAuthStateKeyPrefix is the prefix for Redis keys storing OAuth states
	OAuthStateKeyPrefix = "oauth:state:"
	// OAuthStateIDKeyPrefix is the prefix for Redis keys storing OAuth states by ID
	OAuthStateIDKeyPrefix = "oauth:state_id:"
)

// RedisOAuthStateRepository implements OAuthStateRepository using Redis and PostgreSQL
type RedisOAuthStateRepository struct {
	db                database.Database
	redisClient       cache.Cache
	obs               *observability.ObservabilityProvider
	defaultExpiration time.Duration
}

// NewRedisOAuthStateRepository creates a new RedisOAuthStateRepository
func NewRedisOAuthStateRepository(db database.Database, redisClient cache.Cache, obs *observability.ObservabilityProvider, defaultExpiration time.Duration) *RedisOAuthStateRepository {
	return &RedisOAuthStateRepository{
		db:                db,
		redisClient:       redisClient,
		obs:               obs,
		defaultExpiration: defaultExpiration,
	}
}

// StoreStateData stores the state with associated redirect URI and user data
func (r *RedisOAuthStateRepository) StoreStateData(ctx context.Context, data *domain.OAuthStateData, expiration time.Duration) error {
	// If expiration is 0, use the default expiration
	if expiration == 0 {
		expiration = r.defaultExpiration
	}

	// Create Redis keys
	stateKey := utils.StringsBuilder(OAuthStateKeyPrefix, data.State)
	idKey := utils.StringsBuilder(OAuthStateIDKeyPrefix, data.ID)

	// Convert to JSON for Redis
	jsonData, err := sonic.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal OAuth state data: %w", err)
	}

	// Store in Redis with state key
	if err := r.redisClient.Set(ctx, stateKey, string(jsonData), expiration); err != nil {
		return fmt.Errorf("failed to store OAuth state in Redis: %w", err)
	}

	// Also store with ID key for efficient lookups
	if err := r.redisClient.Set(ctx, idKey, string(jsonData), expiration); err != nil {
		return fmt.Errorf("failed to store OAuth state by ID in Redis: %w", err)
	}

	// Convert domain object to database model
	model, err := r.mapToModel(data)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %w", err)
	}

	// Store in database
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return fmt.Errorf("failed to store OAuth state in database: %w", result.Error)
	}

	return nil
}

// GetStateDataByState gets the data associated with a state
func (r *RedisOAuthStateRepository) GetStateDataByState(ctx context.Context, state string) (*domain.OAuthStateData, error) {
	// First try Redis
	key := utils.StringsBuilder(OAuthStateKeyPrefix, state)
	jsonData, err := r.redisClient.Get(ctx, key)

	// If found in Redis, parse and return
	if err == nil {
		var data domain.OAuthStateData
		if err := sonic.Unmarshal([]byte(jsonData), &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal OAuth state data: %w", err)
		}
		return &data, nil
	}

	// If not in Redis, try database
	var model OAuthStateModel
	result := r.db.WithContext(ctx).Where("state = ?", state).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get OAuth state from database: %w", result.Error)
	}

	// Convert model to domain object
	data, err := r.mapToDomain(&model)
	if err != nil {
		return nil, fmt.Errorf("failed to convert model to domain: %w", err)
	}

	// Cache in Redis for future lookups
	jsonBytes, err := sonic.Marshal(data)
	if err == nil {
		r.redisClient.Set(ctx, key, string(jsonBytes), r.defaultExpiration)
	}

	return data, nil
}

// GetStateDataByID gets the data associated with a state by ID
func (r *RedisOAuthStateRepository) GetStateDataByID(ctx context.Context, id string) (*domain.OAuthStateData, error) {
	// First try Redis (for frequent polling)
	key := utils.StringsBuilder(OAuthStateIDKeyPrefix, id)
	jsonData, err := r.redisClient.Get(ctx, key)

	// If found in Redis, parse and return
	if err == nil {
		var data domain.OAuthStateData
		if err := sonic.Unmarshal([]byte(jsonData), &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal OAuth state data: %w", err)
		}
		return &data, nil
	}

	// If not in Redis, try database
	var model OAuthStateModel
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get OAuth state from database: %w", result.Error)
	}

	// Convert model to domain object
	data, err := r.mapToDomain(&model)
	if err != nil {
		return nil, fmt.Errorf("failed to convert model to domain: %w", err)
	}

	// Cache in Redis for future lookups (will improve polling performance)
	jsonBytes, err := sonic.Marshal(data)
	if err == nil {
		r.redisClient.Set(ctx, key, string(jsonBytes), r.defaultExpiration)
	}

	return data, nil
}

// DeleteState deletes a state
func (r *RedisOAuthStateRepository) DeleteState(ctx context.Context, state string) error {
	// Get data first to find the ID
	data, err := r.GetStateDataByState(ctx, state)
	if err != nil {
		return fmt.Errorf("failed to get OAuth state before deletion: %w", err)
	}

	if data == nil {
		// State not found, nothing to delete
		return nil
	}

	// Create Redis keys
	stateKey := utils.StringsBuilder(OAuthStateKeyPrefix, state)
	idKey := utils.StringsBuilder(OAuthStateIDKeyPrefix, data.ID)

	// Delete both entries from Redis
	if err := r.redisClient.Delete(ctx, stateKey); err != nil {
		return fmt.Errorf("failed to delete OAuth state from Redis: %w", err)
	}

	if err := r.redisClient.Delete(ctx, idKey); err != nil {
		return fmt.Errorf("failed to delete OAuth state by ID from Redis: %w", err)
	}

	// Delete from database (soft delete)
	result := r.db.WithContext(ctx).Where("state = ?", state).Delete(&OAuthStateModel{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete OAuth state from database: %w", result.Error)
	}

	return nil
}

// UpdateStateData updates the data of a state
func (r *RedisOAuthStateRepository) UpdateStateData(ctx context.Context, data *domain.OAuthStateData) error {
	// Verify the state exists without overwriting the input data
	existingData, err := r.GetStateDataByID(ctx, data.ID)
	if err != nil {
		return fmt.Errorf("failed to get OAuth state for status update: %w", err)
	}

	if existingData == nil {
		return fmt.Errorf("OAuth state not found for status update: %s", data.ID)
	}

	// Update Redis
	stateKey := utils.StringsBuilder(OAuthStateKeyPrefix, data.State)
	idKey := utils.StringsBuilder(OAuthStateIDKeyPrefix, data.ID)

	// Convert to JSON
	jsonData, err := sonic.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal updated OAuth state data: %w", err)
	}

	// Update in Redis
	if err := r.redisClient.Set(ctx, stateKey, string(jsonData), r.defaultExpiration); err != nil {
		return fmt.Errorf("failed to update OAuth state in Redis: %w", err)
	}

	if err := r.redisClient.Set(ctx, idKey, string(jsonData), r.defaultExpiration); err != nil {
		return fmt.Errorf("failed to update OAuth state by ID in Redis: %w", err)
	}

	// Map to model
	model, err := r.mapToModel(data)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %w", err)
	}

	// Update in database
	result := r.db.WithContext(ctx).Model(&OAuthStateModel{}).
		Where("id = ?", data.ID).
		Updates(map[string]interface{}{
			"status":          model.Status,
			"json_attributes": model.JSONAttributes,
			"updated_at":      model.UpdatedAt,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update OAuth state data in database: %w", result.Error)
	}

	return nil
}

// mapToModel converts OAuthStateData to OAuthStateModel
func (r *RedisOAuthStateRepository) mapToModel(data *domain.OAuthStateData) (*OAuthStateModel, error) {
	// Create attributes JSON
	attributes := struct {
		Permissions    []string               `json:"permissions"`
		RedirectURL    string                 `json:"redirect_url"`
		UserData       map[string]interface{} `json:"user_data"`
		CallbackParams map[string]interface{} `json:"callback_params"`
	}{
		Permissions:    data.Permissions,
		RedirectURL:    data.RedirectURL,
		UserData:       data.UserData,
		CallbackParams: data.CallbackParams,
	}

	jsonAttributes, err := sonic.Marshal(attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OAuth state attributes: %w", err)
	}

	return &OAuthStateModel{
		ID:                 data.ID,
		State:              data.State,
		Status:             string(data.Status),
		UserID:             data.UserID,
		ProviderIdentifier: data.ProviderIdentifier,
		JSONAttributes:     jsonAttributes,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
		DeletedAt:          parseDomainDeletedAt(data.DeletedAt),
	}, nil
}

// mapToDomain converts OAuthStateModel to OAuthStateData
func (r *RedisOAuthStateRepository) mapToDomain(model *OAuthStateModel) (*domain.OAuthStateData, error) {
	var attributes struct {
		Permissions    []string               `json:"permissions"`
		RedirectURL    string                 `json:"redirect_url"`
		UserData       map[string]interface{} `json:"user_data"`
		CallbackParams map[string]interface{} `json:"callback_params"`
	}

	if err := sonic.Unmarshal(model.JSONAttributes, &attributes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OAuth state attributes: %w", err)
	}

	return &domain.OAuthStateData{
		ID:                 model.ID,
		State:              model.State,
		Status:             domain.OAuthStateStatus(model.Status),
		UserID:             model.UserID,
		ProviderIdentifier: model.ProviderIdentifier,
		Permissions:        attributes.Permissions,
		RedirectURL:        attributes.RedirectURL,
		UserData:           attributes.UserData,
		CallbackParams:     attributes.CallbackParams,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
		DeletedAt:          parseGormDeletedAt(model.DeletedAt),
	}, nil
}
