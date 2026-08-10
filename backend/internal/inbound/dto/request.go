package dto

type InboundItemRequest struct {
	ProductID   int64 `json:"product_id" binding:"required"`
	ExpectedQty int   `json:"expected_qty"`
	ReceivedQty int   `json:"received_qty"`
	RejectedQty int   `json:"rejected_qty"`
	QcPassed    int   `json:"qc_passed"`
}

type CreateInboundRequest struct {
	ReceiptCode  string               `json:"receipt_code" binding:"required"`
	SupplierName string               `json:"supplier_name" binding:"required"`
	WarehouseID  int64                `json:"warehouse_id" binding:"required,min=1"`
	Status       string               `json:"status"`
	CreatedBy    string               `json:"created_by"`
	Items        []InboundItemRequest `json:"items" binding:"required,dive"`
}

type UpdateInboundRequest struct {
	Items  *[]InboundItemRequest `json:"items,omitempty"`
	Status *string               `json:"status,omitempty"`
}
