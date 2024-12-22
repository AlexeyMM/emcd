package handler

import (
	"context"
	"fmt"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	coinPb "code.emcdtech.com/emcd/service/coin/protocol/coin"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/handler/mapping"
	"code.emcdtech.com/emcd/blockchain/address/internal/service"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
	"code.emcdtech.com/emcd/blockchain/address/repository"
)

type AddressHandler struct {
	addressService     service.AddressService
	coinValidator      coinValidatorRepo.CoinValidatorRepository
	coinProtoCli       coinPb.CoinServiceClient
	isNetworkOldWayMap map[nodeCommon.NetworkEnum]bool

	addressPb.UnimplementedAddressServiceServer
}

func NewAddressHandler(
	addressService service.AddressService,
	coinValidator coinValidatorRepo.CoinValidatorRepository,
	coinProtoCli coinPb.CoinServiceClient,
	isNetworkOldWayMap map[nodeCommon.NetworkEnum]bool,
) *AddressHandler {

	return &AddressHandler{
		addressService:                    addressService,
		isNetworkOldWayMap:                isNetworkOldWayMap,
		coinValidator:                     coinValidator,
		coinProtoCli:                      coinProtoCli,
		UnimplementedAddressServiceServer: addressPb.UnimplementedAddressServiceServer{},
	}
}

func (h *AddressHandler) getAddressTypeNewWayByNetworkGroup(networkGroup nodeCommon.NetworkGroupEnum) (*addressPb.AddressType, error) {
	switch networkGroup {
	case nodeCommon.EthNetworkGroupId, nodeCommon.AlphNetworkGroupId:
		return addressPb.AddressType_ADDRESS_TYPE_DERIVED.Enum(), nil
	case nodeCommon.TrxNetworkGroupId:
		return addressPb.AddressType_ADDRESS_TYPE_DIRECT.Enum(), nil
	case nodeCommon.TonNetworkGroupId:
		return addressPb.AddressType_ADDRESS_TYPE_MEMO.Enum(), nil
	case nodeCommon.BtcNetworkGroupId, nodeCommon.BchNetworkGroupId, nodeCommon.DogeNetworkGroupId, nodeCommon.LtcNetworkGroupId,
		nodeCommon.KasNetworkGroupId, nodeCommon.DashNetworkGroupId, nodeCommon.BelNetworkGroupId, nodeCommon.FbNetworkGroupId:
		return addressPb.AddressType_ADDRESS_TYPE_BASED_NONE.Enum(), nil
	default:

		return nil, fmt.Errorf("unsupported network group for create address by new way: %v", networkGroup)
	}
}

func (h *AddressHandler) getAddressTypeOldWayByNetwork(network nodeCommon.NetworkEnum) (*addressPb.AddressType, error) {
	switch network.Group() {
	case nodeCommon.TrxNetworkGroupId:
		return addressPb.AddressType_ADDRESS_TYPE_DIRECT.Enum(), nil
	case nodeCommon.TonNetworkGroupId:
		return addressPb.AddressType_ADDRESS_TYPE_MEMO.Enum(), nil
	case nodeCommon.EthNetworkGroupId:
		return addressPb.AddressType_ADDRESS_TYPE_BASED_ID.Enum(), nil
	case nodeCommon.BtcNetworkGroupId, nodeCommon.BchNetworkGroupId, nodeCommon.DogeNetworkGroupId, nodeCommon.LtcNetworkGroupId,
		nodeCommon.KasNetworkGroupId, nodeCommon.DashNetworkGroupId, nodeCommon.BelNetworkGroupId, nodeCommon.FbNetworkGroupId:
		return addressPb.AddressType_ADDRESS_TYPE_BASED_NONE.Enum(), nil
	default:

		return nil, fmt.Errorf("unsupported network for create address by old way: %v", network)
	}
}

