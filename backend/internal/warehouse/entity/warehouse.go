package models

import "time"

type Warehouse struct {
	ID            string                 `json:"_id,omitempty" bson:"_id,omitempty"`
	WarehouseCode string                 `json:"warehouse_code" bson:"warehouse_code"`
	Name          string                 `json:"name" bson:"name"`
	Address       string                 `json:"address" bson:"address"`
	Location      map[string]interface{} `json:"location,omitempty" bson:"location,omitempty"`
	ContactPhone  *string                `json:"contact_phone,omitempty" bson:"contact_phone,omitempty"`
	ManagerName   *string                `json:"manager_name,omitempty" bson:"manager_name,omitempty"`
	CreatedAt     *time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt     *time.Time             `json:"updated_at" bson:"updated_at"`
	IsActive      bool                   `json:"is_active" bson:"is_active"`
}
