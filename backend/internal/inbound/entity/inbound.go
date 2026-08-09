package entity

import "time"

type Inbound struct {
	ID           int64      `json:"id" db:"id"`
	ReceiptCode  string     `json:"receipt_code" db:"receipt_code"`
	SupplierName string     `json:"supplier_name" db:"supplier_name"`
	Status       string     `json:"status" db:"status"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	CompletedAt  *time.Time `json:"completed_at" db:"completed_at"`
	CreatedBy    string     `json:"created_by" db:"created_by"`
}

type InboundItem struct {
	ID          int64 `json:"id" db:"id"`
	InboundID   int64 `json:"inbound_id" db:"inbound_id"`
	ProductID   int64 `json:"product_id" db:"product_id"`
	ExpectedQty int   `json:"expected_qty" db:"expected_qty"`
	ReceivedQty int   `json:"received_qty" db:"received_qty"`
	RejectedQty int   `json:"rejected_qty" db:"rejected_qty"`
	QcPassed    int   `json:"qc_passed" db:"qc_passed"`
}
