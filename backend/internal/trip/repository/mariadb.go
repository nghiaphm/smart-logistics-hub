package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/trip/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const tripColumns = `id, trip_code, driver_id, vehicle_license_plate, status, total_distance_km, estimated_duration_min, actual_start_at, actual_end_at, created_at, updated_at, created_by`

func (r *Repository) Create(ctx context.Context, t *entity.Trip) error {
	now := time.Now().UTC()
	query := `INSERT INTO trips (trip_code, driver_id, vehicle_license_plate, status, total_distance_km, estimated_duration_min, actual_start_at, actual_end_at, created_at, updated_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		t.TripCode, t.DriverID, t.VehicleLicensePlate, t.Status, t.TotalDistanceKm, t.EstimatedDurationMin,
		t.ActualStartAt, t.ActualEndAt, now, now, t.CreatedBy)
	if err != nil {
		return fmt.Errorf("create trip: %w", err)
	}
	id, _ := result.LastInsertId()
	t.ID = id
	t.CreatedAt = now
	t.UpdatedAt = now
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*entity.Trip, error) {
	query := `SELECT ` + tripColumns + ` FROM trips WHERE id = ?`
	t := &entity.Trip{}
	var driverID sql.NullInt64
	var startAt, endAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.TripCode, &driverID, &t.VehicleLicensePlate, &t.Status, &t.TotalDistanceKm, &t.EstimatedDurationMin,
		&startAt, &endAt, &t.CreatedAt, &t.UpdatedAt, &t.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get trip: %w", err)
	}
	if driverID.Valid {
		t.DriverID = &driverID.Int64
	}
	if startAt.Valid {
		t.ActualStartAt = &startAt.Time
	}
	if endAt.Valid {
		t.ActualEndAt = &endAt.Time
	}
	return t, nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*entity.Trip, error) {
	query := `SELECT ` + tripColumns + ` FROM trips WHERE trip_code = ?`
	t := &entity.Trip{}
	var driverID sql.NullInt64
	var startAt, endAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&t.ID, &t.TripCode, &driverID, &t.VehicleLicensePlate, &t.Status, &t.TotalDistanceKm, &t.EstimatedDurationMin,
		&startAt, &endAt, &t.CreatedAt, &t.UpdatedAt, &t.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, apierrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get trip by code: %w", err)
	}
	if driverID.Valid {
		t.DriverID = &driverID.Int64
	}
	if startAt.Valid {
		t.ActualStartAt = &startAt.Time
	}
	if endAt.Valid {
		t.ActualEndAt = &endAt.Time
	}
	return t, nil
}

func (r *Repository) List(ctx context.Context, offset, limit int) ([]entity.Trip, error) {
	query := `SELECT ` + tripColumns + ` FROM trips ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list trips: %w", err)
	}
	defer rows.Close()

	var trips []entity.Trip
	for rows.Next() {
		var t entity.Trip
		var driverID sql.NullInt64
		var startAt, endAt sql.NullTime
		if err := rows.Scan(
			&t.ID, &t.TripCode, &driverID, &t.VehicleLicensePlate, &t.Status, &t.TotalDistanceKm, &t.EstimatedDurationMin,
			&startAt, &endAt, &t.CreatedAt, &t.UpdatedAt, &t.CreatedBy,
		); err != nil {
			return nil, fmt.Errorf("scan trip: %w", err)
		}
		if driverID.Valid {
			t.DriverID = &driverID.Int64
		}
		if startAt.Valid {
			t.ActualStartAt = &startAt.Time
		}
		if endAt.Valid {
			t.ActualEndAt = &endAt.Time
		}
		trips = append(trips, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if trips == nil {
		trips = []entity.Trip{}
	}
	return trips, nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM trips`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count trips: %w", err)
	}
	return count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, t *entity.Trip) error {
	now := time.Now().UTC()
	query := `UPDATE trips SET trip_code=?, driver_id=?, vehicle_license_plate=?, status=?, total_distance_km=?, estimated_duration_min=?, actual_start_at=?, actual_end_at=?, updated_at=?
		WHERE id=?`
	result, err := r.db.ExecContext(ctx, query,
		t.TripCode, t.DriverID, t.VehicleLicensePlate, t.Status, t.TotalDistanceKm, t.EstimatedDurationMin,
		t.ActualStartAt, t.ActualEndAt, now, id)
	if err != nil {
		return fmt.Errorf("update trip: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	t.UpdatedAt = now
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM trips WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete trip: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apierrors.ErrNotFound
	}
	return nil
}

func (r *Repository) CreateStops(ctx context.Context, tripID int64, stops []entity.TripStop) error {
	query := `INSERT INTO trip_stops (trip_id, order_code, stop_type, address, lat, lng, status, planned_at, arrived_at, departure_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	for i := range stops {
		stop := &stops[i]
		result, err := r.db.ExecContext(ctx, query,
			tripID, stop.OrderCode, stop.StopType, stop.Address, stop.Lat, stop.Lng,
			stop.Status, stop.PlannedAt, stop.ArrivedAt, stop.DepartureAt)
		if err != nil {
			return fmt.Errorf("create trip stop: %w", err)
		}
		id, _ := result.LastInsertId()
		stop.ID = id
		stop.TripID = tripID
	}
	return nil
}

func (r *Repository) GetStops(ctx context.Context, tripID int64) ([]entity.TripStop, error) {
	query := `SELECT id, trip_id, order_code, stop_type, address, lat, lng, status, planned_at, arrived_at, departure_at
		FROM trip_stops WHERE trip_id = ? ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query, tripID)
	if err != nil {
		return nil, fmt.Errorf("get trip stops: %w", err)
	}
	defer rows.Close()

	var stops []entity.TripStop
	for rows.Next() {
		var s entity.TripStop
		var plannedAt, arrivedAt, departureAt sql.NullTime
		if err := rows.Scan(
			&s.ID, &s.TripID, &s.OrderCode, &s.StopType, &s.Address, &s.Lat, &s.Lng,
			&s.Status, &plannedAt, &arrivedAt, &departureAt,
		); err != nil {
			return nil, fmt.Errorf("scan trip stop: %w", err)
		}
		if plannedAt.Valid {
			s.PlannedAt = &plannedAt.Time
		}
		if arrivedAt.Valid {
			s.ArrivedAt = &arrivedAt.Time
		}
		if departureAt.Valid {
			s.DepartureAt = &departureAt.Time
		}
		stops = append(stops, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if stops == nil {
		stops = []entity.TripStop{}
	}
	return stops, nil
}

func (r *Repository) DeleteStops(ctx context.Context, tripID int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM trip_stops WHERE trip_id = ?", tripID)
	if err != nil {
		return fmt.Errorf("delete trip stops: %w", err)
	}
	return nil
}
