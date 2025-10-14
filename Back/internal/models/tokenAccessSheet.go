package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TokenAccessSheet struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Token string             `json:"token" bson:"token"`
}
