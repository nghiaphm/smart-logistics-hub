package repository

import (
	"context"
	"database/sql"
	"fmt"

	"my-web-app.com/smart-logistic-hub/internal/role/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserRolesAndPermissions(ctx context.Context, keycloakUserID string) ([]entity.Role, map[int64][]entity.Permission, error) {
	// First fetch all roles belonging to this user
	queryRoles := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.keycloak_user_id = ?
	`
	rows, err := r.db.QueryContext(ctx, queryRoles, keycloakUserID)
	if err != nil {
		return nil, nil, fmt.Errorf("get user roles query: %w", err)
	}
	defer rows.Close()

	var roles []entity.Role
	var roleIDs []interface{}
	for rows.Next() {
		var role entity.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, nil, err
		}
		roles = append(roles, role)
		roleIDs = append(roleIDs, role.ID)
	}

	if len(roles) == 0 {
		return nil, nil, nil
	}

	// Fetch all permissions linked dynamically for these role IDs
	placeholders := ""
	for i := range roleIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "?"
	}

	queryPerms := fmt.Sprintf(`
		SELECT rp.role_id, p.id, p.resource, p.action, p.created_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id IN (%s)
	`, placeholders)

	pRows, err := r.db.QueryContext(ctx, queryPerms, roleIDs...)
	if err != nil {
		return nil, nil, fmt.Errorf("get role permissions query: %w", err)
	}
	defer pRows.Close()

	permsMap := make(map[int64][]entity.Permission)
	for pRows.Next() {
		var roleID int64
		var p entity.Permission
		if err := pRows.Scan(&roleID, &p.ID, &p.Resource, &p.Action, &p.CreatedAt); err != nil {
			return nil, nil, err
		}
		permsMap[roleID] = append(permsMap[roleID], p)
	}

	return roles, permsMap, nil
}
