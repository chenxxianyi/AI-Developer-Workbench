package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// schemaMigrationsTable records which migration files have been applied.
const schemaMigrationsTable = "schema_migrations"

// migrationFile represents a single .sql migration discovered on disk.
type migrationFile struct {
	Version string // filename without extension, e.g. "0001_init"
	Path    string // absolute path to the .sql file
	Source  []byte // raw SQL content
}

// RunVersionedMigrations executes SQL migration files in version order, each in
// its own transaction. Already-applied migrations are skipped. A failed
// migration is NOT marked as applied.
//
// Migrations are discovered from the directory `migrationsDir` (default
// "./migrations" relative to the working directory, overridable via the
// MIGRATIONS_DIR environment variable). Only files matching `NNNN_*.sql` are
// considered.
//
// The `autoMigrate` flag controls whether GORM AutoMigrate also runs afterwards
// (useful for dev: keeps the model and schema in sync without hand-writing
// every column change). In production prefer versioned SQL only.
func RunVersionedMigrations(db *gorm.DB, autoMigrate bool, models []any) error {
	// Ensure the tracking table exists.
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to ensure %s table: %w", schemaMigrationsTable, err)
	}

	migrationsDir := getMigrationsDir()
	files, err := discoverMigrations(migrationsDir)
	if err != nil {
		// Missing migrations dir in dev is non-fatal: fall back to AutoMigrate.
		slog.Warn("Versioned migrations directory not readable, skipping SQL migrations", "dir", migrationsDir, "error", err)
	} else {
		applied, err := appliedMigrations(db)
		if err != nil {
			return fmt.Errorf("failed to read applied migrations: %w", err)
		}

		for _, mf := range files {
			if _, ok := applied[mf.Version]; ok {
				continue // already applied
			}
			slog.Info("Applying migration", "version", mf.Version)
			if err := applyMigration(db, mf); err != nil {
				return fmt.Errorf("migration %s failed: %w", mf.Version, err)
			}
			slog.Info("Migration applied", "version", mf.Version)
		}
	}

	if autoMigrate && len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			return fmt.Errorf("auto-migrate failed: %w", err)
		}
	}

	return nil
}

// schemaMigration is a GORM model for the migration tracking table.
// Using a model (instead of raw DDL) keeps it portable across MySQL and SQLite.
type schemaMigration struct {
	Version   string `gorm:"type:varchar(255);primaryKey"`
	AppliedAt string `gorm:"type:varchar(40);not null;default:'1970-01-01T00:00:00Z'"`
}

func (schemaMigration) TableName() string { return schemaMigrationsTable }

func ensureMigrationsTable(db *gorm.DB) error {
	return db.AutoMigrate(&schemaMigration{})
}

func appliedMigrations(db *gorm.DB) (map[string]bool, error) {
	type row struct {
		Version string
	}
	var rows []row
	if err := db.Table(schemaMigrationsTable).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]bool, len(rows))
	for _, r := range rows {
		out[r.Version] = true
	}
	return out, nil
}

func applyMigration(db *gorm.DB, mf migrationFile) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(string(mf.Source)).Error; err != nil {
			return err
		}
		return tx.Create(&schemaMigration{Version: mf.Version}).Error
	})
}

func discoverMigrations(dir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []migrationFile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		name := e.Name()
		version := strings.TrimSuffix(name, ".sql")
		// require NNNN_ prefix (4-digit version)
		if len(version) < 5 || version[4] != '_' {
			continue
		}
		full := filepath.Join(dir, name)
		src, err := os.ReadFile(full)
		if err != nil {
			return nil, fmt.Errorf("read migration %s: %w", name, err)
		}
		files = append(files, migrationFile{Version: version, Path: full, Source: src})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Version < files[j].Version })
	return files, nil
}

func getMigrationsDir() string {
	if d := os.Getenv("MIGRATIONS_DIR"); d != "" {
		return d
	}
	return "./migrations"
}

// PingMigrations is a tiny helper for tests: it returns the count of applied
// migrations. Useful to assert a fresh DB reached the expected version.
func PingMigrations(ctx context.Context, db *gorm.DB) (int, error) {
	var n int64
	if err := db.WithContext(ctx).Table(schemaMigrationsTable).Count(&n).Error; err != nil {
		return 0, err
	}
	return int(n), nil
}
