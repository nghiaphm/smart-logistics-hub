package models

import "time"

type Inbound struct {
	ID           string                   `json:"_id,omitempty" bson:"_id,omitempty"`
	ReceiptCode  string                   `json:"receipt_code" bson:"receipt_code"`
	SupplierName string                   `json:"supplier_name" bson:"supplier_name"`
	Items        []map[string]interface{} `json:"items" bson:"items"`
	Status       string                   `json:"status" bson:"status"`
	CreatedAt    *time.Time               `json:"created_at" bson:"created_at"`
	UpdatedAt    *time.Time               `json:"updated_at" bson:"updated_at"`
	CompletedAt  *time.Time               `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	CreatedBy    *string                  `json:"created_by,omitempty" bson:"created_by,omitempty"`
}
