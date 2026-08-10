package dto

import (
	"time"

	"my-web-app.com/smart-logistic-hub/internal/product/entity"
)

type ProductResponse struct {
	ID         int64     `json:"id"`
	Sku        string    `json:"sku"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	Price      float64   `json:"price"`
	WeightGram int       `json:"weight_gram"`
	LengthCm   float64   `json:"length_cm"`
	WidthCm    float64   `json:"width_cm"`
	HeightCm   float64   `json:"height_cm"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  string    `json:"created_by"`
}

type PaginatedResponse struct {
	Items []ProductResponse `json:"items"`
	Total int               `json:"total"`
	Skip  int               `json:"skip"`
	Limit int               `json:"limit"`
}

func ToResponse(p *entity.Product) ProductResponse {
	return ProductResponse{
		ID:         p.ID,
		Sku:        p.Sku,
		Name:       p.Name,
		Category:   p.Category,
		Price:      p.Price,
		WeightGram: p.WeightGram,
		LengthCm:   p.LengthCm,
		WidthCm:    p.WidthCm,
		HeightCm:   p.HeightCm,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		CreatedBy:  p.CreatedBy,
	}
}

func ToResponseList(products []entity.Product) []ProductResponse {
	res := make([]ProductResponse, len(products))
	for i, p := range products {
		res[i] = ToResponse(&p)
	}
	return res
}
