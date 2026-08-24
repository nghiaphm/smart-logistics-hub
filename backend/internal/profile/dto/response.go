package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/profile/entity"
)

type ProfileResponse struct {
	UserSub     string    `json:"user_sub"`
	DisplayName string    `json:"display_name"`
	Phone       string    `json:"phone"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToResponse(p *entity.Profile) ProfileResponse {
	return ProfileResponse{
		UserSub:     p.UserSub,
		DisplayName: p.DisplayName,
		Phone:       p.Phone,
		AvatarURL:   p.AvatarURL,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
