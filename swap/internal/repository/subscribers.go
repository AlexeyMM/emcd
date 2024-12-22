package repository

import (
	"code.emcdtech.com/b2b/swap/model"
	"github.com/google/uuid"
)

type Subscribers interface {
	Add(swapID, clientID uuid.UUID, ch chan model.PublicStatus)
	GetBySwapID(swapID uuid.UUID) []*model.Subscriber
	Delete(swapID, clientID uuid.UUID)
}
