package provider

import "github.com/context-space/context-space/backend/internal/shared/types"

// ProviderDTO Provider data transfer object
type ProviderDTO struct {
	ID          string         `json:"id"`
	Identifier  string         `json:"identifier"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	AuthType    string         `json:"auth_type"`
	Status      string         `json:"status"`
	Tags        []string       `json:"tags"`
	IconURL     string         `json:"icon_url"`
	ApiDocURL   string         `json:"api_doc_url"`
	Categories  []string       `json:"categories"`
	Operations  []OperationDTO `json:"operations"`
	Embedding   []float64      `json:"embedding"`
	CreatedAt   int64          `json:"created_at"`
	UpdatedAt   int64          `json:"updated_at"`
}

// OperationDTO Operation data transfer object
type OperationDTO struct {
	ID                  string             `json:"id"`
	ProviderID          string             `json:"provider_id"`
	Identifier          string             `json:"identifier"`
	Name                string             `json:"name"`
	Description         string             `json:"description"`
	Category            string             `json:"category"`
	RequiredPermissions []types.Permission `json:"required_permissions"`
	Parameters          []ParameterDTO     `json:"parameters"`
	CreatedAt           int64              `json:"created_at"`
	UpdatedAt           int64              `json:"updated_at"`
}

type ParameterDTO struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Enum        []string    `json:"enum,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}
