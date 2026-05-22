package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PurcharseRecord struct {
	ID                     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Key                    string             `json:"key" bson:"key"`
	Name                   string             `json:"name" bson:"name"`
	LastName               string             `json:"lastName" bson:"lastName"`
	Status                 string             `json:"status" bson:"status"`
	FileName               string             `json:"fileName" bson:"fileName"`
	CompanyName            string             `json:"companyName" bson:"companyName"`
	MessageReturn          string             `json:"messageReturn" bson:"messageReturn"`
	CompanyCode            string             `json:"companyCode" bson:"companyCode"`
	Value                  string             `json:"value" bson:"value"`
	OriginData             string             `json:"originData" bson:"originData"`
	Active                 bool               `json:"active" bson:"active,omitempty"`
	CreatedAt              time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt              time.Time          `json:"updatedAt" bson:"updatedAt"`
	IdUserInserted         string             `json:"idUserInserted" bson:"idUserInserted"`
	EmissionDate           time.Time          `json:"emissionDate" bson:"emissionDate"`
	IdUserUpdate           string             `json:"idUserUpdate" bson:"idUserUpdate"`
	DirectionOfDestination string             `json:"directionOfDestination" bson:"directionOfDestination"`
}
