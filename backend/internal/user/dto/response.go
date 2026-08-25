package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/user/entity"
)

type UserResponse struct {
	ID          int64     `json:"id"`
	KeycloakSub string    `json:"keycloak_sub"`
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Role        string    `json:"role"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   string    `json:"created_by"`
}

type PaginatedResponse struct {
	Items []UserResponse `json:"items"`
	Total int            `json:"total"`
	Skip  int            `json:"skip"`
	Limit int            `json:"limit"`
}

func ToResponse(u *entity.User) UserResponse {
	return UserResponse{
		ID:          u.ID,
		KeycloakSub: u.KeycloakSub,
		Username:    u.Username,
		FullName:    u.FullName,
		Email:       u.Email,
		Phone:       u.Phone,
		Role:        u.Role,
		IsActive:    u.IsActive,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		CreatedBy:   u.CreatedBy,
	}
}

func ToResponseList(users []entity.User) []UserResponse {
	res := make([]UserResponse, len(users))
	for i, u := range users {
		res[i] = ToResponse(&u)
	}
	return res
}
