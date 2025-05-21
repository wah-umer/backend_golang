package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID        primitive.ObjectID `json:"_id,omitempty"    bson:"_id,omitempty"`
	User      primitive.ObjectID `json:"user"             bson:"user"`
	Product   primitive.ObjectID `json:"product"          bson:"product"`
	Rating    int                `json:"rating"           bson:"rating"`
	Comment   string             `json:"comment"          bson:"comment"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
