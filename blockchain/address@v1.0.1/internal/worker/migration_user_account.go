package worker

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/repository"
	"code.emcdtech.com/emcd/blockchain/address/internal/repository/repository_migration"
	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/enum"
	"code.emcdtech.com/emcd/blockchain/address/model/model_migration"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

const batchSize = 100
const token1 = 1
const token2 = 2
const token3 = 3
const token4 = 4
const token5 = 5
const token6 = 6

const migrateLogDelay = 10 * time.Minute

type MigrateAddress struct {
	coinValidator          coinValidatorRepo.CoinValidatorRepository
	repo                   repository_migration.MigrationRepository
	repoNew                repository.AddressRepository
	doUserAccountLogLastAt time.Time
	doAddressLogLastAt     time.Time
}

func NewMigrateAddress(coinValidator coinValidatorRepo.CoinValidatorRepository,
	repo repository_migration.MigrationRepository,
	repoNew repository.AddressRepository,
) *MigrateAddress {

	return &MigrateAddress{
		coinValidator:          coinValidator,
		repo:                   repo,
		repoNew:                repoNew,
		doUserAccountLogLastAt: time.Time{},
		doAddressLogLastAt:     time.Time{},
	}
}

func (m *MigrateAddress) DoUserAccount(ctx context.Context) error {
	if lastAt, err := m.repo.GetMigrationLastAt(ctx, model_migration.MigrationTableUsersAccounts); err != nil {
		sdkLog.Error(ctx, "failed get last at migration date of users accounts: %v", err)

		return err
	} else if totalCount, usersAccounts, err := m.repo.GetUserAccountMigrations(ctx, *lastAt, batchSize); err != nil {
		sdkLog.Error(ctx, "failed get users accounts: %v", err)

		return err
	} else {
		if *totalCount > 0 {
			sdkLog.Info(ctx, "migrate users accounts: %d/%d since %v", *totalCount, len(usersAccounts), *lastAt)

		} else if time.Since(m.doUserAccountLogLastAt) > migrateLogDelay {
			sdkLog.Info(ctx, "migrate users accounts: %d/%d since %v", *totalCount, len(usersAccounts), *lastAt)
			m.doUserAccountLogLastAt = time.Now()

		}

		createdAtLast := *lastAt
		for _, userAccount := range usersAccounts {
			if addressOld, err := convertUserAccountMigrationToOldAddress(m.coinValidator, userAccount); err != nil {
				sdkLog.Error(ctx, "failed convert users account to old address: %v", err)

			} else if ok := m.addOldAddressDirectThenCheck(ctx, addressOld); ok && createdAtLast.Before(addressOld.CreatedAt) {
				createdAtLast = addressOld.CreatedAt

			}
		}

		if len(usersAccounts) > 0 {
			sdkLog.Info(ctx, "migrated users accounts until %v", createdAtLast)

		}

		if len(usersAccounts) < batchSize {
			createdAtLast = createdAtLast.Add(time.Millisecond)

		}

		if err := m.repo.UpdateMigrationLastAt(ctx, model_migration.MigrationTableUsersAccounts, createdAtLast); err != nil {

			sdkLog.Error(ctx, "failed update last at migration data: %v", err)
		} else {
			// pass
		}

		return nil
	}
}

