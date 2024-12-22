package model

const (
	btcMinpayFrom  = 0.0001
	bchMinpayFrom  = 0.01
	ltcMinpayFrom  = 0.05
	dashMinpayFrom = 0.1
	etcMinpayFrom  = 0.1
	dogeMinpayFrom = 1
	kasMinPayFrom  = 50
	belMinPayFrom  = 1
	fbMinPayFrom   = 0.05
)

var MinPayDefault = map[int]float64{
	btcCoinID:  btcMinpayFrom,
	bchCoinID:  bchMinpayFrom,
	ltcCoinID:  ltcMinpayFrom,
	dashCoinID: dashMinpayFrom,
	etcCoinID:  etcMinpayFrom,
	dogeCoinID: dogeMinpayFrom,
	kasCoinID:  kasMinPayFrom,
	belCoinID:  belMinPayFrom,
	fbCoinID:   fbMinPayFrom,
}

func IsMinPayValid(amount float32, coinID int) bool {
	var isValid bool
	switch coinID {
	case btcCoinID:
		isValid = amount >= btcMinpayFrom
	case bchCoinID:
		isValid = amount >= bchMinpayFrom
	case ltcCoinID:
		isValid = amount >= ltcMinpayFrom
	case dashCoinID:
		isValid = amount >= dashMinpayFrom
	case etcCoinID:
		isValid = amount >= etcMinpayFrom
	case dogeCoinID:
		isValid = amount >= dogeMinpayFrom
	case kasCoinID:
		isValid = amount >= kasMinPayFrom
	default:
		isValid = amount >= 0
	}
	return isValid
}
