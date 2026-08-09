package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/driver/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, d *entity.Driver) error {
	query := `INSERT INTO drivers (driver_code, full_name, phone, vehicle_type, license_plate, status, current_lat, current_lng, warehouse_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		d.DriverCode, d.FullName, d.Phone, d.VehicleType, d.LicensePlate,
		d.Status, d.CurrentLat, d.CurrentLng, d.WarehouseID, d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create driver: %w", err)
	}
	id, _ := result.LastInsertId()
	d.ID = id
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*entity.Driver, error) {
	query := `SELECT id, driver_code, full_name, phone, vehicle_type, license_plate, status, current_lat, current_lng, warehouse_id, created_at, updated_at, created_by
		FROM drivers WHERE id = ?`
	d := &entity.Driver{}
	var whID sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.DriverCode, &d.FullName, &d.Phone, &d.VehicleType, &d.LicensePlate,
		&d.Status, &d.CurrentLat, &d.CurrentLng, &whID, &d.CreatedAt, &d.UpdatedAt, &d.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get driver: %w", err)
	}
	if whID.Valid {
		d.WarehouseID = &whID.Int64
	}
	return d, nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*entity.Driver, error) {
	query := `SELECT id, driver_code, full_name, phone, vehicle_type, license_plate, status, current_lat, current_lng, warehouse_id, created_at, updated_at, created_by
		FROM drivers WHERE driver_code = ?`
	d := &entity.Driver{}
	var whID sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&d.ID, &d.DriverCode, &d.FullName, &d.Phone, &d.VehicleType, &d.LicensePlate,
		&d.Status, &d.CurrentLat, &d.CurrentLng, &whID, &d.CreatedAt, &d.UpdatedAt, &d.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get driver by code: %w", err)
	}
	if whID.Valid {
		d.WarehouseID = &whID.Int64
	}
	return d, nil
}

func (r *Repository) List(ctx context.Context, status string, offset, limit int) ([]entity.Driver, error) {
	query := `SELECT id, driver_code, full_name, phone, vehicle_type, license_plate, status, current_lat, current_lng, warehouse_id, created_at, updated_at, created_by
		FROM drivers`
	args := []interface{}{}
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list drivers: %w", err)
	}
	defer rows.Close()

	var drivers []entity.Driver
	for rows.Next() {
		var d entity.Driver
		var whID sql.NullInt64
		if err := rows.Scan(
			&d.ID, &d.DriverCode, &d.FullName, &d.Phone, &d.VehicleType, &d.LicensePlate,
			&d.Status, &d.CurrentLat, &d.CurrentLng, &whID, &d.CreatedAt, &d.UpdatedAt, &d.CreatedBy,
		); err != nil {
			return nil, fmt.Errorf("scan driver: %w", err)
		}
		if whID.Valid {
			d.WarehouseID = &whID.Int64
		}
		drivers = append(drivers, d)
	}
	return drivers, rows.Err()
}

func (r *Repository) Count(ctx context.Context, status string) (int, error) {
	query := "SELECT COUNT(*) FROM drivers"
	args := []interface{}{}
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count drivers: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, d *entity.Driver) error {
	query := `UPDATE drivers SET driver_code=?, full_name=?, phone=?, vehicle_type=?, license_plate=?, status=?, current_lat=?, current_lng=?, warehouse_id=?, updated_at=?
		WHERE id=?`
	d.UpdatedAt = time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query,
		d.DriverCode, d.FullName, d.Phone, d.VehicleType, d.LicensePlate,
		d.Status, d.CurrentLat, d.CurrentLng, d.WarehouseID, d.UpdatedAt, id,
	)
	if err != nil {
		return fmt.Errorf("update driver: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM drivers WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete driver: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}
