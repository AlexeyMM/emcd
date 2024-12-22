package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	nodeRepository "code.emcdtech.com/emcd/blockchain/node/repository_external"
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	userAccountModel "code.emcdtech.com/emcd/service/accounting/model"
	userAccountModelEnum "code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountRepository "code.emcdtech.com/emcd/service/accounting/repository"
	profilePb "code.emcdtech.com/emcd/service/profile/protocol/profile"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/repository"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

type AddressService interface {
	CreateNewAddress(ctx context.Context, addressUuid, userUuid uuid.UUID, addressType addressPb.AddressType, networkGroup nodeCommon.NetworkGroupEnum) (*model.Address, error)
	CreateOldAddress(ctx context.Context, addressUuid, userUuid uuid.UUID, addressType addressPb.AddressType, network nodeCommon.NetworkEnum, coin string) (*model.AddressOld, error)
	CreateProcessingAddress(ctx context.Context, addressUuid, userUuid, processingUuid uuid.UUID, addressType addressPb.AddressType, networkGroup nodeCommon.NetworkGroupEnum) (*model.Address, error)
	CreatePersonalAddress(ctx context.Context, addressStr string, userUuid uuid.UUID, network nodeCommon.NetworkEnum, minPayout float64) (*model.AddressPersonal, error)
	UpdatePersonalAddress(ctx context.Context, address *model.AddressPersonal, addressStr string, minPayout *float64) (*model.AddressPersonal, error)
	DeletePersonalAddress(ctx context.Context, userUuid uuid.UUID, network nodeCommon.NetworkEnum) error
	CreateOrUpdateDirtyAddress(ctx context.Context, address *model.AddressDirty) (*model.AddressDirty, error)
	// AddOrUpdatePersonalAddress(ctx context.Context, addressStr string, userUuid uuid.UUID, network nodeCommon.NetworkEnum, minPayout *float64) (*model.AddressPersonal, error)

	GetNewAddressByConstraint(ctx context.Context, userUuid uuid.UUID, addressType addressPb.AddressType, networkGroup nodeCommon.NetworkGroupEnum) (model.Addresses, error)
	GetOldAddressByConstraint(ctx context.Context, userUuid uuid.UUID, network nodeCommon.NetworkEnum, coin string) (model.AddressesOld, error)
	GetPersonalAddressByConstraint(ctx context.Context, userUuid uuid.UUID, network nodeCommon.NetworkEnum) (model.AddressesPersonal, error)

	GetNewAddressByUuid(ctx context.Context, addressUuid uuid.UUID) (model.Addresses, error)
	GetNewAddressByStr(ctx context.Context, addressStr string) (model.Addresses, error)
	GetNewAddressesByUserUuid(ctx context.Context, userUuid uuid.UUID) (model.Addresses, error)

	GetOldAddressByUuid(ctx context.Context, addressUuid uuid.UUID) (model.AddressesOld, error)
	GetOldAddressByStr(ctx context.Context, addressStr string) (model.AddressesOld, error)
	GetOldAddressesByUserUuid(ctx context.Context, userUuid uuid.UUID) (model.AddressesOld, error)

	GetNewAddressesByFilter(ctx context.Context, filter *model.AddressFilter) (*uint64, model.Addresses, error)
	GetOldAddressesByFilter(ctx context.Context, filter *model.AddressOldFilter) (*uint64, model.AddressesOld, error)
	GetPersonalAddressesByFilter(ctx context.Context, filter *model.AddressPersonalFilter) (*uint64, model.AddressesPersonal, error)
	GetDirtyAddressesByFilter(ctx context.Context, filter *model.AddressDirtyFilter) (model.AddressesDirty, error)

	GetDerivedFunc(networkGroup nodeCommon.NetworkGroupEnum) (repository.DerivedFunc, *uint32, error)
}

type addressServiceImp struct {
	addressRepo      repository.AddressRepository
	userAccountRepo  userAccountRepository.UserAccountRepository
	nodeRepo         nodeRepository.AddressNodeRepository
	profileProtoCli  profilePb.ProfileServiceClient
	masterKeysIdMap  map[nodeCommon.NetworkGroupEnum][]string
	useLastMasterKey bool
	rabbitService    RabbitService
}

