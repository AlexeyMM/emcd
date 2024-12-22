package model

import (
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

// CREATE TABLE address_old
// (
// id              UUID      NOT NULL, -- uuid identifier
// address         TEXT      NOT NULL, -- address
// user_uuid       UUID      NOT NULL, -- user (for tables unification)
// address_type    int4      NOT NULL, -- type of address generation
// network         TEXT      NOT NULL, -- network
// user_account_id int4      NOT NULL, -- wallets user account id
// coin            TEXT      NOT NULL, -- coin
// created_at      TIMESTAMP NOT NULL, -- created date
//
// CONSTRAINT address_old_network_user_account_id_idx UNIQUE (network, user_account_id),
// CONSTRAINT address_old_address_idx UNIQUE (address),
//
// PRIMARY KEY (id)
// );

type AddressOld struct {
	Id            uuid.UUID
	Address       string
	UserUuid      uuid.UUID
	AddressType   enum.AddressTypeWrapper
	Network       nodeCommon.NetworkEnumWrapper
	UserAccountId int32
	Coin          string
	CreatedAt     time.Time
}

type AddressesOld []*AddressOld

func NewAddressOld(
	addressUuid uuid.UUID,
	address string,
	userUuid uuid.UUID,
	addressType addressPb.AddressType,
	network nodeCommon.NetworkEnum,
	userAccountId int32,
	coin string,
) *AddressOld {

	return &AddressOld{
		Id:            addressUuid,
		Address:       address,
		UserUuid:      userUuid,
		AddressType:   enum.NewAddressTypeWrapper(addressType),
		Network:       nodeCommon.NewNetworkEnumWrapper(network),
		UserAccountId: userAccountId,
		Coin:          coin,
		CreatedAt:     time.Now().UTC(),
	}
}

type AddressOldFilter struct {
	Id            *uuid.UUID
	Address       *string
	UserUuid      *uuid.UUID
	AddressType   *addressPb.AddressType
	Network       *nodeCommon.NetworkEnum
	UserAccountId *int32
	Coin          *string
	CreatedAtGt   *time.Time
	Pagination    *Pagination
}
