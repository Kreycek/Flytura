package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImagesDB struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FileName        string             `json:"fileName" bson:"fileName"`
	PDFFileName     string             `json:"pdfFileName" bson:"pdfFileName"`
	DownloadPDFDone bool               `json:"downloadPDFDone" bson:"downloadPDFDone"`
	XMLFileName     string             `json:"xmlFileName" bson:"xmlFileName"`
	DownloadXMLDone bool               `json:"downloadXMLDone" bson:"downloadXMLDone"`
	ZipFileName     string             `json:"zipFileName" bson:"zipFileName,omitempty"`
	DtImport        time.Time          `json:"dtImport" bson:"dtImport,omitempty"`
	CompanyCode     string             `json:"companyCode" bson:"companyCode,omitempty"`
	Key             string             `json:"key" bson:"key,omitempty"`
	CompanyName     string             `json:"companyName" bson:"companyName,omitempty"`
	Active          bool               `json:"active" bson:"active,omitempty"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	DownloadDone    bool               `json:"downloadDone" bson:"downloadDone"`
	BilledFlytura   bool               `json:"billedFlytura" bson:"billedFlytura,omitempty"`
	IdUserInserted  string             `json:"idUserInserted" bson:"idUserInserted"`
	UserNameImport  string             `json:"userNameImport" bson:"userNameImport"`
	OriginData      string             `json:"originData" bson:"originData"`
	ServerDate      time.Time          `json:"serverDate" bson:"serverDate"`
}
