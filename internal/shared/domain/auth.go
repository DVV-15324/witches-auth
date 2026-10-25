package domain

import "time"

type Auth struct {
	ID        int64
	Email     string
	Password  string
	UserId    int64
	Salt      string
	Banned    bool
	AuthType  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Auth) TableName() string {
	return "auths"
}
