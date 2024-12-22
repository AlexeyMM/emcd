package mapping

import (
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"

	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapProtoPersonalAddressRequestToArgs(p *addressPb.CreatePersonalAddressRequest) (string, *uuid.UUID, *nodeCommon.NetworkEnum, *float64, error) {
	var addressStr string
	if p.Address == "" {

		return "", nil, nil, nil, fmt.Errorf("address is empty")
	} else {
		addressStr = p.Address

	}

	var userUuid *uuid.UUID
	if userUuidParsed, err := uuid.Parse(p.UserUuid); err != nil {

		return "", nil, nil, nil, fmt.Errorf("failed parse user_uuid: %s, %w", p.UserUuid, err)
	} else {
		userUuid = &userUuidParsed

	}

	var network *nodeCommon.NetworkEnum
	networkNew := nodeCommon.NewNetworkEnum(p.Network)
	if err := networkNew.Validate(); err != nil {

		return "", nil, nil, nil, fmt.Errorf("failed parse network_group: %s, %w", p.Network, err)
	} else {
		network = &networkNew

	}

	return addressStr, userUuid, network, p.MinPayout, nil
}
