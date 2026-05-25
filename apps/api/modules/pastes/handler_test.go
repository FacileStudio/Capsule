package pastes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Capsule/apps/api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func setupRouter(t *testing.T) (chi.Router, *Service) {
	t.Helper()
	db := setupTestDB(t)
	svc := NewService(db, 1024)
	limiter := middleware.NewRateLimiter(5, time.Minute)

	r := chi.NewRouter()
	RegisterRoutes(r, svc, limiter)
	return r, svc
}

func doRequest(r chi.Router, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = &bytes.Buffer{}
	}
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func parseJSON(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&m); err != nil {
		t.Fatalf("decode response: %v (body: %s)", err, rr.Body.String())
	}
	return m
}

func TestHandlerCreateSuccess(t *testing.T) {
	r, _ := setupRouter(t)

	rr := doRequest(r, "POST", "/pastes", map[string]any{
		"content": "handler test",
	}, nil)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	body := parseJSON(t, rr)
	id, ok := body["id"].(string)
	if !ok || !strings.HasPrefix(id, "cap_") {
		t.Errorf("expected id with cap_ prefix, got %v", body["id"])
	}
	if _, ok := body["delete_token"].(string); !ok {
		t.Error("expected delete_token in response")
	}
}

func TestHandlerCreateEmptyContent(t *testing.T) {
	r, _ := setupRouter(t)

	rr := doRequest(r, "POST", "/pastes", map[string]any{
		"content": "",
	}, nil)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandlerCreateContentTooLarge(t *testing.T) {
	r, _ := setupRouter(t)

	rr := doRequest(r, "POST", "/pastes", map[string]any{
		"content": strings.Repeat("x", 2000),
	}, nil)

	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected 413, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandlerCreateInvalidJSON(t *testing.T) {
	r, _ := setupRouter(t)

	req := httptest.NewRequest("POST", "/pastes", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandlerCreateInvalidExpiry(t *testing.T) {
	r, _ := setupRouter(t)

	rr := doRequest(r, "POST", "/pastes", map[string]any{
		"content":    "test",
		"expires_in": "99y",
	}, nil)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandlerGetMeta(t *testing.T) {
	r, _ := setupRouter(t)

	createRR := doRequest(r, "POST", "/pastes", map[string]any{
		"content": "meta test",
		"syntax":  "go",
	}, nil)
	created := parseJSON(t, createRR)
	id := created["id"].(string)

	rr := doRequest(r, "GET", "/pastes/"+id, nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	body := parseJSON(t, rr)
	if body["exists"] != true {
		t.Error("expected exists=true")
	}
	if body["syntax"] != "go" {
		t.Errorf("expected syntax=go, got %v", body["syntax"])
	}
}

func TestHandlerGetMetaNotFound(t *testing.T) {
	r, _ := setupRouter(t)

	rr := doRequest(r, "GET", "/pastes/cap_0000000000000000", nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	body := parseJSON(t, rr)
	if body["exists"] != false {
		t.Error("expected exists=false for non-existent paste")
	}
}

func setupPostgresRouter(t *testing.T) chi.Router {
	t.Helper()
	db := setupPostgresDB(t)
	svc := NewService(db, 1024)
	limiter := middleware.NewRateLimiter(100, time.Minute)
	r := chi.NewRouter()
	RegisterRoutes(r, svc, limiter)
	return r
}

func TestHandlerGetContent(t *testing.T) {
	r := setupPostgresRouter(t)

	createRR := doRequest(r, "POST", "/pastes", map[string]any{
		"content":         "reveal me",
		"burn_after_read": false,
	}, nil)
	created := parseJSON(t, createRR)
	id := created["id"].(string)

	rr := doRequest(r, "POST", "/pastes/"+id+"/content", nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	body := parseJSON(t, rr)
	if body["content"] != "reveal me" {
		t.Errorf("expected content='reveal me', got %v", body["content"])
	}
}

func TestHandlerGetContentBurnAfterRead(t *testing.T) {
	r := setupPostgresRouter(t)

	createRR := doRequest(r, "POST", "/pastes", map[string]any{
		"content":         "burn me",
		"burn_after_read": true,
	}, nil)
	created := parseJSON(t, createRR)
	id := created["id"].(string)

	rr1 := doRequest(r, "POST", "/pastes/"+id+"/content", nil, nil)
	if rr1.Code != http.StatusOK {
		t.Fatalf("first read: expected 200, got %d", rr1.Code)
	}

	rr2 := doRequest(r, "POST", "/pastes/"+id+"/content", nil, nil)
	if rr2.Code != http.StatusNotFound {
		t.Errorf("second read: expected 404, got %d", rr2.Code)
	}
}

func TestHandlerGetContentNotFound(t *testing.T) {
	r := setupPostgresRouter(t)

	rr := doRequest(r, "POST", "/pastes/cap_0000000000000000/content", nil, nil)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestHandlerRevoke(t *testing.T) {
	r, _ := setupRouter(t)

	createRR := doRequest(r, "POST", "/pastes", map[string]any{
		"content": "revoke me",
	}, nil)
	created := parseJSON(t, createRR)
	id := created["id"].(string)
	token := created["delete_token"].(string)

	rr := doRequest(r, "DELETE", "/pastes/"+id, nil, map[string]string{
		"X-Delete-Token": token,
	})
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	metaRR := doRequest(r, "GET", "/pastes/"+id, nil, nil)
	meta := parseJSON(t, metaRR)
	if meta["exists"] != false {
		t.Error("revoked paste should report exists=false")
	}
}

func TestHandlerRevokeMissingToken(t *testing.T) {
	r, _ := setupRouter(t)

	createRR := doRequest(r, "POST", "/pastes", map[string]any{
		"content": "test",
	}, nil)
	created := parseJSON(t, createRR)
	id := created["id"].(string)

	rr := doRequest(r, "DELETE", "/pastes/"+id, nil, nil)
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rr.Code)
	}
}

func TestHandlerRevokeWrongToken(t *testing.T) {
	r, _ := setupRouter(t)

	createRR := doRequest(r, "POST", "/pastes", map[string]any{
		"content": "test",
	}, nil)
	created := parseJSON(t, createRR)
	id := created["id"].(string)

	rr := doRequest(r, "DELETE", "/pastes/"+id, nil, map[string]string{
		"X-Delete-Token": "wrong",
	})
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rr.Code)
	}
}

func TestHandlerRevokeNotFound(t *testing.T) {
	r, _ := setupRouter(t)

	rr := doRequest(r, "DELETE", "/pastes/cap_0000000000000000", nil, map[string]string{
		"X-Delete-Token": "whatever",
	})
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestHandlerRateLimiting(t *testing.T) {
	r, _ := setupRouter(t)

	for i := range 6 {
		rr := doRequest(r, "POST", "/pastes", map[string]any{
			"content": "rate limit test",
		}, nil)

		if i < 5 {
			if rr.Code != http.StatusCreated {
				t.Fatalf("request %d: expected 201, got %d", i+1, rr.Code)
			}
		} else {
			if rr.Code != http.StatusTooManyRequests {
				t.Errorf("request %d: expected 429, got %d", i+1, rr.Code)
			}
		}
	}
}

func TestHandlerResponseHeaders(t *testing.T) {
	r, _ := setupRouter(t)

	rr := doRequest(r, "POST", "/pastes", map[string]any{
		"content": "headers test",
	}, nil)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestHandlerErrorFormat(t *testing.T) {
	r, _ := setupRouter(t)

	rr := doRequest(r, "POST", "/pastes", map[string]any{
		"content": "",
	}, nil)

	body := parseJSON(t, rr)
	errObj, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error object in response, got %v", body)
	}
	if _, ok := errObj["code"].(string); !ok {
		t.Error("expected error.code string")
	}
	if _, ok := errObj["message"].(string); !ok {
		t.Error("expected error.message string")
	}
}
