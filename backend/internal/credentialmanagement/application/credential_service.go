package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	observability "github.com/context-space/cloud-observability"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/shared/events"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/cache"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
)

// Common errors
var (
	ErrCredentialNotFound      = errors.New("credential not found")
	ErrInvalidCredentialType   = errors.New("invalid credential type")
	ErrCredentialAlreadyExists = errors.New("credential already exists for this user and provider")
	ErrRefreshNotSupported     = errors.New("credential does not support refresh")
	ErrCredentialExpired       = errors.New("credential expired")
)

// CredentialEventTypes defines the event types for the credential service
type CredentialEventTypes struct {
	Created   events.EventType
	Updated   events.EventType
	Deleted   events.EventType
	Refreshed events.EventType
}

// DefaultCredentialEventTypes returns the default credential event types
func DefaultCredentialEventTypes() CredentialEventTypes {
	return CredentialEventTypes{
		Created:   events.EventType("credential.created"),
		Updated:   events.EventType("credential.updated"),
		Deleted:   events.EventType("credential.deleted"),
		Refreshed: events.EventType("credential.refreshed"),
	}
}

// CredentialService provides operations for managing credentials
type CredentialService struct {
	oAuthRedirectURL    string
	credentialRepo      domain.CredentialRepository
	credFactory         domain.CredentialFactory
	oauthProvider       domain.OAuthProvider
	eventBus            events.EventBus
	eventTypes          CredentialEventTypes
	obs                 *observability.ObservabilityProvider
	unitOfWorkFactory   database.UnitOfWorkFactory
	redisClient         cache.Cache
	tokenRefreshService domain.TokenRefresh
}

// NewCredentialService creates a new credential service
func NewCredentialService(
	credentialRepo domain.CredentialRepository,
	credFactory domain.CredentialFactory,
	eventBus events.EventBus,
	observabilityProvider *observability.ObservabilityProvider,
	oauthProvider domain.OAuthProvider,
	unitOfWorkFactory database.UnitOfWorkFactory,
	oAuthRedirectURL string,
	redisClient cache.Cache,
	tokenRefreshService domain.TokenRefresh,
) *CredentialService {
	return &CredentialService{
		credentialRepo:      credentialRepo,
		credFactory:         credFactory,
		eventBus:            eventBus,
		eventTypes:          DefaultCredentialEventTypes(),
		obs:                 observabilityProvider,
		oauthProvider:       oauthProvider,
		unitOfWorkFactory:   unitOfWorkFactory,
		oAuthRedirectURL:    oAuthRedirectURL,
		redisClient:         redisClient,
		tokenRefreshService: tokenRefreshService,
	}
}

// CreateOAuthCredential creates a new OAuth credential
func (s *CredentialService) CreateOAuthCredential(
	ctx context.Context,
	userID, providerIdentifier string,
	oauth2Token *oauth2.Token,
	scopes []string,
) (*domain.OAuthCredential, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "CredentialService.CreateOAuthCredential")
	defer span.End()

	// Start a new transaction
	unitOfWork := s.unitOfWorkFactory.Create()
	err := unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// Check if a credential already exists for this user and provider
	cred, err := s.credentialRepo.GetByUserAndProvider(ctx, userID, providerIdentifier)
	if err != nil {
		return nil, err
	}
	if cred != nil && cred.Type == domain.CredentialTypeOAuth {
		if err := s.credentialRepo.Delete(ctx, cred.ID); err != nil {
			unitOfWork.Rollback(ctx)
			return nil, err
		}
	}

	// Create the OAuth credential
	oauthCred, err := s.credFactory.CreateOAuth(ctx, userID, providerIdentifier, oauth2Token, scopes)
	if err != nil {
		unitOfWork.Rollback(ctx)
		return nil, err
	}

	// Emit credential created event
	event := events.NewEvent(
		s.eventTypes.Created,
		events.Payload{
			"credential_id":       oauthCred.ID,
			"user_id":             userID,
			"provider_identifier": providerIdentifier,
			"type":                string(oauthCred.Type),
		},
		events.Metadata{
			UserID:             userID,
			ProviderIdentifier: providerIdentifier,
			TraceID:            span.SpanContext().TraceID().String(),
			SpanID:             span.SpanContext().SpanID().String(),
		},
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		unitOfWork.Rollback(ctx)
		return nil, err
	}

	// Commit the transaction
	if err := unitOfWork.Commit(ctx); err != nil {
		return nil, err
	}

	return oauthCred, nil
}

