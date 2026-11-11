package auth

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/FacileStudio/porte/kit"
	mcp "github.com/FacileStudio/tronc/mcp"
)

// newUsers returns the kit's store over porte_users, which is Capsule's only
// account table. It resolves a session to its profile and erases an account;
// sign-in keeps writing through porte/pg's own user store, which also owns
// facile_id and email_verified. Pastes reference no account, so there is no
// paste data to erase: EraseData drops the person's pending MCP actions,
// and the foreign keys carry the porte_users delete to identities, sessions,
// login codes and magic links.
func newUsers(sqlDB *sql.DB) *kit.Store {
	return &kit.Store{DB: sqlDB, Columns: kit.Columns{Table: "porte_users"}, EraseData: erasePending}
}

// erasePending deletes every MCP action the person asked for, inside the
// erasure's transaction.
func erasePending(ctx context.Context, tx *sql.Tx, userID int64) error {
	return mcp.ErasePending(ctx, tx, strconv.FormatInt(userID, 10))
}
