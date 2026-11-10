package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func createFrom(router chi.Router, remoteAddr string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, "/api/pastes", strings.NewReader(`{"content":"ciphertext"}`))
	request.Header.Set("Content-Type", "application/json")
	request.RemoteAddr = remoteAddr
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestPasteCreationIsLimitedPerVisitorWithTheErrorEnvelope(t *testing.T) {
	router := launchRouter(t)

	for i := range 30 {
		if code := createFrom(router, "203.0.113.7:"+strconv.Itoa(40000+i)).Code; code != http.StatusCreated {
			t.Fatalf("create %d = %d, want 201", i+1, code)
		}
	}
	refused := createFrom(router, "203.0.113.7:50000")
	if refused.Code != http.StatusTooManyRequests {
		t.Fatalf("create 31 from a new port = %d, want 429: the key is the host, not the connection", refused.Code)
	}
	if got := refused.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("429 Content-Type = %q, want application/json", got)
	}
	if want := `{"error":{"code":"rate_limited","message":"too many requests"}}`; strings.TrimSpace(refused.Body.String()) != want {
		t.Errorf("429 body = %q, want %q", refused.Body.String(), want)
	}
	if seconds, err := strconv.Atoi(refused.Header().Get("Retry-After")); err != nil || seconds < 1 {
		t.Errorf("Retry-After = %q, want whole seconds", refused.Header().Get("Retry-After"))
	}
	if code := createFrom(router, "198.51.100.9:1234").Code; code != http.StatusCreated {
		t.Errorf("another visitor = %d, want 201", code)
	}
}
