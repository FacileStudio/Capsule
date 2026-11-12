package auth

import (
	"context"
	"database/sql"
	stderrors "errors"
	"log/slog"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/FacileStudio/Capsule/apps/api/internal/env"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/oidc"
	"github.com/FacileStudio/porte/oidc/jwt"
	portepg "github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/porte/session"
	"github.com/FacileStudio/tronc/mcp/mcptest"
)

type roundTrip func(*http.Request) (*http.Response, error)

func (f roundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func rewriteTo(target string) *http.Client {
	host, _ := url.Parse(target)
	return &http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
		r = r.Clone(r.Context())
		r.URL.Scheme, r.URL.Host, r.Host = host.Scheme, host.Host, ""
		return http.DefaultTransport.RoundTrip(r)
	})}
}

// erasedFixture is the receiver and the OIDC kit wired as New wires them. The
// issuer is a test IdP, and the receiver's polls go to a Registre that lists
// one erasure, subject, every time.
type erasedFixture struct {
	sqlDB     *sql.DB
	issuer    string
	deletions *directory.Deletions
	kit       *oidc.Kit
}

func newErasedFixture(t *testing.T) erasedFixture {
	t.Helper()
	sqlDB := testSQL(t)
	var polls atomic.Int32
	registre := fakeRegistre(t, &polls)
	idp := mcptest.NewIssuer(t)
	logger := slog.New(slog.DiscardHandler)
	store := portepg.New(sqlDB)
	deletions := newDeletions(sqlDB, store, newUsers(sqlDB),
		&env.OIDCConfig{Issuer: idp.URL, ClientID: "capsule", ClientSecret: "secret"}, logger)
	deletions.Client.HTTP = rewriteTo(registre)
	cfg := porte.Config{Issuer: idp.URL, ClientID: "capsule", ClientSecret: "secret",
		RedirectURL: "http://capsule.test/api/auth/oidc/callback", SuccessURL: "http://capsule.test"}
	sessions, err := session.New(cfg, session.Deps{Sessions: store.Sessions(), Logger: logger})
	if err != nil {
		t.Fatalf("session.New: %v", err)
	}
	kit, err := oidc.New(context.Background(), cfg, oidc.Deps{
		Users: store.Users(), Identities: store.Identities(), Sessions: sessions,
		Codes: store.LoginCodes(), Logger: logger, Erased: deletions,
	})
	if err != nil {
		t.Fatalf("oidc.New: %v", err)
	}
	return erasedFixture{sqlDB: sqlDB, issuer: idp.URL, deletions: deletions, kit: kit}
}

func (f erasedFixture) poll(t *testing.T) {
	t.Helper()
	if err := f.deletions.Poll(context.Background()); err != nil {
		t.Fatalf("poll: %v", err)
	}
}

func (f erasedFixture) provision(sub, email string) (porte.Identity, error) {
	verified := true
	return f.kit.ProvisionBearer(context.Background(),
		jwt.Claims{Subject: sub, Email: email, Name: sub, Profile: jwt.Profile{EmailVerified: &verified}})
}

func (f erasedFixture) accounts(t *testing.T, email string) int {
	t.Helper()
	var n int
	if err := f.sqlDB.QueryRow(`SELECT count(*) FROM porte_users WHERE email = $1`, email).Scan(&n); err != nil {
		t.Fatalf("count accounts: %v", err)
	}
	return n
}

// An account re-created for a subject Registre already erased, here written
// past the guard as a bearer minted before the erasure could, is erased again
// by the next poll of the same receiver, with no restart, and nobody else is
// touched.
func TestARecreatedAccountForAnErasedSubjectIsErasedAgain(t *testing.T) {
	f := newErasedFixture(t)
	seedAccount(t, f.sqlDB, f.issuer, "gone@example.test", subject)
	kept := seedAccount(t, f.sqlDB, f.issuer, "kept@example.test", "sub-other")
	f.poll(t)
	if f.accounts(t, "gone@example.test") != 0 {
		t.Fatalf("the first poll did not erase %s", subject)
	}

	again := seedAccount(t, f.sqlDB, f.issuer, "gone@example.test", subject)
	f.poll(t)
	for table, rows := range rowsFor(t, f.sqlDB, again) {
		if rows != 0 {
			t.Errorf("%s keeps %d rows of the re-created account", table, rows)
		}
	}
	for table, rows := range rowsFor(t, f.sqlDB, kept) {
		if rows != 1 {
			t.Errorf("%s has %d rows of the other person, want 1", table, rows)
		}
	}
}

// Once the receiver handled an erasure, the OIDC kit refuses that subject an
// account on the MCP first-call path, and still provisions everyone else.
func TestProvisioningAnErasedSubjectIsRefused(t *testing.T) {
	f := newErasedFixture(t)
	seedAccount(t, f.sqlDB, f.issuer, "gone@example.test", subject)
	f.poll(t)

	_, err := f.provision(subject, "gone@example.test")
	if !stderrors.Is(err, porte.ErrSubjectErased) || !stderrors.Is(err, oidc.ErrProvisionRefused) {
		t.Fatalf("ProvisionBearer(%s) = %v, want a refusal wrapping ErrSubjectErased", subject, err)
	}
	if f.accounts(t, "gone@example.test") != 0 {
		t.Fatalf("a refused provisioning still wrote an account")
	}
	if identity, err := f.provision("sub-new", "new@example.test"); err != nil || identity.UserID == 0 {
		t.Fatalf("ProvisionBearer(sub-new) = %+v, %v; want a new account", identity, err)
	}
}
