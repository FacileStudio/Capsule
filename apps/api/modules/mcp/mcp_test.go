package mcp_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/FacileStudio/tronc/mcp/mcptest"
)

type statusOut struct {
	ID            string `json:"id"`
	Available     bool   `json:"available"`
	BurnAfterRead bool   `json:"burn_after_read"`
	CreatedAt     string `json:"created_at"`
}

type revokeOut struct {
	Status     string `json:"status"`
	ConfirmURL string `json:"confirm_url"`
}

func TestMCPConformance(t *testing.T) {
	f := newFixture(t)
	mcptest.Run(t, mcptest.Config{
		URL:     f.resource,
		Token:   f.token,
		Confirm: []mcptest.ConfirmCase{{Tool: "capsule_revoke", Setup: f.freshCapsule}},
	})
}

func (f *fixture) freshCapsule(tb testing.TB) (map[string]any, func() int) {
	created := f.seal(tb, true)
	return map[string]any{"id": created.ID, "delete_token": created.DeleteToken}, func() int {
		if f.burned(created.ID) {
			return 1
		}
		return 0
	}
}

func TestStatusReadsMetadataAndNeverContent(t *testing.T) {
	f := newFixture(t)
	session := f.connect(t, "other")
	created := f.seal(t, true)
	var out statusOut
	text := call(t, session, "capsule_get_status", map[string]any{"id": created.ID}, &out)
	if !out.Available || !out.BurnAfterRead || out.CreatedAt == "" {
		t.Fatalf("a fresh capsule must be available and burn after reading, got %+v", out)
	}
	if strings.Contains(text, ciphertext) || strings.Contains(text, created.DeleteToken) {
		t.Fatalf("the status must carry neither the ciphertext nor the delete token: %s", text)
	}
	if f.burned(created.ID) {
		t.Fatal("checking a capsule must not open it")
	}
	link := "https://capsule.facile.studio/" + created.ID + "#key"
	if text := refused(t, session, "capsule_get_status", map[string]any{"id": link}); !strings.Contains(text, "#key") {
		t.Fatalf("a full link must be refused with guidance, got %q", text)
	}
	must(t, f.pastes.Burn(created.ID))
	call(t, session, "capsule_get_status", map[string]any{"id": created.ID}, &out)
	if out.Available {
		t.Fatalf("a revoked capsule must not be available, got %+v", out)
	}
}

func TestRevokeRefusesAWrongToken(t *testing.T) {
	f := newFixture(t)
	session := f.connect(t, "owner")
	created := f.seal(t, false)
	for _, token := range []string{"", strings.Repeat("0", 64)} {
		text := refused(t, session, "capsule_revoke", map[string]any{"id": created.ID, "delete_token": token})
		if !strings.Contains(text, "does not match") {
			t.Fatalf("a wrong token must say so, got %q", text)
		}
	}
	if code, body := f.page(t, http.MethodGet, "", "owner", nil); code != http.StatusOK || strings.Contains(body, created.ID) {
		t.Fatalf("a wrong token must open no pending action, got %d %s", code, body)
	}
	if f.burned(created.ID) {
		t.Fatal("a wrong token must burn nothing")
	}
}

func TestConfirmPageRevokesOnceForItsOwnerOnly(t *testing.T) {
	f := newFixture(t)
	created := f.seal(t, false)
	var out revokeOut
	call(t, f.connect(t, "owner"), "capsule_revoke", map[string]any{"id": created.ID, "delete_token": created.DeleteToken}, &out)
	if out.Status != "pending_confirmation" || out.ConfirmURL == "" || f.burned(created.ID) {
		t.Fatalf("a client that cannot ask must get a link and burn nothing, got %+v", out)
	}
	id := actionID(out.ConfirmURL)
	if _, body := f.page(t, http.MethodGet, "", "other", nil); strings.Contains(body, id) {
		t.Fatalf("another person must not see the action: %s", body)
	}
	if _, body := f.page(t, http.MethodGet, "", "owner", nil); !strings.Contains(body, `"capsule":{"id":"`+created.ID) || strings.Contains(body, created.DeleteToken) {
		t.Fatalf("the owner must see the capsule's metadata and never its token: %s", body)
	}
	if code, _ := f.page(t, http.MethodPost, "/"+id+"/confirm", "owner", map[string]string{"Authorization": "Bearer cli"}); code != http.StatusForbidden {
		t.Fatalf("a bearer token must not confirm, got %d", code)
	}
	if code, _ := f.page(t, http.MethodPost, "/"+id+"/confirm", "other", nil); code != http.StatusNotFound || f.burned(created.ID) {
		t.Fatalf("another person must not confirm, got %d", code)
	}
	if code, body := f.page(t, http.MethodPost, "/"+id+"/confirm", "owner", nil); code != http.StatusOK || !strings.Contains(body, "revoked") || !f.burned(created.ID) {
		t.Fatalf("the owner's confirmation must burn the capsule, got %d %s", code, body)
	}
	if code, _ := f.page(t, http.MethodPost, "/"+id+"/confirm", "owner", nil); code != http.StatusNotFound {
		t.Fatalf("a confirmation runs once, got %d", code)
	}
}

func TestOpenedCapsuleLeavesNothingToConfirm(t *testing.T) {
	f := newFixture(t)
	created := f.seal(t, true)
	var out revokeOut
	call(t, f.connect(t, "owner"), "capsule_revoke", map[string]any{"id": created.ID, "delete_token": created.DeleteToken}, &out)
	if _, err := f.pastes.GetContent(created.ID); err != nil {
		t.Fatal(err)
	}
	if _, body := f.page(t, http.MethodGet, "", "owner", nil); !strings.Contains(body, `"state":"missing"`) || strings.Contains(body, `"capsule"`) {
		t.Fatalf("an opened capsule must show as missing, got %s", body)
	}
	if code, _ := f.page(t, http.MethodPost, "/"+actionID(out.ConfirmURL)+"/confirm", "owner", nil); code == http.StatusOK {
		t.Fatal("an opened capsule cannot be confirmed")
	}
}
