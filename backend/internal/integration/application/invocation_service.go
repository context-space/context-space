package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/integration/domain"
	"github.com/context-space/context-space/backend/internal/shared/events"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/cache"
)

// InvocationEventTypes defines the event types for invocation events
type InvocationEventTypes struct {
	Started  events.EventType
	Success  events.EventType
	Failed   events.EventType
	Canceled events.EventType
}

// DefaultInvocationEventTypes returns the default invocation event types
func DefaultInvocationEventTypes() InvocationEventTypes {
	return InvocationEventTypes{
		Started:  events.EventType("invocation.started"),
		Success:  events.EventType("invocation.success"),
		Failed:   events.EventType("invocation.failed"),
		Canceled: events.EventType("invocation.canceled"),
	}
}

// Common errors
var (
	ErrProviderNotFound        = errors.New("provider not found")
	ErrProviderAdapterNotFound = errors.New("provider adapter not found")
	ErrOperationNotFound       = errors.New("operation not found")
	ErrInvalidParameters       = errors.New("invalid parameters")
	ErrCredentialNotFound      = errors.New("credential not found")
	ErrRateLimitExceeded       = errors.New("rate limit exceeded")
	ErrAdapterExecuteFailed    = errors.New("adapter execution failed")
	ErrInvocationNotFound      = errors.New("invocation not found")
)

// Common adapter error codes (duplicated from adapterDomain for safety)
const (
	ErrorCodeRateLimitExceeded    = "rate_limited"
	ErrorCodeAuthenticationFailed = "authentication_failed"
	ErrorCodeProviderError        = "provider_error"
	ErrorCodeInvalidInput         = "invalid_input"
	ErrorCodeOperationNotFound    = "operation_not_found"
)

// InvocationService handles invocation of provider operations
type InvocationService struct {
	providerProvider     domain.ProviderProvider
	adapterProvider      domain.AdapterProvider
	credProvider         domain.CredentialProvider
	invocationRepo       domain.InvocationRepository
	tokenRefreshProvider domain.TokenRefreshProvider
	eventBus             events.EventBus
	eventTypes           InvocationEventTypes
	obs                  *observability.ObservabilityProvider
	redisClient          cache.Cache
}

// NewInvocationService creates a new invocation service
func NewInvocationService(
	providerProvider domain.ProviderProvider,
	adapterProvider domain.AdapterProvider,
	credProvider domain.CredentialProvider,
	invocationRepo domain.InvocationRepository,
	eventBus events.EventBus,
	observabilityProvider *observability.ObservabilityProvider,
	redisClient cache.Cache,
	tokenRefreshProvider domain.TokenRefreshProvider,
) *InvocationService {
	return &InvocationService{
		providerProvider:     providerProvider,
		adapterProvider:      adapterProvider,
		credProvider:         credProvider,
		tokenRefreshProvider: tokenRefreshProvider,
		invocationRepo:       invocationRepo,
		eventBus:             eventBus,
		eventTypes:           DefaultInvocationEventTypes(),
		obs:                  observabilityProvider,
		redisClient:          redisClient,
	}
}

