package inventory

import (
	"context"
	"database/sql"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(ctx context.Context, inv *Inventory) error {
	now := time.Now().UTC()
	query := `INSERT INTO inventory (product_id, warehouse_id, available_qty, reserved_qty, damaged_qty, hold_qty, created_at, updated_at, updated_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.DB.ExecContext(ctx, query,
		inv.ProductID, inv.WarehouseID, inv.AvailableQty, inv.ReservedQty, inv.DamagedQty, inv.HoldQty,
		now, now, inv.UpdatedBy)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	inv.ID = id
	inv.CreatedAt = now
	inv.UpdatedAt = now
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Inventory, error) {
	query := `SELECT id, product_id, warehouse_id, available_qty, reserved_qty, damaged_qty, hold_qty, created_at, updated_at, updated_by
		FROM inventory WHERE id = ?`
	inv := &Inventory{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&inv.ID, &inv.ProductID, &inv.WarehouseID,
		&inv.AvailableQty, &inv.ReservedQty, &inv.DamagedQty, &inv.HoldQty,
		&inv.CreatedAt, &inv.UpdatedAt, &inv.UpdatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *Repository) GetByProductWarehouse(ctx context.Context, productID, warehouseID int64) (*Inventory, error) {
	query := `SELECT id, product_id, warehouse_id, available_qty, reserved_qty, damaged_qty, hold_qty, created_at, updated_at, updated_by
		FROM inventory WHERE product_id = ? AND warehouse_id = ?`
	inv := &Inventory{}
	err := r.DB.QueryRowContext(ctx, query, productID, warehouseID).Scan(
		&inv.ID, &inv.ProductID, &inv.WarehouseID,
		&inv.AvailableQty, &inv.ReservedQty, &inv.DamagedQty, &inv.HoldQty,
		&inv.CreatedAt, &inv.UpdatedAt, &inv.UpdatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *Repository) List(ctx context.Context, offset, limit int) ([]Inventory, error) {
	query := `SELECT id, product_id, warehouse_id, available_qty, reserved_qty, damaged_qty, hold_qty, created_at, updated_at, updated_by
		FROM inventory ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []Inventory
	for rows.Next() {
		var inv Inventory
		if err := rows.Scan(
			&inv.ID, &inv.ProductID, &inv.WarehouseID,
			&inv.AvailableQty, &inv.ReservedQty, &inv.DamagedQty, &inv.HoldQty,
			&inv.CreatedAt, &inv.UpdatedAt, &inv.UpdatedBy,
		); err != nil {
			return nil, err
		}
		inventories = append(inventories, inv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if inventories == nil {
		inventories = []Inventory{}
	}
	return inventories, nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM inventory`
	var count int
	err := r.DB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, inv *Inventory) error {
	now := time.Now().UTC()
	query := `UPDATE inventory SET available_qty = ?, reserved_qty = ?, damaged_qty = ?, hold_qty = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	result, err := r.DB.ExecContext(ctx, query,
		inv.AvailableQty, inv.ReservedQty, inv.DamagedQty, inv.HoldQty,
		now, inv.UpdatedBy, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	inv.ID = id
	inv.UpdatedAt = now
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM inventory WHERE id = ?`
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}
