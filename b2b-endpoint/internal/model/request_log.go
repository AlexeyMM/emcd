package model

import (
	"time"

	"github.com/google/uuid"
)

type RequestLog struct {
	ApiKey      uuid.UUID
	RequestHash string
	CreatedAt   time.Time
}

type RequestLogFilter struct {
	ApiKey      *uuid.UUID
	RequestHash *string
}
