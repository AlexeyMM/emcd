package mapping

import (
	"database/sql"
	"fmt"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapModelAddressPersonalToProto(address *model.AddressPersonal) *addressPb.PersonalAddressResponse {
	var deletedAt *timestamppb.Timestamp
	if address.DeletedAt.Valid {
		deletedAt = timestamppb.New(address.DeletedAt.Time)

	} else {
		deletedAt = nil

	}

	ret := &addressPb.PersonalAddressResponse{
		AddressUuid: address.Id.String(),
		Address:     address.Address,
		UserUuid:    address.UserUuid.String(),
		Network:     address.Network.ToString(),
		MinPayout:   address.MinPayout,
		DeletedAt:   deletedAt,
		UpdatedAt:   timestamppb.New(address.UpdatedAt),
		CreatedAt:   timestamppb.New(address.CreatedAt),
	}

	return ret
}

func MapProtoAddressResponseToModelPersonal(p *addressPb.PersonalAddressResponse) (*model.AddressPersonal, error) {
	var addressUuid uuid.UUID
	if addressUuidParsed, err := uuid.Parse(p.AddressUuid); err != nil {

		return nil, fmt.Errorf("failed parse address id: %s, %w", p.UserUuid, err)
	} else {
		addressUuid = addressUuidParsed

	}

	var userUuid uuid.UUID
	if userUuidParsed, err := uuid.Parse(p.UserUuid); err != nil {

		return nil, fmt.Errorf("failed parse user_uuid: %s, %w", p.UserUuid, err)
	} else {
		userUuid = userUuidParsed

	}

	var network nodeCommon.NetworkEnum
	networkNew := nodeCommon.NewNetworkEnum(p.Network)
	if err := networkNew.Validate(); err != nil {

		return nil, fmt.Errorf("failed parse network: %s, %w", p.Network, err)
	} else {
		network = networkNew

	}

	var deletedAt sql.NullTime
	if p.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: p.DeletedAt.AsTime(), Valid: true}

	} else {
		deletedAt = sql.NullTime{Time: time.Time{}, Valid: false}

	}

	return &model.AddressPersonal{
		Id:        addressUuid,
		Address:   p.Address,
		UserUuid:  userUuid,
		Network:   nodeCommon.NewNetworkEnumWrapper(network),
		MinPayout: p.MinPayout,
		DeletedAt: deletedAt,
		UpdatedAt: p.UpdatedAt.AsTime(),
		CreatedAt: p.CreatedAt.AsTime(),
	}, nil
}

func MapModelAddressesPersonalToProto(totalCount *uint64, addressesPersonal model.AddressesPersonal) *addressPb.PersonalAddressMultiResponse {
	var dumps []*addressPb.PersonalAddressResponse

	for _, addr := range addressesPersonal {
		dumps = append(dumps, MapModelAddressPersonalToProto(addr))

	}

	return &addressPb.PersonalAddressMultiResponse{
		Addresses:  dumps,
		TotalCount: totalCount,
	}
}

func MapProtoPersonalAddressMultiResponsesToModel(p *addressPb.PersonalAddressMultiResponse) (*uint64, model.AddressesPersonal, error) {
	var dumpsPersonal model.AddressesPersonal

	for _, address := range p.Addresses {
		if dump, err := MapProtoAddressResponseToModelPersonal(address); err != nil {

			return nil, nil, fmt.Errorf("multi: %w", err)
		} else {
			dumpsPersonal = append(dumpsPersonal, dump)

		}
	}

	return p.TotalCount, dumpsPersonal, nil
}
