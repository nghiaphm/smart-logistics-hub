package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/ai/entity"
)

type AIEventResponse struct {
	ID              int64     `json:"id"`
	EventCode       string    `json:"event_code"`
	LicensePlate    string    `json:"license_plate"`
	ConfidenceScore float64   `json:"confidence_score"`
	EventType       string    `json:"event_type"`
	GateID          string    `json:"gate_id"`
	Timestamp       time.Time `json:"timestamp"`
	MatchedDriverID *int64    `json:"matched_driver_id"`
	MatchedTripID   *int64    `json:"matched_trip_id"`
	CreatedAt       time.Time `json:"created_at"`
	LowConfidence   bool      `json:"low_confidence"`
}

type PaginatedResponse struct {
	Items []AIEventResponse `json:"items"`
	Total int               `json:"total"`
	Skip  int               `json:"skip"`
	Limit int               `json:"limit"`
}

func ToResponse(e *entity.AIEvent, lowConfidence bool) AIEventResponse {
	return AIEventResponse{
		ID:              e.ID,
		EventCode:       e.EventCode,
		LicensePlate:    e.LicensePlate,
		ConfidenceScore: e.ConfidenceScore,
		EventType:       e.EventType,
		GateID:          e.GateID,
		Timestamp:       e.Timestamp,
		MatchedDriverID: e.MatchedDriverID,
		MatchedTripID:   e.MatchedTripID,
		CreatedAt:       e.CreatedAt,
		LowConfidence:   lowConfidence,
	}
}

func ToResponseList(events []entity.AIEvent) []AIEventResponse {
	res := make([]AIEventResponse, len(events))
	for i, e := range events {
		res[i] = ToResponse(&e, e.ConfidenceScore < 0.7)
	}
	return res
}