// CreateAPIKeyCredential creates a new API key credential
func (s *CredentialService) CreateAPIKeyCredential(
	ctx context.Context,
	userID, providerIdentifier, apiKey string,
) (*domain.APIKeyCredential, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "CredentialService.CreateAPIKeyCredential")
	defer span.End()

	// Start a new transaction
	unitOfWork := s.unitOfWorkFactory.Create()
	err := unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// Check if a credential already exists for this user and provider
	cred, err := s.credentialRepo.GetByUserAndProvider(ctx, userID, providerIdentifier)
	if err != nil {
		return nil, err
	}
	if cred != nil && cred.Type == domain.CredentialTypeAPIKey {
		if err := s.credentialRepo.Delete(ctx, cred.ID); err != nil {
			unitOfWork.Rollback(ctx)
			return nil, err
		}
	}

	// Create the API key credential
	apiKeyCred, err := s.credFactory.CreateAPIKey(ctx, userID, providerIdentifier, apiKey)
	if err != nil {
		unitOfWork.Rollback(ctx)
		return nil, err
	}

	// Emit credential created event
	event := events.NewEvent(
		s.eventTypes.Created,
		events.Payload{
			"credential_id":       apiKeyCred.ID,
			"user_id":             userID,
			"provider_identifier": providerIdentifier,
			"type":                string(apiKeyCred.Type),
		},
		events.Metadata{
			UserID:             userID,
			ProviderIdentifier: providerIdentifier,
			TraceID:            span.SpanContext().TraceID().String(),
			SpanID:             span.SpanContext().SpanID().String(),
		},
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		unitOfWork.Rollback(ctx)
		return nil, err
	}

	// Commit the transaction
	if err := unitOfWork.Commit(ctx); err != nil {
		return nil, err
	}

	return apiKeyCred, nil
}

// GetCredential retrieves a credential by ID
func (s *CredentialService) GetCredential(ctx context.Context, id string) (interface{}, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "CredentialService.GetCredential")
	defer span.End()

	// Get the credential from the factory
	cred, err := s.credFactory.GetCredential(ctx, id)
	if err != nil {
		return nil, err
	}

	if cred == nil {
		return nil, ErrCredentialNotFound
	}

	return cred, nil
}

// GetCredentialByUserAndProvider retrieves a credential by user ID and provider identifier
func (s *CredentialService) GetCredentialByUserAndProvider(ctx context.Context, userID, providerIdentifier string) (interface{}, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "CredentialService.GetCredentialByUserAndProvider")
	defer span.End()

	// lock the credential prevent refresh token when getting credential
	lockKey := fmt.Sprintf(cache.AccessTokenLockKey, providerIdentifier, userID)

	// use retry mechanism to get lock
	const maxRetries = 5
	const retryDelay = 100 * time.Millisecond

	for attempt := 1; attempt <= maxRetries; attempt++ {
		lock, err := s.redisClient.AcquireLock(ctx, lockKey, cache.AccessTokenLockTimeout)
		if err != nil {
			// when get lock failed, log and retry
			s.obs.Logger.Warn(ctx, "Failed to acquire lock, retrying",
				zap.String("lock_key", lockKey),
				zap.Int("attempt", attempt),
				zap.Error(err))
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("%w : %s", ErrCredentialExpired, err.Error())
		}

		if !lock {
			// when lock is held by another process, retry
			s.obs.Logger.Debug(ctx, "Lock held by another process, retrying",
				zap.String("lock_key", lockKey),
				zap.Int("attempt", attempt))
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("%w : %s", ErrCredentialExpired, "please re-authorize")
		}

		// when get lock successfully, break the retry loop
		break
	}

	defer s.redisClient.ReleaseLock(ctx, lockKey)

	// Get the credential from the factory
	cred, err := s.credFactory.GetCredentialByUserAndProvider(ctx, userID, providerIdentifier)
	if err != nil {
		return nil, err
	}

	if cred == nil {
		return nil, ErrCredentialNotFound
	}

	cred, err = s.tokenRefreshService.RefreshAccessTokenIfNeeded(ctx, providerIdentifier, cred)
	if err != nil {
		s.obs.Logger.Error(ctx, "Failed to refresh access token",
			zap.String("provider_identifier", providerIdentifier),
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, fmt.Errorf("%w: %s", ErrCredentialExpired, err.Error())
	}

	return cred, nil
}

