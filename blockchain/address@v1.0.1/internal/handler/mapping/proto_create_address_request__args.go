package mapping

import (
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"github.com/google/uuid"

	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapProtoAddressRequestToArgs(coinValidator coinValidatorRepo.CoinValidatorRepository, p *addressPb.CreateAddressRequest) (*uuid.UUID, *nodeCommon.NetworkEnum, *string, error) {
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

	var coin *string
	if p.Coin != nil {
		if !coinValidator.IsValidCode(*p.Coin) {

			return nil, nil, nil, fmt.Errorf("failed validate coin: %s", *p.Coin)
		}

		coin = p.Coin
	} else {
		coin = nil

	}

	return userUuid, network, coin, nil
}
