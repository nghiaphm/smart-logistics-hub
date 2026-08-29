package dto

type CreateVehicleRequest struct {
	LicensePlate string  `json:"license_plate" binding:"required"`
	Type         string  `json:"type"`
	Capacity     float64 `json:"capacity"`
	Status       string  `json:"status" binding:"omitempty,oneof=ACTIVE MAINTENANCE INACTIVE"`
}

type UpdateVehicleRequest struct {
	LicensePlate *string  `json:"license_plate"`
	Type         *string  `json:"type"`
	Capacity     *float64 `json:"capacity"`
	Status       *string  `json:"status" binding:"omitempty,oneof=ACTIVE MAINTENANCE INACTIVE"`
}
