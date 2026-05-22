package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Perfil struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	ShortName string             `json:"shortName" bson:"shortName"`
}
