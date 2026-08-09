package entity

import "time"

type Trip struct {
	ID                   int64      `json:"id" db:"id"`
	TripCode             string     `json:"trip_code" db:"trip_code"`
	DriverID             *int64     `json:"driver_id" db:"driver_id"`
	VehicleLicensePlate  string     `json:"vehicle_license_plate" db:"vehicle_license_plate"`
	Status               string     `json:"status" db:"status"`
	TotalDistanceKm      float64    `json:"total_distance_km" db:"total_distance_km"`
	EstimatedDurationMin int        `json:"estimated_duration_min" db:"estimated_duration_min"`
	ActualStartAt        *time.Time `json:"actual_start_at" db:"actual_start_at"`
	ActualEndAt          *time.Time `json:"actual_end_at" db:"actual_end_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy            string     `json:"created_by" db:"created_by"`
}
