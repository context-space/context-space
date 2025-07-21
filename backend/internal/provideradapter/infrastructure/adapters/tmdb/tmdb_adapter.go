package tmdb

import (
	"context"
	"fmt"
	"net/http"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
)

// TmdbAdapter implements the provider adapter for The Movie Database (TMDB) API
type TmdbAdapter struct {
	*base.BaseAdapter
	restAdapter domain.Adapter // Uses RESTAdapter for actual execution
	operations  Operations     // Map operation ID to OperationDefinition
}

// NewTmdbAdapter creates a new TMDB adapter instance
func NewTmdbAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	restAdapter domain.Adapter,
) *TmdbAdapter {
	// Create TMDB adapter instance
	baseAdapter := base.NewBaseAdapter(providerInfo, config)
	adapter := &TmdbAdapter{
		BaseAdapter: baseAdapter,
		restAdapter: restAdapter,
		operations:  make(Operations),
	}

	// Register specific TMDB operations
	adapter.registerOperations()

	return adapter
}

// Execute finds the appropriate handler for the operationID, prepares parameters
// including the API key, and delegates execution to the underlying RESTAdapter.
func (a *TmdbAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{}, // Original user parameters
	credential interface{}, // Expected to be *credDomain.APIKeyCredential
) (interface{}, error) {
	apiKeyCred, ok := credential.(*credDomain.APIKeyCredential)
	if !ok || apiKeyCred == nil || apiKeyCred.APIKey == "" {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrCredentialError, "invalid or missing API key credential", http.StatusUnauthorized)
	}

	// 1. Find the operation handler definition
	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "", fmt.Sprintf("unknown operation ID: %s", operationID), http.StatusNotFound)
	}

	// 2. Process and validate parameters
	processedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "PARAMETER_ERROR", fmt.Sprintf("parameter validation failed: %s", err.Error()), http.StatusBadRequest)
	}

	// 3. Call the operation handler
	handler := opDef.Handler
	restParams, err := handler(ctx, processedParams)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "HANDLER_ERROR", fmt.Sprintf("handler execution failed: %s", err.Error()), http.StatusInternalServerError)
	}

	// 4. Add API Key to query parameters (TMDB API uses query parameters)
	queryParams, _ := restParams["query_params"].(map[string]string)
	if queryParams == nil {
		queryParams = make(map[string]string)
	}
	queryParams[apikeyParamName] = apiKeyCred.APIKey
	restParams["query_params"] = queryParams

	// 5. Delegate execution to REST adapter
	// Pass nil credential to restAdapter since authentication is handled here
	result, err := a.restAdapter.Execute(ctx, operationID, restParams, nil)
	if err != nil {
		// Error should already be wrapped by RESTAdapter or domain.AdapterError
		return nil, err
	}

	return result, nil
}
