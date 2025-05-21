package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cart struct {
	ID       primitive.ObjectID `json:"_id,omitempty"    bson:"_id,omitempty"`
	User     primitive.ObjectID `json:"user"             bson:"user"`
	Product  primitive.ObjectID `json:"product"          bson:"product"`
	Quantity int                `json:"quantity"         bson:"quantity"`
}
