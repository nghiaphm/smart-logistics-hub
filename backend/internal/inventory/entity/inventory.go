package entity

import "time"

type Inventory struct {
	ID           int64     `json:"id" db:"id"`
	ProductID    int64     `json:"product_id" db:"product_id"`
	WarehouseID  int64     `json:"warehouse_id" db:"warehouse_id"`
	AvailableQty int       `json:"available_qty" db:"available_qty"`
	ReservedQty  int       `json:"reserved_qty" db:"reserved_qty"`
	DamagedQty   int       `json:"damaged_qty" db:"damaged_qty"`
	HoldQty      int       `json:"hold_qty" db:"hold_qty"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy    string    `json:"updated_by" db:"updated_by"`
}
