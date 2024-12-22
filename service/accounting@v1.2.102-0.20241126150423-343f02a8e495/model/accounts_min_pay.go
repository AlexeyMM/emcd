package model

const (
	CoinBtcID   = 1
	CoinBCHID   = 2
	CoinLTCID   = 4
	CoinDASHID  = 5
	CoinETHID   = 6
	CoinETCID   = 7
	CoinDOGEID  = 8
	CoinUSDTID  = 10
	CoinUSDCID  = 11
	CoinKaspaID = 14

	btcMinpayFrom  = 0.0001
	bchMinpayFrom  = 0.01
	ltcMinpayFrom  = 0.1
	dashMinpayFrom = 0.0001
	etcMinpayFrom  = 0.1
	dogeMinpayFrom = 1
	kasMinpayFrom  = 50
)

var MinPayDefault = map[int]float64{
	CoinBtcID:   btcMinpayFrom,
	CoinBCHID:   bchMinpayFrom,
	CoinLTCID:   ltcMinpayFrom,
	CoinDASHID:  dashMinpayFrom,
	CoinETCID:   etcMinpayFrom,
	CoinDOGEID:  dogeMinpayFrom,
	CoinKaspaID: kasMinpayFrom,
}
