package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/inbound/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const inboundColumns = `id, receipt_code, supplier_name, warehouse_id, status, created_at, updated_at, completed_at, created_by`

func (r *Repository) Create(ctx context.Context, in *entity.Inbound) error {
	now := time.Now().UTC()
	query := `INSERT INTO inbounds (receipt_code, supplier_name, warehouse_id, status, created_at, updated_at, completed_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		in.ReceiptCode, in.SupplierName, in.WarehouseID, in.Status, now, now, in.CompletedAt, in.CreatedBy)
	if err != nil {
		return fmt.Errorf("create inbound: %w", err)
	}
	id, _ := result.LastInsertId()
	in.ID = id
	in.CreatedAt = now
	in.UpdatedAt = now
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*entity.Inbound, error) {
	query := `SELECT ` + inboundColumns + ` FROM inbounds WHERE id = ?`
	in := &entity.Inbound{}
	var completedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&in.ID, &in.ReceiptCode, &in.SupplierName, &in.WarehouseID, &in.Status,
		&in.CreatedAt, &in.UpdatedAt, &completedAt, &in.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get inbound: %w", err)
	}
	if completedAt.Valid {
		in.CompletedAt = &completedAt.Time
	}
	return in, nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*entity.Inbound, error) {
	query := `SELECT ` + inboundColumns + ` FROM inbounds WHERE receipt_code = ?`
	in := &entity.Inbound{}
	var completedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&in.ID, &in.ReceiptCode, &in.SupplierName, &in.WarehouseID, &in.Status,
		&in.CreatedAt, &in.UpdatedAt, &completedAt, &in.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get inbound by code: %w", err)
	}
	if completedAt.Valid {
		in.CompletedAt = &completedAt.Time
	}
	return in, nil
}

func (r *Repository) List(ctx context.Context, offset, limit int) ([]entity.Inbound, error) {
	query := `SELECT ` + inboundColumns + ` FROM inbounds ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list inbounds: %w", err)
	}
	defer rows.Close()

	var inbounds []entity.Inbound
	for rows.Next() {
		var in entity.Inbound
		var completedAt sql.NullTime
		if err := rows.Scan(
			&in.ID, &in.ReceiptCode, &in.SupplierName, &in.WarehouseID, &in.Status,
			&in.CreatedAt, &in.UpdatedAt, &completedAt, &in.CreatedBy,
		); err != nil {
			return nil, fmt.Errorf("scan inbound: %w", err)
		}
		if completedAt.Valid {
			in.CompletedAt = &completedAt.Time
		}
		inbounds = append(inbounds, in)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if inbounds == nil {
		inbounds = []entity.Inbound{}
	}
	return inbounds, nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM inbounds`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count inbounds: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, in *entity.Inbound) error {
	now := time.Now().UTC()
	query := `UPDATE inbounds SET receipt_code=?, supplier_name=?, warehouse_id=?, status=?, completed_at=?, updated_at=?
		WHERE id=?`
	result, err := r.db.ExecContext(ctx, query,
		in.ReceiptCode, in.SupplierName, in.WarehouseID, in.Status, in.CompletedAt, now, id)
	if err != nil {
		return fmt.Errorf("update inbound: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	in.UpdatedAt = now
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM inbounds WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete inbound: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}

func (r *Repository) CreateItems(ctx context.Context, inboundID int64, items []entity.InboundItem) error {
	query := `INSERT INTO inbound_items (inbound_id, product_id, expected_qty, received_qty, rejected_qty, qc_passed)
		VALUES (?, ?, ?, ?, ?, ?)`
	for i := range items {
		item := &items[i]
		result, err := r.db.ExecContext(ctx, query,
			inboundID, item.ProductID, item.ExpectedQty, item.ReceivedQty, item.RejectedQty, item.QcPassed)
		if err != nil {
			return fmt.Errorf("create inbound item: %w", err)
		}
		id, _ := result.LastInsertId()
		item.ID = id
		item.InboundID = inboundID
	}
	return nil
}

func (r *Repository) GetItems(ctx context.Context, inboundID int64) ([]entity.InboundItem, error) {
	query := `SELECT id, inbound_id, product_id, expected_qty, received_qty, rejected_qty, qc_passed
		FROM inbound_items WHERE inbound_id = ? ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query, inboundID)
	if err != nil {
		return nil, fmt.Errorf("get inbound items: %w", err)
	}
	defer rows.Close()

	var items []entity.InboundItem
	for rows.Next() {
		var item entity.InboundItem
		if err := rows.Scan(
			&item.ID, &item.InboundID, &item.ProductID, &item.ExpectedQty,
			&item.ReceivedQty, &item.RejectedQty, &item.QcPassed,
		); err != nil {
			return nil, fmt.Errorf("scan inbound item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if items == nil {
		items = []entity.InboundItem{}
	}
	return items, nil
}

func (r *Repository) DeleteItems(ctx context.Context, inboundID int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM inbound_items WHERE inbound_id = ?", inboundID)
	if err != nil {
		return fmt.Errorf("delete inbound items: %w", err)
	}
	return nil
}
