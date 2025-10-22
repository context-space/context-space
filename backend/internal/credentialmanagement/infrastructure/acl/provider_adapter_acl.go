package acl

import (
	"context"
	"fmt"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	contractAdapter "github.com/context-space/context-space/backend/internal/shared/contract/provideradapter"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type ProviderAdapterACL struct {
	contractReader contractAdapter.ProviderAdapterContract
	obs            *observability.ObservabilityProvider
}

// NewProviderAdapterACL creates a new provider adapter ACL
func NewProviderAdapterACL(
	contractReader contractAdapter.ProviderAdapterContract,
	obs *observability.ObservabilityProvider,
) domain.OAuthProvider {
	return &ProviderAdapterACL{
		contractReader: contractReader,
		obs:            obs,
	}
}

func (a *ProviderAdapterACL) GenerateOAuthURL(ctx context.Context, providerIdentifier, redirectURL, state, codeChallenge string, scopes []string) (string, error) {
	ctx, span := a.obs.Tracer.Start(ctx, "ProviderAdapterACL.GenerateOAuthURL")
	defer span.End()

	a.obs.Logger.Debug(ctx, "Generating OAuth URL",
		zap.String("provider_identifier", providerIdentifier),
		zap.String("redirect_url", redirectURL),
		zap.String("state", state),
		zap.String("code_challenge", codeChallenge),
		zap.Strings("scopes", scopes),
	)

	url, err := a.contractReader.GenerateOAuthURLContract(ctx, providerIdentifier, redirectURL, state, codeChallenge, scopes)
	if err != nil {
		return "", fmt.Errorf("failed to generate OAuth URL: %w", err)
	}

	return url, nil
}

func (a *ProviderAdapterACL) ExchangeCodeForToken(ctx context.Context, providerIdentifier, code, redirectURL, codeVerifier string) (*oauth2.Token, error) {
	ctx, span := a.obs.Tracer.Start(ctx, "ProviderAdapterACL.ExchangeCodeForToken")
	defer span.End()

	a.obs.Logger.Debug(ctx, "Exchanging code for token",
		zap.String("provider_identifier", providerIdentifier),
		zap.String("code", code),
		zap.String("redirect_url", redirectURL),
		zap.String("code_verifier", codeVerifier),
	)

	token, err := a.contractReader.ExchangeCodeForTokenContract(ctx, providerIdentifier, code, redirectURL, codeVerifier)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token, nil
}

func (a *ProviderAdapterACL) ShouldRefreshToken(providerIdentifier string, token *oauth2.Token) (bool, error) {
	shouldRefresh, err := a.contractReader.ShouldRefreshTokenContract(providerIdentifier, token)
	if err != nil {
		return false, fmt.Errorf("failed to check if token should be refreshed: %w", err)
	}

	return shouldRefresh, nil
}

func (a *ProviderAdapterACL) RefreshToken(ctx context.Context, providerIdentifier string, oldToken *oauth2.Token) (*oauth2.Token, error) {
	ctx, span := a.obs.Tracer.Start(ctx, "ProviderAdapterACL.RefreshToken")
	defer span.End()

	a.obs.Logger.Debug(ctx, "Refreshing token",
		zap.String("provider_identifier", providerIdentifier),
		zap.Any("old_token", oldToken),
	)

	token, err := a.contractReader.RefreshTokenContract(ctx, providerIdentifier, oldToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return token, nil
}

func (a *ProviderAdapterACL) GetScopesFromPermissions(ctx context.Context, providerIdentifier string, permissions []string) ([]string, error) {
	ctx, span := a.obs.Tracer.Start(ctx, "ProviderAdapterACL.GetScopesFromPermissions")
	defer span.End()

	a.obs.Logger.Debug(ctx, "Getting scopes from permissions",
		zap.String("provider_identifier", providerIdentifier),
		zap.Strings("permissions", permissions),
	)

	scopes, err := a.contractReader.GetScopesFromPermissionsContract(ctx, providerIdentifier, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to get scopes from permissions: %w", err)
	}

	return scopes, nil
}

func (a *ProviderAdapterACL) GetPermissionIdentifiersFromScopes(ctx context.Context, providerIdentifier string, scopes []string) ([]string, error) {
	ctx, span := a.obs.Tracer.Start(ctx, "ProviderAdapterACL.GetPermissionIdentifiersFromScopes")
	defer span.End()

	a.obs.Logger.Debug(ctx, "Getting permission identifiers from scopes",
		zap.String("provider_identifier", providerIdentifier),
		zap.Strings("scopes", scopes),
	)

	permissionIdentifiers, err := a.contractReader.GetPermissionIdentifiersFromScopesContract(ctx, providerIdentifier, scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission identifiers from scopes: %w", err)
	}

	return permissionIdentifiers, nil
}
