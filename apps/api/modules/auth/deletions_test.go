package auth

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/FacileStudio/Capsule/apps/api/internal/env"
	"github.com/FacileStudio/Capsule/apps/api/internal/testsupport"
	portepg "github.com/FacileStudio/porte/pg"
)

const subject = "sub-123"

func testSQL(t *testing.T) *sql.DB {
	t.Helper()
	sqlDB, err := testsupport.DB(t).DB()
	if err != nil {
		t.Fatalf("sql handle: %v", err)
	}
	return sqlDB
}

func fakeRegistre(t *testing.T, polls *atomic.Int32) string {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/directory/deletions" {
			http.NotFound(w, r)
			return
		}
		polls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		body := map[string]any{"deleted": []map[string]any{{"sub": subject, "deleted_at": time.Now().UTC().Add(-time.Minute)}}}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Errorf("fake Registre: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	return server.URL
}

func seedAccount(t *testing.T, sqlDB *sql.DB, issuer, email, sub string) int64 {
	t.Helper()
	ctx := context.Background()
	var userID int64
	err := sqlDB.QueryRowContext(ctx,
		`INSERT INTO porte_users (email, email_verified) VALUES ($1, true) RETURNING id`, email).Scan(&userID)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	for _, stmt := range []string{
		`INSERT INTO porte_identities (user_id, provider, subject) VALUES ($1, $2, $3)`,
		`INSERT INTO porte_sessions (token_hash, user_id) VALUES ($3::text || $2::text, $1)`,
		`INSERT INTO porte_login_codes (code_hash, user_id, expires_at, provider, subject) VALUES ($3 || '-c', $1, now() + interval '1 minute', $2, $3)`,
		`INSERT INTO porte_magic_links (token_hash, user_id, created_at, expires_at) VALUES ($3::text || $2::text, $1, now(), now() + interval '1 minute')`,
	} {
		if _, err := sqlDB.ExecContext(ctx, stmt, userID, issuer, sub); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	return userID
}

var userTables = []string{"porte_identities", "porte_sessions", "porte_login_codes", "porte_magic_links"}

func rowsFor(t *testing.T, sqlDB *sql.DB, userID int64) map[string]int {
	t.Helper()
	counts := map[string]int{}
	for _, table := range append([]string{"porte_users"}, userTables...) {
		column := "user_id"
		if table == "porte_users" {
			column = "id"
		}
		var rows int
		if err := sqlDB.QueryRow("SELECT count(*) FROM "+table+" WHERE "+column+" = $1", userID).Scan(&rows); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		counts[table] = rows
	}
	return counts
}

func pollTwice(t *testing.T, sqlDB *sql.DB, issuer string) {
	t.Helper()
	oidc := &env.OIDCConfig{Issuer: issuer, ClientID: "capsule", ClientSecret: "secret"}
	deletions := newDeletions(sqlDB, portepg.New(sqlDB), newUsers(sqlDB), oidc, slog.New(slog.DiscardHandler))
	for attempt := range 2 {
		if err := deletions.Poll(context.Background()); err != nil {
			t.Fatalf("poll %d: %v", attempt, err)
		}
	}
}

func TestRegistreErasureRemovesEveryRowOfThePersonAndOnlyThem(t *testing.T) {
	sqlDB := testSQL(t)
	var polls atomic.Int32
	issuer := fakeRegistre(t, &polls)
	erased := seedAccount(t, sqlDB, issuer, "gone@example.test", subject)
	kept := seedAccount(t, sqlDB, issuer, "kept@example.test", "sub-other")
	if _, err := sqlDB.Exec(`INSERT INTO pastes (id, content) VALUES ('cap_keep', 'cipher')`); err != nil {
		t.Fatalf("seed paste: %v", err)
	}

	pollTwice(t, sqlDB, issuer)

	if polls.Load() < 2 {
		t.Fatalf("Registre was polled %d times, want 2", polls.Load())
	}
	for table, rows := range rowsFor(t, sqlDB, erased) {
		if rows != 0 {
			t.Errorf("%s keeps %d rows of the erased person", table, rows)
		}
	}
	for table, rows := range rowsFor(t, sqlDB, kept) {
		if rows != 1 {
			t.Errorf("%s has %d rows of the other person, want 1", table, rows)
		}
	}
	var pastes int
	if err := sqlDB.QueryRow(`SELECT count(*) FROM pastes`).Scan(&pastes); err != nil || pastes != 1 {
		t.Errorf("pastes = %d, %v; want the anonymous paste kept", pastes, err)
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
