package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Email    string
	Language string
}

type UserFilter struct {
	ID *uuid.UUID
}
