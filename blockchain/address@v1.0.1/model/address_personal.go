package model

import (
	"database/sql"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
)

// CREATE TABLE address_personal
// (
// id         UUID      NOT NULL, -- uuid identifier
// address    TEXT      NOT NULL, -- address (possible not uniq)
// user_uuid  UUID      NOT NULL, -- user
// network    TEXT      NOT NULL, -- network group local identifier
// min_payout DOUBLE PRECISION NOT NULL, -- custom minimum payout ge coin.minpay_mining
// deleted_at TIMESTAMP NULL,     -- deleted date
// updated_at TIMESTAMP NOT NULL, -- updated date
// created_at TIMESTAMP NOT NULL, -- created date
//
// CONSTRAINT address_personal_user_uuid_network_group_uniq UNIQUE (user_uuid, network),
//
// PRIMARY KEY (id)
// );

type AddressPersonal struct {
	Id        uuid.UUID
	Address   string
	UserUuid  uuid.UUID
	Network   nodeCommon.NetworkEnumWrapper
	MinPayout float64
	DeletedAt sql.NullTime
	UpdatedAt time.Time
	CreatedAt time.Time
}

type AddressesPersonal []*AddressPersonal

func NewAddressPersonal(
	addressUuid uuid.UUID,
	address string,
	userUuid uuid.UUID,
	network nodeCommon.NetworkEnum,
	minPayout float64,
) *AddressPersonal {
	dt := time.Now().UTC()
	return &AddressPersonal{
		Id:        addressUuid,
		Address:   address,
		UserUuid:  userUuid,
		Network:   nodeCommon.NewNetworkEnumWrapper(network),
		MinPayout: minPayout,
		DeletedAt: sql.NullTime{Time: time.Time{}, Valid: false},
		UpdatedAt: dt,
		CreatedAt: dt,
	}
}

type AddressPersonalFilter struct {
	Id         *uuid.UUID
	Address    *string
	UserUuid   *uuid.UUID
	Network    *nodeCommon.NetworkEnum
	IsDeleted  *bool
	Pagination *Pagination
}

type AddressPersonalPartial struct {
	Address   *string
	MinPayout *float64
	DeletedAt *sql.NullTime
	UpdatedAt *time.Time
}

func (u *AddressPersonal) Update(partial *AddressPersonalPartial) {
	if partial.Address != nil {
		u.Address = *partial.Address

	}

	if partial.MinPayout != nil {
		u.MinPayout = *partial.MinPayout

	}

	if partial.DeletedAt != nil {
		u.DeletedAt = *partial.DeletedAt

	}

	if partial.UpdatedAt != nil {
		u.UpdatedAt = *partial.UpdatedAt

	}
}
