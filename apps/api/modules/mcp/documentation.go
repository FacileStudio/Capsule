package mcp

import (
	"github.com/FacileStudio/tronc/apiref"
	mcpkit "github.com/FacileStudio/tronc/mcp"
)

// Documentation describes the confirm page routes for the API reference. The
// MCP endpoint itself is outside /api: POST /mcp speaks the Model Context
// Protocol to any MCP client holding a Registre access token for this server,
// and /.well-known/oauth-protected-resource/mcp describes how to get one.
var Documentation = apiref.Module{
	Name: "mcp",
	Description: "Capsule revokes MCP clients asked for, waiting for the person to confirm them. A client that cannot " +
		"ask the person itself hands them a link to the pending page, which uses these routes. The MCP endpoint is " +
		"POST /mcp, outside /api.",
	Routes: mcpkit.PendingRoutes[PendingItem](mcpkit.PendingDoc{
		Auth: "session cookie only; a bearer token is refused",
		List: "Open actions, newest first, with their capsule's metadata as it is now. state is ready, changed or " +
			"missing (the capsule was opened, revoked or expired since the request).",
		Confirm: "Runs the action once and destroys the capsule. A capsule that changed since the request is " +
			"refused with 412 and the action closes.",
		Dismiss: "Closes the action without running it.",
		Changed: "The capsule changed after the request.",
	}),
}
