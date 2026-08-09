package tracking

import "time"

type TrackingEvent struct {
	ID           int64     `json:"id" db:"id"`
	OrderCode    string    `json:"order_code" db:"order_code"`
	DriverCode   string    `json:"driver_code" db:"driver_code"`
	StatusUpdate string    `json:"status_update" db:"status_update"`
	Lat          float64   `json:"lat" db:"lat"`
	Lng          float64   `json:"lng" db:"lng"`
	Note         string    `json:"note" db:"note"`
	Timestamp    time.Time `json:"timestamp" db:"timestamp"`
}
