package models

type EmailRequest struct {
	Email string `json:"email"`
}

type PurcharseRecordVerifyExistRequest struct {
	Key string `json:"key"`
}
