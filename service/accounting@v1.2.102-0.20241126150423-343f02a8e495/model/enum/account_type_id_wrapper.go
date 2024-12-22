package enum

import "errors"

type AccountTypeIdWrapper struct {
	AccountTypeId
}

func NewAccountTypeIdWrapper(id AccountTypeId) AccountTypeIdWrapper {

	return AccountTypeIdWrapper{AccountTypeId: id}
}

// func NewAccountTypeIdWrapperInt32(v int32) AccountTypeIdWrapper {
//
// 	return AccountTypeIdWrapper{NewAccountTypeId(v)}
// }

func (w *AccountTypeIdWrapper) Scan(value any) error {
	if value == nil {

		return errors.New("empty value of account type id")
	} else {
		switch v := value.(type) {
		case int8:
			w.AccountTypeId = NewAccountTypeId(int32(v))

			return nil
		case int16:
			w.AccountTypeId = NewAccountTypeId(int32(v))

			return nil
		case int32:
			w.AccountTypeId = NewAccountTypeId(v)

			return nil
		case int64:
			w.AccountTypeId = NewAccountTypeId(int32(v))

			return nil
		case int:
			w.AccountTypeId = NewAccountTypeId(int32(v))

			return nil
		case string:
			if p, e := parseAccountTypeIdFromString(v); e == nil {
				w.AccountTypeId = p

				return nil
			} else {

				return e
			}
		default:
			return errors.New("invalid type of account type id")

		}
	}
}
