package providertranslation

import "github.com/context-space/context-space/backend/internal/shared/types"

type ProviderTranslationDTO struct {
	Identifier   string             `json:"identifier"`
	LanguageCode string             `json:"language_code"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Categories   []string           `json:"categories"`
	Operations   []OperationDTO     `json:"operations"`
	Permissions  []types.Permission `json:"permissions"`
}

type OperationDTO struct {
	Name        string         `json:"name"`
	Identifier  string         `json:"identifier"`
	Description string         `json:"description"`
	Parameters  []ParameterDTO `json:"parameters"`
}

// Parameter represents a parameter for an operation
type ParameterDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
