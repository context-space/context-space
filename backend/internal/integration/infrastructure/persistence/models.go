package persistence

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// InvocationModel is the GORM model for operation invocations
type InvocationModel struct {
	ID                  string          `gorm:"type:uuid;primaryKey"`
	UserID              string          `gorm:"type:uuid;not null;index"`
	ProviderIdentifier  string          `gorm:"type:varchar(50);not null;index"`
	OperationIdentifier string          `gorm:"type:varchar(50);not null;index"`
	Status              string          `gorm:"type:varchar(20);not null"`
	Duration            int64           `gorm:"not null"`
	StartedAt           *time.Time      `gorm:"type:timestamp with time zone"`
	CompletedAt         *time.Time      `gorm:"type:timestamp with time zone"`
	JSONAttributes      json.RawMessage `gorm:"type:jsonb"`
	CreatedAt           time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt           time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt           gorm.DeletedAt  `gorm:"type:timestamp with time zone;index"`
}

// TableName overrides the table name
func (InvocationModel) TableName() string {
	return "invocations"
}
