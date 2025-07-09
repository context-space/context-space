package application

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/providercore/domain"
	"github.com/context-space/context-space/backend/internal/shared/apierrors"
	"github.com/context-space/context-space/backend/internal/shared/events"
	"go.uber.org/zap"
)

// ProviderEventType defines the events related to providers
type ProviderEventType string

const (
	// ProviderCreatedEvent is emitted when a provider is created
	ProviderCreatedEvent ProviderEventType = "provider.created"
	// ProviderUpdatedEvent is emitted when a provider is updated
	ProviderUpdatedEvent ProviderEventType = "provider.updated"
	// ProviderDeletedEvent is emitted when a provider is deleted
	ProviderDeletedEvent ProviderEventType = "provider.deleted"
	// OperationCreatedEvent is emitted when an operation is created
	OperationCreatedEvent ProviderEventType = "operation.created"
	// OperationUpdatedEvent is emitted when an operation is updated
	OperationUpdatedEvent ProviderEventType = "operation.updated"
	// OperationDeletedEvent is emitted when an operation is deleted
	OperationDeletedEvent ProviderEventType = "operation.deleted"
)

// ProviderService provides provider-related application services
type ProviderService struct {
	providerRepo  domain.ProviderRepository
	operationRepo domain.OperationRepository
	eventBus      *events.Bus
	obs           *observability.ObservabilityProvider
}

// NewProviderService creates a new ProviderService
func NewProviderService(
	providerRepo domain.ProviderRepository,
	operationRepo domain.OperationRepository,
	eventBus *events.Bus,
	observabilityProvider *observability.ObservabilityProvider,
) *ProviderService {
	return &ProviderService{
		providerRepo:  providerRepo,
		operationRepo: operationRepo,
		eventBus:      eventBus,
		obs:           observabilityProvider,
	}
}

// GetProviderByID retrieves a provider by ID
func (s *ProviderService) GetProviderByID(ctx context.Context, id string) (*domain.Provider, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.GetProviderByID")
	defer span.End()

	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}
	if provider == nil {
		return nil, apierrors.NewNotFoundError("", err)
	}

	return provider, nil
}

// GetProviderByIdentifier retrieves a provider by identifier
func (s *ProviderService) GetProviderByIdentifier(ctx context.Context, identifier string) (*domain.Provider, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.GetProviderByIdentifier")
	defer span.End()

	provider, err := s.providerRepo.GetByIdentifier(ctx, identifier)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}
	if provider == nil {
		return nil, apierrors.NewNotFoundError("", err)
	}

	return provider, nil
}

// ListProviders retrieves all providers
func (s *ProviderService) ListProviders(ctx context.Context) ([]*domain.Provider, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.ListProviders")
	defer span.End()

	providers, err := s.providerRepo.List(ctx)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}

	return providers, nil
}

// CreateProvider creates a new provider
func (s *ProviderService) CreateProvider(
	ctx context.Context,
	identifier, name, description, iconURL string,
	authType domain.ProviderAuthType,
	status domain.ProviderStatus,
	categories []string,
	operations []domain.Operation,
	permissions []domain.Permission,
) (*domain.Provider, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.CreateProvider")
	defer span.End()

	// Create new provider
	provider := domain.NewProvider(identifier, name, description, authType, status, iconURL, categories, permissions, operations)

	// Save provider
	if err := s.providerRepo.Create(ctx, provider); err != nil {
		return nil, apierrors.NewInternalError("", err)
	}

	// Emit provider created event
	s.emitProviderEvent(ctx, ProviderCreatedEvent, provider)

	return provider, nil
}

// UpdateProvider updates an existing provider
func (s *ProviderService) UpdateProvider(ctx context.Context, provider *domain.Provider) error {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.UpdateProvider")
	defer span.End()

	// Update provider
	if err := s.providerRepo.Update(ctx, provider); err != nil {
		return apierrors.NewInternalError("", err)
	}

	// Emit provider updated event
	s.emitProviderEvent(ctx, ProviderUpdatedEvent, provider)

	return nil
}

// DeleteProvider deletes a provider
func (s *ProviderService) DeleteProvider(ctx context.Context, id string) error {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.DeleteProvider")
	defer span.End()

	// Get existing provider to emit event
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return apierrors.NewInternalError("", err)
	}
	if provider == nil {
		return apierrors.NewNotFoundError("", err)
	}

	// Delete provider
	if err := s.providerRepo.Delete(ctx, id); err != nil {
		return apierrors.NewInternalError("", err)
	}

	// Emit provider deleted event
	s.emitProviderEvent(ctx, ProviderDeletedEvent, provider)

	return nil
}

// ActivateProvider activates a provider
func (s *ProviderService) ActivateProvider(ctx context.Context, id string) error {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.ActivateProvider")
	defer span.End()

	// Get provider
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return apierrors.NewInternalError("", err)
	}
	if provider == nil {
		return apierrors.NewNotFoundError("", err)
	}

	// If provider is already active, do nothing
	if provider.IsActive() {
		return nil
	}

	// Activate provider
	provider.Activate()

	// Update provider
	if err := s.providerRepo.Update(ctx, provider); err != nil {
		return apierrors.NewInternalError("", err)
	}

	// Emit provider updated event
	s.emitProviderEvent(ctx, ProviderUpdatedEvent, provider)

	return nil
}

