package dto

type CreateDriverRequest struct {
	DriverCode   string  `json:"driver_code" binding:"required"`
	FullName     string  `json:"full_name" binding:"required"`
	Phone        string  `json:"phone" binding:"required"`
	VehicleType  string  `json:"vehicle_type" binding:"required"`
	LicensePlate string  `json:"license_plate" binding:"required"`
	Status       string  `json:"status" binding:"oneof=AVAILABLE BUSY OFFLINE"`
	CurrentLat   float64 `json:"current_lat"`
	CurrentLng   float64 `json:"current_lng"`
	WarehouseID  *int64  `json:"warehouse_id"`
}

type UpdateDriverRequest struct {
	DriverCode   *string  `json:"driver_code"`
	FullName     *string  `json:"full_name"`
	Phone        *string  `json:"phone"`
	VehicleType  *string  `json:"vehicle_type"`
	LicensePlate *string  `json:"license_plate"`
	Status       *string  `json:"status" binding:"omitempty,oneof=AVAILABLE BUSY OFFLINE"`
	CurrentLat   *float64 `json:"current_lat"`
	CurrentLng   *float64 `json:"current_lng"`
	WarehouseID  *int64   `json:"warehouse_id"`
}
