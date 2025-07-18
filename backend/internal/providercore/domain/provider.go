package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/language"
)

// ProviderAuthType represents the authentication type for a provider
type ProviderAuthType string

const (
	AuthTypeOAuth  ProviderAuthType = "oauth"
	AuthTypeAPIKey ProviderAuthType = "apikey"
	AuthTypeBasic  ProviderAuthType = "basic"
	AuthTypeNone   ProviderAuthType = "none"
)

// ProviderStatus represents the status of a provider
type ProviderStatus string

const (
	ProviderStatusActive      ProviderStatus = "active"
	ProviderStatusInactive    ProviderStatus = "inactive"
	ProviderStatusMaintenance ProviderStatus = "maintenance"
	ProviderStatusDeprecated  ProviderStatus = "deprecated"
)

type TranslatedProvider struct {
	ID          string           `json:"id"`
	Identifier  string           `json:"identifier"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	AuthType    ProviderAuthType `json:"auth_type"`
	Status      ProviderStatus   `json:"status"`
	IconURL     string           `json:"icon_url"`
	ApiDocURL   string           `json:"api_doc_url"`
	Categories  []string         `json:"categories"`
	Tags        []string         `json:"tags"`
	Permissions []Permission     `json:"permissions"`
	Operations  []Operation      `json:"operations"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   *time.Time       `json:"deleted_at"`
	language    language.Tag
}

// Provider represents a third-party provider integration
type Provider struct {
	ID           string           `json:"id"`
	Identifier   string           `json:"identifier"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	AuthType     ProviderAuthType `json:"auth_type"`
	Status       ProviderStatus   `json:"status"`
	IconURL      string           `json:"icon_url"`
	Categories   []string         `json:"categories"`
	Tags         []string         `json:"tags"`
	Permissions  []Permission     `json:"permissions"`
	Operations   []Operation      `json:"operations"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    *time.Time       `json:"deleted_at"`
	translations map[language.Tag]TranslatedProvider
}

func (p *Provider) GetTranslation(lang language.Tag) TranslatedProvider {
	if p.translations == nil {
		p.translations = make(map[language.Tag]TranslatedProvider)
	}

	// 1. Try to get the requested language translation
	if translation, exists := p.translations[lang]; exists {
		// Ensure ApiDocURL is set for existing translations
		if translation.ApiDocURL == "" {
			translation.ApiDocURL = GetProviderAPIDocURL(p.Identifier)
		}
		return translation
	}

	// 2. If not found, try to get the parent language translation
	if parent := lang.Parent(); parent != language.Und {
		if translation, exists := p.translations[parent]; exists {
			// Ensure ApiDocURL is set for existing translations
			if translation.ApiDocURL == "" {
				translation.ApiDocURL = GetProviderAPIDocURL(p.Identifier)
			}
			return translation
		}
	}

	// 3. If parent language not found, try to get default language (English) translation
	if translation, exists := p.translations[language.English]; exists {
		// Ensure ApiDocURL is set for existing translations
		if translation.ApiDocURL == "" {
			translation.ApiDocURL = GetProviderAPIDocURL(p.Identifier)
		}
		return translation
	}

	// 4. If default language not found, construct TranslatedProvider using basic fields from Provider struct
	fallbackTranslation := TranslatedProvider{
		ID:          p.ID,
		Identifier:  p.Identifier,
		Name:        p.Name,
		Description: p.Description,
		AuthType:    p.AuthType,
		Status:      p.Status,
		IconURL:     p.IconURL,
		ApiDocURL:   GetProviderAPIDocURL(p.Identifier),
		Categories:  p.Categories,
		Tags:        p.Tags,
		Permissions: p.Permissions,
		Operations:  p.Operations,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		DeletedAt:   p.DeletedAt,
		language:    language.English, // Mark as basic data (using English as marker)
	}

	return fallbackTranslation
}

func (p *Provider) SetTranslation(lang language.Tag, translation TranslatedProvider) {
	if p.translations == nil {
		p.translations = make(map[language.Tag]TranslatedProvider)
	}
	translation.language = lang
	p.translations[lang] = translation
}

func (p *Provider) HasTag(tagName string) bool {
	for _, tag := range p.Tags {
		if tag == tagName {
			return true
		}
	}
	return false
}

// NewProvider creates a new provider
func NewProvider(identifier, name, description string, authType ProviderAuthType, status ProviderStatus, iconURL string, categories []string, permissions []Permission, operations []Operation) *Provider {
	return &Provider{
		ID:           uuid.New().String(),
		Identifier:   identifier,
		Name:         name,
		Description:  description,
		AuthType:     authType,
		Status:       status,
		IconURL:      iconURL,
		Categories:   categories,
		Permissions:  permissions,
		Operations:   operations,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		translations: make(map[language.Tag]TranslatedProvider),
	}
}

// Activate activates the provider
func (p *Provider) Activate() {
	p.Status = ProviderStatusActive
	p.UpdatedAt = time.Now()
}

// Deactivate deactivates the provider
func (p *Provider) Deactivate() {
	p.Status = ProviderStatusInactive
	p.UpdatedAt = time.Now()
}

// SetMaintenance puts the provider in maintenance mode
func (p *Provider) SetMaintenance() {
	p.Status = ProviderStatusMaintenance
	p.UpdatedAt = time.Now()
}

// Deprecate marks the provider as deprecated
func (p *Provider) Deprecate() {
	p.Status = ProviderStatusDeprecated
	p.UpdatedAt = time.Now()
}

// IsActive returns true if the provider is active
func (p *Provider) IsActive() bool {
	return p.Status == ProviderStatusActive
}

// IsInactive returns true if the provider is inactive
func (p *Provider) IsInactive() bool {
	return p.Status == ProviderStatusInactive
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
