package model

import (
	"time"

	"github.com/google/uuid"
)

type Secret struct {
	ClientID  uuid.UUID
	ApiKey    uuid.UUID
	ApiSecret uuid.UUID
	IsActive  bool
	CreatedAt time.Time
	LastUsed  time.Time
}

type SecretFilter struct {
	ApiKey   *uuid.UUID
	ClientID *uuid.UUID
	IsActive *bool
}

type SecretPartial struct {
	IsActive *bool
	LastUsed *time.Time
}

func (s *Secret) Update(partial *SecretPartial) {
	if partial.IsActive != nil {
		s.IsActive = *partial.IsActive
	}
	if partial.LastUsed != nil {
		s.LastUsed = *partial.LastUsed
	}
}
