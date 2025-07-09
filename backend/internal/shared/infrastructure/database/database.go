package database

import (
	"context"

	"gorm.io/gorm"
)

// Database is the interface for the Database
type Database interface {
	WithContext(ctx context.Context) *gorm.DB
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	Close() error
	Ping() error
}