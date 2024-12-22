package model

import (
	"context"
	"database/sql"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/accounting/model/enum"
)

type UserAccount struct {
	ID            int32
	UserID        int32
	CoinID        int32
	AccountTypeID enum.AccountTypeIdWrapper
	Minpay        float64
	Address       sql.NullString
	ChangedAt     sql.NullTime
	Img1          sql.NullFloat64
	Img2          sql.NullFloat64
	IsActive      sql.NullBool
	CreatedAt     sql.NullTime
	UpdatedAt     sql.NullTime
	Fee           sql.NullFloat64
	UserIDNew     uuid.NullUUID
	CoinNew       sql.NullString
}

type UserAccounts []*UserAccount

type UserAccountFilter struct {
	ID              *int32
	UserID          *int32
	AccountTypeID   *enum.AccountTypeId
	UserIDNew       *uuid.UUID
	CoinNew         *string
	IsActive        *bool
	Pagination      *Pagination
	UserIDNewIsNull *struct{}
	CoinNewIsNull   *struct{}
}

type UserAccountPartial struct {
	UserIDNew *uuid.UUID
	CoinNew   *string
}

func (u *UserAccount) Update(partial *UserAccountPartial) {
	if partial.UserIDNew != nil {
		u.UserIDNew = uuid.NullUUID{
			UUID:  *partial.UserIDNew,
			Valid: true,
		}
	}

	if partial.CoinNew != nil {
		u.CoinNew = sql.NullString{
			String: *partial.CoinNew,
			Valid:  true,
		}
	}
}

func (u *UserAccount) SqlNoEmptyValidate() error {
	if !u.UserIDNew.Valid {

		log.Warn(context.Background(), "user_id_new is empty: %d, %s, %s", u.UserID, u.CoinNew, u.AccountTypeID.ToString())

		// linter // return fmt.Errorf("user_id_new is empty")
	} else if !u.CoinNew.Valid {

		log.Warn(context.Background(), "coin_new is empty: %d, %d, %s", u.UserID, u.CoinID, u.AccountTypeID.ToString())
		// linter // return fmt.Errorf("coin_new is empty")
	} else {

		// return nil
	}

	return nil
}

func (us UserAccounts) GetIdList() []int32 {
	idList := make([]int32, len(us))
	for i := range us {
		idList[i] = us[i].ID

	}

	return idList
}

func (us UserAccounts) SqlNoEmptyValidate() error {
	for _, u := range us {
		if err := u.SqlNoEmptyValidate(); err != nil {

			return err
		}
	}

	return nil
}
