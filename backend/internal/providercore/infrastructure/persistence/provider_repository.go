package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/providercore/domain"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"github.com/context-space/context-space/backend/internal/shared/types"
	"gorm.io/gorm"
)

// ProviderRepository implements the domain.ProviderRepository interface
type ProviderRepository struct {
	db  database.Database
	obs *observability.ObservabilityProvider
}

// NewProviderRepository creates a new provider repository
func NewProviderRepository(db database.Database, observabilityProvider *observability.ObservabilityProvider) *ProviderRepository {
	return &ProviderRepository{
		db:  db,
		obs: observabilityProvider,
	}
}

// GetByID retrieves a provider by ID
func (r *ProviderRepository) GetByID(ctx context.Context, id string) (*domain.Provider, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.GetByID")
	defer span.End()

	var model ProviderModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Get the provider without operations first
	provider, err := r.mapToDomain(&model)
	if err != nil {
		return nil, err
	}

	// Load operations for this provider
	operations, err := r.loadOperations(ctx, provider.ID)
	if err != nil {
		return nil, err
	}

	// Set operations on the provider
	provider.Operations = operations
	return provider, nil
}

// ListByIDs retrieves a list of providers by IDs
func (r *ProviderRepository) ListByIDs(ctx context.Context, ids []string) ([]*domain.Provider, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.ListByIDs")
	defer span.End()

	var models []ProviderModel
	result := r.db.WithContext(ctx).Where("id IN (?)", ids).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	providers := make([]*domain.Provider, 0, len(models))
	for i := range models {
		provider, err := r.mapToDomain(&models[i])
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	return providers, nil
}

// GetByIdentifier retrieves a provider by identifier
func (r *ProviderRepository) GetByIdentifier(ctx context.Context, identifier string) (*domain.Provider, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.GetByIdentifier")
	defer span.End()

	var model ProviderModel
	result := r.db.WithContext(ctx).First(&model, "identifier = ?", identifier)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Get the provider without operations first
	provider, err := r.mapToDomain(&model)
	if err != nil {
		return nil, err
	}

	// Load operations for this provider
	operations, err := r.loadOperations(ctx, provider.ID)
	if err != nil {
		return nil, err
	}

	// Set operations on the provider
	provider.Operations = operations
	return provider, nil
}

// ListFullProviders returns all providers with operations
func (r *ProviderRepository) ListFullProviders(ctx context.Context) ([]*domain.Provider, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.List")
	defer span.End()

	var models []ProviderModel
	result := r.db.WithContext(ctx).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	providers := make([]*domain.Provider, 0, len(models))
	for i := range models {
		provider, err := r.mapToDomain(&models[i])
		if err != nil {
			return nil, err
		}

		// Load operations for this provider
		operations, err := r.loadOperations(ctx, provider.ID)
		if err != nil {
			return nil, err
		}

		// Set operations on the provider
		provider.Operations = operations
		providers = append(providers, provider)
	}

	return providers, nil
}

func (r *ProviderRepository) ListBasicProviders(ctx context.Context) ([]*domain.Provider, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.ListBasicProviders")
	defer span.End()

	var models []ProviderModel
	result := r.db.WithContext(ctx).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	providers := make([]*domain.Provider, 0, len(models))
	for i := range models {
		provider, err := r.mapToDomain(&models[i])
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	return providers, nil
}

// Create creates a new provider
func (r *ProviderRepository) Create(ctx context.Context, provider *domain.Provider) error {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.Create")
	defer span.End()

	model, err := r.mapToModel(provider)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Create(model)
	return result.Error
}

// Update updates a provider
func (r *ProviderRepository) Update(ctx context.Context, provider *domain.Provider) error {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.Update")
	defer span.End()

	model, err := r.mapToModel(provider)
	if err != nil {
		return err
	}

	// Use Where clause to update by ID to avoid unique constraint issues
	result := r.db.WithContext(ctx).Where("id = ?", provider.ID).Updates(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete soft-deletes a provider
func (r *ProviderRepository) Delete(ctx context.Context, id string) error {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.Delete")
	defer span.End()

	return r.db.Transaction(ctx, func(tx *gorm.DB) error {
		// First, get the provider identifier for deleting translations
		var providerModel ProviderModel
		result := tx.Where("id = ?", id).First(&providerModel)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return fmt.Errorf("provider not found with id: %s", id)
			}
			return result.Error
		}

		// Delete operations associated with this provider
		result = tx.Where("provider_id = ?", id).Delete(&OperationModel{})
		if result.Error != nil {
			return result.Error
		}

		// Delete translations associated with this provider
		result = tx.Where("provider_identifier = ?", providerModel.Identifier).Delete(&ProviderTranslationModel{})
		if result.Error != nil {
			return result.Error
		}

		// Finally, delete the provider
		result = tx.Where("id = ?", id).Delete(&ProviderModel{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("provider not found with id: %s", id)
		}

		return nil
	})
}

// loadOperations loads all operations for a provider
func (r *ProviderRepository) loadOperations(ctx context.Context, providerID string) ([]domain.Operation, error) {
	// Create operation repository
	operationRepo := NewOperationRepository(r.db, r.obs)

	// Get operations for this provider
	operationPtrs, err := operationRepo.ListByProviderID(ctx, providerID)
	if err != nil {
		return nil, err
	}

	// Convert from []*domain.Operation to []domain.Operation
	operations := make([]domain.Operation, len(operationPtrs))
	for i, op := range operationPtrs {
		operations[i] = *op
	}

	return operations, nil
}

// mapToDomain maps a provider model to a domain provider
func (r *ProviderRepository) mapToDomain(model *ProviderModel) (*domain.Provider, error) {
	var jsonAttributes struct {
		Categories []string `json:"categories"`
		Tags       []string `json:"tags"`
	}

	if err := sonic.Unmarshal(model.JSONAttributes, &jsonAttributes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal provider info metadata: %w", err)
	}

	return &domain.Provider{
		ID:          model.ID,
		Identifier:  model.Identifier,
		Name:        model.Name,
		Description: model.Description,
		AuthType:    types.ProviderAuthType(model.AuthType),
		Status:      types.ProviderStatus(model.Status),
		IconURL:     model.IconURL,
		Categories:  jsonAttributes.Categories,
		Tags:        jsonAttributes.Tags,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   parseGormDeletedAt(model.DeletedAt),
	}, nil
}

// mapToModel maps a domain provider to a provider model
func (r *ProviderRepository) mapToModel(provider *domain.Provider) (*ProviderModel, error) {
	jsonAttributes := struct {
		Categories []string `json:"categories"`
	}{
		Categories: provider.Categories,
	}

	jsonAttributesJSON, err := sonic.Marshal(jsonAttributes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal provider info metadata: %w", err)
	}

	return &ProviderModel{
		ID:             provider.ID,
		Identifier:     provider.Identifier,
		Name:           provider.Name,
		Description:    provider.Description,
		AuthType:       string(provider.AuthType),
		Status:         string(provider.Status),
		IconURL:        provider.IconURL,
		JSONAttributes: jsonAttributesJSON,
		CreatedAt:      provider.CreatedAt,
		UpdatedAt:      provider.UpdatedAt,
		DeletedAt:      parseDomainDeletedAt(provider.DeletedAt),
	}, nil
}

// SyncTagsToProvider syncs tags to provider's json_attributes
func (r *ProviderRepository) SyncTagsToProvider(ctx context.Context, providerIdentifier string, tags []string) error {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.SyncTagsToProvider")
	defer span.End()

	// Get existing provider
	var model ProviderModel
	result := r.db.WithContext(ctx).First(&model, "identifier = ?", providerIdentifier)
	if result.Error != nil {
		return result.Error
	}

	// Parse existing json_attributes
	var jsonAttributes struct {
		Categories []string `json:"categories"`
		Tags       []string `json:"tags"`
	}

	if err := sonic.Unmarshal(model.JSONAttributes, &jsonAttributes); err != nil {
		return fmt.Errorf("failed to unmarshal provider json_attributes: %w", err)
	}

	// Update tags field
	jsonAttributes.Tags = tags

	// Re-serialize
	updatedJSON, err := sonic.Marshal(jsonAttributes)
	if err != nil {
		return fmt.Errorf("failed to marshal provider json_attributes: %w", err)
	}

	// Update database
	result = r.db.WithContext(ctx).Model(&model).Where("identifier = ?", providerIdentifier).Update("json_attributes", updatedJSON)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
