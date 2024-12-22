package mapping

import (
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func MapProtoToModelUserAccountFilter(coinValidator coinValidatorRepo.CoinValidatorRepository, p *userAccountPb.UserAccountFilter) (*model.UserAccountFilter, error) {
	var accountTypeId *enum.AccountTypeId
	if p.AccountTypeId == nil {
		accountTypeId = nil

	} else {
		v := enum.AccountTypeId(*p.AccountTypeId)
		if err := v.Validate(); err != nil {

			return nil, err
		} else {
			accountTypeId = &v

		}
	}

	var userIdNew *uuid.UUID
	if p.UserIdNew == nil {
		userIdNew = nil

	} else {
		if userIdNewParsed, err := uuid.Parse(*p.UserIdNew); err != nil {

			return nil, fmt.Errorf("failed parse user_id_new: %s, %w", *p.UserIdNew, err)
		} else {
			userIdNew = &userIdNewParsed

		}
	}

	var coinNew *string
	if p.CoinNew == nil {
		coinNew = nil

	} else {
		if ok := coinValidator.IsValidCode(strings.ToLower(*p.CoinNew)); !ok {

			return nil, fmt.Errorf("invalid coin_new: %s", *p.CoinNew)
		} else {
			coinNew = p.CoinNew

		}
	}

	var pagination *model.Pagination
	if p.Pagination == nil {
		pagination = nil

	} else {
		pagination = &model.Pagination{
			Limit:  p.Pagination.Limit,
			Offset: p.Pagination.Offset,
		}
	}

	return &model.UserAccountFilter{
		ID:              p.Id,
		UserID:          p.UserId,
		AccountTypeID:   accountTypeId,
		UserIDNew:       userIdNew,
		CoinNew:         coinNew,
		IsActive:        p.IsActive,
		Pagination:      pagination,
		UserIDNewIsNull: nil,
		CoinNewIsNull:   nil,
	}, nil
}

func MapModelUserAccountFilterToProto(filter *model.UserAccountFilter) *userAccountPb.UserAccountFilter {
	var accountTypeId *int32
	if filter.AccountTypeID == nil {
		accountTypeId = nil

	} else {
		accountTypeId = filter.AccountTypeID.ToInt32Ptr()

	}

	var userIdNew *string
	if filter.UserIDNew == nil {
		userIdNew = nil

	} else {
		userIdNew = utils.StringToPtr(filter.UserIDNew.String())

	}

	var pagination *userAccountPb.UserAccountPagination
	if filter.Pagination == nil {
		pagination = nil

	} else {
		pagination = &userAccountPb.UserAccountPagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	return &userAccountPb.UserAccountFilter{
		Id:            filter.ID,
		UserId:        filter.UserID,
		AccountTypeId: accountTypeId,
		UserIdNew:     userIdNew,
		CoinNew:       filter.CoinNew,
		IsActive:      filter.IsActive,
		Pagination:    pagination,
	}
}