// InvokeOperation invokes an operation on a provider
func (s *InvocationService) InvokeOperation(
	ctx context.Context,
	userID string,
	providerIdentifier string,
	operationIdentifier string,
	params map[string]interface{},
) (*domain.Invocation, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "InvocationService.InvokeOperation")
	defer span.End()

	// Add tracing attributes
	span.SetAttributes(
		attribute.String("user_id", userID),
		attribute.String("provider_identifier", providerIdentifier),
		attribute.String("operation_identifier", operationIdentifier),
	)

	s.obs.Logger.Debug(ctx, "Invoking operation",
		zap.String("log_key", "invoke_operation"),
		zap.String("user_id", userID),
		zap.String("provider_identifier", providerIdentifier),
		zap.String("operation_identifier", operationIdentifier),
		zap.Any("params", params),
	)

	// Get the provider (not used directly but validates existence)
	_, err := s.providerProvider.GetProviderByIdentifier(ctx, providerIdentifier)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, err.Error())
	}

	// Get the provider adapter
	providerAdapter, err := s.adapterProvider.GetAdapterByProviderIdentifier(ctx, providerIdentifier)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrProviderAdapterNotFound, err.Error())
	}

	// Get credential for the provider (if needed)
	var credential interface{}
	adapterInfo := providerAdapter.GetAdapterInfoContract()
	authType := adapterInfo.AuthType
	s.obs.Logger.Debug(ctx, "Loaded provider adapter",
		zap.String("provider_identifier", adapterInfo.Identifier),
		zap.String("auth_type", authType),
	)

	if authType != "none" {
		credential, err = s.credProvider.GetCredentialByUserAndProvider(ctx, userID, providerIdentifier)
		if err != nil {
			return nil, fmt.Errorf("failed to get credential: %w", err)
		}
		if credential == nil {
			s.obs.Logger.Debug(ctx, "Credential not found", zap.String("provider_identifier", providerIdentifier), zap.Error(err))
			return nil, ErrCredentialNotFound
		}

		credential, err = s.tokenRefreshProvider.RefreshAccessToken(ctx, providerIdentifier, credential)
		if err != nil {
			s.obs.Logger.Error(ctx, "Failed to refresh access token",
				zap.String("provider_identifier", providerIdentifier),
				zap.String("user_id", userID),
				zap.Error(err))
			return nil, fmt.Errorf("failed to refresh access token: %w", err)

		}
	} else {
		credential, err = s.credProvider.CreateNone(ctx, userID, providerIdentifier)
		if err != nil {
			s.obs.Logger.Debug(ctx, "Failed to create none credential", zap.Error(err))
			return nil, fmt.Errorf("failed to create none credential: %w", err)
		}
	}

	// Create a unique ID for this invocation
	invocationID := uuid.New().String()

	// Create invocation record
	invocation := domain.NewInvocation(
		invocationID,
		userID,
		providerIdentifier,
		operationIdentifier,
		params,
	)

	// Set the invocation as started
	invocation.SetStarted()

	// Save initial invocation record
	if err := s.invocationRepo.Create(ctx, invocation); err != nil {
		s.obs.Logger.Debug(ctx, "Failed to create invocation record", zap.Error(err))
		return nil, fmt.Errorf("failed to create invocation record: %w", err)
	}

	// Emit started event
	s.emitInvocationEvent(ctx, s.eventTypes.Started, invocation)

	// Execute the operation directly
	result, execErr := providerAdapter.ExecuteContract(ctx, operationIdentifier, params, credential)

	if execErr != nil {
		errMsg := fmt.Sprintf("Failed to execute operation: %s", execErr.Error())
		s.obs.Logger.Debug(ctx, errMsg, zap.Error(execErr))
		s.handleInvocationError(ctx, invocation, errors.New(errMsg))
		return invocation, fmt.Errorf("%w: %s", ErrAdapterExecuteFailed, execErr)
	}

	// Update credential last used at
	if err := s.credProvider.UpdateCredentialLastUsedAt(ctx, credential); err != nil {
		s.obs.Logger.Error(ctx, "Failed to update credential last used at", zap.Error(err))
	}

	// Serialize result to JSON
	resultJSON, err := sonic.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to marshal result: %s", err.Error())
		s.handleInvocationError(ctx, invocation, errors.New(errMsg))
		return invocation, fmt.Errorf("failed to marshal result: %w", err)
	}

	// Update invocation record with success
	invocation.SetSuccess(resultJSON) // Duration is calculated internally
	if err := s.invocationRepo.Update(ctx, invocation); err != nil {
		s.obs.Logger.Error(ctx, fmt.Sprintf("Failed to update invocation:%+v", invocation), zap.Error(err))
	}

	// Emit success event
	s.emitInvocationEvent(ctx, s.eventTypes.Success, invocation)

	s.obs.Logger.Debug(ctx, "Invocation completed",
		zap.String("invocation_id", invocation.ID),
		zap.String("status", string(invocation.Status)),
	)

	return invocation, nil
}

// GetInvocationByID returns an invocation by ID
func (s *InvocationService) GetInvocationByID(ctx context.Context, id string) (*domain.Invocation, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "InvocationService.GetInvocationByID")
	defer span.End()

	invocation, err := s.invocationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get invocation: %w", err)
	}

	if invocation == nil {
		return nil, ErrInvocationNotFound
	}

	return invocation, nil
}

// ListInvocationsByUserID returns invocations for a user
func (s *InvocationService) ListInvocationsByUserID(
	ctx context.Context,
	userID string,
	limit, offset int,
) ([]*domain.Invocation, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "InvocationService.ListInvocationsByUserID")
	defer span.End()

	invocations, err := s.invocationRepo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list invocations: %w", err)
	}

	return invocations, nil
}

// CountInvocationsByUserID returns the count of invocations for a user
func (s *InvocationService) CountInvocationsByUserID(ctx context.Context, userID string) (int64, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "InvocationService.CountInvocationsByUserID")
	defer span.End()

	count, err := s.invocationRepo.CountByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count invocations: %w", err)
	}

	return count, nil
}

// handleInvocationError updates the invocation record with an error
func (s *InvocationService) handleInvocationError(ctx context.Context, invocation *domain.Invocation, err error) {
	// Update invocation record with error
	invocation.SetFailed(err.Error()) // Duration is calculated internally
	if updateErr := s.invocationRepo.Update(ctx, invocation); updateErr != nil {
		s.obs.Logger.Error(ctx, fmt.Sprintf("Failed to update invocation:%+v", invocation), zap.Error(updateErr))
	}

	// Emit failed event
	s.emitInvocationEvent(ctx, s.eventTypes.Failed, invocation)
}

// emitInvocationEvent emits an event for an invocation
func (s *InvocationService) emitInvocationEvent(ctx context.Context, eventType events.EventType, invocation *domain.Invocation) {
	// Extract trace information from the context
	spanCtx := trace.SpanContextFromContext(ctx)
	traceID := ""
	spanID := ""
	if spanCtx.IsValid() {
		traceID = spanCtx.TraceID().String()
		spanID = spanCtx.SpanID().String()
	}

	// Create metadata
	metadata := events.Metadata{
		UserID:             invocation.UserID,
		ProviderIdentifier: invocation.ProviderIdentifier,
		Operation:          invocation.OperationIdentifier,
		TraceID:            traceID,
		SpanID:             spanID,
		Properties: map[string]string{
			"invocation_id": invocation.ID,
			"status":        string(invocation.Status),
		},
	}

	// Create event
	event := events.NewEvent(eventType, invocation, metadata)

	// Publish event
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.obs.Logger.Error(ctx, "Failed to publish invocation event",
			zap.Error(err),
			zap.String("event_type", string(eventType)),
			zap.String("invocation_id", invocation.ID),
		)
	}
}
