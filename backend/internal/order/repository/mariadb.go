package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/order/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, o *entity.Order) error {
	query := `INSERT INTO orders (order_code, sender_name, sender_phone, sender_address, sender_province, sender_district, sender_ward, sender_postal_code,
		receiver_name, receiver_phone, receiver_address, receiver_province, receiver_district, receiver_ward, receiver_postal_code,
		status, assigned_driver_id, created_at, updated_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		o.OrderCode, o.SenderName, o.SenderPhone, o.SenderAddress, o.SenderProvince, o.SenderDistrict, o.SenderWard, o.SenderPostalCode,
		o.ReceiverName, o.ReceiverPhone, o.ReceiverAddress, o.ReceiverProvince, o.ReceiverDistrict, o.ReceiverWard, o.ReceiverPostalCode,
		o.Status, o.AssignedDriverID, o.CreatedAt, o.UpdatedAt, o.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}
	id, _ := result.LastInsertId()
	o.ID = id
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*entity.Order, error) {
	query := `SELECT id, order_code, sender_name, sender_phone, sender_address, sender_province, sender_district, sender_ward, sender_postal_code,
		receiver_name, receiver_phone, receiver_address, receiver_province, receiver_district, receiver_ward, receiver_postal_code,
		status, assigned_driver_id, created_at, updated_at, created_by FROM orders WHERE id = ?`
	o := &entity.Order{}
	var driverID sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&o.ID, &o.OrderCode, &o.SenderName, &o.SenderPhone, &o.SenderAddress, &o.SenderProvince, &o.SenderDistrict, &o.SenderWard, &o.SenderPostalCode,
		&o.ReceiverName, &o.ReceiverPhone, &o.ReceiverAddress, &o.ReceiverProvince, &o.ReceiverDistrict, &o.ReceiverWard, &o.ReceiverPostalCode,
		&o.Status, &driverID, &o.CreatedAt, &o.UpdatedAt, &o.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	if driverID.Valid {
		o.AssignedDriverID = &driverID.Int64
	}
	return o, nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*entity.Order, error) {
	query := `SELECT id, order_code, sender_name, sender_phone, sender_address, sender_province, sender_district, sender_ward, sender_postal_code,
		receiver_name, receiver_phone, receiver_address, receiver_province, receiver_district, receiver_ward, receiver_postal_code,
		status, assigned_driver_id, created_at, updated_at, created_by FROM orders WHERE order_code = ?`
	o := &entity.Order{}
	var driverID sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&o.ID, &o.OrderCode, &o.SenderName, &o.SenderPhone, &o.SenderAddress, &o.SenderProvince, &o.SenderDistrict, &o.SenderWard, &o.SenderPostalCode,
		&o.ReceiverName, &o.ReceiverPhone, &o.ReceiverAddress, &o.ReceiverProvince, &o.ReceiverDistrict, &o.ReceiverWard, &o.ReceiverPostalCode,
		&o.Status, &driverID, &o.CreatedAt, &o.UpdatedAt, &o.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get order by code: %w", err)
	}
	if driverID.Valid {
		o.AssignedDriverID = &driverID.Int64
	}
	return o, nil
}

func (r *Repository) List(ctx context.Context, offset, limit int) ([]entity.Order, error) {
	query := `SELECT id, order_code, sender_name, sender_phone, sender_address, sender_province, sender_district, sender_ward, sender_postal_code,
		receiver_name, receiver_phone, receiver_address, receiver_province, receiver_district, receiver_ward, receiver_postal_code,
		status, assigned_driver_id, created_at, updated_at, created_by
		FROM orders ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var o entity.Order
		var driverID sql.NullInt64
		if err := rows.Scan(
			&o.ID, &o.OrderCode, &o.SenderName, &o.SenderPhone, &o.SenderAddress, &o.SenderProvince, &o.SenderDistrict, &o.SenderWard, &o.SenderPostalCode,
			&o.ReceiverName, &o.ReceiverPhone, &o.ReceiverAddress, &o.ReceiverProvince, &o.ReceiverDistrict, &o.ReceiverWard, &o.ReceiverPostalCode,
			&o.Status, &driverID, &o.CreatedAt, &o.UpdatedAt, &o.CreatedBy,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		if driverID.Valid {
			o.AssignedDriverID = &driverID.Int64
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM orders").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count orders: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, o *entity.Order) error {
	query := `UPDATE orders SET order_code=?, sender_name=?, sender_phone=?, sender_address=?, sender_province=?, sender_district=?, sender_ward=?, sender_postal_code=?,
		receiver_name=?, receiver_phone=?, receiver_address=?, receiver_province=?, receiver_district=?, receiver_ward=?, receiver_postal_code=?,
		status=?, assigned_driver_id=?, updated_at=? WHERE id=?`
	o.UpdatedAt = time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query,
		o.OrderCode, o.SenderName, o.SenderPhone, o.SenderAddress, o.SenderProvince, o.SenderDistrict, o.SenderWard, o.SenderPostalCode,
		o.ReceiverName, o.ReceiverPhone, o.ReceiverAddress, o.ReceiverProvince, o.ReceiverDistrict, o.ReceiverWard, o.ReceiverPostalCode,
		o.Status, o.AssignedDriverID, o.UpdatedAt, id,
	)
	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete order: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}

func (r *Repository) CreateItem(ctx context.Context, item *entity.OrderItem) error {
	query := `INSERT INTO order_items (order_id, product_id, product_name, quantity, weight_gram) VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, item.OrderID, item.ProductID, item.ProductName, item.Quantity, item.WeightGram)
	if err != nil {
		return fmt.Errorf("create order item: %w", err)
	}
	id, _ := result.LastInsertId()
	item.ID = id
	return nil
}

func (r *Repository) GetItemsByOrder(ctx context.Context, orderID int64) ([]entity.OrderItem, error) {
	query := `SELECT id, order_id, product_id, product_name, quantity, weight_gram FROM order_items WHERE order_id = ?`
	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("get order items: %w", err)
	}
	defer rows.Close()

	var items []entity.OrderItem
	for rows.Next() {
		var item entity.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName, &item.Quantity, &item.WeightGram); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) DeleteItems(ctx context.Context, orderID int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		return fmt.Errorf("delete order items: %w", err)
	}
	return nil
}
