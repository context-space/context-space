package persistence

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

// UserModel represents the users table in the database
type UserModel struct {
	ID          string         `gorm:"type:uuid;primaryKey"`
	SupID       string         `gorm:"type:uuid;not null;uniqueIndex"`
	Email       *string        `gorm:"type:varchar(255);uniqueIndex"`
	IsAnonymous bool           `gorm:"type:boolean;not null;default:false"`
	CreatedAt   time.Time      `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt   time.Time      `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamp with time zone;index"`
}

// TableName returns the table name for the User model
func (UserModel) TableName() string {
	return "users"
}

// UserInfoModel represents the user_infos table in the database
type UserInfoModel struct {
	ID           string          `gorm:"type:uuid;primaryKey"`
	UserID       string          `gorm:"type:uuid;not null;index"`
	InfoMetadata json.RawMessage `gorm:"type:jsonb"`
	CreatedAt    time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt    time.Time       `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt    gorm.DeletedAt  `gorm:"type:timestamp with time zone;index"`

	// Relationship
	User UserModel `gorm:"foreignKey:UserID;references:ID"`
}

// TableName returns the table name for the UserInfo model
func (UserInfoModel) TableName() string {
	return "user_infos"
}

// UserAPIKeyModel represents the api_keys table in the database
type UserAPIKeyModel struct {
	ID          string         `gorm:"type:uuid;primaryKey"`
	UserID      string         `gorm:"type:uuid;not null;index"`
	KeyValue    string         `gorm:"type:varchar(64);not null;uniqueIndex;column:key_value"`
	Name        string         `gorm:"type:varchar(100)"`
	Description string         `gorm:"type:text"`
	LastUsed    *time.Time     `gorm:"type:timestamp with time zone"`
	CreatedAt   time.Time      `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt   time.Time      `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamp with time zone;index"`

	// Relationships
	User UserModel `gorm:"foreignKey:UserID;references:ID"`
}

// TableName returns the table name for the APIKey model
func (UserAPIKeyModel) TableName() string {
	return "user_api_keys"
}

// BeforeCreate is called before creating a new record
func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate is called before updating an existing record
func (u *UserModel) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate is called before creating a new record
func (u *UserInfoModel) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate is called before updating an existing record
func (u *UserInfoModel) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate is called before creating a new record
func (k *UserAPIKeyModel) BeforeCreate(tx *gorm.DB) error {
	if k.ID == "" {
		k.ID = uuid.New().String()
	}
	if k.CreatedAt.IsZero() {
		k.CreatedAt = time.Now()
	}
	if k.UpdatedAt.IsZero() {
		k.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate is called before updating an existing record
func (k *UserAPIKeyModel) BeforeUpdate(tx *gorm.DB) error {
	k.UpdatedAt = time.Now()
	return nil
}
