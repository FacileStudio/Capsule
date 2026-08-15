// Package testsupport gives a test package a migrated PostgreSQL database.
//
// Capsule's schema is the hand-written DDL in migrations/, and its service
// layer runs SELECT ... FOR UPDATE, so the in-memory SQLite these tests used to
// open could neither build the shipped schema nor execute the shipped queries.
package testsupport

import (
	"context"
	"hash/fnv"
	"io"
	"log/slog"
	"sync"
	"testing"

	"github.com/FacileStudio/Capsule/apps/api/migrations"

	"github.com/FacileStudio/tronc/migrate"
	"github.com/FacileStudio/tronc/testdb"
	"gorm.io/gorm"
)

const (
	prefix    = "capsule_test"
	bootstrap = "docker compose up -d capsule-db   # or any PostgreSQL 16"
)

var (
	once     sync.Once
	shared   *gorm.DB
	openErr  error
	settings = testdb.Config{Prefix: prefix, Migrate: runMigrations}
)

// DB returns this package's database, migrated once and emptied per test. It
// skips the test when TEST_DATABASE_URL is unset.
func DB(t *testing.T) *gorm.DB {
	t.Helper()

	url, configured := testdb.URL()
	if !configured {
		testdb.Announce(bootstrap)
		t.Skip(testdb.SkipReason(bootstrap))
	}

	once.Do(func() { shared, openErr = testdb.Open(url, settings) })
	if openErr != nil {
		t.Fatalf("test database: %v", openErr)
	}
	if err := testdb.Truncate(shared, settings); err != nil {
		t.Fatalf("reset test database: %v", err)
	}
	return shared
}

// runMigrations applies schema changes to this package's private schema. The
// advisory lock is scoped to the database, not the schema, so a fixed lock id
// would otherwise make every test package queue on the same default id;
// lockID derives a per-schema id instead.
func runMigrations(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return migrate.Run(context.Background(), migrate.Config{
		DB:     sqlDB,
		FS:     migrations.FS,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		LockID: lockID(),
	})
}

func lockID() int64 {
	digest := fnv.New64a()
	_, _ = digest.Write([]byte(testdb.SchemaName(prefix)))
	return int64(digest.Sum64() >> 1)
}
