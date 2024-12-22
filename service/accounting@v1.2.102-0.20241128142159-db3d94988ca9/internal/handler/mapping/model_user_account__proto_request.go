package mapping

import (
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func MapModelUserAccountToProtoRequest(u *model.UserAccount) *userAccountPb.UserAccountRequest {
	var address *string
	if !u.Address.Valid {
		address = nil

	} else {
		address = &u.Address.String

	}

	var img1 *float64
	if !u.Img1.Valid {
		img1 = nil

	} else {
		img1 = &u.Img1.Float64

	}

	var img2 *float64
	if !u.Img2.Valid {
		img2 = nil

	} else {
		img2 = &u.Img2.Float64

	}

	var fee *float64
	if !u.Fee.Valid {
		fee = nil

	} else {
		fee = &u.Fee.Float64

	}

	return &userAccountPb.UserAccountRequest{
		AccountTypeId: u.AccountTypeID.ToInt32(),
		Minpay:        u.Minpay,
		Address:       address,
		Img1:          img1,
		Img2:          img2,
		Fee:           fee,
		CoinNew:       u.CoinNew.String,
	}
}

func MapModelUserAccountToProtoOneRequest(u *model.UserAccount) *userAccountPb.UserAccountOneRequest {
	var address *string
	if !u.Address.Valid {
		address = nil

	} else {
		address = &u.Address.String

	}

	var img1 *float64
	if !u.Img1.Valid {
		img1 = nil

	} else {
		img1 = &u.Img1.Float64

	}

	var img2 *float64
	if !u.Img2.Valid {
		img2 = nil

	} else {
		img2 = &u.Img2.Float64

	}

	var fee *float64
	if !u.Fee.Valid {
		fee = nil

	} else {
		fee = &u.Fee.Float64

	}

	return &userAccountPb.UserAccountOneRequest{
		UserId:        u.UserID,
		UserIdNew:     u.UserIDNew.UUID.String(),
		AccountTypeId: u.AccountTypeID.ToInt32(),
		Minpay:        u.Minpay,
		Address:       address,
		Img1:          img1,
		Img2:          img2,
		Fee:           fee,
		CoinNew:       u.CoinNew.String,
	}
}

func MapModelUserAccountsToProtoMultiRequest(
	userId int32,
	userIdNew uuid.UUID,
	us model.UserAccounts,
) (*userAccountPb.UserAccountMultiRequest, error) {
	var dumps []*userAccountPb.UserAccountRequest

	for _, u := range us {
		if userId != u.UserID {

			return nil, fmt.Errorf("mismatch user_id: %d, %d", userId, u.UserID)
		} else if !u.UserIDNew.Valid || userIdNew.String() != u.UserIDNew.UUID.String() {

			return nil, fmt.Errorf("mismatch user_id_new: %s, %s", userIdNew.String(), u.UserIDNew.UUID.String())
		} else {

			dumps = append(dumps, MapModelUserAccountToProtoRequest(u))
		}
	}

	return &userAccountPb.UserAccountMultiRequest{
		UserId:       userId,
		UserIdNew:    userIdNew.String(),
		UserAccounts: dumps,
	}, nil
}

func MapProtoRequestToModelUserAccount(
	coinValidator coinValidatorRepo.CoinValidatorRepository,
	userId int32,
	userIdNew uuid.UUID,
	p *userAccountPb.UserAccountRequest,
) (*model.UserAccount, error) {
	accountTypeId := enum.NewAccountTypeId(p.AccountTypeId)
	if err := accountTypeId.Validate(); err != nil {

		return nil, err
	}

	var address sql.NullString
	if p.Address == nil {
		address = sql.NullString{}

	} else {
		address = utils.StringToStringNull(*p.Address)

	}

	var img1 sql.NullFloat64
	if p.Img1 == nil {
		img1 = sql.NullFloat64{}

	} else {
		img1 = utils.Float64ToFloat64Null(*p.Img1)

	}

	var img2 sql.NullFloat64
	if p.Img2 == nil {
		img2 = sql.NullFloat64{}

	} else {
		img2 = utils.Float64ToFloat64Null(*p.Img2)

	}

	var fee sql.NullFloat64
	if p.Fee == nil {
		fee = sql.NullFloat64{}

	} else {
		fee = utils.Float64ToFloat64Null(*p.Fee)

	}

	var coinId int32
	var coinNew string
	if id, ok := coinValidator.GetIdByCode(strings.ToLower(p.CoinNew)); !ok {

		return nil, fmt.Errorf("invalid coin_new: %s", p.CoinNew)
	} else {
		coinNew = p.CoinNew
		coinId = id

	}

	return &model.UserAccount{
		ID:            0,
		UserID:        userId,
		CoinID:        coinId,
		AccountTypeID: enum.NewAccountTypeIdWrapper(accountTypeId),
		Minpay:        p.Minpay,
		Address:       address,
		ChangedAt:     sql.NullTime{},
		Img1:          img1,
		Img2:          img2,
		IsActive:      sql.NullBool{},
		CreatedAt:     sql.NullTime{},
		UpdatedAt:     sql.NullTime{},
		Fee:           fee,
		UserIDNew:     uuid.NullUUID{UUID: userIdNew, Valid: true},
		CoinNew:       utils.StringToStringNull(coinNew),
	}, nil
}

func MapProtoOneRequestToModelUserAccount(
	coinValidator coinValidatorRepo.CoinValidatorRepository,
	p *userAccountPb.UserAccountOneRequest,
) (*model.UserAccount, error) {
	var userIdNew uuid.UUID
	if userIdNewParsed, err := uuid.Parse(p.UserIdNew); err != nil {

		return nil, err
	} else {
		userIdNew = userIdNewParsed

	}

	accountTypeId := enum.NewAccountTypeId(p.AccountTypeId)
	if err := accountTypeId.Validate(); err != nil {

		return nil, err
	}

	var address sql.NullString
	if p.Address == nil {
		address = sql.NullString{}

	} else {
		address = utils.StringToStringNull(*p.Address)

	}

	var img1 sql.NullFloat64
	if p.Img1 == nil {
		img1 = sql.NullFloat64{}

	} else {
		img1 = utils.Float64ToFloat64Null(*p.Img1)

	}

	var img2 sql.NullFloat64
	if p.Img2 == nil {
		img2 = sql.NullFloat64{}

	} else {
		img2 = utils.Float64ToFloat64Null(*p.Img2)

	}

	var fee sql.NullFloat64
	if p.Fee == nil {
		fee = sql.NullFloat64{}

	} else {
		fee = utils.Float64ToFloat64Null(*p.Fee)

	}

	var coinId int32
	var coinNew string
	if id, ok := coinValidator.GetIdByCode(strings.ToLower(p.CoinNew)); !ok {

		return nil, fmt.Errorf("invalid coin_new: %s", p.CoinNew)
	} else {
		coinNew = p.CoinNew
		coinId = id

	}

	return &model.UserAccount{
		ID:            0,
		UserID:        p.UserId,
		CoinID:        coinId,
		AccountTypeID: enum.NewAccountTypeIdWrapper(accountTypeId),
		Minpay:        p.Minpay,
		Address:       address,
		ChangedAt:     sql.NullTime{},
		Img1:          img1,
		Img2:          img2,
		IsActive:      sql.NullBool{},
		CreatedAt:     sql.NullTime{},
		UpdatedAt:     sql.NullTime{},
		Fee:           fee,
		UserIDNew:     uuid.NullUUID{UUID: userIdNew, Valid: true},
		CoinNew:       utils.StringToStringNull(coinNew),
	}, nil
}

func MapProtoMultiRequestToModelUserAccounts(
	coinValidator coinValidatorRepo.CoinValidatorRepository,
	ps *userAccountPb.UserAccountMultiRequest,
) (int32, uuid.UUID, model.UserAccounts, error) {
	var dumps model.UserAccounts

	var userIdNew uuid.UUID
	if userIdNewParsed, err := uuid.Parse(ps.UserIdNew); err != nil {

		return 0, uuid.UUID{}, nil, fmt.Errorf("failed parse user_id_new: %s, %w", ps.UserIdNew, err)
	} else {
		userIdNew = userIdNewParsed

	}

	for _, p := range ps.UserAccounts {
		if d, err := MapProtoRequestToModelUserAccount(coinValidator, ps.UserId, userIdNew, p); err != nil {

			return 0, uuid.UUID{}, nil, err
		} else {
			dumps = append(dumps, d)

		}
	}

	return ps.UserId, userIdNew, dumps, nil
}
