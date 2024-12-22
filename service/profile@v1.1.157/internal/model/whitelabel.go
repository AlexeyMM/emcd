package model

import "github.com/google/uuid"

type WhiteLabel struct {
	ID                    uuid.UUID
	Prefix                string
	Version               int
	IsEmailConfirmEnabled bool
}
