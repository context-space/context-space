package persistence

import (
	"time"

	"gorm.io/gorm"
)

func parseGormDeletedAt(deletedAt gorm.DeletedAt) *time.Time {
	if deletedAt.Valid {
		return &deletedAt.Time
	}
	return nil
}

func parseDomainDeletedAt(deletedAt *time.Time) gorm.DeletedAt {
	if deletedAt == nil {
		return gorm.DeletedAt{}
	}
	return gorm.DeletedAt{Time: *deletedAt, Valid: true}
}
