package repository

import (
	"context"

	transactor "code.emcdtech.com/emcd/sdk/pg"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type DerivedFunc func(uint32) (string, error)

type AddressRepository interface {
	AddOldAddress(ctx context.Context, address *model.AddressOld) error
	AddNewCommonAddress(ctx context.Context, address *model.Address) error
	AddNewDerivedAddress(ctx context.Context, address *model.Address, masterKeyId uint32, derivedFunc DerivedFunc) error

	GetOldAddresses(ctx context.Context, addressOldFilter *model.AddressOldFilter) (*uint64, model.AddressesOld, error)
	GetNewAddresses(ctx context.Context, addressFilter *model.AddressFilter) (*uint64, model.Addresses, error)

	AddPersonalAddress(ctx context.Context, address *model.AddressPersonal) error
	GetPersonalAddresses(ctx context.Context, addressFilter *model.AddressPersonalFilter) (*uint64, model.AddressesPersonal, error)
	UpdatePersonalAddress(ctx context.Context, address *model.AddressPersonal, addressPartial *model.AddressPersonalPartial) error

	AddOrUpdateDirtyAddress(ctx context.Context, address *model.AddressDirty) error
	GetDirtyAddresses(ctx context.Context, addressFilter *model.AddressDirtyFilter) (model.AddressesDirty, error)
	// UpdateDirtyAddress(ctx context.Context, address *model.AddressDirty, partial *model.AddressDirtyPartial) error

	transactor.PgxTransactor
}
