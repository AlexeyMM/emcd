package mapping

import (
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapModelAddressOldToProto(address *model.AddressOld) *addressPb.AddressResponse {
	oldWay := &addressPb.OldWay{
		UserAccountId: address.UserAccountId,
		Network:       address.Network.ToString(),
		Coin:          address.Coin,
	}

	ret := &addressPb.AddressResponse{
		AddressUuid:  address.Id.String(),
		Address:      address.Address,
		UserUuid:     address.UserUuid.String(),
		AddressType:  address.AddressType.AddressType,
		NetworkGroup: address.Network.Group().ToString(),
		Way:          &addressPb.AddressResponse_OldWay{OldWay: oldWay},
		CreatedAt:    timestamppb.New(address.CreatedAt),
	}

	return ret
}

func MapProtoAddressResponseToModelOld(p *addressPb.AddressResponse) (*model.AddressOld, error) {
	oldWay := p.GetOldWay()

	if oldWay == nil {

		return nil, fmt.Errorf("address response empty old way")
	}

	var addressUuid uuid.UUID
	if addressUuidParsed, err := uuid.Parse(p.AddressUuid); err != nil {

		return nil, fmt.Errorf("failed parse address id: %s, %w", p.UserUuid, err)
	} else {
		addressUuid = addressUuidParsed

	}

	var userUuid uuid.UUID
	if userUuidParsed, err := uuid.Parse(p.UserUuid); err != nil {

		return nil, fmt.Errorf("failed parse user_uuid: %s, %w", p.UserUuid, err)
	} else {
		userUuid = userUuidParsed

	}

	var network nodeCommon.NetworkEnum
	networkNew := nodeCommon.NewNetworkEnum(oldWay.Network)
	if err := networkNew.Validate(); err != nil {

		return nil, fmt.Errorf("failed parse network: %s, %w", oldWay.Network, err)
	} else {
		network = networkNew

	}

	address := &model.AddressOld{
		Id:            addressUuid,
		Address:       p.Address,
		UserUuid:      userUuid,
		AddressType:   enum.NewAddressTypeWrapper(p.AddressType),
		Network:       nodeCommon.NewNetworkEnumWrapper(network),
		UserAccountId: oldWay.UserAccountId,
		Coin:          oldWay.Coin,
		CreatedAt:     p.CreatedAt.AsTime(),
	}

	return address, nil
}
