package dto

type Location struct {
	Lat *float64 `json:"lat,omitempty"`
	Lng *float64 `json:"lng,omitempty"`
}

type CreateWarehouseRequest struct {
	WarehouseCode string    `json:"warehouse_code" binding:"required"`
	Name          string    `json:"name" binding:"required"`
	Address       string    `json:"address" binding:"required"`
	Location      *Location `json:"location,omitempty"`
	ContactPhone  *string   `json:"contact_phone,omitempty"`
	ManagerName   *string   `json:"manager_name,omitempty"`
}

type UpdateWarehouseRequest struct {
	Name         *string   `json:"name,omitempty"`
	Address      *string   `json:"address,omitempty"`
	Location     *Location `json:"location,omitempty"`
	ContactPhone *string   `json:"contact_phone,omitempty"`
	ManagerName  *string   `json:"manager_name,omitempty"`
	IsActive     *bool     `json:"is_active,omitempty"`
}
