package enum

import (
	"fmt"
)

type AccountTypeId int32
type AccountTypeIds []AccountTypeId

// DB date
// 1	Platform
// 2	Pool
// 3	Coinhold
// 4	Referral
// 5	Block
// 6	Hedge
// 7	P2P

const (
	// WalletAccountTypeID Идентификатор типа кошельковых счетов.
	WalletAccountTypeID AccountTypeId = iota + 1
	// MiningAccountTypeID Идентификатор типа майнинговых счетов.
	MiningAccountTypeID
	// CoinholdAccountTypeID Идентификатор типа счёта для счетов Коинхолда.
	CoinholdAccountTypeID
	// ReferralAccountTypeID Идентификатор типа реферальных счетов.
	ReferralAccountTypeID
	// BlockUserAccountTypeID Идентификатор типа блокировочных счетов.
	BlockUserAccountTypeID
	// HedgeUserAccountTypeID Идентификатор типа хеджевых счетов (у нас он только у пользователя).
	HedgeUserAccountTypeID
	// P2PUserAccountTypeIDNotUsedBefore Идентификатор типа P2P счетов
	P2PUserAccountTypeIDNotUsedBefore
)

func NewAccountTypeId(v int32) AccountTypeId {

	return AccountTypeId(v)
}

func NewAccountTypeIds(vs []int32) AccountTypeIds {
	var as AccountTypeIds
	for _, v := range vs {
		as = append(as, AccountTypeId(v))

	}

	return as
}

func (a AccountTypeId) Validate() error {
	switch a {
	case WalletAccountTypeID, MiningAccountTypeID, CoinholdAccountTypeID, ReferralAccountTypeID,
		BlockUserAccountTypeID, P2PUserAccountTypeIDNotUsedBefore:

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

func (as AccountTypeIds) Validate() error {
	for _, a := range as {
		if err := a.Validate(); err != nil {

			return err
		}
	}

	return nil
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
	case P2PUserAccountTypeIDNotUsedBefore:
		return "P2P"
	default:

		// panic(fmt.Sprintf("unknown account_type_id: %s", a)) // TODO: may be remove comment?

		return ""

	}
}

func parseAccountTypeIdFromString(a string) (AccountTypeId, error) {
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
		return P2PUserAccountTypeIDNotUsedBefore, nil
	default:
		return 0, fmt.Errorf("invalid account_type_id format: %s", a)

	}
}
