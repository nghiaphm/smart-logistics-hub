package entity

import "time"

type Profile struct {
	ID             int64     `json:"id" db:"id"`
	KeycloakUserID string    `json:"keycloak_user_id" db:"keycloak_user_id"`
	Name           string    `json:"name" db:"name"`
	Phone          string    `json:"phone" db:"phone"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