func NewAddressService(
	addressRepo repository.AddressRepository,
	userAccountRepo userAccountRepository.UserAccountRepository,
	nodeRepo nodeRepository.AddressNodeRepository,
	profileProtoCli profilePb.ProfileServiceClient,
	masterKeyIdMap map[nodeCommon.NetworkGroupEnum][]string,
	useLastMasterKey bool, // otherwise is first
	rabbitService RabbitService,
) AddressService {

	return &addressServiceImp{
		addressRepo:      addressRepo,
		userAccountRepo:  userAccountRepo,
		nodeRepo:         nodeRepo,
		profileProtoCli:  profileProtoCli,
		masterKeysIdMap:  masterKeyIdMap,
		useLastMasterKey: useLastMasterKey,
		rabbitService:    rabbitService,
	}
}

func (s *addressServiceImp) CreateNewAddress(
	ctx context.Context,
	addressUuid uuid.UUID,
	userUuid uuid.UUID,
	addressType addressPb.AddressType,
	networkGroup nodeCommon.NetworkGroupEnum,
) (*model.Address, error) {
	if address, err := s._createNewAddress(ctx, addressUuid, userUuid, nil, addressType, networkGroup); err != nil {

		return nil, fmt.Errorf("create new address: %w", err)

	} else {
		s.sendConsumerNetworkGroupReloadAddr(ctx, networkGroup)

		return address, err
	}
}

func (s *addressServiceImp) CreateProcessingAddress(
	ctx context.Context,
	addressUuid uuid.UUID,
	userUuid uuid.UUID,
	processingUuid uuid.UUID,
	addressType addressPb.AddressType,
	networkGroup nodeCommon.NetworkGroupEnum,
) (*model.Address, error) {
	if address, err := s._createNewAddress(ctx, addressUuid, userUuid, &processingUuid, addressType, networkGroup); err != nil {

		return nil, fmt.Errorf("create processing address: %w", err)

	} else {
		s.sendConsumerNetworkGroupReloadAddr(ctx, networkGroup)

		return address, err
	}
}

func (s *addressServiceImp) _createNewAddress(
	ctx context.Context,
	addressUuid uuid.UUID,
	userUuid uuid.UUID,
	processingUuid *uuid.UUID,
	addressType addressPb.AddressType,
	networkGroup nodeCommon.NetworkGroupEnum,
) (*model.Address, error) {
	switch addressType {
	case addressPb.AddressType_ADDRESS_TYPE_DERIVED:
		return s.createDerivedAddress(ctx, addressUuid, userUuid, processingUuid, networkGroup)
	case addressPb.AddressType_ADDRESS_TYPE_DIRECT, addressPb.AddressType_ADDRESS_TYPE_BASED_NONE:
		return s.createNodeNewAddress(ctx, addressUuid, userUuid, processingUuid, addressType, networkGroup)
	case addressPb.AddressType_ADDRESS_TYPE_MEMO:
		return s.createMemoNewAddress(ctx, addressUuid, userUuid, processingUuid, networkGroup)
	default:
		return nil, fmt.Errorf("unsupported address type %v", addressType)
	}
}

func (s *addressServiceImp) CreateOldAddress(ctx context.Context, addressUuid, userUuid uuid.UUID, addressType addressPb.AddressType, network nodeCommon.NetworkEnum, coin string) (*model.AddressOld, error) {
	if userAccountId, err := s.getUserAccountId(ctx, userUuid, coin); err != nil {

		return nil, fmt.Errorf("get user account id: %w", err)
	} else if addressOld, err := s._createOldAddress(ctx, addressUuid, userUuid, addressType, network, *userAccountId, coin); err != nil {

		return nil, fmt.Errorf("create old address: %w", err)

	} else {
		s.sendConsumerNetworkReloadAddr(ctx, network)

		return addressOld, err
	}
}

