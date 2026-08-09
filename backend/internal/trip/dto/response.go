package trips

import "time"

type TripResponse struct {
	TripCreate  `json:",inline"`
	ID          string     `json:"_id" bson:"_id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedBy   *string    `json:"created_by,omitempty"`
}
