package model

import (
	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
)

// CREATE TABLE address_derived
// (
// address_uuid   uuid NOT NULL REFERENCES address (id), -- address uuid
// network_group  TEXT NOT NULL,                         -- network group local identifier
// master_key_id  int4 NOT NULL,                         -- key id == 1 (reserved column)
// derived_offset int8 NOT NULL,                         -- offset of derived address
//
// CONSTRAINT address_derived_helper_key_offset_network_idx UNIQUE (master_key_id, derived_offset, network_group)
// );

type AddressDerived struct {
	AddressUuid   uuid.UUID
	NetworkGroup  nodeCommon.NetworkGroupEnumWrapper
	MasterKeyId   uint32
	DerivedOffset uint32
}

type AddressDeriveds []*AddressDerived

func NewAddressDerived(addressUuid uuid.UUID, networkGroup nodeCommon.NetworkGroupEnum, masterKeyId uint32) *AddressDerived {

	return &AddressDerived{
		AddressUuid:   addressUuid,
		NetworkGroup:  nodeCommon.NewNetworkGroupEnumWrapper(networkGroup),
		MasterKeyId:   masterKeyId,
		DerivedOffset: 0,
	}
}