// DeactivateProvider deactivates a provider
func (s *ProviderService) DeactivateProvider(ctx context.Context, id string) error {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.DeactivateProvider")
	defer span.End()

	// Get provider
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return apierrors.NewInternalError("", err)
	}
	if provider == nil {
		return apierrors.NewNotFoundError("", err)
	}

	// If provider is already inactive, do nothing
	if provider.IsInactive() {
		return nil
	}

	// Deactivate provider
	provider.Deactivate()

	// Update provider
	if err := s.providerRepo.Update(ctx, provider); err != nil {
		return apierrors.NewInternalError("", err)
	}

	// Emit provider updated event
	s.emitProviderEvent(ctx, ProviderUpdatedEvent, provider)

	return nil
}

// GetOperationsByProviderID retrieves operations for a provider
func (s *ProviderService) GetOperationsByProviderID(ctx context.Context, providerID string) ([]*domain.Operation, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.GetOperationsByProviderID")
	defer span.End()

	operations, err := s.operationRepo.ListByProviderID(ctx, providerID)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}

	return operations, nil
}

// CreateOperation creates a new operation for a provider
func (s *ProviderService) CreateOperation(
	ctx context.Context,
	identifier, providerID, name, description, category string,
	requiredPermissions []domain.Permission,
	parameters []domain.Parameter,
) (*domain.Operation, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.CreateOperation")
	defer span.End()

	// Check if provider exists
	provider, err := s.providerRepo.GetByID(ctx, providerID)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}
	if provider == nil {
		return nil, apierrors.NewNotFoundError("", err)
	}

	// Create operation
	operation := domain.NewOperation(identifier, providerID, name, description, category, requiredPermissions, parameters)

	// Save operation
	if err := s.operationRepo.Create(ctx, operation); err != nil {
		return nil, apierrors.NewInternalError("", err)
	}

	// Emit operation created event
	s.emitOperationEvent(ctx, OperationCreatedEvent, operation)

	return operation, nil
}

// UpdateOperation updates an existing operation
func (s *ProviderService) UpdateOperation(ctx context.Context, operation *domain.Operation) error {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.UpdateOperation")
	defer span.End()

	// Update operation
	if err := s.operationRepo.Update(ctx, operation); err != nil {
		return apierrors.NewInternalError("", err)
	}

	// Emit operation updated event
	s.emitOperationEvent(ctx, OperationUpdatedEvent, operation)

	return nil
}

// DeleteOperation deletes an operation
func (s *ProviderService) DeleteOperation(ctx context.Context, id string) error {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.DeleteOperation")
	defer span.End()

	// Get operation to emit event
	operations, err := s.operationRepo.ListByProviderID(ctx, id)
	if err != nil {
		return apierrors.NewInternalError("", err)
	}

	var operation *domain.Operation
	for _, op := range operations {
		if op.ID == id {
			operation = op
			break
		}
	}

	if operation == nil {
		return apierrors.NewNotFoundError("", nil)
	}

	// Delete operation
	if err := s.operationRepo.Delete(ctx, id); err != nil {
		return apierrors.NewInternalError("", err)
	}

	// Emit operation deleted event
	s.emitOperationEvent(ctx, OperationDeletedEvent, operation)

	return nil
}

// emitProviderEvent emits a provider-related event
func (s *ProviderService) emitProviderEvent(ctx context.Context, eventType ProviderEventType, provider *domain.Provider) {
	// Create event metadata
	metadata := events.Metadata{
		TraceID: observability.GetTraceID(ctx),
		SpanID:  observability.GetSpanID(ctx),
		Properties: map[string]string{
			"provider_id":         provider.ID,
			"provider_identifier": provider.Identifier,
			"provider_name":       provider.Name,
			"provider_status":     string(provider.Status),
			"provider_auth_type":  string(provider.AuthType),
		},
	}

	// Create and publish event
	event := events.NewEvent(events.EventType(eventType), provider, metadata)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.obs.Logger.Error(ctx, "Failed to publish provider event", zap.Error(err),
			zap.String("event_type", string(eventType)),
			zap.String("provider_id", provider.ID),
		)
	}
}

// emitOperationEvent emits an operation-related event
func (s *ProviderService) emitOperationEvent(ctx context.Context, eventType ProviderEventType, operation *domain.Operation) {
	// Create event metadata
	metadata := events.Metadata{
		TraceID: observability.GetTraceID(ctx),
		SpanID:  observability.GetSpanID(ctx),
		Properties: map[string]string{
			"operation_id":         operation.ID,
			"operation_identifier": operation.Identifier,
			"operation_name":       operation.Name,
			"provider_id":          operation.ProviderID,
		},
	}

	// Create and publish event
	event := events.NewEvent(events.EventType(eventType), operation, metadata)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.obs.Logger.Error(ctx, "Failed to publish operation event", zap.Error(err),
			zap.String("event_type", string(eventType)),
			zap.String("operation_id", operation.ID),
		)
	}
}
