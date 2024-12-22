package model

import "github.com/google/uuid"

type ID struct {
	Old int32
	New uuid.UUID
}
