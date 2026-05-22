package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Module           string             `json:"module" bson:"module,omitempty"`
	Class            string             `json:"class" bson:"class,omitempty"`
	Method           string             `json:"method" bson:"method,omitempty"`
	ErrorDescription string             `json:"errorDescription" bson:"errorDescription,omitempty"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
}
