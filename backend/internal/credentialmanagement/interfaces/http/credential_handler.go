package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/application"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	identityDomain "github.com/context-space/context-space/backend/internal/identityaccess/domain"
	httpapi "github.com/context-space/context-space/backend/internal/shared/interfaces/http"
	"github.com/context-space/context-space/backend/internal/shared/security"
	"github.com/context-space/context-space/backend/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CredentialHandler handles HTTP requests for credential management
type CredentialHandler struct {
	credentialService    *application.CredentialService
	oauthStateService    domain.OAuthStateService
	obs                  *observability.ObservabilityProvider
	redirectURLValidator *security.RedirectURLValidator
}

// NewCredentialHandler creates a new instance of CredentialHandler
func NewCredentialHandler(
	credentialService *application.CredentialService,
	oauthStateService domain.OAuthStateService,
	observabilityProvider *observability.ObservabilityProvider,
	redirectURLValidator *security.RedirectURLValidator,
) *CredentialHandler {
	return &CredentialHandler{
		credentialService:    credentialService,
		oauthStateService:    oauthStateService,
		obs:                  observabilityProvider,
		redirectURLValidator: redirectURLValidator,
	}
}

// RegisterRoutes registers the credential routes with the given router
func (h *CredentialHandler) RegisterRoutes(router *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	// Credentials without auth
	router.GET("/credentials/auth/oauth/callback", h.HandleOAuthCallback)

	// Credentials with auth
	credentialsWithAuth := router.Group("/credentials")
	credentialsWithAuth.Use(requireAuth)
	{
		credentialsWithAuth.GET("", h.GetAllCredentialsByUser)
		credentialsWithAuth.GET("/provider/:provider_identifier", h.GetCredentialByUserAndProvider)
		credentialsWithAuth.DELETE("/:id", h.DeleteCredential)

		credentialsWithAuth.POST("/auth/apikey/:provider_identifier", h.CreateAPIKeyCredential)

		credentialsWithAuth.POST("/auth/oauth/:provider_identifier/auth-url", h.CreateOAuthURL)
		credentialsWithAuth.GET("/auth/oauth/state/:oauth_state_id", h.GetOAuthStateData)
	}
}

type CredentialResponse struct {
	ID                 string   `json:"id"`
	UserID             string   `json:"user_id"`
	ProviderIdentifier string   `json:"provider_identifier"`
	Type               string   `json:"type"`
	Permissions        []string `json:"permissions,omitempty"`
	IsValid            bool     `json:"is_valid"`
	CreatedAt          string   `json:"created_at"`
}

type CredentialResponseList struct {
	Credentials []CredentialResponse `json:"credentials"`
}

func (h *CredentialHandler) mapCredentialToResponse(cred *domain.Credential, permissions []string) CredentialResponse {
	return CredentialResponse{
		ID:                 cred.ID,
		UserID:             cred.UserID,
		ProviderIdentifier: cred.ProviderIdentifier,
		Type:               string(cred.Type),
		Permissions:        permissions,
		IsValid:            cred.IsValid,
		CreatedAt:          cred.CreatedAt.Format(time.RFC3339),
	}
}

// GetAllCredentialsByUser godoc
// @Summary Get all credentials for current user
// @Description Returns all credentials associated with the current user
// @Tags credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} httpapi.Response{data=CredentialResponseList} "Success response with credentials list"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /credentials [get]
func (h *CredentialHandler) GetAllCredentialsByUser(c *gin.Context) {
	ctx := c.Request.Context()

	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*identityDomain.User)

	creds, err := h.credentialService.GetAllCredentialsByUser(ctx, user.ID)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to get credentials by user")
		return
	}

	credsResponse := CredentialResponseList{
		Credentials: make([]CredentialResponse, len(creds)),
	}
	for i, cred := range creds {
		credsResponse.Credentials[i] = h.mapCredentialToResponse(cred, nil)
	}

	httpapi.OK(c, credsResponse, "Credentials retrieved successfully")
}