func (h *AddressHandler) GetOrCreateAddress(ctx context.Context, req *addressPb.CreateAddressRequest) (*addressPb.AddressResponse, error) {
	/*
		TODO: после созвона
		GetOrCreate

		enum AddressType {
			// прямая генерация с храннеием приватного ключа в базе ноды
			ADDRESS_TYPE_DIRECT                                 = 0;
			// наследованный адрес от мастер ключа
			ADDRESS_TYPE_DERIVED                                = 1;
			// адрес с приватным ключём на блокчейн-ноде с доступом через user_id соль (старая схема)
			ADDRESS_TYPE_BASED_ID                               = 100[deprecated = true];
			// мемо
			ADDRESS_TYPE_MEMO                                   = 3;
		}

		1.взять из конфига сеть (is_old)
		[ADDRESS_TYPE_DERIVED] - новая
		[ADDRESS_TYPE_DIRECT, ADDRESS_TYPE_BASED_ID, ADDRESS_TYPE_MEMO] - старая

					1.1 если новая, определить группу сетей, если есть вернуть, если нет создать

					1.2 если старая, на основе сети определить метод генеррации,
		на основе койна находится кошельковый акк, если метод его требует.
	*/

	userUuid, network, coin, err := mapping.MapProtoAddressRequestToArgs(h.coinValidator, req)

	if err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1011
	}

	isOldWay := coin != nil                               // if coin is defined then old way
	isOldWay = isOldWay && h.isNetworkOldWayMap[*network] // if old way then old way
	addressCreateUuid := uuid.New()

	if !isOldWay {
		// new way
		// TODO: get new address and insert is not safe
		networkGroup := network.Group()

		if addressTypeNew, err := h.getAddressTypeNewWayByNetworkGroup(networkGroup); err != nil {
			sdkLog.Error(ctx, err.Error())

			return nil, repository.ErrAddr1012
		} else if addressNewGet, err := h.addressService.GetNewAddressByConstraint(ctx, *userUuid, *addressTypeNew, networkGroup); err != nil {
			sdkLog.Error(ctx, err.Error())

			return nil, repository.ErrAddr1013
		} else if len(addressNewGet) == 1 {
			// success

			return mapping.MapModelAddressNewToProto(addressNewGet[0]), nil
		} else if addressNew, err := h.addressService.CreateNewAddress(ctx, addressCreateUuid, *userUuid, *addressTypeNew, networkGroup); err != nil {
			sdkLog.Error(ctx, err.Error())

			return nil, repository.ErrAddr1014
		} else {
			// success

			return mapping.MapModelAddressNewToProto(addressNew), nil
		}
	} else {
		// old way
		// TODO: get old address and insert is not safe
		if addressTypeOld, err := h.getAddressTypeOldWayByNetwork(*network); err != nil {
			sdkLog.Error(ctx, err.Error())

			return nil, repository.ErrAddr1015
		} else if addressOldGet, err := h.addressService.GetOldAddressByConstraint(ctx, *userUuid, *network, *coin); err != nil {
			sdkLog.Error(ctx, err.Error())

			return nil, repository.ErrAddr1016
		} else if len(addressOldGet) == 1 {
			// success

			return mapping.MapModelAddressOldToProto(addressOldGet[0]), nil
		} else if addressOld, err := h.addressService.CreateOldAddress(ctx, addressCreateUuid, *userUuid, *addressTypeOld, *network, *coin); err != nil {
			sdkLog.Error(ctx, err.Error())

			return nil, repository.ErrAddr1017
		} else {

			return mapping.MapModelAddressOldToProto(addressOld), nil
		}
	}
}

func (h *AddressHandler) CreateProcessingAddress(ctx context.Context, req *addressPb.CreateProcessingAddressRequest) (*addressPb.AddressResponse, error) {
	addressCreateUuid := uuid.New()

	if userUuid, network, processingUuid, err := mapping.MapProtoProcessingAddressRequestToArgs(req); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1021
	} else if addressTypeNew, err := h.getAddressTypeNewWayByNetworkGroup(network.Group()); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1022
	} else if addressNew, err := h.addressService.CreateProcessingAddress(ctx, addressCreateUuid, *userUuid, *processingUuid, *addressTypeNew, network.Group()); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1023
	} else {
		// success

		return mapping.MapModelAddressNewToProto(addressNew), nil
	}
}

func (h *AddressHandler) GetAddressByUuid(ctx context.Context, req *addressPb.AddressUuid) (*addressPb.AddressResponse, error) {
	if addressUuid, err := mapping.MapStringToUuid(req.AddressUuid); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1031
	} else if addressesNew, err := h.addressService.GetNewAddressByUuid(ctx, *addressUuid); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1032
	} else if addressesOld, err := h.addressService.GetOldAddressByUuid(ctx, *addressUuid); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1033
	} else if len(addressesOld) > 1 {

		return nil, repository.ErrAddr1034
	} else if len(addressesNew) > 1 {

		return nil, repository.ErrAddr1035
	} else if len(addressesNew) == 0 && len(addressesOld) == 0 {

		return nil, repository.ErrAddr1036
	} else {

		return mapping.MapModelOneOfAddressToProto(ctx, addressesNew, addressesOld), nil
	}
}

