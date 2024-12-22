package repository

import (
	"context"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/email_messages.go -package=mockstore

type EmailMessages interface {
	Create(ctx context.Context, em *model.EmailMessageEvent) error
	ListMessages(
		ctx context.Context,
		email, eventType *string,
		skip, take int32) ([]*model.EmailMessageEvent, int, error)
}
