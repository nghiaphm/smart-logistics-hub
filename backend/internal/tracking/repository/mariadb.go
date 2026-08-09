package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/tracking/entity"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(ctx context.Context, event *entity.TrackingEvent) error {
	now := time.Now().UTC()
	query := `INSERT INTO tracking_events (order_code, driver_code, status_update, lat, lng, note, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := r.DB.ExecContext(ctx, query,
		event.OrderCode, event.DriverCode, event.StatusUpdate, event.Lat, event.Lng, event.Note, now)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	event.ID = id
	event.Timestamp = now
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*entity.TrackingEvent, error) {
	query := `SELECT id, order_code, driver_code, status_update, lat, lng, note, timestamp
		FROM tracking_events WHERE id = ?`
	event := &entity.TrackingEvent{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&event.ID, &event.OrderCode, &event.DriverCode, &event.StatusUpdate,
		&event.Lat, &event.Lng, &event.Note, &event.Timestamp,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r *Repository) List(ctx context.Context, orderCode, driverCode string, offset, limit int) ([]entity.TrackingEvent, error) {
	var conditions []string
	var args []interface{}
	if orderCode != "" {
		conditions = append(conditions, "order_code = ?")
		args = append(args, orderCode)
	}
	if driverCode != "" {
		conditions = append(conditions, "driver_code = ?")
		args = append(args, driverCode)
	}
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	query := fmt.Sprintf(`SELECT id, order_code, driver_code, status_update, lat, lng, note, timestamp
		FROM tracking_events %s ORDER BY timestamp DESC LIMIT ? OFFSET ?`, whereClause)
	args = append(args, limit, offset)
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []entity.TrackingEvent
	for rows.Next() {
		var event entity.TrackingEvent
		if err := rows.Scan(&event.ID, &event.OrderCode, &event.DriverCode, &event.StatusUpdate,
			&event.Lat, &event.Lng, &event.Note, &event.Timestamp); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if events == nil {
		events = []entity.TrackingEvent{}
	}
	return events, nil
}

func (r *Repository) Count(ctx context.Context, orderCode, driverCode string) (int, error) {
	var conditions []string
	var args []interface{}
	if orderCode != "" {
		conditions = append(conditions, "order_code = ?")
		args = append(args, orderCode)
	}
	if driverCode != "" {
		conditions = append(conditions, "driver_code = ?")
		args = append(args, driverCode)
	}
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	query := fmt.Sprintf(`SELECT COUNT(*) FROM tracking_events %s`, whereClause)
	var count int
	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) GetByOrder(ctx context.Context, orderCode string) ([]entity.TrackingEvent, error) {
	query := `SELECT id, order_code, driver_code, status_update, lat, lng, note, timestamp
		FROM tracking_events WHERE order_code = ? ORDER BY timestamp ASC`
	rows, err := r.DB.QueryContext(ctx, query, orderCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []entity.TrackingEvent
	for rows.Next() {
		var event entity.TrackingEvent
		if err := rows.Scan(&event.ID, &event.OrderCode, &event.DriverCode, &event.StatusUpdate,
			&event.Lat, &event.Lng, &event.Note, &event.Timestamp); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if events == nil {
		events = []entity.TrackingEvent{}
	}
	return events, nil
}

func (r *Repository) Update(ctx context.Context, id int64, event *entity.TrackingEvent) error {
	query := `UPDATE tracking_events SET order_code = ?, driver_code = ?, status_update = ?, lat = ?, lng = ?, note = ?, timestamp = ? WHERE id = ?`
	result, err := r.DB.ExecContext(ctx, query,
		event.OrderCode, event.DriverCode, event.StatusUpdate, event.Lat, event.Lng, event.Note, event.Timestamp, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	event.ID = id
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tracking_events WHERE id = ?`
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
