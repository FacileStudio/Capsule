package auth

import (
	"net/http"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/kit"
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

// RegisterRoutes mounts porte's routes through the kit, GET /auth/me, which
// the CLI reads to check a stored token, GET /auth/session, which the browser
// reads to know whether to show the user menu, and, with single sign-on,
// GET /me/invitations for the menu's invitation count.
func RegisterRoutes(router chi.Router, kits *Kits) error {
	if err := kit.Mount(router, kit.Config{Sessions: kits.Sessions, OIDC: kits.OIDC}); err != nil {
		return err
	}
	requireAuth := kit.RequireAuth(kits.Sessions, kits.Users)
	router.With(requireAuth).Get("/auth/me", handleMe)
	router.With(kits.Sessions.Optional).Get("/auth/session", handleSession(kits.Users))
	if kits.Invitations == nil {
		return nil
	}
	return kit.MountInvitations(router.With(requireAuth), *kits.Invitations)
}

func userFrom(identity kit.Identity) *sessionUser {
	return &sessionUser{Email: identity.Email, Name: identity.Name, AvatarURL: identity.PictureURL}
}

// handleSession answers 200 whether or not a session exists, so an anonymous
// visitor's browser never logs a failed request just to learn it is anonymous.
// A cookie session carries the user id and nothing else, so the profile comes
// from the porte user row; a lookup that fails falls back to what the
// identity holds.
func handleSession(users kit.Resolver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, ok := porte.From(r.Context())
		if !ok {
			httpjson.WriteJSON(w, http.StatusOK, sessionResponse{})
			return
		}
		identity, err := users.Identity(r.Context(), authenticated.UserID)
		if err != nil {
			identity = kit.Identity{Email: authenticated.Email, Name: authenticated.Name}
		}
		httpjson.WriteJSON(w, http.StatusOK, sessionResponse{User: userFrom(identity)})
	}
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	identity, _ := kit.FromContext(r.Context())
	httpjson.WriteJSON(w, http.StatusOK, userFrom(identity))
}
