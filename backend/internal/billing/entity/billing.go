package entity

import "time"

type Billing struct {
	ID            int64      `json:"id" db:"id"`
	BillingCode   string     `json:"billing_code" db:"billing_code"`
	OrderCode     string     `json:"order_code" db:"order_code"`
	AmountTotal   float64    `json:"amount_total" db:"amount_total"`
	Currency      string     `json:"currency" db:"currency"`
	PaymentMethod string     `json:"payment_method" db:"payment_method"`
	PaymentStatus string     `json:"payment_status" db:"payment_status"`
	TransactionID string     `json:"transaction_id" db:"transaction_id"`
	PayerName     string     `json:"payer_name" db:"payer_name"`
	PayerPhone    string     `json:"payer_phone" db:"payer_phone"`
	PayerEmail    string     `json:"payer_email" db:"payer_email"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	PaidAt        *time.Time `json:"paid_at" db:"paid_at"`
	CreatedBy     string     `json:"created_by" db:"created_by"`
}
