package dto

import "time"

type DriverResponse struct {
	ID           int64     `json:"id"`
	DriverCode   string    `json:"driver_code"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	VehicleType  string    `json:"vehicle_type"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CurrentLat   float64   `json:"current_lat"`
	CurrentLng   float64   `json:"current_lng"`
	WarehouseID  *int64    `json:"warehouse_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    string    `json:"created_by"`
}
