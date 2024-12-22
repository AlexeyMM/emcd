package model

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/google/uuid"
)

type CodeTemplate string

func (t CodeTemplate) String() string {
	return string(t)
}

func NewCodeTemplate(s string) CodeTemplate {
	return CodeTemplate(s)
}

const (
	WorkerChangedState    CodeTemplate = "worker changed state"
	PasswordRestoration   CodeTemplate = "password restoration"
	WalletChangedAddress  CodeTemplate = "wallet changed address"
	Register              CodeTemplate = "register"
	UserHashrateDecreased CodeTemplate = "user hashrate decreased"
	MobileTwoFaOff        CodeTemplate = "mobile two fa off"
	MobileTwoFaOn         CodeTemplate = "mobile two fa on"
	PasswordChange        CodeTemplate = "password change"
	PhoneDelete           CodeTemplate = "phone delete"
	ReferralRewardPayouts CodeTemplate = "referral reward payouts"
	WorkerReport          CodeTemplate = "worker report"
	SwapSupport           CodeTemplate = "swap support message"
	SwapMessage           CodeTemplate = "swap message"
	SwapSuccessful        CodeTemplate = "swap successful message"
	IncomeReport          CodeTemplate = "income report"
	GoggleTwoFaOn         CodeTemplate = "google two fa on"
	GoggleTwoFaOff        CodeTemplate = "google two fa off"
	ChangeEmail           CodeTemplate = "change email"
	PayoutReport          CodeTemplate = "payout report"
)

type TextTemplate string

func (t TextTemplate) String() string {
	return string(t)
}

func (t TextTemplate) Render(placeholder any) (string, error) {
	tpl, err := template.New("").Parse(t.String())
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	msg := new(bytes.Buffer)
	err = tpl.Execute(msg, placeholder)
	if err != nil {
		return "", fmt.Errorf("template execute: %w", err)
	}
	return msg.String(), nil
}

func NewTextTemplate(s string) TextTemplate {
	return TextTemplate(s)
}

type Template struct {
	WhiteLabelID uuid.UUID
	Type         CodeTemplate
	Language     string
	Template     TextTemplate
	Subject      TextTemplate
	Footer       string
}

func (t *Template) GetFooter() string {
	const defaultFooter = `` +
		`<p style="display: block; padding: 0; margin: 10px 0 0; ` +
		`font-family: sans-serif; font-size: 12px; line-height: 1.4;">— demo-wl.emcd.io</p>`
	if t.Footer == "" {
		return defaultFooter
	}
	return t.Footer
}
