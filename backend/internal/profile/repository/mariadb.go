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

const profileColumns = `id, keycloak_user_id, name, phone, created_at`

func (r *Repository) GetByKeycloakUserID(ctx context.Context, keycloakUserID string) (*entity.Profile, error) {
	query := `SELECT ` + profileColumns + ` FROM user_profiles WHERE keycloak_user_id = ?`
	p := &entity.Profile{}
	err := r.db.QueryRowContext(ctx, query, keycloakUserID).Scan(
		&p.ID, &p.KeycloakUserID, &p.Name, &p.Phone, &p.CreatedAt,
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
	query := `INSERT INTO user_profiles (keycloak_user_id, name, phone, created_at) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		p.KeycloakUserID, p.Name, p.Phone, now)
	if err != nil {
		return fmt.Errorf("create profile: %w", err)
	}
	id, _ := result.LastInsertId()
	p.ID = id
	p.CreatedAt = now
	return nil
}

func (r *Repository) Update(ctx context.Context, p *entity.Profile) error {
	query := `UPDATE user_profiles SET name=?, phone=? WHERE keycloak_user_id=?`
	result, err := r.db.ExecContext(ctx, query, p.Name, p.Phone, p.KeycloakUserID)
	if err != nil {
		return fmt.Errorf("update profile: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}
