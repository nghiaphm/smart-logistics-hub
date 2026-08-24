package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/driver/entity"
)

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

type PaginatedResponse struct {
	Items []DriverResponse `json:"items"`
	Total int              `json:"total"`
	Skip  int              `json:"skip"`
	Limit int              `json:"limit"`
}

func ToResponse(d *entity.Driver) DriverResponse {
	return DriverResponse{
		ID:           d.ID,
		DriverCode:   d.DriverCode,
		FullName:     d.FullName,
		Phone:        d.Phone,
		VehicleType:  d.VehicleType,
		LicensePlate: d.LicensePlate,
		Status:       d.Status,
		CurrentLat:   d.CurrentLat,
		CurrentLng:   d.CurrentLng,
		WarehouseID:  d.WarehouseID,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		CreatedBy:    d.CreatedBy,
	}
}

func ToResponseList(drivers []entity.Driver) []DriverResponse {
	res := make([]DriverResponse, len(drivers))
	for i, d := range drivers {
		res[i] = ToResponse(&d)
	}
	return res
}
