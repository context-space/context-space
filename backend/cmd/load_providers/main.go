package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	adapter_domain "github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/providercore/domain"
	"github.com/context-space/context-space/backend/internal/shared/types"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// Global OpenAI client
var openaiClient *openai.Client

// ProviderJSON represents the structure of a provider JSON file
type ProviderJSON struct {
	Identifier            string                      `json:"identifier"`
	Name                  string                      `json:"name"`
	Description           string                      `json:"description"`
	AuthType              string                      `json:"auth_type"`
	Status                string                      `json:"status"`
	IconURL               string                      `json:"icon_url"`
	Categories            []string                    `json:"categories"`
	Permissions           []PermissionJSON            `json:"permissions"`
	Operations            []OperationJSON             `json:"operations"`
	OAuthConfig           *adapter_domain.OAuthConfig `json:"oauth_config,omitempty"`
	ApiKeyConfig          *ApiKeyConfig               `json:"api_key_config,omitempty"`
	VolcengineCredentials *VolcengineCredentials      `json:"volcengine_credentials,omitempty"`
	OpenaiCredentials     *OpenaiCredentials          `json:"openai_credentials,omitempty"`
	KnowledgebaseConfig   *KnowledgebaseConfig        `json:"knowledgebase_config,omitempty"`
}

