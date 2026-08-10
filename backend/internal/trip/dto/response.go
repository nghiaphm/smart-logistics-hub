package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/trip/entity"
)

type TripStopResponse struct {
	ID          int64      `json:"id"`
	TripID      int64      `json:"trip_id"`
	OrderCode   string     `json:"order_code"`
	StopType    string     `json:"stop_type"`
	Address     string     `json:"address"`
	Lat         float64    `json:"lat"`
	Lng         float64    `json:"lng"`
	Status      string     `json:"status"`
	PlannedAt   *time.Time `json:"planned_at"`
	ArrivedAt   *time.Time `json:"arrived_at"`
	DepartureAt *time.Time `json:"departure_at"`
}

type TripResponse struct {
	ID                   int64              `json:"id"`
	TripCode             string             `json:"trip_code"`
	DriverID             *int64             `json:"driver_id"`
	VehicleLicensePlate  string             `json:"vehicle_license_plate"`
	Status               string             `json:"status"`
	TotalDistanceKm      float64            `json:"total_distance_km"`
	EstimatedDurationMin int                `json:"estimated_duration_min"`
	ActualStartAt        *time.Time         `json:"actual_start_at"`
	ActualEndAt          *time.Time         `json:"actual_end_at"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
	CreatedBy            string             `json:"created_by"`
	Stops                []TripStopResponse `json:"stops"`
}

type PaginatedResponse struct {
	Items []TripResponse `json:"items"`
	Total int            `json:"total"`
	Skip  int            `json:"skip"`
	Limit int            `json:"limit"`
}

func ToStopResponse(s *entity.TripStop) TripStopResponse {
	return TripStopResponse{
		ID:          s.ID,
		TripID:      s.TripID,
		OrderCode:   s.OrderCode,
		StopType:    s.StopType,
		Address:     s.Address,
		Lat:         s.Lat,
		Lng:         s.Lng,
		Status:      s.Status,
		PlannedAt:   s.PlannedAt,
		ArrivedAt:   s.ArrivedAt,
		DepartureAt: s.DepartureAt,
	}
}

func ToResponse(t *entity.Trip, stops []entity.TripStop) TripResponse {
	res := TripResponse{
		ID:                   t.ID,
		TripCode:             t.TripCode,
		DriverID:             t.DriverID,
		VehicleLicensePlate:  t.VehicleLicensePlate,
		Status:               t.Status,
		TotalDistanceKm:      t.TotalDistanceKm,
		EstimatedDurationMin: t.EstimatedDurationMin,
		ActualStartAt:        t.ActualStartAt,
		ActualEndAt:          t.ActualEndAt,
		CreatedAt:            t.CreatedAt,
		UpdatedAt:            t.UpdatedAt,
		CreatedBy:            t.CreatedBy,
		Stops:                []TripStopResponse{},
	}
	if len(stops) > 0 {
		res.Stops = make([]TripStopResponse, len(stops))
		for i, s := range stops {
			res.Stops[i] = ToStopResponse(&s)
		}
	}
	return res
}

func ToResponseList(trips []entity.Trip) []TripResponse {
	res := make([]TripResponse, len(trips))
	for i, t := range trips {
		res[i] = ToResponse(&t, nil)
	}
	return res
}
