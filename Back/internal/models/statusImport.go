package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StatusImport struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
	Code string             `json:"code" bson:"code"`
}
