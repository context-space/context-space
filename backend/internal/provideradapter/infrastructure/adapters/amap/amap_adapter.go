package amap

import (
	"context"
	"fmt"
	"net/http"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
)

// AmapAdapter is the adapter for the Amap API
type AmapAdapter struct {
	*base.BaseAdapter
	restAdapter domain.Adapter // The underlying REST adapter instance
	operations  Operations     // Map of operation ID to definition defined in amap_operations.go
}

// NewAmapAdapter creates a new Amap adapter
func NewAmapAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	restAdapter domain.Adapter,
) *AmapAdapter {
	// Create AmapAdapter instance
	baseAdapter := base.NewBaseAdapter(providerInfo, config)

	adapter := &AmapAdapter{
		BaseAdapter: baseAdapter,
		restAdapter: restAdapter,
		operations:  make(Operations),
	}

	adapter.registerOperations()
	return adapter
}

// Execute handles API calls based on operationID using the REST adapter
func (a *AmapAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{},
	credential interface{},
) (interface{}, error) {
	// Extract API Key from credential
	apiKeyCred, ok := credential.(*credDomain.APIKeyCredential)
	if !ok || apiKeyCred == nil || apiKeyCred.APIKey == "" {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			domain.ErrCredentialError,
			"invalid or missing API key credential",
			http.StatusUnauthorized)
	}

	// 1. Check if the operation exists
	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			domain.ErrOperationNotSupported,
			fmt.Sprintf("Unknown operation ID: %s", operationID),
			http.StatusNotFound,
		)
	}

	// 2.Validate and process parameters
	processedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			domain.ErrInvalidParameters,
			fmt.Sprintf("Parameter validation failed: %s", err.Error()),
			http.StatusBadRequest,
		)
	}

	// 3. Call the operation handler
	handler := opDef.Handler
	restParams, err := handler(ctx, processedParams)
	if err != nil {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			domain.ErrInternal,
			fmt.Sprintf("Handler execution failed: %s", err.Error()),
			http.StatusInternalServerError,
		)
	}

	// 4. Add API Key to query parameters (Amap API uses query parameters to pass the key)
	queryParams, _ := restParams["query_params"].(map[string]string)
	if queryParams == nil {
		queryParams = make(map[string]string)
	}
	queryParams[apikeyParamName] = apiKeyCred.APIKey
	restParams["query_params"] = queryParams

	// 5. Execute REST call
	rawResult, err := a.restAdapter.Execute(ctx, operationID, restParams, nil)
	if err != nil {
		return nil, err
	}

	return rawResult, nil
}
