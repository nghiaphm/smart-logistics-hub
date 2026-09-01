package entity

import "time"

type Order struct {
	ID                 int64     `json:"id" db:"id"`
	OrderCode          string    `json:"order_code" db:"order_code"`
	SenderWorkspaceID  *int64    `json:"sender_workspace_id" db:"sender_workspace_id"`
	SenderName         string    `json:"sender_name" db:"sender_name"`
	SenderPhone        string    `json:"sender_phone" db:"sender_phone"`
	SenderAddress      string    `json:"sender_address" db:"sender_address"`
	SenderProvince     string    `json:"sender_province" db:"sender_province"`
	SenderDistrict     string    `json:"sender_district" db:"sender_district"`
	SenderWard         string    `json:"sender_ward" db:"sender_ward"`
	SenderPostalCode   string    `json:"sender_postal_code" db:"sender_postal_code"`
	ReceiverName       string    `json:"receiver_name" db:"receiver_name"`
	ReceiverPhone      string    `json:"receiver_phone" db:"receiver_phone"`
	ReceiverAddress    string    `json:"receiver_address" db:"receiver_address"`
	ReceiverProvince   string    `json:"receiver_province" db:"receiver_province"`
	ReceiverDistrict   string    `json:"receiver_district" db:"receiver_district"`
	ReceiverWard       string    `json:"receiver_ward" db:"receiver_ward"`
	ReceiverPostalCode string    `json:"receiver_postal_code" db:"receiver_postal_code"`
	Status             string    `json:"status" db:"status"`
	AssignedDriverID   *int64    `json:"assigned_driver_id" db:"assigned_driver_id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy          string    `json:"created_by" db:"created_by"`
}
