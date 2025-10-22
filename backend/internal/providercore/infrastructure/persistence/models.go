package persistence

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ProviderModel represents a provider in the database
type ProviderModel struct {
	ID             string          `gorm:"type:uuid;primaryKey"`
	Identifier     string          `gorm:"type:varchar(50);uniqueIndex;not null"`
	Name           string          `gorm:"type:varchar(100);not null"`
	Description    string          `gorm:"type:text"`
	AuthType       string          `gorm:"type:varchar(20);not null"`
	Status         string          `gorm:"type:varchar(20);not null"`
	IconURL        string          `gorm:"type:text"`
	JSONAttributes json.RawMessage `gorm:"type:jsonb;column:json_attributes"`
	Embedding      string          `gorm:"type:vector(1536)"`
	CreatedAt      time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt      time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt      gorm.DeletedAt  `gorm:"type:timestamp with time zone;index"`
}

// TableName specifies the table name for ProviderModel
func (ProviderModel) TableName() string {
	return "providers"
}

// OperationModel represents an operation in the database
type OperationModel struct {
	ID             string          `gorm:"type:uuid;primaryKey"`
	Identifier     string          `gorm:"type:varchar(50);not null"`
	ProviderID     string          `gorm:"type:uuid;not null;index"`
	Name           string          `gorm:"type:varchar(100);not null"`
	Description    string          `gorm:"type:text"`
	Category       string          `gorm:"type:varchar(50);not null"`
	JSONAttributes json.RawMessage `gorm:"type:jsonb;column:json_attributes"`
	Embedding      string          `gorm:"type:vector(1536)"`
	CreatedAt      time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt      time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt      gorm.DeletedAt  `gorm:"type:timestamp with time zone;index"`

	// Relations
	Provider ProviderModel `gorm:"foreignKey:ProviderID;references:ID"`
}

// TableName specifies the table name for OperationModel
func (OperationModel) TableName() string {
	return "operations"
}

// ProviderTranslationModel represents a provider's translation in the database
type ProviderTranslationModel struct {
	ID                 string          `gorm:"type:uuid;primaryKey"`
	ProviderIdentifier string          `gorm:"type:varchar(255);not null;uniqueIndex:idx_provider_lang"`
	LanguageCode       string          `gorm:"type:varchar(10);not null;uniqueIndex:idx_provider_lang"`
	Translations       json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt          time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt          time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt          gorm.DeletedAt  `gorm:"type:timestamp with time zone;index"`
}

// TableName specifies the table name for ProviderTranslationModel
func (ProviderTranslationModel) TableName() string {
	return "provider_translations"
}
