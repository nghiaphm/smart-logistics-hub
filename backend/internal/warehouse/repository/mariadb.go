package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/warehouse/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, w *entity.Warehouse) error {
	now := time.Now().UTC()
	query := `INSERT INTO warehouses (warehouse_code, name, address, lat, lng, contact_phone, manager_name, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		w.WarehouseCode, w.Name, w.Address, w.Lat, w.Lng,
		w.ContactPhone, w.ManagerName, w.IsActive, now, now)
	if err != nil {
		return fmt.Errorf("create warehouse: %w", err)
	}
	id, _ := result.LastInsertId()
	w.ID = id
	w.CreatedAt = now
	w.UpdatedAt = now
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*entity.Warehouse, error) {
	query := `SELECT id, warehouse_code, name, address, lat, lng, contact_phone, manager_name, is_active, created_at, updated_at
		FROM warehouses WHERE id = ?`
	w := &entity.Warehouse{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&w.ID, &w.WarehouseCode, &w.Name, &w.Address, &w.Lat, &w.Lng,
		&w.ContactPhone, &w.ManagerName, &w.IsActive, &w.CreatedAt, &w.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get warehouse: %w", err)
	}
	return w, nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*entity.Warehouse, error) {
	query := `SELECT id, warehouse_code, name, address, lat, lng, contact_phone, manager_name, is_active, created_at, updated_at
		FROM warehouses WHERE warehouse_code = ?`
	w := &entity.Warehouse{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&w.ID, &w.WarehouseCode, &w.Name, &w.Address, &w.Lat, &w.Lng,
		&w.ContactPhone, &w.ManagerName, &w.IsActive, &w.CreatedAt, &w.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get warehouse by code: %w", err)
	}
	return w, nil
}

func (r *Repository) List(ctx context.Context, offset, limit int) ([]entity.Warehouse, error) {
	query := `SELECT id, warehouse_code, name, address, lat, lng, contact_phone, manager_name, is_active, created_at, updated_at
		FROM warehouses ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list warehouses: %w", err)
	}
	defer rows.Close()

	var warehouses []entity.Warehouse
	for rows.Next() {
		var w entity.Warehouse
		if err := rows.Scan(
			&w.ID, &w.WarehouseCode, &w.Name, &w.Address, &w.Lat, &w.Lng,
			&w.ContactPhone, &w.ManagerName, &w.IsActive, &w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan warehouse: %w", err)
		}
		warehouses = append(warehouses, w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if warehouses == nil {
		warehouses = []entity.Warehouse{}
	}
	return warehouses, nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM warehouses`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count warehouses: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, w *entity.Warehouse) error {
	now := time.Now().UTC()
	query := `UPDATE warehouses SET warehouse_code=?, name=?, address=?, lat=?, lng=?, contact_phone=?, manager_name=?, is_active=?, updated_at=?
		WHERE id=?`
	result, err := r.db.ExecContext(ctx, query,
		w.WarehouseCode, w.Name, w.Address, w.Lat, w.Lng,
		w.ContactPhone, w.ManagerName, w.IsActive, now, id)
	if err != nil {
		return fmt.Errorf("update warehouse: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	w.UpdatedAt = now
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM warehouses WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete warehouse: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}
