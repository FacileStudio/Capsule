package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FacileStudio/Capsule/apps/api/internal/env"
	"github.com/FacileStudio/Capsule/apps/api/internal/testsupport"
	"github.com/FacileStudio/Capsule/apps/api/modules/auth"
	"github.com/go-chi/chi/v5"
)

const inlineScript = "window.__capsule = 1;"

func launchRouter(t *testing.T) chi.Router {
	t.Helper()
	db := testsupport.DB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql handle: %v", err)
	}
	client := t.TempDir()
	page := "<!doctype html><script>" + inlineScript + "</script><div id=app></div>"
	if err := os.WriteFile(filepath.Join(client, "index.html"), []byte(page), 0o600); err != nil {
		t.Fatalf("write client: %v", err)
	}
	t.Setenv("CLIENT_DIR", client)

	appEnv := env.Config{MaxPasteSize: 1024}
	kits, err := auth.New(context.Background(), sqlDB, appEnv, slog.Default())
	if err != nil {
		t.Fatalf("auth.New: %v", err)
	}
	router, err := buildRouter(db, sqlDB, kits, appEnv, slog.Default())
	if err != nil {
		t.Fatalf("buildRouter: %v", err)
	}
	return router
}

func serve(router chi.Router, method, path string, body any) *httptest.ResponseRecorder {
	var payload bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&payload).Encode(body)
	}
	req := httptest.NewRequest(method, path, &payload)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func TestAnonymousRevealNeedsNoSession(t *testing.T) {
	router := launchRouter(t)

	created := serve(router, http.MethodPost, "/api/pastes", map[string]any{"content": "ciphertext"})
	if created.Code != http.StatusCreated {
		t.Fatalf("create = %d: %s", created.Code, created.Body.String())
	}
	var paste struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(created.Body.Bytes(), &paste); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	if meta := serve(router, http.MethodGet, "/api/pastes/"+paste.ID, nil); meta.Code != http.StatusOK {
		t.Errorf("anonymous meta = %d, want 200", meta.Code)
	}
	content := serve(router, http.MethodPost, "/api/pastes/"+paste.ID+"/content", nil)
	if content.Code != http.StatusOK || !strings.Contains(content.Body.String(), "ciphertext") {
		t.Errorf("anonymous content = %d %q, want 200 with the ciphertext", content.Code, content.Body.String())
	}
}

func TestClientPagesKeepTheKeyFragmentPrivate(t *testing.T) {
	router := launchRouter(t)

	recorder := serve(router, http.MethodGet, "/cap_0000000000000000", nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("GET client page = %d, want 200", recorder.Code)
	}
	if got := recorder.Header().Get("Referrer-Policy"); got != "no-referrer" {
		t.Errorf("Referrer-Policy = %q, want no-referrer", got)
	}

	sum := sha256.Sum256([]byte(inlineScript))
	hash := "'sha256-" + base64.StdEncoding.EncodeToString(sum[:]) + "'"
	policy := recorder.Header().Get("Content-Security-Policy-Report-Only")
	for _, want := range []string{hash, analyticsOrigin} {
		if !strings.Contains(policy, want) {
			t.Errorf("CSP %q does not allow %s", policy, want)
		}
	}
	for _, name := range []string{"Content-Security-Policy", "Content-Security-Policy-Report-Only"} {
		if strings.Contains(recorder.Header().Get(name), "unsafe-eval") {
			t.Errorf("%s allows unsafe-eval, which Web Crypto does not need", name)
		}
	}
}

func TestAPIReferrerPolicyNeverSendsThePath(t *testing.T) {
	router := launchRouter(t)

	got := serve(router, http.MethodGet, "/api/pastes/cap_0000000000000000", nil).Header().Get("Referrer-Policy")
	if got != "strict-origin-when-cross-origin" {
		t.Errorf("Referrer-Policy = %q, want tronc's strict-origin-when-cross-origin", got)
	}
}

func TestVersionIsServed(t *testing.T) {
	router := launchRouter(t)

	recorder := serve(router, http.MethodGet, "/version", nil)
	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), version) {
		t.Errorf("GET /version = %d %q, want 200 carrying %q", recorder.Code, recorder.Body.String(), version)
	}
}
