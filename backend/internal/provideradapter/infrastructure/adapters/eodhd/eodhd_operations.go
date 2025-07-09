package eodhd

import (
	"context"
	"fmt"
	"net/http"
)

// Define operation IDs as constants
const (
	operationIDGetUser                   = "get_user"
	operationIDGetEODHistoricalData      = "get_eod_historical_data"
	operationIDGetFundamentalData        = "get_fundamental_data"
	operationIDGetIntradayHistoricalData = "get_intraday_historical_data"
	operationIDGetRealTimeData           = "get_real_time_data"
	operationIDGetExchanges              = "get_exchanges"
)

// GetEODHistoricalDataParams defines parameters for the EOD historical data operation.
// Based on: GET /eod/{symbol}
type GetEODHistoricalDataParams struct {
	Symbol string `mapstructure:"symbol" validate:"required"`                    // Path parameter, handled in handler
	From   string `mapstructure:"from" validate:"omitempty,datetime=2006-01-02"` // Query parameter
	To     string `mapstructure:"to" validate:"omitempty,datetime=2006-01-02"`   // Query parameter
	Period string `mapstructure:"period" validate:"omitempty,oneof=d w m"`       // Query parameter, defaults to 'd' in handler
	Order  string `mapstructure:"order" validate:"omitempty,oneof=a d"`          // Query parameter, defaults to 'a' in handler
}

// GetFundamentalDataParams defines parameters for the fundamental data operation.
// Based on: GET /fundamentals/{symbol}
type GetFundamentalDataParams struct {
	Symbol string `mapstructure:"symbol" validate:"required"`  // Path parameter
	Filter string `mapstructure:"filter" validate:"omitempty"` // Query parameter
}

// GetIntradayHistoricalDataParams defines parameters for the intraday historical data operation.
// Based on: GET /intraday/{symbol}
type GetIntradayHistoricalDataParams struct {
	Symbol   string `mapstructure:"symbol" validate:"required"`                  // Path parameter
	Interval string `mapstructure:"interval" validate:"required,oneof=1m 5m 1h"` // Query parameter
	From     string `mapstructure:"from" validate:"omitempty"`                   // Query parameter (Unix timestamp or YYYY-MM-DD HH:MM:SS)
	To       string `mapstructure:"to" validate:"omitempty"`                     // Query parameter (Unix timestamp or YYYY-MM-DD HH:MM:SS)
}

// GetRealTimeDataParams defines parameters for the live/delayed data operation.
// Based on: GET /real-time/{symbol}
type GetRealTimeDataParams struct {
	Symbol string `mapstructure:"symbol" validate:"required"` // Path parameter (can be single or comma-separated symbols)
}

// OperationHandler now receives the processed parameters (decoded struct pointer)
// It returns a map containing parameters required by the RESTAdapter.
type OperationHandler func(ctx context.Context, processedParams interface{}) (map[string]interface{}, error)

// OperationDefinition combines the parameter schema and handler for an operation
type OperationDefinition struct {
	Schema  interface{}      // Parameter schema (struct pointer, e.g., &GetUserParams{})
	Handler OperationHandler // Operation handler function
}

// Operations maps operation IDs (e.g., "eodhd.get_user") to their corresponding definitions.
type Operations map[string]OperationDefinition // Modified type

// RegisterOperation registers the parameter schema and handler for an operation
// It is now a method on EodhdAdapter, defined in this file.
func (a *EodhdAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler) {
	// Register the parameter schema with the base adapter for validation/decoding
	a.BaseAdapter.RegisterOperation(operationID, schema)

	// Store the definition in the EodhdAdapter's operations map
	if a.operations == nil {
		a.operations = make(Operations) // Ensure map is initialized
	}
	a.operations[operationID] = OperationDefinition{
		Schema:  schema,
		Handler: handler,
	}
}

// registerOperations populates the operations map in the EodhdAdapter.
// It is now a method on EodhdAdapter, defined in this file.
func (a *EodhdAdapter) registerOperations() {
	// Register get_user operation
	a.RegisterOperation(
		operationIDGetUser,
		&struct{}{},
		handleGetUser,
	)

	// Register get_eod_historical_data operation
	a.RegisterOperation(
		operationIDGetEODHistoricalData,
		&GetEODHistoricalDataParams{},
		handleGetEODHistoricalData,
	)

	// Register get_fundamental_data operation
	a.RegisterOperation(
		operationIDGetFundamentalData,
		&GetFundamentalDataParams{},
		handleGetFundamentalData,
	)

	// Register get_intraday_historical_data operation
	a.RegisterOperation(
		operationIDGetIntradayHistoricalData,
		&GetIntradayHistoricalDataParams{},
		handleGetIntradayHistoricalData,
	)

	// Register get_real_time_data operation
	a.RegisterOperation(
		operationIDGetRealTimeData,
		&GetRealTimeDataParams{},
		handleGetRealTimeData,
	)

	// Register get_exchanges operation
	a.RegisterOperation(
		operationIDGetExchanges,
		&struct{}{},
		handleGetExchanges,
	)
}

func handleGetUser(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   "user",
	}
	return restParams, nil
}

func handleGetEODHistoricalData(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*GetEODHistoricalDataParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDGetEODHistoricalData)
	}

	// Build query parameters, always adding fmt=json
	queryParams := map[string]string{
		"fmt": "json", // Always request JSON
	}
	if params.From != "" {
		queryParams["from"] = params.From
	}
	if params.To != "" {
		queryParams["to"] = params.To
	}
	// Set default if not provided
	if params.Period != "" {
		queryParams["period"] = params.Period
	} else {
		queryParams["period"] = "d" // Default period
	}
	if params.Order != "" {
		queryParams["order"] = params.Order
	} else {
		queryParams["order"] = "a" // Default order
	}

	// Construct the path using the Symbol parameter
	// Note: URL escaping for the symbol might be needed depending on the RESTAdapter implementation
	path := fmt.Sprintf("eod/%s", params.Symbol)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         path,
		"query_params": queryParams,
		// No body needed for GET
	}

	return restParams, nil
}

func handleGetFundamentalData(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*GetFundamentalDataParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDGetFundamentalData)
	}

	queryParams := map[string]string{
		"fmt": "json", // Always request JSON
	}
	if params.Filter != "" {
		queryParams["filter"] = params.Filter
	}

	path := fmt.Sprintf("fundamentals/%s", params.Symbol)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         path,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleGetIntradayHistoricalData(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*GetIntradayHistoricalDataParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDGetIntradayHistoricalData)
	}

	queryParams := map[string]string{
		"fmt":      "json", // Always request JSON
		"interval": params.Interval,
	}
	if params.From != "" {
		queryParams["from"] = params.From
	}
	if params.To != "" {
		queryParams["to"] = params.To
	}

	path := fmt.Sprintf("intraday/%s", params.Symbol)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         path,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleGetRealTimeData(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*GetRealTimeDataParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDGetRealTimeData)
	}

	queryParams := map[string]string{
		"fmt": "json", // Always request JSON
	}

	// Symbol in path can contain multiple comma-separated values according to docs
	path := fmt.Sprintf("real-time/%s", params.Symbol)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         path,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleGetExchanges(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	// No specific user parameters needed for this operation

	queryParams := map[string]string{
		"fmt": "json", // Always request JSON
	}

	path := "exchanges-list/"

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         path,
		"query_params": queryParams,
	}

	return restParams, nil
}
