package domain

import "time"

type User struct {
	ID        int64
	Name      string
	Role      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "users"
}