// PermissionJSON represents the structure of a permission in the JSON file
type PermissionJSON struct {
	Identifier  string   `json:"identifier"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	OAuthScopes []string `json:"oauth_scopes,omitempty"`
}

// OperationJSON represents the structure of an operation in the JSON file
type OperationJSON struct {
	Identifier          string          `json:"identifier"`
	Name                string          `json:"name"`
	Description         string          `json:"description"`
	Category            string          `json:"category"`
	RequiredPermissions []string        `json:"required_permissions,omitempty"`
	Parameters          []ParameterJSON `json:"parameters,omitempty"`
}

// ParameterJSON represents the structure of a parameter in the JSON file
type ParameterJSON struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Enum        []string    `json:"enum,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

type Provider struct {
	ID          string                 `json:"id"`
	Identifier  string                 `json:"identifier"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	AuthType    types.ProviderAuthType `json:"auth_type"`
	Status      types.ProviderStatus   `json:"status"`
	IconURL     string                 `json:"icon_url"`
	Categories  []string               `json:"categories"`
	Tags        []string               `json:"tags"`
	Operations  []OperationJSON        `json:"operations"`
	Embedding   []float64              `json:"-"` // Vector embedding for semantic search
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   *time.Time             `json:"deleted_at"`
}
type ApiKeyConfig struct {
	Value string `json:"value"`
}

type VolcengineCredentials struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

type OpenaiCredentials struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

type KnowledgebaseConfig struct {
	Project        string        `json:"project"`
	CollectionName string        `json:"collection_name"`
	Search         *SearchConfig `json:"search,omitempty"`
	Chat           *ChatConfig   `json:"chat,omitempty"`
	Query          *QueryConfig  `json:"query,omitempty"`
}

type SearchConfig struct {
	Limit int `json:"limit"`
}

type ChatConfig struct {
	Model       string  `json:"model"`
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
}

type QueryConfig struct {
	SearchLimit         int     `json:"search_limit"`
	RewriteQuery        bool    `json:"rewrite_query"`
	Rerank              bool    `json:"rerank"`
	RerankRetrieveCount int     `json:"rerank_retrieve_count"`
	RerankModel         string  `json:"rerank_model"`
	LLMModel            string  `json:"llm_model"`
	LLMTemperature      float64 `json:"llm_temperature"`
}

// TranslationData represents the structure of translation data
type TranslationData struct {
	ID                 string `json:"-"`
	ProviderIdentifier string `json:"-"`
	LanguageCode       string `json:"-"`
	Translations       string `json:"-"`
}

// ProviderWithTranslations represents a provider with its translations
type ProviderWithTranslations struct {
	Provider     *domain.Provider
	Translations []TranslationData
	Adapter      *adapter_domain.ProviderAdapterConfig
}

// shouldSkipFile determines if a file/directory should be skipped during provider loading
func shouldSkipFile(fileName string) bool {
	// Skip system files and directories
	systemFiles := []string{
		".DS_Store",  // macOS system file
		".git",       // Git directory
		"Thumbs.db",  // Windows thumbnail cache
		".gitignore", // Git ignore file
		"README.md",  // Documentation files
		"readme.md",
		"README.txt",
		"readme.txt",
		".idea", // IDE directories
		".vscode",
		".env", // Environment files
		".env.local",
		".env.example",
	}

	// Check exact matches
	for _, sysFile := range systemFiles {
		if fileName == sysFile {
			return true
		}
	}

	// Skip files/directories starting with dot (hidden files)
	if strings.HasPrefix(fileName, ".") {
		return true
	}

	return false
}

func main() {
	// Define command-line flags
	providersPath := flag.String("path", "configs/providers", "Path to the directory containing provider JSON files")
	sqlOutputDir := flag.String("sql-output", "generated_sql", "Directory to save generated SQL files")
	update := flag.Bool("update", false, "Generate update provider SQL files")
	providerId := flag.String("provider-id", "", "Provider ID to update")
	loadAll := flag.Bool("all", false, "Load all providers from the providers directory")
	providerNames := flag.String("providers", "", "Specific providers to load (only used when --all is not set)")
	openaiAPIKey := flag.String("openai-key", os.Getenv("OPENAI_API_KEY"), "OpenAI API key for generating embeddings")
	openaiBaseURL := flag.String("openai-base-url", os.Getenv("OPENAI_BASE_URL"), "OpenAI API base URL (optional, defaults to OpenAI's API)")
	generateEmbeddings := flag.Bool("embeddings", true, "Generate embeddings for providers and operations")
	flag.Parse()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize OpenAI client if embeddings are requested
	if *generateEmbeddings {
		if *openaiAPIKey == "" {
			fmt.Printf("Error: OpenAI API key is required when --embeddings flag is used\n")
			fmt.Printf("Set OPENAI_API_KEY environment variable or use --openai-key flag\n")
			return
		}

		config := openai.DefaultConfig(*openaiAPIKey)
		if *openaiBaseURL != "" {
			config.BaseURL = *openaiBaseURL
			fmt.Printf("Using custom OpenAI base URL: %s\n", *openaiBaseURL)
		}

		openaiClient = openai.NewClientWithConfig(config)
		fmt.Println("OpenAI client initialized for embedding generation")
	}

	var providersWithTranslations []*ProviderWithTranslations

	if *loadAll {
		// Load all providers from the providers directory
		fmt.Println("Loading all providers...")
		providersWithTranslations, err = loadProvidersFromJSON(*providersPath)
		if err != nil {
			fmt.Printf("Failed to load providers from JSON: %v\n", err)
			return
		}
	} else {
		// Load only the specified provider
		if *providerNames == "" {
			fmt.Printf("Error: When --all is not specified, you must provide --provider parameter\n")
			fmt.Printf("Usage examples:\n")
			fmt.Printf("  %s --all                         # Generate SQL for all providers\n", os.Args[0])
			fmt.Printf("  %s --provider notion             # Generate SQL for notion provider only\n", os.Args[0])
			return
		}

		fmt.Printf("Loading assigned provider: %s\n", *providerNames)
		providerNames := strings.Split(*providerNames, ",")
		for _, providerName := range providerNames {
			providerDir := filepath.Join(*providersPath, providerName)

			// Check if provider directory exists
			if _, err := os.Stat(providerDir); os.IsNotExist(err) {
				fmt.Printf("Error: Provider directory '%s' does not exist\n", providerDir)
				return
			}

			// Load single provider
			pwt, err := loadProviderFromNewStructure(providerDir)
			if err != nil {
				fmt.Printf("Failed to load provider %s: %v\n", providerName, err)
				return
			}
			providersWithTranslations = append(providersWithTranslations, pwt)
		}
	}

	// Print loaded providers
	for _, pwt := range providersWithTranslations {
		provider := pwt.Provider
		fmt.Printf("Loaded provider: %s (%s)\n", provider.Name, provider.Identifier)
		fmt.Printf("  Description: %s\n", provider.Description)
		fmt.Printf("  Auth Type: %s\n", provider.AuthType)
		fmt.Printf("  Categories: %v\n", provider.Categories)
		fmt.Printf("  Operations: %d\n", len(provider.Operations))
		fmt.Printf("  Permissions: %d\n", len(pwt.Adapter.Permissions))
		fmt.Printf("  Translations: %d\n", len(pwt.Translations))

		fmt.Println()
	}

	if *update {
		if *providerId == "" {
			fmt.Printf("Error: When --update is specified, you must provide --provider-id parameter\n")
			fmt.Printf("Usage examples:\n")
			fmt.Printf("  %s --update --provider-id 123                           # Generate SQL for notion provider only\n", os.Args[0])
			return
		}
		if len(providersWithTranslations) != 1 {
			fmt.Printf("Error: When --update is specified, you must provide only one provider\n")
			return
		}
		// Generate Update SQL files
		if err := generateUpdateSQLFiles(*providerId, providersWithTranslations, *sqlOutputDir); err != nil {
			fmt.Printf("SQL file generation failed: %v\n", err)
		}
	} else {
		// Generate Insert SQL files
		if err := generateInsertSQLFiles(providersWithTranslations, *sqlOutputDir); err != nil {
			fmt.Printf("SQL file generation failed: %v\n", err)
		}
	}
}

// loadProvidersFromJSON loads all provider JSON files from the specified directory
// Supports both old structure (*.json files) and new structure (provider_name/manifest.json + i18n/*.json)
func loadProvidersFromJSON(dirPath string) ([]*ProviderWithTranslations, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var providersWithTranslations []*ProviderWithTranslations

	for _, file := range files {
		if shouldSkipFile(file.Name()) {
			continue
		}

		if file.IsDir() {
			// New structure: provider_name/manifest.json + i18n/*.json
			providerDir := filepath.Join(dirPath, file.Name())
			manifestPath := filepath.Join(providerDir, "manifest.json")

			// Check if manifest.json exists
			if _, err := os.Stat(manifestPath); err == nil {
				pwt, err := loadProviderFromNewStructure(providerDir)
				if err != nil {
					return nil, fmt.Errorf("failed to load provider from new structure %s: %w", providerDir, err)
				}
				providersWithTranslations = append(providersWithTranslations, pwt)
			}
		} else {
			// Skip non-directory files (like old JSON files, README, etc.)
			// Only directories are expected to contain provider configurations
			fmt.Printf("Skipping non-directory file: %s\n", file.Name())
			continue
		}
	}

	return providersWithTranslations, nil
}

// loadProviderFromNewStructure loads a provider from the new directory structure
// (provider_name/manifest.json + i18n/*.json)
func loadProviderFromNewStructure(providerDir string) (*ProviderWithTranslations, error) {
	manifestPath := filepath.Join(providerDir, "manifest.json")

	// Load provider from manifest.json
	provider, adapter, err := loadProviderFromFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load manifest.json: %w", err)
	}

	// Load translations from i18n directory
	translations, err := loadTranslationsFromDir(filepath.Join(providerDir, "i18n"), provider.Identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to load translations: %w", err)
	}

	return &ProviderWithTranslations{
		Provider:     provider,
		Adapter:      adapter,
		Translations: translations,
	}, nil
}

// loadTranslationsFromDir loads all translation files from the i18n directory
func loadTranslationsFromDir(i18nDir string, providerIdentifier string) ([]TranslationData, error) {
	var translations []TranslationData

	// Check if i18n directory exists
	if _, err := os.Stat(i18nDir); os.IsNotExist(err) {
		// No translations directory, return empty slice
		return translations, nil
	}

	files, err := os.ReadDir(i18nDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read i18n directory %s: %w", i18nDir, err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			// Extract language code from filename (e.g., "zh-TW.json" -> "zh-TW")
			languageCode := strings.TrimSuffix(file.Name(), ".json")

			// Read translation file
			translationPath := filepath.Join(i18nDir, file.Name())
			translationData, err := os.ReadFile(translationPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read translation file %s: %w", translationPath, err)
			}

			// Validate JSON format and convert to compact JSON string
			var temp interface{}
			if err := sonic.Unmarshal(translationData, &temp); err != nil {
				return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", translationPath, err)
			}

			// Convert to compact JSON string
			compactJSON, err := sonic.Marshal(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal compact JSON: %w", err)
			}

			translations = append(translations, TranslationData{
				ID:                 uuid.New().String(),
				ProviderIdentifier: providerIdentifier,
				LanguageCode:       languageCode,
				Translations:       string(compactJSON),
			})
		}
	}

	return translations, nil
}

// loadProviderFromFile loads a provider from a JSON file
func loadProviderFromFile(filePath string) (*domain.Provider, *adapter_domain.ProviderAdapterConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	providerID := uuid.New().String()

	var providerJSON ProviderJSON
	if err := sonic.Unmarshal(data, &providerJSON); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", filePath, err)
	}

	// Convert permissions from JSON to domain model
	permissions := make([]types.Permission, 0, len(providerJSON.Permissions))
	for _, permJSON := range providerJSON.Permissions {
		perm := adapter_domain.NewPermission(permJSON.Identifier, permJSON.Name, permJSON.Description, permJSON.OAuthScopes)
		permissions = append(permissions, *perm)
	}

	// Convert operations from JSON to domain model
	operations := make([]domain.Operation, 0, len(providerJSON.Operations))
	for _, opJSON := range providerJSON.Operations {
		// Convert required permissions identifiers to Permission objects
		var requiredPermissions []types.Permission
		if len(opJSON.RequiredPermissions) > 0 {
			// Create a map for quick lookup of permissions by identifier
			permMap := make(map[string]types.Permission)
			for _, perm := range permissions {
				permMap[perm.Identifier] = perm
			}

			// Convert required permission identifiers to Permission objects
			for _, permID := range opJSON.RequiredPermissions {
				if perm, ok := permMap[permID]; ok {
					requiredPermissions = append(requiredPermissions, perm)
				}
			}
		}

		// Convert parameters from JSON to domain model
		parameters := make([]domain.Parameter, 0, len(opJSON.Parameters))
		for _, paramJSON := range opJSON.Parameters {
			// Convert string type to domain ParameterType
			paramType := domain.ParameterType(paramJSON.Type)

			// Create domain Parameter
			param := domain.NewParameter(paramJSON.Name, paramType, paramJSON.Description, paramJSON.Required, paramJSON.Enum, paramJSON.Default)
			parameters = append(parameters, *param)
		}

		op := domain.NewOperation(opJSON.Identifier, providerID, opJSON.Name, opJSON.Description, opJSON.Category, requiredPermissions, parameters)
		operations = append(operations, *op)
	}

	// Use default status if not specified
	status := types.ProviderStatusActive
	if providerJSON.Status != "" {
		status = types.ProviderStatus(providerJSON.Status)
	}

	// Create provider
	provider := &domain.Provider{
		ID:          providerID,
		Identifier:  providerJSON.Identifier,
		Name:        providerJSON.Name,
		Description: providerJSON.Description,
		AuthType:    types.ProviderAuthType(providerJSON.AuthType),
		Status:      status,
		IconURL:     providerJSON.IconURL,
		Categories:  providerJSON.Categories,
		Operations:  operations,
	}

	adapter := &adapter_domain.ProviderAdapterConfig{
		ProviderAdapterInfo: adapter_domain.ProviderAdapterInfo{
			Identifier:  provider.Identifier,
			Name:        provider.Name,
			Description: provider.Description,
		},
		ID:          providerID,
		Permissions: permissions,
	}
	if providerJSON.OAuthConfig != nil {
		adapter.OAuthConfig = providerJSON.OAuthConfig
	}
	adapter.CustomConfig = map[string]interface{}{}
	if providerJSON.ApiKeyConfig != nil {
		adapter.CustomConfig["api_key"] = providerJSON.ApiKeyConfig.Value
	}
	if providerJSON.VolcengineCredentials != nil {
		adapter.CustomConfig["volcengine_credentials"] = providerJSON.VolcengineCredentials
	}
	if providerJSON.OpenaiCredentials != nil {
		adapter.CustomConfig["openai_credentials"] = providerJSON.OpenaiCredentials
	}
	if providerJSON.KnowledgebaseConfig != nil {
		adapter.CustomConfig["knowledgebase_config"] = providerJSON.KnowledgebaseConfig
	}
	if len(adapter.CustomConfig) == 0 {
		adapter.CustomConfig = nil
	}
	return provider, adapter, nil
}

// generateInsertSQLFiles generates SQL insert files for the four tables
func generateInsertSQLFiles(providersWithTranslations []*ProviderWithTranslations, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate SQL for each table
	if err := generateProvidersInsertSQL(providersWithTranslations, nil, outputDir); err != nil {
		return fmt.Errorf("failed to generate providers SQL: %w", err)
	}

	if err := generateOperationsInsertSQL(providersWithTranslations, nil, outputDir); err != nil {
		return fmt.Errorf("failed to generate operations SQL: %w", err)
	}

	if err := generateProviderAdaptersInsertSQL(providersWithTranslations, nil, outputDir); err != nil {
		return fmt.Errorf("failed to generate provider_adapters SQL: %w", err)
	}

	if err := generateProviderTranslationsInsertSQL(providersWithTranslations, nil, outputDir); err != nil {
		return fmt.Errorf("failed to generate provider_translations SQL: %w", err)
	}

	return nil
}

// generateUpdateSQLFiles generates SQL update files for the four tables
func generateUpdateSQLFiles(oldProviderId string, providersWithTranslations []*ProviderWithTranslations, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(filepath.Join(outputDir, "update_providers.sql"))
	if err != nil {
		return fmt.Errorf("failed to create update providers SQL file: %w", err)
	}
	defer file.Close()

	// Write header
	_, err = file.WriteString("-- Generated SQL update statements with transaction\n")
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	_, err = file.WriteString("-- Generated at: " + time.Now().Format(time.RFC3339) + "\n\n")
	if err != nil {
		return fmt.Errorf("failed to write timestamp: %w", err)
	}

	// Start transaction
	_, err = file.WriteString("BEGIN;\n\n")
	if err != nil {
		return fmt.Errorf("failed to write transaction start: %w", err)
	}

	// Delete old provider
	if err := generateProvidersDeleteSQL(oldProviderId, file); err != nil {
		return fmt.Errorf("failed to generate providers delete SQL: %w", err)
	}
	// Delete old provider operations
	if err := generateOperationsDeleteSQL(oldProviderId, file); err != nil {
		return fmt.Errorf("failed to generate operations delete SQL: %w", err)
	}

	// Insert new provider
	if err := generateProvidersInsertSQL(providersWithTranslations, file, outputDir); err != nil {
		return fmt.Errorf("failed to generate providers insert SQL: %w", err)
	}
	// Insert new provider operations
	if err := generateOperationsInsertSQL(providersWithTranslations, file, outputDir); err != nil {
		return fmt.Errorf("failed to generate operations insert SQL: %w", err)
	}

	// Update provider_adapters
	if err := generateProviderAdaptersUpdateSQL(providersWithTranslations[0], file); err != nil {
		return fmt.Errorf("failed to generate provider_adapters update SQL: %w", err)
	}
	// Update provider_translations
	if err := generateProviderTranslationsUpdateSQL(providersWithTranslations[0], file); err != nil {
		return fmt.Errorf("failed to generate provider_translations update SQL: %w", err)
	}

	// Commit transaction
	_, err = file.WriteString("\nCOMMIT;\n")
	if err != nil {
		return fmt.Errorf("failed to write transaction commit: %w", err)
	}

	return nil
}

// generateProvidersSQL generates INSERT statements for the providers table
func generateProvidersInsertSQL(providersWithTranslations []*ProviderWithTranslations, file *os.File, outputDir string) error {
	if file == nil {
		fileName := filepath.Join(outputDir, "providers_inserts.sql")
		var err error
		file, err = os.Create(fileName)
		if err != nil {
			return fmt.Errorf("failed to create providers SQL file: %w", err)
		}
		defer file.Close()
		// Write header
		_, err = file.WriteString("-- Generated SQL insert statements for providers table\n")
		if err != nil {
			return err
		}
		_, err = file.WriteString("-- Generated at: " + time.Now().Format(time.RFC3339) + "\n\n")
		if err != nil {
			return err
		}
	}

	ctx := context.Background()

	for _, pwt := range providersWithTranslations {
		provider := pwt.Provider

		// Create JSON attributes with categories and permissions count
		jsonAttributes := map[string]interface{}{
			"categories": provider.Categories,
		}
		jsonAttributesData, err := sonic.Marshal(jsonAttributes)
		if err != nil {
			return fmt.Errorf("failed to marshal json_attributes: %w", err)
		}

		// Generate embedding for provider description if OpenAI client is available
		var embeddingStr string = "NULL"
		if openaiClient != nil {
			fmt.Printf("Generating embedding for provider: %s\n", provider.Identifier)
			embedding, err := generateEmbedding(ctx, provider.Description)
			if err != nil {
				fmt.Printf("Warning: Failed to generate embedding for provider %s: %v\n", provider.Identifier, err)
				embeddingStr = "NULL"
			} else if embedding != nil {
				embeddingStr = "'" + formatEmbeddingForPostgres(embedding) + "'"
				fmt.Printf("Successfully generated embedding for provider: %s (%d dimensions)\n", provider.Identifier, len(embedding))
			}
		}

		// Generate INSERT statement
		sql := fmt.Sprintf(`INSERT INTO providers (id, identifier, name, description, auth_type, status, icon_url, json_attributes, embedding)
VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %s);

`,
			escapeSQL(provider.ID),
			escapeSQL(provider.Identifier),
			escapeSQL(provider.Name),
			escapeSQL(provider.Description),
			escapeSQL(string(provider.AuthType)),
			escapeSQL(string(provider.Status)),
			escapeSQL(provider.IconURL),
			escapeSQL(string(jsonAttributesData)),
			embeddingStr,
		)

		if _, err := file.WriteString(sql); err != nil {
			return fmt.Errorf("failed to write providers SQL statement: %w", err)
		}
	}

	return nil
}

// generateOperationsSQL generates INSERT statements for the operations table
func generateOperationsInsertSQL(providersWithTranslations []*ProviderWithTranslations, file *os.File, outputDir string) error {
	if file == nil {
		fileName := filepath.Join(outputDir, "operations_inserts.sql")
		var err error
		file, err = os.Create(fileName)
		if err != nil {
			return fmt.Errorf("failed to create operations SQL file: %w", err)
		}
		defer file.Close()
		// Write header
		_, err = file.WriteString("-- Generated SQL insert statements for operations table\n")
		if err != nil {
			return err
		}
		_, err = file.WriteString("-- Generated at: " + time.Now().Format(time.RFC3339) + "\n\n")
		if err != nil {
			return err
		}
	}

	ctx := context.Background()

	for _, pwt := range providersWithTranslations {
		provider := pwt.Provider

		for _, operation := range provider.Operations {
			// Create JSON attributes with required_permissions and parameters
			jsonAttributes := map[string]interface{}{
				"required_permissions": operation.RequiredPermissions,
				"parameters":           operation.Parameters,
			}
			jsonAttributesData, err := sonic.Marshal(jsonAttributes)
			if err != nil {
				return fmt.Errorf("failed to marshal operation json_attributes: %w", err)
			}

			// Generate embedding for operation description if OpenAI client is available
			var embeddingStr string = "NULL"
			if openaiClient != nil {
				fmt.Printf("Generating embedding for operation: %s.%s\n", provider.Identifier, operation.Identifier)
				embedding, err := generateEmbedding(ctx, operation.Description)
				if err != nil {
					fmt.Printf("Warning: Failed to generate embedding for operation %s.%s: %v\n", provider.Identifier, operation.Identifier, err)
					embeddingStr = "NULL"
				} else if embedding != nil {
					embeddingStr = "'" + formatEmbeddingForPostgres(embedding) + "'"
					fmt.Printf("Successfully generated embedding for operation: %s.%s (%d dimensions)\n", provider.Identifier, operation.Identifier, len(embedding))
				}
			}

			// Generate INSERT statement
			sql := fmt.Sprintf(`INSERT INTO operations (id, identifier, provider_id, name, description, category, json_attributes, embedding)
VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', %s);

`,
				escapeSQL(operation.ID),
				escapeSQL(operation.Identifier),
				escapeSQL(operation.ProviderID),
				escapeSQL(operation.Name),
				escapeSQL(operation.Description),
				escapeSQL(operation.Category),
				escapeSQL(string(jsonAttributesData)),
				embeddingStr,
			)

			if _, err := file.WriteString(sql); err != nil {
				return fmt.Errorf("failed to write operations SQL statement: %w", err)
			}
		}
	}

	return nil
}

// generateProviderAdaptersSQL generates INSERT statements for the provider_adapters table
func generateProviderAdaptersInsertSQL(providersWithTranslations []*ProviderWithTranslations, file *os.File, outputDir string) error {
	if file == nil {
		fileName := filepath.Join(outputDir, "provider_adapters_inserts.sql")
		var err error
		file, err = os.Create(fileName)
		if err != nil {
			return fmt.Errorf("failed to create provider_adapters SQL file: %w", err)
		}
		defer file.Close()

		// Write header
		_, err = file.WriteString("-- Generated SQL insert statements for provider_adapters table\n")
		if err != nil {
			return err
		}
		_, err = file.WriteString("-- Generated at: " + time.Now().Format(time.RFC3339) + "\n\n")
		if err != nil {
			return err
		}
	}

	for _, pwt := range providersWithTranslations {
		adapter := pwt.Adapter

		// Create configs JSON with oauth_config and custom_config
		configs := map[string]interface{}{}
		if adapter.OAuthConfig != nil {
			configs["oauth_config"] = adapter.OAuthConfig
		}
		if adapter.CustomConfig != nil {
			configs["custom_config"] = adapter.CustomConfig
		}

		configsData, err := sonic.Marshal(configs)
		if err != nil {
			return fmt.Errorf("failed to marshal adapter configs: %w", err)
		}

		permissionsData, err := sonic.Marshal(adapter.Permissions)
		if err != nil {
			return fmt.Errorf("failed to marshal adapter permissions: %w", err)
		}

		// Generate INSERT statement
		sql := fmt.Sprintf(`INSERT INTO provider_adapters (id, identifier, configs, permissions)
VALUES ('%s', '%s', '%s', '%s');

`,
			escapeSQL(adapter.ID),
			escapeSQL(adapter.Identifier),
			escapeSQL(string(configsData)),
			escapeSQL(string(permissionsData)),
		)

		if _, err := file.WriteString(sql); err != nil {
			return fmt.Errorf("failed to write provider_adapters SQL statement: %w", err)
		}
	}

	return nil
}

// generateProviderTranslationsSQL generates INSERT statements for the provider_translations table
func generateProviderTranslationsInsertSQL(providersWithTranslations []*ProviderWithTranslations, file *os.File, outputDir string) error {
	if file == nil {
		fileName := filepath.Join(outputDir, "provider_translations_inserts.sql")
		var err error
		file, err = os.Create(fileName)
		if err != nil {
			return fmt.Errorf("failed to create provider_translations SQL file: %w", err)
		}
		defer file.Close()

		// Write header
		_, err = file.WriteString("-- Generated SQL insert statements for provider_translations table\n")
		if err != nil {
			return err
		}
		_, err = file.WriteString("-- Generated at: " + time.Now().Format(time.RFC3339) + "\n")
		if err != nil {
			return err
		}
		_, err = file.WriteString("-- Note: translations field contains compact JSON string (no whitespace/indentation)\n\n")
		if err != nil {
			return err
		}
	}

	for _, pwt := range providersWithTranslations {
		for _, translation := range pwt.Translations {
			// Generate INSERT statement (translations is already a compact JSON string)
			sql := fmt.Sprintf(`INSERT INTO provider_translations (id, provider_identifier, language_code, translations)
VALUES ('%s', '%s', '%s', '%s');

`,
				escapeSQL(translation.ID),
				escapeSQL(translation.ProviderIdentifier),
				escapeSQL(translation.LanguageCode),
				escapeSQL(translation.Translations),
			)

			if _, err := file.WriteString(sql); err != nil {
				return fmt.Errorf("failed to write provider_translations SQL statement: %w", err)
			}
		}
	}

	return nil
}

func generateProvidersDeleteSQL(oldProviderId string, file *os.File) error {
	sql := fmt.Sprintf(`UPDATE providers SET deleted_at = NOW() WHERE id = '%s';

`,
		escapeSQL(oldProviderId),
	)
	if _, err := file.WriteString(sql); err != nil {
		return fmt.Errorf("failed to write providers delete SQL statement: %w", err)
	}
	return nil
}

func generateOperationsDeleteSQL(oldProviderId string, file *os.File) error {
	sql := fmt.Sprintf(`UPDATE operations SET deleted_at = NOW() WHERE provider_id = '%s';

`,
		escapeSQL(oldProviderId),
	)
	if _, err := file.WriteString(sql); err != nil {
		return fmt.Errorf("failed to write operations delete SQL statement: %w", err)
	}
	return nil
}

func generateProviderAdaptersUpdateSQL(providersWithTranslations *ProviderWithTranslations, file *os.File) error {
	adapter := providersWithTranslations.Adapter
	configs := map[string]interface{}{}
	if adapter.OAuthConfig != nil {
		configs["oauth_config"] = adapter.OAuthConfig
	}
	if adapter.CustomConfig != nil {
		configs["custom_config"] = adapter.CustomConfig
	}

	configsData, err := sonic.Marshal(configs)
	if err != nil {
		return fmt.Errorf("failed to marshal adapter configs: %w", err)
	}

	permissionsData, err := sonic.Marshal(adapter.Permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal adapter permissions: %w", err)
	}

	sql := fmt.Sprintf(`UPDATE provider_adapters SET configs = '%s', permissions = '%s' WHERE identifier = '%s';

`,
		escapeSQL(string(configsData)),
		escapeSQL(string(permissionsData)),
		escapeSQL(adapter.Identifier),
	)
	if _, err := file.WriteString(sql); err != nil {
		return fmt.Errorf("failed to write provider_adapters update SQL statement: %w", err)
	}

	return nil
}

func generateProviderTranslationsUpdateSQL(providersWithTranslations *ProviderWithTranslations, file *os.File) error {
	for _, translation := range providersWithTranslations.Translations {
		sql := fmt.Sprintf(`UPDATE provider_translations SET translations = '%s' WHERE provider_identifier = '%s' AND language_code = '%s';

`,
			escapeSQL(translation.Translations),
			escapeSQL(translation.ProviderIdentifier),
			escapeSQL(translation.LanguageCode),
		)
		if _, err := file.WriteString(sql); err != nil {
			return fmt.Errorf("failed to write provider_translations update SQL statement: %w", err)
		}
	}
	return nil
}

// escapeSQL escapes single quotes in SQL strings to prevent SQL injection
func escapeSQL(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// generateEmbedding generates embedding for the given text using OpenAI API
func generateEmbedding(ctx context.Context, text string) ([]float64, error) {
	if openaiClient == nil {
		return nil, fmt.Errorf("OpenAI client not initialized")
	}

	if text == "" {
		return nil, nil
	}

	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.SmallEmbedding3,
	}

	resp, err := openaiClient.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}

	// Convert from []float32 to []float64
	embedding := make([]float64, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		embedding[i] = float64(v)
	}

	return embedding, nil
}

// formatEmbeddingForPostgres converts a float64 slice to PostgreSQL vector format
func formatEmbeddingForPostgres(embedding []float64) string {
	if len(embedding) == 0 {
		return "NULL"
	}

	var builder strings.Builder
	builder.WriteString("[")

	for i, val := range embedding {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("%f", val))
	}

	builder.WriteString("]")
	return builder.String()
}
