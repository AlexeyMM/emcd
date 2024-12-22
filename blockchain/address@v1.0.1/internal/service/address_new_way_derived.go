package service

import (
	"context"
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/blockchain/address/internal/repository"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func (s *addressServiceImp) createDerivedAddress(ctx context.Context, addressUuid, userUuid uuid.UUID, processingUuid *uuid.UUID, networkGroup nodeCommon.NetworkGroupEnum) (*model.Address, error) {
	var address *model.Address
	if processingUuid != nil {
		address = model.NewProcessingAddress(addressUuid, "", userUuid, *processingUuid, addressPb.AddressType_ADDRESS_TYPE_DERIVED, networkGroup)

	} else {
		address = model.NewAddress(addressUuid, "", userUuid, addressPb.AddressType_ADDRESS_TYPE_DERIVED, networkGroup)

	}

	if derivedFunc, masterKeyId, err := s.GetDerivedFunc(networkGroup); err != nil {

		return nil, fmt.Errorf("get derived function: %w", err)
	} else if err := s.addressRepo.AddNewDerivedAddress(ctx, address, *masterKeyId, derivedFunc); err != nil {

		return nil, fmt.Errorf("add derived address: %w", err)
	} else {

		return address, nil
	}
}

func (s *addressServiceImp) GetDerivedFunc(networkGroup nodeCommon.NetworkGroupEnum) (repository.DerivedFunc, *uint32, error) {
	if masterKeys, ok := s.masterKeysIdMap[networkGroup]; !ok {

		return nil, nil, fmt.Errorf("master keys not found for network group: %v", networkGroup)
	} else {
		masterKeyId := s.choseMasterKeyId(masterKeys)
		masterKeyPub := s.masterKeysIdMap[networkGroup][masterKeyId]

		rootKey, err := hdkeychain.NewKeyFromString(masterKeyPub)
		if err != nil {

			return nil, nil, fmt.Errorf("get root key: %w", err)
		}

		return s.getDeriveFuncFromRoot(rootKey), &masterKeyId, nil
	}
}

func (s *addressServiceImp) choseMasterKeyId(masterKeys []string) uint32 {
	if s.useLastMasterKey {

		return uint32(len(masterKeys) - 1)
	} else {

		return 0
	}
}

func (*addressServiceImp) getDeriveFuncFromRoot(rootKey *hdkeychain.ExtendedKey) repository.DerivedFunc {

	return func(derivedOffset uint32) (string, error) {
		if derivedKey, err := rootKey.Derive(derivedOffset); err != nil {

			return "", fmt.Errorf("derive key from root: %w", err)
		} else if addrPubKey, err := derivedKey.ECPubKey(); err != nil {

			return "", fmt.Errorf("derived key convert: %w", err)
		} else {
			addrPubECDSAKey := addrPubKey.ToECDSA()
			addressStr := crypto.PubkeyToAddress(*addrPubECDSAKey).Hex()

			return addressStr, nil
		}
	}
}
