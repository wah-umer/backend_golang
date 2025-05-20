package usermodel

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"-" bson:"password"`
	IsVerified bool               `json:"isVerified" bson:"isVerified"`
	IsAdmin    bool               `json:"isAdmin" bson:"isAdmin"`
}
