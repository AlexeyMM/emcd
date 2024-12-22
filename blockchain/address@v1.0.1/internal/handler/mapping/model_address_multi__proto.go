package mapping

import (
	"context"
	"fmt"

	sdkLog "code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

func MapModelOneOfAddressToProto(ctx context.Context, addressesNewOne model.Addresses, addressOldOne model.AddressesOld) *addressPb.AddressResponse {
	if len(addressesNewOne) != 0 && len(addressOldOne) != 0 {
		sdkLog.Warn(ctx, "new and old address are mutually exclusive: %+v, %+v", addressesNewOne[0], addressOldOne[0])

		return MapModelAddressNewToProto(addressesNewOne[0]) // new prefer
	} else if len(addressesNewOne) != 0 {

		return MapModelAddressNewToProto(addressesNewOne[0])
	} else if len(addressOldOne) != 0 {

		return MapModelAddressOldToProto(addressOldOne[0])
	} else {

		return nil
	}
}

func MapModelAddressesToProto(totalCount *uint64, addressesNew model.Addresses, addressesOld model.AddressesOld) *addressPb.AddressMultiResponse {
	var dumps []*addressPb.AddressResponse

	for _, addr := range addressesNew {
		dumps = append(dumps, MapModelAddressNewToProto(addr))

	}

	for _, addr := range addressesOld {
		dumps = append(dumps, MapModelAddressOldToProto(addr))

	}

	return &addressPb.AddressMultiResponse{
		Addresses:  dumps,
		TotalCount: totalCount,
	}
}

func MapModelDirtyAddressesToProto(addressesDirty model.AddressesDirty) *addressPb.DirtyAddressMultiForm {
	var dumps []*addressPb.DirtyAddressForm

	for _, addr := range addressesDirty {
		dumps = append(dumps, MapModelAddressDirtyToProto(addr))

	}

	return &addressPb.DirtyAddressMultiForm{
		Addresses: dumps,
	}
}

func MapProtoAddressResponsesToModel(p *addressPb.AddressResponse) (*model.Address, *model.AddressOld, error) {
	if p.GetNewWay() != nil {
		if addressNew, err := MapProtoAddressResponseToModelNew(p); err != nil {

			return nil, nil, fmt.Errorf("mapping new way: %w", err)
		} else {

			return addressNew, nil, nil
		}
	} else {
		if addressOld, err := MapProtoAddressResponseToModelOld(p); err != nil {

			return nil, nil, fmt.Errorf("mapping old way: %w", err)
		} else {

			return nil, addressOld, nil
		}
	}
}

func MapProtoAddressMultiResponsesToModel(p *addressPb.AddressMultiResponse) (*uint64, model.Addresses, model.AddressesOld, error) {
	var dumpsNew model.Addresses
	var dumpsOld model.AddressesOld

	for _, address := range p.Addresses {
		if dumpNew, dumpOld, err := MapProtoAddressResponsesToModel(address); err != nil {

			return nil, nil, nil, fmt.Errorf("multi: %w", err)
		} else {
			if dumpNew != nil {
				dumpsNew = append(dumpsNew, dumpNew)

			}

			if dumpOld != nil {
				dumpsOld = append(dumpsOld, dumpOld)

			}
		}
	}

	return p.TotalCount, dumpsNew, dumpsOld, nil
}

func MapProtoAddressMultiFormToModelDirty(p *addressPb.DirtyAddressMultiForm) (model.AddressesDirty, error) {
	var dumps model.AddressesDirty

	for _, addr := range p.Addresses {
		if addressModel, err := MapProtoAddressFormToModelDirty(addr); err != nil {

			return nil, fmt.Errorf("mapping multi addresses: %w", err)
		} else {
			dumps = append(dumps, addressModel)

		}
	}

	return dumps, nil
}
