package auth

import (
	"database/sql"

	"github.com/FacileStudio/porte/kit"
)

// newUsers returns the kit's store over porte_users, which is Capsule's only
// account table. It resolves a session to its profile and erases an account;
// sign-in keeps writing through porte/pg's own user store, which also owns
// facile_id and email_verified. Pastes reference no account, so there is no
// EraseData: the porte_users row is all an erased person leaves, and the
// foreign keys carry its delete to identities, sessions, login codes and
// magic links.
func newUsers(sqlDB *sql.DB) *kit.Store {
	return &kit.Store{DB: sqlDB, Columns: kit.Columns{Table: "porte_users"}}
}