func (m *MigrateAddress) DoAddress(ctx context.Context) error {
	if lastAt, err := m.repo.GetMigrationLastAt(ctx, model_migration.MigrationTableAddresses); err != nil {
		sdkLog.Error(ctx, "failed get last at migration date of address: %v", err)

		return err
	} else if totalCount, addresses, err := m.repo.GetAddressMigrations(ctx, *lastAt, batchSize); err != nil {
		sdkLog.Error(ctx, "failed get addresses: %v", err)

		return err
	} else {
		if *totalCount > 0 {
			sdkLog.Info(ctx, "migrate addresses: %d/%d since %v", *totalCount, len(addresses), *lastAt)

		} else if time.Since(m.doAddressLogLastAt) > migrateLogDelay {
			sdkLog.Info(ctx, "migrate addresses: %d/%d since %v", *totalCount, len(addresses), *lastAt)
			m.doAddressLogLastAt = time.Now()

		}

		createdAtLast := *lastAt
		for _, address := range addresses {
			if !address.AddressOffset.Valid {
				if addressOld, err := convertAddressMigrationToOldAddress(m.coinValidator, address); err != nil {
					sdkLog.Error(ctx, "failed convert address to old address: %v", err)

				} else if ok := m.addOldAddressDirectThenCheck(ctx, addressOld); ok && createdAtLast.Before(addressOld.CreatedAt) {
					createdAtLast = addressOld.CreatedAt

				}
			} else {
				addressNewCommon, addressNewDerived := convertAddressMigrationToNewAddress(address)

				if ok := m.addNewAddressDirectThenCheck(ctx, addressNewCommon, addressNewDerived); ok && createdAtLast.Before(addressNewCommon.CreatedAt) {
					createdAtLast = addressNewCommon.CreatedAt

				}
			}
		}

		if len(addresses) > 0 {
			sdkLog.Info(ctx, "migrated addresses until %v", createdAtLast)

		}

		if len(addresses) < batchSize {
			createdAtLast = createdAtLast.Add(time.Millisecond)

		}

		if err := m.repo.UpdateMigrationLastAt(ctx, model_migration.MigrationTableAddresses, createdAtLast); err != nil {

			sdkLog.Error(ctx, "failed update last at migration data: %v", err)
		} else {
			// pass
		}

		return nil
	}
}

func (m *MigrateAddress) DoPersonalAddress(ctx context.Context) error {
	if lastAt, err := m.repo.GetMigrationLastAt(ctx, model_migration.MigrationTablePersonalAddresses); err != nil {
		sdkLog.Error(ctx, "failed get last at migration date of personal address: %v", err)

		return err
	} else if totalCount, addresses, err := m.repo.GetAddressPersonalMigrations(ctx, *lastAt, batchSize); err != nil {
		sdkLog.Error(ctx, "failed get personal addresses: %v", err)

		return err
	} else {
		if *totalCount > 0 {
			sdkLog.Info(ctx, "migrate addresses: %d/%d since %v", *totalCount, len(addresses), *lastAt)

		} else if time.Since(m.doAddressLogLastAt) > migrateLogDelay {
			sdkLog.Info(ctx, "migrate addresses: %d/%d since %v", *totalCount, len(addresses), *lastAt)
			m.doAddressLogLastAt = time.Now()

		}

		createdAtLast := *lastAt
		for _, address := range addresses {
			addressNewCommon := convertPersonalAddressMigrationToNewAddress(address)

			if ok := m.addPersonalAddressDirectThenCheck(ctx, addressNewCommon); ok && createdAtLast.Before(addressNewCommon.CreatedAt) {
				createdAtLast = addressNewCommon.CreatedAt

			}
		}

		if len(addresses) > 0 {
			sdkLog.Info(ctx, "migrated personal addresses until %v", createdAtLast)

		}

		if len(addresses) < batchSize {
			createdAtLast = createdAtLast.Add(time.Millisecond)

		}

		if err := m.repo.UpdateMigrationLastAt(ctx, model_migration.MigrationTablePersonalAddresses, createdAtLast); err != nil {

			sdkLog.Error(ctx, "failed update last at migration data: %v", err)
		} else {
			// pass
		}

		return nil
	}
}

func (m *MigrateAddress) addOldAddressDirectThenCheck(ctx context.Context, addressOld *model.AddressOld) bool {
	if err := m.repo.AddOldAddressDirect(ctx, addressOld); err != nil {
		if ok, errLocal := m.checkOldAddressConstraint(ctx, addressOld); errLocal != nil {
			sdkLog.Error(ctx, "failed add old address: %v, %+v", err, addressOld)
			sdkLog.Error(ctx, "failed check old address: %v", errLocal)

			return false
		} else if !ok {
			sdkLog.Error(ctx, "failed add old address: %+v", err, addressOld)

			return false
		} else {
			sdkLog.Warn(ctx, "old address already exists: %+v", addressOld)

			return true // update last at
		}
	} else {
		sdkLog.Info(ctx, "success old address: %+v", addressOld)

		return true // update last at
	}
}

