package auth

import (
	"net/http"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/httpjson"
	"github.com/go-chi/chi/v5"
)

type meResponse struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// RegisterRoutes mounts porte's routes and GET /auth/me, which the CLI reads to
// check a stored token.
func RegisterRoutes(router chi.Router, kits *Kits) {
	kits.Sessions.Mount(router)
	kits.OIDC.Mount(router)
	router.With(kits.Sessions.RequireAuth).Get("/auth/me", handleMe)
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	identity, _ := porte.From(r.Context())
	httpjson.WriteJSON(w, http.StatusOK, meResponse{Email: identity.Email, Name: identity.Name})
}