func (s *addressServiceImp) _createOldAddress(
	ctx context.Context,
	addressUuid uuid.UUID,
	userUuid uuid.UUID,
	addressType addressPb.AddressType,
	network nodeCommon.NetworkEnum,
	userAccountId int32,
	coin string,
) (*model.AddressOld, error) {
	switch addressType {
	case addressPb.AddressType_ADDRESS_TYPE_BASED_ID,
		addressPb.AddressType_ADDRESS_TYPE_BASED_NONE,
		addressPb.AddressType_ADDRESS_TYPE_DIRECT:
		return s.createNodeOldAddress(ctx, addressUuid, userUuid, addressType, network, userAccountId, coin)
	case addressPb.AddressType_ADDRESS_TYPE_MEMO:
		return s.createMemoOldAddress(ctx, addressUuid, userUuid, network, userAccountId, coin)
	default:
		return nil, fmt.Errorf("unsupported address type %v", addressType)
	}
}

func (s *addressServiceImp) getUserAccountId(ctx context.Context, userUuid uuid.UUID, coin string) (*int32, error) {
	if userId, err := s.getUserIdByUuid(ctx, userUuid); err != nil {

		return nil, fmt.Errorf("get user by uuid: %w", err)
	} else if userAccountId, err := s.getOrCreateUserAccount(ctx, *userId, userUuid, coin); err != nil {

		return nil, fmt.Errorf("get or create user account: %w", err)
	} else {

		return userAccountId, nil
	}
}

func (s *addressServiceImp) getUserIdByUuid(ctx context.Context, userUuid uuid.UUID) (*int32, error) {
	req := &profilePb.GetByUserIDV2Request{
		UserID: userUuid.String(),
	}

	if userResp, err := s.profileProtoCli.GetByUserIDV2(ctx, req); err != nil {

		return nil, fmt.Errorf("get user via profile: %w", err)
	} else {

		return &userResp.Profile.User.OldID, nil
	}
}

func (s *addressServiceImp) getOrCreateUserAccount(ctx context.Context, userId int32, userUuid uuid.UUID, coin string) (*int32, error) {
	userAccount := &userAccountModel.UserAccount{
		ID:            0,
		UserID:        userId,
		CoinID:        0, // will ignore
		AccountTypeID: userAccountModelEnum.NewAccountTypeIdWrapper(userAccountModelEnum.WalletAccountTypeID),
		Minpay:        0.0,
		Address:       sql.NullString{},
		ChangedAt:     sql.NullTime{},
		Img1:          sql.NullFloat64{},
		Img2:          sql.NullFloat64{},
		IsActive:      utils.BoolToBoolNull(true),
		CreatedAt:     sql.NullTime{},
		UpdatedAt:     sql.NullTime{},
		Fee:           utils.Float64ToFloat64Null(0.0),
		UserIDNew:     utils.UuidToUuidNull(userUuid),
		CoinNew:       utils.StringToStringNull(coin),
	}

	if userAccounts, err := s.userAccountRepo.CreateUserAccounts(ctx, userId, userUuid, userAccountModel.UserAccounts{userAccount}); err != nil {

		return nil, fmt.Errorf("get user account via accounting: %w", err)
	} else if len(userAccounts) == 0 {

		return nil, fmt.Errorf("get user account via accounting: no user account found")
	} else if len(userAccounts) > 1 {

		return nil, fmt.Errorf("get user account via accounting: more than one user")
	} else {

		return &userAccounts[0].ID, nil
	}
}

