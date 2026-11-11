package mcp

import "time"

// statusInput is capsule_get_status's input.
type statusInput struct {
	ID string `json:"id" jsonschema:"the capsule id, such as cap_0123456789abcdef: the last part of the link's path. Never pass the #key part of a link"`
}

// statusOutput is what the reveal page knows before a capsule is opened.
type statusOutput struct {
	ID            string     `json:"id"`
	Available     bool       `json:"available"`
	BurnAfterRead bool       `json:"burn_after_read,omitempty"`
	HasPassword   bool       `json:"has_password,omitempty"`
	Syntax        string     `json:"syntax,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	Message       string     `json:"message"`
}

// revokeInput is capsule_revoke's input.
type revokeInput struct {
	ID          string `json:"id" jsonschema:"the capsule id, such as cap_0123456789abcdef. Never pass the #key part of a link"`
	DeleteToken string `json:"delete_token" jsonschema:"the delete token shown once when the capsule was sealed"`
}

// revokeOutput is where a revoke stands.
type revokeOutput struct {
	Status     string `json:"status"`
	ID         string `json:"id,omitempty"`
	ConfirmURL string `json:"confirm_url,omitempty"`
	Message    string `json:"message"`
}
