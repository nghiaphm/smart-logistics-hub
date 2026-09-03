package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/workspace/entity"
)

const workspaceUserColumns = `id, workspace_id, user_id, is_admin, is_active, created_by, updated_by, created_at, updated_at`

type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

// Upsert tạo hoặc cập nhật dòng membership (workspace_id, user_id) với
// is_admin mới. INSERT ghi created_by (người gán đầu tiên); UPDATE sau đó chỉ
// đổi is_admin + updated_by + updated_at (giữ nguyên created_by).
func (r *MemberRepository) Upsert(ctx context.Context, workspaceID, userID int64, isAdmin bool, actorID *int64) error {
	now := time.Now().UTC()
	query := `INSERT INTO workspace_users (workspace_id, user_id, is_admin, created_by, updated_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE is_admin = VALUES(is_admin), updated_by = VALUES(updated_by), updated_at = VALUES(updated_at)`
	_, err := r.db.ExecContext(ctx, query, workspaceID, userID, isAdmin, actorID, actorID, now, now)
	if err != nil {
		return fmt.Errorf("upsert workspace_user: %w", err)
	}
	return nil
}

func (r *MemberRepository) GetByWorkspaceAndUser(ctx context.Context, workspaceID, userID int64) (*entity.WorkspaceUser, error) {
	query := `SELECT ` + workspaceUserColumns + ` FROM workspace_users WHERE workspace_id = ? AND user_id = ?`
	wu := &entity.WorkspaceUser{}
	var createdBy, updatedBy sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, workspaceID, userID).Scan(
		&wu.ID, &wu.WorkspaceID, &wu.UserID, &wu.IsAdmin, &wu.IsActive, &createdBy, &updatedBy, &wu.CreatedAt, &wu.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get workspace_user: %w", err)
	}
	if createdBy.Valid {
		wu.CreatedBy = &createdBy.Int64
	}
	if updatedBy.Valid {
		wu.UpdatedBy = &updatedBy.Int64
	}
	return wu, nil
}

// IsWorkspaceAdmin trả true nếu user là admin của workspace và membership
// đang is_active. User không có dòng membership → false (không phải lỗi).
func (r *MemberRepository) IsWorkspaceAdmin(ctx context.Context, userID, workspaceID int64) (bool, error) {
	query := `SELECT is_admin FROM workspace_users WHERE workspace_id = ? AND user_id = ? AND is_active = TRUE`
	var isAdmin bool
	err := r.db.QueryRowContext(ctx, query, workspaceID, userID).Scan(&isAdmin)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check workspace admin: %w", err)
	}
	return isAdmin, nil
}
