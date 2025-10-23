package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImagesDB struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FileName    string             `json:"fileName" bson:"fileName"`
	PDFFileName string             `json:"pdfFileName" bson:"pdfFileName"`
	XMLFileName string             `json:"xmlFileName" bson:"xmlFileName"`
	ZipFileName string             `json:"zipFileName" bson:"zipFileName,omitempty"`
	DtImport    time.Time          `json:"dtImport" bson:"dtImport,omitempty"`
	CompanyCode string             `json:"companyCode" bson:"companyCode,omitempty"`
	Key         string             `json:"key" bson:"key,omitempty"`
	CompanyName string             `json:"companyName" bson:"companyName,omitempty"`

	Active       bool      `json:"active" bson:"active,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
	DownloadDone bool      `json:"downloadDone" bson:"downloadDone,omitempty"`
}
