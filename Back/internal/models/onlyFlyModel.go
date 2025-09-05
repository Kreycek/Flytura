package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OnlyFlyExcel struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Key      string             `json:"key" bson:"key,omitempty"`
	Name     string             `json:"name" bson:"name,omitempty"`
	LastName string             `json:"lastName" bson:"lastName,omitempty"`
	Status   string             `json:"status" bson:"status,omitempty"`
	FileName string             `json:"fileName" bson:"fileName,omitempty"`
	Value    string             `json:"value" bson:"value,omitempty"`
	// CostCenterSecondary []CostCenterSecondary `json:"costCenterSecondary" bson:"costCenterSecondary,omitempty"`
	Active        bool      `json:"active" bson:"active,omitempty"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
	IdUserCreated string    `json:"idUserCreated" bson:"idUserCreated,omitempty"`
	IdUserUpdate  string    `json:"idUserUpdate" bson:"idUserUpdate,omitempty"`
}
