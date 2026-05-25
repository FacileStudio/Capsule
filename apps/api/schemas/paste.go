package schemas

import "time"

type Paste struct {
	ID            string     `gorm:"primaryKey;type:varchar(24)" json:"id"`
	Content       string     `gorm:"type:text;not null" json:"-"`
	BurnAfterRead bool       `gorm:"not null;default:true" json:"burn_after_read"`
	ExpiresAt     *time.Time `json:"expires_at"`
	MaxViews      *int       `json:"max_views"`
	ViewCount     int        `gorm:"not null;default:0" json:"view_count"`
	HasPassword   bool       `gorm:"not null;default:false" json:"has_password"`
	Syntax        string     `gorm:"type:varchar(50)" json:"syntax"`
	DeleteToken   string     `gorm:"type:varchar(64);uniqueIndex" json:"-"`
	Burned        bool       `gorm:"not null;default:false" json:"burned"`
	CreatedAt     time.Time  `json:"created_at"`
}
