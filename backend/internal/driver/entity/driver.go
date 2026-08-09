package entity

import "time"

type Driver struct {
	ID           int64     `json:"id" db:"id"`
	DriverCode   string    `json:"driver_code" db:"driver_code"`
	FullName     string    `json:"full_name" db:"full_name"`
	Phone        string    `json:"phone" db:"phone"`
	VehicleType  string    `json:"vehicle_type" db:"vehicle_type"`
	LicensePlate string    `json:"license_plate" db:"license_plate"`
	Status       string    `json:"status" db:"status"`
	CurrentLat   float64   `json:"current_lat" db:"current_lat"`
	CurrentLng   float64   `json:"current_lng" db:"current_lng"`
	WarehouseID  *int64    `json:"warehouse_id" db:"warehouse_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy    string    `json:"created_by" db:"created_by"`
}
