package mapping

import (
	"code.emcdtech.com/b2b/endpoint/internal/controller/response"
	"code.emcdtech.com/b2b/swap/model"
)

func MapToCoinsResponseDeprecated(coins []*model.Coin) map[string]*response.Coin {
	resp := make(map[string]*response.Coin, len(coins))

	for _, coin := range coins {
		networks := make([]*response.Network, len(coin.Networks))
		for j, network := range coin.Networks {
			networks[j] = &response.Network{
				Title:             network.Title,
				Accuracy:          network.AccuracyWithdrawAndDeposit,
				WithdrawSupported: network.WithdrawSupported,
			}
		}
		resp[coin.Title] = &response.Coin{
			Title:    coin.Title,
			Networks: networks,
			IconUrl:  coin.Info.IconURL,
		}
	}

	return resp
}

func MapToCoinsResponse(coins []*model.Coin) []*response.Coin {
	resp := make([]*response.Coin, len(coins))

	for i, coin := range coins {
		networks := make([]*response.Network, len(coin.Networks))
		for j, network := range coin.Networks {
			networks[j] = &response.Network{
				Title:             network.Title,
				Accuracy:          network.AccuracyWithdrawAndDeposit,
				WithdrawSupported: network.WithdrawSupported,
			}
		}
		resp[i] = &response.Coin{
			Title:    coin.Title,
			Networks: networks,
			IconUrl:  coin.Info.IconURL,
		}
	}

	return resp
}
