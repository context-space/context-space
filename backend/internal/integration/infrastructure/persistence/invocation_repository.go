package persistence

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/integration/domain"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"gorm.io/gorm"
)

// InvocationRepository implements the domain.InvocationRepository interface using GORM
type InvocationRepository struct {
	db  database.Database
	obs *observability.ObservabilityProvider
}

// NewInvocationRepository creates a new invocation repository
func NewInvocationRepository(db database.Database, observabilityProvider *observability.ObservabilityProvider) *InvocationRepository {
	return &InvocationRepository{
		db:  db,
		obs: observabilityProvider,
	}
}

// Create creates a new invocation
func (r *InvocationRepository) Create(ctx context.Context, invocation *domain.Invocation) error {
	ctx, span := r.obs.Tracer.Start(ctx, "InvocationRepository.Create")
	defer span.End()

	model, err := r.mapToModel(invocation)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Create(model)
	return result.Error
}

// Update updates an invocation
func (r *InvocationRepository) Update(ctx context.Context, invocation *domain.Invocation) error {
	ctx, span := r.obs.Tracer.Start(ctx, "InvocationRepository.Update")
	defer span.End()

	model, err := r.mapToModel(invocation)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Save(model)
	return result.Error
}

// GetByID returns an invocation by ID
func (r *InvocationRepository) GetByID(ctx context.Context, id string) (*domain.Invocation, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "InvocationRepository.GetByID")
	defer span.End()

	var model InvocationModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.mapToDomain(&model)
}

// ListByUserID returns invocations by user ID
func (r *InvocationRepository) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Invocation, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "InvocationRepository.ListByUserID")
	defer span.End()

	var models []InvocationModel
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	invocations := make([]*domain.Invocation, 0, len(models))
	for i := range models {
		invocation, err := r.mapToDomain(&models[i])
		if err != nil {
			return nil, err
		}
		invocations = append(invocations, invocation)
	}

	return invocations, nil
}

// CountByUserID returns the count of invocations by user ID
func (r *InvocationRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "InvocationRepository.CountByUserID")
	defer span.End()

	var count int64
	result := r.db.WithContext(ctx).Model(&InvocationModel{}).Where("user_id = ?", userID).Count(&count)
	return count, result.Error
}

// CountByProviderIdentifier returns the count of invocations by provider identifier
func (r *InvocationRepository) CountByProviderIdentifier(ctx context.Context, providerIdentifier string) (int64, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "InvocationRepository.CountByProviderIdentifier")
	defer span.End()

	var count int64
	result := r.db.WithContext(ctx).Model(&InvocationModel{}).Where("provider_identifier = ?", providerIdentifier).Count(&count)
	return count, result.Error
}

// CountByOperationIdentifier returns the count of invocations by operation identifier
func (r *InvocationRepository) CountByOperationIdentifier(ctx context.Context, providerIdentifier, operationIdentifier string) (int64, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "InvocationRepository.CountByOperationIdentifier")
	defer span.End()

	var count int64
	result := r.db.WithContext(ctx).Model(&InvocationModel{}).
		Where("provider_identifier = ? AND operation_identifier = ?", providerIdentifier, operationIdentifier).
		Count(&count)
	return count, result.Error
}

// mapToDomain converts an invocation model to a domain invocation
func (r *InvocationRepository) mapToDomain(model *InvocationModel) (*domain.Invocation, error) {
	var jsonAttributes struct {
		Parameters   string `json:"parameters"`    // Base64 encoded string
		ResponseData string `json:"response_data"` // Base64 encoded string
		ErrorMessage string `json:"error_message"` // Base64 encoded string
	}

	if err := sonic.Unmarshal(model.JSONAttributes, &jsonAttributes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal invocation JSON attributes: %w", err)
	}

	parameters, err := base64.StdEncoding.DecodeString(jsonAttributes.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to decode parameters: %w", err)
	}

	var parametersMap map[string]interface{}
	if err := sonic.Unmarshal(parameters, &parametersMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameters: %w", err)
	}

	responseData, err := base64.StdEncoding.DecodeString(jsonAttributes.ResponseData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response data: %w", err)
	}

	errorMessage, err := base64.StdEncoding.DecodeString(jsonAttributes.ErrorMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to decode error message: %w", err)
	}

	return &domain.Invocation{
		ID:                  model.ID,
		UserID:              model.UserID,
		ProviderIdentifier:  model.ProviderIdentifier,
		OperationIdentifier: model.OperationIdentifier,
		Status:              domain.InvocationStatus(model.Status),
		Duration:            model.Duration,
		StartedAt:           model.StartedAt,
		CompletedAt:         model.CompletedAt,
		ErrorMessage:        string(errorMessage),
		Parameters:          parametersMap,
		ResponseData:        responseData,
		CreatedAt:           model.CreatedAt,
		UpdatedAt:           model.UpdatedAt,
		DeletedAt:           parseGormDeletedAt(model.DeletedAt),
	}, nil
}

// mapToModel converts a domain invocation to an invocation model
func (r *InvocationRepository) mapToModel(invocation *domain.Invocation) (*InvocationModel, error) {
	parameters, err := sonic.Marshal(invocation.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters: %w", err)
	}

	responseData, err := sonic.Marshal(invocation.ResponseData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response data: %w", err)
	}

	jsonAttributes := struct {
		Parameters   string `json:"parameters"`    // Base64 encoded string
		ResponseData string `json:"response_data"` // Base64 encoded string
		ErrorMessage string `json:"error_message"` // Base64 encoded string
	}{
		Parameters:   base64.StdEncoding.EncodeToString(parameters),
		ResponseData: base64.StdEncoding.EncodeToString(responseData),
		ErrorMessage: base64.StdEncoding.EncodeToString([]byte(invocation.ErrorMessage)),
	}

	jsonAttributesBytes, err := sonic.Marshal(jsonAttributes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal invocation JSON attributes: %w", err)
	}

	return &InvocationModel{
		ID:                  invocation.ID,
		UserID:              invocation.UserID,
		ProviderIdentifier:  invocation.ProviderIdentifier,
		OperationIdentifier: invocation.OperationIdentifier,
		Status:              string(invocation.Status),
		Duration:            invocation.Duration,
		StartedAt:           invocation.StartedAt,
		CompletedAt:         invocation.CompletedAt,
		JSONAttributes:      jsonAttributesBytes,
		CreatedAt:           invocation.CreatedAt,
		UpdatedAt:           invocation.UpdatedAt,
		DeletedAt:           parseDomainDeletedAt(invocation.DeletedAt),
	}, nil
}
