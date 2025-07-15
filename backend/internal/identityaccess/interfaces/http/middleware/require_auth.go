package middleware

import (
	"strings"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/identityaccess/application"
	httpapi "github.com/context-space/context-space/backend/internal/shared/interfaces/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequireAuth middleware ensures the user is authenticated and sets domain.User in context
// This can be used by any module that needs auth with user information
func RequireAuth(
	authService *application.AuthService,
	userService *application.UserService,
	obs *observability.ObservabilityProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			obs.Logger.Debug(ctx, "No authorization header provided")
			httpapi.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			obs.Logger.Debug(ctx, "Invalid authorization header", zap.String("auth_header", authHeader))
			httpapi.Unauthorized(c, "Authorization header must start with Bearer")
			c.Abort()
			return
		}

		authString := parts[1]
		// Check if it's a API key
		if strings.HasPrefix(authString, "cs-") {
			// API key authentication is only allowed for MCP requests
			requestPath := c.Request.URL.Path
			if !strings.HasPrefix(requestPath, "/v1/mcp") && !strings.HasPrefix(requestPath, "/v1/invocations") {
				httpapi.Unauthorized(c, "API key authentication is not allowed for this request")
				c.Abort()
				return
			}

			user, apiKey, err := userService.ValidateAPIKey(ctx, authString)
			if err != nil {
				obs.Logger.Debug(ctx, "Invalid API key", zap.Error(err))
				httpapi.Unauthorized(c, "Invalid API key")
				c.Abort()
				return
			}
			c.Set("api_key", apiKey)
			c.Set("user", user)
			c.Set("auth_type", "api_key")
		} else {
			// Validate token with Supabase
			claims, err := authService.ValidateToken(ctx, authString)
			if err != nil {
				obs.Logger.Debug(ctx, "Invalid token", zap.Error(err))
				httpapi.Unauthorized(c, "Invalid token")
				c.Abort()
				return
			}

			// Get or create user from token
			user, err := authService.FindOrCreateUser(ctx, claims)
			if err != nil {
				obs.Logger.Debug(ctx, "Failed to get or create user from token", zap.Error(err))
				httpapi.Unauthorized(c, "Failed to authenticate user")
				c.Abort()
				return
			}

			c.Set("supabase_claims", claims)
			c.Set("user", user)
			c.Set("auth_type", "supabase")
		}

		c.Next()
	}
}
