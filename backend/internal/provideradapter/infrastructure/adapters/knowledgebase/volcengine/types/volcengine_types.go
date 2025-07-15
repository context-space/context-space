package types

import "encoding/json"

// VolcengineCredential defines the expected credential structure for Volcengine.
// Moved from knowledgebase package to break import cycle.
type VolcengineCredential struct {
	AccessKeyID     string `json:"access_key_id" validate:"required"`
	SecretAccessKey string `json:"secret_access_key" validate:"required"`
}

// ApiResponse represents the standard wrapper for Volcengine API responses.
type ApiResponse struct {
	Code      int             `json:"code"`       // Business status code (0 for success)
	Message   string          `json:"message"`    // Response message
	RequestID string          `json:"request_id"` // Unique request ID
	Data      json.RawMessage `json:"data"`       // The actual data payload (type depends on the specific API call)
}

// Message represents a single message in a conversation.
// Used in both ChatCompletionsRequest and SearchKnowledgeRequest.PreProcessing.
type Message struct {
	Role    string `json:"role"`    // e.g., "system", "user", "assistant"
	Content string `json:"content"` // Message content (can be complex for multimodal, but basic string for now)
	// Note: Content can be an array for multimodal input, need adjustment if supporting that.
}

// TokenUsage provides information about token consumption.
// Used in ChatCompletionsResponseData (if structure is confirmed)
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`     // Tokens in the input prompt
	CompletionTokens int `json:"completion_tokens"` // Tokens in the generated completion
	TotalTokens      int `json:"total_tokens"`      // Total tokens used
}

// EmbeddingTokenUsage provides token usage specific to embedding operations.
// Used in SearchKnowledgeResponseData.TokenUsage for embedding_token_usage.
type EmbeddingTokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`     // Tokens from the query for embedding
	CompletionTokens int `json:"completion_tokens"` // Usually 0 for embedding
	TotalTokens      int `json:"total_tokens"`      // Total tokens for embedding
}

// SearchKnowledgeTokenUsage holds token usage details for the search knowledge API.
// Used in SearchKnowledgeResponseData.
type SearchKnowledgeTokenUsage struct {
	EmbeddingTokenUsage EmbeddingTokenUsage `json:"embedding_token_usage,omitempty"` // Token usage during embedding phase
	RerankTokenUsage    *int                `json:"rerank_token_usage,omitempty"`    // Token usage during rerank phase (pointer for optionality)
	RewriteTokenUsage   *int                `json:"rewrite_token_usage,omitempty"`   // Token usage during query rewrite phase (pointer for optionality)
}

// SearchKnowledgeRequest represents the request body for the search_knowledge API.
type SearchKnowledgeRequest struct {
	Project        string          `json:"project,omitempty"`         // Knowledge base project name
	Name           string          `json:"name,omitempty"`            // Knowledge base name
	ResourceID     string          `json:"resource_id,omitempty"`     // Knowledge base resource ID (use either name+project or resource_id)
	Query          string          `json:"query"`                     // Search query (Required)
	Limit          int             `json:"limit,omitempty"`           // Max number of results (Default: 10, Range: [1, 200])
	QueryParam     *QueryParam     `json:"query_param,omitempty"`     // Filtering and advanced query settings
	DenseWeight    *float64        `json:"dense_weight,omitempty"`    // Weight for dense vector in hybrid search (Range: [0.2, 1], Default: 0.5)
	PreProcessing  *PreProcessing  `json:"pre_processing,omitempty"`  // Query pre-processing options
	PostProcessing *PostProcessing `json:"post_processing,omitempty"` // Result post-processing options
}

// QueryParam contains filtering parameters for the search query.
type QueryParam struct {
	DocFilter map[string]interface{} `json:"doc_filter,omitempty"` // Document filter map (refer to filter expression docs)
	// Add other potential query params if discovered
}

