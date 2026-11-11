package pastes

import (
	"crypto/subtle"

	"github.com/FacileStudio/Capsule/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

// Revoke burns a paste for whoever holds its delete token.
func (s *Service) Revoke(id string, token string) error {
	var paste schemas.Paste
	if err := s.db.Where("id = ?", id).First(&paste).Error; err != nil {
		return errors.NotFound("paste not found")
	}

	if !TokenMatches(&paste, token) {
		return errors.Forbidden("invalid delete token")
	}

	return s.burn(s.db.Model(&paste))
}

// Live returns a paste that can still be opened: stored, not burned and not
// expired. Anything else is not found, as GetMeta reports it.
func (s *Service) Live(id string) (*schemas.Paste, error) {
	var paste schemas.Paste
	if err := s.db.Where("id = ?", id).First(&paste).Error; err != nil {
		return nil, errors.NotFound("capsule not found, already opened, revoked or expired")
	}
	if paste.Burned || isExpired(paste.ExpiresAt) {
		return nil, errors.NotFound("capsule not found, already opened, revoked or expired")
	}
	return &paste, nil
}

// Burn destroys a paste that is not burned yet, as a revoke does. Whoever
// calls it has already checked the delete token.
func (s *Service) Burn(id string) error {
	return s.burn(s.db.Model(&schemas.Paste{}).Where("id = ? AND burned = false", id))
}

// burn marks the pastes query selects as burned and drops their content.
func (s *Service) burn(query *gorm.DB) error {
	if err := query.Updates(map[string]any{"burned": true, "content": ""}).Error; err != nil {
		return errors.Internal("failed to revoke paste", err)
	}
	return nil
}

// TokenMatches reports, in constant time, whether token is the paste's
// delete token. An empty token never matches.
func TokenMatches(paste *schemas.Paste, token string) bool {
	return token != "" && subtle.ConstantTimeCompare([]byte(paste.DeleteToken), []byte(token)) == 1
}
