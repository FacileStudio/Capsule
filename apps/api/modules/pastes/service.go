package pastes

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"strings"
	"time"

	"github.com/FacileStudio/Capsule/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

type Service struct {
	db           *gorm.DB
	maxPasteSize int
}

func NewService(db *gorm.DB, maxPasteSize int) *Service {
	return &Service{db: db, maxPasteSize: maxPasteSize}
}

func (s *Service) Create(req CreateRequest) (*CreateResponse, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.Invalid("content is required")
	}
	if len(content) > s.maxPasteSize {
		return nil, errors.TooLarge("content exceeds maximum paste size")
	}

	id, err := generateID()
	if err != nil {
		return nil, errors.Internal("failed to generate paste id", err)
	}

	deleteToken, err := generateDeleteToken()
	if err != nil {
		return nil, errors.Internal("failed to generate delete token", err)
	}

	now := time.Now().UTC()

	burnAfterRead := true
	if req.BurnAfterRead != nil {
		burnAfterRead = *req.BurnAfterRead
	}

	var expiresAt *time.Time
	if req.ExpiresIn != "" {
		d, parseErr := parseExpiresIn(req.ExpiresIn)
		if parseErr != nil {
			return nil, parseErr
		}
		t := now.Add(d)
		expiresAt = &t
	}

	paste := schemas.Paste{
		ID:            id,
		Content:       content,
		BurnAfterRead: burnAfterRead,
		ExpiresAt:     expiresAt,
		MaxViews:      req.MaxViews,
		HasPassword:   req.HasPassword,
		Syntax:        req.Syntax,
		DeleteToken:   deleteToken,
		CreatedAt:     now,
	}

	if err := s.db.Create(&paste).Error; err != nil {
		return nil, errors.Internal("failed to create paste", err)
	}

	return &CreateResponse{
		ID:          paste.ID,
		DeleteToken: paste.DeleteToken,
		ExpiresAt:   paste.ExpiresAt,
		CreatedAt:   paste.CreatedAt,
	}, nil
}

func (s *Service) GetMeta(id string) (*MetaResponse, error) {
	var paste schemas.Paste
	if err := s.db.Where("id = ?", id).First(&paste).Error; err != nil {
		return &MetaResponse{ID: id, Exists: false}, nil
	}

	if paste.Burned || isExpired(paste.ExpiresAt) {
		return &MetaResponse{ID: id, Exists: false}, nil
	}

	return &MetaResponse{
		ID:            paste.ID,
		Exists:        true,
		Burned:        paste.Burned,
		HasPassword:   paste.HasPassword,
		Syntax:        paste.Syntax,
		ExpiresAt:     paste.ExpiresAt,
		CreatedAt:     &paste.CreatedAt,
		BurnAfterRead: paste.BurnAfterRead,
	}, nil
}

func (s *Service) GetContent(id string) (*ContentResponse, error) {
	var content string

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var paste schemas.Paste
		if err := tx.Raw("SELECT * FROM pastes WHERE id = ? FOR UPDATE", id).Scan(&paste).Error; err != nil || paste.ID == "" {
			return errors.NotFound("paste not found")
		}

		if paste.Burned || isExpired(paste.ExpiresAt) {
			return errors.NotFound("paste not found")
		}

		content = paste.Content

		updates := map[string]any{
			"view_count": gorm.Expr("view_count + 1"),
		}

		shouldBurn := paste.BurnAfterRead
		if !shouldBurn && paste.MaxViews != nil && paste.ViewCount+1 >= *paste.MaxViews {
			shouldBurn = true
		}

		if shouldBurn {
			updates["burned"] = true
			updates["content"] = ""
		}

		return tx.Model(&paste).Updates(updates).Error
	})

	if err != nil {
		return nil, err
	}

	return &ContentResponse{Content: content}, nil
}

func (s *Service) Revoke(id string, token string) error {
	var paste schemas.Paste
	if err := s.db.Where("id = ?", id).First(&paste).Error; err != nil {
		return errors.NotFound("paste not found")
	}

	if subtle.ConstantTimeCompare([]byte(paste.DeleteToken), []byte(token)) != 1 {
		return errors.Forbidden("invalid delete token")
	}

	if err := s.db.Model(&paste).Updates(map[string]any{
		"burned":  true,
		"content": "",
	}).Error; err != nil {
		return errors.Internal("failed to revoke paste", err)
	}

	return nil
}

func generateID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "cap_" + hex.EncodeToString(b), nil
}

func generateDeleteToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func parseExpiresIn(s string) (time.Duration, error) {
	switch s {
	case "1h":
		return time.Hour, nil
	case "24h":
		return 24 * time.Hour, nil
	case "7d":
		return 7 * 24 * time.Hour, nil
	case "30d":
		return 30 * 24 * time.Hour, nil
	default:
		return 0, errors.Invalid("expires_in must be one of: 1h, 24h, 7d, 30d")
	}
}

func isExpired(expiresAt *time.Time) bool {
	if expiresAt == nil {
		return false
	}
	return time.Now().UTC().After(*expiresAt)
}
