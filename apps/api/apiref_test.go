package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/Capsule/apps/api/internal/env"
	"github.com/FacileStudio/Capsule/apps/api/modules/docs"
	"github.com/FacileStudio/tronc/apiref"
	"github.com/go-chi/chi/v5"
)

func testRouter(t *testing.T) chi.Router {
	t.Helper()
	appEnv := env.Config{}
	router, err := buildRouter(nil, nil, appEnv, slog.Default())
	if err != nil {
		t.Fatalf("buildRouter: %v", err)
	}
	return router
}

func TestEveryRouteIsDocumented(t *testing.T) {
	router := testRouter(t)
	if missing := apiref.Undocumented(router, docs.Reference()); len(missing) > 0 {
		t.Errorf("routes missing from the API registry: %v", missing)
	}
}

func TestRegistryIsComplete(t *testing.T) {
	if issues := apiref.Incomplete(docs.Reference(), "/pastes/{id}/content"); len(issues) > 0 {
		t.Errorf("incomplete documentation routes: %v", issues)
	}
}

func TestReferenceIsServedAtDocs(t *testing.T) {
	router := testRouter(t)

	page := httptest.NewRecorder()
	router.ServeHTTP(page, httptest.NewRequest(http.MethodGet, "/docs", nil))
	if page.Code != http.StatusOK {
		t.Fatalf("GET /docs = %d, want 200", page.Code)
	}

	spec := httptest.NewRecorder()
	router.ServeHTTP(spec, httptest.NewRequest(http.MethodGet, "/docs/openapi.json", nil))
	if spec.Code != http.StatusOK {
		t.Fatalf("GET /docs/openapi.json = %d, want 200", spec.Code)
	}
	var document struct {
		OpenAPI string         `json:"openapi"`
		Paths   map[string]any `json:"paths"`
	}
	if err := json.Unmarshal(spec.Body.Bytes(), &document); err != nil {
		t.Fatalf("spec is not JSON: %v", err)
	}
	if document.OpenAPI == "" || len(document.Paths) == 0 {
		t.Fatalf("spec is empty: %+v", document)
	}
}
