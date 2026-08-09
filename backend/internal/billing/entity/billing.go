package models

import "time"

type Billing struct {
	ID            string                 `json:"_id,omitempty" bson:"_id,omitempty"`
	BillingCode   string                 `json:"billing_code" bson:"billing_code"`
	OrderCode     string                 `json:"order_code" bson:"order_code"`
	AmountTotal   float64                `json:"amount_total" bson:"amount_total"`
	Currency      string                 `json:"currency" bson:"currency"`
	PaymentMethod string                 `json:"payment_method" bson:"payment_method"`
	PaymentStatus string                 `json:"payment_status" bson:"payment_status"`
	TransactionID *string                `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	PayerInfo     map[string]interface{} `json:"payer_info" bson:"payer_info"`
	CreatedAt     *time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt     *time.Time             `json:"updated_at" bson:"updated_at"`
	PaidAt        *time.Time             `json:"paid_at,omitempty" bson:"paid_at,omitempty"`
	CreatedBy     *string                `json:"created_by,omitempty" bson:"created_by,omitempty"`
}
