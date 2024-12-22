package model

import "github.com/google/uuid"

type EmailSettings struct {
	WhiteLabelID        uuid.UUID
	PasswordRestoration *PasswordRestorationSetting
}

type PasswordRestorationSetting struct {
	WhiteLabelID uuid.UUID
	Sender       string
	Title        string
	Body         string
	Login        string
	Password     string
	Domain       string
	ApiKey       string
	RedirectUrl  string
	Provider     string
}
