package dto

import "time"

type Location struct {
	Lat *float64 `json:"lat,omitempty"`
	Lng *float64 `json:"lng,omitempty"`
}

type TripStopRequest struct {
	OrderCode   string     `json:"order_code" binding:"required"`
	StopType    string     `json:"stop_type"`
	Address     string     `json:"address" binding:"required"`
	Location    *Location  `json:"location,omitempty"`
	Status      string     `json:"status"`
	PlannedAt   *time.Time `json:"planned_at,omitempty"`
	ArrivedAt   *time.Time `json:"arrived_at,omitempty"`
	DepartureAt *time.Time `json:"departure_at,omitempty"`
}

type CreateTripRequest struct {
	TripCode            string            `json:"trip_code" binding:"required"`
	DriverCode          string            `json:"driver_code" binding:"required"`
	VehicleLicensePlate *string           `json:"vehicle_license_plate,omitempty"`
	Status              string            `json:"status"`
	TotalDistanceKm     *float64          `json:"total_distance_km,omitempty"`
	Stops               []TripStopRequest `json:"stops" binding:"required,dive"`
}

type UpdateTripRequest struct {
	Status      *string            `json:"status,omitempty"`
	Stops       *[]TripStopRequest `json:"stops,omitempty"`
	StartedAt   *time.Time         `json:"started_at,omitempty"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
}

type AssignDriverRequest struct {
	DriverCode string `json:"driver_code" binding:"required"`
}
