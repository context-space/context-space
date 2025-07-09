package database

import (
	"context"

	"gorm.io/gorm"
)

type UnitOfWork interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	GetTx() *gorm.DB
}

type UnitOfWorkFactory interface {
	Create() UnitOfWork
}
