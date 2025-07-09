package credential

import (
	"context"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
)

// CredentialProviderAdapter implements the integration.domain.CredentialProvider interface
// by adapting the credentialmanagement.domain.CredentialFactory
type CredentialProviderAdapter struct {
	credentialFactory credDomain.CredentialFactory
}

// NewCredentialProviderAdapter creates a new adapter for credential provider
func NewCredentialProviderAdapter(
	credentialFactory credDomain.CredentialFactory,
) *CredentialProviderAdapter {
	return &CredentialProviderAdapter{
		credentialFactory: credentialFactory,
	}
}

// GetCredentialByUserAndProvider delegates to the credential factory
func (a *CredentialProviderAdapter) GetCredentialByUserAndProvider(
	ctx context.Context,
	userID,
	providerIdentifier string,
) (interface{}, error) {
	// Simply delegate to the credential factory
	return a.credentialFactory.GetCredentialByUserAndProvider(ctx, userID, providerIdentifier)
}

// CreateNone delegates to the credential factory
func (a *CredentialProviderAdapter) CreateNone(ctx context.Context, userID, providerIdentifier string) (*credDomain.NoneCredential, error) {
	return a.credentialFactory.CreateNone(ctx, userID, providerIdentifier)
}

// UpdateCredentialLastUsedAt delegates to the credential factory
func (a *CredentialProviderAdapter) UpdateCredentialLastUsedAt(ctx context.Context, credential interface{}) error {
	return a.credentialFactory.UpdateCredentialLastUsedAt(ctx, credential)
}
