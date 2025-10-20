package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImagesDB struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FileName     string             `json:"fileName" bson:"fileName"`
	DtImport     time.Time          `json:"dtImport" bson:"dtImport,omitempty"`
	CompanyCode  string             `json:"companyCode" bson:"companyCode,omitempty"`
	CompanyName  string             `json:"companyName" bson:"companyName,omitempty"`
	FileURL      string             `json:"fileURL" bson:"fileURL,omitempty"`
	Active       bool               `json:"active" bson:"active,omitempty"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	DownloadDone bool               `json:"downloadDone" bson:"downloadDone,omitempty"`
}
