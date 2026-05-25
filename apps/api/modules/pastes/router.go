package pastes

import (
	"github.com/FacileStudio/Capsule/apps/api/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router, service *Service, createLimiter *middleware.RateLimiter) {
	handler := newHandler(service)

	router.Route("/pastes", func(r chi.Router) {
		r.With(createLimiter.Middleware).Post("/", handler.create)
		r.Get("/{id}", handler.getMeta)
		r.Post("/{id}/content", handler.getContent)
		r.Delete("/{id}", handler.revoke)
	})
}
