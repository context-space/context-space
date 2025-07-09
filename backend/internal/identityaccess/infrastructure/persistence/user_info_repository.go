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

// UserInfoRepository implements the domain.UserInfoRepository interface
type UserInfoRepository struct {
	db  database.Database
	obs *observability.ObservabilityProvider
}

// NewUserInfoRepository creates a new user info repository
func NewUserInfoRepository(db database.Database, observability *observability.ObservabilityProvider) *UserInfoRepository {
	return &UserInfoRepository{
		db:  db,
		obs: observability,
	}
}

// Get gets a user info by id
func (r *UserInfoRepository) Get(ctx context.Context, id string) (*domain.UserInfo, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "UserInfoRepository.Get")
	defer span.End()

	var model UserInfoModel
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.mapToDomain(&model), nil
}

// GetByUserID gets a user info by user id
func (r *UserInfoRepository) GetByUserID(ctx context.Context, userID string) (*domain.UserInfo, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "UserInfoRepository.GetByUserID")
	defer span.End()

	var model UserInfoModel
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.mapToDomain(&model), nil
}

// Create creates a user info
func (r *UserInfoRepository) Create(ctx context.Context, info *domain.UserInfo) error {
	ctx, span := r.obs.Tracer.Start(ctx, "UserInfoRepository.Create")
	defer span.End()

	model := r.mapToModel(info)

	result := r.db.WithContext(ctx).Create(&model)
	return result.Error
}

// Update updates a user info
func (r *UserInfoRepository) Update(ctx context.Context, info *domain.UserInfo) error {
	ctx, span := r.obs.Tracer.Start(ctx, "UserInfoRepository.Update")
	defer span.End()

	model := r.mapToModel(info)

	result := r.db.WithContext(ctx).Model(&model).Updates(map[string]interface{}{
		"info_metadata": model.InfoMetadata,
		"updated_at":    model.UpdatedAt,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user info not found with id: %s", info.ID)
	}

	return nil
}

// Delete deletes a user info by id
func (r *UserInfoRepository) Delete(ctx context.Context, id string) error {
	ctx, span := r.obs.Tracer.Start(ctx, "UserInfoRepository.Delete")
	defer span.End()

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&UserInfoModel{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user info not found with id: %s", id)
	}

	return nil
}

// mapToDomain maps a user info model to a domain user info
func (r *UserInfoRepository) mapToDomain(model *UserInfoModel) *domain.UserInfo {
	return &domain.UserInfo{
		ID:           model.ID,
		UserID:       model.UserID,
		InfoMetadata: mustUnmarshalJSON(model.InfoMetadata),
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
		DeletedAt:    parseGormDeletedAt(model.DeletedAt),
	}
}

// mapToModel maps a domain user info to a user info model
func (r *UserInfoRepository) mapToModel(info *domain.UserInfo) *UserInfoModel {
	return &UserInfoModel{
		ID:           info.ID,
		UserID:       info.UserID,
		InfoMetadata: mustMarshalJSON(info.InfoMetadata),
		CreatedAt:    info.CreatedAt,
		UpdatedAt:    info.UpdatedAt,
		DeletedAt:    parseDomainDeletedAt(info.DeletedAt),
	}
}
