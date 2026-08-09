package entity

import "time"

type Product struct {
	ID         int64     `json:"id" db:"id"`
	Sku        string    `json:"sku" db:"sku"`
	Name       string    `json:"name" db:"name"`
	Category   string    `json:"category" db:"category"`
	Price      float64   `json:"price" db:"price"`
	WeightGram int       `json:"weight_gram" db:"weight_gram"`
	LengthCm   float64   `json:"length_cm" db:"length_cm"`
	WidthCm    float64   `json:"width_cm" db:"width_cm"`
	HeightCm   float64   `json:"height_cm" db:"height_cm"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
}
