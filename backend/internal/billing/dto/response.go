package billings

import "time"

type BillingResponse struct {
	BillingCreate `json:",inline"`
	ID            string     `json:"_id"`
	TransactionID *string    `json:"transaction_id,omitempty"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	CreatedBy     *string    `json:"created_by,omitempty"`
}
