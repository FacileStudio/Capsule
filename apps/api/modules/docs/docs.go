package docs

import (
	"net/http"

	"github.com/FacileStudio/tronc/apiref"
	"github.com/go-chi/chi/v5"
)

type (
	Registry = apiref.Registry
	Module   = apiref.Module
	Route    = apiref.Route
	Field    = apiref.Field
	Error    = apiref.Error
)

type CreateRequest struct {
	Content       string `json:"content" validate:"required"`
	BurnAfterRead *bool  `json:"burn_after_read,omitempty"`
	ExpiresIn     string `json:"expires_in,omitempty"`
	MaxViews      int    `json:"max_views,omitempty"`
	HasPassword   bool   `json:"has_password,omitempty"`
	Syntax        string `json:"syntax,omitempty"`
}

type CreateResponse struct {
	ID          string  `json:"id"`
	DeleteToken string  `json:"delete_token"`
	ExpiresAt   *string `json:"expires_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

type MetaResponse struct {
	ID          string  `json:"id"`
	Exists      bool    `json:"exists"`
	Burned      bool    `json:"burned,omitempty"`
	HasPassword bool    `json:"has_password,omitempty"`
	Syntax      string  `json:"syntax,omitempty"`
	ExpiresAt   *string `json:"expires_at,omitempty"`
	CreatedAt   *string `json:"created_at,omitempty"`
}

type ContentResponse struct {
	Content string `json:"content"`
}

type DeleteResponse struct {
	Deleted bool `json:"deleted"`
}

// Mount mounts the reference page and the OpenAPI document on router.
func Mount(router chi.Router) {
	apiref.Mount(router, Reference())
}

// Reference returns Capsule's API reference configuration.
func Reference() apiref.Config {
	return apiref.Config{
		Title:       "Capsule API",
		Description: "Zero-knowledge encrypted paste sharing. The server stores ciphertext — it never sees your plaintext.",
		Servers:     []string{"/api"},
		Registry: Registry{
			Modules: []Module{
				{
					Name:        "pastes",
					Description: "Create, read, and revoke encrypted pastes",
					Routes: []Route{
						{
							Method:       "POST",
							Path:         "/pastes",
							Summary:      "Store encrypted content on the server",
							RequestBody:  CreateRequest{},
							ResponseBody: CreateResponse{},
							Status:       http.StatusCreated,
						},
						{
							Method:       "GET",
							Path:         "/pastes/{id}",
							Summary:      "Check if a paste exists and retrieve metadata",
							PathParams:   []Field{{Name: "id", Type: "string", Description: "Paste ID"}},
							ResponseBody: MetaResponse{},
						},
						{
							Method:       "DELETE",
							Path:         "/pastes/{id}",
							Summary:      "Permanently destroy an encrypted paste",
							Auth:         "X-Delete-Token",
							PathParams:   []Field{{Name: "id", Type: "string", Description: "Paste ID"}},
							ResponseBody: DeleteResponse{},
						},
						{
							Method:       "POST",
							Path:         "/pastes/{id}/content",
							Summary:      "Retrieve encrypted paste ciphertext",
							PathParams:   []Field{{Name: "id", Type: "string", Description: "Paste ID"}},
							ResponseBody: ContentResponse{},
						},
					},
				},
			},
		},
	}
}
