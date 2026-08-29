package dto

type CreateTrackingEventRequest struct {
	OrderID      *int64  `json:"order_id"`
	OrderCode    string  `json:"order_code" binding:"required"`
	DriverCode   string  `json:"driver_code" binding:"required"`
	StatusUpdate string  `json:"status_update" binding:"required"`
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
	Note         string  `json:"note"`
}

type UpdateTrackingEventRequest struct {
	OrderID      *int64   `json:"order_id,omitempty"`
	OrderCode    *string  `json:"order_code,omitempty"`
	DriverCode   *string  `json:"driver_code,omitempty"`
	StatusUpdate *string  `json:"status_update,omitempty"`
	Lat          *float64 `json:"lat,omitempty"`
	Lng          *float64 `json:"lng,omitempty"`
	Note         *string  `json:"note,omitempty"`
}
