package model

const (
	btcCoinID  = 1
	bchCoinID  = 2
	ltcCoinID  = 4
	dashCoinID = 5
	etcCoinID  = 7
	dogeCoinID = 8
	kasCoinID  = 14
	belCoinID  = 20
	fbCoinID   = 21
)

var CoinsWithMiningAccounts = map[int]struct{}{
	btcCoinID:  {},
	bchCoinID:  {},
	ltcCoinID:  {},
	dashCoinID: {},
	etcCoinID:  {},
	dogeCoinID: {},
}
