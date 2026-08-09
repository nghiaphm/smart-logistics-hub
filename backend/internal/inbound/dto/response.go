package inbounds

import "time"

type InboundResponse struct {
	InboundCreate `json:",inline"`
	ID            string     `json:"_id" bson:"_id"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CreatedBy     *string    `json:"created_by,omitempty"`
}
