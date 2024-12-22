package model

import (
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
)

// CREATE TABLE address_dirty
// (
// address         TEXT      NOT NULL, -- address
// network         TEXT      NOT NULL, -- network
// is_dirty        BOOLEAN   NOT NULL, -- is_dirty
// updated_at      TIMESTAMP NOT NULL, -- updated date
// created_at      TIMESTAMP NOT NULL, -- created date
//
// CONSTRAINT address_dirty_address_network_uniq UNIQUE (address, network)
//
// );

type AddressDirty struct {
	Address   string
	Network   nodeCommon.NetworkEnumWrapper
	IsDirty   bool
	UpdatedAt time.Time
	CreatedAt time.Time
}

type AddressesDirty []*AddressDirty

func NewAddressDirty(address string, network nodeCommon.NetworkEnum, isDirty bool, currentTime time.Time) *AddressDirty {

	return &AddressDirty{
		Address:   address,
		Network:   nodeCommon.NewNetworkEnumWrapper(network),
		IsDirty:   isDirty,
		UpdatedAt: currentTime,
		CreatedAt: currentTime,
	}
}

type AddressDirtyFilter struct {
	Address *string
	Network *nodeCommon.NetworkEnum
}

// type AddressDirtyPartial struct {
// 	IsDirty *bool
// }
//
// func (u *AddressDirty) Update(partial *AddressDirtyPartial) {
// 	if partial.IsDirty != nil {
// 		u.IsDirty = *partial.IsDirty
//
// 	}
// }
