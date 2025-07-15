package application

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/identityaccess/domain"
	"github.com/context-space/context-space/backend/internal/shared/apierrors"
	"github.com/context-space/context-space/backend/internal/shared/events"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"go.uber.org/zap"
)

// UserEventType defines the events related to users
type UserEventType string

const (
	// UserCreatedEvent is emitted when a user is created
	UserCreatedEvent UserEventType = "user.created"
	// UserUpdatedEvent is emitted when a user is updated
	UserUpdatedEvent UserEventType = "user.updated"
	// UserDeletedEvent is emitted when a user is deleted
	UserDeletedEvent UserEventType = "user.deleted"
	// APIKeyCreatedEvent is emitted when an API key is created
	APIKeyCreatedEvent UserEventType = "apikey.created"
	// APIKeyDeactivatedEvent is emitted when an API key is deactivated
	APIKeyDeactivatedEvent UserEventType = "apikey.deactivated"
	// APIKeyDeletedEvent is emitted when an API key is deleted
	APIKeyDeletedEvent UserEventType = "apikey.deleted"
)

const (
	maxAPIKeysPerUser = 3
)

// UserService provides user-related application services
type UserService struct {
	userRepo          domain.UserRepository
	userinfoRepo      domain.UserInfoRepository
	apiKeyRepo        domain.APIKeyRepository
	unitOfWorkFactory database.UnitOfWorkFactory
	eventBus          *events.Bus
	obs               *observability.ObservabilityProvider
}

// NewUserService creates a new UserService
func NewUserService(
	userRepo domain.UserRepository,
	userInfoRepo domain.UserInfoRepository,
	apiKeyRepo domain.APIKeyRepository,
	unitOfWorkFactory database.UnitOfWorkFactory,
	eventBus *events.Bus,
	observabilityProvider *observability.ObservabilityProvider,
) *UserService {
	return &UserService{
		userRepo:          userRepo,
		userinfoRepo:      userInfoRepo,
		apiKeyRepo:        apiKeyRepo,
		unitOfWorkFactory: unitOfWorkFactory,
		eventBus:          eventBus,
		obs:               observabilityProvider,
	}
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "UserService.GetUserByID")
	defer span.End()

	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}
	if user == nil {
		return nil, apierrors.NewNotFoundError("", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserBySupID(ctx context.Context, supID string) (*domain.User, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "UserService.GetUserBySupID")
	defer span.End()

	user, err := s.userRepo.GetBySupID(ctx, supID)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}
	if user == nil {
		return nil, apierrors.NewNotFoundError("", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "UserService.GetUserByEmail")
	defer span.End()

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}
	if user == nil {
		return nil, apierrors.NewNotFoundError("", err)
	}

	return user, nil
}

// ListUsers retrieves a list of users with pagination
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int) ([]*domain.User, int64, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "UserService.ListUsers")
	defer span.End()

	// Calculate offset
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// Get total count
	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, apierrors.NewInternalError("", err)
	}

	// Get users for current page
	users, err := s.userRepo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, apierrors.NewInternalError("", err)
	}

	return users, total, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	ctx, span := s.obs.Tracer.Start(ctx, "UserService.DeleteUser")
	defer span.End()

	// Get existing user to emit event
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return apierrors.NewInternalError("", err)
	}
	if user == nil {
		return apierrors.NewNotFoundError("", err)
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return apierrors.NewInternalError("", err)
	}

	// Emit user deleted event
	s.emitUserEvent(ctx, UserDeletedEvent, user)

	return nil
}

// CreateAPIKey creates a new API key for a user
func (s *UserService) CreateAPIKey(ctx context.Context, userID, name, description string) (*domain.APIKey, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "UserService.CreateAPIKey")
	defer span.End()

	// Get user
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}
	if user == nil {
		return nil, apierrors.NewNotFoundError("", err)
	}

	apiKeys, err := s.apiKeyRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}

	// Check if user has reached the maximum number of API keys
	if len(apiKeys) >= maxAPIKeysPerUser {
		return nil, apierrors.NewForbiddenError("Maximum number of API keys reached", nil)
	}

	// Create API key
	apiKey := domain.NewAPIKey(userID, name, description)

	// Save to database
	if err := s.apiKeyRepo.Create(ctx, apiKey); err != nil {
		return nil, apierrors.NewInternalError("", err)
	}

	// Emit API key created event
	s.emitAPIKeyEvent(ctx, APIKeyCreatedEvent, apiKey)

	return apiKey, nil
}

