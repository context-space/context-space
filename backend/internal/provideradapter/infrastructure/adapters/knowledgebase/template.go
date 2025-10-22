package knowledgebase

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	openaiclient "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/knowledgebase/openai/client"
	openaitypes "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/knowledgebase/openai/types"
	volcclient "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/knowledgebase/volcengine/client"
	volcenginetypes "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/knowledgebase/volcengine/types"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

const (
	region      = "cn-north-1"
	endpoint    = "https://api-knowledgebase.mlp.cn-beijing.volces.com"
	serviceName = "air"
)

var DefaultKnowledgebaseTemplates = []string{
	"cfa_knowledgebase",
}

// Register the Volcengine Knowledge Base adapter template
func init() {
	for _, identifier := range DefaultKnowledgebaseTemplates {
		template := &KnowledgeBaseTemplate{
			Identifier: identifier,
		}
		registry.RegisterAdapterTemplate(identifier, template)
	}
}

// KnowledgeBaseTemplate is a template for creating Volcengine Knowledge Base adapters
type KnowledgeBaseTemplate struct {
	Identifier string
}

// CreateAdapter creates a new Volcengine Knowledge Base adapter based on the provided configuration.
// It handles extracting configuration, creating the internal Volcengine client, and instantiating the adapter.
func (t *KnowledgeBaseTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
	providerInfo := &domain.ProviderAdapterInfo{
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    provider.AuthType,
	}

	// Define default adapter settings using correct types
	defaultTimeout := 30 * time.Second
	defaultMaxRetries := 3
	defaultRetryBackoff := 1 * time.Second
	defaultFailureThreshold := 5 // Use int for CircuitBreakerConfig fields
	defaultResetTimeout := 60    // Use int
	defaultHalfOpenMaxCalls := 2 // Use int

	// Create AdapterConfig using defaults and correct types
	adapterConfig := &domain.AdapterConfig{
		Timeout:      defaultTimeout,
		MaxRetries:   defaultMaxRetries,
		RetryBackoff: defaultRetryBackoff,
		CircuitBreaker: domain.CircuitBreakerConfig{
			FailureThreshold: defaultFailureThreshold,
			ResetTimeout:     defaultResetTimeout,
			HalfOpenMaxCalls: defaultHalfOpenMaxCalls,
		},
	}

	jsonBytes, err := sonic.Marshal(provider.CustomConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal provider: %w", err)
	}

	var jsonAttributes struct {
		VolcengineCredentials *volcenginetypes.VolcengineCredential `json:"volcengine_credentials"`
		OpenaiCredentials     *openaitypes.OpenaiCredential         `json:"openai_credentials"`
		KnowledgebaseConfig   *KnowledgebaseAdapterConfig           `json:"knowledgebase_config"`
	}
	err = sonic.Unmarshal(jsonBytes, &jsonAttributes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal provider: %w", err)
	}

	volcengineCreds := jsonAttributes.VolcengineCredentials

	baseConfig := DefaultKnowledgebaseAdapterConfig

	mergeKnowledgebaseConfig(&baseConfig, jsonAttributes.KnowledgebaseConfig)

	internalClient := volcclient.NewVolcengineClient(
		providerInfo.Identifier, // Pass providerIdentifier
		endpoint,                // Correct: Pass endpoint as apiBaseURL
		serviceName,             // Correct: Pass serviceName as service
		region,                  // Correct: Pass region as region
		adapterConfig.Timeout,   // Pass timeout from AdapterConfig
		volcengineCreds,         // Pass the credentials
	)

	openaiCreds := jsonAttributes.OpenaiCredentials
	openaiClient := openaiclient.NewOpenaiClient(
		openaiCreds.APIKey,
		openaiCreds.BaseURL,
	)

	adapter, err := NewKnowledgeBaseAdapter(
		providerInfo,
		adapterConfig,
		*internalClient, // Pass the dereferenced struct value
		*openaiClient,
		&baseConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create KnowledgeBaseAdapter: %w", err)
	}

	return adapter, nil
}

