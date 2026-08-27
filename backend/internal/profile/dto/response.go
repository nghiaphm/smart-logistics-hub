package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/profile/entity"
)

type ProfileResponse struct {
	ID             int64     `json:"id"`
	KeycloakUserID string    `json:"keycloak_user_id"`
	Name           string    `json:"name"`
	UserSub        string    `json:"user_sub"`     // For frontend compatibility
	DisplayName    string    `json:"display_name"` // For frontend compatibility
	Phone          string    `json:"phone"`
	CreatedAt      time.Time `json:"created_at"`
}

func ToResponse(p *entity.Profile) ProfileResponse {
	return ProfileResponse{
		ID:             p.ID,
		KeycloakUserID: p.KeycloakUserID,
		Name:           p.Name,
		UserSub:        p.KeycloakUserID,
		DisplayName:    p.Name,
		Phone:          p.Phone,
		CreatedAt:      p.CreatedAt,
	}
}
