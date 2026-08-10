package dto

type CreateOrderRequest struct {
	OrderCode          string             `json:"order_code" binding:"required"`
	WarehouseID        int64              `json:"warehouse_id" binding:"required,min=1"`
	SenderName         string             `json:"sender_name" binding:"required"`
	SenderPhone        string             `json:"sender_phone" binding:"required"`
	SenderAddress      string             `json:"sender_address" binding:"required"`
	SenderProvince     string             `json:"sender_province"`
	SenderDistrict     string             `json:"sender_district"`
	SenderWard         string             `json:"sender_ward"`
	SenderPostalCode   string             `json:"sender_postal_code"`
	ReceiverName       string             `json:"receiver_name" binding:"required"`
	ReceiverPhone      string             `json:"receiver_phone" binding:"required"`
	ReceiverAddress    string             `json:"receiver_address" binding:"required"`
	ReceiverProvince   string             `json:"receiver_province"`
	ReceiverDistrict   string             `json:"receiver_district"`
	ReceiverWard       string             `json:"receiver_ward"`
	ReceiverPostalCode string             `json:"receiver_postal_code"`
	Items              []OrderItemRequest `json:"items" binding:"dive"`
	Status             string             `json:"status" binding:"oneof=PENDING RESERVED PICKING PACKING SORTING SHIPPING COMPLETED PICKING_UP"`
	AssignedDriverID   *int64             `json:"assigned_driver_id"`
}

type OrderItemRequest struct {
	ProductID   *int64 `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
	WeightGram  int    `json:"weight_gram"`
}

type UpdateOrderRequest struct {
	OrderCode          *string `json:"order_code"`
	SenderName         *string `json:"sender_name"`
	SenderPhone        *string `json:"sender_phone"`
	SenderAddress      *string `json:"sender_address"`
	SenderProvince     *string `json:"sender_province"`
	SenderDistrict     *string `json:"sender_district"`
	SenderWard         *string `json:"sender_ward"`
	SenderPostalCode   *string `json:"sender_postal_code"`
	ReceiverName       *string `json:"receiver_name"`
	ReceiverPhone      *string `json:"receiver_phone"`
	ReceiverAddress    *string `json:"receiver_address"`
	ReceiverProvince   *string `json:"receiver_province"`
	ReceiverDistrict   *string `json:"receiver_district"`
	ReceiverWard       *string `json:"receiver_ward"`
	ReceiverPostalCode *string `json:"receiver_postal_code"`
	Status             *string `json:"status" binding:"omitempty,oneof=PENDING RESERVED PICKING PACKING SORTING SHIPPING COMPLETED PICKING_UP"`
	AssignedDriverID   *int64  `json:"assigned_driver_id"`
}
