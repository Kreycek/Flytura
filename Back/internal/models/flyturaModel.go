package models

type EmailRequest struct {
	Email string `json:"email"`
}

type OnlyFlyExcelVerifyExistRequest struct {
	Key string `json:"key"`
}
