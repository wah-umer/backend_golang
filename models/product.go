package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID                 primitive.ObjectID `json:"_id,omitempty"          bson:"_id,omitempty"`
	Title              string             `json:"title"                  bson:"title"`
	Description        string             `json:"description"            bson:"description"`
	Price              float64            `json:"price"                  bson:"price"`
	DiscountPercentage float64            `json:"discountPercentage"     bson:"discountPercentage"`
	Category           primitive.ObjectID `json:"category"               bson:"category"`
	Brand              primitive.ObjectID `json:"brand"                  bson:"brand"`
	StockQuantity      int                `json:"stockQuantity"          bson:"stockQuantity"`
	Thumbnail          string             `json:"thumbnail"              bson:"thumbnail"`
	Images             []string           `json:"images"                 bson:"images"`
	IsDeleted          bool               `json:"isDeleted"              bson:"isDeleted"`
	UpdatedAt          time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	CreatedAt          time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"  `
}
