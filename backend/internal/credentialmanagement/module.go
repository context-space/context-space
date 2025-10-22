package credentialmanagement

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/application"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/infrastructure/acl"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/infrastructure/persistence"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/infrastructure/vault"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/interfaces/contract"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/interfaces/http"
	"github.com/context-space/context-space/backend/internal/shared/config"
	contractCredential "github.com/context-space/context-space/backend/internal/shared/contract/credentialmanagement"
	contractAdapter "github.com/context-space/context-space/backend/internal/shared/contract/provideradapter"
	"github.com/context-space/context-space/backend/internal/shared/events"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/cache"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"github.com/context-space/context-space/backend/internal/shared/security"
	"github.com/gin-gonic/gin"
)

// Module holds all components for the credential management bounded context
type Module struct {
	CredentialService        *application.CredentialService
	CredentialHandler        *http.CredentialHandler
	OAuthStateService        domain.OAuthStateService
	tokenRefreshService      domain.TokenRefresh
	credentialRepository     domain.CredentialRepository
	credentialFactory        domain.CredentialFactory
	redisClient              cache.Cache
	credentialContractFacade contractCredential.CredentialManagementContract
}

// NewModule initializes a new credential management module
func NewModule(
	ctx context.Context,
	db database.Database,
	config *config.Config,
	eventBus events.EventBus,
	observabilityProvider *observability.ObservabilityProvider,
	providerAdapterContract contractAdapter.ProviderAdapterContract,
	redisClient cache.Cache,
) (*Module, error) {
	unitOfWorkFactory := database.NewDefaultUnitOfWorkFactory(db, observabilityProvider)

	// Initialize repositories
	credentialRepo := persistence.NewCredentialRepository(db, observabilityProvider)
	oauthRepo := persistence.NewOAuthCredentialRepository(db, observabilityProvider)
	apiKeyRepo := persistence.NewAPIKeyCredentialRepository(db, observabilityProvider)

	// Initialize OAuth state repository
	oauthStateRepo := persistence.NewRedisOAuthStateRepository(db, redisClient, observabilityProvider, application.DefaultStateExpiration)

	// Initialize vault service
	vaultService, err := vault.NewVaultService(ctx, &vault.VaultConfig{
		Regions: convertVaultRegionsMap(config.Vault.Regions),
	}, domain.VaultRegion(config.Vault.DefaultRegion))
	if err != nil {
		return nil, err
	}

	// Initialize credential factory
	credentialFactory := persistence.NewCredentialFactory(
		credentialRepo,
		oauthRepo,
		apiKeyRepo,
		vaultService,
	)

	providerAdapterACL := acl.NewProviderAdapterACL(providerAdapterContract, observabilityProvider)

	// Initialize OAuth state service
	oauthStateService := application.NewOAuthStateService(oauthStateRepo, observabilityProvider)

	// Initialize token refresh service
	tokenRefreshService := persistence.NewTokenRefreshService(
		redisClient,
		providerAdapterACL,
		credentialRepo,
		oauthRepo,
		vaultService,
		observabilityProvider,
	)

	// Initialize credential service
	credentialService := application.NewCredentialService(
		credentialRepo,
		credentialFactory,
		eventBus,
		observabilityProvider,
		providerAdapterACL,
		unitOfWorkFactory,
		config.Provider.OAuthRedirectURL,
		redisClient,
		tokenRefreshService,
	)

	// Initialize OAuth redirect URL validator
	redirectURLValidator := security.NewRedirectURLValidator(
		config.Security.RedirectURLValidator.AllowedDomains,
		config.Security.RedirectURLValidator.AllowedSchemes,
	)

	// Initialize credential handler with OAuth state service
	credentialHandler := http.NewCredentialHandler(
		credentialService,
		oauthStateService,
		observabilityProvider,
		redirectURLValidator,
	)

	// Initialize credential contract facade
	credentialContractFacade := contract.NewCredentialContractFacade(
		credentialService,
		credentialFactory,
		tokenRefreshService,
		observabilityProvider,
	)

	return &Module{
		CredentialService:        credentialService,
		CredentialHandler:        credentialHandler,
		OAuthStateService:        oauthStateService,
		tokenRefreshService:      tokenRefreshService,
		credentialRepository:     credentialRepo,
		credentialFactory:        credentialFactory,
		redisClient:              redisClient,
		credentialContractFacade: credentialContractFacade,
	}, nil
}

// RegisterRoutes registers the credential management routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	m.CredentialHandler.RegisterRoutes(router, requireAuth)
}

func (m *Module) GetCredentialContractFacade() contractCredential.CredentialManagementContract {
	return m.credentialContractFacade
}

// TokenRefreshService returns the token refresh service for use by other modules
func (m *Module) TokenRefreshService() domain.TokenRefresh {
	return m.tokenRefreshService
}

// Initialize initializes the credential management module
func (m *Module) Initialize(ctx context.Context) error {
	return nil
}

// Close closes any resources managed by the module
func (m *Module) Close() error {
	if m.redisClient != nil {
		return m.redisClient.Close()
	}
	return nil
}

// Helper function to convert config.VaultRegionalConfig to vault.VaultRegionalConfig
func convertVaultRegionsMap(regions map[string]*config.VaultRegionalConfig) map[domain.VaultRegion]vault.VaultRegionalConfig {
	result := make(map[domain.VaultRegion]vault.VaultRegionalConfig)
	for region, cfg := range regions {
		result[domain.VaultRegion(region)] = vault.VaultRegionalConfig{
			Address:     cfg.Address,
			Token:       cfg.Token,
			TransitPath: cfg.TransitPath,
		}
	}
	return result
}