func (s *addressServiceImp) validatePreparedNodeAddress(addressType addressPb.AddressType, networkGroup nodeCommon.NetworkGroupEnum, isUserAccountDefined bool) error {
	// BtcNetworkGroupId  NetworkGroupEnum = "btc"
	// BchNetworkGroupId  NetworkGroupEnum = "bch"
	// LtcNetworkGroupId  NetworkGroupEnum = "ltc"
	// DashNetworkGroupId NetworkGroupEnum = "dash"
	// EthNetworkGroupId  NetworkGroupEnum = "eth"
	// DogeNetworkGroupId NetworkGroupEnum = "doge"
	// TrxNetworkGroupId  NetworkGroupEnum = "trx"
	// TonNetworkGroupId  NetworkGroupEnum = "ton"
	// KasNetworkGroupId  NetworkGroupEnum = "kas"

	switch addressType {
	case addressPb.AddressType_ADDRESS_TYPE_DIRECT:
		if networkGroup != nodeCommon.TrxNetworkGroupId {

			return fmt.Errorf("unsupported network group %v for address type %v", networkGroup, addressType)
		} else {

			return nil
		}
	case addressPb.AddressType_ADDRESS_TYPE_DERIVED:

		return fmt.Errorf("unsupported network type %v for this way", networkGroup)
	case addressPb.AddressType_ADDRESS_TYPE_BASED_ID:
		if !isUserAccountDefined {

			return fmt.Errorf("required old user id for this address type %v", addressType)
		} else if networkGroup != nodeCommon.EthNetworkGroupId {

			return fmt.Errorf("unsupported network group %v for address type %v", networkGroup, addressType)
		} else {

			return nil
		}
	case addressPb.AddressType_ADDRESS_TYPE_BASED_NONE:
		if !isUserAccountDefined {

			return fmt.Errorf("required old user id for this address type %v", addressType)
		} else if networkGroup != nodeCommon.BtcNetworkGroupId &&
			networkGroup != nodeCommon.BchNetworkGroupId &&
			networkGroup != nodeCommon.LtcNetworkGroupId &&
			networkGroup != nodeCommon.DogeNetworkGroupId &&
			networkGroup != nodeCommon.KasNetworkGroupId {

			return fmt.Errorf("unsupported network group %v for address type %v", networkGroup, addressType)
		} else {

			return nil
		}
	case addressPb.AddressType_ADDRESS_TYPE_MEMO:

		return fmt.Errorf("unsupported network type %v for this way", networkGroup)
	default:
		return fmt.Errorf("unsupported address type %v", addressType)

	}
}

func (s *addressServiceImp) sendConsumerNetworkReloadAddr(ctx context.Context, network nodeCommon.NetworkEnum) {
	message := &AdminMessage{
		Command:      reloadAddressesCommand,
		TxID:         "",
		Hash:         "",
		ToBlockScore: 0,
	}

	if msgByte, err := json.Marshal(message); err != nil {

		sdkLog.Error(ctx, "marshal admin message: %w", err)
	} else {
		routingKey := fmt.Sprintf("%s.admin", network.ToString()) // TODO: id?

		if err := s.rabbitService.Publish(ctx, routingKey, msgByte); err != nil {

			sdkLog.Error(ctx, "publish admin message: %w", err)
		}
	}
}

func (s *addressServiceImp) sendConsumerNetworkGroupReloadAddr(ctx context.Context, networkGroup nodeCommon.NetworkGroupEnum) {
	for _, network := range networkGroup.GetNetworks() {
		s.sendConsumerNetworkReloadAddr(ctx, network)

	}
}

func (s *addressServiceImp) GetNewAddressByConstraint(ctx context.Context, userUuid uuid.UUID, addressType addressPb.AddressType, networkGroup nodeCommon.NetworkGroupEnum) (model.Addresses, error) {
	filter := &model.AddressFilter{
		Id:           nil,
		Address:      nil,
		UserUuid:     &userUuid,
		IsProcessing: utils.BoolToPtr(false),
		AddressType:  addressType.Enum(),
		NetworkGroup: &networkGroup,
		CreatedAtGt:  nil,
		Pagination:   nil,
	}

	if _, addresses, err := s.addressRepo.GetNewAddresses(ctx, filter); err != nil {

		return nil, fmt.Errorf("get new addresses: %w", err)
	} else if len(addresses) > 1 {

		return nil, fmt.Errorf("can not determinate new address by constraint: %v, %v, %v", userUuid, addressType, networkGroup)
	} else {

		return addresses, nil
	}
}

