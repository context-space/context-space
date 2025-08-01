package mcp

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const credentialKey contextKey = "credential"

// MCPAdapter implements base adapter for MCP (Model Context Protocol) providers
type MCPAdapter struct {
	*base.BaseAdapter
	config        *MCPAdapterConfig
	configBuilder *MCPConfigBuilder // Add config builder
	tools         []MCPTool
	operations    Operations // Keep operations field
	once          sync.Once
	initError     error
	mu            sync.RWMutex
}

// MCPAdapterConfig represents the configuration for MCP adapter in manifest.json
type MCPAdapterConfig struct {
	// Connection configuration
	Command            string            `json:"command"`             // Base command to run MCP server
	Args               []string          `json:"args"`                // Arguments for the command
	Envs               map[string]string `json:"envs"`                // Environment variables for server
	Timeout            time.Duration     `json:"timeout"`             // Timeout for MCP operations
	CredentialMappings map[string]string `json:"credential_mappings"` // Maps credential keys to environment variables
	DummyCredentials   map[string]string `json:"dummy_credentials"`   // Dummy values for initialization
	ParameterMappings  map[string]string `json:"parameter_mappings"`  // Maps parameter keys to environment variables
	DummyParameters    map[string]string `json:"dummy_parameters"`    // Default values for parameters
}

// UnmarshalJSON provides custom JSON unmarshaling for MCPAdapterConfig
func (c *MCPAdapterConfig) UnmarshalJSON(data []byte) error {
	// Define a temporary struct with timeout as interface{}
	type TempConfig struct {
		Command              string            `json:"command"`
		Args                 []string          `json:"args"`
		Envs                 map[string]string `json:"envs"`
		Timeout              interface{}       `json:"timeout"`
		CredentialMappings   map[string]string `json:"credential_mappings"`
		DummyCredentials     map[string]string `json:"dummy_credentials"`
		ParameterMappings    map[string]string `json:"parameter_mappings"`
		DummyParameters      map[string]string `json:"dummy_parameters"`
		CustomConfigTemplate map[string]string `json:"custom_config"`
	}

	var temp TempConfig
	if err := sonic.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Copy simple fields
	c.Command = temp.Command
	c.Args = temp.Args
	c.Envs = temp.Envs
	c.CredentialMappings = temp.CredentialMappings
	c.DummyCredentials = temp.DummyCredentials
	c.ParameterMappings = temp.ParameterMappings
	c.DummyParameters = temp.DummyParameters

	// Handle timeout conversion
	if temp.Timeout != nil {
		switch v := temp.Timeout.(type) {
		case string:
			if duration, err := time.ParseDuration(v); err == nil {
				c.Timeout = duration
			}
		case float64:
			// Handle numeric values as seconds
			c.Timeout = time.Duration(v) * time.Second
		case int:
			c.Timeout = time.Duration(v) * time.Second
		case int64:
			c.Timeout = time.Duration(v) * time.Second
		}
	}

	return nil
}

// NewMCPAdapter creates a new MCP adapter
func NewMCPAdapter(baseAdapter *base.BaseAdapter, config *MCPAdapterConfig) *MCPAdapter {
	adapter := &MCPAdapter{
		BaseAdapter:   baseAdapter,
		config:        config,
		configBuilder: NewMCPConfigBuilder(config), // Initialize config builder
		tools:         []MCPTool{},
		operations:    make(Operations),
	}
	return adapter
}

// initializeOperations discovers tools from MCP server and registers them as operations
func (a *MCPAdapter) initializeOperations() error {
	tools, err := a.getTools(nil) // Use dummy values for initialization
	if err != nil {
		return fmt.Errorf("failed to get tools from MCP server: %w", err)
	}

	a.registerOperationsFromTools(tools)
	return nil
}

// getTools retrieves tools from MCP server (only once, no expiration)
func (a *MCPAdapter) getTools(credential interface{}) ([]MCPTool, error) {
	// Fast path: check if tools are already cached (read lock)
	a.mu.RLock()
	if len(a.tools) > 0 { // Changed from toolsInitialized to tools
		tools := a.tools
		a.mu.RUnlock()
		return tools, nil
	}
	a.mu.RUnlock()

	// Slow path: need to initialize tools (write lock)
	a.mu.Lock()
	defer a.mu.Unlock()

	// Double-check pattern: verify tools are still not initialized
	if len(a.tools) > 0 { // Changed from toolsInitialized to tools
		return a.tools, nil
	}

	// Create MCP client configuration with dynamic credential injection
	clientConfig := a.buildMCPClientConfig(credential, nil) // Pass nil for parameters during initialization

	// Get tools from MCP server
	tools, err := GetMCPServerTools(context.Background(), clientConfig, a.GetProviderAdapterInfo().Identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get MCP server tools: %w", err)
	}

	// Cache the results permanently
	a.tools = tools

	return tools, nil
}

// buildMCPClientConfig creates MCP client configuration with credential and parameter injection
func (a *MCPAdapter) buildMCPClientConfig(credential interface{}, parameters map[string]interface{}) MCPClientConfig {
	return a.configBuilder.Build(credential, parameters)
}

// registerOperationsFromTools registers MCP tools as operations
func (a *MCPAdapter) registerOperationsFromTools(tools []MCPTool) {
	// Since this is called during initialization (protected by sync.Once),
	// the write lock should already be held by getTools, but we'll be safe
	for _, tool := range tools {
		// Create operation handler for this tool
		handler := a.createOperationHandler(tool)

		// Register operation with parameter schema (thread-safe)
		a.registerOperation(tool.Name, &MCPOperationParams{}, handler)
	}
}

