package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/workspace/entity"
)

type WorkspaceResponse struct {
	ID            int64     `json:"id"`
	WorkspaceCode string    `json:"workspace_code"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedBy     string    `json:"created_by"`
}

type PaginatedResponse struct {
	Items []WorkspaceResponse `json:"items"`
	Total int                 `json:"total"`
	Skip  int                 `json:"skip"`
	Limit int                 `json:"limit"`
}

func ToResponse(w *entity.Workspace) WorkspaceResponse {
	return WorkspaceResponse{
		ID:            w.ID,
		WorkspaceCode: w.WorkspaceCode,
		Name:          w.Name,
		Description:   w.Description,
		IsActive:      w.IsActive,
		CreatedAt:     w.CreatedAt,
		UpdatedAt:     w.UpdatedAt,
		CreatedBy:     w.CreatedBy,
	}
}

func ToResponseList(workspaces []entity.Workspace) []WorkspaceResponse {
	res := make([]WorkspaceResponse, len(workspaces))
	for i, w := range workspaces {
		res[i] = ToResponse(&w)
	}
	return res
}
