package cleanup

import (
	"context"
	"log/slog"
	"time"

	"github.com/FacileStudio/Capsule/apps/api/schemas"

	"gorm.io/gorm"
)

func Start(ctx context.Context, db *gorm.DB, logger *slog.Logger) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("cleanup goroutine stopped")
			return
		case <-ticker.C:
			run(db, logger)
		}
	}
}

func run(db *gorm.DB, logger *slog.Logger) {
	now := time.Now().UTC()

	expired := db.Model(&schemas.Paste{}).
		Where("expires_at < ? AND burned = false", now).
		Updates(map[string]any{"burned": true, "content": ""})
	if expired.Error != nil {
		logger.Error("cleanup: failed to burn expired pastes", slog.Any("error", expired.Error))
	} else if expired.RowsAffected > 0 {
		logger.Info("cleanup: burned expired pastes", slog.Int64("count", expired.RowsAffected))
	}

	cutoff := now.Add(-30 * 24 * time.Hour)
	purged := db.Where("burned = true AND created_at < ?", cutoff).Delete(&schemas.Paste{})
	if purged.Error != nil {
		logger.Error("cleanup: failed to purge old burned pastes", slog.Any("error", purged.Error))
	} else if purged.RowsAffected > 0 {
		logger.Info("cleanup: purged old burned pastes", slog.Int64("count", purged.RowsAffected))
	}
}
