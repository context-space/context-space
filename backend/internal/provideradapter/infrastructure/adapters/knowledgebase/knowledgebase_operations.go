package knowledgebase

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	domain "github.com/context-space/context-space/backend/internal/provideradapter/domain"
	volcenginetypes "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/knowledgebase/volcengine/types"
	"github.com/context-space/context-space/backend/internal/shared/utils"
	"github.com/openai/openai-go"
)

// Define constants for operation IDs used by handlers.
const (
	operationIDSearchKnowledge = "search_knowledge"
	operationIDChatCompletions = "chat_completions"
	operationIDQuery           = "query"
	// Add other potential Volcengine KB operation IDs here
)

// Define constants for API paths used by handlers.
const (
	endpointSearchKnowledge = "/api/knowledge/collection/search_knowledge"
	endpointChatCompletions = "/api/knowledge/chat/completions"
	// Add other endpoints as needed
)

// SearchConfig holds defaults specifically for the 'volcengine_kb.search_knowledge' operation.
type SearchConfig struct {
	Limit *int // Default limit for search results
}

// ChatConfig holds defaults specifically for the 'volcengine_kb.chat_completions' operation.
type ChatConfig struct {
	Model       *string  // Default LLM model ID or endpoint ID
	Stream      *bool    // Default stream setting
	Temperature *float64 // Default sampling temperature
}

// QueryConfig holds defaults specifically for the RAG 'volcengine_kb.query' operation.
type QueryConfig struct {
	// Defaults for the internal Search Knowledge step
	SearchLimit         *int    // Limit for initial knowledge retrieval (before rerank)
	RewriteQuery        *bool   // Enable query rewriting based on history?
	Rerank              *bool   // Enable reranking of search results?
	RerankRetrieveCount *int    // How many chunks to retrieve for reranking (must be >= SearchLimit if Rerank is true)
	RerankModel         *string // Specific rerank model name (optional, uses API default if empty)

	// Defaults for the internal LLM call step
	LLMModel       *string  // LLM model ID for the final answer generation
	LLMTemperature *float64 // Temperature for the LLM generation step
}

// OperationHandlerOutput defines what the handler returns to KnowledgeBaseAdapter.Execute.
// It includes the necessary details to make the internal API call.
type OperationHandlerOutput struct {
	Method string      // HTTP Method (e.g., http.MethodPost)
	Path   string      // API Path (e.g., "/api/knowledge/collection/search_knowledge")
	Query  url.Values  // URL Query parameters (if any)
	Body   interface{} // Fully constructed internal API request struct (e.g., *volcenginetypes.SearchKnowledgeRequest)
}

// OperationHandler translates user params (already validated and decoded into a struct)
// and credentials into either the internal API request details OR the final result.
// It returns an interface{} which could be either:
// 1. *OperationHandlerOutput: For operations that require KnowledgeBaseAdapter.Execute to make the final API call.
// 2. The final result directly: For operations (like RAG query) that perform multiple API calls internally.
type OperationHandler func(ctx context.Context, processedParams interface{}) (interface{}, error)

// OperationDefinition combines the user-facing parameter schema and the handler function.
type OperationDefinition struct {
	Schema  interface{}      // Pointer to the user-facing parameter struct (e.g., &SearchKnowledgeParams{})
	Handler OperationHandler // The handler function for this operation
}

// Operations maps operation IDs to their definitions.
type Operations map[string]OperationDefinition

// MessageRole defines the possible roles in a conversation.
type MessageRole string