// DeleteCredential deletes a credential
func (s *CredentialService) DeleteCredential(ctx context.Context, id string) error {
	ctx, span := s.obs.Tracer.Start(ctx, "CredentialService.DeleteCredential")
	defer span.End()

	// Get the credential to determine its type and user/provider IDs
	baseCred, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if baseCred == nil {
		return ErrCredentialNotFound
	}

	// Delete the base credential
	if err := s.credentialRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Emit credential deleted event
	event := events.NewEvent(
		s.eventTypes.Deleted,
		events.Payload{
			"credential_id":       id,
			"user_id":             baseCred.UserID,
			"provider_identifier": baseCred.ProviderIdentifier,
			"type":                string(baseCred.Type),
		},
		events.Metadata{
			UserID:             baseCred.UserID,
			ProviderIdentifier: baseCred.ProviderIdentifier,
			TraceID:            span.SpanContext().TraceID().String(),
			SpanID:             span.SpanContext().SpanID().String(),
		},
	)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		return err
	}

	return nil
}

// GetAllCredentialsByUser retrieves all credentials for a user
func (s *CredentialService) GetAllCredentialsByUser(ctx context.Context, userID string) ([]*domain.Credential, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "CredentialService.GetAllCredentialsByUser")
	defer span.End()

	// Get all base credentials for the user
	baseCreds, err := s.credentialRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return baseCreds, nil
}

// GetOAuthURL generates an OAuth authorization URL for a provider
func (s *CredentialService) GetOAuthURL(ctx context.Context, providerIdentifier, state, codeChallenge string, permissionIdentifiers []string) (string, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "CredentialService.GetOAuthURL")
	defer span.End()

	scopes, err := s.oauthProvider.GetScopesFromPermissions(ctx, providerIdentifier, permissionIdentifiers)
	if err != nil {
		return "", err
	}

	oauthURL, err := s.oauthProvider.GenerateOAuthURL(ctx, providerIdentifier, s.oAuthRedirectURL, state, codeChallenge, scopes)
	if err != nil {
		return "", err
	}

	return oauthURL, nil
}

// HandleOAuthCallback processes an OAuth callback and stores the credentials
func (s *CredentialService) HandleOAuthCallback(ctx context.Context, code, providerIdentifier, userID string, permissions []string, codeVerifier string) (*domain.OAuthCredential, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "CredentialService.HandleOAuthCallback")
	defer span.End()

	token, err := s.oauthProvider.ExchangeCodeForToken(ctx, providerIdentifier, code, s.oAuthRedirectURL, codeVerifier)
	s.obs.Logger.Debug(ctx, "ExchangeCodeForToken",
		zap.String("provider_identifier", providerIdentifier),
		zap.String("user_id", userID),
		zap.Any("token", token),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for tokens: %w", err)
	}

	scopes, err := s.oauthProvider.GetScopesFromPermissions(ctx, providerIdentifier, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to get scopes from permissions: %w", err)
	}

	oauthCred, err := s.CreateOAuthCredential(ctx, userID, providerIdentifier, token, scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth credential: %w", err)
	}

	return oauthCred, nil
}

// GetPermissionIdentifiersFromScopes gets the permission identifiers from the scopes
func (s *CredentialService) GetPermissionIdentifiersFromScopes(ctx context.Context, providerIdentifier string, scopes []string) ([]string, error) {
	return s.oauthProvider.GetPermissionIdentifiersFromScopes(ctx, providerIdentifier, scopes)
}

func (s *CredentialService) UpdateLastUsedAt(ctx context.Context, credential interface{}) error {
	switch cred := credential.(type) {
	case *domain.OAuthCredential:
		return s.credFactory.UpdateCredentialLastUsedAt(ctx, cred)
	case *domain.APIKeyCredential:
		return s.credFactory.UpdateCredentialLastUsedAt(ctx, cred)
	}

	return nil
}