// GetAPIKeyByID retrieves an API key by ID
func (s *UserService) GetAPIKeyByID(ctx context.Context, id string) (*domain.APIKey, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "UserService.GetAPIKeyByID")
	defer span.End()

	apiKey, err := s.apiKeyRepo.Get(ctx, id)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}
	if apiKey == nil {
		return nil, apierrors.NewNotFoundError("", err)
	}

	return apiKey, nil
}

// DeleteAPIKey deletes an API key
func (s *UserService) DeleteAPIKey(ctx context.Context, id string) error {
	ctx, span := s.obs.Tracer.Start(ctx, "UserService.DeleteAPIKey")
	defer span.End()

	// Get API key to emit event
	apiKey, err := s.apiKeyRepo.Get(ctx, id)
	if err != nil {
		return apierrors.NewInternalError("", err)
	}
	if apiKey == nil {
		return apierrors.NewNotFoundError("", err)
	}

	// Delete API key
	if err := s.apiKeyRepo.Delete(ctx, id); err != nil {
		return apierrors.NewInternalError("", err)
	}

	// Emit API key deleted event
	s.emitAPIKeyEvent(ctx, APIKeyDeletedEvent, apiKey)

	return nil
}

// ValidateAPIKey validates an API key
func (s *UserService) ValidateAPIKey(ctx context.Context, keyValue string) (*domain.User, *domain.APIKey, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "UserService.ValidateAPIKey")
	defer span.End()

	// Get API key by value
	apiKey, err := s.apiKeyRepo.GetByKeyValue(ctx, keyValue)
	if err != nil {
		return nil, nil, apierrors.NewUnauthorizedError("", err)
	}

	if apiKey == nil {
		return nil, nil, apierrors.NewUnauthorizedError("api key is invalid", nil)
	}

	// Get associated user
	user, err := s.userRepo.Get(ctx, apiKey.UserID)
	if err != nil {
		return nil, nil, apierrors.NewInternalError("", err)
	}
	if user == nil {
		return nil, nil, apierrors.NewNotFoundError("", err)
	}

	// Update last used timestamp
	apiKey.UpdateLastUsed()
	if err := s.apiKeyRepo.Update(ctx, apiKey); err != nil {
		return nil, nil, apierrors.NewInternalError("", err)
	}

	return user, apiKey, nil
}

// ListAPIKeys retrieves API keys for a user with pagination
func (s *UserService) ListAPIKeys(ctx context.Context, userID string) ([]*domain.APIKey, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "UserService.ListAPIKeys")
	defer span.End()

	// List API keys
	apiKeys, err := s.apiKeyRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}

	return apiKeys, nil
}

// emitUserEvent emits a user-related event
func (s *UserService) emitUserEvent(ctx context.Context, eventType UserEventType, user *domain.User) {
	// Create event metadata
	metadata := events.Metadata{
		UserID:  user.ID,
		TraceID: observability.GetTraceID(ctx),
		SpanID:  observability.GetSpanID(ctx),
		Properties: map[string]string{
			"email": user.Email,
		},
	}

	// Create and publish event
	event := events.NewEvent(events.EventType(eventType), user, metadata)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.obs.Logger.Error(ctx, "Failed to publish user event", zap.Error(err),
			zap.String("event_type", string(eventType)),
			zap.String("user_id", user.ID),
		)
	}
}

// emitAPIKeyEvent emits an API key-related event
func (s *UserService) emitAPIKeyEvent(ctx context.Context, eventType UserEventType, apiKey *domain.APIKey) {
	// Create event metadata
	metadata := events.Metadata{
		UserID:  apiKey.UserID,
		TraceID: observability.GetTraceID(ctx),
		SpanID:  observability.GetSpanID(ctx),
		Properties: map[string]string{
			"api_key_id": apiKey.ID,
		},
	}

	// Create and publish event
	event := events.NewEvent(events.EventType(eventType), apiKey, metadata)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.obs.Logger.Error(ctx, "Failed to publish API key event", zap.Error(err),
			zap.String("event_type", string(eventType)),
			zap.String("api_key_id", apiKey.ID),
		)
	}
}
