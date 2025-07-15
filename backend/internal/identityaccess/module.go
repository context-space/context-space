package identityaccess

import (
	"context"
	"fmt"

	"github.com/context-space/context-space/backend/internal/identityaccess/domain"
	"github.com/context-space/context-space/backend/internal/identityaccess/infrastructure/supabase"
	"github.com/context-space/context-space/backend/internal/identityaccess/interfaces/http/middleware"

	"github.com/gin-gonic/gin"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/identityaccess/application"
	"github.com/context-space/context-space/backend/internal/identityaccess/infrastructure/persistence"
	iahttp "github.com/context-space/context-space/backend/internal/identityaccess/interfaces/http"
	"github.com/context-space/context-space/backend/internal/shared/config"
	"github.com/context-space/context-space/backend/internal/shared/events"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
)

// Module encapsulates all identity and access components
type Module struct {
	userService *application.UserService
	userHandler *iahttp.UserHandler
	authService *application.AuthService
	userRepo    domain.UserRepository
	obs         *observability.ObservabilityProvider
}

// NewModule creates a new identity and access module
func NewModule(
	db database.Database,
	eventBus *events.Bus,
	observabilityProvider *observability.ObservabilityProvider,
	cfg *config.Config,
) (*Module, error) {
	// Create Unit of Work
	unitOfWorkFactory := database.NewDefaultUnitOfWorkFactory(db, observabilityProvider)

	// Create repositories
	userRepo := persistence.NewUserRepository(db, observabilityProvider)
	userInfoRepo := persistence.NewUserInfoRepository(db, observabilityProvider)
	apiKeyRepo := persistence.NewUserAPIKeyRepository(db, observabilityProvider)

	// Initialize Supabase auth service
	supabaseConfig := &supabase.SupabaseConfig{
		ProjectRef:  cfg.Supabase.ProjectRef,
		ServiceRole: cfg.Supabase.ServiceRole,
		JWTSecret:   cfg.Supabase.JWTSecret,
	}

	supabaseAuthService, err := supabase.NewSupabaseAuthService(supabaseConfig, observabilityProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Supabase auth service: %w", err)
	}

	// Create application services
	authService := application.NewAuthService(
		userRepo,
		userInfoRepo,
		supabaseAuthService,
		unitOfWorkFactory,
		eventBus,
		observabilityProvider,
	)

	// Create the user service
	userService := application.NewUserService(
		userRepo,
		userInfoRepo,
		apiKeyRepo,
		unitOfWorkFactory,
		eventBus,
		observabilityProvider,
	)

	// Create HTTP handler
	userHandler := iahttp.NewUserHandler(userService, observabilityProvider)

	return &Module{
		userService: userService,
		userHandler: userHandler,
		authService: authService,
		userRepo:    userRepo,
		obs:         observabilityProvider,
	}, nil
}

// Initialize initializes the integration module
func (m *Module) Initialize(ctx context.Context) error {
	// Initialize the registry service to load providers
	return nil
}

// RegisterRoutes registers all Identity and Access HTTP routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	// Register routes with appropriate middleware
	m.userHandler.RegisterRoutes(router, requireAuth)
}

// GetRequireAuthMiddleware returns a middleware that authenticates requests and extracts domain.User
// Other modules can use this to secure their routes and get access to the domain.User object
func (m *Module) GetRequireAuthMiddleware() gin.HandlerFunc {
	return middleware.RequireAuth(m.authService, m.userService, m.obs)
}
