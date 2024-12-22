package repository

import (
	"context"

	emailPb "code.emcdtech.com/emcd/service/email/protocol/email"
	"github.com/google/uuid"
)

type Email interface {
	SendWalletChangedAddress(ctx context.Context, userID uuid.UUID, token, coinCode string) error
}

type email struct {
	cli emailPb.EmailServiceClient
}

func NewEmail(cli emailPb.EmailServiceClient) Email {
	return &email{
		cli: cli,
	}
}

func (e *email) SendWalletChangedAddress(ctx context.Context, userID uuid.UUID, token, coinCode string) error {
	_, err := e.cli.SendWalletChangedAddress(ctx, &emailPb.SendWalletChangedAddressRequest{
		Token:    token,
		CoinCode: coinCode,
		UserId:   userID.String(),
	})
	return err
}
