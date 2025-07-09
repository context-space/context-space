package integration

import (
	"context"

	"github.com/gin-gonic/gin"

	observability "github.com/context-space/cloud-observability"
	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	credentialDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/integration/application"
	"github.com/context-space/context-space/backend/internal/integration/domain"
	"github.com/context-space/context-space/backend/internal/integration/infrastructure/adapter"
	"github.com/context-space/context-space/backend/internal/integration/infrastructure/credential"
	"github.com/context-space/context-space/backend/internal/integration/infrastructure/persistence"
	"github.com/context-space/context-space/backend/internal/integration/interfaces/http"
	providerAdapterApp "github.com/context-space/context-space/backend/internal/provideradapter/application"
	providercoreApp "github.com/context-space/context-space/backend/internal/providercore/application"
	"github.com/context-space/context-space/backend/internal/shared/events"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/cache"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
)

// Module encapsulates all integration components
type Module struct {
	InvocationService *application.InvocationService
	InvocationHandler *http.InvocationHandler
	McpHandler        *http.McpHandler
	obs               *observability.ObservabilityProvider
}

// NewModule creates a new integration module
func NewModule(
	db database.Database,
	eventBus *events.Bus,
	observabilityProvider *observability.ObservabilityProvider,
	providerProvider domain.ProviderProvider,
	adapterFactory *providerAdapterApp.AdapterFactory,
	credentialFactory credDomain.CredentialFactory,
	providerService *providercoreApp.ProviderService,
	redisClient cache.Cache,
	tokenRefreshService credentialDomain.TokenRefresh,
) (*Module, error) {
	// Create repositories
	invocationRepo := persistence.NewInvocationRepository(db, observabilityProvider)

	// Create provider adapter
	adapterProvider := adapter.NewProviderAdapter(adapterFactory)
	credProvider := credential.NewCredentialProviderAdapter(credentialFactory)

	// Create application service
	invocationService := application.NewInvocationService(
		providerProvider,
		adapterProvider,
		credProvider,
		invocationRepo,
		eventBus,
		observabilityProvider,
		redisClient,
		tokenRefreshService,
	)

	// Create HTTP handler
	invocationHandler := http.NewInvocationHandler(invocationService, observabilityProvider)
	mcpHandler := http.NewMcpHandler(invocationService, providerService, observabilityProvider)

	return &Module{
		InvocationService: invocationService,
		InvocationHandler: invocationHandler,
		McpHandler:        mcpHandler,
		obs:               observabilityProvider,
	}, nil
}

// Initialize initializes the integration module
func (m *Module) Initialize(ctx context.Context) error {
	m.obs.Logger.Info(ctx, "Initializing Integration module")
	// Nothing to initialize for now
	return nil
}

// RegisterRoutes registers all integration HTTP routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	m.InvocationHandler.RegisterRoutes(router, requireAuth)
	m.McpHandler.RegisterRoutes(router, requireAuth)
}

// GetInvocationService returns the invocation service
func (m *Module) GetInvocationService() *application.InvocationService {
	return m.InvocationService
}
