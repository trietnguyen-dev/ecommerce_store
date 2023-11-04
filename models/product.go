package models

import "time"

type ProductPrice struct {
	Size     string  `json:"size" bson:"size"`
	Price    float64 `json:"price" bson:"price"`
	Currency string  `json:"currency" bson:"currency"`
}
type ProductType string

const (
	Coffee ProductType = "coffee"
	Bean   ProductType = "bean"
)

type Product struct {
	Id                string         `json:"id" bson:"id"`
	Name              string         `json:"name" bson:"name"`
	ImageLinkSquare   string         `json:"image_link_square" bson:"image_link_square"`
	Type              ProductType    `json:"type" bson:"type"`
	Quantity          int            `json:"quantity" bson:"quantity"`
	CreatedAt         time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at" bson:"updated_at"`
	Description       string         `json:"description" bson:"description"`
	Roasted           string         `json:"roasted" bson:"roasted"`
	ImageLinkPortrait string         `json:"image_link_portrait" bson:"image_link_portrait"`
	Ingredients       string         `json:"ingredients" bson:"ingredients"`
	SpecialIngredient string         `json:"special_ingredient" bson:"special_ingredient"`
	Prices            []ProductPrice `json:"prices" bson:"prices"`
	AverageRating     float64        `json:"average_rating" bson:"average_rating"`
	RatingsCount      float64        `json:"ratings_count" bson:"ratings_count"`
	Favourite         bool           `json:"favourite" bson:"favourite"`
}
