package mapping

import (
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	userAccountPb "code.emcdtech.com/emcd/service/accounting/protocol/user_account"
	coinValidatorRepo "code.emcdtech.com/emcd/service/coin/repository"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func MapProtoUserAccountConstraintRequestToArgs(coinValidator coinValidatorRepo.CoinValidatorRepository, p *userAccountPb.UserAccountConstraintRequest) (uuid.UUID, string, enum.AccountTypeId, error) {
	var userIdNew uuid.UUID
	if uuidParsed, err := uuid.Parse(p.UserIdNew); err != nil {

		return uuid.UUID{}, "", 0, fmt.Errorf("failed parse user_id_new: %s, %w", p.UserIdNew, err)
	} else {
		userIdNew = uuidParsed

	}

	userAccountTypeId := enum.AccountTypeId(p.AccountTypeId)
	if err := userAccountTypeId.Validate(); err != nil {

		return uuid.UUID{}, "", 0, fmt.Errorf("invalid account_type_id: %d, %w", p.AccountTypeId, err)
	} else {
		// pass
	}

	var coinNew string
	if ok := coinValidator.IsValidCode(strings.ToLower(p.CoinNew)); !ok {

		return uuid.UUID{}, "", 0, fmt.Errorf("invalid coin_new: %s", p.CoinNew)
	} else {
		coinNew = p.CoinNew

	}

	return userIdNew, coinNew, userAccountTypeId, nil
}
