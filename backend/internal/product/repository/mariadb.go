package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/product/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, p *entity.Product) error {
	now := time.Now().UTC()
	query := `INSERT INTO products (sku, name, category, price, weight_gram, length_cm, width_cm, height_cm, created_at, updated_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		p.Sku, p.Name, p.Category, p.Price, p.WeightGram,
		p.LengthCm, p.WidthCm, p.HeightCm, now, now, p.CreatedBy)
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	id, _ := result.LastInsertId()
	p.ID = id
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*entity.Product, error) {
	query := `SELECT id, sku, name, category, price, weight_gram, length_cm, width_cm, height_cm, created_at, updated_at, created_by
		FROM products WHERE id = ?`
	p := &entity.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Sku, &p.Name, &p.Category, &p.Price, &p.WeightGram,
		&p.LengthCm, &p.WidthCm, &p.HeightCm, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}
	return p, nil
}

func (r *Repository) GetBySku(ctx context.Context, sku string) (*entity.Product, error) {
	query := `SELECT id, sku, name, category, price, weight_gram, length_cm, width_cm, height_cm, created_at, updated_at, created_by
		FROM products WHERE sku = ?`
	p := &entity.Product{}
	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&p.ID, &p.Sku, &p.Name, &p.Category, &p.Price, &p.WeightGram,
		&p.LengthCm, &p.WidthCm, &p.HeightCm, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get product by sku: %w", err)
	}
	return p, nil
}

func (r *Repository) List(ctx context.Context, offset, limit int) ([]entity.Product, error) {
	query := `SELECT id, sku, name, category, price, weight_gram, length_cm, width_cm, height_cm, created_at, updated_at, created_by
		FROM products ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(
			&p.ID, &p.Sku, &p.Name, &p.Category, &p.Price, &p.WeightGram,
			&p.LengthCm, &p.WidthCm, &p.HeightCm, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy,
		); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if products == nil {
		products = []entity.Product{}
	}
	return products, nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM products`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count products: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, p *entity.Product) error {
	now := time.Now().UTC()
	query := `UPDATE products SET sku=?, name=?, category=?, price=?, weight_gram=?, length_cm=?, width_cm=?, height_cm=?, updated_at=?, created_by=?
		WHERE id=?`
	result, err := r.db.ExecContext(ctx, query,
		p.Sku, p.Name, p.Category, p.Price, p.WeightGram,
		p.LengthCm, p.WidthCm, p.HeightCm, now, p.CreatedBy, id)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	p.UpdatedAt = now
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}