// PreProcessing defines options for query pre-processing.
type PreProcessing struct {
	NeedInstruction  *bool     `json:"need_instruction,omitempty"`   // Add instruction before query (Default: false)
	Rewrite          *bool     `json:"rewrite,omitempty"`            // Rewrite query based on history (Default: false, requires Messages)
	ReturnTokenUsage *bool     `json:"return_token_usage,omitempty"` // Return token usage details (Default: false)
	Messages         []Message `json:"messages,omitempty"`           // Conversation history (Required if Rewrite is true)
}

// PostProcessing defines options for result post-processing.
type PostProcessing struct {
	RerankSwitch        *bool   `json:"rerank_switch,omitempty"`         // Enable reranking (Default: false)
	RetrieveCount       *int    `json:"retrieve_count,omitempty"`        // Number of chunks for reranking (Default: 25, must be >= Limit)
	ChunkDiffusionCount *int    `json:"chunk_diffusion_count,omitempty"` // Context chunk diffusion count (Default: 0, Range: [0, 5])
	ChunkGroup          *bool   `json:"chunk_group,omitempty"`           // Group chunks by original document order (Default: false)
	RerankModel         *string `json:"rerank_model,omitempty"`          // Rerank model name (Default: "m3-v2-rerank")
	RerankOnlyChunk     *bool   `json:"rerank_only_chunk,omitempty"`     // Rerank using only chunk content (Default: false)
	GetAttachmentLink   *bool   `json:"get_attachment_link,omitempty"`   // Get temporary links for attachments (Default: false)
}

// SearchKnowledgeResponseData represents the structure within the 'data' field of a successful search_knowledge response.
// This needs significant expansion based on the docs.
type SearchKnowledgeResponseData struct {
	CollectionName string                    `json:"collection_name"`         // Name of the searched knowledge base
	Count          int                       `json:"count"`                   // Number of results returned (<= limit)
	RewriteQuery   *string                   `json:"rewrite_query,omitempty"` // The rewritten query, if rewrite was enabled
	TokenUsage     SearchKnowledgeTokenUsage `json:"token_usage,omitempty"`   // Token usage details, if requested
	ResultList     []ResultItem              `json:"result_list"`             // List of search result chunks/items
	// The fields 'total', 'limit', 'offset' from the previous Go struct are NOT in the markdown 'data' field.
	// The previous 'search_result' field is now named 'result_list'.
}

// ResultItem represents a single item (chunk) in the search_knowledge result list.
type ResultItem struct {
	ID                 string           `json:"id"`                            // Index primary key
	Content            string           `json:"content"`                       // Chunk content / Answer for FAQ / K:V for structured
	Score              float64          `json:"score"`                         // Original retrieval score (before rerank)
	PointID            string           `json:"point_id"`                      // Unique ID of the chunk
	ChunkTitle         string           `json:"chunk_title"`                   // Chunk title
	ChunkID            int              `json:"chunk_id"`                      // Chunk sequence ID within the document
	ProcessTime        int              `json:"process_time"`                  // Search time (ms?)
	RerankScore        *float64         `json:"rerank_score,omitempty"`        // Rerank score, if rerank enabled
	DocInfo            DocInfo          `json:"doc_info"`                      // Information about the original document
	RecallPosition     int              `json:"recall_position"`               // Original rank before rerank
	RerankPosition     *int             `json:"rerank_position,omitempty"`     // Rank after rerank, if rerank enabled
	TableChunkFields   []FieldValuePair `json:"table_chunk_fields,omitempty"`  // Full row data for structured results
	OriginalQuestion   *string          `json:"original_question,omitempty"`   // Original question for FAQ results
	ChunkType          string           `json:"chunk_type"`                    // Type of chunk (e.g., "text", "image", "table")
	ChunkAttachment    []AttachmentInfo `json:"chunk_attachment,omitempty"`    // Attachment info, if requested and present
	OriginalCoordinate *CoordinateInfo  `json:"original_coordinate,omitempty"` // Coordinate in original PDF/PPT
	// Previous structure only had DocInfo and Score. Added many fields from markdown.
}

