package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/workspace/entity"
)

// SetMemberAdminRequest: PUT /workspaces/{id}/members/{user_id}
type SetMemberAdminRequest struct {
	IsAdmin *bool `json:"is_admin"`
}

type WorkspaceUserResponse struct {
	ID          int64     `json:"id"`
	WorkspaceID int64     `json:"workspace_id"`
	UserID      int64     `json:"user_id"`
	IsAdmin     bool      `json:"is_admin"`
	IsActive    bool      `json:"is_active"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToWorkspaceUserResponse(wu *entity.WorkspaceUser) WorkspaceUserResponse {
	return WorkspaceUserResponse{
		ID:          wu.ID,
		WorkspaceID: wu.WorkspaceID,
		UserID:      wu.UserID,
		IsAdmin:     wu.IsAdmin,
		IsActive:    wu.IsActive,
		CreatedBy:   wu.CreatedBy,
		UpdatedBy:   wu.UpdatedBy,
		CreatedAt:   wu.CreatedAt,
		UpdatedAt:   wu.UpdatedAt,
	}
}
