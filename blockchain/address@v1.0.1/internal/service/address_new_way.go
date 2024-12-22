package service

import (
	"context"
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func (s *addressServiceImp) createMemoNewAddress(ctx context.Context, addressUuid, userUuid uuid.UUID, processingUuid *uuid.UUID, networkGroup nodeCommon.NetworkGroupEnum) (*model.Address, error) {
	addressMemo := uuid.New().String()
	var address *model.Address

	if processingUuid != nil {
		address = model.NewProcessingAddress(addressUuid, addressMemo, userUuid, *processingUuid, addressPb.AddressType_ADDRESS_TYPE_MEMO, networkGroup)

	} else {
		address = model.NewAddress(addressUuid, addressMemo, userUuid, addressPb.AddressType_ADDRESS_TYPE_MEMO, networkGroup)

	}

	if err := s.addressRepo.AddNewCommonAddress(ctx, address); err != nil {

		return nil, fmt.Errorf("add memo address: %w", err)
	} else {

		return address, nil
	}
}

func (s *addressServiceImp) createNodeNewAddress(
	ctx context.Context,
	addressUuid uuid.UUID,
	userUuid uuid.UUID,
	processingUuid *uuid.UUID,
	addressType addressPb.AddressType,
	networkGroup nodeCommon.NetworkGroupEnum,
) (*model.Address, error) {
	if err := s.validatePreparedNodeAddress(addressType, networkGroup, false); err != nil {

		return nil, fmt.Errorf("debug validate: %w", err)
	} else if len(networkGroup.GetNetworks()) > 1 {

		return nil, fmt.Errorf("debug validate: network group contain more then one networks: %v, %+v", networkGroup, networkGroup.GetNetworks())
	} else if addressStr, err := s.nodeRepo.GenerateAddress(ctx, networkGroup.GetNetworks()[0], nil, &userUuid); err != nil {

		return nil, fmt.Errorf("generate address: %w", err)
	} else {
		var addressNew *model.Address
		if processingUuid != nil {
			addressNew = model.NewProcessingAddress(addressUuid, addressStr, userUuid, *processingUuid, addressType, networkGroup)

		} else {
			addressNew = model.NewAddress(addressUuid, addressStr, userUuid, addressType, networkGroup)

		}

		if err := s.addressRepo.AddNewCommonAddress(ctx, addressNew); err != nil {

			return nil, fmt.Errorf("add common address: %w", err)
		} else {

			return addressNew, nil
		}
	}
}
