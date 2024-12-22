package model

import (
	"time"
)

type KycStatus int

const (
	Unknown KycStatus = iota
	Processing
	Approved
	Declined
)

type Kyc struct {
	RetryAfter   time.Time
	DelayMinutes int
	Status       KycStatus
	IsAllowed    bool
	Overall      string
	DocCheck     string
	FaceCheck    string
}
