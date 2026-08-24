package entity

import "time"

type Workspace struct {
	ID            int64     `json:"id" db:"id"`
	WorkspaceCode string    `json:"workspace_code" db:"workspace_code"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy     string    `json:"created_by" db:"created_by"`
}
