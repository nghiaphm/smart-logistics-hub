package entity

type OrderItem struct {
	ID          int64  `json:"id" db:"id"`
	OrderID     int64  `json:"order_id" db:"order_id"`
	ProductID   *int64 `json:"product_id" db:"product_id"`
	ProductName string `json:"product_name" db:"product_name"`
	Quantity    int    `json:"quantity" db:"quantity"`
	WeightGram  int    `json:"weight_gram" db:"weight_gram"`
}