func (m *MigrateAddress) addPersonalAddressDirectThenCheck(ctx context.Context, addressPersonal *model.AddressPersonal) bool {
	if err := m.repo.AddPersonalAddressDirect(ctx, addressPersonal); err != nil {
		if ok, errLocal := m.checkPersonalAddressConstraint(ctx, addressPersonal); errLocal != nil {
			sdkLog.Error(ctx, "failed add personal address: %v, %+v", err, addressPersonal)
			sdkLog.Error(ctx, "failed check personal address: %v", errLocal)

			return false
		} else if !ok {
			sdkLog.Error(ctx, "failed add personal address: %+v", err, addressPersonal)

			return false
		} else {
			sdkLog.Warn(ctx, "personal address already exists: %+v", addressPersonal)

			return true // update last at
		}
	} else {
		sdkLog.Info(ctx, "success personal address: %+v", addressPersonal)

		return true // update last at
	}
}

func (m *MigrateAddress) addNewAddressDirectThenCheck(ctx context.Context, addressNewCommon *model.Address, addressNewDerived *model.AddressDerived) bool {
	if err := m.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := m.repo.AddNewAddressDirect(ctx, addressNewCommon); err != nil {

			return err
		} else if err := m.repo.AddNewDerivedAddressDirect(ctx, addressNewDerived); err != nil {

			return err
		} else {

			return nil
		}
	}); err != nil {
		if ok, errLocal := m.checkNewAddressConstraint(ctx, addressNewCommon); errLocal != nil {
			sdkLog.Error(ctx, "failed add new address: %v, %+v", err, addressNewCommon)
			sdkLog.Error(ctx, "failed check new address: %v", errLocal)

			return false
		} else if !ok {
			sdkLog.Error(ctx, "failed add new address: %v, %+v", err, addressNewCommon)

			return false
		} else {
			sdkLog.Warn(ctx, "new address already exists: %+v", addressNewCommon)

			return true // update last at
		}
	} else {
		sdkLog.Info(ctx, "success new address: %+v", addressNewCommon)

		return true // update last at
	}
}

func (m *MigrateAddress) checkOldAddressConstraint(ctx context.Context, addressOld *model.AddressOld) (bool, error) {
	filter := &model.AddressOldFilter{
		Id:            nil,
		Address:       nil,
		UserUuid:      &addressOld.UserUuid,
		AddressType:   nil,
		Network:       addressOld.Network.ToPtr(),
		UserAccountId: nil,
		Coin:          &addressOld.Coin,
		CreatedAtGt:   nil,
		Pagination:    nil,
	}

	if _, addressesOldResp, err := m.repoNew.GetOldAddresses(ctx, filter); err != nil {

		return false, err
	} else if len(addressesOldResp) > 1 {

		return false, fmt.Errorf("too many old addresses by constraint: %d", len(addressesOldResp))
	} else if len(addressesOldResp) == 0 {

		return false, nil
	} else {

		return addressesOldResp[0].Address == addressOld.Address, nil
	}
}

func (m *MigrateAddress) checkPersonalAddressConstraint(ctx context.Context, addressPersonal *model.AddressPersonal) (bool, error) {
	filter := &model.AddressPersonalFilter{
		Id:         nil,
		Address:    nil,
		UserUuid:   &addressPersonal.UserUuid,
		Network:    addressPersonal.Network.ToPtr(),
		IsDeleted:  nil,
		Pagination: nil,
	}

	if _, addressesPersonalResp, err := m.repoNew.GetPersonalAddresses(ctx, filter); err != nil {

		return false, err
	} else if len(addressesPersonalResp) > 1 {

		return false, fmt.Errorf("too many personal addresses by constraint: %d", len(addressesPersonalResp))
	} else if len(addressesPersonalResp) == 0 {

		return false, nil
	} else {

		return addressesPersonalResp[0].Address == addressPersonal.Address, nil
	}
}

