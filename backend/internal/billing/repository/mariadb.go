package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"my-web-app.com/smart-logistic-hub/internal/billing/entity"
	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const billingColumns = `id, billing_code, order_code, amount_total, currency, payment_method, payment_status, transaction_id, payer_name, payer_phone, payer_email, created_at, updated_at, paid_at, created_by`

func (r *Repository) Create(ctx context.Context, b *entity.Billing) error {
	now := time.Now().UTC()
	query := `INSERT INTO billing (billing_code, order_code, amount_total, currency, payment_method, payment_status, transaction_id, payer_name, payer_phone, payer_email, created_at, updated_at, paid_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		b.BillingCode, b.OrderCode, b.AmountTotal, b.Currency, b.PaymentMethod, b.PaymentStatus, b.TransactionID,
		b.PayerName, b.PayerPhone, b.PayerEmail, now, now, b.PaidAt, b.CreatedBy)
	if err != nil {
		return fmt.Errorf("create billing: %w", err)
	}
	id, _ := result.LastInsertId()
	b.ID = id
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*entity.Billing, error) {
	query := `SELECT ` + billingColumns + ` FROM billing WHERE id = ?`
	b := &entity.Billing{}
	var paidAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.BillingCode, &b.OrderCode, &b.AmountTotal, &b.Currency, &b.PaymentMethod, &b.PaymentStatus,
		&b.TransactionID, &b.PayerName, &b.PayerPhone, &b.PayerEmail, &b.CreatedAt, &b.UpdatedAt, &paidAt, &b.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get billing: %w", err)
	}
	if paidAt.Valid {
		b.PaidAt = &paidAt.Time
	}
	return b, nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*entity.Billing, error) {
	query := `SELECT ` + billingColumns + ` FROM billing WHERE billing_code = ?`
	b := &entity.Billing{}
	var paidAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&b.ID, &b.BillingCode, &b.OrderCode, &b.AmountTotal, &b.Currency, &b.PaymentMethod, &b.PaymentStatus,
		&b.TransactionID, &b.PayerName, &b.PayerPhone, &b.PayerEmail, &b.CreatedAt, &b.UpdatedAt, &paidAt, &b.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get billing by code: %w", err)
	}
	if paidAt.Valid {
		b.PaidAt = &paidAt.Time
	}
	return b, nil
}

func (r *Repository) GetByOrderCode(ctx context.Context, orderCode string) (*entity.Billing, error) {
	query := `SELECT ` + billingColumns + ` FROM billing WHERE order_code = ? ORDER BY id DESC LIMIT 1`
	b := &entity.Billing{}
	var paidAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, orderCode).Scan(
		&b.ID, &b.BillingCode, &b.OrderCode, &b.AmountTotal, &b.Currency, &b.PaymentMethod, &b.PaymentStatus,
		&b.TransactionID, &b.PayerName, &b.PayerPhone, &b.PayerEmail, &b.CreatedAt, &b.UpdatedAt, &paidAt, &b.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get billing by order code: %w", err)
	}
	if paidAt.Valid {
		b.PaidAt = &paidAt.Time
	}
	return b, nil
}

func (r *Repository) List(ctx context.Context, offset, limit int) ([]entity.Billing, error) {
	query := `SELECT ` + billingColumns + ` FROM billing ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list billings: %w", err)
	}
	defer rows.Close()

	var billings []entity.Billing
	for rows.Next() {
		var b entity.Billing
		var paidAt sql.NullTime
		if err := rows.Scan(
			&b.ID, &b.BillingCode, &b.OrderCode, &b.AmountTotal, &b.Currency, &b.PaymentMethod, &b.PaymentStatus,
			&b.TransactionID, &b.PayerName, &b.PayerPhone, &b.PayerEmail, &b.CreatedAt, &b.UpdatedAt, &paidAt, &b.CreatedBy,
		); err != nil {
			return nil, fmt.Errorf("scan billing: %w", err)
		}
		if paidAt.Valid {
			b.PaidAt = &paidAt.Time
		}
		billings = append(billings, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if billings == nil {
		billings = []entity.Billing{}
	}
	return billings, nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM billing`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count billings: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, b *entity.Billing) error {
	now := time.Now().UTC()
	query := `UPDATE billing SET billing_code=?, order_code=?, amount_total=?, currency=?, payment_method=?, payment_status=?, transaction_id=?, payer_name=?, payer_phone=?, payer_email=?, paid_at=?, updated_at=?
		WHERE id=?`
	result, err := r.db.ExecContext(ctx, query,
		b.BillingCode, b.OrderCode, b.AmountTotal, b.Currency, b.PaymentMethod, b.PaymentStatus, b.TransactionID,
		b.PayerName, b.PayerPhone, b.PayerEmail, b.PaidAt, now, id)
	if err != nil {
		return fmt.Errorf("update billing: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	b.UpdatedAt = now
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM billing WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete billing: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}
