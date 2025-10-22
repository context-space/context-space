package persistence

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/providercore/domain"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/cache"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"github.com/context-space/context-space/backend/internal/shared/types"
	"github.com/context-space/context-space/backend/internal/shared/utils"
	"gorm.io/gorm"
)

const (
	cacheKeyOperationByID                      = "operation_by_id:"
	cacheKeyOperationByProviderIDAndIdentifier = "operation_by_provider_id_and_identifier:"
)

// OperationRepository implements the domain.OperationRepository interface
type OperationRepository struct {
	db        database.Database
	obs       *observability.ObservabilityProvider
	cache     *cache.LRUCache[string, *domain.Operation]
	cacheOnce sync.Once
}

// NewOperationRepository creates a new operation repository
func NewOperationRepository(db database.Database, observabilityProvider *observability.ObservabilityProvider) *OperationRepository {
	repo := &OperationRepository{
		db:  db,
		obs: observabilityProvider,
	}
	repo.cacheOnce.Do(func() {
		repo.cache = cache.NewLRUCache[string, *domain.Operation](500, 1*time.Hour)
	})
	return repo
}

// ListByProviderID returns all operations for a provider
func (r *OperationRepository) ListByProviderID(ctx context.Context, providerID string) ([]*domain.Operation, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "OperationRepository.ListByProviderID")
	defer span.End()

	var models []OperationModel
	result := r.db.WithContext(ctx).Where("provider_id = ?", providerID).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	operations := make([]*domain.Operation, 0, len(models))
	for i := range models {
		operation, err := r.mapToDomain(&models[i])
		if err != nil {
			return nil, err
		}
		operations = append(operations, operation)
	}

	batchCache := make(map[string]*domain.Operation, len(operations))
	for _, operation := range operations {
		batchCache[utils.StringsBuilder(cacheKeyOperationByProviderIDAndIdentifier, operation.ProviderID, ":", operation.Identifier)] = operation
	}
	r.cache.BatchSet(batchCache)

	return operations, nil
}

// GetByID returns an operation by ID
func (r *OperationRepository) GetByID(ctx context.Context, id string) (*domain.Operation, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "OperationRepository.GetByID")
	defer span.End()

	cachekey := utils.StringsBuilder(cacheKeyOperationByID, id)
	if operation, ok := r.cache.Get(cachekey); ok {
		return operation, nil
	}

	var model OperationModel
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	operation, err := r.mapToDomain(&model)
	if err != nil {
		return nil, err
	}

	r.cache.Set(cachekey, operation)
	return operation, nil
}

// ListByIDs returns a list of operations by IDs
func (r *OperationRepository) ListByIDs(ctx context.Context, ids []string) ([]*domain.Operation, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "OperationRepository.ListByIDs")
	defer span.End()

	var models []OperationModel
	result := r.db.WithContext(ctx).Where("id IN (?)", ids).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	operations := make([]*domain.Operation, 0, len(models))
	for i := range models {
		operation, err := r.mapToDomain(&models[i])
		if err != nil {
			return nil, err
		}
		operations = append(operations, operation)
	}

	batchCache := make(map[string]*domain.Operation, len(operations))
	batchIDCache := make(map[string]*domain.Operation, len(operations))
	for _, operation := range operations {
		batchCache[utils.StringsBuilder(cacheKeyOperationByProviderIDAndIdentifier, operation.ProviderID, ":", operation.Identifier)] = operation
		batchIDCache[utils.StringsBuilder(cacheKeyOperationByID, operation.ID)] = operation
	}
	r.cache.BatchSet(batchCache)
	r.cache.BatchSet(batchIDCache)
	return operations, nil
}

// GetByProviderIDAndIdentifier returns an operation by provider ID and identifier
func (r *OperationRepository) GetByProviderIDAndIdentifier(ctx context.Context, providerID, identifier string) (*domain.Operation, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "OperationRepository.GetByProviderIDAndIdentifier")
	defer span.End()

	cachekey := utils.StringsBuilder(cacheKeyOperationByProviderIDAndIdentifier, providerID, ":", identifier)
	if operation, ok := r.cache.Get(cachekey); ok {
		return operation, nil
	}

	var model OperationModel
	result := r.db.WithContext(ctx).Where("provider_id = ? AND identifier = ?", providerID, identifier).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	operation, err := r.mapToDomain(&model)
	if err != nil {
		return nil, err
	}

	r.cache.Set(cachekey, operation)
	return operation, nil
}

// Create creates a new operation
func (r *OperationRepository) Create(ctx context.Context, operation *domain.Operation) error {
	ctx, span := r.obs.Tracer.Start(ctx, "OperationRepository.Create")
	defer span.End()

	model, err := r.mapToModel(operation)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Create(model)
	return result.Error
}

// Update updates an operation
func (r *OperationRepository) Update(ctx context.Context, operation *domain.Operation) error {
	ctx, span := r.obs.Tracer.Start(ctx, "OperationRepository.Update")
	defer span.End()

	model, err := r.mapToModel(operation)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Save(model)
	return result.Error
}

// Delete deletes an operation
func (r *OperationRepository) Delete(ctx context.Context, id string) error {
	ctx, span := r.obs.Tracer.Start(ctx, "OperationRepository.Delete")
	defer span.End()

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&OperationModel{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("operation not found with id: %s", id)
	}

	return nil
}

// mapToDomain maps an operation model to a domain operation
func (r *OperationRepository) mapToDomain(model *OperationModel) (*domain.Operation, error) {
	var jsonAttributes struct {
		RequiredPermissions []types.Permission `json:"required_permissions"`
		Parameters          []domain.Parameter `json:"parameters"`
	}

	if err := sonic.Unmarshal(model.JSONAttributes, &jsonAttributes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal operation info metadata: %w", err)
	}

	return &domain.Operation{
		ID:                  model.ID,
		Identifier:          model.Identifier,
		ProviderID:          model.ProviderID,
		Name:                model.Name,
		Description:         model.Description,
		Category:            model.Category,
		RequiredPermissions: jsonAttributes.RequiredPermissions,
		Parameters:          jsonAttributes.Parameters,
		CreatedAt:           model.CreatedAt,
		UpdatedAt:           model.UpdatedAt,
		DeletedAt:           parseGormDeletedAt(model.DeletedAt),
	}, nil
}

// mapToModel maps a domain operation to an operation model
func (r *OperationRepository) mapToModel(operation *domain.Operation) (*OperationModel, error) {
	jsonAttributes := struct {
		RequiredPermissions []types.Permission `json:"required_permissions"`
		Parameters          []domain.Parameter `json:"parameters"`
	}{
		RequiredPermissions: operation.RequiredPermissions,
		Parameters:          operation.Parameters,
	}

	jsonAttributesJSON, err := sonic.Marshal(jsonAttributes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operation info metadata: %w", err)
	}

	return &OperationModel{
		ID:             operation.ID,
		Identifier:     operation.Identifier,
		ProviderID:     operation.ProviderID,
		Name:           operation.Name,
		Description:    operation.Description,
		Category:       operation.Category,
		JSONAttributes: jsonAttributesJSON,
		CreatedAt:      operation.CreatedAt,
		UpdatedAt:      operation.UpdatedAt,
		DeletedAt:      parseDomainDeletedAt(operation.DeletedAt),
	}, nil
}