func (s *addressServiceImp) GetOldAddressByConstraint(ctx context.Context, userUuid uuid.UUID, network nodeCommon.NetworkEnum, coin string) (model.AddressesOld, error) {
	filter := &model.AddressOldFilter{
		Id:            nil,
		Address:       nil,
		UserUuid:      &userUuid,
		Network:       &network,
		AddressType:   nil,
		UserAccountId: nil,
		Coin:          &coin,
		CreatedAtGt:   nil,
		Pagination:    nil,
	}

	if _, addressesOld, err := s.addressRepo.GetOldAddresses(ctx, filter); err != nil {

		return nil, fmt.Errorf("get new addresses: %w", err)
	} else if len(addressesOld) > 1 {

		return nil, fmt.Errorf("can not determinate old address by constraint: %v, %v, %v", userUuid, network, coin)
	} else {

		return addressesOld, nil
	}
}

func (s *addressServiceImp) GetNewAddressByUuid(ctx context.Context, addressUuid uuid.UUID) (model.Addresses, error) {
	filter := &model.AddressFilter{
		Id:           &addressUuid,
		Address:      nil,
		UserUuid:     nil,
		IsProcessing: utils.BoolToPtr(false),
		AddressType:  nil,
		NetworkGroup: nil,
		CreatedAtGt:  nil,
		Pagination:   nil,
	}

	if _, addresses, err := s.addressRepo.GetNewAddresses(ctx, filter); err != nil {

		return nil, fmt.Errorf("get new addresses: %w", err)
	} else {

		return addresses, nil
	}
}

func (s *addressServiceImp) GetNewAddressByStr(ctx context.Context, addressStr string) (model.Addresses, error) {
	filter := &model.AddressFilter{
		Id:           nil,
		Address:      &addressStr,
		UserUuid:     nil,
		IsProcessing: utils.BoolToPtr(false),
		AddressType:  nil,
		NetworkGroup: nil,
		CreatedAtGt:  nil,
		Pagination:   nil,
	}

	if _, addresses, err := s.addressRepo.GetNewAddresses(ctx, filter); err != nil {

		return nil, fmt.Errorf("get new addresses: %w", err)
	} else {

		return addresses, nil
	}
}

func (s *addressServiceImp) GetNewAddressesByUserUuid(ctx context.Context, userUuid uuid.UUID) (model.Addresses, error) {
	filter := &model.AddressFilter{
		Id:           nil,
		Address:      nil,
		UserUuid:     &userUuid,
		IsProcessing: utils.BoolToPtr(false),
		AddressType:  nil,
		NetworkGroup: nil,
		CreatedAtGt:  nil,
		Pagination:   nil,
	}

	if _, addresses, err := s.addressRepo.GetNewAddresses(ctx, filter); err != nil {

		return nil, fmt.Errorf("get new addresses: %w", err)
	} else {

		return addresses, nil
	}
}

func (s *addressServiceImp) GetOldAddressByUuid(ctx context.Context, addressUuid uuid.UUID) (model.AddressesOld, error) {
	filter := &model.AddressOldFilter{
		Id:            &addressUuid,
		Address:       nil,
		UserUuid:      nil,
		AddressType:   nil,
		Network:       nil,
		UserAccountId: nil,
		Coin:          nil,
		CreatedAtGt:   nil,
		Pagination:    nil,
	}

	if _, addresses, err := s.addressRepo.GetOldAddresses(ctx, filter); err != nil {

		return nil, fmt.Errorf("get old addresses: %w", err)
	} else {

		return addresses, nil
	}
}

func (s *addressServiceImp) GetOldAddressByStr(ctx context.Context, addressStr string) (model.AddressesOld, error) {
	filter := &model.AddressOldFilter{
		Id:            nil,
		Address:       &addressStr,
		UserUuid:      nil,
		AddressType:   nil,
		Network:       nil,
		UserAccountId: nil,
		Coin:          nil,
		CreatedAtGt:   nil,
		Pagination:    nil,
	}

	if _, addresses, err := s.addressRepo.GetOldAddresses(ctx, filter); err != nil {

		return nil, fmt.Errorf("get old addresses: %w", err)
	} else {

		return addresses, nil
	}
}

