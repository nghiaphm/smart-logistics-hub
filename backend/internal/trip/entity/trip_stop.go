package entity

import "time"

type TripStop struct {
	ID          int64      `json:"id" db:"id"`
	TripID      int64      `json:"trip_id" db:"trip_id"`
	OrderCode   string     `json:"order_code" db:"order_code"`
	StopType    string     `json:"stop_type" db:"stop_type"`
	Address     string     `json:"address" db:"address"`
	Lat         float64    `json:"lat" db:"lat"`
	Lng         float64    `json:"lng" db:"lng"`
	Status      string     `json:"status" db:"status"`
	PlannedAt   *time.Time `json:"planned_at" db:"planned_at"`
	ArrivedAt   *time.Time `json:"arrived_at" db:"arrived_at"`
	DepartureAt *time.Time `json:"departure_at" db:"departure_at"`
}
