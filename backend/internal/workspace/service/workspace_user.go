package service

import (
	"context"
	"fmt"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	userentity "my-web-app.com/smart-logistic-hub/internal/user/entity"
	"my-web-app.com/smart-logistic-hub/internal/workspace/entity"
)

// WorkspaceUserRepository: thao tác bảng workspace_users (upsert membership).
type WorkspaceUserRepository interface {
	Upsert(ctx context.Context, workspaceID, userID int64, isAdmin bool, actorID *int64) error
	GetByWorkspaceAndUser(ctx context.Context, workspaceID, userID int64) (*entity.WorkspaceUser, error)
	IsWorkspaceAdmin(ctx context.Context, userID, workspaceID int64) (bool, error)
}

// WorkspaceLookup: kiểm tra workspace tồn tại (báo 404 thay vì lỗi FK).
type WorkspaceLookup interface {
	GetByID(ctx context.Context, id int64) (*entity.Workspace, error)
}

// UserLookup: kiểm tra user mục tiêu tồn tại + resolve actor (Super Admin)
// từ keycloak_sub (JWT sub) để ghi created_by/updated_by.
type UserLookup interface {
	GetByID(ctx context.Context, id int64) (*userentity.User, error)
	GetByKeycloakSub(ctx context.Context, sub string) (*userentity.User, error)
}

type MembershipService struct {
	members   WorkspaceUserRepository
	workspace WorkspaceLookup
	users     UserLookup
}

func NewMembershipService(members WorkspaceUserRepository, workspace WorkspaceLookup, users UserLookup) *MembershipService {
	return &MembershipService{members: members, workspace: workspace, users: users}
}

// SetIsAdmin gán/thu hồi quyền admin của 1 user trong 1 workspace (upsert
// membership). actorSub là JWT "sub" của Super Admin đang thao tác — nếu
// resolve được users.id thì ghi created_by/updated_by, không thì NULL.
func (s *MembershipService) SetIsAdmin(ctx context.Context, workspaceID, userID int64, isAdmin bool, actorSub string) (*entity.WorkspaceUser, error) {
	if _, err := s.workspace.GetByID(ctx, workspaceID); err != nil {
		if err == apierrors.ErrNotFound {
			return nil, fmt.Errorf("%w: workspace %d does not exist", apierrors.ErrNotFound, workspaceID)
		}
		return nil, err
	}
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		if err == apierrors.ErrNotFound {
			return nil, fmt.Errorf("%w: user %d does not exist", apierrors.ErrNotFound, userID)
		}
		return nil, err
	}

	var actorID *int64
	if actorSub != "" {
		if actor, err := s.users.GetByKeycloakSub(ctx, actorSub); err == nil {
			actorID = &actor.ID
		} else if err != apierrors.ErrNotFound {
			return nil, err
		}
	}

	if err := s.members.Upsert(ctx, workspaceID, userID, isAdmin, actorID); err != nil {
		return nil, err
	}
	return s.members.GetByWorkspaceAndUser(ctx, workspaceID, userID)
}

// IsWorkspaceAdmin — hàm kiểm tra dùng cho middleware "Workspace Admin" sau
// này. User không có membership / is_active=false → false.
func (s *MembershipService) IsWorkspaceAdmin(ctx context.Context, userID, workspaceID int64) (bool, error) {
	return s.members.IsWorkspaceAdmin(ctx, userID, workspaceID)
}