func (s *addressServiceImp) GetOldAddressesByUserUuid(ctx context.Context, userUuid uuid.UUID) (model.AddressesOld, error) {
	filter := &model.AddressOldFilter{
		Id:            nil,
		Address:       nil,
		UserUuid:      &userUuid,
		AddressType:   nil,
		Network:       nil,
		UserAccountId: nil,
		Coin:          nil,
		CreatedAtGt:   nil,
		Pagination:    nil,
	}

	if _, addresses, err := s.addressRepo.GetOldAddresses(ctx, filter); err != nil {

		return nil, fmt.Errorf("get old addresses: %w", err)
	} else {

		return addresses, nil
	}
}

func (s *addressServiceImp) GetNewAddressesByFilter(ctx context.Context, filter *model.AddressFilter) (*uint64, model.Addresses, error) {
	if totalCount, addresses, err := s.addressRepo.GetNewAddresses(ctx, filter); err != nil {

		return nil, nil, fmt.Errorf("get new addresses: %w", err)
	} else {

		return totalCount, addresses, nil
	}
}

func (s *addressServiceImp) GetOldAddressesByFilter(ctx context.Context, filter *model.AddressOldFilter) (*uint64, model.AddressesOld, error) {
	if totalCount, addresses, err := s.addressRepo.GetOldAddresses(ctx, filter); err != nil {

		return nil, nil, fmt.Errorf("get old addresses: %w", err)
	} else {

		return totalCount, addresses, nil
	}
}

func (s *addressServiceImp) CreatePersonalAddress(ctx context.Context, addressStr string, userUuid uuid.UUID, network nodeCommon.NetworkEnum, minPayout float64) (*model.AddressPersonal, error) {
	addressUuid := uuid.New()

	address := model.NewAddressPersonal(addressUuid, addressStr, userUuid, network, minPayout)
	if err := s.addressRepo.AddPersonalAddress(ctx, address); err != nil {

		return nil, fmt.Errorf("add personal address: %w", err)
	} else {

		return address, nil
	}
}

func (s *addressServiceImp) UpdatePersonalAddress(ctx context.Context, address *model.AddressPersonal, addressStr string, minPayout *float64) (*model.AddressPersonal, error) {
	partial := &model.AddressPersonalPartial{
		Address:   &addressStr,
		MinPayout: minPayout,
		DeletedAt: &sql.NullTime{Time: time.Time{}, Valid: false},
		UpdatedAt: utils.TimeToPtr(time.Now().UTC()),
	}

	if err := s.addressRepo.UpdatePersonalAddress(ctx, address, partial); err != nil {

		return nil, fmt.Errorf("update personal address: %w", err)
	} else {

		return address, nil
	}
}

// func (s *addressServiceImp) AddOrUpdatePersonalAddress(ctx context.Context, addressStr string, userUuid uuid.UUID, network nodeCommon.NetworkEnum, minPayout *float64) (*model.AddressPersonal, error) {
// 	filter := &model.AddressPersonalFilter{
// 		Id:         nil,
// 		Address:    nil,
// 		UserUuid:   &userUuid,
// 		Network:    network.ToPtr(),
// 		IsDeleted:  nil,
// 		Pagination: nil,
// 	}
//
// 	if _, addressList, err := s.addressRepo.GetPersonalAddresses(ctx, filter); err != nil {
//
// 		return nil, fmt.Errorf("get personal addresses: %w", err)
// 	} else if len(addressList) > 1 {
//
// 		return nil, fmt.Errorf("multiple personal addresses") // panic, unexpected
// 	} else if len(addressList) == 1 { // update one
// 		address := addressList[0]
//
// 		partial := &model.AddressPersonalPartial{
// 			Address:   &addressStr,
// 			MinPayout: minPayout,
// 			DeletedAt: &sql.NullTime{Time: time.Time{}, Valid: false},
// 			UpdatedAt: utils.TimeToPtr(time.Now().UTC()),
// 		}
//
// 		if err := s.addressRepo.UpdatePersonalAddress(ctx, address, partial); err != nil {
//
// 			return nil, fmt.Errorf("update personal address: %w", err)
// 		} else {
//
// 			return address, nil
// 		}
// 	} else { // add new one
// 		addressUuid := uuid.New()
//
// 		if defaultMinPayout, err := s.getDefaultMinPayout(ctx, network); err != nil {
//
// 			return nil, fmt.Errorf("get min payout: %w", err)
// 		} else {
// 			var addressMinPayout float64
// 			if minPayout == nil {
// 				addressMinPayout = defaultMinPayout
//
// 			} else {
// 				addressMinPayout = math.Max(*minPayout, defaultMinPayout)
//
// 			}
//
// 			address := model.NewAddressPersonal(addressUuid, addressStr, userUuid, network, addressMinPayout)
// 			if err := s.addressRepo.AddPersonalAddress(ctx, address); err != nil {
//
// 				return nil, fmt.Errorf("add personal address: %w", err)
// 			} else {
//
// 				return address, nil
// 			}
// 		}
// 	}
// }

