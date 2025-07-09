package serper

import (
	"context"
	"fmt"
	"net/http"
)

// Define constants for operation IDs used by handlers.
const (
	operationIDSearch = "search"
	operationIDScrape = "scrape"
)

// Define constants for API paths used by handlers.
const (
	endpointSearch = "/search"
	endpointScrape = "https://scrape.serper.dev"
)

// SearchParams defines parameters for the search operation.
// Based on: POST /search
type SearchParams struct {
	Query       string `mapstructure:"q" validate:"required"`
	Type        string `mapstructure:"type" validate:"omitempty"`               // e.g., search, images, news, videos
	Country     string `mapstructure:"gl" validate:"omitempty"`                 // Country code
	Location    string `mapstructure:"location" validate:"omitempty"`           // More precise location
	Language    string `mapstructure:"hl" validate:"omitempty"`                 // Language code
	Timeframe   string `mapstructure:"tbs" validate:"omitempty"`                // Date range (qdr:h, qdr:d, etc.)
	Autocorrect *bool  `mapstructure:"autocorrect" validate:"omitempty"`        // Pointer for optional boolean
	NumResults  *int   `mapstructure:"num" validate:"omitempty,min=10,max=100"` // Pointer for optional int
	Page        *int   `mapstructure:"page" validate:"omitempty,min=1"`         // Pointer for optional int
}

// ScrapeParams defines parameters for the scrape operation.
// Based on: POST https://scrape.serper.dev
type ScrapeParams struct {
	URL             string `mapstructure:"url" validate:"required,url"`
	IncludeMarkdown *bool  `mapstructure:"includeMarkdown" validate:"omitempty"`
}

// OperationHandler receives the processed parameters (decoded struct pointer)
// and returns a map containing parameters required by the RESTAdapter.
type OperationHandler func(ctx context.Context, processedParams interface{}) (map[string]interface{}, error)

// OperationDefinition combines the parameter schema and handler for an operation.
type OperationDefinition struct {
	Schema  interface{}      // Parameter schema (struct pointer, e.g., &SearchParams{})
	Handler OperationHandler // Operation handler function
}

// Operations maps operation IDs to their corresponding definitions.
type Operations map[string]OperationDefinition

// RegisterOperation registers the parameter schema and handler for an operation.
func (a *SerperAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler) {
	a.BaseAdapter.RegisterOperation(operationID, schema)
	if a.operations == nil {
		a.operations = make(Operations)
	}
	a.operations[operationID] = OperationDefinition{
		Schema:  schema,
		Handler: handler,
	}
}

// registerOperations populates the operations map in the SerperAdapter.
func (a *SerperAdapter) registerOperations() {
	// Register search operation
	a.RegisterOperation(
		operationIDSearch,
		&SearchParams{},
		handleSearch,
	)

	// Register scrape operation
	a.RegisterOperation(
		operationIDScrape,
		&ScrapeParams{},
		handleScrape,
	)
}

// handleSearch prepares parameters for the Serper search operation.
func handleSearch(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*SearchParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDSearch)
	}

	// Construct the request body map from the SearchParams struct
	body := make(map[string]interface{})
	body["q"] = params.Query
	if params.Type != "" {
		body["type"] = params.Type
	}
	if params.Country != "" {
		body["gl"] = params.Country
	}
	if params.Location != "" {
		body["location"] = params.Location
	}
	if params.Language != "" {
		body["hl"] = params.Language
	}
	if params.Timeframe != "" {
		body["tbs"] = params.Timeframe
	}
	if params.Autocorrect != nil {
		body["autocorrect"] = *params.Autocorrect
	}
	if params.NumResults != nil {
		body["num"] = *params.NumResults
	}
	if params.Page != nil {
		body["page"] = *params.Page
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointSearch,
		"body":   body,
		// Serper API expects Content-Type: application/json, which should be handled by restAdapter or serperAdapter New method
	}
	return restParams, nil
}

// handleScrape prepares parameters for the Serper scrape operation.
func handleScrape(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*ScrapeParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDScrape)
	}

	body := make(map[string]interface{})
	body["url"] = params.URL
	if params.IncludeMarkdown != nil {
		body["includeMarkdown"] = *params.IncludeMarkdown
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointScrape,
		"body":   body,
	}
	return restParams, nil
}
