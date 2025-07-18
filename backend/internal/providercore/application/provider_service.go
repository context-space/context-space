package application

import (
	"context"
	"sort"

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

// ProviderFilterParams represents the parameters for filtering providers
type ProviderFilterParams struct {
	Tag          string
	AuthType     domain.ProviderAuthType
	ProviderName string
	Page         int
	PageSize     int
	SortField    string
	SortOrder    string
}

// ProviderFilterResult represents the result of filtering providers
type ProviderFilterResult struct {
	Providers   []*domain.Provider
	TotalCount  int
	CurrentPage int
	PageSize    int
	TotalPages  int
	HasNext     bool
	HasPrev     bool
}

// FilterProviders filters providers based on various criteria with pagination and sorting
func (s *ProviderService) FilterProviders(ctx context.Context, params ProviderFilterParams) (*ProviderFilterResult, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderService.FilterProviders")
	defer span.End()

	// Get all providers first
	allProviders, err := s.providerRepo.List(ctx)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}

	// Apply filters
	var filteredProviders []*domain.Provider

	// Filter by tag if specified
	if params.Tag != "" {
		// This will be implemented in step 7
		filteredProviders = s.filterByTag(allProviders, params.Tag)
	} else {
		filteredProviders = allProviders
	}

	// Filter by auth_type and provider_identifier in memory
	filteredProviders = s.applyInMemoryFilters(filteredProviders, params)

	// Calculate total count before pagination
	totalCount := len(filteredProviders)

	// Apply sorting
	s.applySorting(filteredProviders, params.SortField, params.SortOrder)

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	if offset < 0 {
		offset = 0
	}

	end := offset + params.PageSize
	if end > len(filteredProviders) {
		end = len(filteredProviders)
	}

	var paginatedProviders []*domain.Provider
	if offset < len(filteredProviders) {
		paginatedProviders = filteredProviders[offset:end]
	} else {
		paginatedProviders = []*domain.Provider{}
	}

	// Calculate pagination metadata
	totalPages := (totalCount + params.PageSize - 1) / params.PageSize
	if totalPages == 0 {
		totalPages = 1
	}

	hasNext := params.Page < totalPages
	hasPrev := params.Page > 1

	return &ProviderFilterResult{
		Providers:   paginatedProviders,
		TotalCount:  totalCount,
		CurrentPage: params.Page,
		PageSize:    params.PageSize,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}, nil
}

// filterByTag filters providers by tag using TagService
func (s *ProviderService) filterByTag(providers []*domain.Provider, tagName string) []*domain.Provider {
	var filtered []*domain.Provider
	for _, provider := range providers {
		if provider.HasTag(tagName) {
			filtered = append(filtered, provider)
		}
	}
	return filtered
}

// applyInMemoryFilters applies auth_type and provider_identifier filters in memory
func (s *ProviderService) applyInMemoryFilters(providers []*domain.Provider, params ProviderFilterParams) []*domain.Provider {
	var filtered []*domain.Provider

	for _, provider := range providers {
		// Skip deprecated providers
		if provider.Status == domain.ProviderStatusDeprecated {
			continue
		}

		// Filter by auth_type if specified
		if params.AuthType != "" && provider.AuthType != params.AuthType {
			continue
		}

		// Filter by provider_identifier if specified
		if params.ProviderName != "" && provider.Name != params.ProviderName {
			continue
		}

		filtered = append(filtered, provider)
	}

	return filtered
}

// applySorting applies sorting to the providers slice
func (s *ProviderService) applySorting(providers []*domain.Provider, sortField, sortOrder string) {
	switch sortField {
	case "created_at":
		if sortOrder == "asc" {
			sort.Slice(providers, func(i, j int) bool {
				return providers[i].CreatedAt.Before(providers[j].CreatedAt)
			})
		} else {
			sort.Slice(providers, func(i, j int) bool {
				return providers[i].CreatedAt.After(providers[j].CreatedAt)
			})
		}
	case "updated_at":
		if sortOrder == "asc" {
			sort.Slice(providers, func(i, j int) bool {
				return providers[i].UpdatedAt.Before(providers[j].UpdatedAt)
			})
		} else {
			sort.Slice(providers, func(i, j int) bool {
				return providers[i].UpdatedAt.After(providers[j].UpdatedAt)
			})
		}
	case "name":
		if sortOrder == "asc" {
			sort.Slice(providers, func(i, j int) bool {
				return providers[i].Name < providers[j].Name
			})
		} else {
			sort.Slice(providers, func(i, j int) bool {
				return providers[i].Name > providers[j].Name
			})
		}
	case "identifier":
		if sortOrder == "asc" {
			sort.Slice(providers, func(i, j int) bool {
				return providers[i].Identifier < providers[j].Identifier
			})
		} else {
			sort.Slice(providers, func(i, j int) bool {
				return providers[i].Identifier > providers[j].Identifier
			})
		}
	}
}
