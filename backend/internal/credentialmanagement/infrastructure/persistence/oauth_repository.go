package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"gorm.io/gorm"
)

// OAuthCredentialRepository implements the domain.OAuthCredentialRepository interface
type OAuthCredentialRepository struct {
	db  database.Database
	obs *observability.ObservabilityProvider
}

// NewOAuthCredentialRepository creates a new OAuth credential repository
func NewOAuthCredentialRepository(db database.Database, observabilityProvider *observability.ObservabilityProvider) *OAuthCredentialRepository {
	return &OAuthCredentialRepository{
		db:  db,
		obs: observabilityProvider,
	}
}

// GetByCredentialID retrieves an OAuth credential by credential ID
func (r *OAuthCredentialRepository) GetByCredentialID(ctx context.Context, credentialID string) (*domain.OAuthCredential, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "OAuthCredentialRepository.GetByCredentialID")
	defer span.End()

	var credModel CredentialModel
	result := r.db.WithContext(ctx).First(&credModel, "id = ? AND is_valid = true", credentialID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	var oauthModel OAuthCredentialModel
	result = r.db.WithContext(ctx).Where("credential_id = ?", credentialID).First(&oauthModel)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.mapToDomain(&credModel, &oauthModel)
}

// Create creates a new OAuth credential
func (r *OAuthCredentialRepository) Create(ctx context.Context, credential *domain.OAuthCredential) error {
	ctx, span := r.obs.Tracer.Start(ctx, "OAuthCredentialRepository.Create")
	defer span.End()

	oauthModel, err := r.mapToModel(credential)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Create(oauthModel)
	return result.Error
}

// ListByExpiryWithin lists OAuth credentials by expiry within the given time
func (r *OAuthCredentialRepository) ListByExpiryWithin(ctx context.Context, expiry time.Time) ([]*domain.OAuthCredential, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "OAuthCredentialRepository.ListByExpiryWithin")
	defer span.End()

	var oauthCredentialModels []OAuthCredentialModel
	result := r.db.WithContext(ctx).
		Where("expiry <= ?", expiry).
		Where("expiry > ?", time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)).
		Find(&oauthCredentialModels)
	if result.Error != nil {
		return nil, result.Error
	}

	oauthCredentials := make([]*domain.OAuthCredential, 0, len(oauthCredentialModels))
	for _, oauthCredentialModel := range oauthCredentialModels {
		credModel := CredentialModel{
			ID: oauthCredentialModel.CredentialID,
		}
		cred, err := r.mapToDomain(&credModel, &oauthCredentialModel)
		if err != nil {
			return nil, err
		}
		oauthCredentials = append(oauthCredentials, cred)
	}
	return oauthCredentials, nil
}

// Update updates an OAuth credential
func (r *OAuthCredentialRepository) Update(ctx context.Context, credential *domain.OAuthCredential) error {
	ctx, span := r.obs.Tracer.Start(ctx, "OAuthCredentialRepository.Update")
	defer span.End()

	oauthModel, err := r.mapToModel(credential)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Save(oauthModel).Error
}

// mapToDomain maps credential models to a domain OAuth credential
func (r *OAuthCredentialRepository) mapToDomain(credModel *CredentialModel, oauthCredentialModel *OAuthCredentialModel) (*domain.OAuthCredential, error) {
	var jsonAttributes struct {
		EncryptionMetadata *domain.EncryptionMetadata `json:"encryption_metadata"`
		Scopes             []string                   `json:"scopes"`
	}

	if err := sonic.Unmarshal(oauthCredentialModel.JSONAttributes, &jsonAttributes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal oauth credential json attributes: %w", err)
	}

	cred := &domain.Credential{
		ID:                 credModel.ID,
		UserID:             credModel.UserID,
		ProviderIdentifier: credModel.ProviderIdentifier,
		Type:               domain.CredentialType(credModel.CredentialType),
		IsValid:            credModel.IsValid,
		CreatedAt:          credModel.CreatedAt,
		UpdatedAt:          credModel.UpdatedAt,
		DeletedAt:          parseGormDeletedAt(credModel.DeletedAt),
	}

	return &domain.OAuthCredential{
		Credential:         cred,
		EncryptionMetadata: jsonAttributes.EncryptionMetadata,
		Token:              nil,
		Scopes:             jsonAttributes.Scopes,
	}, nil
}

// mapToModel maps a domain OAuth credential to an OAuth credential model
func (r *OAuthCredentialRepository) mapToModel(credential *domain.OAuthCredential) (*OAuthCredentialModel, error) {
	jsonAttributes := struct {
		EncryptionMetadata *domain.EncryptionMetadata `json:"encryption_metadata"`
		Scopes             []string                   `json:"scopes"`
	}{
		EncryptionMetadata: credential.EncryptionMetadata,
		Scopes:             credential.Scopes,
	}

	jsonAttributesJSON, err := sonic.Marshal(jsonAttributes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal oauth credential json attributes: %w", err)
	}

	return &OAuthCredentialModel{
		CredentialID:   credential.ID,
		Expiry:         credential.Token.Expiry.Local(),
		JSONAttributes: jsonAttributesJSON,
		CreatedAt:      credential.CreatedAt,
		UpdatedAt:      credential.UpdatedAt,
		DeletedAt:      parseDomainDeletedAt(credential.DeletedAt),
	}, nil
}
