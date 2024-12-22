package model

import "github.com/google/uuid"

type User struct {
	UserID       uuid.UUID
	WhitelableID uuid.UUID
	ReferralID   uuid.UUID
}