func (h *AddressHandler) GetAddressByStr(ctx context.Context, req *addressPb.AddressStrId) (*addressPb.AddressResponse, error) {
	if addressesNew, err := h.addressService.GetNewAddressByStr(ctx, req.Address); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1041
	} else if addressesOld, err := h.addressService.GetOldAddressByStr(ctx, req.Address); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1042
	} else if len(addressesOld) > 1 {

		return nil, repository.ErrAddr1043
	} else if len(addressesNew) > 1 {

		return nil, repository.ErrAddr1044
	} else if len(addressesNew) == 0 && len(addressesOld) == 0 {

		return nil, repository.ErrAddr1045
	} else {

		return mapping.MapModelOneOfAddressToProto(ctx, addressesNew, addressesOld), nil
	}
}

func (h *AddressHandler) GetAddressesByUserUuid(ctx context.Context, req *addressPb.UserUuid) (*addressPb.AddressMultiResponse, error) {
	if userUuid, err := mapping.MapStringToUuid(req.UserUuid); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1051
	} else if addressesNew, err := h.addressService.GetNewAddressesByUserUuid(ctx, *userUuid); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1052
	} else if addressesOld, err := h.addressService.GetOldAddressesByUserUuid(ctx, *userUuid); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1053
	} else {

		return mapping.MapModelAddressesToProto(nil, addressesNew, addressesOld), nil
	}
}

func (h *AddressHandler) GetAddressesOldByFilter(ctx context.Context, req *addressPb.AddressOldFilter) (*addressPb.AddressMultiResponse, error) {
	if filter, err := mapping.MapProtoToModelAddressOldFilter(h.coinValidator, req); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1061

	} else if totalCount, addressesOld, err := h.addressService.GetOldAddressesByFilter(ctx, filter); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1062
	} else {

		return mapping.MapModelAddressesToProto(totalCount, nil, addressesOld), nil
	}
}

func (h *AddressHandler) GetAddressesNewByFilter(ctx context.Context, req *addressPb.AddressNewFilter) (*addressPb.AddressMultiResponse, error) {
	if filter, err := mapping.MapProtoToModelAddressNewFilter(req); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1071

	} else if totalCount, addressesNew, err := h.addressService.GetNewAddressesByFilter(ctx, filter); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1072
	} else {

		return mapping.MapModelAddressesToProto(totalCount, addressesNew, nil), nil
	}
}

func (h *AddressHandler) AddOrUpdatePersonalAddress(ctx context.Context, req *addressPb.CreatePersonalAddressRequest) (*addressPb.PersonalAddressResponse, error) {
	if addressStr, userUuid, network, minPayout, err := mapping.MapProtoPersonalAddressRequestToArgs(req); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1081
	} else if addresses, err := h.addressService.GetPersonalAddressByConstraint(ctx, *userUuid, *network); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1082
	} else {
		if len(addresses) == 0 { // create new one
			if minPayout == nil {
				sdkLog.Error(ctx, repository.ErrAddr1083.Error())

				return nil, repository.ErrAddr1083
			}

			if defaultMinPayout, err := h.getDefaultMinPayout(ctx, *network); err != nil {
				sdkLog.Error(ctx, err.Error())

				return nil, repository.ErrAddr1084
			} else if *minPayout < defaultMinPayout {
				sdkLog.Error(ctx, repository.ErrAddr1085.Error())

				return nil, repository.ErrAddr1085
			} else if address, err := h.addressService.CreatePersonalAddress(ctx, addressStr, *userUuid, *network, *minPayout); err != nil {
				sdkLog.Error(ctx, err.Error())

				return nil, repository.ErrAddr1086

			} else {
				// success create
				return mapping.MapModelAddressPersonalToProto(address), nil
			}
		} else { // update one
			var addressMinPayout *float64
			if minPayout != nil {
				if defaultMinPayout, err := h.getDefaultMinPayout(ctx, *network); err != nil {
					sdkLog.Error(ctx, err.Error())

					return nil, repository.ErrAddr1084
				} else if *minPayout < defaultMinPayout {
					sdkLog.Error(ctx, repository.ErrAddr1085.Error())

					return nil, repository.ErrAddr1085

				} else {
					addressMinPayout = minPayout

				}
			}

			if address, err := h.addressService.UpdatePersonalAddress(ctx, addresses[0], addressStr, addressMinPayout); err != nil {
				sdkLog.Error(ctx, err.Error())

				return nil, repository.ErrAddr1087
			} else {
				// success update
				return mapping.MapModelAddressPersonalToProto(address), nil
			}
		}
	}
}

