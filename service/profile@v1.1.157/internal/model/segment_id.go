package model

import (
	"github.com/google/uuid"
)

type UserInfoBySegmentID struct {
	ID       int32
	UUID     uuid.UUID
	UserName string
	Email    string
}
