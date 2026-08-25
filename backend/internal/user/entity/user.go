package entity

import "time"

type User struct {
	ID          int64     `json:"id" db:"id"`
	KeycloakSub string    `json:"keycloak_sub" db:"keycloak_sub"`
	Username    string    `json:"username" db:"username"`
	FullName    string    `json:"full_name" db:"full_name"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	Role        string    `json:"role" db:"role"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
}
