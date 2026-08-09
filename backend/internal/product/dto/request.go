package products

type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type ProductCreate struct {
	Sku        string      `json:"sku" binding:"required"`
	Name       string      `json:"name" binding:"required"`
	Category   string      `json:"category" binding:"required"`
	Price      float64     `json:"price" binding:"required"`
	WeightGram float64     `json:"weight_gram" binding:"required"`
	Dimensions *Dimensions `json:"dimensions,omitempty"`
}

type ProductUpdate struct {
	Name       *string     `json:"name,omitempty"`
	Category   *string     `json:"category,omitempty"`
	Price      *float64    `json:"price,omitempty"`
	WeightGram *float64    `json:"weight_gram,omitempty"`
	Dimensions *Dimensions `json:"dimensions,omitempty"`
}