// ValidateConfig validates the configuration structure for the Volcengine Knowledge Base template.
func (t *KnowledgeBaseTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider is required")
	}

	if provider.Identifier != t.Identifier {
		return fmt.Errorf("invalid provider identifier, must be '%s'", t.Identifier)
	}

	if provider.AuthType != types.AuthTypeNone {
		return fmt.Errorf("missing or invalid auth_type, expected 'none'")
	}

	jsonBytes, err := sonic.Marshal(provider.CustomConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal provider: %w", err)
	}

	var jsonAttributes struct {
		VolcengineCredentials *volcenginetypes.VolcengineCredential `json:"volcengine_credentials"`
		OpenaiCredentials     *openaitypes.OpenaiCredential         `json:"openai_credentials"`
		KnowledgebaseConfig   *KnowledgebaseAdapterConfig           `json:"knowledgebase_config"`
	}
	err = sonic.Unmarshal(jsonBytes, &jsonAttributes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal provider: %w", err)
	}

	if jsonAttributes.VolcengineCredentials == nil {
		return fmt.Errorf("volcengine_credentials is required")
	}

	if jsonAttributes.VolcengineCredentials.AccessKeyID == "" || jsonAttributes.VolcengineCredentials.SecretAccessKey == "" {
		return fmt.Errorf("volcengine_credentials access_key_id or secret_access_key is required")
	}

	if jsonAttributes.OpenaiCredentials == nil {
		return fmt.Errorf("openai_credentials is required")
	}

	if jsonAttributes.KnowledgebaseConfig == nil {
		return fmt.Errorf("knowledgebase_config is required")
	}

	if jsonAttributes.KnowledgebaseConfig.Project == "" {
		return fmt.Errorf("knowledgebase_config project is required")
	}

	if jsonAttributes.KnowledgebaseConfig.CollectionName == "" {
		return fmt.Errorf("knowledgebase_config collection_name is required")
	}

	if jsonAttributes.KnowledgebaseConfig.Search != nil {
		if jsonAttributes.KnowledgebaseConfig.Search.Limit != nil && *jsonAttributes.KnowledgebaseConfig.Search.Limit <= 0 {
			return fmt.Errorf("knowledgebase_config search limit must be greater than 0")
		}
	}

	if jsonAttributes.KnowledgebaseConfig.Chat != nil {
		if jsonAttributes.KnowledgebaseConfig.Chat.Model != nil && *jsonAttributes.KnowledgebaseConfig.Chat.Model == "" {
			return fmt.Errorf("knowledgebase_config chat model is required")
		}
		if jsonAttributes.KnowledgebaseConfig.Chat.Temperature != nil && *jsonAttributes.KnowledgebaseConfig.Chat.Temperature < 0 {
			return fmt.Errorf("knowledgebase_config chat temperature must be greater than 0")
		}
	}

	if jsonAttributes.KnowledgebaseConfig.Query != nil {
		if jsonAttributes.KnowledgebaseConfig.Query.SearchLimit != nil && *jsonAttributes.KnowledgebaseConfig.Query.SearchLimit <= 0 {
			return fmt.Errorf("knowledgebase_config query search limit must be greater than 0")
		}
		if jsonAttributes.KnowledgebaseConfig.Query.RerankRetrieveCount != nil && *jsonAttributes.KnowledgebaseConfig.Query.RerankRetrieveCount <= 0 {
			return fmt.Errorf("knowledgebase_config query rerank retrieve count must be greater than 0")
		}
		if jsonAttributes.KnowledgebaseConfig.Query.RerankModel != nil && *jsonAttributes.KnowledgebaseConfig.Query.RerankModel == "" {
			return fmt.Errorf("knowledgebase_config query rerank model is required")
		}
		if jsonAttributes.KnowledgebaseConfig.Query.LLMModel != nil && *jsonAttributes.KnowledgebaseConfig.Query.LLMModel == "" {
			return fmt.Errorf("knowledgebase_config query llm model is required")
		}
		if jsonAttributes.KnowledgebaseConfig.Query.LLMTemperature != nil && *jsonAttributes.KnowledgebaseConfig.Query.LLMTemperature < 0 {
			return fmt.Errorf("knowledgebase_config query llm temperature must be greater than 0")
		}
	}

	return nil
}

// mergeKnowledgebaseConfig merge KnowledgebaseAdapterConfig, only non-nil values will override default values
func mergeKnowledgebaseConfig(base *KnowledgebaseAdapterConfig, input *KnowledgebaseAdapterConfig) {
	if input == nil {
		return
	}

	base.Project = input.Project
	base.CollectionName = input.CollectionName

	if input.Search != nil {
		if base.Search == nil {
			base.Search = &SearchConfig{}
		}
		mergeSearchConfig(base.Search, input.Search)
	}

	if input.Chat != nil {
		if base.Chat == nil {
			base.Chat = &ChatConfig{}
		}
		mergeChatConfig(base.Chat, input.Chat)
	}

	if input.Query != nil {
		if base.Query == nil {
			base.Query = &QueryConfig{}
		}
		mergeQueryConfig(base.Query, input.Query)
	}
}

// mergeSearchConfig merge SearchConfig, only non-nil values will override default values
func mergeSearchConfig(base *SearchConfig, input *SearchConfig) {
	if input.Limit != nil && *input.Limit > 0 {
		base.Limit = input.Limit
	}
}

// mergeChatConfig merge ChatConfig, only non-nil values will override default values
func mergeChatConfig(base *ChatConfig, input *ChatConfig) {
	if input.Model != nil {
		base.Model = input.Model
	}
	if input.Stream != nil {
		base.Stream = input.Stream
	}
	if input.Temperature != nil {
		base.Temperature = input.Temperature
	}
}

// mergeQueryConfig merge QueryConfig, only non-nil values will override default values
func mergeQueryConfig(base *QueryConfig, input *QueryConfig) {
	if input.SearchLimit != nil {
		base.SearchLimit = input.SearchLimit
	}
	if input.RerankRetrieveCount != nil {
		base.RerankRetrieveCount = input.RerankRetrieveCount
	}
	if input.RerankModel != nil {
		base.RerankModel = input.RerankModel
	}
	if input.LLMModel != nil {
		base.LLMModel = input.LLMModel
	}
	if input.LLMTemperature != nil {
		base.LLMTemperature = input.LLMTemperature
	}

	if input.RewriteQuery != nil {
		base.RewriteQuery = input.RewriteQuery
	}
	if input.Rerank != nil {
		base.Rerank = input.Rerank
	}
}
