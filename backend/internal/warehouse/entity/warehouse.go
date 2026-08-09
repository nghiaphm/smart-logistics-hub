package entity

import "time"

type Warehouse struct {
	ID            int64     `json:"id" db:"id"`
	WarehouseCode string    `json:"warehouse_code" db:"warehouse_code"`
	Name          string    `json:"name" db:"name"`
	Address       string    `json:"address" db:"address"`
	Lat           float64   `json:"lat" db:"lat"`
	Lng           float64   `json:"lng" db:"lng"`
	ContactPhone  string    `json:"contact_phone" db:"contact_phone"`
	ManagerName   string    `json:"manager_name" db:"manager_name"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
