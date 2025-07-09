package persistence

import (
	"context"
	"fmt"
	"sync"
	"time"

	observability "github.com/context-space/cloud-observability"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/cache"
	"github.com/context-space/context-space/backend/internal/shared/utils"
)

// TokenRefreshService manages automatic OAuth token refresh
type TokenRefreshService struct {
	redisClient         cache.Cache
	oauthProvider       domain.OAuthProvider
	credentialRepo      domain.CredentialRepository
	oauthCredentialRepo domain.OAuthCredentialRepository
	vaultService        domain.VaultService
	obs                 *observability.ObservabilityProvider
}

// NewTokenRefreshService creates a new token refresh service
func NewTokenRefreshService(
	redisClient cache.Cache,
	oauthProvider domain.OAuthProvider,
	credentialRepo domain.CredentialRepository,
	oauthCredentialRepo domain.OAuthCredentialRepository,
	vaultService domain.VaultService,
	obs *observability.ObservabilityProvider,
) *TokenRefreshService {
	return &TokenRefreshService{
		redisClient:         redisClient,
		oauthProvider:       oauthProvider,
		credentialRepo:      credentialRepo,
		oauthCredentialRepo: oauthCredentialRepo,
		vaultService:        vaultService,
		obs:                 obs,
	}
}

// RefreshToken begins the token refresh service
func (s *TokenRefreshService) RefreshAccessTokens(ctx context.Context) error {
	ctx, span := s.obs.Tracer.Start(ctx, "TokenRefreshService.RefreshAccessTokens")
	defer span.End()

	nowTime := time.Now()
	// Get expiry within 30 minutes oauth credentials
	oauthCredentials, err := s.oauthCredentialRepo.ListByExpiryWithin(ctx, nowTime.Add(30*time.Minute))
	if err != nil {
		return fmt.Errorf("failed to list OAuth credentials: %w", err)
	}

	credentialsID := make([]string, 0, len(oauthCredentials))
	oauthCredentialMap := make(map[string]*domain.OAuthCredential)
	for _, oauthCredential := range oauthCredentials {
		credentialsID = append(credentialsID, oauthCredential.ID)
		oauthCredentialMap[oauthCredential.ID] = oauthCredential
	}

	//Get credentials by credentialID
	credentials, err := s.credentialRepo.ListByID(ctx, credentialsID)
	if err != nil {
		return fmt.Errorf("failed to list credentials: %w", err)
	}

	//Get last used in 24 hours credentials
	providerRefreshMap := make(map[string][]*domain.OAuthCredential)
	for _, credential := range credentials {
		if credential.LastUsedAt.Compare(nowTime.Add(-24*time.Hour)) == -1 {
			continue
		}

		// check if the credential exists in oauthCredentialMap
		oauthCred, exists := oauthCredentialMap[credential.ID]
		if !exists {
			s.obs.Logger.Warn(ctx, "OAuth credential not found in map",
				zap.String("credential_id", credential.ID))
			continue
		}

		// decrypt token
		token := &oauth2.Token{}
		if err := s.vaultService.DecryptJSON(ctx, oauthCred.EncryptionMetadata, token); err != nil {
			s.obs.Logger.Error(ctx, utils.StringsBuilder("Failed to decrypt credential id: ", credential.ID, " metadata to token"), zap.Error(err))
			continue
		}

		oauthCred.Token = token
		oauthCred.UserID = credential.UserID
		oauthCred.ProviderIdentifier = credential.ProviderIdentifier

		providerRefreshMap[credential.ProviderIdentifier] = append(
			providerRefreshMap[credential.ProviderIdentifier],
			oauthCred,
		)
	}

	// process each provider's token refresh
	// use goroutines to process each provider, limit max 10 goroutines at the same time
	var wg sync.WaitGroup
	const maxConcurrentGoroutines = 10
	semaphore := make(chan struct{}, maxConcurrentGoroutines)

	for providerIdentifier, oauthCreds := range providerRefreshMap {
		providerIdentifierCopy := providerIdentifier
		oauthCredsCopy := oauthCreds
		wg.Add(1)
		go func(providerIdentifierCopy string, oauthCredsCopy []*domain.OAuthCredential) {
			defer wg.Done()

			// add panic recovery
			defer func() {
				if r := recover(); r != nil {
					s.obs.Logger.Error(ctx, "Panic recovered in token refresh goroutine",
						zap.String("provider", providerIdentifier),
						zap.Any("panic", r))
				}
			}()

			// get semaphore
			select {
			case semaphore <- struct{}{}:
				defer func() {
					// release semaphore
					<-semaphore
				}()
			case <-ctx.Done():
				s.obs.Logger.Warn(ctx, "Context cancelled before acquiring semaphore",
					zap.String("provider", providerIdentifier))
				return
			case <-time.After(30 * time.Second):
				s.obs.Logger.Error(ctx, "Timeout acquiring semaphore",
					zap.String("provider", providerIdentifier))
				return
			}

			s.refreshTokensForProvider(ctx, providerIdentifierCopy, oauthCredsCopy)
		}(providerIdentifierCopy, oauthCredsCopy)
	}

	wg.Wait()
	return nil
}

