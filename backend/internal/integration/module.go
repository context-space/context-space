package integration

import (
	"context"

	"github.com/gin-gonic/gin"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/integration/application"
	"github.com/context-space/context-space/backend/internal/integration/infrastructure/acl"
	"github.com/context-space/context-space/backend/internal/integration/infrastructure/persistence"
	"github.com/context-space/context-space/backend/internal/integration/interfaces/http"
	providercoreApp "github.com/context-space/context-space/backend/internal/providercore/application"
	contractCredential "github.com/context-space/context-space/backend/internal/shared/contract/credentialmanagement"
	contractAdapter "github.com/context-space/context-space/backend/internal/shared/contract/provideradapter"
	contractProvider "github.com/context-space/context-space/backend/internal/shared/contract/providercore"
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
	providerContract contractProvider.ProviderCoreReader,
	adapterContract contractAdapter.ProviderAdapterContract,
	credentialContract contractCredential.CredentialManagementContract,
	providerService *providercoreApp.ProviderService,
	redisClient cache.Cache,
) (*Module, error) {
	// Create repositories
	invocationRepo := persistence.NewInvocationRepository(db, observabilityProvider)

	// Create ACL for provider operations
	providerProvider := acl.NewProviderACL(providerContract, observabilityProvider)

	// Create ACL that uses contract facade for provider adapter
	adapterProvider := acl.NewProviderAdapterACL(
		adapterContract,
		observabilityProvider,
	)

	// Create ACL for credential management
	credProvider := acl.NewCredentialACL(
		credentialContract,
		observabilityProvider,
	)

	// Create application service
	invocationService := application.NewInvocationService(
		providerProvider,
		adapterProvider,
		credProvider,
		invocationRepo,
		eventBus,
		observabilityProvider,
		redisClient,
		credProvider, // Same ACL instance implements both interfaces
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
