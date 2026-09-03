package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/user/entity"
)

const userColumns = `id, keycloak_sub, username, full_name, email, phone, role, is_active, created_at, updated_at, created_by`

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *entity.User) error {
	now := time.Now().UTC()
	query := `INSERT INTO users (keycloak_sub, username, full_name, email, phone, role, is_active, created_at, updated_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		u.KeycloakSub, u.Username, u.FullName, u.Email, u.Phone, u.Role, u.IsActive, now, now, u.CreatedBy)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	id, _ := result.LastInsertId()
	u.ID = id
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE id = ?`
	u := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.KeycloakSub, &u.Username, &u.FullName, &u.Email, &u.Phone, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE username = ?`
	u := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.KeycloakSub, &u.Username, &u.FullName, &u.Email, &u.Phone, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return u, nil
}

// GetByKeycloakSub tra user theo keycloak_sub (JWT "sub"). Lưu ý WN-012:
// nhiều user chưa link Keycloak có keycloak_sub = ” — gọi với sub rỗng sẽ
// trả ErrNotFound (không khớp hàng ” nào có ý nghĩa).
func (r *Repository) GetByKeycloakSub(ctx context.Context, sub string) (*entity.User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE keycloak_sub = ?`
	u := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, sub).Scan(
		&u.ID, &u.KeycloakSub, &u.Username, &u.FullName, &u.Email, &u.Phone, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by keycloak sub: %w", err)
	}
	return u, nil
}

func (r *Repository) List(ctx context.Context, offset, limit int) ([]entity.User, error) {
	query := `SELECT ` + userColumns + ` FROM users ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(
			&u.ID, &u.KeycloakSub, &u.Username, &u.FullName, &u.Email, &u.Phone, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.CreatedBy,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if users == nil {
		users = []entity.User{}
	}
	return users, nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, u *entity.User) error {
	now := time.Now().UTC()
	query := `UPDATE users SET keycloak_sub=?, full_name=?, email=?, phone=?, role=?, is_active=?, updated_at=?
		WHERE id=?`
	result, err := r.db.ExecContext(ctx, query,
		u.KeycloakSub, u.FullName, u.Email, u.Phone, u.Role, u.IsActive, now, id)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	u.UpdatedAt = now
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}
