package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID          primitive.ObjectID `json:"_id,omitempty"    bson:"_id,omitempty"`
	User        primitive.ObjectID `json:"user"             bson:"user"`
	Items       []interface{}      `json:"item"             bson:"item"`
	Address     []interface{}      `json:"address"          bson:"address"`
	Status      string             `json:"status"           bson:"status"`
	PaymentMode string             `json:"paymentMode"      bsno:"paymentMode"`
	Total       float64            `json:"total"            bson:"total"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