const (
	MessageRoleSystem    MessageRole = "system"
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

// Message structure for conversation history (common for chat-like operations)
type Message struct {
	Role    MessageRole `mapstructure:"role" json:"role" validate:"required,oneof=system user assistant"`
	Content string      `mapstructure:"content" json:"content" validate:"required"` // Assuming string content for now
	// Note: Volcengine API might support richer content types (ToolCalls etc.) - adapt if needed
}

// Parameters for the search_knowledge operation.
// The collection name is now implicitly defined by the adapter instance's defaults.
type SearchKnowledgeParams struct {
	Query string `mapstructure:"query" json:"query" validate:"required"`
}

type ChatCompletionsParams struct {
	Model    string    `mapstructure:"model" json:"model" validate:"required"` // Model ID or Endpoint ID
	Messages []Message `mapstructure:"messages" json:"messages" validate:"required,min=1,dive"`
}

type QueryParams struct {
	Query    string    `mapstructure:"query" json:"query" validate:"required"` // User query
	Messages []Message `mapstructure:"messages" json:"messages" validate:"omitempty,min=1,dive"`
}

// RegisterOperation registers a single operation, linking its ID, user parameter schema,
// and the handler function that translates user params to internal API details.
// It stores the schema in the BaseAdapter for parameter processing and the full
// definition (including the handler) in the KnowledgeBaseAdapter's operations map.
func (a *KnowledgeBaseAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler) {
	a.BaseAdapter.RegisterOperation(operationID, schema)

	if a.operations == nil {
		// This check might be redundant if constructor always initializes, but safer
		a.operations = make(Operations)
	}
	a.operations[operationID] = OperationDefinition{
		Schema:  schema, // Storing schema here aligns with scratchpad plan
		Handler: handler,
	}
}

// registerOperations registers all supported operations for the KnowledgeBaseAdapter.
// It calls RegisterOperation for each specific operation, defining the handler logic inline.
func (a *KnowledgeBaseAdapter) registerOperations() {
	// Register search_knowledge operation
	a.RegisterOperation(
		operationIDSearchKnowledge,
		&SearchKnowledgeParams{},
		a.handleSearchKnowledge,
	)

	// Register query operation
	// This handler will now perform the full RAG logic: search -> augment -> generate
	a.RegisterOperation(
		operationIDQuery,
		&QueryParams{},
		a.handleQuery,
	)
	// Add registrations for other operations here
}

func (a *KnowledgeBaseAdapter) handleSearchKnowledge(ctx context.Context, processedParams interface{}) (interface{}, error) {
	params, ok := processedParams.(*SearchKnowledgeParams)
	if !ok {
		return nil, fmt.Errorf("invalid parameters type for %s: expected *SearchKnowledgeParams, got %T", operationIDSearchKnowledge, processedParams)
	}
	// Credential is not used directly in this handler anymore
	// The internal client will use stored credentials

	// Construct Internal API Request Body using actual types
	apiRequestBody := &volcenginetypes.SearchKnowledgeRequest{
		Name:    a.baseConfig.CollectionName, // Use CollectionName from adapter defaults
		Project: a.baseConfig.Project,        // Use Project from adapter defaults
		Query:   params.Query,
		Limit:   *a.baseConfig.Search.Limit, // Use adapter's Search default limit
		// PreProcessing/PostProcessing: Rely on Volcengine API defaults for now, or expose via SearchKnowledgeParams
	}

	return &OperationHandlerOutput{
		Method: http.MethodPost,
		Path:   endpointSearchKnowledge, // Use constant
		Query:  nil,
		Body:   apiRequestBody,
	}, nil
}

func (a *KnowledgeBaseAdapter) handleChatCompletions(ctx context.Context, processedParams interface{}) (interface{}, error) {
	params, ok := processedParams.(*ChatCompletionsParams)
	if !ok {
		return nil, fmt.Errorf("invalid parameters type for %s: expected *ChatCompletionsParams, got %T", operationIDChatCompletions, processedParams)
	}
	// Credential is not used directly in this handler anymore
	// The internal client will use stored credentials

	// Convert user messages to internal API messages
	apiMessages := make([]volcenginetypes.Message, len(params.Messages))
	for i, msg := range params.Messages {
		// Explicitly cast MessageRole back to string for the volcengine type
		apiMessages[i] = volcenginetypes.Message{Role: string(msg.Role), Content: msg.Content}
	}

	// Prepare optional parameters using Chat defaults
	var temp *float64
	// Use Chat.Temperature default
	if a.baseConfig.Chat.Temperature != nil && *a.baseConfig.Chat.Temperature >= 0 { // Assuming 0 is a valid temp, check range if needed
		t := *a.baseConfig.Chat.Temperature
		temp = &t
	}
	// Use Chat.Stream default
	stream := *a.baseConfig.Chat.Stream
	// MaxTokens is not currently configurable via defaults, rely on API default

	// Construct Internal API Request Body using actual types
	apiRequestBody := &volcenginetypes.ChatCompletionsRequest{
		Model:       params.Model, // Model comes from user params for this operation
		Messages:    apiMessages,
		Temperature: temp, // Use pointer prepared above
		Stream:      stream,
		// Add fields like ModelVersion, APIKey, ReturnTokenUsage if they should be configurable via ChatCompletionsParams
	}

	return &OperationHandlerOutput{
		Method: http.MethodPost,
		Path:   endpointChatCompletions,
		Query:  nil,
		Body:   apiRequestBody,
	}, nil
}

func (a *KnowledgeBaseAdapter) handleQuery(ctx context.Context, processedParams interface{}) (interface{}, error) {
	params, ok := processedParams.(*QueryParams)
	if !ok {
		return nil, fmt.Errorf("invalid parameters type for %s: expected *QueryParams, got %T", operationIDQuery, processedParams)
	}

	// ** RAG Steps **
	// 1. Search Knowledge
	// 2. Process Search Results
	// 3. Prepare OpenAI Prompt
	// 4. Call OpenAI
	// 5. Process OpenAI Response and return updated message list

	// Step 1: Search Knowledge
	searchRequestBody := &volcenginetypes.SearchKnowledgeRequest{
		Name:    a.baseConfig.CollectionName,
		Project: a.baseConfig.Project,
		Query:   params.Query,
		Limit:   *a.baseConfig.Query.SearchLimit, // Use Query default SearchLimit for RAG context fetching
	}

	// Configure PreProcessing based on Query defaults and messages
	if len(params.Messages) > 0 || *a.baseConfig.Query.RewriteQuery {
		searchRequestBody.PreProcessing = &volcenginetypes.PreProcessing{}
		// Add messages if they exist
		if len(params.Messages) > 0 {
			// Determine capacity and if prepending system message is needed
			prependSystemMessage := params.Messages[0].Role != MessageRoleSystem
			capacity := len(params.Messages)
			if prependSystemMessage {
				capacity++
			}

			// Create the slice for Volcengine API messages
			volcApiMessages := make([]volcenginetypes.Message, 0, capacity)

			// Prepend empty system message ONLY if the original list doesn't start with one
			if prependSystemMessage {
				volcApiMessages = append(volcApiMessages, volcenginetypes.Message{Role: "system", Content: ""})
			}

			// Append converted user/assistant messages (and original system message if present)
			for _, msg := range params.Messages {
				volcApiMessages = append(volcApiMessages, volcenginetypes.Message{
					Role:    string(msg.Role),
					Content: msg.Content,
				})
			}
			searchRequestBody.PreProcessing.Messages = volcApiMessages

			// Enable rewrite only if RewriteQuery default is true AND original user messages exist
			if *a.baseConfig.Query.RewriteQuery {
				rewrite := true
				searchRequestBody.PreProcessing.Rewrite = &rewrite
			}
		}
	}

	// Configure PostProcessing based on Query defaults
	if *a.baseConfig.Query.Rerank {
		searchRequestBody.PostProcessing = &volcenginetypes.PostProcessing{}
		rerankSwitch := true
		searchRequestBody.PostProcessing.RerankSwitch = &rerankSwitch
		searchRequestBody.PostProcessing.RetrieveCount = a.baseConfig.Query.RerankRetrieveCount
		// Set rerank model only if specified in defaults
		if a.baseConfig.Query.RerankModel != nil {
			searchRequestBody.PostProcessing.RerankModel = a.baseConfig.Query.RerankModel
		}
		// Other PostProcessing fields (ChunkDiffusionCount, ChunkGroup, RerankOnlyChunk, GetAttachmentLink)
		// could be added to QueryDefaults and configured here if needed.
	}

	searchOperationID := operationIDSearchKnowledge // For context in error reporting

	searchResultRaw, err := a.internalVolcengineClient.Execute(
		ctx,
		searchOperationID, // Use search op ID for context
		http.MethodPost,
		endpointSearchKnowledge, // Use constant
		nil,                     // No query params for search
		searchRequestBody,
	)
	if err != nil {
		// Error already wrapped by internal client or needs wrapping
		if _, ok := err.(*domain.AdapterError); !ok {
			return nil, domain.NewAdapterError(a.ProviderAdapterInfo.Identifier, searchOperationID, domain.ErrProviderAPIError, fmt.Sprintf("internal knowledge search failed: %v", err), http.StatusInternalServerError)
		}
		return nil, err // Return existing AdapterError
	}

	// Step 2: Process Search Results
	// Extract relevant text chunks from searchResultRaw.
	knowledgeChunks := ""
	searchResultMap, ok := searchResultRaw.(map[string]interface{})
	if !ok {
		// Log or handle unexpected result type. Proceed without chunks if parsing fails.
		// Consider logging: a.Logger.Warnf("Unexpected search result type: %T", searchResultRaw)
	} else {
		if candidates, ok := searchResultMap["candidates"].([]interface{}); ok {
			for i, candRaw := range candidates {
				if candidate, ok := candRaw.(map[string]interface{}); ok {
					if chunk, ok := candidate["chunk"].(string); ok {
						// TODO: Consider adding metadata or refining chunk formatting
						knowledgeChunks = utils.StringsBuilder(knowledgeChunks, fmt.Sprintf("--- Document %d ---\n%s\n\n", i+1, chunk))
					}
				}
			}
		}
	}

	// Step 3: Prepare OpenAI Prompt
	systemPrompt := SystemPromptRAG // Use constant from prompts.go

	// Estimate capacity: System + Context + History + User Query
	initialCapacity := len(params.Messages) + 3
	apiMessages := make([]openai.ChatCompletionMessageParamUnion, 0, initialCapacity)

	// Add System Prompt
	apiMessages = append(apiMessages, openai.SystemMessage(systemPrompt))

	// Convert history messages (params.Messages)
	// params.Messages is []knowledgebase.Message
	for _, msg := range params.Messages {
		switch msg.Role {
		case MessageRoleUser:
			apiMessages = append(apiMessages, openai.UserMessage(msg.Content))
		case MessageRoleAssistant:
			apiMessages = append(apiMessages, openai.AssistantMessage(msg.Content))
		case MessageRoleSystem:
			continue
		}
	}

	// Add Context Documents if available
	if knowledgeChunks != "" {
		contextSection := fmt.Sprintf("--- Context Documents ---\n%s", knowledgeChunks)
		// Add Context Documents as a user message. Alternative: System message, but user often works well.
		apiMessages = append(apiMessages, openai.UserMessage(contextSection))
	}

	// Append the current user query
	apiMessages = append(apiMessages, openai.UserMessage(params.Query))

	// Step 4: Call OpenAI
	// Check for LLM Model configuration removed from here
	openaiRequest := openai.ChatCompletionNewParams{
		Model:    *a.baseConfig.Query.LLMModel, // Use Query default LLM model
		Messages: apiMessages,
	}
	// Apply Query default LLM temperature if set
	if a.baseConfig.Query.LLMTemperature != nil && *a.baseConfig.Query.LLMTemperature >= 0 { // Assuming 0 is valid
		openaiRequest.Temperature = openai.Float(*a.baseConfig.Query.LLMTemperature)
	}
	// Apply default stream setting if needed (currently not used for RAG response)
	// if a.defaults.Stream { // Which stream default? Chat? Query?

	openaiResponse, err := a.openaiClient.SdkClient.Chat.Completions.New(ctx, openaiRequest)
	if err != nil {
		// Map OpenAI errors (e.g., rate limits, auth) to AdapterError codes if possible.
		// TODO: Implement more specific error mapping based on OpenAI error types/codes.
		return nil, domain.NewAdapterError(a.ProviderAdapterInfo.Identifier, operationIDQuery, domain.ErrLLMProviderError, fmt.Sprintf("OpenAI API call failed: %v", err), http.StatusInternalServerError)
	}

	// Step 5: Process OpenAI Response and Return Updated Messages
	if len(openaiResponse.Choices) == 0 || openaiResponse.Choices[0].Message.Content == "" {
		// Handle cases where OpenAI returns no response or empty content
		return nil, domain.NewAdapterError(a.ProviderAdapterInfo.Identifier, operationIDQuery, domain.ErrLLMEmptyResponse, "OpenAI returned no usable response", http.StatusInternalServerError)
	}

	// Extract the assistant's response content
	assistantContent := openaiResponse.Choices[0].Message.Content

	// Create the assistant message using the knowledgebase.Message type
	assistantMessage := Message{
		Role:    MessageRoleAssistant,
		Content: assistantContent,
	}

	// Append the new assistant message to the original messages list
	// Note: params.Messages might be nil if it's the first turn.
	updatedMessages := append(params.Messages, assistantMessage)

	// Return the updated list of messages (including the new assistant response)
	return updatedMessages, nil
}
