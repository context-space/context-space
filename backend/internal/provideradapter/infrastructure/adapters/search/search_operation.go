package search

import (
	"context"
	"fmt"
	"net/http"
)

// Define constants for API paths used by handlers.
const (
	endpointSearch = "search"
)

// Define constants for operation IDs used by handlers.
const (
	operationIDSearch = "search"
)

var operationDefaults = OperationDefaults{
	Search: SearchDefaults{
		Type:        "search",
		Country:     "us",
		Language:    "en",
		AutoCorrect: true,
		Num:         10,
		Page:        1,
	},
}

// SearchDefaults holds default values for 'serper.search' operation
type SearchDefaults struct {
	Type        string // Search type (search, images, news, videos)
	Country     string // Country code
	Language    string // Language code
	AutoCorrect bool   // Whether to auto-correct
	Num         int    // Number of results
	Page        int    // Page number
}

// OperationDefaults holds default parameter values for various Serper operations
// These are typically configured once when the adapter is created
type OperationDefaults struct {
	Search SearchDefaults // Default values for search operation
}

// SearchParams defines user parameters for the serper.search operation
type SearchParams struct {
	Query     string `mapstructure:"query" validate:"required"`
	DataRange string `mapstructure:"data_range" validate:"omitempty,oneof=hour day week month year"`
}

// OperationHandler receives processed parameters and returns REST parameters
type OperationHandler func(ctx context.Context, processedParams interface{}) (map[string]interface{}, error)

// OperationDefinition combines parameter structure and handler
type OperationDefinition struct {
	Schema  interface{}
	Handler OperationHandler
}

// Operations maps operation IDs to their definitions
type Operations map[string]OperationDefinition

// RegisterOperation registers operation's parameter structure and handler
func (a *SearchAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler) {
	a.BaseAdapter.RegisterOperation(operationID, schema)

	if a.operations == nil {
		a.operations = make(Operations)
	}
	a.operations[operationID] = OperationDefinition{
		Schema:  schema,
		Handler: handler,
	}
}

// registerOperations populates the operations map
func (a *SearchAdapter) registerOperations() {
	a.RegisterOperation(
		operationIDSearch,
		&SearchParams{},
		handleSearch,
	)
}

func handleSearch(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*SearchParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDSearch)
	}

	body := map[string]interface{}{
		"q":           params.Query,
		"type":        operationDefaults.Search.Type,
		"gl":          operationDefaults.Search.Country,
		"hl":          operationDefaults.Search.Language,
		"autocorrect": operationDefaults.Search.AutoCorrect,
		"num":         operationDefaults.Search.Num,
		"page":        operationDefaults.Search.Page,
	}

	if params.DataRange != "" {
		body["tbs"] = fmt.Sprintf("qdr:%s", params.DataRange[0:1])
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointSearch,
		"body":   body,
		"headers": map[string]string{
			"Content-Type": "application/json",
		},
	}

	return restParams, nil
}
