package repository

import (
	"code.emcdtech.com/emcd/sdk/log"
	coinPb "code.emcdtech.com/emcd/service/coin/protocol/coin"
	"context"
	"fmt"
	"sync"
	"time"
)

type CoinValidatorRepository interface {
	Serve(ctx context.Context, wg *sync.WaitGroup, cacheUpdateInterval time.Duration)

	GetCodeById(coinIdLegacy int32) (string, bool)
	GetIdByCode(coinCode string) (int32, bool)
	IsValidIdLegacy(coinIdLegacy int32) bool
	IsValidCode(coinCode string) bool
	GetCodes() []string
	GetIdsLegacy() []int32
}

type coinValidatorImp struct {
	coinCli     coinPb.CoinServiceClient
	mapCodeById sync.Map // map[int32]string
	mapIdByCode sync.Map // map[string]int32
}

func NewCoinValidatorRepository(coinHandler coinPb.CoinServiceClient) CoinValidatorRepository {

	return &coinValidatorImp{
		coinCli:     coinHandler,
		mapCodeById: sync.Map{},
		mapIdByCode: sync.Map{},
	}
}

func (v *coinValidatorImp) Serve(ctx context.Context, wg *sync.WaitGroup, cacheUpdateInterval time.Duration) {
	defer wg.Done()
	defer func() {
		log.Info(context.Background(), "coin validator stopped")

	}()

	log.Info(context.Background(), "coin validator started")

	if err := v.updateMaps(ctx); err != nil {
		log.Error(ctx, err.Error())

	}

	t := time.NewTicker(cacheUpdateInterval)

	for {
		select {
		case <-ctx.Done():
			t.Stop()

			return
		case <-t.C:
			if err := v.updateMaps(ctx); err != nil {
				log.Error(ctx, err.Error())

			}
		}
	}
}

func (v *coinValidatorImp) GetCodeById(coinIdLegacy int32) (string, bool) {
	if v, ok := v.mapCodeById.Load(coinIdLegacy); ok {

		return v.(string), true
	} else {

		return "", false
	}
}

func (v *coinValidatorImp) GetIdByCode(coinCode string) (int32, bool) {
	if v, ok := v.mapIdByCode.Load(coinCode); ok {

		return v.(int32), true
	} else {

		return 0, false
	}
}

func (v *coinValidatorImp) IsValidIdLegacy(coinIdLegacy int32) bool {
	_, ok := v.mapCodeById.Load(coinIdLegacy)

	return ok
}

func (v *coinValidatorImp) IsValidCode(coinCode string) bool {
	_, ok := v.mapIdByCode.Load(coinCode)

	return ok
}

func (v *coinValidatorImp) GetCodes() []string {
	codes := make([]string, 0)

	v.mapIdByCode.Range(func(key, value any) bool {
		codes = append(codes, key.(string))

		return true
	})

	return codes
}

func (v *coinValidatorImp) GetIdsLegacy() []int32 {
	ids := make([]int32, 0)

	v.mapCodeById.Range(func(key, value any) bool {
		ids = append(ids, key.(int32))

		return true
	})

	return ids
}

func (v *coinValidatorImp) updateMaps(ctx context.Context) error {
	req := &coinPb.GetCoinsRequest{
		Limit:  9999999,
		Offset: 0,
		WlId:   nil,
	}

	if coinsResponse, err := v.coinCli.GetCoins(ctx, req); err != nil {

		return fmt.Errorf("failed update maps: %w", err)
	} else {
		for _, coin := range coinsResponse.Coins {
			idLegacy := coin.GetLegacyCoinId()
			code := coin.GetId()

			v.mapCodeById.Store(idLegacy, code)
			v.mapIdByCode.Store(code, idLegacy)
		}

		return nil
	}
}
