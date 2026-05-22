package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OutPutInvoices struct {
	ID                   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Key                  string             `json:"key" bson:"key,omitempty"`
	Status               string             `json:"status" bson:"status,omitempty"`
	DtProcess            time.Time          `json:"dtProcess" bson:"dtProcess,omitempty"`
	MonthProcess         int                `json:"monthProcess" bson:"monthProcess,omitempty"`
	DtFLy                string             `json:"dtFly" bson:"dtFly,omitempty"`
	RFC                  string             `json:"rfc" bson:"rfc,omitempty"`
	TransferredBaseValue float64            `json:"transferredBaseValue" bson:"transferredBaseValue,omitempty"`
	IVAValue             float64            `json:"ivaValue" bson:"ivaValue,omitempty"`
	Tax                  string             `json:"tax" bson:"tax,omitempty"`
	Rate                 string             `json:"rate" bson:"rate,omitempty"`
	FactorType           string             `json:"factorType" bson:"factorType,omitempty"`
	SubTotalValue        float64            `json:"subTotalValue" bson:"subTotalValue,omitempty"`
	TotalValue           float64            `json:"totalValue" bson:"totalValue,omitempty"`
	CreatedAt            time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	Active               bool               `json:"active" bson:"active,omitempty"`
	Ruta                 string             `json:"ruta" bson:"ruta,omitempty"`
	CompanyName          string             `json:"companyName" bson:"companyName,omitempty"`
	CompanyCode          string             `json:"companyCode" bson:"companyCode,omitempty"`
	TUA                  float64            `json:"tua" bson:"tua,omitempty"`
	OtherValues          float64            `json:"otherValues" bson:"otherValues,omitempty"`
	Segment              int                `json:"segment" bson:"segment,omitempty"`
	Ticket               string             `json:"ticket" bson:"ticket,omitempty"`
}
