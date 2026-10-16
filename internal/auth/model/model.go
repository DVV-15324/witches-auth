package model

import "time"

type Auth struct {
	ID        int64     `gorm:"column:id" json:"id"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	UserId    int64     `gorm:"column:user_id" json:"user_id"`
	Salt      string    `gorm:"column:salt" json:"salt"`
	Banned    bool      `gorm:"column:banned" json:"banned"`
	AuthType  string    `gorm:"column:auth_type" json:"auth_type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Auth) TableName() string {
	return "auths"
}
