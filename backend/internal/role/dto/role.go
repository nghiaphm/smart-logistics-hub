package dto

type RoleResponse struct {
	ID          int64                `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Permissions []PermissionResponse `json:"permissions"`
}

type PermissionResponse struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type UserRolesAndPermissionsResponse struct {
	KeycloakUserID string         `json:"keycloak_user_id"`
	Roles          []RoleResponse `json:"roles"`
}
