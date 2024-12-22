package repository

import (
	"code.emcdtech.com/b2b/swap/model"
	"github.com/google/uuid"
)

// StatusCache отвечает за хранение последних отправленных статусов подписчикам.
// Нужен, что бы не отправлять одни и те же статусы повторно.
type StatusCache interface {
	Add(swapID uuid.UUID, status model.PublicStatus)
	Get(swapID uuid.UUID) model.PublicStatus
	Delete(swapID uuid.UUID)
}
