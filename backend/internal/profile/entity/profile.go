package entity

import "time"

type Profile struct {
	ID          int64     `json:"id" db:"id"`
	UserSub     string    `json:"user_sub" db:"user_sub"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Phone       string    `json:"phone" db:"phone"`
	AvatarURL   string    `json:"avatar_url" db:"avatar_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
