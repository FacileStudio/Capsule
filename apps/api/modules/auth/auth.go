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

	"github.com/FacileStudio/Capsule/apps/api/internal/env"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/oidc"
	portepg "github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/porte/session"
)

// Kits is porte, assembled: one session manager over porte's tables and the
// OIDC kit that issues through it.
type Kits struct {
	DB       *sql.DB
	Sessions *session.Manager
	OIDC     *oidc.Kit

	// Deletions erases the people Registre erased. It is nil without single
	// sign-on. Run it in exactly one replica.
	Deletions *directory.Deletions
}

// New builds the session manager and runs OIDC discovery. A deployment with no
// OIDC_ISSUER still serves /auth/config, answering oidc_enabled false.
func New(ctx context.Context, sqlDB *sql.DB, appEnv env.Config, logger *slog.Logger) (*Kits, error) {
	cfg := appEnv.Porte()
	store := portepg.New(sqlDB)

	sessions, err := session.New(cfg, session.Deps{Sessions: store.Sessions(), Logger: logger})
	if err != nil {
		return nil, err
	}

	kit, err := oidc.New(ctx, cfg, oidc.Deps{
		Users:      store.Users(),
		Identities: store.Identities(),
		Sessions:   sessions,
		Codes:      store.LoginCodes(),
		Logger:     logger,
	})
	if err != nil {
		return nil, err
	}
	deletions := newDeletions(sqlDB, store, appEnv.OIDC, logger)
	return &Kits{DB: sqlDB, Sessions: sessions, OIDC: kit, Deletions: deletions}, nil
}
