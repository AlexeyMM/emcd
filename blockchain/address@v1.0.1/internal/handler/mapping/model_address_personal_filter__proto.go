package mapping

import (
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapModelAddressPersonalFilterToProto(filter *model.AddressPersonalFilter) *addressPb.AddressPersonalFilter {
	var addressUuid *string
	if filter.Id != nil {
		addressUuid = utils.StringToPtr(filter.Id.String())

	} else {
		addressUuid = nil

	}

	var userUuid *string
	if filter.UserUuid != nil {
		userUuid = utils.StringToPtr(filter.UserUuid.String())

	} else {
		userUuid = nil

	}

	var network *string
	if filter.Network != nil {
		network = utils.StringToPtr(filter.Network.ToString())

	} else {
		network = nil

	}

	var pagination *addressPb.AddressPagination
	if filter.Pagination != nil {
		pagination = &addressPb.AddressPagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}

	} else {
		pagination = nil

	}

	ret := &addressPb.AddressPersonalFilter{
		AddressUuid: addressUuid,
		Address:     filter.Address,
		UserUuid:    userUuid,
		Network:     network,
		IsDeleted:   filter.IsDeleted,
		Pagination:  pagination,
	}

	return ret
}

func MapProtoToModelAddressPersonalFilter(p *addressPb.AddressPersonalFilter) (*model.AddressPersonalFilter, error) {
	var addressUuid *uuid.UUID
	if p.AddressUuid != nil {
		if addressUuidParsed, err := uuid.Parse(*p.AddressUuid); err != nil {

			return nil, fmt.Errorf("invalid address_uuid: %v, %w", *p.AddressUuid, err)
		} else {
			addressUuid = &addressUuidParsed

		}
	} else {
		addressUuid = nil

	}

	var userUuid *uuid.UUID
	if p.UserUuid != nil {
		if userUuidParsed, err := uuid.Parse(*p.UserUuid); err != nil {

			return nil, fmt.Errorf("invalid user_uuid: %v, %w", *p.UserUuid, err)
		} else {
			userUuid = &userUuidParsed

		}
	} else {
		userUuid = nil

	}

	var network *nodeCommon.NetworkEnum
	if p.Network != nil {
		networkParsed := nodeCommon.NewNetworkEnum(*p.Network)
		if err := networkParsed.Validate(); err != nil {

			return nil, fmt.Errorf("invalid network_group: %v, %w", *p.Network, err)
		} else {
			network = &networkParsed

		}
	} else {
		network = nil

	}

	var pagination *model.Pagination
	if p.Pagination != nil {
		pagination = &model.Pagination{
			Limit:  p.Pagination.Limit,
			Offset: p.Pagination.Offset,
		}

	} else {
		pagination = nil

	}

	return &model.AddressPersonalFilter{
		Id:         addressUuid,
		Address:    p.Address,
		UserUuid:   userUuid,
		Network:    network,
		IsDeleted:  p.IsDeleted,
		Pagination: pagination,
	}, nil
}
