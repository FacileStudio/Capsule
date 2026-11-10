package auth

import (
	"log/slog"

	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/kit"
	portepg "github.com/FacileStudio/porte/pg"
)

// newInvitations configures GET /me/invitations for the user menu with the
// directory client the deletion receiver already holds. It returns nil
// without single sign-on, which leaves the route unmounted.
func newInvitations(deletions *directory.Deletions, store *portepg.Store, logger *slog.Logger) *kit.InvitationsConfig {
	if deletions == nil {
		return nil
	}
	return &kit.InvitationsConfig{Directory: deletions.Client, Identities: store.Identities(), Logger: logger}
}
