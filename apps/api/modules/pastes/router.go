package pastes

import "github.com/go-chi/chi/v5"

func RegisterRoutes(router chi.Router, service *Service) {
	handler := newHandler(service)

	router.Route("/api/pastes", func(r chi.Router) {
		r.Post("/", handler.create)
		r.Get("/{id}", handler.getMeta)
		r.Post("/{id}/content", handler.getContent)
		r.Delete("/{id}", handler.revoke)
	})
}
