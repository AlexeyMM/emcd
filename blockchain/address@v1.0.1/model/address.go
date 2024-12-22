package model

import (
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

// CREATE TABLE address
// (
// id            UUID      NOT NULL, -- uuid identifier
// address       TEXT      NOT NULL, -- address
// user_uuid     UUID      NOT NULL, -- user
// address_type  int4      NOT NULL, -- type of address generation
// network_group TEXT      NOT NULL, -- network group local identifier
// created_at    TIMESTAMP NOT NULL, -- created date
//
// CONSTRAINT address_user_uuid_network_group_idx UNIQUE (user_uuid, network_group),
// CONSTRAINT address_address_idx UNIQUE (address),
//
// PRIMARY KEY (id)
// );

type Address struct {
	Id             uuid.UUID
	Address        string
	UserUuid       uuid.UUID
	ProcessingUuid uuid.UUID
	AddressType    enum.AddressTypeWrapper
	NetworkGroup   nodeCommon.NetworkGroupEnumWrapper
	CreatedAt      time.Time

	addressDerived AddressDerived
}

type Addresses []*Address

func NewAddress(
	addressUuid uuid.UUID,
	address string,
	userUuid uuid.UUID,
	addressType addressPb.AddressType,
	networkGroup nodeCommon.NetworkGroupEnum,
) *Address {

	return NewProcessingAddress(addressUuid, address, userUuid, userUuid, addressType, networkGroup)
}

func NewProcessingAddress(
	addressUuid uuid.UUID,
	address string,
	userUuid uuid.UUID,
	processingUuid uuid.UUID,
	addressType addressPb.AddressType,
	networkGroup nodeCommon.NetworkGroupEnum,
) *Address {

	return &Address{
		Id:             addressUuid,
		Address:        address,
		UserUuid:       userUuid,
		ProcessingUuid: processingUuid,
		AddressType:    enum.NewAddressTypeWrapper(addressType),
		NetworkGroup:   nodeCommon.NewNetworkGroupEnumWrapper(networkGroup),
		CreatedAt:      time.Now().UTC(),
		addressDerived: AddressDerived{
			AddressUuid:   uuid.UUID{},
			NetworkGroup:  nodeCommon.NewNetworkGroupEnumWrapper(networkGroup),
			MasterKeyId:   0,
			DerivedOffset: 0,
		},
	}
}

type AddressFilter struct {
	Id           *uuid.UUID
	Address      *string
	UserUuid     *uuid.UUID
	IsProcessing *bool
	AddressType  *addressPb.AddressType
	NetworkGroup *nodeCommon.NetworkGroupEnum
	CreatedAtGt  *time.Time
	Pagination   *Pagination
}

type AddressPartial struct {
	Address *string
}

func (u *Address) Update(partial *AddressPartial) {
	if partial.Address != nil {
		u.Address = *partial.Address

	}
}

func (u *Address) SetAddressDerived(addressDerived *AddressDerived) {
	u.addressDerived = *addressDerived

}

func (u *Address) GetAddressDerived() *AddressDerived {

	return &u.addressDerived
}

func (us Addresses) GetDerivedUuids() uuid.UUIDs {
	var dumps uuid.UUIDs

	for _, u := range us {
		if u.AddressType.Number() == addressPb.AddressType_ADDRESS_TYPE_DERIVED.Number() {
			dumps = append(dumps, u.Id)

		}
	}

	return dumps
}
