package schemas

import "time"

// Paste is a single encrypted-at-rest capsule: content plus its burn-after-read
// and view-count controls.
type Paste struct {
	ID      string `gorm:"primaryKey;type:varchar(24)" json:"id"`
	Content string `gorm:"type:text;not null" json:"-"`
	// No `default:` here on purpose. GORM omits a zero-valued field from the
	// INSERT when the tag declares a default, so `default:true` silently turned
	// every BurnAfterRead=false paste into a burning one. The column keeps its
	// default in migrations/, which is where the schema now lives.
	BurnAfterRead bool       `gorm:"not null" json:"burn_after_read"`
	ExpiresAt     *time.Time `json:"expires_at"`
	MaxViews      *int       `json:"max_views"`
	ViewCount     int        `gorm:"not null;default:0" json:"view_count"`
	HasPassword   bool       `gorm:"not null;default:false" json:"has_password"`
	Syntax        string     `gorm:"type:varchar(50)" json:"syntax"`
	DeleteToken   string     `gorm:"type:varchar(64);uniqueIndex" json:"-"`
	Burned        bool       `gorm:"not null;default:false" json:"burned"`
	CreatedAt     time.Time  `json:"created_at"`
}
