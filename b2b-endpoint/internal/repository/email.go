package repository

import (
	"context"
	"fmt"

	emailCli "code.emcdtech.com/emcd/service/email/protocol/email"
)

type Email interface {
	SendSupportMessage(ctx context.Context, name, email, text string) error
}

type EmailImp struct {
	cli emailCli.EmailServiceClient
}

func NewEmailImp(cli emailCli.EmailServiceClient) *EmailImp {
	return &EmailImp{
		cli: cli,
	}
}

func (e *EmailImp) SendSupportMessage(ctx context.Context, name, email, text string) error {
	_, err := e.cli.SendSwapSupportMessage(ctx, &emailCli.SendSwapSupportMessageRequest{
		Name:      name,
		UserEmail: email,
		Text:      text,
	})
	if err != nil {
		return fmt.Errorf("sendSwapSupportMessage: %w", err)
	}
	return nil
}
