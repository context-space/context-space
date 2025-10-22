package domain

import (
	"time"

	"github.com/context-space/context-space/backend/internal/shared/types"
	"github.com/google/uuid"
)

// Operation represents an operation that can be performed on a provider
type Operation struct {
	ID                  string
	Identifier          string
	ProviderID          string
	Name                string
	Description         string
	Category            string
	RequiredPermissions []types.Permission
	Parameters          []Parameter
	Embedding           []float64 // Vector embedding for semantic search
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

// NewOperation creates a new operation
func NewOperation(identifier, providerID, name, description, category string, requiredPermissions []types.Permission, parameters []Parameter) *Operation {
	return &Operation{
		ID:                  uuid.New().String(),
		Identifier:          identifier,
		ProviderID:          providerID,
		Name:                name,
		Description:         description,
		Category:            category,
		RequiredPermissions: requiredPermissions,
		Parameters:          parameters,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}
