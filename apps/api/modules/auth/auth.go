// Package auth is Capsule's sign-in, delegated entirely to porte.
//
// Pastes stay anonymous: nothing here gates them. Signing in exists so the CLI
// can hold a session, through the same one-step browser flow as the rest of the
// suite. Sign-in is single sign-on only, so no password kit is built.
package auth

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/FacileStudio/Capsule/apps/api/internal/env"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/kit"
	"github.com/FacileStudio/porte/oidc"
	"github.com/FacileStudio/porte/session"
)

// Kits is porte, assembled: one session manager over porte's tables, the
// OIDC kit that issues through it, and the kit store that resolves and erases
// accounts.
type Kits struct {
	Sessions *session.Manager
	OIDC     *oidc.Kit
	Users    *kit.Store

	// Identities resolves a Registre subject to its account, for the MCP
	// token verifier.
	Identities porte.IdentityStore

	// Deletions erases the people Registre erased. It is nil without single
	// sign-on. Run it in exactly one replica.
	Deletions *directory.Deletions

	// Invitations configures the user menu's invitation count. It is nil
	// without single sign-on.
	Invitations *kit.InvitationsConfig
}

// New builds the session manager and runs OIDC discovery. A deployment with no
// OIDC_ISSUER still serves /auth/config, answering oidc_enabled false.
func New(ctx context.Context, sqlDB *sql.DB, appEnv env.Config, logger *slog.Logger) (*Kits, error) {
	cfg := appEnv.Porte()
	store, err := newPorteStore(sqlDB)
	if err != nil {
		return nil, err
	}
	resealInBackground(ctx, store, logger)

	sessions, err := session.New(cfg, session.Deps{Sessions: store.Sessions(), Logger: logger})
	if err != nil {
		return nil, err
	}

	oidcKit, err := oidc.New(ctx, cfg, oidc.Deps{
		Users:      store.Users(),
		Identities: store.Identities(),
		Sessions:   sessions,
		Codes:      store.LoginCodes(),
		Logger:     logger,
	})
	if err != nil {
		return nil, err
	}
	users := newUsers(sqlDB)
	deletions := newDeletions(sqlDB, store, users, appEnv.OIDC, logger)
	return &Kits{
		Sessions:    sessions,
		OIDC:        oidcKit,
		Users:       users,
		Identities:  store.Identities(),
		Deletions:   deletions,
		Invitations: newInvitations(deletions, store, logger),
	}, nil
}

// RequireAuth is the middleware that admits a signed-in person, resolved to
// their account as /auth/me resolves them.
func (k *Kits) RequireAuth() func(http.Handler) http.Handler {
	return kit.RequireAuth(k.Sessions, k.Users)
}
