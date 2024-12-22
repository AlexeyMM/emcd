package model

import "time"

type Worker struct {
	Name           string
	Username       string
	Coin           string
	IsOn           bool
	StateChangedAt time.Time
}
