package http

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	observability "github.com/context-space/cloud-observability"
	identityDomain "github.com/context-space/context-space/backend/internal/identityaccess/domain"
	"github.com/context-space/context-space/backend/internal/integration/application"
	integrationDomain "github.com/context-space/context-space/backend/internal/integration/domain"
	providercoreApp "github.com/context-space/context-space/backend/internal/providercore/application"
	httpapi "github.com/context-space/context-space/backend/internal/shared/interfaces/http"
	"github.com/context-space/context-space/backend/internal/shared/types"
	"github.com/context-space/context-space/backend/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CallToolRequest defines the expected request body for the call_tool endpoint.
// This struct is no longer directly bound from the main request body with the new path parameter approach.
// The request body will be directly bound to a map[string]interface{} for the tool's input parameters.
type CallToolRequest struct {
	ToolName string                 `json:"tool_name"` // This will be derived from path parameters
	Input    map[string]interface{} `json:"input"`     // This will be the entire request body
}

// CallToolResponse defines the response structure for the call_tool endpoint.
type CallToolResponse struct {
	ToolResult json.RawMessage `json:"tool_result,omitempty"`
	Error      string          `json:"error,omitempty"` // Error message if the tool call failed
}

// ListToolsRequest defines the request body for the list_tools endpoint.
type ListToolsRequest struct {
	Query         *string `json:"query"`          // Optional: A query string to filter tools by name or description.
	Context       *string `json:"context"`        // Optional: Contextual information to help recommend relevant tools.
	AllowDisabled *bool   `json:"allow_disabled"` // Optional: Whether to include tools from disabled providers. Defaults to false.
}

// ToolParameterDefinition describes a single parameter of a tool.
type ToolParameterDefinition struct {
	Type        string   `json:"type"`                  // Parameter type (e.g., "string", "integer", "boolean", "object", "array")
	Description string   `json:"description,omitempty"` // Description of the parameter.
	Required    bool     `json:"required"`              // Whether the parameter is required.
	Enum        []string `json:"enum,omitempty"`        // Optional: Possible values for the parameter.
	Default     any      `json:"default,omitempty"`     // Optional: Default value for the parameter.
}

// ToolDefinition describes a single tool available to the MCP client.
type ToolDefinition struct {
	Name                string                             `json:"name"` // Format: "provider_identifier.operation_identifier"
	Description         string                             `json:"description"`
	ParametersSchema    map[string]ToolParameterDefinition `json:"parameters_schema"`           // JSON Schema for the tool's input parameters.
	RequiredPermissions []string                           `json:"required_permissions"`        // List of permission identifiers required to use this tool.
	ProviderName        string                             `json:"provider_name"`               // Display name of the provider.
	ProviderIconURL     string                             `json:"provider_icon_url,omitempty"` // Icon URL of the provider.
}

// ListToolsResponse defines the response structure for the list_tools endpoint.
type ListToolsResponse struct {
	Tools []ToolDefinition `json:"tools"`
}

// McpToolParameters represents the "parameters" object for a tool method.
// Example: { "query": "string", "maxResults": "number" }
// type McpToolParameters map[string]string // Key: param name, Value: param type // This type will be removed

// McpToolMethod represents a single method for a provider.
// Example: { "parameters": { "query": "string", "maxResults": "number" } }
type McpToolMethod struct {
	Description string                             `json:"description"`
	Parameters  map[string]ToolParameterDefinition `json:"parameters"`
}

// McpProviderTools represents all operations for a given provider.
// Example: { "operations": { "listEmails": { ... }, "sendEmail": { ... } } }
type McpProviderTools struct {
	Operations map[string]McpToolMethod `json:"operations"` // Key: operation_identifier (e.g., "listEmails")
}

// NewListToolsResponse defines the new response structure for the list_tools endpoint.
type NewListToolsResponse struct {
	Tools map[string]McpProviderTools `json:"tools"` // Key: provider_identifier (e.g., "gmail")
}

// McpHandler handles HTTP requests for MCP (Meta Call Protocol) endpoints.
type McpHandler struct {
	invocationService *application.InvocationService
	providerService   *providercoreApp.ProviderService
	obs               *observability.ObservabilityProvider
}

// NewMcpHandler creates a new McpHandler.
func NewMcpHandler(
	invocationService *application.InvocationService,
	providerService *providercoreApp.ProviderService,
	observabilityProvider *observability.ObservabilityProvider,
) *McpHandler {
	return &McpHandler{
		invocationService: invocationService,
		providerService:   providerService,
		obs:               observabilityProvider,
	}
}

