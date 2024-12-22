package mapping

import (
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapModelAddressDirtyFilterToProto(filter *model.AddressDirtyFilter) *addressPb.DirtyAddressFilter {
	var network *string
	if filter.Network != nil {
		network = utils.StringToPtr(filter.Network.ToString())

	} else {
		network = nil

	}

	return &addressPb.DirtyAddressFilter{
		Address: filter.Address,
		Network: network,
	}
}

func MapProtoToModelAddressDirtyFilter(p *addressPb.DirtyAddressFilter) (*model.AddressDirtyFilter, error) {
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

	return &model.AddressDirtyFilter{
		Address: p.Address,
		Network: network,
	}, nil
}
