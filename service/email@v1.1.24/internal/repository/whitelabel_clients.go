package repository

import (
	"context"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/whitelabeleventclients.go -package=mockstore

type WhiteLabelEventClients interface {
	SendWLEvent(ctx context.Context, wlID uuid.UUID, req *model.WLEventRequest)
}
