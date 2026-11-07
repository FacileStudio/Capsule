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

func TestCookieSessionMeReadsTheUserRow(t *testing.T) {
	sqlDB := testSQL(t)
	var userID int64
	err := sqlDB.QueryRow(
		`INSERT INTO porte_users (email, email_verified, name, avatar_url) VALUES ('row@example.test', true, 'Row Name', 'https://img.example.test/a.png') RETURNING id`,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	request := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	request = request.WithContext(porte.WithIdentity(request.Context(), porte.Identity{UserID: userID}))
	recorder := httptest.NewRecorder()
	handleMe(sqlDB)(recorder, request)
	var body sessionUser
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	want := sessionUser{Email: "row@example.test", Name: "Row Name", AvatarURL: "https://img.example.test/a.png"}
	if body != want {
		t.Fatalf("me: got %+v, want %+v", body, want)
	}
}
