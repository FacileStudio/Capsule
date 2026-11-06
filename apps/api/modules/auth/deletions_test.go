package auth

import (
	"bytes"
	"context"
	"database/sql"
	"testing"

	"github.com/FacileStudio/Capsule/apps/api/internal/testsupport"
	portepg "github.com/FacileStudio/porte/pg"
)

const (
	issuer  = "https://sso.example.test"
	subject = "sub-123"
)

func testSQL(t *testing.T) *sql.DB {
	t.Helper()
	sqlDB, err := testsupport.DB(t).DB()
	if err != nil {
		t.Fatalf("sql handle: %v", err)
	}
	return sqlDB
}

func seedAccount(t *testing.T, sqlDB *sql.DB) {
	t.Helper()
	ctx := context.Background()
	var userID int64
	err := sqlDB.QueryRowContext(ctx,
		`INSERT INTO porte_users (email, email_verified) VALUES ('a@example.test', true) RETURNING id`).Scan(&userID)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	for _, stmt := range []string{
		`INSERT INTO porte_identities (user_id, provider, subject) VALUES ($1, '` + issuer + `', '` + subject + `')`,
		`INSERT INTO porte_sessions (token_hash, user_id) VALUES ('h', $1)`,
		`INSERT INTO porte_login_codes (code_hash, user_id, expires_at) VALUES ('c', $1, now() + interval '1 minute')`,
	} {
		if _, err := sqlDB.ExecContext(ctx, stmt, userID); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	if _, err := sqlDB.ExecContext(ctx, `INSERT INTO pastes (id, content) VALUES ('cap_keep', 'cipher')`); err != nil {
		t.Fatalf("seed paste: %v", err)
	}
}

func rowCount(t *testing.T, sqlDB *sql.DB, table string) int {
	t.Helper()
	var rows int
	if err := sqlDB.QueryRow("SELECT count(*) FROM " + table).Scan(&rows); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return rows
}

func TestErasedPersonLosesTheirAccountAndASecondCallIsANoOp(t *testing.T) {
	sqlDB := testSQL(t)
	seedAccount(t, sqlDB)
	ctx := context.Background()

	erase := eraseAccount(sqlDB)
	store := portepg.New(sqlDB)
	for attempt := range 2 {
		if err := erase(ctx, issuer, subject); err != nil {
			t.Fatalf("erase attempt %d: %v", attempt, err)
		}
		if err := store.EraseIdentity(ctx, issuer, subject); err != nil {
			t.Fatalf("EraseIdentity attempt %d: %v", attempt, err)
		}
	}

	for _, table := range []string{"porte_users", "porte_identities", "porte_sessions", "porte_login_codes"} {
		if rows := rowCount(t, sqlDB, table); rows != 0 {
			t.Errorf("%s has %d rows, want 0", table, rows)
		}
	}
	if rows := rowCount(t, sqlDB, "pastes"); rows != 1 {
		t.Errorf("pastes has %d rows, want the anonymous paste kept", rows)
	}
}

func TestDeletionStateRoundTrips(t *testing.T) {
	state := deletionState{db: testSQL(t)}
	ctx := context.Background()

	if got, err := state.Load(ctx); err != nil || got != nil {
		t.Fatalf("Load on an empty table = %q, %v; want nil, nil", got, err)
	}
	for _, want := range [][]byte{[]byte(`{"v":1}`), []byte(`{"v":2}`)} {
		if err := state.Save(ctx, want); err != nil {
			t.Fatalf("Save: %v", err)
		}
		if got, err := state.Load(ctx); err != nil || !bytes.Equal(got, want) {
			t.Fatalf("Load = %q, %v; want %q", got, err, want)
		}
	}
}
