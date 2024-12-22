package mapping

import (
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"

	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapProtoProcessingAddressRequestToArgs(p *addressPb.CreateProcessingAddressRequest) (*uuid.UUID, *nodeCommon.NetworkEnum, *uuid.UUID, error) {
	var userUuid *uuid.UUID
	if userUuidParsed, err := uuid.Parse(p.UserUuid); err != nil {

		return nil, nil, nil, fmt.Errorf("failed parse user_uuid: %s, %w", p.UserUuid, err)
	} else {
		userUuid = &userUuidParsed

	}

	var network *nodeCommon.NetworkEnum
	networkNew := nodeCommon.NewNetworkEnum(p.Network)
	if err := networkNew.Validate(); err != nil {

		return nil, nil, nil, fmt.Errorf("failed parse network_group: %s, %w", p.Network, err)
	} else {
		network = &networkNew

	}

	var processingUuid *uuid.UUID
	if processingUuidParsed, err := uuid.Parse(p.ProcessingUuid); err != nil {

		return nil, nil, nil, fmt.Errorf("failed parse processing_uuid: %s, %w", p.ProcessingUuid, err)
	} else {
		processingUuid = &processingUuidParsed

	}

	return userUuid, network, processingUuid, nil
}
