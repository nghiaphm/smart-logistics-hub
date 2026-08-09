package ai_events

type AIEventResponse struct {
	AIEventCreate   `json:",inline"`
	ID              string  `json:"_id"`
	EventCode       string  `json:"event_code"`
	MatchedDriverID *string `json:"matched_driver_id,omitempty"`
	MatchedTripID   *string `json:"matched_trip_id,omitempty"`
}
