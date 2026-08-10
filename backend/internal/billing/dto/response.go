package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/billing/entity"
)

type BillingResponse struct {
	ID            int64      `json:"id"`
	BillingCode   string     `json:"billing_code"`
	OrderCode     string     `json:"order_code"`
	AmountTotal   float64    `json:"amount_total"`
	Currency      string     `json:"currency"`
	PaymentMethod string     `json:"payment_method"`
	PaymentStatus string     `json:"payment_status"`
	TransactionID string     `json:"transaction_id"`
	PayerName     string     `json:"payer_name"`
	PayerPhone    string     `json:"payer_phone"`
	PayerEmail    string     `json:"payer_email"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	PaidAt        *time.Time `json:"paid_at"`
	CreatedBy     string     `json:"created_by"`
}

type PaginatedResponse struct {
	Items []BillingResponse `json:"items"`
	Total int               `json:"total"`
	Skip  int               `json:"skip"`
	Limit int               `json:"limit"`
}

func ToResponse(b *entity.Billing) BillingResponse {
	return BillingResponse{
		ID:            b.ID,
		BillingCode:   b.BillingCode,
		OrderCode:     b.OrderCode,
		AmountTotal:   b.AmountTotal,
		Currency:      b.Currency,
		PaymentMethod: b.PaymentMethod,
		PaymentStatus: b.PaymentStatus,
		TransactionID: b.TransactionID,
		PayerName:     b.PayerName,
		PayerPhone:    b.PayerPhone,
		PayerEmail:    b.PayerEmail,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		PaidAt:        b.PaidAt,
		CreatedBy:     b.CreatedBy,
	}
}

func ToResponseList(billings []entity.Billing) []BillingResponse {
	res := make([]BillingResponse, len(billings))
	for i, b := range billings {
		res[i] = ToResponse(&b)
	}
	return res
}
