package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"gorm.io/gorm"
)

// CredentialRepository implements the domain.CredentialRepository interface
type CredentialRepository struct {
	db  database.Database
	obs *observability.ObservabilityProvider
}

// NewCredentialRepository creates a new credential repository
func NewCredentialRepository(db database.Database, observabilityProvider *observability.ObservabilityProvider) *CredentialRepository {
	return &CredentialRepository{
		db:  db,
		obs: observabilityProvider,
	}
}

// GetByID retrieves a credential by ID
func (r *CredentialRepository) GetByID(ctx context.Context, id string) (*domain.Credential, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "CredentialRepository.GetByID")
	defer span.End()

	var model CredentialModel
	result := r.db.WithContext(ctx).First(&model, "id = ? AND is_valid = ?", id, true)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.mapToDomain(&model), nil
}

// GetByUserAndProvider retrieves a credential by user ID and provider ID
func (r *CredentialRepository) GetByUserAndProvider(ctx context.Context, userID, providerID string) (*domain.Credential, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "CredentialRepository.GetByUserAndProvider")
	defer span.End()

	var model CredentialModel
	result := r.db.WithContext(ctx).First(&model, "user_id = ? AND provider_identifier = ? AND is_valid = ?", userID, providerID, true)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.mapToDomain(&model), nil
}

// ListByUser retrieves all credentials for a user
func (r *CredentialRepository) ListByUser(ctx context.Context, userID string) ([]*domain.Credential, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "CredentialRepository.ListByUser")
	defer span.End()

	var models []CredentialModel
	result := r.db.WithContext(ctx).Find(&models, "user_id = ? AND is_valid = ?", userID, true)
	if result.Error != nil {
		return nil, result.Error
	}

	credentials := make([]*domain.Credential, len(models))
	for i, model := range models {
		credentials[i] = r.mapToDomain(&model)
	}

	return credentials, nil
}

// ListByID retrieves credentials by IDs
func (r *CredentialRepository) ListByID(ctx context.Context, ids []string) ([]*domain.Credential, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "CredentialRepository.ListByID")
	defer span.End()

	var models []CredentialModel
	result := r.db.WithContext(ctx).Find(&models, "id IN (?)", ids)
	if result.Error != nil {
		return nil, result.Error
	}

	credentials := make([]*domain.Credential, len(models))
	for i, model := range models {
		credentials[i] = r.mapToDomain(&model)
	}
	return credentials, nil
}

// Create creates a new credential
func (r *CredentialRepository) Create(ctx context.Context, credential *domain.Credential) error {
	ctx, span := r.obs.Tracer.Start(ctx, "CredentialRepository.Create")
	defer span.End()

	model := r.mapToModel(credential)

	result := r.db.WithContext(ctx).Create(model)
	return result.Error
}

// Delete soft-deletes a credential
func (r *CredentialRepository) Delete(ctx context.Context, id string) error {
	ctx, span := r.obs.Tracer.Start(ctx, "CredentialRepository.Delete")
	defer span.End()

	err := r.db.Transaction(ctx, func(tx *gorm.DB) error {
		// Get the credential model
		var model CredentialModel
		result := tx.Where("id = ?", id).First(&model)
		if result.Error != nil {
			return result.Error
		}

		// Delete the specific credential model
		switch model.CredentialType {
		case string(domain.CredentialTypeOAuth):
			result = tx.Where("credential_id = ?", id).Delete(&OAuthCredentialModel{})
		case string(domain.CredentialTypeAPIKey):
			result = tx.Where("credential_id = ?", id).Delete(&APIKeyCredentialModel{})
		default:
			return fmt.Errorf("credential type not found with id: %s", id)
		}
		if result.Error != nil {
			return result.Error
		}

		result = tx.Where("id = ?", id).Delete(&CredentialModel{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("credential not found with id: %s", id)
		}

		return nil
	})

	return err
}

// UpdateLastUsedAt updates the last used at time of a credential
func (r *CredentialRepository) UpdateLastUsedAt(ctx context.Context, id string) error {
	ctx, span := r.obs.Tracer.Start(ctx, "CredentialRepository.UpdateLastUsedAt")
	defer span.End()

	return r.db.WithContext(ctx).Model(&CredentialModel{}).Where("id = ?", id).Update("last_used_at", time.Now()).Error
}

// mapToDomain maps a credential model to a domain credential
func (r *CredentialRepository) mapToDomain(model *CredentialModel) *domain.Credential {
	return &domain.Credential{
		ID:                 model.ID,
		UserID:             model.UserID,
		ProviderIdentifier: model.ProviderIdentifier,
		Type:               domain.CredentialType(model.CredentialType),
		IsValid:            model.IsValid,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
		LastUsedAt:         model.LastUsedAt,
		DeletedAt:          parseGormDeletedAt(model.DeletedAt),
	}
}

// mapToModel maps a domain credential to a credential model
func (r *CredentialRepository) mapToModel(credential *domain.Credential) *CredentialModel {
	return &CredentialModel{
		ID:                 credential.ID,
		UserID:             credential.UserID,
		ProviderIdentifier: credential.ProviderIdentifier,
		CredentialType:     string(credential.Type),
		IsValid:            credential.IsValid,
		CreatedAt:          credential.CreatedAt,
		UpdatedAt:          credential.UpdatedAt,
		DeletedAt:          parseDomainDeletedAt(credential.DeletedAt),
	}
}
