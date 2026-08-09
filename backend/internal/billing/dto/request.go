package billings

import "time"

type PayerInfo struct {
	Name  string  `json:"name" binding:"required"`
	Phone string  `json:"phone" binding:"required"`
	Email *string `json:"email,omitempty"`
}

type BillingCreate struct {
	BillingCode   string    `json:"billing_code" binding:"required"`
	OrderCode     string    `json:"order_code" binding:"required"`
	AmountTotal   float64   `json:"amount_total" binding:"required"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method" binding:"oneof=COD VNPAY BANK_TRANSFER"`
	PaymentStatus string    `json:"payment_status" binding:"oneof=UNPAID PENDING PAID FAILED"`
	PayerInfo     PayerInfo `json:"payer_info" binding:"required"`
}

type BillingUpdate struct {
	PaymentStatus *string    `json:"payment_status,omitempty" binding:"omitempty,oneof=UNPAID PENDING PAID FAILED"`
	TransactionID *string    `json:"transaction_id,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
}
