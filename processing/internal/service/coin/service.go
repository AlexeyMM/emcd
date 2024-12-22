package coin

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/processing/internal/client"
	"code.emcdtech.com/b2b/processing/internal/repository"
	"code.emcdtech.com/b2b/processing/internal/service"
	"code.emcdtech.com/b2b/processing/model"
)

type Service struct {
	coinRepository repository.Coin
	coinClient     client.CoinClient
	availableCoins map[string]struct{}
}

func NewCoinService(
	coinRepository repository.Coin,
	coinClient client.CoinClient,
	availableCoins []string,
) service.CoinService {
	availableCoinsSet := make(map[string]struct{}, len(availableCoins))
	for _, coin := range availableCoins {
		availableCoinsSet[coin] = struct{}{}
	}

	return &Service{coinRepository: coinRepository, coinClient: coinClient, availableCoins: availableCoinsSet}
}

func (s *Service) GetCoins() []*model.Coin {
	return s.coinRepository.GetCoins()
}

func (s *Service) UpdateCoins(ctx context.Context) error {
	coins, err := s.coinClient.GetCoins(ctx)
	if err != nil {
		return fmt.Errorf("updateCoins: %w", err)
	}

	coinsToStore := make([]*model.Coin, 0, len(s.availableCoins))

	for _, coin := range coins {
		if _, ok := s.availableCoins[coin.ID]; ok {
			coinsToStore = append(coinsToStore, coin)
		}
	}

	s.coinRepository.SetCoins(coinsToStore)

	return nil
}
