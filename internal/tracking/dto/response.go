package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/tracking/entity"
)

type TrackingEventResponse struct {
	ID           int64     `json:"id"`
	OrderCode    string    `json:"order_code"`
	DriverCode   string    `json:"driver_code"`
	StatusUpdate string    `json:"status_update"`
	Lat          float64   `json:"lat"`
	Lng          float64   `json:"lng"`
	Note         string    `json:"note"`
	Timestamp    time.Time `json:"timestamp"`
}

type PaginatedResponse struct {
	Items []TrackingEventResponse `json:"items"`
	Total int                     `json:"total"`
	Skip  int                     `json:"skip"`
	Limit int                     `json:"limit"`
}

func ToResponse(event *entity.TrackingEvent) TrackingEventResponse {
	return TrackingEventResponse{
		ID: event.ID, OrderCode: event.OrderCode, DriverCode: event.DriverCode,
		StatusUpdate: event.StatusUpdate, Lat: event.Lat, Lng: event.Lng,
		Note: event.Note, Timestamp: event.Timestamp,
	}
}

func ToResponseList(events []entity.TrackingEvent) []TrackingEventResponse {
	res := make([]TrackingEventResponse, len(events))
	for i, e := range events {
		res[i] = ToResponse(&e)
	}
	return res
}
