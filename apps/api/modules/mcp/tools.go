package mcp

import (
	"context"
	stderrors "errors"
	"regexp"

	"github.com/FacileStudio/Capsule/apps/api/modules/pastes"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	toolStatus = "capsule_get_status"
	toolRevoke = "capsule_revoke"
	hintGone   = "Nothing is left to read or revoke; to share the secret again, the user seals a new capsule in the browser."
)

var (
	idPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,24}$`)
	errID     = stderrors.New("pass the capsule id only, such as cap_0123456789abcdef: the last part of the link's path, " +
		"without the domain and without the #key part, which must never be sent to any tool")
	errToken = stderrors.New("this delete token does not match the capsule. It is the token shown once when the " +
		"capsule was sealed; without it the capsule cannot be revoked")
)

// toolset holds what the tool handlers share.
type toolset struct {
	Deps
	confirm *mcpkit.Confirmations
}

// list is the curated tool set: one read and the revoke the user confirms.
func (t *toolset) list() []mcpkit.Tool {
	return []mcpkit.Tool{
		mcpkit.NewTool(mcpkit.ToolSpec{
			Name: toolStatus, Title: "Check a capsule",
			Description: "Check whether a capsule can still be opened, without opening it: available, when it was sealed " +
				"and expires, whether opening it destroys it and whether it needs a password. It never returns the " +
				"content, which only the link's key can decrypt in a browser. A capsule that was opened, revoked or " +
				"expired comes back available false.",
			Annotations: &mcpkit.Annotations{ReadOnly: true, Idempotent: true, OpenWorld: true},
		}, t.status),
		mcpkit.NewTool(mcpkit.ToolSpec{
			Name: toolRevoke, Title: "Revoke a capsule",
			Description: "Destroy a capsule for good before anyone opens it, with the delete token shown when it was " +
				"sealed. The user always confirms: in your client when it can ask, otherwise the result has status " +
				"pending_confirmation and a confirm_url to give the user. Nothing is destroyed before they confirm.",
			Annotations: &mcpkit.Annotations{Destructive: true, Idempotent: true},
		}, t.revoke),
	}
}

// status is capsule_get_status: the reveal page's metadata for one capsule,
// read through the same service call as GET /pastes/{id}.
func (t *toolset) status(_ context.Context, _ *sdk.CallToolRequest, in statusInput) (*sdk.CallToolResult, statusOutput, error) {
	if !idPattern.MatchString(in.ID) {
		return nil, statusOutput{}, errID
	}
	meta, err := t.Pastes.GetMeta(in.ID)
	if err != nil {
		return nil, statusOutput{}, mcpkit.Failure(err, hintGone)
	}
	return nil, describeStatus(meta), nil
}

// describeStatus words the reveal page's metadata for a model.
func describeStatus(m *pastes.MetaResponse) statusOutput {
	out := statusOutput{
		ID: m.ID, Available: m.Exists, BurnAfterRead: m.BurnAfterRead, HasPassword: m.HasPassword, Syntax: m.Syntax,
		ExpiresAt: m.ExpiresAt, CreatedAt: m.CreatedAt,
		Message: "This capsule cannot be opened: it was opened, revoked or expired, or the id is wrong. " + hintGone,
	}
	if !m.Exists {
		return out
	}
	out.Message = "The capsule can be opened with its full link."
	if m.BurnAfterRead {
		out.Message += " Opening it destroys it, so only open it when the recipient is ready."
	}
	return out
}
