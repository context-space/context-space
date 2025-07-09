package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	adapter_domain "github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/providercore/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

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
	generateSQL := flag.Bool("sql", false, "Generate SQL insert files for the four tables")
	sqlOutputDir := flag.String("sql-output", "generated_sql", "Directory to save generated SQL files")
	flag.Parse()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Load providers from JSON files (supporting both old and new structure)
	providersWithTranslations, err := loadProvidersFromJSON(*providersPath)
	if err != nil {
		fmt.Printf("Failed to load providers from JSON: %v\n", err)
		return
	}

	// Print loaded providers
	for _, pwt := range providersWithTranslations {
		provider := pwt.Provider
		fmt.Printf("Loaded provider: %s (%s)\n", provider.Name, provider.Identifier)
		fmt.Printf("  Description: %s\n", provider.Description)
		fmt.Printf("  Auth Type: %s\n", provider.AuthType)
		fmt.Printf("  Categories: %v\n", provider.Categories)
		fmt.Printf("  Operations: %d\n", len(provider.Operations))
		fmt.Printf("  Permissions: %d\n", len(provider.Permissions))
		fmt.Printf("  Translations: %d\n", len(pwt.Translations))

		// Print operations and their parameters
		for i, op := range provider.Operations {
			fmt.Printf("  Operation %d: %s (%s)\n", i+1, op.Name, op.Identifier)
			fmt.Printf("    Parameters: %d\n", len(op.Parameters))
			for j, param := range op.Parameters {
				fmt.Printf("      Param %d: %s (Type: %s, Required: %v)\n",
					j+1, param.Name, param.Type, param.Required)
			}
		}

		// Print translations
		for _, translation := range pwt.Translations {
			fmt.Printf("  Translation: %s\n", translation.LanguageCode)
		}

		fmt.Println()
	}

	// Generate SQL files if requested
	if *generateSQL {
		fmt.Println("Generating SQL insert files...")
		if err := generateSQLFiles(providersWithTranslations, *sqlOutputDir); err != nil {
			fmt.Printf("SQL file generation failed: %v\n", err)
			return
		} else {
			fmt.Printf("Successfully generated SQL files in directory: %s\n", *sqlOutputDir)
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
			if err := json.Unmarshal(translationData, &temp); err != nil {
				return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", translationPath, err)
			}

			// Convert to compact JSON string
			compactJSON, err := json.Marshal(temp)
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
	if err := json.Unmarshal(data, &providerJSON); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", filePath, err)
	}

	// Convert permissions from JSON to domain model
	permissions := make([]domain.Permission, 0, len(providerJSON.Permissions))
	for _, permJSON := range providerJSON.Permissions {
		perm := domain.NewPermission(permJSON.Identifier, permJSON.Name, permJSON.Description, permJSON.OAuthScopes)
		permissions = append(permissions, *perm)
	}

	// Convert operations from JSON to domain model
	operations := make([]domain.Operation, 0, len(providerJSON.Operations))
	for _, opJSON := range providerJSON.Operations {
		// Convert required permissions identifiers to Permission objects
		var requiredPermissions []domain.Permission
		if len(opJSON.RequiredPermissions) > 0 {
			// Create a map for quick lookup of permissions by identifier
			permMap := make(map[string]domain.Permission)
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
	status := domain.ProviderStatusActive
	if providerJSON.Status != "" {
		status = domain.ProviderStatus(providerJSON.Status)
	}

	// Create provider
	provider := &domain.Provider{
		ID:          providerID,
		Identifier:  providerJSON.Identifier,
		Name:        providerJSON.Name,
		Description: providerJSON.Description,
		AuthType:    domain.ProviderAuthType(providerJSON.AuthType),
		Status:      status,
		IconURL:     providerJSON.IconURL,
		Categories:  providerJSON.Categories,
		Permissions: permissions,
		Operations:  operations,
	}

	adapter := &adapter_domain.ProviderAdapterConfig{
		ProviderAdapterInfo: adapter_domain.ProviderAdapterInfo{
			Identifier:  provider.Identifier,
			Name:        provider.Name,
			Description: provider.Description,
		},
		ID: providerID,
	}
	if providerJSON.OAuthConfig != nil {
		adapter.OAuthConfig = providerJSON.OAuthConfig
	}
	if providerJSON.ApiKeyConfig != nil {
		adapter.CustomConfig = map[string]interface{}{
			"api_key": providerJSON.ApiKeyConfig.Value,
		}
	}
	if providerJSON.VolcengineCredentials != nil {
		adapter.CustomConfig = map[string]interface{}{
			"volcengine_credentials": providerJSON.VolcengineCredentials,
		}
	}
	if providerJSON.OpenaiCredentials != nil {
		adapter.CustomConfig = map[string]interface{}{
			"openai_credentials": providerJSON.OpenaiCredentials,
		}
	}
	return provider, adapter, nil
}

// generateSQLFiles generates SQL insert files for the four tables
func generateSQLFiles(providersWithTranslations []*ProviderWithTranslations, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate SQL for each table
	if err := generateProvidersSQL(providersWithTranslations, outputDir); err != nil {
		return fmt.Errorf("failed to generate providers SQL: %w", err)
	}

	if err := generateOperationsSQL(providersWithTranslations, outputDir); err != nil {
		return fmt.Errorf("failed to generate operations SQL: %w", err)
	}

	if err := generateProviderAdaptersSQL(providersWithTranslations, outputDir); err != nil {
		return fmt.Errorf("failed to generate provider_adapters SQL: %w", err)
	}

	if err := generateProviderTranslationsSQL(providersWithTranslations, outputDir); err != nil {
		return fmt.Errorf("failed to generate provider_translations SQL: %w", err)
	}

	return nil
}

// generateProvidersSQL generates INSERT statements for the providers table
func generateProvidersSQL(providersWithTranslations []*ProviderWithTranslations, outputDir string) error {
	fileName := filepath.Join(outputDir, "providers_inserts.sql")
	file, err := os.Create(fileName)
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

	for _, pwt := range providersWithTranslations {
		provider := pwt.Provider

		// Create JSON attributes with categories and permissions count
		jsonAttributes := map[string]interface{}{
			"categories":       provider.Categories,
			"permissions":      provider.Permissions,
			"operations_count": len(provider.Operations),
		}
		jsonAttributesData, err := json.Marshal(jsonAttributes)
		if err != nil {
			return fmt.Errorf("failed to marshal json_attributes: %w", err)
		}

		// Generate INSERT statement
		sql := fmt.Sprintf(`INSERT INTO providers (id, identifier, name, description, auth_type, status, icon_url, json_attributes)
VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');

`,
			escapeSQL(provider.ID),
			escapeSQL(provider.Identifier),
			escapeSQL(provider.Name),
			escapeSQL(provider.Description),
			escapeSQL(string(provider.AuthType)),
			escapeSQL(string(provider.Status)),
			escapeSQL(provider.IconURL),
			escapeSQL(string(jsonAttributesData)),
		)

		if _, err := file.WriteString(sql); err != nil {
			return fmt.Errorf("failed to write providers SQL statement: %w", err)
		}
	}

	return nil
}

// generateOperationsSQL generates INSERT statements for the operations table
func generateOperationsSQL(providersWithTranslations []*ProviderWithTranslations, outputDir string) error {
	fileName := filepath.Join(outputDir, "operations_inserts.sql")
	file, err := os.Create(fileName)
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

	for _, pwt := range providersWithTranslations {
		provider := pwt.Provider

		for _, operation := range provider.Operations {
			// Create JSON attributes with required_permissions and parameters
			jsonAttributes := map[string]interface{}{
				"required_permissions": operation.RequiredPermissions,
				"parameters":           operation.Parameters,
			}
			jsonAttributesData, err := json.Marshal(jsonAttributes)
			if err != nil {
				return fmt.Errorf("failed to marshal operation json_attributes: %w", err)
			}

			// Generate INSERT statement
			sql := fmt.Sprintf(`INSERT INTO operations (id, identifier, provider_id, name, description, category, json_attributes)
VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');

`,
				escapeSQL(operation.ID),
				escapeSQL(operation.Identifier),
				escapeSQL(operation.ProviderID),
				escapeSQL(operation.Name),
				escapeSQL(operation.Description),
				escapeSQL(operation.Category),
				escapeSQL(string(jsonAttributesData)),
			)

			if _, err := file.WriteString(sql); err != nil {
				return fmt.Errorf("failed to write operations SQL statement: %w", err)
			}
		}
	}

	return nil
}

// generateProviderAdaptersSQL generates INSERT statements for the provider_adapters table
func generateProviderAdaptersSQL(providersWithTranslations []*ProviderWithTranslations, outputDir string) error {
	fileName := filepath.Join(outputDir, "provider_adapters_inserts.sql")
	file, err := os.Create(fileName)
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

		configsData, err := json.Marshal(configs)
		if err != nil {
			return fmt.Errorf("failed to marshal adapter configs: %w", err)
		}

		// Generate INSERT statement
		sql := fmt.Sprintf(`INSERT INTO provider_adapters (id, identifier, configs)
VALUES ('%s', '%s', '%s');

`,
			escapeSQL(adapter.ID),
			escapeSQL(adapter.Identifier),
			escapeSQL(string(configsData)),
		)

		if _, err := file.WriteString(sql); err != nil {
			return fmt.Errorf("failed to write provider_adapters SQL statement: %w", err)
		}
	}

	return nil
}

// generateProviderTranslationsSQL generates INSERT statements for the provider_translations table
func generateProviderTranslationsSQL(providersWithTranslations []*ProviderWithTranslations, outputDir string) error {
	fileName := filepath.Join(outputDir, "provider_translations_inserts.sql")
	file, err := os.Create(fileName)
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

// escapeSQL escapes single quotes in SQL strings to prevent SQL injection
func escapeSQL(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
