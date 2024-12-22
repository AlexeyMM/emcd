package service

import (
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"code.emcdtech.com/emcd/service/accounting/internal/client"
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
)

type LimitPayouts interface {
	CheckLimit(ctx context.Context, userID int64, coinID string, amount float64) error
	GetBlockStatus(ctx context.Context, userID int64) (int, error)
	SetBlockStatus(ctx context.Context, userID int64, status int) error
}

type limitPayouts struct {
	LimitPayouts repository.LimitPayouts
	slack        client.SlackNotifier
	SlackUrl     string
}

func NewLimitPayouts(limitPayoutsRepository repository.LimitPayouts, slackUrl string) LimitPayouts {
	return &limitPayouts{
		LimitPayouts: limitPayoutsRepository,
		slack:        client.SlackNotifier{},
		SlackUrl:     slackUrl,
	}
}

func (s *limitPayouts) CheckLimit(ctx context.Context, userID int64, coinID string, amount float64) error {

	coinId, err := strconv.Atoi(coinID)
	if err != nil {
		return fmt.Errorf("invalid coinId: %w", err)
	}

	if userID == 0 {
		return fmt.Errorf("empty UserId")
	}

	mainUserId := userID

	parentId, err := s.LimitPayouts.GetMainUserId(ctx, userID)
	if err != nil {
		return fmt.Errorf("srv get main user id: %w", err)
	}
	var mainAccMessage string
	if parentId > 0 {
		mainUserId = parentId
		mainAccMessage = fmt.Sprintf("(Main Acc: %d)", mainUserId)
	}

	status, err := s.LimitPayouts.GetRedisBlockStatus(ctx, mainUserId)
	if err != nil {
		return fmt.Errorf("repo redis: %w", err)
	}

	switch status {
	case model.StatusBlocked:
		return fmt.Errorf("user is already blocked")
	case model.StatusUnblocked:
		return nil
	}

	result, err := s.checkUserPayoutsSum(ctx, userID, coinId, amount)

	if err != nil {
		return fmt.Errorf("srv check user payouts sum: %w", err)
	}

	if !result.ChkResult {
		err = s.LimitPayouts.SetRedisBlockStatus(ctx, mainUserId, model.StatusBlocked, 0)
		if err != nil {
			return fmt.Errorf("user %d set blocked status is fail: %w", userID, err)
		}

		err = s.LimitPayouts.SetUserNopay(ctx, mainUserId)
		if err != nil {
			return fmt.Errorf("user %d set nopay is fail: %w", userID, err)
		}

		blockMes := fmt.Sprintf("user %d %s has been blocked (nopay). CoinId: %d, Limit: %f, Paid per day: %f (including unprocessed), Requested amount: %f",
			userID, mainAccMessage, coinId, result.Limit, result.PayoutsSum, amount)
		m := fmt.Sprintf(
			"NEW PAYOUTS LIMIT BLOCK:\n"+
				"Info: %s", blockMes)
		err := s.slack.SendMessage([]string{s.SlackUrl}, m)
		if err != nil {
			sdkLog.Error(ctx, "send slack message: %v", err)
		}

		return fmt.Errorf("new limit block: %s", blockMes)
	}

	return nil
}

func (s *limitPayouts) GetBlockStatus(ctx context.Context, userID int64) (int, error) {

	if userID == 0 {
		return 0, fmt.Errorf("empty UserId")
	}

	parentId, err := s.LimitPayouts.GetMainUserId(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("srv get main user id: %w", err)
	}
	if parentId > 0 {
		userID = parentId
	}

	result, err := s.LimitPayouts.GetRedisBlockStatus(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("repo GetRedisBlockStatus: %w", err)
	}

	return result, nil
}

func (s *limitPayouts) SetBlockStatus(ctx context.Context, userID int64, status int) error {
	if userID == 0 {
		return errors.New("empty userId")
	}

	parentId, err := s.LimitPayouts.GetMainUserId(ctx, userID)
	if err != nil {
		return fmt.Errorf("srv get main user id: %w", err)
	}
	if parentId > 0 {
		userID = parentId
	}

	exp := time.Duration(0)
	if status == model.StatusUnblocked {
		exp = time.Hour * 24
	}

	err = s.LimitPayouts.SetRedisBlockStatus(ctx, userID, status, exp)
	if err != nil {
		return fmt.Errorf("repo SetRedisBlockStatus: %w", err)
	}
	return nil
}

type checkResult struct {
	Limit      float64
	PayoutsSum float64
	ChkResult  bool
}

func (s *limitPayouts) checkUserPayoutsSum(ctx context.Context, userID int64, coinId int, _ float64) (checkResult, error) {
	var checkResult checkResult
	var err error

	limit, err := s.LimitPayouts.GetLimit(ctx, coinId)
	if err != nil {
		checkResult.ChkResult = false
		return checkResult, err
	}
	checkResult.Limit = limit

	if limit == 0 {
		checkResult.ChkResult = true
		return checkResult, nil
	}

	payoutsSum, err := s.LimitPayouts.GetUserPayoutsSum(ctx, userID, coinId)
	if err != nil {
		return checkResult, fmt.Errorf("user %d get paytous sum: %w", userID, err)
	}
	checkResult.PayoutsSum = payoutsSum

	if payoutsSum > limit {
		checkResult.ChkResult = false
		return checkResult, nil
	}

	checkResult.ChkResult = true
	return checkResult, nil
}
