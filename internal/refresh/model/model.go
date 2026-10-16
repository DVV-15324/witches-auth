package model

import "time"

type RefreshToken struct {
	ID            int64      `gorm:"column:id" json:"id"`
	UserID        int64      `gorm:"column:user_id" json:"user_id"`
	Token         string     `gorm:"column:token" json:"token"`
	JKT           string     `gorm:"column:jkt" json:"jkt"`
	IPAddress     string     `gorm:"column:ip_address" json:"ip_address"`
	UserAgent     string     `gorm:"column:user_agent" json:"user_agent"`
	ExpiresAt     time.Time  `gorm:"column:expires_at" json:"expires_at"`
	Revoked       bool       `gorm:"column:revoked" json:"revoked"`
	RevokedReason string     `gorm:"column:revoked_reason" json:"revoked_reason"`
	Locale        string     `gorm:"column:locale" json:"locale"`
	Timezone      string     `gorm:"column:timezone" json:"timezone"`
	RevokedAt     *time.Time `gorm:"column:revoked_at" json:"revoked_at"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