// registerOperation is a thread-safe version of RegisterOperation for MCP operations
func (a *MCPAdapter) registerOperation(operationID string, schema interface{}, handler OperationHandler) {
	// Initialize operations map if nil
	if a.operations == nil {
		a.operations = make(Operations)
	}

	// Register the operation
	a.operations[operationID] = OperationDefinition{
		Schema:  schema,
		Handler: handler,
	}

	// Also register with base adapter for compatibility
	a.BaseAdapter.RegisterOperation(operationID, schema)
}

// createOperationHandler creates a handler function for a specific MCP tool
func (a *MCPAdapter) createOperationHandler(tool MCPTool) OperationHandler {
	return func(ctx context.Context, processedParams interface{}) (interface{}, error) {
		params, ok := processedParams.(*MCPOperationParams)
		if !ok {
			return nil, fmt.Errorf("invalid parameter type for MCP operation")
		}

		// Get credential from context or adapter state
		credential := ctx.Value(credentialKey)

		// Create MCP client configuration with credential injection
		clientConfig := a.buildMCPClientConfig(credential, params.Parameters)

		// Call the MCP operation
		result, err := CallMCPOperation(
			ctx,
			clientConfig,
			tool.Name,
			params.Parameters,
			a.GetProviderAdapterInfo().Identifier,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to call MCP operation %s: %w", tool.Name, err)
		}

		return result, nil
	}
}

// Execute handles the execution of a specific MCP operation
func (a *MCPAdapter) Execute(ctx context.Context, operationID string, parameters map[string]interface{}, credential interface{}) (interface{}, error) {
	// Thread-safe initialization using sync.Once
	a.once.Do(func() {
		a.initError = a.initializeOperations()
	})

	// Check if initialization failed
	if a.initError != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInternal, fmt.Sprintf("failed to initialize operations: %v", a.initError), http.StatusInternalServerError)
	}

	// Get operation definition with read lock
	a.mu.RLock()
	operation, exists := a.operations[operationID]
	a.mu.RUnlock()

	if !exists {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrOperationNotSupported, fmt.Sprintf("operation %s not found", operationID), http.StatusNotFound)
	}

	// Process parameters
	processedParams := &MCPOperationParams{
		Parameters: parameters,
	}

	// Validate parameters against the operation's input schema
	if err := a.validateParameters(operationID, parameters, credential); err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInvalidParameters, fmt.Sprintf("parameter validation failed: %v", err), http.StatusBadRequest)
	}

	// Add credential to context for handler access
	ctx = context.WithValue(ctx, credentialKey, credential)

	// Execute the operation
	result, err := operation.Handler(ctx, processedParams)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrProviderAPIError, fmt.Sprintf("operation execution failed: %v", err), http.StatusInternalServerError)
	}

	return result, nil
}

// validateParameters validates operation parameters against the tool's input schema
func (a *MCPAdapter) validateParameters(operationID string, parameters map[string]interface{}, credential interface{}) error {
	// Get the tool definition for validation
	tools, err := a.getTools(credential)
	if err != nil {
		return fmt.Errorf("failed to get tools for validation: %w", err)
	}

	var tool *MCPTool
	for _, t := range tools {
		if t.Name == operationID {
			tool = &t
			break
		}
	}

	if tool == nil {
		return fmt.Errorf("tool %s not found", operationID)
	}

	// Basic validation - check required parameters
	for _, required := range tool.InputSchema.Required {
		if _, exists := parameters[required]; !exists {
			return fmt.Errorf("required parameter %s is missing", required)
		}
	}

	// Additional validation can be added here based on the input schema
	return nil
}

// RegisterOperation registers an operation with the adapter
func (a *MCPAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler) {
	a.BaseAdapter.RegisterOperation(operationID, schema)

	if a.operations == nil {
		a.operations = make(Operations)
	}

	a.operations[operationID] = OperationDefinition{
		Schema:  schema,
		Handler: handler,
	}
}

// ListOperations returns a list of available operations
func (a *MCPAdapter) ListOperations() []string {
	tools, err := a.getTools(nil) // Use dummy values for listing
	if err != nil {
		return []string{}
	}

	operations := make([]string, len(tools))
	for i, tool := range tools {
		operations[i] = tool.Name
	}

	return operations
}

// GetOperationSchema returns the schema for a specific operation
func (a *MCPAdapter) GetOperationSchema(operationID string) (interface{}, error) {
	tools, err := a.getTools(nil) // Use dummy values for schema retrieval
	if err != nil {
		return nil, fmt.Errorf("failed to get tools: %w", err)
	}

	for _, tool := range tools {
		if tool.Name == operationID {
			return tool.InputSchema, nil
		}
	}

	return nil, fmt.Errorf("operation %s not found", operationID)
}

// Health checks the health of the MCP server
func (a *MCPAdapter) Health(ctx context.Context) error {
	// Try to get tools as a health check
	if a.initError != nil {
		return fmt.Errorf("MCP server initialization failed: %w", a.initError)
	}
	return nil
}

// GetProviderAdapterInfo returns provider information (implements domain.Adapter interface)
func (a *MCPAdapter) GetProviderAdapterInfo() *domain.ProviderAdapterInfo {
	return a.BaseAdapter.GetProviderAdapterInfo()
}
