package entity

import "time"

type Vehicle struct {
	ID           int64     `json:"id" db:"id"`
	LicensePlate string    `json:"license_plate" db:"license_plate"`
	Type         string    `json:"type" db:"type"`
	Capacity     float64   `json:"capacity" db:"capacity"`
	Status       string    `json:"status" db:"status"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy    string    `json:"created_by" db:"created_by"`
}
