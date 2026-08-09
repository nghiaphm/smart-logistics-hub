package dto

type CreateInventoryRequest struct {
	ProductID    int64  `json:"product_id" binding:"required"`
	WarehouseID  int64  `json:"warehouse_id" binding:"required"`
	AvailableQty int    `json:"available_qty"`
	ReservedQty  int    `json:"reserved_qty"`
	DamagedQty   int    `json:"damaged_qty"`
	HoldQty      int    `json:"hold_qty"`
	UpdatedBy    string `json:"updated_by"`
}

type UpdateInventoryRequest struct {
	AvailableQty *int    `json:"available_qty,omitempty"`
	ReservedQty  *int    `json:"reserved_qty,omitempty"`
	DamagedQty   *int    `json:"damaged_qty,omitempty"`
	HoldQty      *int    `json:"hold_qty,omitempty"`
	Reason       *string `json:"reason,omitempty"`
	UpdatedBy    string  `json:"updated_by"`
}
