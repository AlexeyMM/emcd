package enum

import (
	"fmt"

	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

type AddressTypeWrapper struct {
	addressPb.AddressType
}

func NewAddressTypeWrapper(addressType addressPb.AddressType) AddressTypeWrapper {

	return AddressTypeWrapper{AddressType: addressType}
}

func (w *AddressTypeWrapper) Scan(value any) error {
	if value == nil {

		return fmt.Errorf("empty value of address type")
	} else {
		var e int32

		switch v := value.(type) {
		case int8:
			e = int32(v)
		case int16:
			e = int32(v)
		case int32:
			e = v
		case int64:
			e = int32(v)
		case int:
			e = int32(v)
		default:
			return fmt.Errorf("invalid type of address type")

		}

		if _, ok := addressPb.AddressType_name[e]; !ok {

			return fmt.Errorf("invalid value of address type: %v", value)
		} else {
			w.AddressType = addressPb.AddressType(e)

			return nil
		}
	}
}
