package domain

import "time"

type RefreshToken struct {
	ID            int64
	UserID        int64
	Token         string
	JKT           string
	IPAddress     string
	UserAgent     string
	ExpiresAt     time.Time
	Revoked       bool
	RevokedAt     *time.Time
	RevokedReason string
	Locale        string
	Timezone      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
