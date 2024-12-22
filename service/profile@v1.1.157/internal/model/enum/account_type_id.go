package enum

import (
	"errors"
	"fmt"
)

type AccountTypeId int32
type AccountTypeIdWrapper struct {
	AccountTypeId
}

//DB date
//1	Platform
//2	Pool
//3	Coinhold
//4	Referral
//5	Block
//6	Hedge
//7	P2P

const (
	// WalletAccountTypeID wallet accounts identifier.
	WalletAccountTypeID AccountTypeId = iota + 1
	// MiningAccountTypeID mining accounts identifier.
	MiningAccountTypeID
	// CoinholdAccountTypeID coinhold accounts identifier..
	CoinholdAccountTypeID
	// ReferralAccountTypeID referral accounts identifier.
	ReferralAccountTypeID
	// BlockUserAccountTypeID blockusers accounts identifier.
	BlockUserAccountTypeID
	// HedgeUserAccountTypeID hedge accounts identifier.
	HedgeUserAccountTypeID
	// P2P p2p accounts identifier.
	P2PUserAccountTypeID_NotUsedBefore
)

func NewAccountTypeId(v int32) AccountTypeId {

	return AccountTypeId(v)
}
func NewAccountTypeIdWrapper(v int32) AccountTypeIdWrapper {

	return AccountTypeIdWrapper{NewAccountTypeId(v)}
}

func (a AccountTypeId) Validate() error {
	switch a {
	case WalletAccountTypeID, MiningAccountTypeID, CoinholdAccountTypeID, ReferralAccountTypeID, BlockUserAccountTypeID, P2PUserAccountTypeID_NotUsedBefore:

		return nil
	default:

		return fmt.Errorf("invalid account type: %d", a)
	}
}

func (a AccountTypeId) ToInt() int {

	return int(a)
}

func (a AccountTypeId) ToInt32() int32 {

	return int32(a)
}

func (a AccountTypeId) ToPtr() *AccountTypeId {

	return &a
}

func (a AccountTypeId) ToInt32Ptr() *int32 {
	v := int32(a)

	return &v
}

func (a AccountTypeId) ToString() string {
	switch a {
	case WalletAccountTypeID:
		return "Platform"
	case MiningAccountTypeID:
		return "Pool"
	case CoinholdAccountTypeID:
		return "Coinhold"
	case ReferralAccountTypeID:
		return "Referral"
	case BlockUserAccountTypeID:
		return "Block"
	case HedgeUserAccountTypeID:
		return "Hedge"
	case P2PUserAccountTypeID_NotUsedBefore:
		return "P2P"
	default:

		//panic(fmt.Sprintf("unknown account_type_id: %s", a)) // TODO: may be remove comment?

		return ""

	}
}

func ParseAccountTypeIdFromString(a string) (AccountTypeId, error) {
	switch a {
	case "Wallet":
		return WalletAccountTypeID, nil
	case "Mining":
		return MiningAccountTypeID, nil
	case "Coinhold":
		return CoinholdAccountTypeID, nil
	case "Referral":
		return ReferralAccountTypeID, nil
	case "Block":
		return BlockUserAccountTypeID, nil
	case "Hedge":
		return HedgeUserAccountTypeID, nil
	case "P2P":
		return P2PUserAccountTypeID_NotUsedBefore, nil
	default:
		return 0, fmt.Errorf("invalid account_type_id format: %s", a)

	}
}

func (w *AccountTypeIdWrapper) Scan(value any) error {
	if value == nil {

		return errors.New("empty value of account type id")
	} else if v, ok := value.(int8); ok {
		w.AccountTypeId = NewAccountTypeId(int32(v))

		return w.AccountTypeId.Validate()
	} else if v, ok := value.(int16); ok {
		w.AccountTypeId = NewAccountTypeId(int32(v))

		return w.AccountTypeId.Validate()
	} else if v, ok := value.(int32); ok {
		w.AccountTypeId = NewAccountTypeId(v)

		return w.AccountTypeId.Validate()
	} else if v, ok := value.(int64); ok {
		w.AccountTypeId = NewAccountTypeId(int32(v))

		return w.AccountTypeId.Validate()
	} else if v, ok := value.(string); ok {
		if p, e := ParseAccountTypeIdFromString(v); e == nil {
			w.AccountTypeId = p

			return nil
		} else {

			return e
		}
	} else {

		return errors.New("invalid type of account type id")
	}
}
