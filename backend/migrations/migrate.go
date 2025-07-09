package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/context-space/context-space/backend/internal/shared/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Parse command-line flags
	upFlag := flag.Bool("up", false, "Run migrations up")
	downFlag := flag.Bool("down", false, "Run migrations down")
	versionFlag := flag.Bool("version", false, "Print current migration version")
	stepsFlag := flag.Int("steps", 0, "Number of migrations to apply (up or down)")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Get migration source path
	migrationsPath := "file://migrations/postgresql"
	if path := os.Getenv("MIGRATIONS_PATH"); path != "" {
		migrationsPath = fmt.Sprintf("file://%s", path)
	}

	// Get database connection string
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.MigrationUsername,
		cfg.Database.MigrationPassword,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
		cfg.Database.SSLMode,
	)

	// Create migration instance
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	defer m.Close()

	// Execute the requested command
	if *versionFlag {
		version, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				log.Println("No migrations have been applied yet")
			} else {
				log.Fatalf("Failed to get migration version: %v", err)
			}
		} else {
			log.Printf("Current migration version: %d (dirty: %v)", version, dirty)
		}
		return
	}

	if *upFlag {
		if *stepsFlag > 0 {
			if err := m.Steps(*stepsFlag); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to apply migrations up %d steps: %v", *stepsFlag, err)
			}
			log.Printf("Applied %d migrations up", *stepsFlag)
		} else {
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to apply migrations up: %v", err)
			}
			log.Println("Applied all migrations up")
		}
		return
	}

	if *downFlag {
		if *stepsFlag > 0 {
			if err := m.Steps(-*stepsFlag); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to apply migrations down %d steps: %v", *stepsFlag, err)
			}
			log.Printf("Applied %d migrations down", *stepsFlag)
		} else {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to apply migrations down: %v", err)
			}
			log.Println("Applied all migrations down")
		}
		return
	}

	// If no command specified, print usage
	flag.Usage()
}
