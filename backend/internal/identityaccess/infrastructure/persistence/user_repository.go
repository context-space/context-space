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

// UserRepository implements the domain.UserRepository interface
type UserRepository struct {
	db  database.Database
	obs *observability.ObservabilityProvider
}

// NewUserRepository creates a new user repository
func NewUserRepository(db database.Database, observabilityProvider *observability.ObservabilityProvider) *UserRepository {
	return &UserRepository{
		db:  db,
		obs: observabilityProvider,
	}
}

// Get retrieves a user by ID
func (r *UserRepository) Get(ctx context.Context, id string) (*domain.User, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "UserRepository.Get")
	defer span.End()

	var model UserModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Map to domain
	return r.mapToDomain(&model), nil
}

// GetBySupID retrieves a user by sup id
func (r *UserRepository) GetBySupID(ctx context.Context, supID string) (*domain.User, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "UserRepository.GetBySupID")
	defer span.End()

	var model UserModel
	result := r.db.WithContext(ctx).First(&model, "sup_id = ?", supID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.mapToDomain(&model), nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "UserRepository.GetByEmail")
	defer span.End()

	var model UserModel
	result := r.db.WithContext(ctx).First(&model, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Map to domain
	return r.mapToDomain(&model), nil
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "UserRepository.List")
	defer span.End()

	var models []UserModel
	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = r.mapToDomain(&model)
	}

	return users, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	ctx, span := r.obs.Tracer.Start(ctx, "UserRepository.Create")
	defer span.End()

	model := r.mapToModel(user)

	result := r.db.WithContext(ctx).Create(&model)
	return result.Error
}

// Delete soft-deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	ctx, span := r.obs.Tracer.Start(ctx, "UserRepository.Delete")
	defer span.End()

	err := r.db.Transaction(ctx, func(tx *gorm.DB) error {
		result := tx.Where("user_id = ?", id).Delete(&UserAPIKeyModel{})
		if result.Error != nil {
			return result.Error
		}

		result = tx.Where("user_id = ?", id).Delete(&UserInfoModel{})
		if result.Error != nil {
			return result.Error
		}

		result = tx.Where("id = ?", id).Delete(&UserModel{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("user not found with id: %s", id)
		}

		return nil
	})

	return err
}

// Count returns the total number of users
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "UserRepository.Count")
	defer span.End()

	var count int64
	result := r.db.WithContext(ctx).Model(&UserModel{}).Count(&count)
	return count, result.Error
}

// mapToDomain maps a user model to a domain user
func (r *UserRepository) mapToDomain(model *UserModel) *domain.User {
	return &domain.User{
		ID:          model.ID,
		SupID:       model.SupID,
		Email:       parseGormEmail(model.Email),
		IsAnonymous: model.IsAnonymous,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   parseGormDeletedAt(model.DeletedAt),
	}
}

// mapToModel maps a domain user to a user model
func (r *UserRepository) mapToModel(user *domain.User) *UserModel {
	return &UserModel{
		ID:          user.ID,
		SupID:       user.SupID,
		Email:       parseDomainEmail(user.Email),
		IsAnonymous: user.IsAnonymous,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   parseDomainDeletedAt(user.DeletedAt),
	}
}