// DocInfo contains details about the original document a chunk belongs to.
// Adjusted based on markdown 'doc_info' structure inside 'result_list'.
type DocInfo struct {
	DocID      string `json:"doc_id"`      // Document ID
	DocName    string `json:"doc_name"`    // Document name
	CreateTime int64  `json:"create_time"` // Document creation timestamp
	DocType    string `json:"doc_type"`    // Document type (e.g., "pdf")
	DocMeta    string `json:"doc_meta"`    // Document metadata (JSON string)
	Source     string `json:"source"`      // Knowledge source (e.g., "tos_fe")
	Title      string `json:"title"`       // Document title
	// Previous structure had DocContent, DocURL, DocMetaInfo (map). Replaced based on markdown.
}

// FieldValuePair represents a key-value pair for structured data results.
type FieldValuePair struct {
	FieldName  string      `json:"field_name"`  // Field name
	FieldValue interface{} `json:"field_value"` // Field value (type varies)
}

// AttachmentInfo holds information about an attachment (e.g., image) within a chunk.
type AttachmentInfo struct {
	UUID    string `json:"uuid"`    // Unique attachment identifier
	Caption string `json:"caption"` // Image caption (or "\n")
	Type    string `json:"type"`    // Attachment type (e.g., "image")
	Link    string `json:"link"`    // Temporary download link (10 min expiry)
}

// CoordinateInfo represents the location of the knowledge point in the original document.
type CoordinateInfo struct {
	PageNo []int       `json:"page_no"` // [start_page, end_page]
	BBox   [][]float64 `json:"bbox"`    // List of bounding boxes [[x1,y1,x2,y2], ...]
}

// ChatCompletionsRequest represents the request body for the generic chat_completions API.
// Updated based on markdown, removing knowledge base specific fields.
type ChatCompletionsRequest struct {
	Model            string    `json:"model"`                        // LLM model ID or endpoint ID (Required)
	ModelVersion     *string   `json:"model_version,omitempty"`      // Specific model version (Optional)
	APIKey           *string   `json:"api_key,omitempty"`            // API Key for private endpoints (Required if model is endpoint ID)
	Messages         []Message `json:"messages"`                     // Conversation history (Required, min 1)
	Stream           bool      `json:"stream,omitempty"`             // Whether to stream the response (Default: false)
	MaxTokens        *int      `json:"max_tokens,omitempty"`         // Max tokens in the response (Default: 4096, pointer for zero value)
	ReturnTokenUsage *bool     `json:"return_token_usage,omitempty"` // Whether to return token usage (Default: false)
	Temperature      *float64  `json:"temperature,omitempty"`        // Sampling temperature (Range: [0, 1], Default: 0.1, pointer for zero value)
	// Add other chat parameters like top_p, etc. if needed, based on full Doubao API docs.
	// Removed: Project, Name, ResourceID, DocFilter, Limit as they are not in the generic chat API description.
}

// ChatCompletionsResponseData represents the structure within the 'data' field of a successful
// non-streamed chat_completions response, based on the provided markdown.
type ChatCompletionsResponseData struct {
	GeneratedAnswer  string  `json:"generated_answer"`            // The final answer generated by the LLM.
	Usage            *string `json:"usage,omitempty"`             // Token usage statistics string (structure unclear from docs, marked optional)
	ReasoningContent *string `json:"reasoning_content,omitempty"` // Model's reasoning steps (optional, model-dependent)
	// This structure is significantly different from the previous OpenAI-like structure.
	// Removed: ID, Object, Created, Model, Choices, TokenUsage (struct), SearchResult
}

// Note: Structures for streamed Chat Completion responses are not defined here,
// as the markdown only shows conceptual examples without exact chunk structure.