func (s *TokenRefreshService) RefreshAccessTokenIfNeeded(
	ctx context.Context,
	providerIdentifier string,
	credential interface{},
) (interface{}, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "TokenRefreshService.RefreshUserTokens")
	defer span.End()

	oauthCred, ok := credential.(*domain.OAuthCredential)
	if !ok { // not oauth credential, skip
		return credential, nil
	}
	if oauthCred == nil || oauthCred.Token == nil {
		return nil, fmt.Errorf("invalid or missing OAuth credential")
	}

	shouldRefresh, err := s.oauthProvider.ShouldRefreshToken(providerIdentifier, oauthCred.Token)
	if err != nil {
		return nil, err
	}
	if !shouldRefresh {
		return oauthCred, nil
	}

	// use core token refresh logic
	err = s.refreshTokenCore(ctx, providerIdentifier, oauthCred)
	if err != nil {
		return nil, err
	}

	return oauthCred, nil
}

// refreshTokensForProvider refresh all tokens for a provider
func (s *TokenRefreshService) refreshTokensForProvider(
	ctx context.Context,
	providerIdentifier string,
	oauthCreds []*domain.OAuthCredential,
) {
	providerCtx, providerSpan := s.obs.Tracer.Start(ctx, fmt.Sprintf("TokenRefreshProvider-%s", providerIdentifier))
	defer providerSpan.End()

	startTime := time.Now()
	s.obs.Logger.Info(providerCtx, "Processing tokens for provider",
		zap.String("provider", providerIdentifier),
		zap.Int("token_count", len(oauthCreds)))

	// call OAuth provider to refresh tokens
	var successCount, failCount int
	for _, oauthCred := range oauthCreds {
		// use distributed lock and retry mechanism to refresh single token
		err := s.refreshSingleTokenWithLockAndRetry(providerCtx, providerIdentifier, oauthCred)
		if err != nil {
			failCount++
			s.obs.Logger.Error(providerCtx, "Failed to refresh token",
				zap.String("provider", providerIdentifier),
				zap.String("credential_id", oauthCred.ID),
				zap.Error(err))
			continue
		}

		successCount++
		s.obs.Logger.Debug(providerCtx, "Token refreshed successfully",
			zap.String("provider", providerIdentifier),
			zap.String("credential_id", oauthCred.ID))

		// Sleep 100ms to avoid triggering rate limit
		time.Sleep(100 * time.Millisecond)
	}

	duration := time.Since(startTime)
	s.obs.Logger.Info(providerCtx, "Token refresh completed",
		zap.String("provider", providerIdentifier),
		zap.Int("success_count", successCount),
		zap.Int("fail_count", failCount),
		zap.Duration("duration", duration))
}

// refreshSingleTokenWithLockAndRetry use distributed lock and retry mechanism to refresh single token
func (s *TokenRefreshService) refreshSingleTokenWithLockAndRetry(
	ctx context.Context,
	providerID string,
	oauthCred *domain.OAuthCredential,
) error {
	const maxRetries = 5
	const retryDelay = 500 * time.Millisecond

	initialToken := oauthCred.Token.AccessToken

	lockKey := fmt.Sprintf(cache.AccessTokenLockKey, oauthCred.ProviderIdentifier, oauthCred.UserID)
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// try to get lock
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
			return fmt.Errorf("failed to acquire lock after %d attempts: %w", maxRetries, err)
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
			return fmt.Errorf("failed to acquire lock after %d attempts", maxRetries)
		}

		// when get lock successfully, execute token refresh
		defer func() {
			// release lock
			if err := s.redisClient.ReleaseLock(ctx, lockKey); err != nil {
				s.obs.Logger.Error(ctx, "Failed to release lock",
					zap.String("credential_id", oauthCred.ID),
					zap.Error(err))
			}
		}()

		// get latest credential information, check if token has been updated by another process
		latestCred, err := s.oauthCredentialRepo.GetByCredentialID(ctx, oauthCred.ID)
		if err != nil {
			return fmt.Errorf("failed to get latest credential: %w", err)
		}

		token := &oauth2.Token{}
		if err := s.vaultService.DecryptJSON(ctx, latestCred.EncryptionMetadata, token); err != nil {
			return fmt.Errorf("failed to decrypt credential: %w", err)
		}

		// check if token has been updated by another process
		if token.AccessToken != initialToken {
			s.obs.Logger.Info(ctx, "Token was updated by another process, skipping refresh",
				zap.String("credential_id", oauthCred.ID),
				zap.String("original_token", initialToken),
				zap.String("current_token", token.AccessToken))
			return nil
		}

		// execute token refresh
		return s.refreshTokenCore(ctx, providerID, oauthCred)
	}

	return fmt.Errorf("unexpected error: should not reach here") // should not reach here
}

// refreshTokenCore execute token refresh core logic
func (s *TokenRefreshService) refreshTokenCore(ctx context.Context, providerID string, oauthCred *domain.OAuthCredential) error {
	// execute token refresh
	newToken, err := s.oauthProvider.RefreshToken(ctx, providerID, oauthCred.Token)
	if err != nil {
		return fmt.Errorf("failed to refresh OAuth token: %w", err)
	}

	// update token
	oauthCred.Token = newToken

	metadata, err := s.vaultService.EncryptJSON(
		ctx,
		oauthCred.Token,
		oauthCred.EncryptionMetadata.Region,
		oauthCred.EncryptionMetadata.CredentialType,
	)
	if err != nil {
		return fmt.Errorf("failed to encrypt OAuth token: %w", err)
	}
	oauthCred.EncryptionMetadata = metadata
	oauthCred.UpdatedAt = time.Now()

	// update database
	err = s.oauthCredentialRepo.Update(ctx, oauthCred)
	if err != nil {
		return fmt.Errorf("failed to update OAuth credential: %w", err)
	}

	return nil
}
