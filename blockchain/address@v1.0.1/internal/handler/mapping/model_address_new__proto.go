package mapping

import (
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapModelAddressNewToProto(address *model.Address) *addressPb.AddressResponse {
	newWay := &addressPb.NewWay{
		MasterKeyId:    nil,
		DerivedOffset:  nil,
		Networks:       address.NetworkGroup.GetNetworks().ToStrings(),
		ProcessingUuid: nil,
	}

	if addressDerived := address.GetAddressDerived(); addressDerived != nil {
		newWay.MasterKeyId = &addressDerived.MasterKeyId
		newWay.DerivedOffset = &addressDerived.DerivedOffset

	}

	if address.ProcessingUuid != address.UserUuid {
		newWay.ProcessingUuid = utils.StringToPtr(address.ProcessingUuid.String())

	}

	ret := &addressPb.AddressResponse{
		AddressUuid:  address.Id.String(),
		Address:      address.Address,
		UserUuid:     address.UserUuid.String(),
		AddressType:  address.AddressType.AddressType,
		NetworkGroup: address.NetworkGroup.ToString(),
		Way:          &addressPb.AddressResponse_NewWay{NewWay: newWay},
		CreatedAt:    timestamppb.New(address.CreatedAt),
	}

	return ret
}

func MapProtoAddressResponseToModelNew(p *addressPb.AddressResponse) (*model.Address, error) {
	newWay := p.GetNewWay()

	if newWay == nil {

		return nil, fmt.Errorf("address response empty new way")
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

	var processingUuid uuid.UUID
	if newWay.ProcessingUuid != nil {

		if processingUuidParsed, err := uuid.Parse(*newWay.ProcessingUuid); err != nil {

			return nil, fmt.Errorf("failed parse processing_uuid: %s, %w", p.UserUuid, err)
		} else {
			processingUuid = processingUuidParsed

		}
	} else {
		processingUuid = userUuid

	}

	var networkGroup nodeCommon.NetworkGroupEnum
	networkGroupNew := nodeCommon.NewNetworkGroupEnum(p.NetworkGroup)
	if err := networkGroupNew.Validate(); err != nil {

		return nil, fmt.Errorf("failed parse network_group: %s, %w", p.NetworkGroup, err)
	} else {
		networkGroup = networkGroupNew

	}

	address := &model.Address{
		Id:             addressUuid,
		Address:        p.Address,
		UserUuid:       userUuid,
		ProcessingUuid: processingUuid,
		AddressType:    enum.NewAddressTypeWrapper(p.AddressType),
		NetworkGroup:   nodeCommon.NewNetworkGroupEnumWrapper(networkGroup),
		CreatedAt:      p.CreatedAt.AsTime(),
	}

	if newWay.DerivedOffset != nil && newWay.MasterKeyId != nil {
		derivedAddress := &model.AddressDerived{
			AddressUuid:   address.Id,
			MasterKeyId:   *newWay.MasterKeyId,
			DerivedOffset: *newWay.DerivedOffset,
			NetworkGroup:  address.NetworkGroup,
		}

		address.SetAddressDerived(derivedAddress)
	}

	return address, nil
}
