package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/identityaccess"
	"github.com/context-space/context-space/backend/internal/integration"
	"github.com/context-space/context-space/backend/internal/provideradapter"
	"github.com/context-space/context-space/backend/internal/providercore"
	"github.com/context-space/context-space/backend/internal/shared/config"
	"github.com/context-space/context-space/backend/internal/shared/cron"
	"github.com/context-space/context-space/backend/internal/shared/events"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/cache"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"github.com/context-space/context-space/backend/internal/shared/interfaces/http/middleware"
	"github.com/context-space/context-space/backend/internal/translation"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	// Import generated swagger docs
	_ "github.com/context-space/context-space/backend/docs/swagger"
	"github.com/swaggo/swag"
)

var (
	// Version information - set during build time via ldflags
	Version   = "dev"     // Git tag version
	GitCommit = "unknown" // Git commit hash
	BuildTime = "unknown" // Build timestamp
	GoVersion = "unknown" // Go version used to build
)

// VersionInfo represents the version information structure
type VersionInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
}

// getVersion returns the version string with additional build info
func getVersion() string {
	if Version == "dev" {
		return "dev"
	}
	return Version
}

// getVersionInfo returns complete version information
func getVersionInfo() VersionInfo {
	return VersionInfo{
		Version:   getVersion(),
		GitCommit: GitCommit,
		BuildTime: BuildTime,
		GoVersion: GoVersion,
	}
}

// printVersionInfo prints version information to stdout
func printVersionInfo() {
	info := getVersionInfo()
	fmt.Printf("Context-Space-Backend\n")
	fmt.Printf("Version: %s\n", info.Version)
	fmt.Printf("Git Commit: %s\n", info.GitCommit)
	fmt.Printf("Build Time: %s\n", info.BuildTime)
	fmt.Printf("Go Version: %s\n", info.GoVersion)
}

// checkVersionFlag checks if version flag is provided
func checkVersionFlag() bool {
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" || arg == "version" {
			return true
		}
	}
	return false
}

// @title Context-Space-Backend API
// @version 1.0
// @description API for the Context-Space-Backend integration platform.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host devapi.context.space
// @BasePath /v1
// @schemes https http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format "Bearer {token}"

