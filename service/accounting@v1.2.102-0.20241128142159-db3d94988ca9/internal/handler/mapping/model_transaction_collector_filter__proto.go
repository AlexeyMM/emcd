package mapping

import (
	"context"
	"fmt"
	"time"

	sdkLog "code.emcdtech.com/emcd/sdk/log"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"google.golang.org/protobuf/types/known/timestamppb"

	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	"code.emcdtech.com/emcd/service/accounting/model"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
)

func MapProtoToModelTransactionCollectorFilter(coinValidator coinValidatorRepo.CoinValidatorRepository, p *accountingPb.TransactionCollectorFilter) (*model.TransactionCollectorFilter, error) {
	var coinId *int32
	if p.CoinCode == nil {
		coinId = nil

	} else {
		if coinIdConverted, ok := coinValidator.GetIdByCode(*p.CoinCode); !ok {

			return nil, fmt.Errorf("invalid coin_code: %s", *p.CoinCode)
		} else {
			coinId = &coinIdConverted

		}
	}

	var createdAtGt *time.Time
	if p.CreatedAtGt == nil {
		createdAtGt = nil

	} else {
		createdAtGt = utils.TimeToPtr(p.CreatedAtGt.AsTime())

	}

	var createdAtLte *time.Time
	if p.CreatedAtLte == nil {
		createdAtLte = nil

	} else {
		createdAtLte = utils.TimeToPtr(p.CreatedAtLte.AsTime())

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

	return &model.TransactionCollectorFilter{
		Types:        p.Types,
		CoinId:       coinId,
		CreatedAtGt:  createdAtGt,
		CreatedAtLte: createdAtLte,
		Pagination:   pagination,
	}, nil
}

func MapModelToProtoTransactionCollectorFilter(ctx context.Context, coinValidator coinValidatorRepo.CoinValidatorRepository, filter *model.TransactionCollectorFilter) *accountingPb.TransactionCollectorFilter {
	var coinCode *string
	if filter.CoinId == nil {
		coinCode = nil

	} else {
		if coinCodeConverted, ok := coinValidator.GetCodeById(*filter.CoinId); !ok {
			sdkLog.Warn(ctx, "failed to convert coin_id to code: %d", *filter.CoinId)

			coinCode = nil
		} else {
			coinCode = &coinCodeConverted

		}
	}

	var createdAtGt *timestamppb.Timestamp
	if filter.CreatedAtGt == nil {
		createdAtGt = nil

	} else {
		createdAtGt = timestamppb.New(*filter.CreatedAtGt)

	}

	var createdAtLte *timestamppb.Timestamp
	if filter.CreatedAtLte == nil {
		createdAtLte = nil

	} else {
		createdAtLte = timestamppb.New(*filter.CreatedAtLte)

	}

	var pagination *accountingPb.Pagination
	if filter.Pagination == nil {
		pagination = nil

	} else {
		pagination = &accountingPb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	return &accountingPb.TransactionCollectorFilter{
		Types:        filter.Types,
		CoinCode:     coinCode,
		CreatedAtGt:  createdAtGt,
		CreatedAtLte: createdAtLte,
		Pagination:   pagination,
	}
}
