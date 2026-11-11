package mcp

import (
	"context"
	"time"

	"github.com/FacileStudio/Capsule/apps/api/schemas"
	mcpkit "github.com/FacileStudio/tronc/mcp"
)

// kindRevoke is the one sensitive action: destroying a capsule, which cannot
// be undone.
const kindRevoke = "revoke_capsule"

// kinds are the sensitive actions. The fingerprint covers the capsule as it
// was when the person was asked, and a capsule opened or expired in between
// is gone, so the action closes instead of running.
func (t *toolset) kinds() map[string]mcpkit.Kind {
	return map[string]mcpkit.Kind{
		kindRevoke: {Tool: toolRevoke, Fingerprint: t.fingerprint, Run: t.runRevoke},
	}
}

// fingerprint digests the targeted capsule while it can still be opened.
func (t *toolset) fingerprint(_ context.Context, a mcpkit.Action) (string, error) {
	paste, err := t.Pastes.Live(a.Target)
	if err != nil {
		return "", err
	}
	return digest(paste), nil
}

// digest fingerprints a capsule by what identifies it, never by its token or
// its ciphertext.
func digest(p *schemas.Paste) string {
	expires := "-"
	if p.ExpiresAt != nil {
		expires = p.ExpiresAt.UTC().Format(time.RFC3339Nano)
	}
	return mcpkit.Digest(p.ID, p.CreatedAt.UTC().Format(time.RFC3339Nano), expires)
}

// runRevoke burns the targeted capsule.
func (t *toolset) runRevoke(_ context.Context, a mcpkit.Action) error {
	return t.Pastes.Burn(a.Target)
}