func (h *AddressHandler) DeletePersonalAddress(ctx context.Context, req *addressPb.DeletePersonalAddressRequest) (*emptypb.Empty, error) {
	if userUuid, network, err := mapping.MapProtoDeletePersonalAddressRequestToArgs(req); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1091
	} else if err := h.addressService.DeletePersonalAddress(ctx, *userUuid, *network); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1092
	} else {
		// success

		return &emptypb.Empty{}, nil
	}
}

func (h *AddressHandler) GetPersonalAddressesByFilter(ctx context.Context, req *addressPb.AddressPersonalFilter) (*addressPb.PersonalAddressMultiResponse, error) {
	if filter, err := mapping.MapProtoToModelAddressPersonalFilter(req); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1101

	} else if totalCount, addressesNew, err := h.addressService.GetPersonalAddressesByFilter(ctx, filter); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1102
	} else {

		return mapping.MapModelAddressesPersonalToProto(totalCount, addressesNew), nil
	}
}

func (h *AddressHandler) getDefaultMinPayout(ctx context.Context, network nodeCommon.NetworkEnum) (float64, error) {
	coinReq := &coinPb.GetCoinRequest{
		CoinId: network.ToString(),
	}
	if coin, err := h.coinProtoCli.GetCoin(ctx, coinReq); err != nil {

		return 0.0, fmt.Errorf("get coin: %w", err)
	} else {
		for _, coinNetwork := range coin.Coin.Networks {
			if coinNetwork.IsMining {

				return coinNetwork.MinpayMining, nil
			}
		}

		return 0.0, fmt.Errorf("coin not found")
	}
}

// rpc GetPersonalAddressesByUserUuid(UserUuid) returns (PersonalAddressMultiResponse);

func (h *AddressHandler) GetPersonalAddressesByUserUuid(ctx context.Context, req *addressPb.UserUuid) (*addressPb.PersonalAddressMultiResponse, error) {
	filterProto := &addressPb.AddressPersonalFilter{
		AddressUuid: nil,
		Address:     nil,
		UserUuid:    &req.UserUuid,
		Network:     nil,
		IsDeleted:   utils.BoolToPtr(false),
		Pagination:  nil,
	}

	if filter, err := mapping.MapProtoToModelAddressPersonalFilter(filterProto); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1201
	} else if _, addressesPersonal, err := h.addressService.GetPersonalAddressesByFilter(ctx, filter); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1202
	} else if len(addressesPersonal) > 1 {
		sdkLog.Error(ctx, repository.ErrAddr1203.Error())

		return nil, repository.ErrAddr1203
	} else { // len in [0, 1]

		return mapping.MapModelAddressesPersonalToProto(nil, addressesPersonal), nil
	}
}

func (h *AddressHandler) CreateOrUpdateDirtyAddress(ctx context.Context, address *addressPb.DirtyAddressForm) (*addressPb.DirtyAddressForm, error) {
	if addressModel, err := mapping.MapProtoAddressFormToModelDirty(address); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1301
	} else if _, err := h.addressService.CreateOrUpdateDirtyAddress(ctx, addressModel); err != nil {
		sdkLog.Error(ctx, "failed create or update: %v", err)

		return nil, repository.ErrAddr1302
	} else {

		return address, nil
	}
}

func (h *AddressHandler) GetDirtyAddressesByFilter(ctx context.Context, filter *addressPb.DirtyAddressFilter) (*addressPb.DirtyAddressMultiForm, error) {
	if filterModel, err := mapping.MapProtoToModelAddressDirtyFilter(filter); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1401
	} else if addressResp, err := h.addressService.GetDirtyAddressesByFilter(ctx, filterModel); err != nil {
		sdkLog.Error(ctx, err.Error())

		return nil, repository.ErrAddr1402
	} else {

		return mapping.MapModelDirtyAddressesToProto(addressResp), nil
	}
}
