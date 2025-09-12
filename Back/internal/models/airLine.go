package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AirLine struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name,omitempty"`
	Code          string             `json:"code" bson:"code,omitempty"`
	FileName      string             `json:"fileName" bson:"fileName,omitempty"`
	Active        bool               `json:"active" bson:"active,omitempty"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	IdUserCreated string             `json:"idUserCreated" bson:"idUserCreated,omitempty"`
	IdUserUpdate  string             `json:"idUserUpdate" bson:"idUserUpdate,omitempty"`
}
