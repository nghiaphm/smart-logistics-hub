package models

import "time"

type Product struct {
	ID         string                 `json:"_id,omitempty" bson:"_id,omitempty"`
	Sku        string                 `json:"sku" bson:"sku"`
	Name       string                 `json:"name" bson:"name"`
	Category   string                 `json:"category" bson:"category"`
	Price      float64                `json:"price" bson:"price"`
	WeightGram float64                `json:"weight_gram" bson:"weight_gram"`
	Dimensions map[string]interface{} `json:"dimensions,omitempty" bson:"dimensions,omitempty"`
	CreatedAt  *time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt  *time.Time             `json:"updated_at" bson:"updated_at"`
	CreatedBy  *string                `json:"created_by,omitempty" bson:"created_by,omitempty"`
}
