package contract

import (
	"context"
	"fmt"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/provideradapter/application"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	contractAdapter "github.com/context-space/context-space/backend/internal/shared/contract/provideradapter"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// AdapterReaderFacade implements the Contract interface at the Interface layer
// Responsibility: Calls the Application layer and handles Domain to DTO conversion
type AdapterContractFacade struct {
	adapterFactory *application.AdapterFactory
	obs            *observability.ObservabilityProvider
}

// Ensure implementation of contract interface
var _ contractAdapter.ProviderAdapterContract = (*AdapterContractFacade)(nil)

// NewAdapterReaderFacade creates a new AdapterReaderFacade
func NewAdapterContractFacade(
	adapterFactory *application.AdapterFactory,
	obs *observability.ObservabilityProvider,
) *AdapterContractFacade {
	return &AdapterContractFacade{
		adapterFactory: adapterFactory,
		obs:            obs,
	}
}

// DomainAdapterWrapper wraps domain adapter to implement contract AdapterDTO
type DomainAdapterWrapper struct {
	domainAdapter domain.Adapter
	obs           *observability.ObservabilityProvider
}

// Execute delegates to domain adapter
func (w *DomainAdapterWrapper) ExecuteContract(
	ctx context.Context,
	operationID string,
	params map[string]interface{},
	credential interface{},
) (interface{}, error) {
	return w.domainAdapter.Execute(ctx, operationID, params, credential)
}

// GetAdapterInfo converts domain adapter info to contract DTO
func (w *DomainAdapterWrapper) GetAdapterInfoContract() *contractAdapter.AdapterInfoDTO {
	domainInfo := w.domainAdapter.GetProviderAdapterInfo()

	return &contractAdapter.AdapterInfoDTO{
		Identifier:  domainInfo.Identifier,
		Name:        domainInfo.Name,
		Description: domainInfo.Description,
		AuthType:    string(domainInfo.AuthType),
		Status:      string(domainInfo.Status),
	}
}

// GetAdapter gets the specified provider adapter and converts it to contract DTO
func (f *AdapterContractFacade) GetAdapterContract(ctx context.Context, providerIdentifier string) (contractAdapter.AdapterContract, error) {
	ctx, span := f.obs.Tracer.Start(ctx, "AdapterReaderFacade.GetAdapter")
	defer span.End()

	f.obs.Logger.Debug(ctx, "Getting adapter through contract facade",
		zap.String("provider_identifier", providerIdentifier))

	domainAdapter, err := f.adapterFactory.GetAdapter(providerIdentifier)
	if err != nil {
		f.obs.Logger.Error(ctx, "Failed to get adapter from factory",
			zap.String("provider_identifier", providerIdentifier),
			zap.Error(err))
		return nil, err
	}

	// Wrap domain adapter to implement contract DTO
	adapterWrapper := &DomainAdapterWrapper{
		domainAdapter: domainAdapter,
		obs:           f.obs,
	}

	f.obs.Logger.Debug(ctx, "Successfully got adapter through contract facade",
		zap.String("provider_identifier", providerIdentifier))

	return adapterWrapper, nil
}

func (f *AdapterContractFacade) GenerateOAuthURLContract(ctx context.Context, providerIdentifier, redirectURL, state, codeChallenge string, scopes []string) (string, error) {
	adapter, err := f.adapterFactory.GetOAuthAdapter(providerIdentifier)
	if err != nil {
		return "", fmt.Errorf("failed to get OAuth adapter: %w", err)
	}

	return adapter.GenerateOAuthURL(ctx, redirectURL, state, codeChallenge, scopes)
}

func (f *AdapterContractFacade) ExchangeCodeForTokenContract(ctx context.Context, providerIdentifier, code, redirectURL, codeVerifier string) (*oauth2.Token, error) {
	adapter, err := f.adapterFactory.GetOAuthAdapter(providerIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth adapter: %w", err)
	}

	return adapter.ExchangeCodeForTokens(ctx, code, redirectURL, codeVerifier)
}

func (f *AdapterContractFacade) ShouldRefreshTokenContract(providerIdentifier string, token *oauth2.Token) (bool, error) {
	adapter, err := f.adapterFactory.GetOAuthAdapter(providerIdentifier)
	if err != nil {
		return false, fmt.Errorf("failed to get OAuth adapter: %w", err)
	}

	return adapter.ShouldRefreshToken(token), nil
}

func (f *AdapterContractFacade) RefreshTokenContract(ctx context.Context, providerIdentifier string, oldToken *oauth2.Token) (*oauth2.Token, error) {
	adapter, err := f.adapterFactory.GetOAuthAdapter(providerIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth adapter: %w", err)
	}

	return adapter.RefreshOAuthToken(ctx, oldToken)
}

func (f *AdapterContractFacade) GetScopesFromPermissionsContract(ctx context.Context, providerIdentifier string, permissions []string) ([]string, error) {
	adapter, err := f.adapterFactory.GetOAuthAdapter(providerIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth adapter: %w", err)
	}

	return adapter.GetScopesFromPermissions(permissions)
}

func (f *AdapterContractFacade) GetPermissionIdentifiersFromScopesContract(ctx context.Context, providerIdentifier string, scopes []string) ([]string, error) {
	adapter, err := f.adapterFactory.GetOAuthAdapter(providerIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth adapter: %w", err)
	}

	return adapter.GetPermissionIdentifiersFromScopes(scopes)
}
