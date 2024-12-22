package mapping

import (
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/b2b/swap/protocol/swapCoin"
)

func MapCoinsToProtoGetAllResponse(coins []*model.Coin) *swapCoin.GetAllResponse {
	coinsResp := make([]*swapCoin.Coin, 0, len(coins))

	for i := range coins {
		networks := make([]*swapCoin.Network, 0, len(coins[i].Networks))
		for x := range coins[i].Networks {
			networks = append(networks, &swapCoin.Network{
				Title:             coins[i].Networks[x].Title,
				Accuracy:          int32(coins[i].Networks[x].AccuracyWithdrawAndDeposit),
				WithdrawSupported: coins[i].Networks[x].WithdrawSupported,
			})
		}

		coinsResp = append(coinsResp, &swapCoin.Coin{
			Title:    coins[i].Title,
			Networks: networks,
			IconUrl:  coins[i].Info.IconURL,
		})
	}

	return &swapCoin.GetAllResponse{
		Coins: coinsResp,
	}
}

func MapProtoGetAllResponseToCoins(resp *swapCoin.GetAllResponse) []*model.Coin {
	coins := make([]*model.Coin, 0, len(resp.Coins))

	for i := range resp.Coins {
		networks := make([]*model.Network, 0, len(resp.Coins[i].Networks))
		for x := range resp.Coins[i].Networks {
			networks = append(networks, &model.Network{
				Title:                      resp.Coins[i].Networks[x].Title,
				AccuracyWithdrawAndDeposit: int(resp.Coins[i].Networks[x].Accuracy),
				WithdrawSupported:          resp.Coins[i].Networks[x].WithdrawSupported,
			})
		}

		coins = append(coins, &model.Coin{
			Title:    resp.Coins[i].Title,
			Networks: networks,
			Info: &model.CoinInfo{
				IconURL: resp.Coins[i].IconUrl,
			},
		})
	}

	return coins
}
