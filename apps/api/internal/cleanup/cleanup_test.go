package cleanup

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/FacileStudio/Capsule/apps/api/schemas"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.AutoMigrate(&schemas.Paste{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestRunBurnsExpiredPastes(t *testing.T) {
	db := setupTestDB(t)

	past := time.Now().UTC().Add(-time.Hour)
	future := time.Now().UTC().Add(time.Hour)

	db.Create(&schemas.Paste{
		ID:          "cap_expired_1",
		Content:     "should be burned",
		ExpiresAt:   &past,
		DeleteToken: "tok_expired_1",
		CreatedAt:   time.Now().UTC(),
	})
	db.Create(&schemas.Paste{
		ID:          "cap_valid_1",
		Content:     "should survive",
		ExpiresAt:   &future,
		DeleteToken: "tok_valid_1",
		CreatedAt:   time.Now().UTC(),
	})
	db.Create(&schemas.Paste{
		ID:          "cap_no_expiry",
		Content:     "no expiry",
		DeleteToken: "tok_no_expiry",
		CreatedAt:   time.Now().UTC(),
	})

	run(db, testLogger())

	var expired schemas.Paste
	db.First(&expired, "id = ?", "cap_expired_1")
	if !expired.Burned {
		t.Error("expired paste should be burned")
	}
	if expired.Content != "" {
		t.Error("expired paste content should be cleared")
	}

	var valid schemas.Paste
	db.First(&valid, "id = ?", "cap_valid_1")
	if valid.Burned {
		t.Error("non-expired paste should not be burned")
	}
	if valid.Content != "should survive" {
		t.Error("non-expired paste content should be intact")
	}

	var noExpiry schemas.Paste
	db.First(&noExpiry, "id = ?", "cap_no_expiry")
	if noExpiry.Burned {
		t.Error("paste without expiry should not be burned")
	}
}

func TestRunPurgesOldBurnedPastes(t *testing.T) {
	db := setupTestDB(t)

	oldTime := time.Now().UTC().Add(-31 * 24 * time.Hour)
	recentTime := time.Now().UTC().Add(-1 * 24 * time.Hour)

	db.Create(&schemas.Paste{
		ID:          "cap_old_burned",
		Content:     "",
		Burned:      true,
		DeleteToken: "tok_old_burned",
		CreatedAt:   oldTime,
	})
	db.Create(&schemas.Paste{
		ID:          "cap_recent_burned",
		Content:     "",
		Burned:      true,
		DeleteToken: "tok_recent_burned",
		CreatedAt:   recentTime,
	})
	db.Create(&schemas.Paste{
		ID:          "cap_old_active",
		Content:     "still active",
		Burned:      false,
		DeleteToken: "tok_old_active",
		CreatedAt:   oldTime,
	})

	run(db, testLogger())

	var count int64
	db.Model(&schemas.Paste{}).Where("id = ?", "cap_old_burned").Count(&count)
	if count != 0 {
		t.Error("old burned paste should be purged")
	}

	db.Model(&schemas.Paste{}).Where("id = ?", "cap_recent_burned").Count(&count)
	if count != 1 {
		t.Error("recently burned paste should not be purged")
	}

	db.Model(&schemas.Paste{}).Where("id = ?", "cap_old_active").Count(&count)
	if count != 1 {
		t.Error("old active paste should not be purged")
	}
}

func TestRunIdempotent(t *testing.T) {
	db := setupTestDB(t)

	past := time.Now().UTC().Add(-time.Hour)
	db.Create(&schemas.Paste{
		ID:          "cap_idem_1",
		Content:     "expire me",
		ExpiresAt:   &past,
		DeleteToken: "tok_idem_1",
		CreatedAt:   time.Now().UTC(),
	})

	run(db, testLogger())
	run(db, testLogger())

	var paste schemas.Paste
	db.First(&paste, "id = ?", "cap_idem_1")
	if !paste.Burned {
		t.Error("paste should be burned")
	}
}
