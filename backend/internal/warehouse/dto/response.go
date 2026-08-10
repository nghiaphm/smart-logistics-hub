package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/warehouse/entity"
)

type WarehouseResponse struct {
	ID            int64     `json:"id"`
	WarehouseCode string    `json:"warehouse_code"`
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	Lat           float64   `json:"lat"`
	Lng           float64   `json:"lng"`
	ContactPhone  string    `json:"contact_phone"`
	ManagerName   string    `json:"manager_name"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PaginatedResponse struct {
	Items []WarehouseResponse `json:"items"`
	Total int                 `json:"total"`
	Skip  int                 `json:"skip"`
	Limit int                 `json:"limit"`
}

func ToResponse(w *entity.Warehouse) WarehouseResponse {
	return WarehouseResponse{
		ID:            w.ID,
		WarehouseCode: w.WarehouseCode,
		Name:          w.Name,
		Address:       w.Address,
		Lat:           w.Lat,
		Lng:           w.Lng,
		ContactPhone:  w.ContactPhone,
		ManagerName:   w.ManagerName,
		IsActive:      w.IsActive,
		CreatedAt:     w.CreatedAt,
		UpdatedAt:     w.UpdatedAt,
	}
}

func ToResponseList(warehouses []entity.Warehouse) []WarehouseResponse {
	res := make([]WarehouseResponse, len(warehouses))
	for i, w := range warehouses {
		res[i] = ToResponse(&w)
	}
	return res
}
