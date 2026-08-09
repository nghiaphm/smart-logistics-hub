package entity

import "time"

type AIEvent struct {
	ID              int64     `json:"id" db:"id"`
	EventCode       string    `json:"event_code" db:"event_code"`
	LicensePlate    string    `json:"license_plate" db:"license_plate"`
	ConfidenceScore float64   `json:"confidence_score" db:"confidence_score"`
	EventType       string    `json:"event_type" db:"event_type"`
	GateID          string    `json:"gate_id" db:"gate_id"`
	Timestamp       time.Time `json:"timestamp" db:"timestamp"`
	MatchedDriverID *int64    `json:"matched_driver_id" db:"matched_driver_id"`
	MatchedTripID   *int64    `json:"matched_trip_id" db:"matched_trip_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
