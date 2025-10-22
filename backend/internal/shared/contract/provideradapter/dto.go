package provideradapter

import "github.com/context-space/context-space/backend/internal/shared/types"

// AdapterInfoDTO contains the adapter information for cross-module communication
// This is the contract DTO, used for communication between modules
type AdapterInfoDTO struct {
	Identifier  string `json:"identifier"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AuthType    string `json:"auth_type"`
	Status      string `json:"status"`
}

type ProviderAdapterInfoDTO struct {
	Identifier  string             `json:"identifier"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	AuthType    string             `json:"auth_type"`
	Status      string             `json:"status"`
	IconURL     string             `json:"icon_url"`
	Categories  []string           `json:"categories"`
	Permissions []types.Permission `json:"permissions"`
	Operations  []OperationDTO     `json:"operations"`
}

type OperationDTO struct {
	Identifier          string             `json:"identifier"`
	Name                string             `json:"name"`
	Description         string             `json:"description"`
	Category            string             `json:"category"`
	RequiredPermissions []types.Permission `json:"required_permissions"`
	Parameters          []ParameterDTO     `json:"parameters"`
}

type ParameterDTO struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Enum        []string    `json:"enum,omitempty"`
	Default     interface{} `json:"default"`
}
