package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/vehicle/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, v *entity.Vehicle) error {
	query := `INSERT INTO vehicles (license_plate, type, capacity, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		v.LicensePlate, v.Type, v.Capacity, v.Status, v.CreatedAt, v.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create vehicle: %w", err)
	}
	id, _ := result.LastInsertId()
	v.ID = id
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*entity.Vehicle, error) {
	query := `SELECT id, license_plate, type, capacity, status, created_at, updated_at, created_by
		FROM vehicles WHERE id = ?`
	v := &entity.Vehicle{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&v.ID, &v.LicensePlate, &v.Type, &v.Capacity, &v.Status, &v.CreatedAt, &v.UpdatedAt, &v.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get vehicle: %w", err)
	}
	return v, nil
}

func (r *Repository) List(ctx context.Context, status string, offset, limit int) ([]entity.Vehicle, error) {
	query := `SELECT id, license_plate, type, capacity, status, created_at, updated_at, created_by
		FROM vehicles`
	args := []interface{}{}
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []entity.Vehicle
	for rows.Next() {
		var v entity.Vehicle
		if err := rows.Scan(
			&v.ID, &v.LicensePlate, &v.Type, &v.Capacity, &v.Status, &v.CreatedAt, &v.UpdatedAt, &v.CreatedBy,
		); err != nil {
			return nil, fmt.Errorf("scan vehicle: %w", err)
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, rows.Err()
}

func (r *Repository) Count(ctx context.Context, status string) (int, error) {
	query := "SELECT COUNT(*) FROM vehicles"
	args := []interface{}{}
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count vehicles: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, v *entity.Vehicle) error {
	query := `UPDATE vehicles SET license_plate=?, type=?, capacity=?, status=?, updated_at=?
		WHERE id=?`
	v.UpdatedAt = time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query,
		v.LicensePlate, v.Type, v.Capacity, v.Status, v.UpdatedAt, id,
	)
	if err != nil {
		return fmt.Errorf("update vehicle: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM vehicles WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete vehicle: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}
