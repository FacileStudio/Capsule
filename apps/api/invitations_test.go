package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/Capsule/apps/api/internal/env"
	"github.com/FacileStudio/Capsule/apps/api/internal/testsupport"
	"github.com/FacileStudio/Capsule/apps/api/modules/auth"
	"github.com/go-chi/chi/v5"
)

func fakeRegistre(t *testing.T) string {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		base := "http://" + r.Host
		body := map[string]any{
			"issuer": base, "authorization_endpoint": base + "/auth", "token_endpoint": base + "/token", "jwks_uri": base + "/keys",
		}
		if r.URL.Path == "/api/directory/invitations" {
			body = map[string]any{"count": 1, "invitations": []map[string]any{{"id": "inv-1", "url": base + "/invitations/inv-1"}}}
		}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Errorf("fake Registre: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	return server.URL
}

func ssoRouter(t *testing.T, issuer string) (chi.Router, *auth.Kits) {
	t.Helper()
	db := testsupport.DB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql handle: %v", err)
	}
	appEnv := env.Config{MaxPasteSize: 1024, OIDC: &env.OIDCConfig{
		Issuer: issuer, ClientID: "capsule", ClientSecret: "secret",
		RedirectURL: "http://localhost:4000/api/auth/oidc/callback", SuccessURL: "http://localhost:4000/",
	}}
	kits, err := auth.New(context.Background(), sqlDB, appEnv, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatalf("auth.New: %v", err)
	}
	router, err := buildRouter(db, sqlDB, kits, appEnv, slog.New(slog.DiscardHandler), nil)
	if err != nil {
		t.Fatalf("buildRouter: %v", err)
	}
	return router, kits
}

func signIn(t *testing.T, kits *auth.Kits, issuer string) string {
	t.Helper()
	sqlDB := kits.Users.DB
	var userID int64
	if err := sqlDB.QueryRow(`INSERT INTO porte_users (email, email_verified) VALUES ('inv@example.test', true) RETURNING id`).Scan(&userID); err != nil {
		t.Fatalf("insert user: %v", err)
	}
	if _, err := sqlDB.Exec(`INSERT INTO porte_identities (user_id, provider, subject) VALUES ($1, $2, 'sub-inv')`, userID, issuer); err != nil {
		t.Fatalf("insert identity: %v", err)
	}
	token, _, err := kits.Sessions.Issue(context.Background(), userID, "test")
	if err != nil {
		t.Fatalf("issue session: %v", err)
	}
	return token
}

func TestInvitationsAreMountedWithSingleSignOnBehindTheSession(t *testing.T) {
	if code := serve(launchRouter(t), http.MethodGet, "/api/me/invitations", nil).Code; code != http.StatusNotFound {
		t.Errorf("without SSO = %d, want the route unmounted", code)
	}

	issuer := fakeRegistre(t)
	router, kits := ssoRouter(t, issuer)
	if code := serve(router, http.MethodGet, "/api/me/invitations", nil).Code; code != http.StatusUnauthorized {
		t.Errorf("without a session = %d, want 401", code)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/me/invitations", nil)
	request.Header.Set("Authorization", "Bearer "+signIn(t, kits, issuer))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	var body struct {
		Count int    `json:"count"`
		Href  string `json:"href"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil || recorder.Code != http.StatusOK {
		t.Fatalf("signed in = %d %q", recorder.Code, recorder.Body.String())
	}
	if body.Count != 1 || body.Href != issuer+"/invitations/inv-1" {
		t.Errorf("invitations = %+v, want one pointing at its Registre page", body)
	}
}
