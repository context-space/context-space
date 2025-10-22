package domain

import (
	"github.com/context-space/context-space/backend/internal/shared/types"
)

type ProviderTranslation struct {
	Identifier   string             `json:"identifier"`
	LanguageCode string             `json:"language_code"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Categories   []string           `json:"categories"`
	Operations   []Operation        `json:"operations"`
	Permissions  []types.Permission `json:"permissions"`
}

type Operation struct {
	Name        string      `json:"name"`
	Identifier  string      `json:"identifier"`
	Description string      `json:"description"`
	Parameters  []Parameter `json:"parameters"`
}

// Parameter represents a parameter for an operation
type Parameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
