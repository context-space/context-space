package acl

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/integration/domain"
	contractCredential "github.com/context-space/context-space/backend/internal/shared/contract/credentialmanagement"
	"go.uber.org/zap"
)

// CredentialACL implements integration/domain.CredentialProvider interface
// Acts as an Anti-Corruption Layer between integration and credentialmanagement modules
// Also implements TokenRefreshProvider interface for token refresh operations
type CredentialACL struct {
	credentialContract contractCredential.CredentialManagementContract
	obs                *observability.ObservabilityProvider
}

// Ensure CredentialACL implements both interfaces
var _ domain.CredentialProvider = (*CredentialACL)(nil)
var _ domain.TokenRefreshProvider = (*CredentialACL)(nil)

// NewCredentialACL creates a new credential ACL that implements both CredentialProvider and TokenRefreshProvider
func NewCredentialACL(
	credentialContract contractCredential.CredentialManagementContract,
	obs *observability.ObservabilityProvider,
) *CredentialACL {
	return &CredentialACL{
		credentialContract: credentialContract,
		obs:                obs,
	}
}

// GetCredentialByUserAndProvider retrieves a credential through the contract layer
func (acl *CredentialACL) GetCredentialByUserAndProvider(ctx context.Context, userID, providerIdentifier string) (interface{}, error) {
	ctx, span := acl.obs.Tracer.Start(ctx, "CredentialACL.GetCredentialByUserAndProvider")
	defer span.End()

	acl.obs.Logger.Debug(ctx, "Getting credential through ACL",
		zap.String("user_id", userID),
		zap.String("provider_identifier", providerIdentifier))

	// Call through contract layer - this is the key isolation point
	credential, err := acl.credentialContract.GetCredentialByUserAndProviderContract(ctx, userID, providerIdentifier)
	if err != nil {
		acl.obs.Logger.Error(ctx, "Failed to get credential through contract",
			zap.String("user_id", userID),
			zap.String("provider_identifier", providerIdentifier),
			zap.Error(err))
		return nil, err
	}

	acl.obs.Logger.Debug(ctx, "Successfully got credential through ACL",
		zap.String("user_id", userID),
		zap.String("provider_identifier", providerIdentifier))

	return credential, nil
}

// CreateNone creates a none credential through the contract layer
func (acl *CredentialACL) CreateNone(ctx context.Context, userID, providerIdentifier string) (*contractCredential.CredentialDTO, error) {
	ctx, span := acl.obs.Tracer.Start(ctx, "CredentialACL.CreateNone")
	defer span.End()

	acl.obs.Logger.Debug(ctx, "Creating none credential through ACL",
		zap.String("user_id", userID),
		zap.String("provider_identifier", providerIdentifier))

	// Call through contract layer
	credentialDTO, err := acl.credentialContract.CreateNoneCredentialContract(ctx, userID, providerIdentifier)
	if err != nil {
		acl.obs.Logger.Error(ctx, "Failed to create none credential through contract",
			zap.String("user_id", userID),
			zap.String("provider_identifier", providerIdentifier),
			zap.Error(err))
		return nil, err
	}

	acl.obs.Logger.Debug(ctx, "Successfully created none credential through ACL",
		zap.String("user_id", userID),
		zap.String("provider_identifier", providerIdentifier))

	// Return the DTO directly (contract already returns *CredentialDTO)
	return credentialDTO, nil
}

// UpdateCredentialLastUsedAt updates credential through the contract layer
func (acl *CredentialACL) UpdateCredentialLastUsedAt(ctx context.Context, credential interface{}) error {
	ctx, span := acl.obs.Tracer.Start(ctx, "CredentialACL.UpdateCredentialLastUsedAt")
	defer span.End()

	acl.obs.Logger.Debug(ctx, "Updating credential last used time through ACL")

	// Call through contract layer
	err := acl.credentialContract.UpdateCredentialLastUsedAtContract(ctx, credential)
	if err != nil {
		acl.obs.Logger.Error(ctx, "Failed to update credential last used time through contract",
			zap.Error(err))
		return err
	}

	acl.obs.Logger.Debug(ctx, "Successfully updated credential last used time through ACL")
	return nil
}

// RefreshAccessToken refreshes OAuth token through the contract layer
func (acl *CredentialACL) RefreshAccessToken(ctx context.Context, providerIdentifier string, credential interface{}) (interface{}, error) {
	ctx, span := acl.obs.Tracer.Start(ctx, "CredentialACL.RefreshAccessToken")
	defer span.End()

	acl.obs.Logger.Debug(ctx, "Refreshing access token through ACL",
		zap.String("provider_identifier", providerIdentifier))

	// Call through contract layer
	refreshedCredential, err := acl.credentialContract.RefreshAccessTokenContract(ctx, providerIdentifier, credential)
	if err != nil {
		acl.obs.Logger.Error(ctx, "Failed to refresh access token through contract",
			zap.String("provider_identifier", providerIdentifier),
			zap.Error(err))
		return nil, err
	}

	acl.obs.Logger.Debug(ctx, "Successfully refreshed access token through ACL",
		zap.String("provider_identifier", providerIdentifier))

	return refreshedCredential, nil
}
