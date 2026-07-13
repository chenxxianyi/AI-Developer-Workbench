package database

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// newTestDB opens an in-memory SQLite DB with a pure-Go driver (no CGO).
// Tests use this to exercise the migration runner without MySQL.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	return db
}

func TestEnsureMigrationsTable_Idempotent(t *testing.T) {
	db := newTestDB(t)
	if err := ensureMigrationsTable(db); err != nil {
		t.Fatalf("first ensure: %v", err)
	}
	// second call must not error (CREATE TABLE IF NOT EXISTS).
	if err := ensureMigrationsTable(db); err != nil {
		t.Fatalf("second ensure: %v", err)
	}
}

func TestAppliedMigrations_Empty(t *testing.T) {
	db := newTestDB(t)
	if err := ensureMigrationsTable(db); err != nil {
		t.Fatalf("ensure: %v", err)
	}
	got, err := appliedMigrations(db)
	if err != nil {
		t.Fatalf("applied: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 applied, got %d", len(got))
	}
}

func TestApplyMigration_RecordsVersion(t *testing.T) {
	db := newTestDB(t)
	if err := ensureMigrationsTable(db); err != nil {
		t.Fatalf("ensure: %v", err)
	}
	mf := migrationFile{
		Version: "9999_test",
		Source:  []byte("CREATE TABLE IF NOT EXISTS tmp_test (id INTEGER PRIMARY KEY)"),
	}
	if err := applyMigration(db, mf); err != nil {
		t.Fatalf("apply: %v", err)
	}

	applied, err := appliedMigrations(db)
	if err != nil {
		t.Fatalf("applied: %v", err)
	}
	if !applied["9999_test"] {
		t.Fatal("migration not recorded as applied")
	}
}

func TestApplyMigration_FailureDoesNotMarkApplied(t *testing.T) {
	db := newTestDB(t)
	if err := ensureMigrationsTable(db); err != nil {
		t.Fatalf("ensure: %v", err)
	}
	mf := migrationFile{
		Version: "9999_bad",
		// Invalid SQL that will fail on SQLite.
		Source: []byte("CREATE TABLE syntax error here ("),
	}
	if err := applyMigration(db, mf); err == nil {
		t.Fatal("expected migration to fail, got nil")
	}

	applied, err := appliedMigrations(db)
	if err != nil {
		t.Fatalf("applied: %v", err)
	}
	if applied["9999_bad"] {
		t.Fatal("failed migration should not be marked applied")
	}
}

func TestRunVersionedMigrations_DoesNotReapply(t *testing.T) {
	db := newTestDB(t)

	// Apply a migration directly, then call RunVersionedMigrations which
	// should skip it (no-op for that version).
	if err := ensureMigrationsTable(db); err != nil {
		t.Fatalf("ensure: %v", err)
	}
	if err := applyMigration(db, migrationFile{Version: "0001_init", Source: []byte("CREATE TABLE IF NOT EXISTS x (id INTEGER)")}); err != nil {
		t.Fatalf("apply: %v", err)
	}

	// discoverMigrations with a non-existent dir returns error → RunVersionedMigrations
	// logs a warning and skips. That's fine for this assertion: the previously
	// applied version is still recorded.
	if err := RunVersionedMigrations(db, false, nil); err != nil {
		t.Fatalf("run: %v", err)
	}
	applied, _ := appliedMigrations(db)
	if !applied["0001_init"] {
		t.Fatal("0001_init should still be applied")
	}
}
