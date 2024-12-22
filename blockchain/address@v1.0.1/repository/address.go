package repository

import (
	"context"
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

type AddressRepository interface {
	GetOrCreateAddress(ctx context.Context, userUuid uuid.UUID, network nodeCommon.NetworkEnum, coin *string) (*model.Address, *model.AddressOld, error)
	CreateProcessingAddress(ctx context.Context, userUuid, processingUuid uuid.UUID, network nodeCommon.NetworkEnum) (*model.Address, *model.AddressOld, error)
	GetAddressByUuid(ctx context.Context, addressUuid uuid.UUID) (*model.Address, *model.AddressOld, error)
	GetAddressByStr(ctx context.Context, addressStr string) (*model.Address, *model.AddressOld, error)
	GetAddressesByUserUuid(ctx context.Context, userUuid uuid.UUID) (model.Addresses, model.AddressesOld, error)
	GetAddressesOldByFilter(ctx context.Context, filter *model.AddressOldFilter) (*uint64, model.AddressesOld, error)
	GetAddressesNewByFilter(ctx context.Context, filter *model.AddressFilter) (*uint64, model.Addresses, error)
	AddOrUpdatePersonalAddress(ctx context.Context, addressStr string, userUuid uuid.UUID, network nodeCommon.NetworkEnum, minPayout *float64) (*model.AddressPersonal, error)
	DeletePersonalAddress(ctx context.Context, userUuid uuid.UUID, network nodeCommon.NetworkEnum) error
	GetPersonalAddressesByFilter(ctx context.Context, filter *model.AddressPersonalFilter) (*uint64, model.AddressesPersonal, error)
	GetPersonalAddressByUserUuid(ctx context.Context, userUuid uuid.UUID) (model.AddressesPersonal, error)
	CreateOrUpdateDirtyAddress(ctx context.Context, address *model.AddressDirty) (*model.AddressDirty, error)
	GetDirtyAddressesByFilter(ctx context.Context, filter *model.AddressDirtyFilter) (model.AddressesDirty, error)
}

type addressRepositoryImp struct {
	handler addressPb.AddressServiceClient
}

func NewAddressRepository(
	handler addressPb.AddressServiceClient,
) AddressRepository {

	return &addressRepositoryImp{
		handler: handler,
	}
}

func (r *addressRepositoryImp) GetOrCreateAddress(ctx context.Context, userUuid uuid.UUID, network nodeCommon.NetworkEnum, coin *string) (*model.Address, *model.AddressOld, error) {
	req := &addressPb.CreateAddressRequest{
		UserUuid: userUuid.String(),
		Network:  network.ToString(),
		Coin:     coin,
	}

	if resp, err := r.handler.GetOrCreateAddress(ctx, req); err != nil {

		return nil, nil, fmt.Errorf("handler get or create: %w", err)
	} else if addressNew, addressOld, err := mapping.MapProtoAddressResponsesToModel(resp); err != nil {

		return nil, nil, fmt.Errorf("mapping address response: %w", err)
	} else {

		return addressNew, addressOld, nil
	}
}

func (r *addressRepositoryImp) CreateProcessingAddress(ctx context.Context, userUuid, processingUuid uuid.UUID, network nodeCommon.NetworkEnum) (*model.Address, *model.AddressOld, error) {
	req := &addressPb.CreateProcessingAddressRequest{
		UserUuid:       userUuid.String(),
		Network:        network.ToString(),
		ProcessingUuid: processingUuid.String(),
	}

	if resp, err := r.handler.CreateProcessingAddress(ctx, req); err != nil {

		return nil, nil, fmt.Errorf("handler processing address: %w", err)
	} else if addressNew, addressOld, err := mapping.MapProtoAddressResponsesToModel(resp); err != nil {

		return nil, nil, fmt.Errorf("mapping address response: %w", err)
	} else {

		return addressNew, addressOld, nil
	}
}

func (r *addressRepositoryImp) GetAddressByUuid(ctx context.Context, addressUuid uuid.UUID) (*model.Address, *model.AddressOld, error) {
	req := &addressPb.AddressUuid{
		AddressUuid: addressUuid.String(),
	}

	if resp, err := r.handler.GetAddressByUuid(ctx, req); err != nil {

		return nil, nil, fmt.Errorf("handler get address by uuid: %w", err)
	} else if addressesNew, addressesOld, err := mapping.MapProtoAddressResponsesToModel(resp); err != nil {

		return nil, nil, fmt.Errorf("mapping address response: %w", err)
	} else {

		return addressesNew, addressesOld, nil
	}
}

func (r *addressRepositoryImp) GetAddressByStr(ctx context.Context, addressStr string) (*model.Address, *model.AddressOld, error) {
	req := &addressPb.AddressStrId{
		Address: addressStr,
	}

	if resp, err := r.handler.GetAddressByStr(ctx, req); err != nil {

		return nil, nil, fmt.Errorf("handler get address by str: %w", err)
	} else if addressNew, addressOld, err := mapping.MapProtoAddressResponsesToModel(resp); err != nil {

		return nil, nil, fmt.Errorf("mapping address response: %w", err)
	} else {

		return addressNew, addressOld, nil
	}
}

func (r *addressRepositoryImp) GetAddressesByUserUuid(ctx context.Context, userUuid uuid.UUID) (model.Addresses, model.AddressesOld, error) {
	req := &addressPb.UserUuid{
		UserUuid: userUuid.String(),
	}

	if resp, err := r.handler.GetAddressesByUserUuid(ctx, req); err != nil {

		return nil, nil, fmt.Errorf("handler get address by user uuid: %w", err)
	} else if _, addressesNew, addressesOld, err := mapping.MapProtoAddressMultiResponsesToModel(resp); err != nil {

		return nil, nil, fmt.Errorf("mapping address response: %w", err)
	} else {

		return addressesNew, addressesOld, nil
	}
}

func (r *addressRepositoryImp) GetAddressesOldByFilter(ctx context.Context, filter *model.AddressOldFilter) (*uint64, model.AddressesOld, error) {
	filterProto := mapping.MapModelAddressOldFilterToProto(filter)

	if resp, err := r.handler.GetAddressesOldByFilter(ctx, filterProto); err != nil {

		return nil, nil, fmt.Errorf("handler get address by new filter: %w", err)
	} else if totalCount, _, addressesOld, err := mapping.MapProtoAddressMultiResponsesToModel(resp); err != nil {

		return nil, nil, fmt.Errorf("mapping address response: %w", err)
	} else {

		return totalCount, addressesOld, nil
	}
}

func (r *addressRepositoryImp) GetAddressesNewByFilter(ctx context.Context, filter *model.AddressFilter) (*uint64, model.Addresses, error) {
	filterProto := mapping.MapModelAddressNewFilterToProto(filter)

	if resp, err := r.handler.GetAddressesNewByFilter(ctx, filterProto); err != nil {

		return nil, nil, fmt.Errorf("handler get address by old filter: %w", err)
	} else if totalCount, addressesNew, _, err := mapping.MapProtoAddressMultiResponsesToModel(resp); err != nil {

		return nil, nil, fmt.Errorf("mapping new addresses response: %w", err)
	} else {

		return totalCount, addressesNew, nil
	}
}

func (r *addressRepositoryImp) AddOrUpdatePersonalAddress(ctx context.Context, addressStr string, userUuid uuid.UUID, network nodeCommon.NetworkEnum, minPayout *float64) (*model.AddressPersonal, error) {
	req := &addressPb.CreatePersonalAddressRequest{
		Address:   addressStr,
		UserUuid:  userUuid.String(),
		Network:   network.ToString(),
		MinPayout: minPayout,
	}

	if resp, err := r.handler.AddOrUpdatePersonalAddress(ctx, req); err != nil {

		return nil, fmt.Errorf("handler add or update personal address: %w", err)
	} else if addressesPersonal, err := mapping.MapProtoAddressResponseToModelPersonal(resp); err != nil {

		return nil, fmt.Errorf("mapping personal addresses response: %w", err)
	} else {

		return addressesPersonal, nil
	}
}

func (r *addressRepositoryImp) DeletePersonalAddress(ctx context.Context, userUuid uuid.UUID, network nodeCommon.NetworkEnum) error {
	req := &addressPb.DeletePersonalAddressRequest{
		UserUuid: userUuid.String(),
		Network:  network.ToString(),
	}

	if _, err := r.handler.DeletePersonalAddress(ctx, req); err != nil {

		return fmt.Errorf("handler delete personal address: %w", err)
	} else {

		return nil
	}
}

func (r *addressRepositoryImp) GetPersonalAddressesByFilter(ctx context.Context, filter *model.AddressPersonalFilter) (*uint64, model.AddressesPersonal, error) {
	filterProto := mapping.MapModelAddressPersonalFilterToProto(filter)

	if resp, err := r.handler.GetPersonalAddressesByFilter(ctx, filterProto); err != nil {

		return nil, nil, fmt.Errorf("handler get addresses by personal filter: %w", err)
	} else if totalCount, addressesNew, err := mapping.MapProtoPersonalAddressMultiResponsesToModel(resp); err != nil {

		return nil, nil, fmt.Errorf("mapping personal addresses response: %w", err)
	} else {

		return totalCount, addressesNew, nil
	}
}

func (r *addressRepositoryImp) GetPersonalAddressByUserUuid(ctx context.Context, userUuid uuid.UUID) (model.AddressesPersonal, error) {
	req := &addressPb.UserUuid{
		UserUuid: userUuid.String(),
	}

	if resp, err := r.handler.GetPersonalAddressesByUserUuid(ctx, req); err != nil {

		return nil, fmt.Errorf("handler get address by personal filter: %w", err)
	} else if resp == nil {
		// success
		return nil, nil
	} else if _, addressesPersonal, err := mapping.MapProtoPersonalAddressMultiResponsesToModel(resp); err != nil {

		return nil, fmt.Errorf("mapping personal address response: %w", err)
	} else {
		// success
		return addressesPersonal, nil
	}
}

func (r *addressRepositoryImp) CreateOrUpdateDirtyAddress(ctx context.Context, address *model.AddressDirty) (*model.AddressDirty, error) {
	req := mapping.MapModelAddressDirtyToProto(address)

	if _, err := r.handler.CreateOrUpdateDirtyAddress(ctx, req); err != nil {

		return nil, fmt.Errorf("handler create or update dirty address: %w", err)
	} else {

		return address, nil
	}
}

func (r *addressRepositoryImp) GetDirtyAddressesByFilter(ctx context.Context, filter *model.AddressDirtyFilter) (model.AddressesDirty, error) {
	req := mapping.MapModelAddressDirtyFilterToProto(filter)

	if addressesProto, err := r.handler.GetDirtyAddressesByFilter(ctx, req); err != nil {

		return nil, fmt.Errorf("handler get dirty address by filter: %w", err)
	} else if addressesResp, err := mapping.MapProtoAddressMultiFormToModelDirty(addressesProto); err != nil {

		return nil, fmt.Errorf("mapping dirty address response: %w", err)
	} else {

		return addressesResp, nil
	}
}
