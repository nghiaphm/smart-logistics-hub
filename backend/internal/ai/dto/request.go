package dto

import "time"

type CreateAIEventRequest struct {
	EventCode       string    `json:"event_code,omitempty"`
	LicensePlate    string    `json:"license_plate" binding:"required"`
	ConfidenceScore float64   `json:"confidence_score" binding:"gte=0,lte=1"`
	EventType       string    `json:"event_type" binding:"required,oneof=INBOUND OUTBOUND"`
	GateID          string    `json:"gate_id" binding:"required"`
	Timestamp       time.Time `json:"timestamp,omitempty"`
}

type UpdateAIEventRequest struct {
	LicensePlate    *string    `json:"license_plate,omitempty"`
	ConfidenceScore *float64   `json:"confidence_score,omitempty"`
	EventType       *string    `json:"event_type,omitempty"`
	GateID          *string    `json:"gate_id,omitempty"`
	Timestamp       *time.Time `json:"timestamp,omitempty"`
}