// GetCredentialByUserAndProvider godoc
// @Summary Get credential by provider
// @Description Returns a specific credential for the current user and provider
// @Tags credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param provider_identifier path string true "Provider Identifier"
// @Success 200 {object} httpapi.Response{data=CredentialResponse} "Success response with credential data"
// @Failure 400 {object} httpapi.SwaggerErrorResponse "Bad request error response"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Failure 404 {object} httpapi.SwaggerErrorResponse "Not found error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /credentials/provider/{provider_identifier} [get]
func (h *CredentialHandler) GetCredentialByUserAndProvider(c *gin.Context) {
	ctx := c.Request.Context()

	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*identityDomain.User)

	providerIdentifier := c.Param("provider_identifier")

	if providerIdentifier == "" {
		httpapi.BadRequest(c, "Provider identifier is required")
		return
	}

	cred, err := h.credentialService.GetCredentialByUserAndProvider(ctx, user.ID, providerIdentifier)
	if err != nil {
		if errors.Is(err, application.ErrCredentialNotFound) {
			httpapi.NotFound(c, "Credential not found")
			return
		}
		httpapi.InternalServerError(c, "Failed to get credential by user and provider")
		h.obs.Logger.Error(ctx, "Failed to get credential by user and provider", zap.Error(err))
		return
	}

	var credResponse CredentialResponse
	switch cred := cred.(type) {
	case *domain.OAuthCredential:
		permissions, err := h.credentialService.GetPermissionIdentifiersFromScopes(ctx, cred.ProviderIdentifier, cred.Scopes)
		if err != nil {
			httpapi.InternalServerError(c, "Failed to get permission identifiers from scopes")
			h.obs.Logger.Error(ctx, "Failed to get permission identifiers from scopes", zap.Error(err))
			return
		}
		credResponse = h.mapCredentialToResponse(cred.Credential, permissions)
	case *domain.APIKeyCredential:
		credResponse = h.mapCredentialToResponse(cred.Credential, nil)
	default:
		httpapi.InternalServerError(c, "Invalid credential type")
		return
	}

	// Update credential last used at
	err = h.credentialService.UpdateLastUsedAt(ctx, cred)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to update credential last used at")
		h.obs.Logger.Error(ctx, "Failed to update credential last used at", zap.Error(err))
		return
	}
	httpapi.OK(c, credResponse, "Credential retrieved successfully")
}

