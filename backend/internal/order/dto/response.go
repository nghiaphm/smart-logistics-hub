package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/order/entity"
)

type OrderResponse struct {
	ID                 int64               `json:"id"`
	OrderCode          string              `json:"order_code"`
	SenderName         string              `json:"sender_name"`
	SenderPhone        string              `json:"sender_phone"`
	SenderAddress      string              `json:"sender_address"`
	SenderProvince     string              `json:"sender_province"`
	SenderDistrict     string              `json:"sender_district"`
	SenderWard         string              `json:"sender_ward"`
	SenderPostalCode   string              `json:"sender_postal_code"`
	ReceiverName       string              `json:"receiver_name"`
	ReceiverPhone      string              `json:"receiver_phone"`
	ReceiverAddress    string              `json:"receiver_address"`
	ReceiverProvince   string              `json:"receiver_province"`
	ReceiverDistrict   string              `json:"receiver_district"`
	ReceiverWard       string              `json:"receiver_ward"`
	ReceiverPostalCode string              `json:"receiver_postal_code"`
	Status             string              `json:"status"`
	AssignedDriverID   *int64              `json:"assigned_driver_id"`
	Items              []OrderItemResponse `json:"items,omitempty"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	CreatedBy          string              `json:"created_by"`
}

type OrderItemResponse struct {
	ID          int64  `json:"id"`
	OrderID     int64  `json:"order_id"`
	ProductID   *int64 `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	WeightGram  int    `json:"weight_gram"`
}

type PaginatedResponse struct {
	Items []OrderResponse `json:"items"`
	Total int             `json:"total"`
	Skip  int             `json:"skip"`
	Limit int             `json:"limit"`
}

func ToResponse(o *entity.Order) OrderResponse {
	return OrderResponse{
		ID:                 o.ID,
		OrderCode:          o.OrderCode,
		SenderName:         o.SenderName,
		SenderPhone:        o.SenderPhone,
		SenderAddress:      o.SenderAddress,
		SenderProvince:     o.SenderProvince,
		SenderDistrict:     o.SenderDistrict,
		SenderWard:         o.SenderWard,
		SenderPostalCode:   o.SenderPostalCode,
		ReceiverName:       o.ReceiverName,
		ReceiverPhone:      o.ReceiverPhone,
		ReceiverAddress:    o.ReceiverAddress,
		ReceiverProvince:   o.ReceiverProvince,
		ReceiverDistrict:   o.ReceiverDistrict,
		ReceiverWard:       o.ReceiverWard,
		ReceiverPostalCode: o.ReceiverPostalCode,
		Status:             o.Status,
		AssignedDriverID:   o.AssignedDriverID,
		CreatedAt:          o.CreatedAt,
		UpdatedAt:          o.UpdatedAt,
		CreatedBy:          o.CreatedBy,
	}
}

func ToResponseList(orders []entity.Order) []OrderResponse {
	res := make([]OrderResponse, len(orders))
	for i, o := range orders {
		res[i] = ToResponse(&o)
	}
	return res
}
