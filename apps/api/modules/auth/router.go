package auth

import (
	"database/sql"
	"net/http"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/httpjson"
	"github.com/go-chi/chi/v5"
)

type sessionUser struct {
	Email     string `json:"email"`
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type sessionResponse struct {
	User *sessionUser `json:"user"`
}

// RegisterRoutes mounts porte's routes, GET /auth/me, which the CLI reads to
// check a stored token, and GET /auth/session, which the browser reads to know
// whether to show the user menu.
func RegisterRoutes(router chi.Router, kits *Kits) {
	kits.Sessions.Mount(router)
	kits.OIDC.Mount(router)
	router.With(kits.Sessions.RequireAuth).Get("/auth/me", handleMe(kits.DB))
	router.With(kits.Sessions.Optional).Get("/auth/session", handleSession(kits.DB))
}

// accountUser returns the caller as the porte user row holds it. A cookie
// session carries the user id and nothing else, so the name, the email and
// Registre's picture come from that row; a lookup that fails falls back to what
// the identity holds.
func accountUser(r *http.Request, db *sql.DB, identity porte.Identity) *sessionUser {
	user := &sessionUser{Email: identity.Email, Name: identity.Name}
	if db != nil {
		_ = db.QueryRowContext(r.Context(), `SELECT email, name, avatar_url FROM porte_users WHERE id = $1`, identity.UserID).
			Scan(&user.Email, &user.Name, &user.AvatarURL)
	}
	return user
}

// handleSession answers 200 whether or not a session exists, so an anonymous
// visitor's browser never logs a failed request just to learn it is anonymous.
func handleSession(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identity, ok := porte.From(r.Context())
		if !ok {
			httpjson.WriteJSON(w, http.StatusOK, sessionResponse{})
			return
		}
		httpjson.WriteJSON(w, http.StatusOK, sessionResponse{User: accountUser(r, db, identity)})
	}
}

func handleMe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identity, _ := porte.From(r.Context())
		httpjson.WriteJSON(w, http.StatusOK, accountUser(r, db, identity))
	}
}
