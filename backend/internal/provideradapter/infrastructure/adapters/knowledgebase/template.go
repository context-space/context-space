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
	providercore "github.com/context-space/context-space/backend/internal/providercore/domain"
)

const (
	identifier  = "cfa_knowledgebase"
	region      = "cn-north-1"
	endpoint    = "https://api-knowledgebase.mlp.cn-beijing.volces.com"
	serviceName = "air"
)

// Register the Volcengine Knowledge Base adapter template
func init() {
	template := &KnowledgeBaseTemplate{}
	registry.RegisterAdapterTemplate("cfa_knowledgebase", template)
}

// KnowledgeBaseTemplate is a template for creating Volcengine Knowledge Base adapters
type KnowledgeBaseTemplate struct{}

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
	}
	err = sonic.Unmarshal(jsonBytes, &jsonAttributes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal provider: %w", err)
	}

	volcengineCreds := jsonAttributes.VolcengineCredentials

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
		&opDefaults,
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

	if provider.Identifier != identifier {
		return fmt.Errorf("invalid provider identifier, must be '%s'", identifier)
	}

	if provider.AuthType != providercore.AuthTypeNone {
		return fmt.Errorf("missing or invalid auth_type, expected 'none'")
	}

	jsonBytes, err := sonic.Marshal(provider.CustomConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal provider: %w", err)
	}

	var jsonAttributes struct {
		VolcengineCredentials *volcenginetypes.VolcengineCredential `json:"volcengine_credentials"`
		OpenaiCredentials     *openaitypes.OpenaiCredential         `json:"openai_credentials"`
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

	return nil
}
