package mapping

import (
	"fmt"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapModelAddressOldFilterToProto(filter *model.AddressOldFilter) *addressPb.AddressOldFilter {
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

	ret := &addressPb.AddressOldFilter{
		AddressUuid:   addressUuid,
		Address:       filter.Address,
		UserUuid:      userUuid,
		AddressType:   filter.AddressType,
		Network:       network,
		UserAccountId: filter.UserAccountId,
		Coin:          filter.Coin,
		CreatedAtGt:   createdAtGt,
		Pagination:    pagination,
	}

	return ret
}

func MapProtoToModelAddressOldFilter(coinValidator coinValidatorRepo.CoinValidatorRepository, p *addressPb.AddressOldFilter) (*model.AddressOldFilter, error) {
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

			return nil, fmt.Errorf("invalid network: %v, %w", *p.Network, err)
		} else {
			network = &networkParsed

		}
	} else {
		network = nil

	}

	var coin *string
	if p.Coin != nil {
		if ok := coinValidator.IsValidCode(*p.Coin); !ok {

			return nil, fmt.Errorf("invalid coin: %v", *p.Coin)
		} else {
			coin = p.Coin

		}
	} else {
		coin = nil

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

	return &model.AddressOldFilter{
		Id:            addressUuid,
		Address:       p.Address,
		UserUuid:      userUuid,
		AddressType:   p.AddressType,
		Network:       network,
		UserAccountId: p.UserAccountId,
		Coin:          coin,
		CreatedAtGt:   createdAtGt,
		Pagination:    pagination,
	}, nil
}
