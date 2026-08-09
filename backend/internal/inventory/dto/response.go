package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/inventory/entity"
)

type InventoryResponse struct {
	ID           int64     `json:"id"`
	ProductID    int64     `json:"product_id"`
	WarehouseID  int64     `json:"warehouse_id"`
	AvailableQty int       `json:"available_qty"`
	ReservedQty  int       `json:"reserved_qty"`
	DamagedQty   int       `json:"damaged_qty"`
	HoldQty      int       `json:"hold_qty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
}

type PaginatedResponse struct {
	Items []InventoryResponse `json:"items"`
	Total int                 `json:"total"`
	Skip  int                 `json:"skip"`
	Limit int                 `json:"limit"`
}

func ToResponse(inv *entity.Inventory) InventoryResponse {
	return InventoryResponse{
		ID:           inv.ID,
		ProductID:    inv.ProductID,
		WarehouseID:  inv.WarehouseID,
		AvailableQty: inv.AvailableQty,
		ReservedQty:  inv.ReservedQty,
		DamagedQty:   inv.DamagedQty,
		HoldQty:      inv.HoldQty,
		CreatedAt:    inv.CreatedAt,
		UpdatedAt:    inv.UpdatedAt,
		UpdatedBy:    inv.UpdatedBy,
	}
}

func ToResponseList(invs []entity.Inventory) []InventoryResponse {
	res := make([]InventoryResponse, len(invs))
	for i, inv := range invs {
		res[i] = ToResponse(&inv)
	}
	return res
}
