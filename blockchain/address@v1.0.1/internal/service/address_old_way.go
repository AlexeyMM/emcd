package service

import (
	"context"
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func (s *addressServiceImp) createMemoOldAddress(ctx context.Context, addressUuid, userUuid uuid.UUID, network nodeCommon.NetworkEnum, userAccountId int32, coin string) (*model.AddressOld, error) {
	addressMemo := uuid.New().String()
	addressOld := model.NewAddressOld(addressUuid, addressMemo, userUuid, addressPb.AddressType_ADDRESS_TYPE_MEMO, network, userAccountId, coin)

	if err := s.addressRepo.AddOldAddress(ctx, addressOld); err != nil {

		return nil, fmt.Errorf("add memo address: %w", err)
	} else {

		return addressOld, nil
	}
}

func (s *addressServiceImp) createNodeOldAddress(
	ctx context.Context,
	addressUuid, userUuid uuid.UUID,
	addressType addressPb.AddressType,
	network nodeCommon.NetworkEnum,
	userAccountId int32,
	coin string,
) (*model.AddressOld, error) {
	if err := s.validatePreparedNodeAddress(addressType, network.Group(), true); err != nil {

		return nil, fmt.Errorf("debug validate: %w", err)
	} else if addressStr, err := s.nodeRepo.GenerateAddress(ctx, network, &userAccountId, &userUuid); err != nil {

		return nil, fmt.Errorf("generate address: %w", err)
	} else {
		addressOld := model.NewAddressOld(addressUuid, addressStr, userUuid, addressType, network, userAccountId, coin)

		if err := s.addressRepo.AddOldAddress(ctx, addressOld); err != nil {

			return nil, fmt.Errorf("add common address: %w", err)
		} else {

			return addressOld, nil
		}
	}
}
