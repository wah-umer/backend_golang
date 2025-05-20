package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"    json:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"user"             json:"user"`
	Street      string             `json:"street" bson:"street"`
	City        string             `json:"city" bson:"city"`
	State       string             `json:"state" bson:"state"`
	PhoneNumber string             `json:"phoneNumber" bson:"phoneNumber"`
	PostalCode  string             `json:"postalCode" bson:"postalCode"`
	Country     string             `json:"country" bson:"country"`
	Type        string             `json:"type" bson:"type"`
}
