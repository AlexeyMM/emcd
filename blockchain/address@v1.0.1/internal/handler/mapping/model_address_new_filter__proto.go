package mapping

import (
	"fmt"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapModelAddressNewFilterToProto(filter *model.AddressFilter) *addressPb.AddressNewFilter {
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

	var networkGroup *string
	if filter.NetworkGroup != nil {
		networkGroup = utils.StringToPtr(filter.NetworkGroup.ToString())

	} else {
		networkGroup = nil

	}

	var createdAtGt *timestamppb.Timestamp
	if filter.CreatedAtGt != nil {
		createdAtGt = timestamppb.New(*filter.CreatedAtGt)

	} else {
		createdAtGt = nil

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

	ret := &addressPb.AddressNewFilter{
		AddressUuid:  addressUuid,
		Address:      filter.Address,
		UserUuid:     userUuid,
		AddressType:  filter.AddressType,
		NetworkGroup: networkGroup,
		IsProcessing: filter.IsProcessing,
		CreatedAtGt:  createdAtGt,
		Pagination:   pagination,
	}

	return ret
}

func MapProtoToModelAddressNewFilter(p *addressPb.AddressNewFilter) (*model.AddressFilter, error) {
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

	var networkGroup *nodeCommon.NetworkGroupEnum
	if p.NetworkGroup != nil {
		networkGroupParsed := nodeCommon.NewNetworkGroupEnum(*p.NetworkGroup)
		if err := networkGroupParsed.Validate(); err != nil {

			return nil, fmt.Errorf("invalid network_group: %v, %w", *p.NetworkGroup, err)
		} else {
			networkGroup = &networkGroupParsed

		}
	} else {
		networkGroup = nil

	}

	var createdAtGt *time.Time
	if p.CreatedAtGt != nil {
		createdAtGt = utils.TimeToPtr(p.CreatedAtGt.AsTime())

	} else {
		createdAtGt = nil

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

	return &model.AddressFilter{
		Id:           addressUuid,
		Address:      p.Address,
		UserUuid:     userUuid,
		AddressType:  p.AddressType,
		NetworkGroup: networkGroup,
		IsProcessing: p.IsProcessing,
		CreatedAtGt:  createdAtGt,
		Pagination:   pagination,
	}, nil
}