func (m *MigrateAddress) checkNewAddressConstraint(ctx context.Context, addressNew *model.Address) (bool, error) {
	filter := &model.AddressFilter{
		Id:           nil,
		Address:      nil,
		UserUuid:     &addressNew.UserUuid,
		IsProcessing: utils.BoolToPtr(false),
		AddressType:  addressNew.AddressType.Enum(),
		NetworkGroup: addressNew.NetworkGroup.ToPtr(),
		CreatedAtGt:  nil,
		Pagination:   nil,
	}

	if _, addressesNewResp, err := m.repoNew.GetNewAddresses(ctx, filter); err != nil {

		return false, err
	} else if len(addressesNewResp) > 1 {

		return false, fmt.Errorf("too many new addresses by constraint: %d", len(addressesNewResp))
	} else if len(addressesNewResp) == 0 {

		return false, nil
	} else {

		return addressesNewResp[0].Address == addressNew.Address, nil
	}
}

func convertUserAccountMigrationToOldAddress(coinValidator coinValidatorRepo.CoinValidatorRepository, ua *model_migration.UserAccountMigration) (*model.AddressOld, error) {
	var address string
	if !ua.Address.Valid {

		return nil, fmt.Errorf("address is invalid string")
	} else if ua.Address.String == "" {

		return nil, fmt.Errorf("address is empty space")
	} else {
		address = ua.Address.String

	}

	var network nodeCommon.NetworkEnum
	var coinCode string
	if coinCodeParsed, ok := coinValidator.GetCodeById(ua.CoinId); !ok {

		return nil, fmt.Errorf("coin code is invalid: %v", ua.CoinId)
	} else {
		coinCode = coinCodeParsed
		network = nodeCommon.NewNetworkEnum(coinCode)
		if err := network.Validate(); err != nil {

			return nil, fmt.Errorf("can not convert coin code to network enum %v: %w", coinCode, err)
		}
	}

	var addressType addressPb.AddressType
	switch network {
	case nodeCommon.TrxNetworkId:
		addressType = addressPb.AddressType_ADDRESS_TYPE_DIRECT

	case nodeCommon.TonNetworkId:
		addressType = addressPb.AddressType_ADDRESS_TYPE_MEMO

	case nodeCommon.EthNetworkId,
		nodeCommon.EtcNetworkId,
		nodeCommon.BnbNetworkId:
		addressType = addressPb.AddressType_ADDRESS_TYPE_BASED_ID

	case nodeCommon.BtcNetworkId,
		nodeCommon.BchNetworkId,
		nodeCommon.LtcNetworkId,
		nodeCommon.DashNetworkId,
		nodeCommon.DogeNetworkId,
		nodeCommon.KasNetworkId:
		addressType = addressPb.AddressType_ADDRESS_TYPE_BASED_NONE

	case nodeCommon.Erc20NetworkId,
		nodeCommon.Erc20NewNetworkId,
		nodeCommon.Bep20NetworkId,
		nodeCommon.Trc20NetworkId:

		return nil, fmt.Errorf("unused legacy network type: %v", network)
	case nodeCommon.AvaxNetworkId,
		nodeCommon.OpNetworkId,
		nodeCommon.PolygonNetworkId,
		nodeCommon.ArbNetworkId:

		return nil, fmt.Errorf("unsupported network type by user account way: %v", network)
	default:
		return nil, fmt.Errorf("unexpected network type by user account way: %v", network)
	}

	return &model.AddressOld{
		Id:            uuid.New(),
		Address:       address,
		UserUuid:      ua.UserUuid,
		AddressType:   enum.NewAddressTypeWrapper(addressType),
		Network:       nodeCommon.NewNetworkEnumWrapper(network),
		UserAccountId: ua.Id,
		Coin:          coinCode,
		CreatedAt:     ua.CreatedAt.Time,
	}, nil
}