// RegisterRoutes registers the MCP routes for this handler.
func (h *McpHandler) RegisterRoutes(router *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	mcpGroup := router.Group("/mcp")
	mcpGroup.Use(requireAuth)
	// Note: The decision to apply requireAuth middleware can be managed here or passed as a parameter.
	// For now, matching the temporary removal in module.go:
	// mcpGroup.Use(requireAuth)
	{
		mcpGroup.POST("/call_tool/:provider_identifier/:operation_identifier", h.HandleMcpCallTool)
		mcpGroup.POST("/list_tools", h.HandleMcpListTools)
	}
}

// HandleMcpCallTool godoc
// @Summary Call a tool (Provider Operation)
// @Description Executes a specific tool (Provider Operation) with the given input. The input parameters should be provided as the JSON request body.
// @Tags mcp
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param provider_identifier path string true "Identifier of the provider (e.g., 'gmail')"
// @Param operation_identifier path string true "Identifier of the operation (e.g., 'sendEmail')"
// @Param request_body body map[string]interface{} true "Input parameters for the tool method"
// @Success 200 {object} httpapi.Response{data=CallToolResponse} "Success response with tool execution result"
// @Failure 400 {object} httpapi.SwaggerErrorResponse "Bad request if input is invalid"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized if JWT is missing or invalid"
// @Failure 403 {object} httpapi.SwaggerErrorResponse "Forbidden if user does not have credential or permission"
// @Failure 404 {object} httpapi.SwaggerErrorResponse "Not found if provider or operation does not exist"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error or tool execution error"
// @Router /mcp/call_tool/{provider_identifier}/{operation_identifier} [post]
func (h *McpHandler) HandleMcpCallTool(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.obs.Logger.With(zap.String("handler", "McpHandler"), zap.String("method", "HandleMcpCallTool"))

	// 1. Get User from context (set by require_auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		logger.Warn(ctx, "User object not found in context")
		httpapi.Unauthorized(c, "Authentication required, user context missing")
		return
	}

	user, ok := userInterface.(*identityDomain.User)
	if !ok || user == nil {
		logger.Error(ctx, "User object in context is not of expected type *identityDomain.User or is nil")
		httpapi.InternalServerError(c, "Internal authentication error, invalid user context")
		return
	}
	userID := user.ID
	logger = logger.With(zap.String("userID", userID))

	// 2. Get provider_identifier and operation_identifier from path parameters
	providerIdentifier := c.Param("provider_identifier")
	operationIdentifier := c.Param("operation_identifier")

	if providerIdentifier == "" {
		logger.Warn(ctx, "provider_identifier path parameter is missing")
		httpapi.BadRequest(c, "Missing provider_identifier path parameter")
		return
	}
	if operationIdentifier == "" {
		logger.Warn(ctx, "operation_identifier path parameter is missing")
		httpapi.BadRequest(c, "Missing operation_identifier path parameter")
		return
	}
	logger = logger.With(zap.String("providerIdentifier", providerIdentifier), zap.String("operationIdentifier", operationIdentifier))

	// 3. Bind request body directly to a map for parameters
	var params map[string]interface{}
	if err := c.ShouldBindJSON(&params); err != nil {
		// Handle cases where body is empty or not valid JSON, but allow empty map if body is effectively empty.
		// EOF means empty body, which can be treated as empty params.
		if err == io.EOF {
			logger.Info(ctx, "Request body is empty, treating as empty parameters for tool call.")
			params = make(map[string]interface{}) // Ensure params is not nil
		} else {
			logger.Warn(ctx, "Failed to bind JSON for tool parameters", zap.Error(err))
			httpapi.BadRequest(c, utils.StringsBuilder("Invalid request format for tool parameters: ", err.Error()))
			return
		}
	}
	logger.Info(ctx, "Received CallTool request", zap.Any("input_keys", keys(params)))

	// 4. Call InvocationService
	// Note: domain.Invocation might be returned even if err is not nil, e.g. if adapter execution fails
	// but the invocation record itself was created.
	invocation, err := h.invocationService.InvokeOperation(ctx, userID, providerIdentifier, operationIdentifier, params) // Use params directly

	// Handle errors from InvokeOperation or failed invocation status
	if err != nil {
		logger.Error(ctx, "InvocationService.InvokeOperation returned an error", zap.Error(err))

		// Specific error handling based on errors from application layer
		if errors.Is(err, application.ErrProviderNotFound) || errors.Is(err, application.ErrOperationNotFound) {
			httpapi.NotFound(c, utils.StringsBuilder(providerIdentifier, ".", operationIdentifier, " not found."))
			return
		}
		if errors.Is(err, application.ErrCredentialNotFound) {
			h.obs.Logger.Info(ctx, utils.StringsBuilder("Access denied: missing or invalid credentials for ", providerIdentifier), zap.Error(err))
			httpapi.Forbidden(c, utils.StringsBuilder("Access denied: missing or invalid credentials for ", providerIdentifier))
			return
		}
		if errors.Is(err, application.ErrInvalidParameters) {
			logger.Warn(ctx, "Invalid parameters reported by InvocationService", zap.Error(err))
			httpapi.BadRequest(c, utils.StringsBuilder("Invalid parameters for tool: ", err.Error()))
			return
		}
		if errors.Is(err, application.ErrAdapterExecuteFailed) {
			errMsg := "Tool execution failed."
			if invocation != nil && invocation.ErrorMessage != "" {
				errMsg = utils.StringsBuilder("Tool execution failed: ", invocation.ErrorMessage)
			} else {
				cause := errors.Unwrap(err)
				if cause != nil {
					errMsg = utils.StringsBuilder("Tool execution failed: ", cause.Error())
				} else {
					errMsg = utils.StringsBuilder("Tool execution failed: ", err.Error())
				}
			}
			// Consider 502 Bad Gateway if it's truly an upstream issue vs. 500 Internal Server Error
			httpapi.InternalServerError(c, errMsg)
			return
		}

		// Fallback for other errors from InvokeOperation or if invocation itself indicates failure
		defaultErrMsg := "Failed to call tool."
		if invocation != nil && invocation.Status == integrationDomain.InvocationStatusFailed && invocation.ErrorMessage != "" {
			logger.Warn(ctx, "Tool invocation failed, reported via invocation status after InvokeOperation error", zap.String("invocationID", invocation.ID), zap.String("errorMessage", invocation.ErrorMessage), zap.Error(err))
			defaultErrMsg = utils.StringsBuilder("Tool execution failed: ", invocation.ErrorMessage)
			httpapi.InternalServerError(c, defaultErrMsg)
		} else {
			// Generic fallback if no specific info from invocation
			logger.Error(ctx, "Unhandled error from InvokeOperation", zap.Error(err))
			httpapi.InternalServerError(c, utils.StringsBuilder(defaultErrMsg, " Please check logs."))
		}
		return
	}

	// If err is nil, but invocation status is Failed (should be rare if errors are propagated correctly by InvokeOperation)
	if invocation.Status == integrationDomain.InvocationStatusFailed {
		logger.Warn(ctx, "Tool invocation reported failure status despite no error from InvokeOperation", zap.String("invocationID", invocation.ID), zap.String("errorMessage", invocation.ErrorMessage))
		errMsg := "Tool execution failed"
		if invocation.ErrorMessage != "" {
			errMsg = utils.StringsBuilder(errMsg, ": ", invocation.ErrorMessage)
		}
		httpapi.InternalServerError(c, errMsg)
		return
	}

	// 5. Prepare response if successful
	response := CallToolResponse{
		ToolResult: invocation.ResponseData, // This is json.RawMessage
	}

	logger.Info(ctx, "Successfully called tool", zap.String("invocationID", invocation.ID))
	httpapi.OK(c, response, "Tool called successfully")
}

