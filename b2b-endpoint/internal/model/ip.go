package model

import (
	"time"

	"github.com/google/uuid"
)

type IP struct {
	ApiKey    uuid.UUID
	Address   string
	CreatedAt time.Time
}

type IPFilter struct {
	ApiKey  *uuid.UUID
	Address *string
}