func (s *addressServiceImp) DeletePersonalAddress(ctx context.Context, userUuid uuid.UUID, network nodeCommon.NetworkEnum) error {
	filter := &model.AddressPersonalFilter{
		Id:         nil,
		Address:    nil,
		UserUuid:   &userUuid,
		Network:    network.ToPtr(),
		IsDeleted:  nil,
		Pagination: nil,
	}

	if _, addressList, err := s.addressRepo.GetPersonalAddresses(ctx, filter); err != nil {

		return fmt.Errorf("get personal addresses: %w", err)
	} else if len(addressList) > 1 {

		return fmt.Errorf("multiple personal addresses") // panic, unexpected
	} else if len(addressList) == 1 { // update one
		address := addressList[0]
		dt := time.Now().UTC()

		partial := &model.AddressPersonalPartial{
			Address:   utils.StringToPtr(""),
			DeletedAt: &sql.NullTime{Time: dt, Valid: true},
			MinPayout: nil,
			UpdatedAt: &dt,
		}

		if err := s.addressRepo.UpdatePersonalAddress(ctx, address, partial); err != nil {

			return fmt.Errorf("update personal address: %w", err)
		} else {

			return nil
		}
	} else { // not found

		return nil
	}
}

func (s *addressServiceImp) GetPersonalAddressesByFilter(ctx context.Context, filter *model.AddressPersonalFilter) (*uint64, model.AddressesPersonal, error) {
	if totalCount, addresses, err := s.addressRepo.GetPersonalAddresses(ctx, filter); err != nil {

		return nil, nil, fmt.Errorf("get personal addresses: %w", err)
	} else {

		return totalCount, addresses, nil
	}
}

func (s *addressServiceImp) GetPersonalAddressByConstraint(ctx context.Context, userUuid uuid.UUID, network nodeCommon.NetworkEnum) (model.AddressesPersonal, error) {
	filter := &model.AddressPersonalFilter{
		Id:         nil,
		Address:    nil,
		UserUuid:   &userUuid,
		Network:    network.ToPtr(),
		IsDeleted:  nil,
		Pagination: nil,
	}

	if _, addressList, err := s.addressRepo.GetPersonalAddresses(ctx, filter); err != nil {

		return nil, fmt.Errorf("get personal addresses: %w", err)
	} else if len(addressList) > 1 {

		return nil, fmt.Errorf("multiple personal addresses") // panic, unexpected
	} else if len(addressList) == 0 {
		// success
		return nil, nil
	} else {
		// success
		return addressList, nil
	}
}

func (s *addressServiceImp) CreateOrUpdateDirtyAddress(ctx context.Context, address *model.AddressDirty) (*model.AddressDirty, error) {
	if err := s.addressRepo.AddOrUpdateDirtyAddress(ctx, address); err != nil {

		return nil, fmt.Errorf("add or update dirty address: %w", err)
	} else {

		return address, nil
	}
}

func (s *addressServiceImp) GetDirtyAddressesByFilter(ctx context.Context, filter *model.AddressDirtyFilter) (model.AddressesDirty, error) {
	if addresses, err := s.addressRepo.GetDirtyAddresses(ctx, filter); err != nil {

		return nil, fmt.Errorf("failed get dirt addresses: %w", err)
	} else {

		return addresses, nil
	}
}