func convertAddressMigrationToOldAddress(coinValidator coinValidatorRepo.CoinValidatorRepository, a *model_migration.AddressMigration) (*model.AddressOld, error) {
	var network nodeCommon.NetworkEnum
	var coinCode string
	if coinCodeParsed, ok := coinValidator.GetCodeById(a.CoinId); !ok {

		return nil, fmt.Errorf("coin code is invalid: %v", a.CoinId)
	} else {
		coinCode = coinCodeParsed
		if a.TokenId.Valid {
			switch a.TokenId.Int32 {
			case token1, token2:
				network = nodeCommon.BnbNetworkId
			case token3, token4:
				network = nodeCommon.TrxNetworkId
			case token5, token6:
				network = nodeCommon.EthNetworkId
			default:

				return nil, fmt.Errorf("unsupported token type: %v", a.TokenId.Int32)
			}
		} else {
			network = nodeCommon.NewNetworkEnum(coinCode)
			if err := network.Validate(); err != nil {

				return nil, fmt.Errorf("can not convert coin code to network enum %v: %w", coinCode, err)
			}
		}
	}

	var addressType addressPb.AddressType
	switch network {
	case nodeCommon.TrxNetworkId:
		addressType = addressPb.AddressType_ADDRESS_TYPE_DIRECT

	case nodeCommon.TonNetworkId:
		addressType = addressPb.AddressType_ADDRESS_TYPE_MEMO

	case nodeCommon.EthNetworkId,
		nodeCommon.EtcNetworkId,
		nodeCommon.BnbNetworkId:
		addressType = addressPb.AddressType_ADDRESS_TYPE_BASED_ID

	case nodeCommon.BtcNetworkId,
		nodeCommon.BchNetworkId,
		nodeCommon.LtcNetworkId,
		nodeCommon.DashNetworkId,
		nodeCommon.DogeNetworkId,
		nodeCommon.KasNetworkId:
		addressType = addressPb.AddressType_ADDRESS_TYPE_BASED_NONE

	case nodeCommon.Erc20NetworkId,
		nodeCommon.Erc20NewNetworkId,
		nodeCommon.Bep20NetworkId,
		nodeCommon.Trc20NetworkId:

		return nil, fmt.Errorf("unused legacy network type: %v", network)
	case nodeCommon.AvaxNetworkId,
		nodeCommon.OpNetworkId,
		nodeCommon.PolygonNetworkId,
		nodeCommon.ArbNetworkId:

		return nil, fmt.Errorf("unsupported network type by user account way: %v", network)
	default:
		return nil, fmt.Errorf("unexpected network type by user account way: %v", network)
	}

	return &model.AddressOld{
		Id:            uuid.New(),
		Address:       a.Address,
		UserUuid:      a.UserUuid,
		AddressType:   enum.NewAddressTypeWrapper(addressType),
		Network:       nodeCommon.NewNetworkEnumWrapper(network),
		UserAccountId: a.UserAccountId,
		Coin:          coinCode,
		CreatedAt:     a.CreatedAt,
	}, nil
}

func convertAddressMigrationToNewAddress(a *model_migration.AddressMigration) (*model.Address, *model.AddressDerived) {
	addressUuid := uuid.New()

	return &model.Address{
			Id:             addressUuid,
			Address:        a.Address,
			UserUuid:       a.UserUuid,
			ProcessingUuid: a.UserUuid,
			AddressType:    enum.NewAddressTypeWrapper(addressPb.AddressType_ADDRESS_TYPE_DERIVED),
			NetworkGroup:   nodeCommon.NewNetworkGroupEnumWrapper(nodeCommon.EthNetworkGroupId),
			CreatedAt:      a.CreatedAt,
		},
		&model.AddressDerived{
			AddressUuid:   addressUuid,
			NetworkGroup:  nodeCommon.NewNetworkGroupEnumWrapper(nodeCommon.EthNetworkGroupId),
			MasterKeyId:   0,
			DerivedOffset: uint32(a.AddressOffset.Int32),
		}
}

func convertPersonalAddressMigrationToNewAddress(a *model_migration.AddressPersonalMigration) *model.AddressPersonal {
	addressPersonalUuid := uuid.New()
	network := nodeCommon.NewNetworkEnum(a.CoinCode)

	return &model.AddressPersonal{
		Id:        addressPersonalUuid,
		Address:   a.AaAddress,
		UserUuid:  a.UserUuid,
		Network:   nodeCommon.NewNetworkEnumWrapper(network),
		MinPayout: a.MinPayout,
		DeletedAt: sql.NullTime{},
		UpdatedAt: time.Now().UTC(),
		CreatedAt: a.CreatedAt,
	}
}
