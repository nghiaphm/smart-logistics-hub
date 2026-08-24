package dto

type CreateWorkspaceRequest struct {
	WorkspaceCode string  `json:"workspace_code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Description   *string `json:"description,omitempty"`
}

type UpdateWorkspaceRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
