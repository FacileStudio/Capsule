package mcp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/errors"
	mcpkit "github.com/FacileStudio/tronc/mcp"
)

const webActor = "capsule-web"

// pendingUser is the signed-in person of a confirm page request. Only the
// person's own browser session passes: a bearer token, a CLI's included, and
// a delegated token never confirm what an agent asked for.
func pendingUser(r *http.Request) (string, error) {
	identity, ok := porte.From(r.Context())
	if !ok {
		return "", errors.Unauthorized("missing auth")
	}
	if r.Header.Get("Authorization") != "" || identity.SessionID == 0 || identity.Actor != "" {
		return "", mcpkit.ErrDelegated
	}
	return strconv.FormatInt(identity.UserID, 10), nil
}

// describe adds the targeted capsule, as it is now, to a confirm page item.
func (t *toolset) describe(_ context.Context, item mcpkit.PendingItem, a mcpkit.Action) PendingItem {
	out := PendingItem{PendingItem: item}
	if paste, err := t.Pastes.Live(a.Target); err == nil {
		out.Capsule = &PendingCapsule{
			ID: paste.ID, CreatedAt: paste.CreatedAt, ExpiresAt: paste.ExpiresAt,
			BurnAfterRead: paste.BurnAfterRead, HasPassword: paste.HasPassword,
		}
	}
	return out
}

// done is the confirm page's answer once an action ran: the capsule burned.
func done(context.Context, mcpkit.Outcome) mcpkit.PendingResult {
	return mcpkit.PendingResult{Status: statusRevoked}
}
