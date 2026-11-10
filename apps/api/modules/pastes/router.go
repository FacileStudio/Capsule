package pastes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts the paste endpoints under /pastes: create, meta, content
// and revoke. Create is rate-limited by address.
func RegisterRoutes(router chi.Router, service *Service, createLimiter func(http.Handler) http.Handler) {
	handler := newHandler(service)

	router.Route("/pastes", func(r chi.Router) {
		r.With(createLimiter).Post("/", handler.create)
		r.Get("/{id}", handler.getMeta)
		r.Post("/{id}/content", handler.getContent)
		r.Delete("/{id}", handler.revoke)
	})
}
