package models

import "time"

type Trip struct {
	ID                   string                 `json:"_id,omitempty" bson:"_id,omitempty"`
	TripCode             string                 `json:"trip_code" bson:"trip_code"`
	DriverID             *string                `json:"driver_id,omitempty" bson:"driver_id,omitempty"`
	OrderIDs             []string               `json:"order_ids" bson:"order_ids"`
	VehicleInfo          map[string]interface{} `json:"vehicle_info,omitempty" bson:"vehicle_info,omitempty"`
	Status               string                 `json:"status" bson:"status"`
	TotalDistanceKm      float64                `json:"total_distance_km" bson:"total_distance_km"`
	EstimatedDurationMin int                    `json:"estimated_duration_min" bson:"estimated_duration_min"`
	ActualStartAt        *time.Time             `json:"actual_start_at,omitempty" bson:"actual_start_at,omitempty"`
	ActualEndAt          *time.Time             `json:"actual_end_at,omitempty" bson:"actual_end_at,omitempty"`
	CreatedAt            *time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt            *time.Time             `json:"updated_at" bson:"updated_at"`
	CreatedBy            *string                `json:"created_by,omitempty" bson:"created_by,omitempty"`
}
