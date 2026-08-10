package dto

import "time"

type PayerInfo struct {
	Name  string  `json:"name" binding:"required"`
	Phone string  `json:"phone" binding:"required"`
	Email *string `json:"email,omitempty"`
}

type CreateBillingRequest struct {
	BillingCode   string    `json:"billing_code" binding:"required"`
	OrderCode     string    `json:"order_code" binding:"required"`
	AmountTotal   float64   `json:"amount_total"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `json:"payment_status"`
	PayerInfo     PayerInfo `json:"payer_info" binding:"required"`
	CreatedBy     string    `json:"created_by"`
}

type UpdateBillingRequest struct {
	PaymentStatus *string    `json:"payment_status,omitempty"`
	TransactionID *string    `json:"transaction_id,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
}
