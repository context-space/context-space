package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/shared/config"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Postgres wraps a GORM DB connection
type Postgres struct {
	DB              *gorm.DB
	obs             *observability.ObservabilityProvider
	traceOperations bool
}

// NewPostgresClient creates a new Postgres client
func NewPostgresClient(cfg *config.Config, obs *observability.ObservabilityProvider) (*Postgres, error) {
	// Configure GORM logger based on the app's log level
	logLevel := logger.Silent
	if cfg.Logging.Level == "debug" {
		logLevel = logger.Info
	}

	logFile, err := os.OpenFile("context-space-backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	customLogger := logger.New(
		log.New(logFile, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(cfg.GetDatabaseDSN()), &gorm.Config{
		Logger: customLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}

	// Set sane connection pool limits
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Postgres{
		DB:              db,
		obs:             obs,
		traceOperations: cfg.Logging.Level == "debug",
	}, nil
}

// WithContext returns a new DB instance with context
func (c *Postgres) WithContext(ctx context.Context) *gorm.DB {
	return c.DB.WithContext(ctx)
}

// Transaction executes a function within a database transaction
func (c *Postgres) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	var span trace.Span
	if c.traceOperations {
		ctx, span = c.obs.Tracer.Start(ctx, "database.Transaction")
		defer span.End()
	}

	return c.DB.WithContext(ctx).Transaction(fn)
}

// Close closes the database connection
func (c *Postgres) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %w", err)
	}
	return sqlDB.Close()
}

// Ping checks if the database connection is alive
func (c *Postgres) Ping() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %w", err)
	}
	return sqlDB.Ping()
}

func (c *Postgres) Where(query interface{}, args ...interface{}) *gorm.DB {
	return c.DB.Where(query, args...)
}

func (c *Postgres) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return c.DB.Find(dest, conds...)
}

func (c *Postgres) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return c.DB.First(dest, conds...)
}

func (c *Postgres) Offset(offset int) *gorm.DB {
	return c.DB.Offset(offset)
}

func (c *Postgres) Count(count *int64) *gorm.DB {
	return c.DB.Count(count)
}

func (c *Postgres) Order(value interface{}) *gorm.DB {
	return c.DB.Order(value)
}

func (c *Postgres) Model(value interface{}) *gorm.DB {
	return c.DB.Model(value)
}

func (c *Postgres) Create(value interface{}) *gorm.DB {
	return c.DB.Create(value)
}

func (c *Postgres) Save(value interface{}) *gorm.DB {
	return c.DB.Save(value)
}

func (c *Postgres) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return c.DB.Delete(value, conds...)
}

func (c *Postgres) Begin() *gorm.DB {
	return c.DB.Begin()
}

func (c *Postgres) Commit() *gorm.DB {
	return c.DB.Commit()
}

func (c *Postgres) Rollback() *gorm.DB {
	return c.DB.Rollback()
}
