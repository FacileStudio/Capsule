package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/porte"
)

func TestHandleSession(t *testing.T) {
	anonymous := httptest.NewRecorder()
	handleSession(nil)(anonymous, httptest.NewRequest(http.MethodGet, "/auth/session", nil))
	if anonymous.Code != http.StatusOK || anonymous.Body.String() != "{\"user\":null}\n" {
		t.Fatalf("anonymous: %d %q", anonymous.Code, anonymous.Body.String())
	}

	request := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	request = request.WithContext(porte.WithIdentity(request.Context(), porte.Identity{Email: "a@b.c", Name: "A"}))
	signedIn := httptest.NewRecorder()
	handleSession(nil)(signedIn, request)
	var body sessionResponse
	if err := json.Unmarshal(signedIn.Body.Bytes(), &body); err != nil || body.User == nil || body.User.Email != "a@b.c" || body.User.Name != "A" {
		t.Fatalf("signed in: %d %q", signedIn.Code, signedIn.Body.String())
	}
}
