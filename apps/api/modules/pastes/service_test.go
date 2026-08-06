package pastes

import (
	"github.com/FacileStudio/Capsule/apps/api/internal/testsupport"
	"github.com/FacileStudio/Capsule/apps/api/schemas"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	return testsupport.DB(t)
}

func boolPtr(v bool) *bool { return &v }
func intPtr(v int) *int    { return &v }

func TestGenerateID(t *testing.T) {
	id, err := generateID()
	if err != nil {
		t.Fatalf("generateID: %v", err)
	}
	if !strings.HasPrefix(id, "cap_") {
		t.Errorf("expected prefix cap_, got %s", id)
	}
	if len(id) != 4+16 {
		t.Errorf("expected length 20, got %d (%s)", len(id), id)
	}

	ids := make(map[string]bool)
	for range 100 {
		id, _ := generateID()
		if ids[id] {
			t.Fatalf("duplicate id: %s", id)
		}
		ids[id] = true
	}
}

func TestGenerateDeleteToken(t *testing.T) {
	token, err := generateDeleteToken()
	if err != nil {
		t.Fatalf("generateDeleteToken: %v", err)
	}
	if len(token) != 64 {
		t.Errorf("expected 64 hex chars, got %d (%s)", len(token), token)
	}
}

func TestParseExpiresIn(t *testing.T) {
	tests := []struct {
		input   string
		want    time.Duration
		wantErr bool
	}{
		{"1h", time.Hour, false},
		{"24h", 24 * time.Hour, false},
		{"7d", 7 * 24 * time.Hour, false},
		{"30d", 30 * 24 * time.Hour, false},
		{"5m", 0, true},
		{"", 0, true},
		{"forever", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseExpiresIn(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseExpiresIn(%q) err=%v, wantErr=%v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("parseExpiresIn(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsExpired(t *testing.T) {
	past := time.Now().UTC().Add(-time.Hour)
	future := time.Now().UTC().Add(time.Hour)

	if !isExpired(&past) {
		t.Error("past time should be expired")
	}
	if isExpired(&future) {
		t.Error("future time should not be expired")
	}
	if isExpired(nil) {
		t.Error("nil should not be expired")
	}
}

func TestCreatePaste(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	tests := []struct {
		name    string
		req     CreateRequest
		wantErr bool
	}{
		{
			name:    "empty content",
			req:     CreateRequest{Content: ""},
			wantErr: true,
		},
		{
			name:    "whitespace only",
			req:     CreateRequest{Content: "   "},
			wantErr: true,
		},
		{
			name:    "content too large",
			req:     CreateRequest{Content: strings.Repeat("a", 2000)},
			wantErr: true,
		},
		{
			name: "valid minimal",
			req:  CreateRequest{Content: "hello world"},
		},
		{
			name: "burn after read explicit true",
			req:  CreateRequest{Content: "secret", BurnAfterRead: boolPtr(true)},
		},
		{
			name: "burn after read false",
			req:  CreateRequest{Content: "persistent", BurnAfterRead: boolPtr(false)},
		},
		{
			name: "with expiry",
			req:  CreateRequest{Content: "temp", ExpiresIn: "1h"},
		},
		{
			name: "with max views",
			req:  CreateRequest{Content: "limited", BurnAfterRead: boolPtr(false), MaxViews: intPtr(5)},
		},
		{
			name: "with syntax",
			req:  CreateRequest{Content: "fn main(){}", Syntax: "rust"},
		},
		{
			name: "with password flag",
			req:  CreateRequest{Content: "locked", HasPassword: true},
		},
		{
			name:    "invalid expiry",
			req:     CreateRequest{Content: "nope", ExpiresIn: "999y"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.Create(tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Create() err=%v, wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !strings.HasPrefix(resp.ID, "cap_") {
				t.Errorf("response id should start with cap_, got %s", resp.ID)
			}
			if resp.DeleteToken == "" {
				t.Error("delete token should not be empty")
			}
		})
	}
}

func TestCreateDefaults(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	resp, err := svc.Create(CreateRequest{Content: "test defaults"})
	if err != nil {
		t.Fatal(err)
	}

	var paste schemas.Paste
	db.First(&paste, "id = ?", resp.ID)

	if !paste.BurnAfterRead {
		t.Error("burn_after_read should default to true")
	}
	if paste.ExpiresAt != nil {
		t.Error("expires_at should be nil when no expiry set")
	}
}

func TestCreateWithExpiry(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	resp, err := svc.Create(CreateRequest{Content: "timed", ExpiresIn: "1h"})
	if err != nil {
		t.Fatal(err)
	}

	var paste schemas.Paste
	db.First(&paste, "id = ?", resp.ID)

	if paste.ExpiresAt == nil {
		t.Fatal("expires_at should be set")
	}
	diff := time.Until(*paste.ExpiresAt)
	if diff < 59*time.Minute || diff > 61*time.Minute {
		t.Errorf("expires_at should be ~1h from now, got %v", diff)
	}
}

func TestGetMetaExists(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	resp, _ := svc.Create(CreateRequest{Content: "look at me", Syntax: "go"})

	meta, err := svc.GetMeta(resp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !meta.Exists {
		t.Error("paste should exist")
	}
	if meta.ID != resp.ID {
		t.Errorf("id mismatch: %s vs %s", meta.ID, resp.ID)
	}
	if meta.Syntax != "go" {
		t.Errorf("syntax should be go, got %s", meta.Syntax)
	}
}

func TestGetMetaNotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	meta, err := svc.GetMeta("cap_doesnotexist1234")
	if err != nil {
		t.Fatal(err)
	}
	if meta.Exists {
		t.Error("non-existent paste should not exist")
	}
}

func TestGetMetaBurned(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	resp, _ := svc.Create(CreateRequest{Content: "burn me"})
	db.Model(&schemas.Paste{}).Where("id = ?", resp.ID).Updates(map[string]any{
		"burned": true, "content": "",
	})

	meta, err := svc.GetMeta(resp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Exists {
		t.Error("burned paste should report exists=false")
	}
}

func TestGetMetaExpired(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	resp, _ := svc.Create(CreateRequest{Content: "expired"})
	past := time.Now().UTC().Add(-time.Hour)
	db.Model(&schemas.Paste{}).Where("id = ?", resp.ID).Update("expires_at", past)

	meta, err := svc.GetMeta(resp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Exists {
		t.Error("expired paste should report exists=false")
	}
}

func TestGetContentBurnAfterRead(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	resp, _ := svc.Create(CreateRequest{Content: "one-time secret", BurnAfterRead: boolPtr(true)})

	content, err := svc.GetContent(resp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if content.Content != "one-time secret" {
		t.Errorf("expected 'one-time secret', got %q", content.Content)
	}

	_, err = svc.GetContent(resp.ID)
	if err == nil {
		t.Error("second GetContent on burn-after-read should fail")
	}
}

func TestGetContentNoBurn(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	resp, _ := svc.Create(CreateRequest{Content: "reusable", BurnAfterRead: boolPtr(false)})

	for i := range 3 {
		content, err := svc.GetContent(resp.ID)
		if err != nil {
			t.Fatalf("read %d: %v", i+1, err)
		}
		if content.Content != "reusable" {
			t.Errorf("read %d: expected 'reusable', got %q", i+1, content.Content)
		}
	}

	var paste schemas.Paste
	db.First(&paste, "id = ?", resp.ID)
	if paste.ViewCount != 3 {
		t.Errorf("view_count should be 3, got %d", paste.ViewCount)
	}
}

func TestGetContentMaxViews(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	resp, _ := svc.Create(CreateRequest{
		Content:       "limited views",
		BurnAfterRead: boolPtr(false),
		MaxViews:      intPtr(2),
	})

	content1, err := svc.GetContent(resp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if content1.Content != "limited views" {
		t.Errorf("first read: expected content, got %q", content1.Content)
	}

	content2, err := svc.GetContent(resp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if content2.Content != "limited views" {
		t.Errorf("second read: expected content, got %q", content2.Content)
	}

	_, err = svc.GetContent(resp.ID)
	if err == nil {
		t.Error("third read after max_views=2 should fail")
	}
}

func TestGetContentNotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	_, err := svc.GetContent("cap_doesnotexist1234")
	if err == nil {
		t.Error("GetContent for non-existent paste should fail")
	}
}

func TestGetContentExpired(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	resp, _ := svc.Create(CreateRequest{Content: "expired content", BurnAfterRead: boolPtr(false)})
	past := time.Now().UTC().Add(-time.Hour)
	db.Model(&schemas.Paste{}).Where("id = ?", resp.ID).Update("expires_at", past)

	_, err := svc.GetContent(resp.ID)
	if err == nil {
		t.Error("GetContent for expired paste should fail")
	}
}

func TestRevokeValid(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	resp, _ := svc.Create(CreateRequest{Content: "revocable"})

	if err := svc.Revoke(resp.ID, resp.DeleteToken); err != nil {
		t.Fatalf("Revoke: %v", err)
	}

	var paste schemas.Paste
	db.First(&paste, "id = ?", resp.ID)
	if !paste.Burned {
		t.Error("paste should be burned after revoke")
	}
	if paste.Content != "" {
		t.Error("content should be empty after revoke")
	}
}

func TestRevokeWrongToken(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	resp, _ := svc.Create(CreateRequest{Content: "protected"})

	err := svc.Revoke(resp.ID, "badtoken")
	if err == nil {
		t.Error("Revoke with wrong token should fail")
	}
}

func TestRevokeNotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db, 1024)

	err := svc.Revoke("cap_doesnotexist1234", "whatever")
	if err == nil {
		t.Error("Revoke for non-existent paste should fail")
	}
}
