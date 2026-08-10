package dto

type Dimensions struct {
	Length *float64 `json:"length,omitempty"`
	Width  *float64 `json:"width,omitempty"`
	Height *float64 `json:"height,omitempty"`
}

type CreateProductRequest struct {
	Sku        string      `json:"sku" binding:"required"`
	Name       string      `json:"name" binding:"required"`
	Category   string      `json:"category"`
	Price      float64     `json:"price"`
	WeightGram int         `json:"weight_gram"`
	Dimensions *Dimensions `json:"dimensions,omitempty"`
	CreatedBy  string      `json:"created_by"`
}

type UpdateProductRequest struct {
	Name       *string     `json:"name,omitempty"`
	Category   *string     `json:"category,omitempty"`
	Price      *float64    `json:"price,omitempty"`
	WeightGram *int        `json:"weight_gram,omitempty"`
	Dimensions *Dimensions `json:"dimensions,omitempty"`
	CreatedBy  string      `json:"created_by"`
}
