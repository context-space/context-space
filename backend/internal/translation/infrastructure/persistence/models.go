package persistence

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// TranslationModel 翻译数据库模型
type TranslationModel struct {
	ID                 string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProviderIdentifier string          `gorm:"type:varchar(255);not null" json:"provider_identifier"`
	LanguageCode       string          `gorm:"type:varchar(10);not null" json:"language_code"`
	Translations       json.RawMessage `gorm:"type:jsonb" json:"translations"`
	CreatedAt          time.Time       `gorm:"type:timestamp with time zone;not null;default:now();index" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"type:timestamp with time zone;not null;default:now();index" json:"updated_at"`
	DeletedAt          gorm.DeletedAt  `gorm:"type:timestamp with time zone;index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (TranslationModel) TableName() string {
	return "provider_translations"
}
