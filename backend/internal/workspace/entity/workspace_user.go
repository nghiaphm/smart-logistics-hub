package entity

import "time"

type WorkspaceUser struct {
	ID          int64     `json:"id" db:"id"`
	WorkspaceID int64     `json:"workspace_id" db:"workspace_id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	IsAdmin     bool      `json:"is_admin" db:"is_admin"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedBy   *int64    `json:"created_by" db:"created_by"`
	UpdatedBy   *int64    `json:"updated_by" db:"updated_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
