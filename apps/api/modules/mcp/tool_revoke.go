package mcp

import (
	"context"
	"fmt"

	"github.com/FacileStudio/Capsule/apps/api/modules/pastes"
	"github.com/FacileStudio/Capsule/apps/api/schemas"
	mcpkit "github.com/FacileStudio/tronc/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	statusRevoked = "revoked"
	waiting       = "Nothing happened yet: your client cannot ask the user to confirm. Give the user confirm_url and " +
		"ask them to confirm there, signed in to Capsule. Do not call this tool again for the same capsule; the " +
		"request stays open for 7 days or until the capsule is opened or expires."
)

// revoke is capsule_revoke: the capsule burns once the user confirms. The
// delete token is checked here and never stored: the pending action names
// the capsule by id alone.
func (t *toolset) revoke(ctx context.Context, req *sdk.CallToolRequest, in revokeInput) (*sdk.CallToolResult, revokeOutput, error) {
	if mcpkit.Retrying(req) {
		return t.resumeRevoke(ctx, req)
	}
	if !idPattern.MatchString(in.ID) {
		return nil, revokeOutput{}, errID
	}
	paste, err := t.Pastes.Live(in.ID)
	if err != nil {
		return nil, revokeOutput{}, mcpkit.Failure(err, hintGone)
	}
	if !pastes.TokenMatches(paste, in.DeleteToken) {
		return nil, revokeOutput{}, errToken
	}
	result, outcome, err := t.confirm.Request(ctx, req, kindRevoke, paste.ID, revokeQuestion(paste))
	if err != nil {
		return nil, revokeOutput{}, mcpkit.Failure(err, hintGone)
	}
	if result != nil {
		return result, revokeOutput{}, nil
	}
	return nil, revokeOutput{Status: outcome.Status, ID: paste.ID, ConfirmURL: outcome.ConfirmURL, Message: waiting}, nil
}

// resumeRevoke handles the retry that carries the person's answer.
func (t *toolset) resumeRevoke(ctx context.Context, req *sdk.CallToolRequest) (*sdk.CallToolResult, revokeOutput, error) {
	outcome, err := t.confirm.Resume(ctx, req, kindRevoke)
	if err != nil {
		return nil, revokeOutput{}, mcpkit.Failure(err, "Call capsule_revoke again to ask the user afresh.")
	}
	if outcome.Status == mcpkit.StatusDeclined {
		return nil, revokeOutput{Status: mcpkit.StatusDeclined, Message: "The user declined. Nothing happened, and the capsule can still be opened."}, nil
	}
	return nil, revokeOutput{Status: statusRevoked, ID: outcome.Action.Target, Message: "The user confirmed and the capsule is destroyed: its link no longer opens."}, nil
}

// revokeQuestion is what the person is asked before a capsule burns.
func revokeQuestion(p *schemas.Paste) string {
	return fmt.Sprintf("Revoke capsule %s, sealed %s UTC? Its content is destroyed for good and its link stops working.",
		p.ID, p.CreatedAt.UTC().Format("2006-01-02 15:04"))
}
