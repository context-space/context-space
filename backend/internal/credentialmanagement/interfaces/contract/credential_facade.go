package contract

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/application"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	contractCredential "github.com/context-space/context-space/backend/internal/shared/contract/credentialmanagement"
	"go.uber.org/zap"
)

// CredentialContractFacade implements the Contract interface at the Interface layer
// Responsibility: Calls the Application layer and handles Domain to DTO conversion
type CredentialContractFacade struct {
	credentialService *application.CredentialService
	credentialFactory domain.CredentialFactory
	tokenRefresh      domain.TokenRefresh
	obs               *observability.ObservabilityProvider
}

// Ensure implementation of contract interface
var _ contractCredential.CredentialManagementContract = (*CredentialContractFacade)(nil)

// NewCredentialContractFacade creates a new Contract Facade
func NewCredentialContractFacade(
	credentialService *application.CredentialService,
	credentialFactory domain.CredentialFactory,
	tokenRefresh domain.TokenRefresh,
	obs *observability.ObservabilityProvider,
) contractCredential.CredentialManagementContract {
	return &CredentialContractFacade{
		credentialService: credentialService,
		credentialFactory: credentialFactory,
		tokenRefresh:      tokenRefresh,
		obs:               obs,
	}
}

// GetCredentialByUserAndProviderContract gets the credential for the user and provider and converts it to contract response
func (f *CredentialContractFacade) GetCredentialByUserAndProviderContract(ctx context.Context, userID, providerIdentifier string) (interface{}, error) {
	ctx, span := f.obs.Tracer.Start(ctx, "CredentialContractFacade.GetCredentialByUserAndProviderContract")
	defer span.End()

	f.obs.Logger.Debug(ctx, "Getting credential through contract facade",
		zap.String("user_id", userID),
		zap.String("provider_identifier", providerIdentifier))

	// Call application service
	credential, err := f.credentialService.GetCredentialByUserAndProvider(ctx, userID, providerIdentifier)
	if err != nil {
		f.obs.Logger.Error(ctx, "Failed to get credential from service",
			zap.String("user_id", userID),
			zap.String("provider_identifier", providerIdentifier),
			zap.Error(err))

		return nil, err
	}

	return credential, nil
}

// CreateNoneCredentialContract creates a none credential and converts it to contract DTO
func (f *CredentialContractFacade) CreateNoneCredentialContract(ctx context.Context, userID, providerIdentifier string) (*contractCredential.CredentialDTO, error) {
	ctx, span := f.obs.Tracer.Start(ctx, "CredentialContractFacade.CreateNoneCredentialContract")
	defer span.End()

	f.obs.Logger.Debug(ctx, "Creating none credential through contract facade",
		zap.String("user_id", userID),
		zap.String("provider_identifier", providerIdentifier))

	// Call credential factory
	noneCredential, err := f.credentialFactory.CreateNone(ctx, userID, providerIdentifier)
	if err != nil {
		f.obs.Logger.Error(ctx, "Failed to create none credential",
			zap.String("user_id", userID),
			zap.String("provider_identifier", providerIdentifier),
			zap.Error(err))
		return nil, err
	}

	f.obs.Logger.Debug(ctx, "Successfully created none credential through contract facade",
		zap.String("user_id", userID),
		zap.String("provider_identifier", providerIdentifier))

	return noneCredential, nil
}

// UpdateCredentialLastUsedAtContract updates the last used time of the credential
func (f *CredentialContractFacade) UpdateCredentialLastUsedAtContract(ctx context.Context, credential interface{}) error {
	ctx, span := f.obs.Tracer.Start(ctx, "CredentialContractFacade.UpdateCredentialLastUsedAtContract")
	defer span.End()

	f.obs.Logger.Debug(ctx, "Updating credential last used time through contract facade")

	// Call credential factory
	err := f.credentialFactory.UpdateCredentialLastUsedAt(ctx, credential)
	if err != nil {
		f.obs.Logger.Error(ctx, "Failed to update credential last used time",
			zap.Error(err))
		return err
	}

	f.obs.Logger.Debug(ctx, "Successfully updated credential last used time through contract facade")
	return nil
}

func (f *CredentialContractFacade) RefreshAccessTokenContract(ctx context.Context, providerIdentifier string, credential interface{}) (interface{}, error) {
	ctx, span := f.obs.Tracer.Start(ctx, "CredentialContractFacade.RefreshAccessTokenContract")
	defer span.End()

	return f.tokenRefresh.RefreshAccessTokenIfNeeded(ctx, providerIdentifier, credential)
}
