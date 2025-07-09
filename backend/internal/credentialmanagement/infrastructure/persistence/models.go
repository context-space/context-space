package persistence

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// CredentialModel represents the credentials table in the database
type CredentialModel struct {
	ID                 string         `gorm:"type:uuid;primary_key"`
	UserID             string         `gorm:"type:uuid;not null;index"`
	ProviderIdentifier string         `gorm:"type:varchar(50);not null;index"`
	CredentialType     string         `gorm:"type:credential_type;not null;index"`
	IsValid            bool           `gorm:"not null;default:true"`
	CreatedAt          time.Time      `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt          time.Time      `gorm:"type:timestamp with time zone;not null;default:now()"`
	LastUsedAt         time.Time      `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt          gorm.DeletedAt `gorm:"type:timestamp with time zone;index"`
}

// TableName returns the table name for the Credential model
func (CredentialModel) TableName() string {
	return "credentials"
}

// OAuthCredentialModel represents the oauth_credentials table in the database
type OAuthCredentialModel struct {
	CredentialID   string          `gorm:"type:uuid;primary_key"`
	Expiry         time.Time       `gorm:"type:timestamp with time zone; not null;index"`
	JSONAttributes json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt      time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt      time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt      gorm.DeletedAt  `gorm:"type:timestamp with time zone;index"`

	// Relationships
	Credential CredentialModel `gorm:"foreignKey:CredentialID;references:ID"`
}

// TableName returns the table name for the OAuthCredential model
func (OAuthCredentialModel) TableName() string {
	return "oauth_credentials"
}

// APIKeyCredentialModel represents the apikey_credentials table in the database
type APIKeyCredentialModel struct {
	CredentialID   string          `gorm:"type:uuid;primary_key"`
	JSONAttributes json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt      time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt      time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt      gorm.DeletedAt  `gorm:"type:timestamp with time zone;index"`

	// Relationships
	Credential CredentialModel `gorm:"foreignKey:CredentialID;references:ID"`
}

// TableName returns the table name for the APIKeyCredential model
func (APIKeyCredentialModel) TableName() string {
	return "apikey_credentials"
}

// OAuthStateModel represents the oauth_states table in the database
type OAuthStateModel struct {
	ID                 string          `gorm:"type:uuid;primary_key"`
	State              string          `gorm:"type:varchar(255);not null;uniqueIndex"`
	Status             string          `gorm:"type:varchar(20);not null;index"`
	UserID             string          `gorm:"type:uuid;not null;index"`
	ProviderIdentifier string          `gorm:"type:varchar(50);not null;index"`
	JSONAttributes     json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt          time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt          time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt          gorm.DeletedAt  `gorm:"type:timestamp with time zone;index"`
}

// TableName returns the table name for the OAuthState model
func (OAuthStateModel) TableName() string {
	return "oauth_states"
}

// BeforeCreate is called before creating a new record
func (c *CredentialModel) BeforeCreate(tx *gorm.DB) error {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate is called before updating an existing record
func (c *CredentialModel) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate is called before creating a new record
func (o *OAuthCredentialModel) BeforeCreate(tx *gorm.DB) error {
	if o.CreatedAt.IsZero() {
		o.CreatedAt = time.Now()
	}
	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate is called before updating an existing record
func (o *OAuthCredentialModel) BeforeUpdate(tx *gorm.DB) error {
	o.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate is called before creating a new record
func (a *APIKeyCredentialModel) BeforeCreate(tx *gorm.DB) error {
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	if a.UpdatedAt.IsZero() {
		a.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate is called before updating an existing record
func (a *APIKeyCredentialModel) BeforeUpdate(tx *gorm.DB) error {
	a.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate is called before creating a new record
func (o *OAuthStateModel) BeforeCreate(tx *gorm.DB) error {
	if o.CreatedAt.IsZero() {
		o.CreatedAt = time.Now()
	}
	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate is called before updating an existing record
func (o *OAuthStateModel) BeforeUpdate(tx *gorm.DB) error {
	o.UpdatedAt = time.Now()
	return nil
}
