package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/vehicle/entity"
)

type VehicleResponse struct {
	ID           int64     `json:"id"`
	LicensePlate string    `json:"license_plate"`
	Type         string    `json:"type"`
	Capacity     float64   `json:"capacity"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    string    `json:"created_by"`
}

type PaginatedResponse struct {
	Items []VehicleResponse `json:"items"`
	Total int               `json:"total"`
	Skip  int               `json:"skip"`
	Limit int               `json:"limit"`
}

func ToResponse(v *entity.Vehicle) VehicleResponse {
	return VehicleResponse{
		ID:           v.ID,
		LicensePlate: v.LicensePlate,
		Type:         v.Type,
		Capacity:     v.Capacity,
		Status:       v.Status,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
		CreatedBy:    v.CreatedBy,
	}
}

func ToResponseList(vehicles []entity.Vehicle) []VehicleResponse {
	res := make([]VehicleResponse, len(vehicles))
	for i, v := range vehicles {
		res[i] = ToResponse(&v)
	}
	return res
}
