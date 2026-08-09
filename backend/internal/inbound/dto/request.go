package inbounds

type InboundItem struct {
	ProductID   string `json:"product_id" binding:"required"`
	ExpectedQty int    `json:"expected_qty" binding:"required"`
	ReceivedQty int    `json:"received_qty"`
	RejectedQty int    `json:"rejected_qty"`
	QCPassed    bool   `json:"qc_passed"`
}

type InboundCreate struct {
	ReceiptCode  string        `json:"receipt_code" binding:"required"`
	SupplierName string        `json:"supplier_name" binding:"required"`
	Items        []InboundItem `json:"items" binding:"required,dive"`
	Status       string        `json:"status" binding:"oneof=PENDING RECEIVING QC_CHECKING COMPLETED"`
}

type InboundUpdate struct {
	Items  *[]InboundItem `json:"items,omitempty" binding:"omitempty,dive"`
	Status *string        `json:"status,omitempty" binding:"omitempty,oneof=PENDING RECEIVING QC_CHECKING COMPLETED"`
}
