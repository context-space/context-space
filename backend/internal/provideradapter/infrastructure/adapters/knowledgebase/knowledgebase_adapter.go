package knowledgebase

import (
	"context"
	"fmt"
	"net/http"

	domain "github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"

	openaiclient "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/knowledgebase/openai/client"
	volcclient "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/knowledgebase/volcengine/client"
)

var DefaultKnowledgebaseAdapterConfig = KnowledgebaseAdapterConfig{
	Search: &SearchConfig{
		Limit: intPtr(10),
	},
	Chat: &ChatConfig{
		Model:       stringPtr("gpt-4o"),
		Stream:      boolPtr(false),
		Temperature: float64Ptr(0.1),
	},
	Query: &QueryConfig{
		SearchLimit:         intPtr(10),
		RewriteQuery:        boolPtr(false),
		Rerank:              boolPtr(false),
		RerankRetrieveCount: intPtr(50),
		RerankModel:         nil,
		LLMModel:            stringPtr("gpt-4o"),
		LLMTemperature:      float64Ptr(0.1),
	},
}

type KnowledgebaseAdapterConfig struct {
	Project        string        `json:"project"`
	CollectionName string        `json:"collection_name"`
	Search         *SearchConfig `json:"search"`
	Chat           *ChatConfig   `json:"chat"`
	Query          *QueryConfig  `json:"query"`
}

// KnowledgeBaseAdapter is the user-facing adapter for Volcengine Knowledge Base operations.
// It abstracts the specific Volcengine API calls behind a common interface.
type KnowledgeBaseAdapter struct {
	*base.BaseAdapter
	internalVolcengineClient volcclient.VolcengineClient
	operations               Operations // Defined in knowledgebase_operations.go
	openaiClient             openaiclient.OpenaiClient
	baseConfig               *KnowledgebaseAdapterConfig // Grouped default parameters
}

// NewKnowledgeBaseAdapter creates a new instance of the KnowledgeBaseAdapter.
func NewKnowledgeBaseAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	internalClient volcclient.VolcengineClient,
	openaiClient openaiclient.OpenaiClient,
	baseConfig *KnowledgebaseAdapterConfig,
) (*KnowledgeBaseAdapter, error) { // Return error for validation

	baseAdapter := base.NewBaseAdapter(providerInfo, config)
	adapter := &KnowledgeBaseAdapter{
		BaseAdapter:              baseAdapter,
		internalVolcengineClient: internalClient,
		openaiClient:             openaiClient,
		operations:               make(Operations),
		baseConfig:               baseConfig,
	}
	adapter.registerOperations()

	return adapter, nil
}

// Execute handles the execution of a specific Volcengine Knowledge Base operation.
func (a *KnowledgeBaseAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{}, // User-provided parameters
	credential interface{}, // Should be *volcenginetypes.VolcengineCredential, can be nil now
) (interface{}, error) {
	// 1. Find Operation Definition (Schema and Handler)
	opDef, exists := a.operations[operationID]
	if !exists {
		// Operation not defined/supported by this adapter
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "OPERATION_NOT_FOUND", fmt.Sprintf("unknown operation ID: %s", operationID), http.StatusNotFound)
	}

	// 2. Process User Parameters (Validation and Type Conversion)
	// This uses the schema found in opDef from step 1.
	processedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		// ProcessParams error likely means bad user input (failed validation or decoding)
		// Ensure this returns PARAMETER_ERROR and Bad Request status.
		// The error message from ProcessParams should contain validation details.
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "PARAMETER_ERROR", fmt.Sprintf("parameter processing failed: %v", err), http.StatusBadRequest)
	}

	// 3. Call Handler to get API Request Details or Final Result
	// The handler takes the user-processed parameters.
	// It returns either *OperationHandlerOutput for single Volcengine calls
	// OR the final result for complex operations like RAG.
	handlerResult, err := opDef.Handler(ctx, processedParams)
	if err != nil {
		// Error from the handler itself (e.g., failed mapping, internal search/LLM call failed)
		// Try to return an AdapterError if possible
		if ae, ok := err.(*domain.AdapterError); ok {
			return nil, ae
		}
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "HANDLER_ERROR", fmt.Sprintf("operation handler failed: %v", err), http.StatusInternalServerError) // Or map handler error code if available
	}

	// 4. Check if Handler returned *OperationHandlerOutput or final result
	handlerOutput, isOutput := handlerResult.(*OperationHandlerOutput)
	if !isOutput {
		// Handler returned the final result directly (e.g., RAG Query)
		return handlerResult, nil
	}

	// Handler returned details for a Volcengine API call
	if handlerOutput == nil { // Safety check
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "HANDLER_ERROR", "handler returned nil output but indicated an API call was needed", http.StatusInternalServerError)
	}

	// 5. Execute via Internal Volcengine Client (Only if handler returned Output)
	result, err := a.internalVolcengineClient.Execute(
		ctx,
		operationID, // Pass operationID for error context
		handlerOutput.Method,
		handlerOutput.Path,
		handlerOutput.Query,
		handlerOutput.Body, // The internal API request struct from the handler
	)

	// 6. Handle Result/Error from Internal Client
	if err != nil {
		// The internal client should ideally return a domain.AdapterError or an error
		// that can be mapped to one.
		// Example: check if err is already *domain.AdapterError, otherwise wrap it.
		if _, ok := err.(*domain.AdapterError); !ok {
			// Wrap non-AdapterError from the internal client to ensure consistent error type.
			return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "PROVIDER_API_ERROR", fmt.Sprintf("internal client execution failed: %v", err), http.StatusInternalServerError)
		}
		return nil, err
	}

	return result, nil
}

func intPtr(v int) *int {
	return &v
}

func stringPtr(v string) *string {
	return &v
}

func boolPtr(v bool) *bool {
	return &v
}

func float64Ptr(v float64) *float64 {
	return &v
}
