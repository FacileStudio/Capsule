package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"

	"github.com/FacileStudio/Capsule/apps/api/internal/env"
	"github.com/FacileStudio/tronc/apiref"
	"github.com/go-chi/chi/v5"
)

// testRouter builds the real router against nil dependencies. Route
// registration never touches the database, so the shape of the router is
// faithful even though no handler could serve a request.
func testRouter(t *testing.T) chi.Router {
	t.Helper()
	appEnv := env.Config{}
	router, err := buildRouter(nil, nil, appEnv, slog.Default())
	if err != nil {
		t.Fatalf("buildRouter: %v", err)
	}
	return router
}

func TestReferenceAndSpecAreServed(t *testing.T) {
	router := testRouter(t)

	page := httptest.NewRecorder()
	router.ServeHTTP(page, httptest.NewRequest(http.MethodGet, apiref.BasePath, nil))
	if page.Code != http.StatusOK {
		t.Fatalf("GET %s = %d, want 200", apiref.BasePath, page.Code)
	}
	body := page.Body.String()
	if !strings.Contains(body, `data-url="`+apiref.SpecPath+`"`) {
		t.Errorf("reference page does not point at its own spec:\n%s", body)
	}
	if !strings.Contains(body, apiref.ScalarScriptURL) {
		t.Errorf("reference page does not load the pinned Scalar bundle:\n%s", body)
	}

	if document := fetchSpec(t, router); document["openapi"] != "3.1.0" {
		t.Errorf("openapi = %v, want 3.1.0", document["openapi"])
	}
}

// TestEveryRouteIsDocumented is the drift guard. Capsule's spec is hand-written
// rather than generated from a route registry, so nothing but this test stops it
// from describing an API the binary no longer serves. It caught the /api prefix
// missing from every documented path.
func TestEveryRouteIsDocumented(t *testing.T) {
	documented := documentedRoutes(t, fetchSpec(t, testRouter(t)))

	var missing []string
	err := chi.Walk(testRouter(t), func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		route = strings.TrimSuffix(route, "/")
		if route == "" || strings.HasSuffix(route, "*") || ignored(route) {
			return nil
		}
		if !documented[method+" "+route] {
			missing = append(missing, method+" "+route)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking the router: %v", err)
	}

	sort.Strings(missing)
	if len(missing) > 0 {
		t.Errorf("routes missing from openapi.yaml: %v", missing)
	}
}

// The reference page documents itself, and the probes are infrastructure served
// at both mount points.
func ignored(route string) bool {
	switch route {
	case apiref.BasePath, apiref.SpecPath, "/health", "/ready", "/api/health", "/api/ready":
		return true
	}
	return false
}

func fetchSpec(t *testing.T, router chi.Router) map[string]any {
	t.Helper()
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, apiref.SpecPath, nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("GET %s = %d, want 200", apiref.SpecPath, recorder.Code)
	}
	var document map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &document); err != nil {
		t.Fatalf("spec is not valid JSON: %v", err)
	}
	return document
}

// documentedRoutes expands the spec's paths against every server base path, so
// an inventory written as /pastes matches a router serving /api/pastes.
func documentedRoutes(t *testing.T, document map[string]any) map[string]bool {
	t.Helper()

	prefixes := []string{""}
	servers, _ := document["servers"].([]any)
	for _, entry := range servers {
		raw, _ := entry.(map[string]any)["url"].(string)
		parsed, err := url.Parse(raw)
		if err != nil {
			t.Fatalf("server url %q: %v", raw, err)
		}
		if base := strings.TrimSuffix(parsed.Path, "/"); base != "" {
			prefixes = append(prefixes, base)
		}
	}

	routes := map[string]bool{}
	paths, _ := document["paths"].(map[string]any)
	for path, item := range paths {
		operations, _ := item.(map[string]any)
		for method := range operations {
			for _, prefix := range prefixes {
				routes[strings.ToUpper(method)+" "+prefix+path] = true
			}
		}
	}
	return routes
}
