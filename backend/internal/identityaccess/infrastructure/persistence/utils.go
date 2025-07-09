package persistence

import (
	"time"

	"github.com/bytedance/sonic"
	"gorm.io/gorm"
)

func mustMarshalJSON(v interface{}) []byte {
	jsonBytes, err := sonic.Marshal(v)
	if err != nil {
		panic(err)
	}
	return jsonBytes
}

func mustUnmarshalJSON(data []byte) map[string]interface{} {
	var v map[string]interface{}
	if err := sonic.Unmarshal(data, &v); err != nil {
		panic(err)
	}
	return v
}

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

func parseGormEmail(email *string) string {
	var result string
	if email != nil {
		result = *email
	}
	return result
}

func parseDomainEmail(email string) *string {
	var result *string
	if email != "" {
		result = &email
	}
	return result
}