func main() {
	// Check for version flag first
	if checkVersionFlag() {
		printVersionInfo()
		return
	}

	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			// log to stderr
			fmt.Fprintf(os.Stderr, "[PANIC] %s: %v\n", time.Now().Format("2006-01-02 15:04:05"), r)

			// log to system log
			log.Printf("PANIC: %v", r)

			// log to file
			f, err := os.OpenFile("panic.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				defer f.Close()
				fmt.Fprintf(f, "[%s] PANIC: %v\n", time.Now().Format("2006-01-02 15:04:05"), r)
			}

			// force flush all outputs
			os.Stderr.Sync()
			os.Stdout.Sync()

			os.Exit(1)
		}
	}()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Log version information
	logger.Info("Starting Context-Space-Backend",
		zap.String("version", getVersion()),
		zap.String("git_commit", GitCommit),
		zap.String("build_time", BuildTime),
		zap.String("go_version", GoVersion),
	)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.Info("Configuration loaded",
		zap.String("environment", cfg.Environment),
		zap.String("server_address", cfg.Server.Address),
		zap.String("supabase_project_ref", cfg.Supabase.ProjectRef),
		zap.String("database_host", cfg.Database.Host),
		zap.Int("database_port", cfg.Database.Port),
		zap.String("database_username", cfg.Database.Username),
		zap.String("database_database", cfg.Database.Database),
		zap.String("database_ssl_mode", cfg.Database.SSLMode),
		zap.String("redis_host", cfg.Redis.Host),
		zap.Int("redis_port", cfg.Redis.Port),
		zap.String("redis_username", cfg.Redis.Username),
		zap.Int("redis_db", cfg.Redis.DB),
		zap.String("vault_default_region", cfg.Vault.DefaultRegion),
		zap.String("logging_level", cfg.Logging.Level),
		zap.String("logging_format", cfg.Logging.Format),
		zap.String("logging_output_path", cfg.Logging.OutputPath),
	)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize observability provider
	ctx := context.Background()
	logConfig := &observability.LogConfig{
		Level:       observability.ParseLogLevel(cfg.Logging.Level),
		Format:      observability.ParseLogFormat(cfg.Logging.Format),
		OutputPaths: []string{cfg.Logging.OutputPath},
		Development: cfg.Environment != "production",
	}

	serviceVersion := getVersion()
	tracingConfig := &observability.TracingConfig{
		ServiceName:    "context-space-backend",
		ServiceVersion: serviceVersion,
		Environment:    cfg.Environment,
		Endpoint:       cfg.Observability.Tracing.Endpoint,
		Enabled:        cfg.Observability.Tracing.Enabled,
		SamplingRate:   cfg.Observability.Tracing.SamplingRate,
	}

	metricsConfig := &observability.MetricsConfig{
		ServiceName:    "context-space-backend",
		ServiceVersion: serviceVersion,
		Environment:    cfg.Environment,
		Enabled:        cfg.Observability.Metrics.Enabled,
		Endpoint:       cfg.Observability.Metrics.Endpoint,
	}

	observabilityProvider, obsCleanup, err := observability.InitializeObservabilityProvider(
		ctx,
		logConfig,
		tracingConfig,
		metricsConfig,
	)
	if err != nil {
		log.Fatalf("Failed to initialize observability: %v", err)
	}
	defer obsCleanup()

	// Initialize database
	postgresClient, err := database.NewPostgresClient(cfg, observabilityProvider)
	if err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize database", zap.Error(err))
	}

	// Initialize Redis client
	redisClient, err := cache.NewRedisClient(cfg, observabilityProvider)
	if err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize Redis client", zap.Error(err))
	}

	// Initialize event bus
	eventBus := events.NewBus()

	// Initialize identity and access module
	identityAccessModule, err := identityaccess.NewModule(
		postgresClient,
		eventBus,
		observabilityProvider,
		cfg,
	)
	if err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize identity and access module", zap.Error(err))
	}

	// Initialize provider translation module
	providerTranslationModule, err := translation.NewModule(
		postgresClient,
		observabilityProvider,
	)
	if err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize provider translation module", zap.Error(err))
	}

	// Initialize provider core module
	providerCoreModule, err := providercore.NewModule(
		postgresClient,
		eventBus,
		observabilityProvider,
		providerTranslationModule.GetProviderTranslationService(),
	)
	if err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize provider core module", zap.Error(err))
	}

	// Initialize provider adapter module
	providerAdapterModule, err := provideradapter.NewModule(
		postgresClient,
		observabilityProvider,
		providerCoreModule.GetProviderService(),
		providerTranslationModule.GetProviderTranslationService(),
	)
	if err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize provider adapter module", zap.Error(err))
	}

	// Initialize credential management module
	credentialManagementModule, err := credentialmanagement.NewModule(
		ctx,
		postgresClient,
		cfg,
		eventBus,
		observabilityProvider,
		providerAdapterModule.GetAdapterContractFacade(),
		redisClient,
	)
	if err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize credential management module", zap.Error(err))
	}

	// Initialize integration module
	integrationModule, err := integration.NewModule(
		postgresClient,
		eventBus,
		observabilityProvider,
		providerCoreModule.GetProviderService(),
		providerAdapterModule.GetAdapterContractFacade(),
		credentialManagementModule.GetCredentialContractFacade(),
		providerCoreModule.GetProviderService(),
		redisClient,
	)
	if err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize integration module", zap.Error(err))
	}

	router := gin.New()

	// Configure CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:           cfg.Security.CORS.AllowedOrigins,
		AllowMethods:           []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:           []string{"Origin", "Authorization", "Content-Type", "Accept", "Content-Length", "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders:          []string{"Content-Length", "Content-Type"},
		AllowCredentials:       true,
		MaxAge:                 12 * time.Hour,
		AllowWildcard:          true,
		AllowBrowserExtensions: true,
	}))

	// Add recovery middleware
	router.Use(gin.Recovery())

	// Register observability middleware
	middleware.RegisterObservabilityMiddleware(router, observabilityProvider)

	// Initialize routes
	initializeRoutes(
		router,
		providerAdapterModule,
		providerCoreModule,
		identityAccessModule,
		credentialManagementModule,
		integrationModule,
	)

	// Initialize modules
	ctx = context.Background()

	// Initialize provider adapter module
	if err := providerAdapterModule.Initialize(ctx); err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize provider adapter module", zap.Error(err))
	}

	// Initialize identity and access module
	if err := identityAccessModule.Initialize(ctx); err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize identity and access module", zap.Error(err))
	}

	// Initialize provider core module
	if err := providerCoreModule.Initialize(ctx); err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize provider core module", zap.Error(err))
	}

	// Initialize credential management module
	if err := credentialManagementModule.Initialize(ctx); err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize credential management module", zap.Error(err))
	}

	// Initialize integration module
	if err := integrationModule.Initialize(ctx); err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize integration module", zap.Error(err))
	}

	// Initialize cron jobs system
	if err := initializeCronJobs(ctx, credentialManagementModule.TokenRefreshService(), observabilityProvider, redisClient); err != nil {
		observabilityProvider.Logger.Fatal(ctx, "Failed to initialize cron jobs system", zap.Error(err))
	}

	// Create an HTTP server
	server := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}

	// Start the HTTP server in a goroutine
	go func() {
		logger.Info("Starting HTTP server", zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			observabilityProvider.Logger.Fatal(ctx, "Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down servers...")

	// log memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logger.Info("Server shutdown - final memory stats",
		zap.Uint64("alloc_mb", m.Alloc/1024/1024),
		zap.Uint64("sys_mb", m.Sys/1024/1024),
		zap.Int("goroutines", runtime.NumGoroutine()),
	)

	// Create a deadline for the shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		observabilityProvider.Logger.Fatal(ctx, "HTTP server forced to shutdown", zap.Error(err))
	}

	logger.Info("Servers exiting")
}

