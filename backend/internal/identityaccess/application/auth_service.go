package application

import (
	"context"
	"time"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/identityaccess/domain"
	"github.com/context-space/context-space/backend/internal/identityaccess/infrastructure/supabase"
	"github.com/context-space/context-space/backend/internal/shared/events"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"go.uber.org/zap"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo          domain.UserRepository
	userInfoRepo      domain.UserInfoRepository
	supabaseService   *supabase.SupabaseAuthService
	unitOfWorkFactory database.UnitOfWorkFactory
	eventBus          *events.Bus
	obs               *observability.ObservabilityProvider
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo domain.UserRepository,
	userInfoRepo domain.UserInfoRepository,
	supabaseService *supabase.SupabaseAuthService,
	unitOfWorkFactory database.UnitOfWorkFactory,
	eventBus *events.Bus,
	observabilityProvider *observability.ObservabilityProvider,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		userInfoRepo:      userInfoRepo,
		supabaseService:   supabaseService,
		unitOfWorkFactory: unitOfWorkFactory,
		eventBus:          eventBus,
		obs:               observabilityProvider,
	}
}

// ValidateToken validates an access token
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*supabase.Claims, error) {
	return s.supabaseService.ValidateToken(ctx, tokenString)
}

// FindOrCreateUser finds or creates a user from a supabase claims
func (s *AuthService) FindOrCreateUser(ctx context.Context, claims *supabase.Claims) (*domain.User, error) {
	return s.findOrCreateUser(ctx, claims)
}

// findOrCreateUser finds or creates a user based on Supabase user info
func (s *AuthService) findOrCreateUser(ctx context.Context, claims *supabase.Claims) (*domain.User, error) {
	user, err := s.userRepo.GetBySupID(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	// If found by supID, return user
	if user != nil {
		return user, nil
	}

	// User not found, create new user
	supUserInfo, err := s.supabaseService.GetUserInfo(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	newUser := domain.NewUser(supUserInfo.ID, supUserInfo.Email, supUserInfo.IsAnonymous)

	// Begin transaction
	unitOfWork := s.unitOfWorkFactory.Create()
	if err := unitOfWork.Begin(ctx); err != nil {
		return nil, err
	}

	// Save new user
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		unitOfWork.Rollback(ctx)
		return nil, err
	}

	// Store user info
	newUserInfo := domain.NewUserInfo(newUser.ID, supUserInfo.InfoMetadata)
	if err := s.userInfoRepo.Create(ctx, newUserInfo); err != nil {
		unitOfWork.Rollback(ctx)
		return nil, err
	}

	// Commit transaction
	if err := unitOfWork.Commit(ctx); err != nil {
		unitOfWork.Rollback(ctx)
		return nil, err
	}

	// Emit user created event
	s.emitUserCreatedEvent(ctx, newUser)

	return newUser, nil
}

// emitUserCreatedEvent emits a user created event
func (s *AuthService) emitUserCreatedEvent(ctx context.Context, user *domain.User) {
	event := events.Event{
		Type:      "user.created",
		Timestamp: time.Now(),
		Payload: events.Payload{
			"user_id": user.ID,
			"email":   user.Email,
		},
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.obs.Logger.Warn(ctx, "Failed to publish user created event", zap.Error(err))
	}
}
