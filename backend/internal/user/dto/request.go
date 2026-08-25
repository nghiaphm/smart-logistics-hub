package dto

type CreateUserRequest struct {
	KeycloakSub *string `json:"keycloak_sub,omitempty"`
	Username    string  `json:"username" binding:"required"`
	FullName    *string `json:"full_name,omitempty"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Role        *string `json:"role,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type UpdateUserRequest struct {
	KeycloakSub *string `json:"keycloak_sub,omitempty"`
	FullName    *string `json:"full_name,omitempty"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Role        *string `json:"role,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
