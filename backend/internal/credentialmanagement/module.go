package credentialmanagement

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/application"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/infrastructure/adapter"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/infrastructure/persistence"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/infrastructure/vault"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/interfaces/http"
	provideradapterApp "github.com/context-space/context-space/backend/internal/provideradapter/application"
	"github.com/context-space/context-space/backend/internal/shared/config"
	"github.com/context-space/context-space/backend/internal/shared/events"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/cache"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"github.com/context-space/context-space/backend/internal/shared/security"
	"github.com/gin-gonic/gin"
)

// Module holds all components for the credential management bounded context
type Module struct {
	CredentialService    *application.CredentialService
	CredentialHandler    *http.CredentialHandler
	OAuthStateService    domain.OAuthStateService
	tokenRefreshService  domain.TokenRefresh
	credentialRepository domain.CredentialRepository
	credentialFactory    domain.CredentialFactory
	redisClient          cache.Cache
}

// NewModule initializes a new credential management module
func NewModule(
	ctx context.Context,
	db database.Database,
	config *config.Config,
	eventBus events.EventBus,
	observabilityProvider *observability.ObservabilityProvider,
	adapterFactory *provideradapterApp.AdapterFactory,
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

	// Create OAuth provider adapter as an anti-corruption layer
	oauthProviderAdapter := adapter.NewOAuthProviderAdapter(adapterFactory)

	// Initialize OAuth state service
	oauthStateService := application.NewOAuthStateService(oauthStateRepo, observabilityProvider)

	// Initialize token refresh service
	tokenRefreshService := persistence.NewTokenRefreshService(
		redisClient,
		oauthProviderAdapter,
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
		oauthProviderAdapter,
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

	return &Module{
		CredentialService:    credentialService,
		CredentialHandler:    credentialHandler,
		OAuthStateService:    oauthStateService,
		tokenRefreshService:  tokenRefreshService,
		credentialRepository: credentialRepo,
		credentialFactory:    credentialFactory,
		redisClient:          redisClient,
	}, nil
}

// RegisterRoutes registers the credential management routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	m.CredentialHandler.RegisterRoutes(router, requireAuth)
}

// CredentialRepository returns the credential repository for use by other modules
func (m *Module) CredentialRepository() domain.CredentialRepository {
	return m.credentialRepository
}

// CredentialFactory returns the credential factory for use by other modules
func (m *Module) CredentialFactory() domain.CredentialFactory {
	return m.credentialFactory
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
