package model

import "github.com/google/uuid"

type Profile struct {
	Email        string
	WhiteLabelID uuid.UUID
	Language     string
	OldID        int32
}
