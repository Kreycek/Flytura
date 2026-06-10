package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OutPutInvoices struct {
	ID                   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Key                  string             `json:"key" bson:"key"`
	Status               string             `json:"status" bson:"status"`
	DtProcess            time.Time          `json:"dtProcess" bson:"dtProcess"`
	MonthProcess         int                `json:"monthProcess" bson:"monthProcess"`
	DtFLy                string             `json:"dtFly" bson:"dtFly"`
	RFC                  string             `json:"rfc" bson:"rfc"`
	TransferredBaseValue float64            `json:"transferredBaseValue" bson:"transferredBaseValue"`
	IVAValue             float64            `json:"ivaValue" bson:"ivaValue"`
	Tax                  string             `json:"tax" bson:"tax"`
	Rate                 string             `json:"rate" bson:"rate"`
	FactorType           string             `json:"factorType" bson:"factorType"`
	SubTotalValue        float64            `json:"subTotalValue" bson:"subTotalValue"`
	TotalValue           float64            `json:"totalValue" bson:"totalValue"`
	CreatedAt            time.Time          `json:"createdAt" bson:"createdAt"`
	Active               bool               `json:"active" bson:"active"`
	Ruta                 string             `json:"ruta" bson:"ruta"`
	CompanyName          string             `json:"companyName" bson:"companyName"`
	CompanyCode          string             `json:"companyCode" bson:"companyCode"`
	TUA                  float64            `json:"tua" bson:"tua"`
	OtherValues          float64            `json:"otherValues" bson:"otherValues"`
	Segment              int                `json:"segment" bson:"segment"`
	Ticket               string             `json:"ticket" bson:"ticket"`
	IdUserInserted       string             `json:"idUserInserted" bson:"idUserInserted"`
	UserNameImport       string             `json:"userNameImport" bson:"userNameImport"`
	OriginData           string             `json:"originData" bson:"originData"`
	ServerDate           time.Time          `json:"serverDate" bson:"serverDate"`
}
