package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/profile/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const profileColumns = `id, user_sub, display_name, phone, avatar_url, created_at, updated_at`

func (r *Repository) GetByUserSub(ctx context.Context, userSub string) (*entity.Profile, error) {
	query := `SELECT ` + profileColumns + ` FROM profiles WHERE user_sub = ?`
	p := &entity.Profile{}
	err := r.db.QueryRowContext(ctx, query, userSub).Scan(
		&p.ID, &p.UserSub, &p.DisplayName, &p.Phone, &p.AvatarURL, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return p, nil
}

func (r *Repository) Create(ctx context.Context, p *entity.Profile) error {
	now := time.Now().UTC()
	query := `INSERT INTO profiles (user_sub, display_name, phone, avatar_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		p.UserSub, p.DisplayName, p.Phone, p.AvatarURL, now, now)
	if err != nil {
		return fmt.Errorf("create profile: %w", err)
	}
	id, _ := result.LastInsertId()
	p.ID = id
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

func (r *Repository) Update(ctx context.Context, p *entity.Profile) error {
	now := time.Now().UTC()
	query := `UPDATE profiles SET display_name=?, phone=?, avatar_url=?, updated_at=? WHERE user_sub=?`
	result, err := r.db.ExecContext(ctx, query,
		p.DisplayName, p.Phone, p.AvatarURL, now, p.UserSub)
	if err != nil {
		return fmt.Errorf("update profile: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	p.UpdatedAt = now
	return nil
}