// initializeRoutes registers all API routes
func initializeRoutes(
	router *gin.Engine,
	providerAdapterModule *provideradapter.Module,
	providerCoreModule *providercore.Module,
	identityAccessModule *identityaccess.Module,
	credentialManagementModule *credentialmanagement.Module,
	integrationModule *integration.Module,
) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		versionInfo := getVersionInfo()

		healthInfo := gin.H{
			"status":     "ok",
			"timestamp":  time.Now().Unix(),
			"version":    versionInfo.Version,
			"git_commit": versionInfo.GitCommit,
			"build_time": versionInfo.BuildTime,
			"go_version": versionInfo.GoVersion,
			"memory": gin.H{
				"alloc_mb":       m.Alloc / 1024 / 1024,
				"total_alloc_mb": m.TotalAlloc / 1024 / 1024,
				"sys_mb":         m.Sys / 1024 / 1024,
			},
			"goroutines": runtime.NumGoroutine(),
			"gc_count":   m.NumGC,
		}

		c.JSON(http.StatusOK, healthInfo)
	})

	// Serve static files from resources directory
	router.Static("/resources", "./resources")

	// API Version group
	v1 := router.Group("/v1")
	{
		// Get authentication middleware that extracts domain.User
		requireAuthMiddleware := identityAccessModule.GetRequireAuthMiddleware()

		// Register provider adapter routes
		providerAdapterModule.RegisterRoutes(v1, requireAuthMiddleware)

		// Register provider core routes
		providerCoreModule.RegisterRoutes(v1, requireAuthMiddleware)

		// Register identity and access routes
		identityAccessModule.RegisterRoutes(v1, requireAuthMiddleware)

		// Register credential management routes
		credentialManagementModule.RegisterRoutes(v1, requireAuthMiddleware)

		// Register integration routes
		integrationModule.RegisterRoutes(v1, requireAuthMiddleware)

		// Custom handler for docs endpoints with path parameter
		v1.GET("/docs/*any", func(c *gin.Context) {
			// Get the current request host and scheme
			host := c.Request.Host

			isLocal := strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") || host == "" || strings.HasPrefix(host, "192.168.")

			scheme := "https"
			if isLocal {
				scheme = "http"
			} else if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
				scheme = "https"
			} else if c.GetHeader("X-Forwarded-Proto") == "http" {
				scheme = "http"
			}

			// Create swagger handler with dynamic URL
			dynamicHandler := ginSwagger.WrapHandler(
				swaggerFiles.Handler,
				ginSwagger.URL(fmt.Sprintf("%s://%s/v1/swagger.json", scheme, host)),
			)

			// Get the path after /docs/
			path := c.Param("any")

			// If it's empty or just a slash, redirect to index.html
			if path == "" || path == "/" {
				c.Redirect(http.StatusMovedPermanently, "/v1/docs/index.html")
				return
			}

			// Use the dynamic swagger handler
			dynamicHandler(c)
		})

		// Handle exact /docs endpoint (redirect to /docs/index.html)
		v1.GET("/docs", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/v1/docs/index.html")
		})

		// Serve swagger.json directly with dynamic host
		v1.GET("/swagger.json", func(c *gin.Context) {
			spec, err := swag.ReadDoc()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var jsonSpec map[string]interface{}
			if err := json.Unmarshal([]byte(spec), &jsonSpec); err == nil {
				host := c.Request.Host
				if host != "" {
					jsonSpec["host"] = host
				}

				isLocal := strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") || host == "" || strings.HasPrefix(host, "192.168.")

				if isLocal {
					jsonSpec["schemes"] = []string{"http", "https"}
				} else {
					jsonSpec["schemes"] = []string{"https", "http"}
				}

				if updatedSpec, err := json.Marshal(jsonSpec); err == nil {
					spec = string(updatedSpec)
				}
			}

			c.Header("Content-Type", "application/json")
			c.String(http.StatusOK, spec)
		})
	}
}

func initializeCronJobs(ctx context.Context,
	tokenRefreshService domain.TokenRefresh,
	observabilityProvider *observability.ObservabilityProvider,
	redisClient cache.Cache,
) error {
	observabilityProvider.Logger.Info(ctx, "Initializing cron jobs system")

	cronManager := cron.NewCronManager(observabilityProvider, redisClient)

	taskBuilder := cron.NewCronTaskBuilder(tokenRefreshService, observabilityProvider)

	taskGroups := taskBuilder.CreateAllTaskGroups()
	for _, group := range taskGroups {
		if err := cronManager.RegisterTaskGroup(ctx, group); err != nil {
			observabilityProvider.Logger.Error(ctx, "Failed to register task group",
				zap.String("group", group.Name),
				zap.Error(err))
			return err
		}
	}

	if err := cronManager.Start(ctx); err != nil {
		observabilityProvider.Logger.Error(ctx, "Failed to start cron manager", zap.Error(err))
		return err
	}

	observabilityProvider.Logger.Info(ctx, "Cron jobs system initialization completed",
		zap.Strings("task_groups", cronManager.ListTaskGroups()))

	return nil
}
