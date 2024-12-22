package repository

import "github.com/google/uuid"

type ErrorCounter interface {
	Inc(swapID uuid.UUID) int
	Delete(swapID uuid.UUID)
}
