package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"gorm.io/gorm"
)

// APIKeyCredentialRepository implements the domain.APIKeyCredentialRepository interface
type APIKeyCredentialRepository struct {
	db  database.Database
	obs *observability.ObservabilityProvider
}

// NewAPIKeyCredentialRepository creates a new API key credential repository
func NewAPIKeyCredentialRepository(db database.Database, observabilityProvider *observability.ObservabilityProvider) *APIKeyCredentialRepository {
	return &APIKeyCredentialRepository{
		db:  db,
		obs: observabilityProvider,
	}
}

// GetByCredentialID retrieves an API key credential by credential ID
func (r *APIKeyCredentialRepository) GetByCredentialID(ctx context.Context, credentialID string) (*domain.APIKeyCredential, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "APIKeyCredentialRepository.GetByCredentialID")
	defer span.End()

	var credModel CredentialModel
	result := r.db.WithContext(ctx).First(&credModel, "id = ? AND is_valid = true", credentialID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	var apiKeyModel APIKeyCredentialModel
	result = r.db.WithContext(ctx).Where("credential_id = ?", credentialID).First(&apiKeyModel)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.mapToDomain(&credModel, &apiKeyModel)
}

// Create creates a new API key credential
func (r *APIKeyCredentialRepository) Create(ctx context.Context, credential *domain.APIKeyCredential) error {
	ctx, span := r.obs.Tracer.Start(ctx, "APIKeyCredentialRepository.Create")
	defer span.End()

	apiKeyModel, err := r.mapToModel(credential)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Create(apiKeyModel)
	return result.Error
}

// mapToDomain maps credential models to a domain API key credential
func (r *APIKeyCredentialRepository) mapToDomain(credModel *CredentialModel, apiKeyModel *APIKeyCredentialModel) (*domain.APIKeyCredential, error) {
	var jsonAttributes struct {
		EncryptionMetadata *domain.EncryptionMetadata `json:"encryption_metadata"`
	}

	if err := sonic.Unmarshal(apiKeyModel.JSONAttributes, &jsonAttributes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal api key credential json attributes: %w", err)
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

	return &domain.APIKeyCredential{
		Credential:         cred,
		EncryptionMetadata: jsonAttributes.EncryptionMetadata,
	}, nil
}

func (r *APIKeyCredentialRepository) mapToModel(credential *domain.APIKeyCredential) (*APIKeyCredentialModel, error) {
	jsonAttributes := struct {
		EncryptionMetadata *domain.EncryptionMetadata `json:"encryption_metadata"`
	}{
		EncryptionMetadata: credential.EncryptionMetadata,
	}

	jsonAttributesJSON, err := sonic.Marshal(jsonAttributes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal api key credential json attributes: %w", err)
	}

	return &APIKeyCredentialModel{
		CredentialID:   credential.ID,
		JSONAttributes: jsonAttributesJSON,
		CreatedAt:      credential.CreatedAt,
		UpdatedAt:      credential.UpdatedAt,
		DeletedAt:      parseDomainDeletedAt(credential.DeletedAt),
	}, nil
}
