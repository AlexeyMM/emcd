package mapping

import (
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

func MapModelUserAccountToProtoResponse(u *model.UserAccount) *userAccountPb.UserAccountResponse {
	var address *string
	if !u.Address.Valid {
		address = nil

	} else {
		address = &u.Address.String

	}

	var changedAt *timestamppb.Timestamp
	if !u.ChangedAt.Valid {
		changedAt = nil

	} else {
		changedAt = timestamppb.New(u.ChangedAt.Time)

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

	var isActive *bool
	if !u.IsActive.Valid {
		isActive = nil

	} else {
		isActive = &u.IsActive.Bool

	}

	var createdAt *timestamppb.Timestamp
	if !u.CreatedAt.Valid {
		createdAt = nil

	} else {
		createdAt = timestamppb.New(u.CreatedAt.Time)

	}

	var updatedAt *timestamppb.Timestamp
	if !u.UpdatedAt.Valid {
		updatedAt = nil

	} else {
		updatedAt = timestamppb.New(u.UpdatedAt.Time)

	}

	var fee *float64
	if !u.Fee.Valid {
		fee = nil

	} else {
		fee = &u.Fee.Float64

	}

	return &userAccountPb.UserAccountResponse{
		Id:            u.ID,
		UserId:        u.UserID,
		CoinId:        u.CoinID,
		AccountTypeId: u.AccountTypeID.ToInt32(),
		Minpay:        u.Minpay,
		Address:       address,
		ChangedAt:     changedAt,
		Img1:          img1,
		Img2:          img2,
		IsActive:      isActive,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Fee:           fee,
		UserIdNew:     u.UserIDNew.UUID.String(),
		CoinNew:       u.CoinNew.String,
	}
}

func MapModelUserAccountsToProtoMultiResponse(totalCount *uint64, us model.UserAccounts) *userAccountPb.UserAccountMultiResponse {
	var dumps []*userAccountPb.UserAccountResponse

	for _, u := range us {
		dumps = append(dumps, MapModelUserAccountToProtoResponse(u))

	}

	return &userAccountPb.UserAccountMultiResponse{
		UserAccounts: dumps,
		TotalCount:   totalCount,
	}
}

func MapProtoResponseToModelUserAccount(ctx context.Context, coinValidator coinValidatorRepo.CoinValidatorRepository, p *userAccountPb.UserAccountResponse) *model.UserAccount {
	accountTypeId := enum.NewAccountTypeId(p.AccountTypeId)
	if err := accountTypeId.Validate(); err != nil {
		sdkLog.Warn(ctx, "migrate warning: failed to validate account_type_id: %d", p.AccountTypeId)

	}

	var address sql.NullString
	if p.Address == nil {
		address = sql.NullString{}

	} else {
		address = utils.StringToStringNull(*p.Address)

	}

	var changedAt sql.NullTime
	if p.ChangedAt == nil {
		changedAt = sql.NullTime{}

	} else {
		changedAt = utils.TimeToTimeNull(p.ChangedAt.AsTime())

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

	var isActive sql.NullBool
	if p.IsActive == nil {
		isActive = sql.NullBool{}

	} else {
		isActive = utils.BoolToBoolNull(*p.IsActive)

	}

	var createdAt sql.NullTime
	if p.CreatedAt == nil {
		createdAt = sql.NullTime{}

	} else {
		createdAt = utils.TimeToTimeNull(p.CreatedAt.AsTime())

	}

	var updatedAt sql.NullTime
	if p.UpdatedAt == nil {
		updatedAt = sql.NullTime{}

	} else {
		updatedAt = utils.TimeToTimeNull(p.UpdatedAt.AsTime())

	}

	var fee sql.NullFloat64
	if p.Fee == nil {
		fee = sql.NullFloat64{}

	} else {
		fee = utils.Float64ToFloat64Null(*p.Fee)

	}

	var coinId int32
	var coinNew string
	// coinValidator mapping required for get coinId, this request will be deprecated
	if id, ok := coinValidator.GetIdByCode(strings.ToLower(p.CoinNew)); !ok {
		sdkLog.Warn(ctx, "migrate warning: failed to convert coin_new: %d, %s", p.Id, p.CoinNew)

	} else {
		coinNew = p.CoinNew
		coinId = id

	}

	var userIdNew uuid.UUID
	if userIdNewParsed, err := uuid.Parse(p.UserIdNew); err != nil {
		sdkLog.Warn(ctx, "migrate warning: failed to convert user_id_new: %d, %s", p.Id, p.UserIdNew)

	} else {
		userIdNew = userIdNewParsed

	}

	return &model.UserAccount{
		ID:            p.Id,
		UserID:        p.UserId,
		CoinID:        coinId,
		AccountTypeID: enum.NewAccountTypeIdWrapper(accountTypeId),
		Minpay:        p.Minpay,
		Address:       address,
		ChangedAt:     changedAt,
		Img1:          img1,
		Img2:          img2,
		IsActive:      isActive,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Fee:           fee,
		UserIDNew:     utils.UuidToUuidNull(userIdNew),
		CoinNew:       utils.StringToStringNull(coinNew),
	}
}

func MapProtoMultiResponseToModelUserAccounts(ctx context.Context, coinValidator coinValidatorRepo.CoinValidatorRepository, ps []*userAccountPb.UserAccountResponse) model.UserAccounts {
	var dumps model.UserAccounts

	for _, p := range ps {
		dumps = append(dumps, MapProtoResponseToModelUserAccount(ctx, coinValidator, p))
	}

	return dumps
}
