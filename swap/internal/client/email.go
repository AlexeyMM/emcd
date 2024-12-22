package client

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/emcd/service/email/protocol/email"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

// https://emcd.io/b2b-swap/transaction/36ed70b9-9a37-4d3a-aa73-ee2b5e34fba5
// https://stage.mytstnv.site/b2b-swap/transaction/bcab4747-70db-49fa-93fd-5633940593f5

type Email interface {
	SendInitialSwapMessage(ctx context.Context, userEmail, language string, swapID uuid.UUID)
	SendSuccessfulSwapMessage(ctx context.Context, swap *model.Swap, user *model.User)
}

type EmailImp struct {
	emailCli    email.EmailServiceClient
	environment string
}

func NewEmail(emailCli email.EmailServiceClient, environment string) *EmailImp {
	return &EmailImp{
		emailCli:    emailCli,
		environment: environment,
	}
}

func (e *EmailImp) SendInitialSwapMessage(ctx context.Context, userEmail, language string, swapID uuid.UUID) {
	var link string
	if e.environment == "production" {
		link = fmt.Sprintf("https://emcd.io/b2b-swap/transaction/%s", swapID.String())
	} else if e.environment == "development" {
		link = fmt.Sprintf("https://stage.mytstnv.site/b2b-swap/transaction/%s", swapID.String())
	}
	if link == "" {
		log.Error("admin: sendInitialSwapMessage: link doesn't found")
		return
	}

	_, err := e.emailCli.SendInitialSwapMessage(ctx, &email.SendInitialSwapMessageRequest{
		Email:    userEmail,
		Language: language,
		Link:     link,
	})
	if err != nil {
		log.Error(ctx, "email: sendInitialSwapMessage: %s", err.Error())
		return
	}
}

func (e *EmailImp) SendSuccessfulSwapMessage(ctx context.Context, swap *model.Swap, user *model.User) {
	from := fmt.Sprintf("%s %s", swap.AmountFrom.String(), swap.CoinFrom)
	to := fmt.Sprintf("%s %s", swap.AmountTo.String(), swap.CoinTo)
	formatTime := swap.EndTime.Format("02/01/2006, 15:04:05")

	_, err := e.emailCli.SendSuccessfulSwapMessage(ctx, &email.SendSuccessfulSwapMessageRequest{
		Email:         user.Email,
		Language:      user.Language,
		SwapId:        swap.ID.String(),
		From:          from,
		To:            to,
		Address:       swap.AddressTo,
		ExecutionTime: formatTime,
	})
	if err != nil {
		log.Error("email: sendSuccessfulSwapMessage: %w", err.Error())
		return
	}
}
