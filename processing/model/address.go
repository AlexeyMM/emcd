package model

import "github.com/google/uuid"

type Address struct {
	Address    string
	NetworkID  string
	MerchantID uuid.UUID
	Available  bool
}
