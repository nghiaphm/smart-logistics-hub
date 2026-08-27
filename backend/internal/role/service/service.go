package service

import (
	"context"

	"my-web-app.com/smart-logistic-hub/internal/role/dto"
	"my-web-app.com/smart-logistic-hub/internal/role/entity"
)

type RoleRepository interface {
	GetUserRolesAndPermissions(ctx context.Context, keycloakUserID string) ([]entity.Role, map[int64][]entity.Permission, error)
}

type Service struct {
	repo RoleRepository
}

func NewService(repo RoleRepository) *Service {
	return &Service{repo: repo}
}

func NewServiceWithRepo(repo RoleRepository) *Service {
	return NewService(repo)
}

func (s *Service) GetUserRolesAndPermissions(ctx context.Context, keycloakUserID string) (*dto.UserRolesAndPermissionsResponse, error) {
	roles, permsMap, err := s.repo.GetUserRolesAndPermissions(ctx, keycloakUserID)
	if err != nil {
		return nil, err
	}

	resRoles := make([]dto.RoleResponse, 0, len(roles))
	for _, r := range roles {
		rolePerms := permsMap[r.ID]
		resPerms := make([]dto.PermissionResponse, 0, len(rolePerms))
		for _, p := range rolePerms {
			resPerms = append(resPerms, dto.PermissionResponse{
				Resource: p.Resource,
				Action:   p.Action,
			})
		}
		resRoles = append(resRoles, dto.RoleResponse{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Permissions: resPerms,
		})
	}

	return &dto.UserRolesAndPermissionsResponse{
		KeycloakUserID: keycloakUserID,
		Roles:          resRoles,
	}, nil
}
