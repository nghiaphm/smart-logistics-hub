package warehouses

import "time"

type WarehouseResponse struct {
	WarehouseCreate `json:",inline"`
	ID              string     `json:"_id" bson:"_id"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
