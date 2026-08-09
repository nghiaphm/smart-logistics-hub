package models

import "time"

type AIEvent struct {
	ID              string     `json:"_id,omitempty" bson:"_id,omitempty"`
	EventCode       string     `json:"event_code" bson:"event_code"`
	LicensePlate    string     `json:"license_plate" bson:"license_plate"`
	ConfidenceScore float64    `json:"confidence_score" bson:"confidence_score"`
	EventType       string     `json:"event_type" bson:"event_type"`
	GateID          string     `json:"gate_id" bson:"gate_id"`
	Timestamp       string     `json:"timestamp" bson:"timestamp"`
	CreatedAt       *time.Time `json:"created_at" bson:"created_at"`
	MatchedDriverID *string    `json:"matched_driver_id,omitempty" bson:"matched_driver_id,omitempty"`
	MatchedTripID   *string    `json:"matched_trip_id,omitempty" bson:"matched_trip_id,omitempty"`
}
