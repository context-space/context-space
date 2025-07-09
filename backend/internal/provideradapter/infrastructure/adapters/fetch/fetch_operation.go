package fetch

import (
	"context"
	"fmt"
	"net/http"
)

// API path constants
const (
	endpointFetch = "/" // Serper Scraping API root path
)

// Define constants for operation IDs used by handlers.
const (
	operationIDFetchContent = "fetch_content"
)

var opDefaults = OperationDefaults{
	FetchContent: FetchDefaults{
		IncludeMarkdown: true,
	},
}

// FetchDefaults stores default values for 'fetch.fetch_content' operation
type FetchDefaults struct {
	IncludeMarkdown bool // Whether to include Markdown format in the response
}

// OperationDefaults stores default parameter values for various Fetch operations
// These are typically configured once when the adapter is created
type OperationDefaults struct {
	FetchContent FetchDefaults // Default values for fetch_content operation
}

// FetchParams defines user parameters for the fetch_content operation
// Only include parameters that the user must provide or may want to customize
// Other parameters are automatically set by the adapter's defaults
type FetchParams struct {
	URL string `mapstructure:"url" validate:"required"` // Website URL to fetch (required)
}

// OperationHandler receives processed parameters and returns REST parameters
type OperationHandler func(ctx context.Context, processedParams interface{}) (map[string]interface{}, error)

// OperationDefinition combines parameter structure and handler
type OperationDefinition struct {
	Schema  interface{}      // Parameter structure (struct pointer, e.g., &FetchParams{})
	Handler OperationHandler // Operation processing function
}

// Operations maps operation IDs to their definitions
type Operations map[string]OperationDefinition

// RegisterOperation registers operation's parameter structure and handler
func (a *FetchAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler) {
	// Register parameter structure in base adapter for validation/decoding
	a.BaseAdapter.RegisterOperation(operationID, schema)

	// Store definition in FetchAdapter's operations map
	if a.operations == nil {
		a.operations = make(Operations)
	}
	a.operations[operationID] = OperationDefinition{
		Schema:  schema,
		Handler: handler,
	}
}

// registerOperations populates the operations map
func (a *FetchAdapter) registerOperations() {
	// Register fetch.fetch_content operation
	a.RegisterOperation(
		operationIDFetchContent,
		&FetchParams{},
		handleFetchContent,
	)
}

func handleFetchContent(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*FetchParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDFetchContent)
	}

	// Build request body
	body := map[string]interface{}{
		"url":             params.URL, // User-provided URL
		"includeMarkdown": opDefaults.FetchContent.IncludeMarkdown,
	}

	// Create REST parameters according to Serper Scraping API
	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointFetch,
		"body":   body,
		"headers": map[string]string{
			"Content-Type": "application/json",
		},
	}

	return restParams, nil
}
