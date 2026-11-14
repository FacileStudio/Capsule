package mcp_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/FacileStudio/Capsule/apps/api/internal/testsupport"
	"github.com/FacileStudio/Capsule/apps/api/modules/mcp"
	"github.com/FacileStudio/Capsule/apps/api/modules/pastes"
	"github.com/FacileStudio/Capsule/apps/api/schemas"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/mcpauth"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	"github.com/FacileStudio/tronc/mcp/mcptest"
	"github.com/go-chi/chi/v5"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"gorm.io/gorm"
)

const ciphertext = "c2VhbGVkLWNpcGhlcnRleHQ"

type identities struct {
	porte.IdentityStore
	users map[string]int64
}

func (s identities) Find(_ context.Context, provider, sub string) (porte.StoredIdentity, error) {
	id, ok := s.users[sub]
	if !ok {
		return porte.StoredIdentity{}, porte.ErrNotFound
	}
	return porte.StoredIdentity{UserID: id, Provider: provider, Subject: sub}, nil
}

type fixture struct {
	db       *gorm.DB
	sqlDB    *sql.DB
	pastes   *pastes.Service
	issuer   *mcptest.Issuer
	resource string
	server   *httptest.Server
	users    map[string]int64
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	db := testsupport.DB(t)
	sqlDB, err := db.DB()
	must(t, err)
	f := &fixture{
		db: db, sqlDB: sqlDB, pastes: pastes.NewService(db, 1<<20), issuer: mcptest.NewIssuer(t),
		users: map[string]int64{"owner": 1, "other": 2},
	}
	var handler http.Handler
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handler.ServeHTTP(w, r) }))
	t.Cleanup(f.server.Close)
	f.resource = f.server.URL + "/mcp"
	service := f.service(t)
	router := chi.NewRouter()
	mcp.RegisterRoutes(router, service)
	router.Route("/api", func(r chi.Router) { mcp.RegisterPendingRoutes(r, service, signedIn) })
	handler = router
	return f
}

func (f *fixture) service(t *testing.T) *mcp.Service {
	t.Helper()
	verifier, err := mcpauth.Verifier(context.Background(), mcpauth.Config{
		Issuer: f.issuer.URL, Resource: f.resource, ScopePrefix: "capsule", Identities: identities{users: f.users},
	})
	must(t, err)
	service, err := mcp.NewService(mcp.Deps{Pastes: f.pastes}, mcp.Config{
		Resource: f.resource, Issuer: f.issuer.URL, Version: "test", StateKey: []byte("state-key"),
		Verifier: verifier, Pending: mcpkit.PostgresPending{DB: f.sqlDB}, PerMinute: 100_000,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	must(t, err)
	return service
}

func (f *fixture) token(t testing.TB, scopes ...string) string {
	if len(scopes) == 0 {
		scopes = []string{"capsule.read", "capsule.write"}
	}
	return f.issuer.Token(t, f.resource, "owner", scopes...)
}

func (f *fixture) connect(t *testing.T, subject string) *sdk.ClientSession {
	t.Helper()
	return mcptest.Connect(t, f.resource, f.issuer.Token(t, f.resource, subject, "capsule.read", "capsule.write"), mcptest.ClientOptions{})
}

func (f *fixture) seal(t testing.TB, burnAfterRead bool) *pastes.CreateResponse {
	t.Helper()
	created, err := f.pastes.Create(pastes.CreateRequest{Content: ciphertext, BurnAfterRead: &burnAfterRead, ExpiresIn: "24h"})
	must(t, err)
	return created
}

func (f *fixture) burned(id string) bool {
	var paste schemas.Paste
	f.db.Where("id = ?", id).First(&paste)
	return paste.Burned && paste.Content == ""
}

func (f *fixture) page(t *testing.T, method, path, user string, header map[string]string) (int, string) {
	t.Helper()
	req, err := http.NewRequest(method, f.server.URL+"/api/mcp/pending"+path, nil)
	must(t, err)
	req.Header.Set("X-Test-User", strconv.FormatInt(f.users[user], 10))
	for key, value := range header {
		req.Header.Set(key, value)
	}
	resp, err := http.DefaultClient.Do(req)
	must(t, err)
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	must(t, err)
	return resp.StatusCode, string(body)
}

func signedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.Header.Get("X-Test-User"), 10, 64)
		if err != nil {
			http.Error(w, "unauthenticated", http.StatusUnauthorized)
			return
		}
		identity := porte.Identity{UserID: id, SessionID: 1, BrowserSession: r.Header.Get("X-Test-Session") != "cli"}
		next.ServeHTTP(w, r.WithContext(porte.WithIdentity(r.Context(), identity)))
	})
}

func call(t *testing.T, session *sdk.ClientSession, tool string, args map[string]any, out any) string {
	t.Helper()
	result := mcptest.Call(t, session, tool, args)
	if result.IsError {
		t.Fatalf("%s: %s", tool, mcptest.Text(result))
	}
	mcptest.Decode(t, result, out)
	return mcptest.Text(result)
}

func refused(t *testing.T, session *sdk.ClientSession, tool string, args map[string]any) string {
	t.Helper()
	result := mcptest.Call(t, session, tool, args)
	if !result.IsError {
		t.Fatalf("%s must be refused, got %s", tool, mcptest.Text(result))
	}
	return mcptest.Text(result)
}

func actionID(confirmURL string) string {
	return confirmURL[strings.LastIndex(confirmURL, "=")+1:]
}

func must(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