// DeleteCredential godoc
// @Summary Delete credential
// @Description Deletes a specific credential by ID
// @Tags credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Credential ID"
// @Success 204 "No content success response"
// @Failure 400 {object} httpapi.SwaggerErrorResponse "Bad request error response"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Failure 403 {object} httpapi.SwaggerErrorResponse "Forbidden error response"
// @Failure 404 {object} httpapi.SwaggerErrorResponse "Not found error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /credentials/{id} [delete]
func (h *CredentialHandler) DeleteCredential(c *gin.Context) {
	ctx := c.Request.Context()

	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*identityDomain.User)

	id := c.Param("id")
	if id == "" {
		httpapi.BadRequest(c, "Credential ID is required")
		return
	}

	credI, err := h.credentialService.GetCredential(ctx, id)
	if err != nil {
		if errors.Is(err, application.ErrCredentialNotFound) {
			httpapi.NotFound(c, "Credential not found")
			return
		}
		httpapi.InternalServerError(c, "Failed to get credential")
		return
	}

	var cred *domain.Credential
	switch credI := credI.(type) {
	case *domain.OAuthCredential:
		cred = credI.Credential
	case *domain.APIKeyCredential:
		cred = credI.Credential
	default:
		httpapi.InternalServerError(c, "Invalid credential type")
		return
	}

	if cred.UserID != user.ID {
		httpapi.Forbidden(c, "You are not allowed to delete this credential")
		return
	}

	err = h.credentialService.DeleteCredential(ctx, cred.ID)
	if err != nil {
		if errors.Is(err, application.ErrCredentialNotFound) {
			httpapi.NotFound(c, "Credential not found")
			return
		}
		httpapi.InternalServerError(c, "Failed to delete credential")
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateAPIKeyCredentialRequest represents the request body for creating an API key credential
type CreateAPIKeyCredentialRequest struct {
	APIKey string `json:"api_key" binding:"required"`
}

// CreateAPIKeyCredential godoc
// @Summary Create API key credential
// @Description Creates a new API key credential for a provider
// @Tags credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param provider_identifier path string true "Provider Identifier"
// @Param request body CreateAPIKeyCredentialRequest true "Create API key credential request"
// @Success 201 {object} httpapi.Response{data=CredentialResponse} "Success response with created credential data"
// @Failure 400 {object} httpapi.SwaggerErrorResponse "Bad request error response"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /credentials/auth/apikey/{provider_identifier} [post]
func (h *CredentialHandler) CreateAPIKeyCredential(c *gin.Context) {
	ctx := c.Request.Context()

	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*identityDomain.User)

	providerIdentifier := c.Param("provider_identifier")
	if providerIdentifier == "" {
		httpapi.BadRequest(c, "Provider identifier is required")
		return
	}

	var req CreateAPIKeyCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.BadRequest(c, "Invalid request format")
		return
	}

	apiKeyCred, err := h.credentialService.CreateAPIKeyCredential(
		ctx,
		user.ID,
		providerIdentifier,
		req.APIKey,
	)

	if err != nil {
		httpapi.InternalServerError(c, "Failed to create credential")
		fmt.Println("Error creating credential:", err)
		return
	}

	credResponse := h.mapCredentialToResponse(apiKeyCred.Credential, nil)
	httpapi.Created(c, credResponse, "Credential created successfully")
}

type CreateOAuthURLRequest struct {
	Permissions []string `json:"permissions" binding:"required"`
	RedirectURL string   `json:"redirect_url" binding:"required"`
}

type CreateOAuthURLResponse struct {
	AuthURL      string `json:"auth_url"`
	OAuthStateID string `json:"oauth_state_id"`
}

// CreateOAuthURL godoc
// @Summary Create OAuth URL
// @Description Generates an OAuth authorization URL for a provider
// @Tags credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param provider_identifier path string true "Provider Identifier"
// @Param request body CreateOAuthURLRequest true "Create OAuth URL request"
// @Success 200 {object} httpapi.Response{data=CreateOAuthURLResponse} "Success response with OAuth URL data"
// @Failure 400 {object} httpapi.SwaggerErrorResponse "Bad request error response"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /credentials/auth/oauth/{provider_identifier}/auth-url [post]
func (h *CredentialHandler) CreateOAuthURL(c *gin.Context) {
	ctx := c.Request.Context()

	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*identityDomain.User)

	providerIdentifier := c.Param("provider_identifier")
	if providerIdentifier == "" {
		httpapi.BadRequest(c, "Provider identifier is required")
		return
	}

	var req CreateOAuthURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.BadRequest(c, "Invalid request format")
		return
	}

	// SECURITY: Validate redirect URL to prevent open redirect attacks
	if err := h.redirectURLValidator.ValidateRedirectURL(req.RedirectURL); err != nil {
		h.obs.Logger.Error(ctx, "Invalid redirect URL provided in OAuth URL request",
			zap.String("redirect_url", req.RedirectURL),
			zap.Error(err))
		httpapi.BadRequest(c, fmt.Sprintf("Invalid redirect URL: %s", err.Error()))
		return
	}

	// Store additional data if needed (like IP address, user agent, etc.)
	oAuthStateData, err := domain.NewOAuthStateData(
		user.ID,
		providerIdentifier,
		req.RedirectURL,
		req.Permissions,
		map[string]interface{}{
			"ip_address": c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		},
	)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to create OAuth state data")
		h.obs.Logger.Error(ctx, fmt.Sprintf("Failed to create OAuth state data for provider %s", providerIdentifier), zap.Error(err))
		return
	}

	// Store state data
	err = h.oauthStateService.StoreStateData(
		ctx,
		oAuthStateData,
	)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to store OAuth state data")
		h.obs.Logger.Error(ctx, fmt.Sprintf("Failed to store OAuth state data for provider %s error:%+v", providerIdentifier, err))
		return
	}

	// Generate OAuth URL
	oauthURL, err := h.credentialService.GetOAuthURL(
		ctx,
		providerIdentifier,
		oAuthStateData.State,
		oAuthStateData.CodeChallenge,
		req.Permissions,
	)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to generate OAuth URL")
		h.obs.Logger.Error(ctx, fmt.Sprintf("Failed to generate OAuth URL for provider %s error:%+v", providerIdentifier, err))
		return
	}

	response := CreateOAuthURLResponse{
		AuthURL:      oauthURL,
		OAuthStateID: oAuthStateData.ID,
	}

	httpapi.OK(c, response, "OAuth URL created successfully")
}

