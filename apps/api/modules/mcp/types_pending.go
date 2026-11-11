package mcp

import (
	"time"

	mcpkit "github.com/FacileStudio/tronc/mcp"
)

// PendingItem is one action waiting on the confirm page.
type PendingItem struct {
	mcpkit.PendingItem
	Capsule *PendingCapsule `json:"capsule,omitempty"`
}

// PendingCapsule is the capsule an action targets, as it is now. It carries
// metadata only: the server never holds a capsule's plaintext.
type PendingCapsule struct {
	ID            string     `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	BurnAfterRead bool       `json:"burn_after_read"`
	HasPassword   bool       `json:"has_password"`
}
