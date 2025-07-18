package mcp

import (
	"context"
)

// MCPOperationParams represents parameters for MCP operations
type MCPOperationParams struct {
	Parameters map[string]interface{} `json:"parameters"`
}

// OperationHandler defines the function signature for handling MCP operations
type OperationHandler func(ctx context.Context, processedParams interface{}) (interface{}, error)

// OperationDefinition combines the parameter schema and handler function
type OperationDefinition struct {
	Schema  interface{}      // Parameter schema
	Handler OperationHandler // Operation handler function
}

// Operations maps operation IDs to their definitions
type Operations map[string]OperationDefinition
