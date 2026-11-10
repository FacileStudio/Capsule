package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/kit"
)

type missingUsers struct{}

func (missingUsers) Identity(context.Context, int64) (kit.Identity, error) {
	return kit.Identity{}, errors.New("no row")
}

func TestHandleSession(t *testing.T) {
	anonymous := httptest.NewRecorder()
	handleSession(missingUsers{})(anonymous, httptest.NewRequest(http.MethodGet, "/auth/session", nil))
	if anonymous.Code != http.StatusOK || anonymous.Body.String() != "{\"user\":null}\n" {
		t.Fatalf("anonymous: %d %q", anonymous.Code, anonymous.Body.String())
	}

	request := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	request = request.WithContext(porte.WithIdentity(request.Context(), porte.Identity{Email: "a@b.c", Name: "A"}))
	signedIn := httptest.NewRecorder()
	handleSession(missingUsers{})(signedIn, request)
	var body sessionResponse
	if err := json.Unmarshal(signedIn.Body.Bytes(), &body); err != nil || body.User == nil || body.User.Email != "a@b.c" || body.User.Name != "A" {
		t.Fatalf("signed in: %d %q", signedIn.Code, signedIn.Body.String())
	}
}

func TestSessionAndMeReadTheUserRow(t *testing.T) {
	sqlDB := testSQL(t)
	var userID int64
	err := sqlDB.QueryRow(
		`INSERT INTO porte_users (email, email_verified, name, avatar_url) VALUES ('row@example.test', true, 'Row Name', 'https://img.example.test/a.png') RETURNING id`,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	want := sessionUser{Email: "row@example.test", Name: "Row Name", AvatarURL: "https://img.example.test/a.png"}

	request := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	request = request.WithContext(porte.WithIdentity(request.Context(), porte.Identity{UserID: userID}))
	session := httptest.NewRecorder()
	handleSession(newUsers(sqlDB))(session, request)
	var body sessionResponse
	if err := json.Unmarshal(session.Body.Bytes(), &body); err != nil || body.User == nil || *body.User != want {
		t.Fatalf("session: got %q, want %+v", session.Body.String(), want)
	}

	identity, err := newUsers(sqlDB).Identity(context.Background(), userID)
	if err != nil {
		t.Fatalf("Identity: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	request = request.WithContext(kit.WithIdentity(request.Context(), identity))
	me := httptest.NewRecorder()
	handleMe(me, request)
	var got sessionUser
	if err := json.Unmarshal(me.Body.Bytes(), &got); err != nil || got != want {
		t.Fatalf("me: got %q, want %+v", me.Body.String(), want)
	}
}
