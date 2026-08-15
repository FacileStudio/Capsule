package pastes

import (
	"net/http"

	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

// Handler wires the paste module's HTTP handlers onto a Service.
type Handler struct {
	service *Service
}

func newHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := httpjson.DecodeJSON(w, r, &req); err != nil {
		httpjson.WriteError(w, err)
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) getMeta(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp, err := h.service.GetMeta(id)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) getContent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp, err := h.service.GetContent(id)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) revoke(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	token := r.Header.Get("X-Delete-Token")

	if token == "" {
		httpjson.WriteError(w, errors.Forbidden("X-Delete-Token header is required"))
		return
	}

	if err := h.service.Revoke(id, token); err != nil {
		httpjson.WriteError(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}
