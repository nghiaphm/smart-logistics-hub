package trips

import "time"

type Location struct {
	Lat float64 `json:"lat" binding:"required"`
	Lng float64 `json:"lng" binding:"required"`
}

type TripStop struct {
	OrderCode            string     `json:"order_code" binding:"required"`
	StopType             string     `json:"stop_type" binding:"required"`
	Location             Location   `json:"location" binding:"required"`
	Address              string     `json:"address" binding:"required"`
	Status               string     `json:"status"`
	EstimatedArrivalTime *time.Time `json:"estimated_arrival_time,omitempty"`
	ActualArrivalTime    *time.Time `json:"actual_arrival_time,omitempty"`
}

type TripCreate struct {
	TripCode            string     `json:"trip_code" binding:"required"`
	DriverCode          string     `json:"driver_code" binding:"required"`
	VehicleLicensePlate *string    `json:"vehicle_license_plate,omitempty"`
	Stops               []TripStop `json:"stops" binding:"required,dive"`
	Status              string     `json:"status"`
	TotalDistanceKm     *float64   `json:"total_distance_km,omitempty"`
}

type TripUpdate struct {
	Status      *string     `json:"status,omitempty"`
	Stops       *[]TripStop `json:"stops,omitempty" binding:"omitempty,dive"`
	StartedAt   *time.Time  `json:"started_at,omitempty"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
}
