package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/invopop/jsonschema"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	DefaultOpenAIBaseURL      = "https://api.openai.com/v1"
	DefaultTranslationTimeout = 60 * time.Second
)

// TranslationConfig holds OpenAI configuration for translation
type TranslationConfig struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

// LanguagePrompts contains translation prompts for different languages
type LanguagePrompts struct {
	ZhCN string // Chinese Simplified prompt
	ZhTW string // Chinese Traditional prompt
}

// TranslationResponse represents the structured output from OpenAI for translations
type TranslationResponse struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Categories  []string                `json:"categories"`
	Permissions []TranslationPermission `json:"permissions"`
	Operations  []TranslationOperation  `json:"operations"`
}

// GenerateSchema generates JSON schema for structured output
func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

// Generate the JSON schema at initialization time
var TranslationResponseSchema = GenerateSchema[TranslationResponse]()

// LoadTranslationConfig loads OpenAI configuration from environment variables or .env file
func LoadTranslationConfig() (*TranslationConfig, error) {
	// Try to load .env file, ignore error if file doesn't exist
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = DefaultOpenAIBaseURL // Default OpenAI URL
	}

	timeoutStr := os.Getenv("OPENAI_TIMEOUT")
	timeout, err := strconv.ParseInt(timeoutStr, 10, 64)
	if err != nil {
		timeout = int64(DefaultTranslationTimeout / time.Second) // Default to seconds if not set
	}

	return &TranslationConfig{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Timeout: time.Duration(timeout) * time.Second,
	}, nil
}

// GetDefaultPrompts returns default translation prompts for each language
func GetDefaultPrompts() LanguagePrompts {
	return LanguagePrompts{
		ZhCN: TranslatorPromptZhCN,
		ZhTW: TranslatorPromptZhTW,
	}
}

// AITranslator handles AI-powered translation
type AITranslator struct {
	client  openai.Client
	config  *TranslationConfig
	prompts LanguagePrompts
}

// NewAITranslator creates a new AI translator instance
func NewAITranslator() (*AITranslator, error) {
	config, err := LoadTranslationConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load translation config: %w", err)
	}

	// Create OpenAI client with custom base URL if provided
	clientOptions := []option.RequestOption{
		option.WithAPIKey(config.APIKey),
	}

	if config.BaseURL != DefaultOpenAIBaseURL {
		clientOptions = append(clientOptions, option.WithBaseURL(config.BaseURL))
	}

	client := openai.NewClient(clientOptions...)

	return &AITranslator{
		client:  client,
		config:  config,
		prompts: GetDefaultPrompts(),
	}, nil
}

// TranslateI18nData translates i18n data to specified language using AI
func (t *AITranslator) TranslateI18nData(ctx context.Context, sourceData TranslationResponse, targetLang string) (TranslationResponse, error) {
	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, t.config.Timeout)
	defer cancel()

	// Get the appropriate prompt for target language
	var prompt string
	switch targetLang {
	case "zh-CN":
		prompt = t.prompts.ZhCN
	case "zh-TW":
		prompt = t.prompts.ZhTW
	default:
		return TranslationResponse{}, fmt.Errorf("unsupported target language: %s", targetLang)
	}

	fmt.Printf("Preparing translation request for %s...\n", targetLang)

	// Convert source data to JSON string
	sourceJSON, err := json.MarshalIndent(sourceData, "", "  ")
	if err != nil {
		return TranslationResponse{}, fmt.Errorf("failed to marshal source data: %w", err)
	}

	// Prepare the full prompt with source JSON
	fullPrompt := fmt.Sprintf("%s\n\nSource JSON to translate:\n```json\n%s\n```", prompt, string(sourceJSON))

	// Create schema parameter for structured output
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "translation_response",
		Description: openai.String("Translation of i18n data to target language"),
		Schema:      TranslationResponseSchema,
		Strict:      openai.Bool(true),
	}

	fmt.Printf("Calling OpenAI API for %s translation (timeout: %v)...\n", targetLang, t.config.Timeout)

	// Call OpenAI API with structured output
	response, err := t.client.Chat.Completions.New(timeoutCtx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini, // Use a model that supports structured output
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(fullPrompt),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
		Temperature: openai.Float(0.1), // Low temperature for consistent translations
	})

	// Check for timeout or other context errors
	if err != nil {
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return TranslationResponse{}, fmt.Errorf("translation request timed out after %v", t.config.Timeout)
		}
		return TranslationResponse{}, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(response.Choices) == 0 {
		return TranslationResponse{}, fmt.Errorf("no response from OpenAI API")
	}

	fmt.Printf("Received response from OpenAI, parsing structured output for %s...\n", targetLang)

	// Extract the structured response
	content := response.Choices[0].Message.Content
	if content == "" {
		return TranslationResponse{}, fmt.Errorf("empty response from OpenAI API")
	}

	// Parse the structured response directly
	var translationResponse TranslationResponse
	if err := json.Unmarshal([]byte(content), &translationResponse); err != nil {
		return TranslationResponse{}, fmt.Errorf("failed to parse structured response: %w\nResponse: %s", err, content)
	}

	return translationResponse, nil
}

// IsTranslationEnabled checks if translation is enabled (API key is available)
func IsTranslationEnabled() bool {
	_ = godotenv.Load() // Try to load .env file
	return os.Getenv("OPENAI_API_KEY") != ""
}
