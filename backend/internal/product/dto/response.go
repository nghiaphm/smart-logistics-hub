package products

import "time"

type ProductResponse struct {
	ProductCreate `json:",inline"`
	ID            string     `json:"_id"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	CreatedBy     *string    `json:"created_by,omitempty"`
}
