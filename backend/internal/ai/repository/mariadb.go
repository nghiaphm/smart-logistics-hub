package repository

import (
	"context"
	"database/sql"
	"fmt"

	"my-web-app.com/smart-logistic-hub/internal/ai/entity"
	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const aiEventColumns = `id, event_code, license_plate, confidence_score, event_type, gate_id, timestamp, matched_driver_id, matched_trip_id, created_at`

func (r *Repository) Create(ctx context.Context, e *entity.AIEvent) error {
	query := `INSERT INTO ai_events (event_code, license_plate, confidence_score, event_type, gate_id, timestamp, matched_driver_id, matched_trip_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		e.EventCode, e.LicensePlate, e.ConfidenceScore, e.EventType, e.GateID, e.Timestamp,
		e.MatchedDriverID, e.MatchedTripID, e.CreatedAt)
	if err != nil {
		return fmt.Errorf("create ai event: %w", err)
	}
	id, _ := result.LastInsertId()
	e.ID = id
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*entity.AIEvent, error) {
	query := `SELECT ` + aiEventColumns + ` FROM ai_events WHERE id = ?`
	e := &entity.AIEvent{}
	var ts sql.NullTime
	var driverID, tripID sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.EventCode, &e.LicensePlate, &e.ConfidenceScore, &e.EventType, &e.GateID,
		&ts, &driverID, &tripID, &e.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get ai event: %w", err)
	}
	if ts.Valid {
		e.Timestamp = ts.Time
	}
	if driverID.Valid {
		e.MatchedDriverID = &driverID.Int64
	}
	if tripID.Valid {
		e.MatchedTripID = &tripID.Int64
	}
	return e, nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*entity.AIEvent, error) {
	query := `SELECT ` + aiEventColumns + ` FROM ai_events WHERE event_code = ?`
	e := &entity.AIEvent{}
	var ts sql.NullTime
	var driverID, tripID sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&e.ID, &e.EventCode, &e.LicensePlate, &e.ConfidenceScore, &e.EventType, &e.GateID,
		&ts, &driverID, &tripID, &e.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get ai event by code: %w", err)
	}
	if ts.Valid {
		e.Timestamp = ts.Time
	}
	if driverID.Valid {
		e.MatchedDriverID = &driverID.Int64
	}
	if tripID.Valid {
		e.MatchedTripID = &tripID.Int64
	}
	return e, nil
}

func (r *Repository) List(ctx context.Context, licensePlate, gateID, eventType string, offset, limit int) ([]entity.AIEvent, error) {
	query := `SELECT ` + aiEventColumns + ` FROM ai_events`
	args := []interface{}{}
	conds := []string{}
	if licensePlate != "" {
		conds = append(conds, "license_plate = ?")
		args = append(args, licensePlate)
	}
	if gateID != "" {
		conds = append(conds, "gate_id = ?")
		args = append(args, gateID)
	}
	if eventType != "" {
		conds = append(conds, "event_type = ?")
		args = append(args, eventType)
	}
	for i, c := range conds {
		if i == 0 {
			query += " WHERE " + c
		} else {
			query += " AND " + c
		}
	}
	query += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list ai events: %w", err)
	}
	defer rows.Close()

	var events []entity.AIEvent
	for rows.Next() {
		var e entity.AIEvent
		var ts sql.NullTime
		var driverID, tripID sql.NullInt64
		if err := rows.Scan(
			&e.ID, &e.EventCode, &e.LicensePlate, &e.ConfidenceScore, &e.EventType, &e.GateID,
			&ts, &driverID, &tripID, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan ai event: %w", err)
		}
		if ts.Valid {
			e.Timestamp = ts.Time
		}
		if driverID.Valid {
			e.MatchedDriverID = &driverID.Int64
		}
		if tripID.Valid {
			e.MatchedTripID = &tripID.Int64
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if events == nil {
		events = []entity.AIEvent{}
	}
	return events, nil
}

func (r *Repository) Count(ctx context.Context, licensePlate, gateID, eventType string) (int, error) {
	query := `SELECT COUNT(*) FROM ai_events`
	args := []interface{}{}
	conds := []string{}
	if licensePlate != "" {
		conds = append(conds, "license_plate = ?")
		args = append(args, licensePlate)
	}
	if gateID != "" {
		conds = append(conds, "gate_id = ?")
		args = append(args, gateID)
	}
	if eventType != "" {
		conds = append(conds, "event_type = ?")
		args = append(args, eventType)
	}
	for i, c := range conds {
		if i == 0 {
			query += " WHERE " + c
		} else {
			query += " AND " + c
		}
	}
	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count ai events: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, e *entity.AIEvent) error {
	query := `UPDATE ai_events SET license_plate=?, confidence_score=?, event_type=?, gate_id=?, timestamp=?, matched_driver_id=?, matched_trip_id=?
		WHERE id=?`
	result, err := r.db.ExecContext(ctx, query,
		e.LicensePlate, e.ConfidenceScore, e.EventType, e.GateID, e.Timestamp,
		e.MatchedDriverID, e.MatchedTripID, id)
	if err != nil {
		return fmt.Errorf("update ai event: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM ai_events WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete ai event: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}