// HandleMcpListTools godoc
// @Summary List available tools (Provider Operations)
// @Description Retrieves a list of tools that the MCP client can call.
// @Tags mcp
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request_body body ListToolsRequest false "Tool list request parameters"
// @Success 200 {object} httpapi.Response{data=ListToolsResponse} "Success response with a list of tools"
// @Failure 400 {object} httpapi.SwaggerErrorResponse "Bad request if input is invalid"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized if JWT is missing or invalid"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error"
// @Router /mcp/list_tools [post]
func (h *McpHandler) HandleMcpListTools(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.obs.Logger.With(zap.String("handler", "McpHandler"), zap.String("method", "HandleMcpListTools"))

	// 1. Get User from context (set by require_auth middleware)
	userInterface, exists := c.Get("user")
	var userID string // Declare userID here

	if !exists {
		logger.Warn(ctx, "User object not found in context for ListTools. Proceeding with default userID for testing.")
		userID = "test_user_mcp_list_tools" // Default userID for testing when no auth
		// No early return here, allow to proceed for testing
	} else {
		user, ok := userInterface.(*identityDomain.User)
		if !ok || user == nil {
			logger.Error(ctx, "User object in context is not of expected type *identityDomain.User or is nil for ListTools")
			httpapi.InternalServerError(c, "Internal authentication error, invalid user context")
			return
		}
		userID = user.ID // For logging and potential future use
	}
	logger = logger.With(zap.String("userID", userID))

	// 2. Bind request body
	var req ListToolsRequest
	// Use BindJSON for POST with body, or BindQuery for GET with query params.
	// Since it's POST as per plan, BindJSON is appropriate.
	// If the request body is optional, or some fields are, ShouldBindJSON is fine.
	// For an empty body or if an error occurs, it's handled.
	if err := c.ShouldBindJSON(&req); err != nil {
		// Check if it's an EOF error, which means empty body. This is acceptable.
		// For other errors, it's a bad request.
		if errors.Is(err, io.EOF) {
			// EOF means empty JSON body {} or no body, which is fine as fields are optional.
			// Initialize req with defaults if needed, though pointers will be nil.
			logger.Info(ctx, "ListToolsRequest is empty or not provided (EOF), proceeding with defaults.")
		} else {
			// For any other error during binding
			logger.Warn(ctx, "Failed to bind JSON for ListToolsRequest", zap.Error(err))
			httpapi.BadRequest(c, utils.StringsBuilder("Invalid request format: ", err.Error()))
			return
		}
	}

	// Log the received request (being mindful of potentially nil pointers)
	var queryVal, contextVal string
	var allowDisabledVal bool
	if req.Query != nil {
		queryVal = *req.Query
	}
	if req.Context != nil {
		contextVal = *req.Context
	}
	if req.AllowDisabled != nil {
		allowDisabledVal = *req.AllowDisabled
	}
	logger.Info(ctx, "Received ListTools request",
		zap.String("query", queryVal),
		zap.String("context", contextVal),
		zap.Bool("allow_disabled", allowDisabledVal))

	// 3. Call providerService.ListProviders(ctx)
	providers, err := h.providerService.ListProviders(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to list providers from ProviderService", zap.Error(err))
		httpapi.InternalServerError(c, utils.StringsBuilder("Failed to retrieve tool list: ", err.Error()))
		return
	}

	// 4. Transform providers and operations into ToolDefinition list
	toolDefinitions := make([]ToolDefinition, 0)
	allowDisabled := false // Default to false if not provided or explicitly set
	if req.AllowDisabled != nil {
		allowDisabled = *req.AllowDisabled
	}

	for _, provider := range providers {
		// Filter by provider status if allowDisabled is false
		if !allowDisabled && provider.Status != string(types.ProviderStatusActive) { // Used pcdDomain alias
			continue
		}

		for _, operation := range provider.Operations {
			parametersSchema := make(map[string]ToolParameterDefinition)
			for _, param := range operation.Parameters {
				parametersSchema[param.Name] = ToolParameterDefinition{
					Type:        string(param.Type),
					Description: param.Description,
					Required:    param.Required,
					Enum:        param.Enum,
					Default:     param.Default,
				}
			}

			requiredPermissions := make([]string, 0, len(operation.RequiredPermissions))
			for _, perm := range operation.RequiredPermissions {
				requiredPermissions = append(requiredPermissions, perm.Identifier)
			}

			toolDef := ToolDefinition{
				Name:                utils.StringsBuilder(provider.Identifier, ".", operation.Identifier),
				Description:         operation.Description,
				ParametersSchema:    parametersSchema,
				RequiredPermissions: requiredPermissions,
				ProviderName:        provider.Name,
				ProviderIconURL:     provider.IconURL,
			}
			toolDefinitions = append(toolDefinitions, toolDef)
		}
	}

	// 5. Transform toolDefinitions into the new response structure
	finalResponseTools := make(map[string]McpProviderTools)

	for _, toolDef := range toolDefinitions {
		// Extract providerIdentifier and operationIdentifier
		parts := strings.SplitN(toolDef.Name, ".", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			logger.Warn(ctx, "Invalid tool name format in ToolDefinition, skipping", zap.String("toolName", toolDef.Name))
			continue // Skip this malformed tool definition
		}
		providerIdentifier := parts[0]
		operationIdentifier := parts[1]

		// Prepare parameters for this method
		// The logic for requiredParamsList and parametersSchemaForMethod is no longer needed.
		// toolDef.ParametersSchema is already map[string]ToolParameterDefinition

		// Get or create the provider entry in the response map
		providerTools, ok := finalResponseTools[providerIdentifier]
		if !ok {
			providerTools = McpProviderTools{
				Operations: make(map[string]McpToolMethod),
			}
		}

		// Add the method to this provider
		providerTools.Operations[operationIdentifier] = McpToolMethod{
			Description: toolDef.Description,
			Parameters:  toolDef.ParametersSchema, // Directly use toolDef.ParametersSchema
		}

		finalResponseTools[providerIdentifier] = providerTools
	}

	// req.Context filtering is for future enhancement (Phase 4)

	// 6. Prepare and send response
	response := NewListToolsResponse{Tools: finalResponseTools}
	httpapi.OK(c, response, "Tools listed successfully")
}

// Helper function to get keys from a map for logging (to avoid logging sensitive values directly)
func keys(m map[string]interface{}) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}
