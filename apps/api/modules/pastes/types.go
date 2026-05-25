package pastes

import "time"

type CreateRequest struct {
	Content       string `json:"content"`
	BurnAfterRead *bool  `json:"burn_after_read"`
	ExpiresIn     string `json:"expires_in"`
	MaxViews      *int   `json:"max_views"`
	HasPassword   bool   `json:"has_password"`
	Syntax        string `json:"syntax"`
}

type CreateResponse struct {
	ID          string     `json:"id"`
	DeleteToken string     `json:"delete_token"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type MetaResponse struct {
	ID          string     `json:"id"`
	Exists      bool       `json:"exists"`
	Burned      bool       `json:"burned,omitempty"`
	HasPassword bool       `json:"has_password,omitempty"`
	Syntax      string     `json:"syntax,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

type ContentResponse struct {
	Content string `json:"content"`
}
