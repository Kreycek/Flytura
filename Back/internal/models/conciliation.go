package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conciliation struct {
	ID                         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Protocol                   string             `json:"protocol" bson:"protocol"`
	EmissionDate               time.Time          `json:"emissionDate" bson:"emissionDate"`
	OriginDate                 time.Time          `json:"originDate" bson:"originDate"`
	ReturnDate                 time.Time          `json:"returnDate" bson:"returnDate"`
	OriginLocator              string             `json:"originLocator" bson:"originLocator"`
	ReturnLocator              string             `json:"returnLocator" bson:"returnLocator"`
	OriginETicket              string             `json:"originETicket" bson:"originETicket"`
	ReturnETicket              string             `json:"returnETicket" bson:"returnETicket"`
	OriginAirline              string             `json:"originAirline" bson:"originAirline"`
	ReturnAirline              string             `json:"returnAirline" bson:"returnAirline"`
	TravelerName               string             `json:"travelerName" bson:"travelerName"`
	TravelerFirstName          string             `json:"travelerFirstName" bson:"travelerFirstName"`
	TravelerLastName           string             `json:"travelerLastName" bson:"travelerLastName"`
	BookingStatus              string             `json:"bookingStatus" bson:"bookingStatus"`
	BookingStatusOld           string             `json:"bookingStatusOld" bson:"bookingStatusOld"`
	CancelledReason            string             `json:"cancelledReason" bson:"cancelledReason"`
	CurrencyCode               string             `json:"currencyCode" bson:"currencyCode"`
	AmountOrigin               float64            `json:"amountOrigin" bson:"amountOrigin"`
	AmountReturn               float64            `json:"amountReturn" bson:"amountReturn"`
	OriginCancelledReason      string             `json:"originCancelledReason" bson:"originCancelledReason"`
	ReturnCancelledReason      string             `json:"returnCancelledReason" bson:"returnCancelledReason"`
	Active                     bool               `json:"active" bson:"active"`
	CreatedAtContractedCountry time.Time          `json:"createdAtContractedCountry" bson:"createdAtContractedCountry"`
	CreatedAtLocalCountry      time.Time          `json:"createdAtLocalCountry" bson:"createdAtLocalCountry"`
	FlightOrigin               string             `json:"flightOrigin" bson:"flightOrigin"`
	FlightDestination          string             `json:"flightDestination" bson:"flightDestination"`
	OriginCountryCode          string             `json:"originCountryCode" bson:"originCountryCode"`
	OriginCity                 string             `json:"originCity" bson:"originCity"`
	DestinationCountryCode     string             `json:"destinationCountryCode" bson:"destinationCountryCode"`
	DestinationCity            string             `json:"destinationCity" bson:"destinationCity"`
	OriginAirlineCommercial    string             `json:"originAirlineCommercial" bson:"originAirlineCommercial"`
	ReturnAirlineCommercial    string             `json:"returnAirlineCommercial" bson:"returnAirlineCommercial"`
}
