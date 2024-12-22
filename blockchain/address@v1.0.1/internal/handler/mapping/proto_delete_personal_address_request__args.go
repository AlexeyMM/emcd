package mapping

import (
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"

	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapProtoDeletePersonalAddressRequestToArgs(p *addressPb.DeletePersonalAddressRequest) (*uuid.UUID, *nodeCommon.NetworkEnum, error) {
	var userUuid *uuid.UUID
	if userUuidParsed, err := uuid.Parse(p.UserUuid); err != nil {

		return nil, nil, fmt.Errorf("failed parse user_uuid: %s, %w", p.UserUuid, err)
	} else {
		userUuid = &userUuidParsed

	}

	var network *nodeCommon.NetworkEnum
	networkNew := nodeCommon.NewNetworkEnum(p.Network)
	if err := networkNew.Validate(); err != nil {

		return nil, nil, fmt.Errorf("failed parse network_group: %s, %w", p.Network, err)
	} else {
		network = &networkNew

	}

	return userUuid, network, nil
}
