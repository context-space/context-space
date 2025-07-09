package eodhd

import (
	"context"
	"fmt"
	"net/http"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
)

// EodhdAdapter is an adapter for the EODHD API
type EodhdAdapter struct {
	*base.BaseAdapter
	restAdapter domain.Adapter
	operations  Operations
}

// NewEodhdAdapter creates a new EodhdAdapter
func NewEodhdAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	restAdapter domain.Adapter,
) *EodhdAdapter {

	// Create the EodhdAdapter instance, using the passed-in apiKeyCfg
	baseAdapter := base.NewBaseAdapter(providerInfo, config)
	adapter := &EodhdAdapter{
		BaseAdapter: baseAdapter,
		restAdapter: restAdapter,
		operations:  make(Operations),
	}

	// Register specific EODHD operations
	adapter.registerOperations()

	return adapter
}

// Execute finds the appropriate handler for the operationID, prepares parameters
// including the API key, and delegates execution to the underlying RESTAdapter.
func (a *EodhdAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{},
	credential interface{},
) (interface{}, error) {

	cred, ok := credential.(*credDomain.APIKeyCredential)

	if !ok || cred == nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "", "invalid or missing API key credential", http.StatusUnauthorized)
	}

	apiKey := cred.APIKey

	if apiKey == "" {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "", "invalid or missing API key", http.StatusUnauthorized)
	}

	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "", fmt.Sprintf("unknown operation ID: %s", operationID), http.StatusNotFound)
	}

	// BaseAdapter.ProcessParams uses the registered schema for the operationID.
	processedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		// Handle parameter processing error (e.g., validation failed)
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "PARAMETER_ERROR", fmt.Sprintf("parameter validation failed: %s", err.Error()), http.StatusBadRequest)
	}

	handler := opDef.Handler
	restParams, err := handler(ctx, processedParams) // Pass processedParams (struct pointer)
	if err != nil {
		// Consider adding specific error code from handler if possible
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "HANDLER_ERROR", fmt.Sprintf("handler execution failed: %s", err.Error()), http.StatusInternalServerError)
	}

	// Inject the API key into the restParams based on apiKeyConfig
	queryParams, _ := restParams["query_params"].(map[string]string)
	if queryParams == nil {
		queryParams = make(map[string]string)
	}
	queryParams[apikeyParamName] = apiKey
	restParams["query_params"] = queryParams

	// Delegate execution to the REST adapter
	// The credential passed to restAdapter is nil because authentication is handled here.
	result, err := a.restAdapter.Execute(ctx, operationID, restParams, nil)
	if err != nil {
		// Error should already be wrapped by RESTAdapter or domain.AdapterError
		return nil, err
	}

	return result, nil
}
