package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/workspace/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, w *entity.Workspace) error {
	now := time.Now().UTC()
	query := `INSERT INTO workspaces (workspace_code, name, description, is_active, created_at, updated_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		w.WorkspaceCode, w.Name, w.Description, w.IsActive, now, now, w.CreatedBy)
	if err != nil {
		return fmt.Errorf("create workspace: %w", err)
	}
	id, _ := result.LastInsertId()
	w.ID = id
	w.CreatedAt = now
	w.UpdatedAt = now
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*entity.Workspace, error) {
	query := `SELECT id, workspace_code, name, description, is_active, created_at, updated_at, created_by
		FROM workspaces WHERE id = ?`
	w := &entity.Workspace{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&w.ID, &w.WorkspaceCode, &w.Name, &w.Description, &w.IsActive, &w.CreatedAt, &w.UpdatedAt, &w.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get workspace: %w", err)
	}
	return w, nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*entity.Workspace, error) {
	query := `SELECT id, workspace_code, name, description, is_active, created_at, updated_at, created_by
		FROM workspaces WHERE workspace_code = ?`
	w := &entity.Workspace{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&w.ID, &w.WorkspaceCode, &w.Name, &w.Description, &w.IsActive, &w.CreatedAt, &w.UpdatedAt, &w.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get workspace by code: %w", err)
	}
	return w, nil
}

func (r *Repository) List(ctx context.Context, offset, limit int) ([]entity.Workspace, error) {
	query := `SELECT id, workspace_code, name, description, is_active, created_at, updated_at, created_by
		FROM workspaces ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	defer rows.Close()

	var workspaces []entity.Workspace
	for rows.Next() {
		var w entity.Workspace
		if err := rows.Scan(
			&w.ID, &w.WorkspaceCode, &w.Name, &w.Description, &w.IsActive, &w.CreatedAt, &w.UpdatedAt, &w.CreatedBy,
		); err != nil {
			return nil, fmt.Errorf("scan workspace: %w", err)
		}
		workspaces = append(workspaces, w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if workspaces == nil {
		workspaces = []entity.Workspace{}
	}
	return workspaces, nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM workspaces`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count workspaces: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, w *entity.Workspace) error {
	now := time.Now().UTC()
	query := `UPDATE workspaces SET name=?, description=?, is_active=?, updated_at=?
		WHERE id=?`
	result, err := r.db.ExecContext(ctx, query,
		w.Name, w.Description, w.IsActive, now, id)
	if err != nil {
		return fmt.Errorf("update workspace: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	w.UpdatedAt = now
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM workspaces WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete workspace: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}
