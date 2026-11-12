package auth

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/FacileStudio/Capsule/apps/api/internal/env"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/kit"
	portepg "github.com/FacileStudio/porte/pg"
)

// deletionState keeps the GDPR receiver's cursor in the one row of
// deletion_state.
type deletionState struct {
	db *sql.DB
}

// Load returns the saved state, or nil when nothing was saved yet.
func (s deletionState) Load(ctx context.Context) ([]byte, error) {
	var state []byte
	err := s.db.QueryRowContext(ctx, `SELECT state FROM deletion_state WHERE id = 1`).Scan(&state)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return state, err
}

// Save replaces the saved state.
func (s deletionState) Save(ctx context.Context, state []byte) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO deletion_state (id, state) VALUES (1, $1)
		ON CONFLICT (id) DO UPDATE SET state = EXCLUDED.state, updated_at = now()`, state)
	return err
}

// newDeletions builds the receiver that polls Registre for erased people, with
// the same client credentials the app signs in with. The kit store erases the
// account first, then porte erases its own rows. A handled subject that an
// identity holds again is erased again (Holds). It returns nil when single
// sign-on is not configured, since there is no Registre to poll.
func newDeletions(sqlDB *sql.DB, store *portepg.Store, users *kit.Store, oidc *env.OIDCConfig, logger *slog.Logger) *directory.Deletions {
	if oidc == nil {
		return nil
	}
	return &directory.Deletions{
		Client: directory.Client{
			Issuer:       oidc.Issuer,
			ClientID:     oidc.ClientID,
			ClientSecret: oidc.ClientSecret,
		},
		OnPersonDeleted: users.EraseAccount,
		Erase:           store.EraseIdentity,
		Holds:           store.HoldsIdentity,
		State:           deletionState{db: sqlDB},
		Logger:          logger,
	}
}

// RunDeletions polls for erased people until ctx ends. It does nothing without
// single sign-on. Call it in exactly one replica.
func (k *Kits) RunDeletions(ctx context.Context) {
	if k.Deletions != nil {
		k.Deletions.Run(ctx)
	}
}
