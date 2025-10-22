package domain

import (
	"time"

	"github.com/context-space/context-space/backend/internal/shared/types"
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type TranslatedProvider struct {
	ID          string                 `json:"id"`
	Identifier  string                 `json:"identifier"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	AuthType    types.ProviderAuthType `json:"auth_type"`
	Status      types.ProviderStatus   `json:"status"`
	IconURL     string                 `json:"icon_url"`
	ApiDocURL   string                 `json:"api_doc_url"`
	Categories  []string               `json:"categories"`
	Tags        []string               `json:"tags"`
	Permissions []types.Permission     `json:"permissions"`
	Operations  []Operation            `json:"operations"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   *time.Time             `json:"deleted_at"`
	language    language.Tag
}

// Provider represents a third-party provider integration
type Provider struct {
	ID           string                 `json:"id"`
	Identifier   string                 `json:"identifier"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	AuthType     types.ProviderAuthType `json:"auth_type"`
	Status       types.ProviderStatus   `json:"status"`
	IconURL      string                 `json:"icon_url"`
	Categories   []string               `json:"categories"`
	Tags         []string               `json:"tags"`
	Operations   []Operation            `json:"operations"`
	Embedding    []float64              `json:"-"` // Vector embedding for semantic search
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	DeletedAt    *time.Time             `json:"deleted_at"`
	translations map[language.Tag]TranslatedProvider
}

// NewProvider creates a new provider
func NewProvider(identifier, name, description string, authType types.ProviderAuthType, status types.ProviderStatus, iconURL string, categories []string, operations []Operation) *Provider {
	return &Provider{
		ID:           uuid.New().String(),
		Identifier:   identifier,
		Name:         name,
		Description:  description,
		AuthType:     authType,
		Status:       status,
		IconURL:      iconURL,
		Categories:   categories,
		Operations:   operations,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		translations: make(map[language.Tag]TranslatedProvider),
	}
}

// Activate activates the provider
func (p *Provider) Activate() {
	p.Status = types.ProviderStatusActive
	p.UpdatedAt = time.Now()
}

// Deactivate deactivates the provider
func (p *Provider) Deactivate() {
	p.Status = types.ProviderStatusInactive
	p.UpdatedAt = time.Now()
}

// SetMaintenance puts the provider in maintenance mode
func (p *Provider) SetMaintenance() {
	p.Status = types.ProviderStatusMaintenance
	p.UpdatedAt = time.Now()
}

// Deprecate marks the provider as deprecated
func (p *Provider) Deprecate() {
	p.Status = types.ProviderStatusDeprecated
	p.UpdatedAt = time.Now()
}

// IsActive returns true if the provider is active
func (p *Provider) IsActive() bool {
	return p.Status == types.ProviderStatusActive
}

// IsInactive returns true if the provider is inactive
func (p *Provider) IsInactive() bool {
	return p.Status == types.ProviderStatusInactive
}

// HasCategory returns true if the provider belongs to the given category
func (p *Provider) HasCategory(category string) bool {
	for _, c := range p.Categories {
		if c == category {
			return true
		}
	}
	return false
}
