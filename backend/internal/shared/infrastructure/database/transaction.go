package database

import (
	"context"
	"fmt"

	observability "github.com/context-space/cloud-observability"
	"gorm.io/gorm"
)

// DefaultUnitOfWork represents a database transaction
type DefaultUnitOfWork struct {
	db  Database
	tx  *gorm.DB
	obs *observability.ObservabilityProvider
}

// NewUnitOfWork creates a new unit of work
func NewDefaultUnitOfWork(db Database, observabilityProvider *observability.ObservabilityProvider) *DefaultUnitOfWork {
	return &DefaultUnitOfWork{
		db:  db,
		obs: observabilityProvider,
	}
}

// Begin begins a transaction
func (u *DefaultUnitOfWork) Begin(ctx context.Context) error {
	ctx, span := u.obs.Tracer.Start(ctx, "UnitOfWork.Begin")
	defer span.End()

	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	u.tx = tx
	return nil
}

// Commit commits a transaction
func (u *DefaultUnitOfWork) Commit(ctx context.Context) error {
	_, span := u.obs.Tracer.Start(ctx, "UnitOfWork.Commit")
	defer span.End()

	if u.tx == nil {
		return fmt.Errorf("no active transaction to commit")
	}

	if err := u.tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	u.tx = nil
	return nil
}

// Rollback rolls back a transaction
func (u *DefaultUnitOfWork) Rollback(ctx context.Context) error {
	_, span := u.obs.Tracer.Start(ctx, "UnitOfWork.Rollback")
	defer span.End()

	if u.tx == nil {
		return fmt.Errorf("no active transaction to rollback")
	}

	if err := u.tx.Rollback().Error; err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	u.tx = nil
	return nil
}

// GetTx gets the transaction
func (u *DefaultUnitOfWork) GetTx() *gorm.DB {
	return u.tx
}

// UnitOfWorkFactory is responsible for creating UnitOfWork instances
type DefaultUnitOfWorkFactory struct {
	db  Database
	obs *observability.ObservabilityProvider
}

// NewUnitOfWorkFactory creates a new UnitOfWorkFactory
func NewDefaultUnitOfWorkFactory(db Database, observabilityProvider *observability.ObservabilityProvider) *DefaultUnitOfWorkFactory {
	return &DefaultUnitOfWorkFactory{
		db:  db,
		obs: observabilityProvider,
	}
}

// Create returns a new UnitOfWork instance
func (f *DefaultUnitOfWorkFactory) Create() UnitOfWork {
	return NewDefaultUnitOfWork(f.db, f.obs)
}