type GetOAuthStateDataResponse struct {
	ID                 string   `json:"id"`
	Status             string   `json:"status"`
	UserID             string   `json:"user_id"`
	ProviderIdentifier string   `json:"provider_identifier"`
	Permissions        []string `json:"permissions"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
}

// GetOAuthStateData godoc
// @Summary Get OAuth state data
// @Description Retrieves OAuth state data by ID
// @Tags credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param oauth_state_id path string true "OAuth State ID"
// @Success 200 {object} httpapi.Response{data=GetOAuthStateDataResponse} "Success response with OAuth state data"
// @Failure 400 {object} httpapi.SwaggerErrorResponse "Bad request error response"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Failure 403 {object} httpapi.SwaggerErrorResponse "Forbidden error response"
// @Failure 404 {object} httpapi.SwaggerErrorResponse "Not found error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /credentials/auth/oauth/state/{oauth_state_id} [get]
func (h *CredentialHandler) GetOAuthStateData(c *gin.Context) {
	ctx := c.Request.Context()

	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*identityDomain.User)

	oauthStateID := c.Param("oauth_state_id")
	if oauthStateID == "" {
		httpapi.BadRequest(c, "OAuth state ID is required")
		return
	}

	oauthStateData, err := h.oauthStateService.GetByID(ctx, oauthStateID)
	if err != nil {
		h.obs.Logger.Error(ctx, "Failed to get OAuth state data", zap.Error(err))
		httpapi.InternalServerError(c, "Failed to get OAuth state data")
		return
	}

	if oauthStateData.UserID != user.ID {
		httpapi.Forbidden(c, "You are not allowed to access this OAuth state data")
		return
	}

	response := GetOAuthStateDataResponse{
		ID:                 oauthStateData.ID,
		Status:             string(oauthStateData.Status),
		UserID:             oauthStateData.UserID,
		ProviderIdentifier: oauthStateData.ProviderIdentifier,
		Permissions:        oauthStateData.Permissions,
		CreatedAt:          oauthStateData.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          oauthStateData.UpdatedAt.Format(time.RFC3339),
	}

	httpapi.OK(c, response, "OAuth state data retrieved successfully")
}

type OAuthCallbackResponse struct{}

// HandleOAuthCallback godoc
// @Summary Handle OAuth callback
// @Description Processes OAuth callback and creates credential
// @Tags credentials
// @Accept json
// @Produce json
// @Param state query string true "OAuth state parameter for verification"
// @Param code query string true "OAuth code for token exchange"
// @Param error query string false "Error message from OAuth provider"
// @Success 302 "Redirect to frontend with success/error parameters"
// @Failure 400 {object} httpapi.SwaggerErrorResponse "Bad request error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /credentials/auth/oauth/callback [get]
func (h *CredentialHandler) HandleOAuthCallback(c *gin.Context) {
	ctx := c.Request.Context()

	callbackState := c.Query("state")
	if callbackState == "" {
		h.obs.Logger.Error(ctx, "State is required")
		httpapi.BadRequest(c, "State is required")
		return
	}

	callbackCode := c.Query("code")
	if callbackCode == "" {
		h.obs.Logger.Error(ctx, "Code is required")
		httpapi.BadRequest(c, "Code is required")
		return
	}

	callbackError := c.Query("error")

	callbackParams := make(map[string]interface{})
	for k, v := range c.Request.URL.Query() {
		if k == "state" || k == "code" {
			continue
		}
		callbackParams[k] = v[0]
	}

	oauthStateData, err := h.oauthStateService.GetByState(ctx, callbackState)
	if err != nil {
		h.obs.Logger.Error(ctx, "Failed to get OAuth state data", zap.Error(err))
		httpapi.InternalServerError(c, "Failed to get OAuth state data")
		return
	}

	oauthStateData.SetCallbackParams(callbackParams)
	redirectURL := oauthStateData.RedirectURL

	// SECURITY: Validate redirect URL to prevent open redirect attacks
	if err := h.redirectURLValidator.ValidateRedirectURL(redirectURL); err != nil {
		h.obs.Logger.Error(ctx, "Invalid redirect URL detected",
			zap.String("redirect_url", redirectURL),
			zap.Error(err))
		httpapi.BadRequest(c, "Invalid redirect URL")
		return
	}

	// Helper function to build redirect URL with parameters
	buildRedirectURL := func(success bool, oauthStateID, errorString string, code int) string {
		values := url.Values{}
		values.Add("success", fmt.Sprintf("%t", success))
		values.Add("oauth_state_id", oauthStateID)
		values.Add("message", errorString)
		values.Add("code", fmt.Sprintf("%d", code))

		if strings.Contains(redirectURL, "?") {
			return utils.StringsBuilder(redirectURL, "&", values.Encode())
		}
		return utils.StringsBuilder(redirectURL, "?", values.Encode())
	}

	// Handle OAuth callback error message
	if callbackCode == "" && callbackError != "" {
		// Set failed status
		oauthStateData.SetStatus(domain.OAuthStateStatusFailed)
		// Update state data
		err = h.oauthStateService.UpdateStateData(ctx, oauthStateData)
		if err != nil {
			h.obs.Logger.Error(ctx, "Failed to update OAuth state status", zap.Error(err))
			httpapi.InternalServerError(c, "Failed to update OAuth state status")
			return
		}
		h.obs.Logger.Warn(ctx, "OAuth callback failed", zap.String("error", callbackError))
		// Redirect with error parameter
		errorRedirectURL := buildRedirectURL(false, oauthStateData.ID, callbackError, http.StatusOK)
		c.Redirect(http.StatusFound, errorRedirectURL)
		return
	}

	// Handle OAuth callback
	_, err = h.credentialService.HandleOAuthCallback(
		ctx,
		callbackCode,
		oauthStateData.ProviderIdentifier,
		oauthStateData.UserID,
		oauthStateData.Permissions,
		oauthStateData.CodeVerifier,
	)
	if err != nil {
		h.obs.Logger.Error(ctx, "Failed to handle OAuth callback", zap.Error(err))
		// Redirect with error parameter
		errorRedirectURL := buildRedirectURL(false, oauthStateData.ID, "Failed to handle OAuth callback", http.StatusOK)
		c.Redirect(http.StatusFound, errorRedirectURL)
		return
	}

	// Update OAuth state status
	oauthStateData.SetStatus(domain.OAuthStateStatusSuccess)
	err = h.oauthStateService.UpdateStateData(ctx, oauthStateData)
	if err != nil {
		h.obs.Logger.Error(ctx, "Failed to update OAuth state status", zap.Error(err))
		// Redirect with error parameter
		errorRedirectURL := buildRedirectURL(false, oauthStateData.ID, "Failed to update OAuth state status", http.StatusOK)
		c.Redirect(http.StatusFound, errorRedirectURL)
		return
	}

	h.obs.Logger.Info(ctx, "OAuth callback handled successfully")
	successRedirectURL := buildRedirectURL(true, oauthStateData.ID, "ok", http.StatusOK)
	c.Redirect(http.StatusFound, successRedirectURL)
}
