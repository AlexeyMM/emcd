package model

import "github.com/google/uuid"

type Subscriber struct {
	ClientID uuid.UUID
	Ch       chan PublicStatus
}
