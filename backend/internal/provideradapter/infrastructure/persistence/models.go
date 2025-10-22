package persistence

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ProviderModel represents a provider in the database
type ProviderAdapterModel struct {
	ID          string          `gorm:"type:uuid;primaryKey"`
	Identifier  string          `gorm:"type:varchar(50);uniqueIndex;not null"`
	Configs     json.RawMessage `gorm:"type:jsonb;column:configs"`
	Permissions json.RawMessage `gorm:"type:jsonb;column:permissions"`
	CreatedAt   time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt   time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt   gorm.DeletedAt  `gorm:"type:timestamp with time zone;index"`
}

// TableName specifies the table name for ProviderModel
func (ProviderAdapterModel) TableName() string {
	return "provider_adapters"
}
