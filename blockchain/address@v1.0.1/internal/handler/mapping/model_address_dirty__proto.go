package mapping

import (
	"fmt"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"

	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapModelAddressDirtyToProto(address *model.AddressDirty) *addressPb.DirtyAddressForm {

	return &addressPb.DirtyAddressForm{
		Address: address.Address,
		Network: address.Network.ToString(),
		IsDirty: address.IsDirty,
	}
}

func MapProtoAddressFormToModelDirty(p *addressPb.DirtyAddressForm) (*model.AddressDirty, error) {
	var network nodeCommon.NetworkEnum
	networkNew := nodeCommon.NewNetworkEnum(p.Network)
	if err := networkNew.Validate(); err != nil {

		return nil, fmt.Errorf("failed parse network: %s, %w", p.Network, err)
	} else {
		network = networkNew

	}

	currentTime := time.Now().UTC()

	return &model.AddressDirty{
		Address:   p.Address,
		Network:   nodeCommon.NewNetworkEnumWrapper(network),
		IsDirty:   p.IsDirty,
		UpdatedAt: currentTime,
		CreatedAt: currentTime,
	}, nil
}
