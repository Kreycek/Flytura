package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Exercise struct {
	Year       int       `json:"year" bson:"year"`
	StartMonth string    `json:"startMonth" bson:"startMonth"`
	EndMonth   string    `json:"endMonth" bson:"endMonth"`
	DtAdd      time.Time `json:"dtAdd" bson:"dtAdd"`
}

type Movements struct {
	Order          int             `json:"order" bson:"order,omitempty"`
	Date           time.Time       `json:"date" bson:"date,omitempty"`
	Description    string          `json:"description" bson:"description,omitempty"`
	Active         bool            `json:"active" bson:"active"`
	MovementsItens []MovementItens `json:"movementsItens" bson:"movementsItens,omitempty"`
}

type MovementItens struct {
	Date          time.Time            `json:"date" bson:"date"`
	CodAccount    string               `json:"codAccount" bson:"codAccount,omitempty"`
	DebitValue    primitive.Decimal128 `json:"debitValue" bson:"debitValue,omitempty"`
	CreditValue   primitive.Decimal128 `json:"creditValue" bson:"creditValue,omitempty"`
	Active        bool                 `json:"active" bson:"active,omitempty"`
	CodAccountIva string               `json:"codAccountIva" bson:"codAccountIva,omitempty"`
}

type CostCenterSecondary struct {
	CodCostCenterSecondary string `json:"codCostCenterSecondary" bson:"codCostCenterSecondary"`
	Description            string `json:"description" bson:"description"`
}

type CostCenterCOA struct {
	IdCostCenter         string   `json:"idCostCenter" bson:"idCostCenter"`
	CostCentersSecondary []string `json:"costCentersSecondary" bson:"costCentersSecondary"`
}
