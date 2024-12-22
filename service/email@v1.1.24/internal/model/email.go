package model

import (
	"time"

	"github.com/google/uuid"
)

type EmailMessageEvent struct {
	ID        uuid.UUID
	Email     string
	Type      CodeTemplate
	CreatedAt time.Time
}

type Attachment struct {
	Filename string
	Data     []byte
}

type Mail struct {
	To          string
	Subject     string
	Message     string
	Attachments []Attachment
}
