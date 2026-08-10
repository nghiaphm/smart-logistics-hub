package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/inbound/entity"
)

type InboundItemResponse struct {
	ID          int64 `json:"id"`
	InboundID   int64 `json:"inbound_id"`
	ProductID   int64 `json:"product_id"`
	ExpectedQty int   `json:"expected_qty"`
	ReceivedQty int   `json:"received_qty"`
	RejectedQty int   `json:"rejected_qty"`
	QcPassed    int   `json:"qc_passed"`
}

type InboundResponse struct {
	ID           int64                 `json:"id"`
	ReceiptCode  string                `json:"receipt_code"`
	SupplierName string                `json:"supplier_name"`
	WarehouseID  int64                 `json:"warehouse_id"`
	Status       string                `json:"status"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	CompletedAt  *time.Time            `json:"completed_at"`
	CreatedBy    string                `json:"created_by"`
	Items        []InboundItemResponse `json:"items"`
}

type PaginatedResponse struct {
	Items []InboundResponse `json:"items"`
	Total int               `json:"total"`
	Skip  int               `json:"skip"`
	Limit int               `json:"limit"`
}

func ToItemResponse(i *entity.InboundItem) InboundItemResponse {
	return InboundItemResponse{
		ID:          i.ID,
		InboundID:   i.InboundID,
		ProductID:   i.ProductID,
		ExpectedQty: i.ExpectedQty,
		ReceivedQty: i.ReceivedQty,
		RejectedQty: i.RejectedQty,
		QcPassed:    i.QcPassed,
	}
}

func ToResponse(in *entity.Inbound, items []entity.InboundItem) InboundResponse {
	res := InboundResponse{
		ID:           in.ID,
		ReceiptCode:  in.ReceiptCode,
		SupplierName: in.SupplierName,
		WarehouseID:  in.WarehouseID,
		Status:       in.Status,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		CompletedAt:  in.CompletedAt,
		CreatedBy:    in.CreatedBy,
		Items:        []InboundItemResponse{},
	}
	if len(items) > 0 {
		res.Items = make([]InboundItemResponse, len(items))
		for i, it := range items {
			res.Items[i] = ToItemResponse(&it)
		}
	}
	return res
}

func ToResponseList(inbounds []entity.Inbound) []InboundResponse {
	res := make([]InboundResponse, len(inbounds))
	for i, in := range inbounds {
		res[i] = ToResponse(&in, nil)
	}
	return res
}
