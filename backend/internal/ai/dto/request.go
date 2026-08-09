package ai_events

type AIEventCreate struct {
	LicensePlate    string  `json:"license_plate" binding:"required"`
	ConfidenceScore float64 `json:"confidence_score" binding:"required,gte=0.0,lte=1.0"`
	EventType       string  `json:"event_type" binding:"required,oneof=INBOUND OUTBOUND"`
	GateID          string  `json:"gate_id" binding:"required"`
	Timestamp       string  `json:"timestamp" binding:"required"`
}
