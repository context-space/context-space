package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/mcp"
	"github.com/context-space/context-space/backend/internal/providercore/domain"
)

// ManifestData represents the structure of manifest.json
type ManifestData struct {
	Identifier  string               `json:"identifier"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	AuthType    string               `json:"auth_type"`
	IconURL     string               `json:"icon_url"`
	Categories  []string             `json:"categories"`
	Permissions []ManifestPermission `json:"permissions"`
	Operations  []ManifestOperation  `json:"operations"`
}

// ManifestPermission represents a permission in the manifest
type ManifestPermission struct {
	Identifier  string   `json:"identifier"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	OAuthScopes []string `json:"oauth_scopes"`
}

// ManifestOperation represents an operation in the manifest
type ManifestOperation struct {
	Identifier          string              `json:"identifier"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	Category            string              `json:"category"`
	RequiredPermissions []string            `json:"required_permissions"`
	Parameters          []ManifestParameter `json:"parameters"`
}

// ManifestParameter represents a parameter in the manifest
type ManifestParameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Enum        []string    `json:"enum,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

// TranslationPermission represents a permission in the i18n data
type TranslationPermission struct {
	Identifier  string `json:"identifier"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TranslationOperation represents an operation in the i18n data
type TranslationOperation struct {
	Identifier  string                 `json:"identifier"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  []TranslationParameter `json:"parameters"`
}

// TranslationParameter represents a parameter in the i18n data
type TranslationParameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// convertMCPToolsToOperations converts MCP tools to domain operations
func convertMCPToolsToOperations(tools []mcp.MCPTool) []domain.Operation {
	var operations []domain.Operation

	for _, tool := range tools {
		// Convert parameters
		var parameters []domain.Parameter

		if tool.InputSchema.Properties != nil {
			for name, prop := range tool.InputSchema.Properties {
				required := contains(tool.InputSchema.Required, name)

				// Convert enum values to strings
				var enumValues []string
				for _, val := range prop.Enum {
					if str, ok := val.(string); ok {
						enumValues = append(enumValues, str)
					}
				}

				param := domain.Parameter{
					Name:        name,
					Type:        convertParameterType(prop.Type),
					Description: prop.Description,
					Required:    required,
					Enum:        enumValues,
					Default:     prop.Default,
				}

				parameters = append(parameters, param)
			}
		}

		// Create operation
		operation := domain.Operation{
			Identifier:          tool.Name,
			Name:                formatOperationName(tool.Name),
			Description:         tool.Description,
			Category:            "mcp_tools",           // Default category for MCP tools
			RequiredPermissions: []domain.Permission{}, // Will be set later if needed
			Parameters:          parameters,
		}

		operations = append(operations, operation)
	}

	return operations
}

func convertParameterType(mcpType string) domain.ParameterType {
	switch mcpType {
	case "string":
		return domain.ParameterTypeString
	case "integer":
		return domain.ParameterTypeInteger
	case "number":
		return domain.ParameterTypeNumber
	case "boolean":
		return domain.ParameterTypeBoolean
	case "array":
		return domain.ParameterTypeArray
	case "object":
		return domain.ParameterTypeObject
	default:
		return domain.ParameterTypeString // Default to string for unknown types
	}
}

func formatOperationName(name string) string {
	// Convert snake_case or kebab-case to Title Case
	words := strings.FieldsFunc(name, func(c rune) bool {
		return c == '_' || c == '-'
	})

	var result []string
	for _, word := range words {
		if word != "" {
			result = append(result, strings.Title(strings.ToLower(word)))
		}
	}

	return strings.Join(result, " ")
}

// generateManifest creates manifest with specified provider details
func generateManifest(identifier, name, description, authType string, categories []string, operations []domain.Operation) *ManifestData {
	// Create default permission for MCP tools access
	defaultPermission := ManifestPermission{
		Identifier:  "mcp_tools_access",
		Name:        "MCP Tools Access",
		Description: "Access to MCP server tools",
		OAuthScopes: []string{},
	}

	// Convert operations to manifest format
	var manifestOps []ManifestOperation
	for _, op := range operations {
		var manifestParams []ManifestParameter
		for _, param := range op.Parameters {
			manifestParam := ManifestParameter{
				Name:        param.Name,
				Type:        string(param.Type),
				Description: param.Description,
				Required:    param.Required,
				Enum:        param.Enum,
				Default:     param.Default,
			}
			manifestParams = append(manifestParams, manifestParam)
		}

		manifestOp := ManifestOperation{
			Identifier:          op.Identifier,
			Name:                op.Name,
			Description:         op.Description,
			Category:            op.Category,
			RequiredPermissions: []string{"mcp_tools_access"},
			Parameters:          manifestParams,
		}
		manifestOps = append(manifestOps, manifestOp)
	}

	manifest := &ManifestData{
		Identifier:  identifier,
		Name:        name,
		Description: description,
		AuthType:    authType,
		IconURL:     "",
		Categories:  categories,
		Permissions: []ManifestPermission{defaultPermission},
		Operations:  manifestOps,
	}

	return manifest
}

// writeJSONFile writes data to a JSON file
func writeJSONFile(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(filename, jsonData, 0644)
}

// generateI18nFile creates i18n data for a specific language
func generateI18nFile(name, description string, operations []domain.Operation, categories []string, lang string, enableTranslation bool) (*TranslationResponse, error) {
	// Create permissions for MCP tools
	permissions := []TranslationPermission{
		{
			Identifier:  "mcp_tools_access",
			Name:        "MCP Tools Access",
			Description: "Access to MCP server tools",
		},
	}

	// Convert operations to translation format
	var translationOps []TranslationOperation
	for _, op := range operations {
		var translationParams []TranslationParameter
		for _, param := range op.Parameters {
			translationParams = append(translationParams, TranslationParameter{
				Name:        param.Name,
				Description: param.Description,
			})
		}

		translationOps = append(translationOps, TranslationOperation{
			Identifier:  op.Identifier,
			Name:        op.Name,
			Description: op.Description,
			Parameters:  translationParams,
		})
	}

	// Create base i18n data structure matching airtable format
	i18nData := &TranslationResponse{
		Name:        name,
		Description: description,
		Categories:  categories,
		Permissions: permissions,
		Operations:  translationOps,
	}

	// If English or translation is disabled, return the base data
	if lang == "en" || !enableTranslation {
		return i18nData, nil
	}

	// For non-English languages with translation enabled, use AI translation
	if !IsTranslationEnabled() {
		fmt.Printf("Warning: Translation requested for %s but OPENAI_API_KEY not found. Falling back to English.\n", lang)
		return i18nData, nil
	}

	fmt.Printf("Translating i18n data to %s using AI...\n", lang)

	// Create AI translator
	translator, err := NewAITranslator()
	if err != nil {
		fmt.Printf("Warning: Failed to create AI translator: %v. Falling back to English.\n", err)
		return i18nData, nil
	}

	// Translate the data
	translatedData, err := translator.TranslateI18nData(context.Background(), *i18nData, lang)
	if err != nil {
		fmt.Printf("Warning: AI translation failed for %s: %v. Falling back to English.\n", lang, err)
		return i18nData, nil
	}

	fmt.Printf("Successfully translated i18n data to %s\n", lang)
	return &translatedData, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
