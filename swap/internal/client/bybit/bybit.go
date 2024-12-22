package bybit

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/b2b/swap/model"
)

type returnCode int

const (
	recvWindow = "5000"

	returnCodeSuccess returnCode = 0

	returnCodeInvalidSymbol returnCode = 10001

	// Нужно повторить запрос
	returnCodeInvalidTimestamp returnCode = 10002
)

type decimalWrapper decimal.Decimal

func (d *decimalWrapper) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		*d = decimalWrapper(decimal.Zero)
		return nil
	}

	return json.Unmarshal(data, (*decimal.Decimal)(d))
}

type intWrapper int

func (i *intWrapper) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		*i = intWrapper(0)
		return nil
	}

	var (
		jsNum json.Number
		num   int64
	)

	err := json.Unmarshal(data, &jsNum)
	if err != nil {
		return err
	}

	num, err = jsNum.Int64()
	if err != nil {
		return err
	}
	*i = intWrapper(num)

	return nil
}

type ByBit struct {
	client          *resty.Client
	apiUrl          string
	masterUid       int
	masterApiKey    string
	masterApiSecret string
}

func NewByBit(
	masterUid int,
	apiUrl string,
	apiKey string,
	apiSecret string,
) *ByBit {
	return &ByBit{
		client:          resty.New(),
		apiUrl:          apiUrl,
		masterUid:       masterUid,
		masterApiKey:    apiKey,
		masterApiSecret: apiSecret,
	}
}

type (
	getFeeResponse struct {
		RetCode returnCode   `json:"retCode"`
		RetMsg  string       `json:"retMsg"`
		Result  getFeeResult `json:"result"`
		Time    int64        `json:"time"`
	}

	getFeeResult struct {
		List []symbolInfo `json:"list"`
	}

	symbolInfo struct {
		Symbol       string         `json:"symbol"`
		TakerFeeRate decimalWrapper `json:"takerFeeRate"`
		MakerFeeRate decimalWrapper `json:"makerFeeRate"`
	}
)

func (b *ByBit) GetFeeRate(ctx context.Context, symbol string) (*model.Fee, error) {
	url := fmt.Sprintf("%s/v5/account/fee-rate/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	params := map[string]string{
		"symbol":   symbol,
		"category": strings.ToLower(model.Spot),
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, fmt.Sprintf("symbol=%s&category=spot", symbol), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetContext(ctx).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getFeeResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	if result.RetCode != returnCodeSuccess {

		//10002, msg: invalid request, please check your server timestamp or recv_window param
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getFeeRate: error: 10002")
			return b.GetFeeRate(ctx, symbol)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	if len(result.Result.List) == 0 {
		return nil, fmt.Errorf("result.Result.List is empty")
	}

	sym := result.Result.List[0]
	if sym.Symbol != symbol {
		return nil, fmt.Errorf("unexpected symbol: %s", symbol)
	}

	return &model.Fee{
		Symbol:   sym.Symbol,
		MakerFee: decimal.Decimal(sym.MakerFeeRate),
		TakerFee: decimal.Decimal(sym.TakerFeeRate),
	}, nil
}

func (b *ByBit) GetAllFeeRate(ctx context.Context) (map[string]*model.Fee, error) {
	url := fmt.Sprintf("%s/v5/account/fee-rate/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	params := map[string]string{
		"category": strings.ToLower(model.Spot),
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, "category=spot", timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getFeeResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	if result.RetCode != returnCodeSuccess {

		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getAllFeeRate: error: 10002")
			return b.GetAllFeeRate(ctx)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	if len(result.Result.List) == 0 {
		return nil, fmt.Errorf("result.Result.List is empty")
	}

	m := make(map[string]*model.Fee, len(result.Result.List))
	for i := range result.Result.List {
		m[result.Result.List[i].Symbol] = &model.Fee{
			Symbol:   result.Result.List[i].Symbol,
			MakerFee: decimal.Decimal(result.Result.List[i].MakerFeeRate),
			TakerFee: decimal.Decimal(result.Result.List[i].TakerFeeRate),
		}
	}

	return m, nil
}

type getTickerResponse struct {
	RetCode    returnCode      `json:"retCode"`
	RetMsg     string          `json:"retMsg"`
	Result     getTickerResult `json:"result"`
	RetExtInfo struct{}        `json:"retExtInfo"`
	Time       int             `json:"time"`
}

type getTickerResult struct {
	Category string              `json:"category"`
	List     []getTickerListItem `json:"list"`
}

type getTickerListItem struct {
	Symbol        string         `json:"symbol"`
	Bid1Price     decimalWrapper `json:"bid1Price"`
	Bid1Size      string         `json:"bid1Size"`
	Ask1Price     decimalWrapper `json:"ask1Price"`
	Ask1Size      string         `json:"ask1Size"`
	LastPrice     string         `json:"lastPrice"`
	PrevPrice24h  string         `json:"prevPrice24h"`
	Price24hPcnt  string         `json:"price24hPcnt"`
	HighPrice24h  string         `json:"highPrice24h"`
	LowPrice24h   string         `json:"lowPrice24h"`
	Turnover24h   string         `json:"turnover24h"`
	Volume24h     string         `json:"volume24h"`
	UsdIndexPrice string         `json:"usdIndexPrice"`
}

func (b *ByBit) GetTicker(ctx context.Context, symbol string) (*model.Ticker, error) {
	url := fmt.Sprintf("%s/v5/market/tickers/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	params := map[string]string{
		"category": strings.ToLower(model.Spot),
		"symbol":   symbol,
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, fmt.Sprintf("symbol=%s&category=spot", symbol), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getTickerResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getTicker: error: 10002")
			return b.GetTicker(ctx, symbol)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	if len(result.Result.List) == 0 {
		return nil, fmt.Errorf("result.Result.List is empty")
	}

	if result.Result.List[0].Symbol != symbol {
		return nil, fmt.Errorf("unexpected symbol: %s", result.Result.List[0].Symbol)
	}

	return &model.Ticker{
		Symbol: result.Result.List[0].Symbol,
		Ask:    decimal.Decimal(result.Result.List[0].Ask1Price),
		Bid:    decimal.Decimal(result.Result.List[0].Bid1Price),
	}, nil
}

func (b *ByBit) IsTickerExist(ctx context.Context, symbol string) (bool, error) {
	url := fmt.Sprintf("%s/v5/market/tickers/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	params := map[string]string{
		"category": strings.ToLower(model.Spot),
		"symbol":   symbol,
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, fmt.Sprintf("symbol=%s&category=spot", symbol), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetContext(ctx).
		Get(url)

	if err != nil {
		return false, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return false, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getTickerResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return false, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess && result.RetCode != returnCodeInvalidSymbol {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "isTickerExist: error: 10002")
			return b.IsTickerExist(ctx, symbol)
		}

		return false, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	} else if result.RetCode == returnCodeInvalidSymbol {
		return false, nil
	}

	return true, nil
}

type getCoinInfoResponse struct {
	RetCode    returnCode        `json:"retCode"`
	RetMsg     string            `json:"retMsg"`
	Result     getCoinInfoResult `json:"result"`
	RetExtInfo interface{}       `json:"retExtInfo"`
	Time       int64             `json:"time"`
}

type getCoinInfoResult struct {
	Rows []getCoinInfoRow `json:"rows"`
}

type getCoinInfoRow struct {
	Name         string             `json:"name"`
	Coin         string             `json:"coin"`
	RemainAmount string             `json:"remainAmount"`
	Chains       []getCoinInfoChain `json:"chains"`
}

type getCoinInfoChain struct {
	ChainType             string         `json:"chainType"`
	Confirmation          string         `json:"confirmation"`
	WithdrawFee           decimalWrapper `json:"withdrawFee"`
	DepositMin            decimalWrapper `json:"depositMin"`
	WithdrawMin           decimalWrapper `json:"withdrawMin"`
	Chain                 string         `json:"chain"`
	ChainDeposit          string         `json:"chainDeposit"`
	ChainWithdraw         string         `json:"chainWithdraw"`
	MinAccuracy           intWrapper     `json:"minAccuracy"`
	WithdrawPercentageFee decimalWrapper `json:"withdrawPercentageFee"`
}

func (b *ByBit) GetCoinInfo(ctx context.Context) ([]*model.Coin, error) {
	url := fmt.Sprintf("%s/v5/asset/coin/query-info", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, "", timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getCoinInfoResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getCoinInfo: error: 10002")
			return b.GetCoinInfo(ctx)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	coins := make([]*model.Coin, len(result.Result.Rows))
	for i := range result.Result.Rows {
		var networks []*model.Network
		coins[i] = &model.Coin{
			Title:    result.Result.Rows[i].Coin,
			Networks: networks,
		}

		for x := range result.Result.Rows[i].Chains {
			var (
				withdrawSupported bool
			)
			if !decimal.Decimal(result.Result.Rows[i].Chains[x].WithdrawFee).IsZero() {
				withdrawSupported = true
			}

			coins[i].Networks = append(coins[i].Networks, &model.Network{
				Title: result.Result.Rows[i].Chains[x].Chain,
				WithdrawFee: &model.WithdrawFee{
					Fee:           decimal.Decimal(result.Result.Rows[i].Chains[x].WithdrawFee),
					PercentageFee: decimal.Decimal(result.Result.Rows[i].Chains[x].WithdrawPercentageFee),
				},
				WithdrawMin:                decimal.Decimal(result.Result.Rows[i].Chains[x].WithdrawMin),
				DepositMin:                 decimal.Decimal(result.Result.Rows[i].Chains[x].DepositMin),
				AccuracyWithdrawAndDeposit: int(result.Result.Rows[i].Chains[x].MinAccuracy),
				WithdrawSupported:          withdrawSupported,
			})
		}
	}

	return coins, nil
}

type getConvertCoinListResponse struct {
	RetCode    returnCode               `json:"retCode"`
	RetMsg     string                   `json:"retMsg"`
	Result     getConvertCoinListResult `json:"result"`
	RetExtInfo interface{}              `json:"retExtInfo"`
	Time       int64                    `json:"time"`
}

type getConvertCoinListResult struct {
	Coins []getConvertCoinListCoin `json:"coins"`
}

type getConvertCoinListCoin struct {
	Coin               string     `json:"coin"`
	FullName           string     `json:"fullName"`
	Icon               string     `json:"icon"`
	IconNight          string     `json:"iconNight"`
	AccuracyLength     intWrapper `json:"accuracyLength"`
	CoinType           string     `json:"coinType"`
	Balance            string     `json:"balance"`
	UBalance           string     `json:"uBalance"`
	TimePeriod         intWrapper `json:"timePeriod"`
	SingleFromMinLimit string     `json:"singleFromMinLimit"`
	SingleFromMaxLimit string     `json:"singleFromMaxLimit"`
	SingleToMinLimit   string     `json:"singleToMinLimit"`
	SingleToMaxLimit   string     `json:"singleToMaxLimit"`
	DailyFromMinLimit  string     `json:"dailyFromMinLimit"`
	DailyFromMaxLimit  string     `json:"dailyFromMaxLimit"`
	DailyToMinLimit    string     `json:"dailyToMinLimit"`
	DailyToMaxLimit    string     `json:"dailyToMaxLimit"`
	DisableFrom        bool       `json:"disableFrom"`
	DisableTo          bool       `json:"disableTo"`
}

func (b *ByBit) GetConvertCoinList(ctx context.Context) (map[string]int, error) {
	url := fmt.Sprintf("%s/v5/asset/exchange/query-coin-list", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, "accountType=eb_convert_funding", timeStamp),
	}

	params := map[string]string{
		"accountType": "eb_convert_funding",
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetContext(ctx).
		SetQueryParams(params).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getConvertCoinListResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getConvertCoinList: error: 10002")
			return b.GetConvertCoinList(ctx)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	m := make(map[string]int)

	for _, l := range result.Result.Coins {
		m[l.Coin] = int(l.AccuracyLength)
	}
	return m, nil
}

type createSubAccountResponse struct {
	RetCode    returnCode             `json:"retCode"`
	RetMsg     string                 `json:"retMsg"`
	Result     subAccountResult       `json:"result"`
	RetExtInfo map[string]interface{} `json:"retExtInfo"`
	Time       int64                  `json:"time"`
}

type subAccountResult struct {
	UID        intWrapper `json:"uid"`
	Username   string     `json:"username"`
	MemberType intWrapper `json:"memberType"`
	Status     intWrapper `json:"status"`
	Remark     string     `json:"remark"`
}

func (b *ByBit) CreateSubAccount(ctx context.Context) (*model.Account, error) {
	url := fmt.Sprintf("%s/v5/user/create-sub-member/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	username, err := generateRandomString(timeStamp)
	if err != nil {
		return nil, fmt.Errorf("generateRandomString: %w", err)
	}

	params := map[string]string{
		"username":   username,
		"memberType": "1",
		"isUta":      "false",
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, string(body), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetBody(body).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result createSubAccountResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "createSubAccount: error: 10002")
			return b.CreateSubAccount(ctx)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	// 1 - normal, 2 -login banned, 4 - frozen
	isValid := int(result.Result.Status) == 1
	if !isValid {
		log.Error(ctx, "create sub account: %d, status isn't normal: %d", result.Result.UID, result.Result.Status)
	}

	return &model.Account{
		ID:      int64(result.Result.UID),
		IsValid: isValid,
	}, nil
}

type createSubAPIKeyResponse struct {
	RetCode    returnCode            `json:"retCode"`
	RetMsg     string                `json:"retMsg"`
	Result     createSubAPIKeyResult `json:"result"`
	RetExtInfo struct{}              `json:"retExtInfo"`
	Time       int64                 `json:"time"`
}

type createSubAPIKeyResult struct {
	ID          string                     `json:"id"`
	Note        string                     `json:"note"`
	ApiKey      string                     `json:"apiKey"`
	ReadOnly    intWrapper                 `json:"readOnly"`
	Secret      string                     `json:"secret"`
	Permissions createSubAPIKeyPermissions `json:"permissions"`
}

type createSubAPIKeyPermissions struct {
	ContractTrade []string `json:"ContractTrade"`
	Spot          []string `json:"Spot"`
	Wallet        []string `json:"Wallet"`
	Options       []string `json:"Options"`
	CopyTrading   []string `json:"CopyTrading"`
	BlockTrade    []string `json:"BlockTrade"`
	Exchange      []string `json:"Exchange"`
	NFT           []string `json:"NFT"`
}

func (b *ByBit) CreateSubAPIKey(ctx context.Context, accountID int) (*model.Secrets, error) {
	url := fmt.Sprintf("%s/v5/user/create-sub-api/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	type permissions struct {
		ContractTrade []string `json:"ContractTrade"`
		Spot          []string `json:"Spot"`
		Wallet        []string `json:"Wallet"`
		Exchange      []string `json:"Exchange"`
	}

	params := map[string]interface{}{
		"subuid":   strconv.Itoa(accountID),
		"readOnly": "0",
		"permissions": permissions{
			ContractTrade: []string{
				"Order",
				"Position",
			},
			Spot: []string{
				"SpotTrade",
			},
			Wallet: []string{
				"AccountTransfer",
				"SubMemberTransferList",
			},
			Exchange: []string{
				"ExchangeHistory",
			},
		},
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, string(body), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetBody(body).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result createSubAPIKeyResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "createSubApiKey: error: 10002")
			return b.CreateSubAPIKey(ctx, accountID)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	if result.Result.ApiKey == "" {
		return nil, fmt.Errorf("result.Result.ApiKey is empty")
	}

	return &model.Secrets{
		ApiKey:    result.Result.ApiKey,
		ApiSecret: result.Result.Secret,
	}, nil
}

type subMemberListResponse struct {
	RetCode returnCode `json:"retCode"`
	RetMsg  string     `json:"retMsg"`
	Result  struct {
		SubMembers []subMember `json:"subMembers"`
	} `json:"result"`
	RetExtInfo map[string]interface{} `json:"retExtInfo"`
	Time       int64                  `json:"time"`
}

type subMember struct {
	UID         intWrapper `json:"uid"`
	Username    string     `json:"username"`
	MemberType  intWrapper `json:"memberType"`
	Status      intWrapper `json:"status"`
	Remark      string     `json:"remark"`
	AccountMode intWrapper `json:"accountMode"`
}

func (b *ByBit) GetSubAccounts(ctx context.Context) ([]*model.Account, error) {
	url := fmt.Sprintf("%s/v5/user/query-sub-members/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, "", timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result subMemberListResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getSubAccounts: error: 10002")
			return b.GetSubAccounts(ctx)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	accounts := make([]*model.Account, len(result.Result.SubMembers))

	for i := range result.Result.SubMembers {
		// 1 - normal, 2 -login banned, 4 - frozen
		isValid := int(result.Result.SubMembers[i].Status) == 1

		accounts[i] = &model.Account{
			ID:      int64(result.Result.SubMembers[i].UID),
			IsValid: isValid,
		}
	}
	return accounts, nil
}

type deleteSubAccountResponse struct {
	RetCode    returnCode  `json:"retCode"`
	RetMsg     string      `json:"retMsg"`
	Result     interface{} `json:"result"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

func (b *ByBit) DeleteSubAccount(ctx context.Context, id int) error {
	url := fmt.Sprintf("%s/v5/user/del-submember/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	params := map[string]string{
		"subMemberId": strconv.Itoa(id),
	}

	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, string(body), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetBody(body).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result deleteSubAccountResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "deleteSubAccount: error: 10002")
			return b.DeleteSubAccount(ctx, id)
		}
		return wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	return nil
}

type deleteSubAPIKeyResponse struct {
	RetCode    returnCode     `json:"retCode"`
	RetMsg     string         `json:"retMsg"`
	Result     map[string]any `json:"result"`
	RetExtInfo map[string]any `json:"retExtInfo"`
	Time       int64          `json:"time"`
}

func (b *ByBit) DeleteSubAPIKey(ctx context.Context, apiKey string) error {
	url := fmt.Sprintf("%s/v5/user/delete-sub-api/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	params := map[string]string{
		"apikey": apiKey,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, string(body), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetBody(body).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result deleteSubAPIKeyResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "deleteSubApiKey: error: 10002")
			return b.DeleteSubAPIKey(ctx, apiKey)
		}
		return wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	return nil
}

type depositAddressResponse struct {
	RetCode returnCode `json:"retCode"`
	RetMsg  string     `json:"retMsg"`
	Result  struct {
		Coin   string `json:"coin"`
		Chains struct {
			ChainType         string `json:"chainType"`
			AddressDeposit    string `json:"addressDeposit"`
			TagDeposit        string `json:"tagDeposit"`
			Chain             string `json:"chain"`
			BatchReleaseLimit string `json:"batchReleaseLimit"`
		} `json:"chains"`
	} `json:"result"`
	RetExtInfo map[string]interface{} `json:"retExtInfo"`
	Time       int64                  `json:"time"`
}

func (b *ByBit) GetSubDepositAddress(ctx context.Context, accountID int, coin, network string) (*model.AddressData, error) {
	url := fmt.Sprintf("%s/v5/asset/deposit/query-sub-member-address/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	type queryParams struct {
		Coin        string `url:"coin"`
		ChainType   string `url:"chainType"`
		SubMemberId int    `url:"subMemberId"`
	}

	params := queryParams{
		Coin:        coin,
		ChainType:   network,
		SubMemberId: accountID,
	}

	sortedParams, err := query.Values(params)
	if err != nil {
		return nil, fmt.Errorf("query.Values: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, sortedParams.Encode(), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetContext(ctx).
		Get(fmt.Sprintf("%s?%s", url, sortedParams.Encode()))

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result depositAddressResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getSubDepositAddress: error: 10002")
			return b.GetSubDepositAddress(ctx, accountID, coin, network)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	return &model.AddressData{
		Address: result.Result.Chains.AddressDeposit,
		Tag:     result.Result.Chains.TagDeposit,
	}, nil
}

type getSubBalanceByCoinResponse struct {
	RetCode returnCode `json:"retCode"`
	RetMsg  string     `json:"retMsg"`
	Result  struct {
		AccountType string     `json:"accountType"`
		BizType     intWrapper `json:"bizType"`
		AccountId   string     `json:"accountId"`
		MemberId    intWrapper `json:"memberId"`
		Balance     struct {
			Coin                  string         `json:"coin"`
			WalletBalance         decimalWrapper `json:"walletBalance"`
			TransferBalance       decimalWrapper `json:"transferBalance"`
			Bonus                 string         `json:"bonus"`
			TransferSafeAmount    decimalWrapper `json:"transferSafeAmount"`
			LtvTransferSafeAmount string         `json:"ltvTransferSafeAmount"`
		} `json:"balance"`
	} `json:"result"`
	RetExtInfo struct{} `json:"retExtInfo"`
	Time       int64    `json:"time"`
}

func (b *ByBit) GetBalanceByCoin(ctx context.Context, accountID int, coin, accountType string) (*model.Balance, error) {
	url := fmt.Sprintf("%s/v5/asset/transfer/query-account-coin-balance/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	type queryParams struct {
		MemberID               string `url:"memberId"`
		AccountType            string `url:"accountType"`
		Coin                   string `url:"coin"`
		WithTransferSafeAmount string `url:"withTransferSafeAmount"`
	}

	params := queryParams{
		MemberID:               strconv.Itoa(accountID),
		AccountType:            accountType,
		Coin:                   coin,
		WithTransferSafeAmount: "1",
	}

	sortedParams, err := query.Values(params)
	if err != nil {
		return nil, fmt.Errorf("query.Values: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, sortedParams.Encode(), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetContext(ctx).
		Get(fmt.Sprintf("%s?%s", url, sortedParams.Encode()))

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getSubBalanceByCoinResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getSubBalanceByCoin: error: 10002")
			return b.GetBalanceByCoin(ctx, accountID, coin, accountType)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	if int(result.Result.MemberId) != accountID {
		return nil, fmt.Errorf("received account id: %d, expected: %d", result.Result.MemberId, accountID)
	}

	return &model.Balance{
		WalletBalance:   decimal.Decimal(result.Result.Balance.WalletBalance),
		TransferBalance: decimal.Decimal(result.Result.Balance.TransferSafeAmount),
	}, nil
}

type getWithdrawalAmountResponse struct {
	RetCode returnCode                `json:"retCode"`
	RetMsg  string                    `json:"retMsg"`
	Result  getWithdrawalAmountResult `json:"result"`
	Time    int64                     `json:"time"`
}

type getWithdrawalAmountResult struct {
	LimitAmountUsd     string                                     `json:"limitAmountUsd"`
	WithdrawableAmount map[string]getWithdrawalAmountWithdrawable `json:"withdrawableAmount"`
}

type getWithdrawalAmountWithdrawable struct {
	Coin               string         `json:"coin"`
	WithdrawableAmount decimalWrapper `json:"withdrawableAmount"`
	AvailableBalance   string         `json:"availableBalance"`
}

func (b *ByBit) GetWithdrawalAmount(ctx context.Context, coin string) (decimal.Decimal, error) {
	url := fmt.Sprintf("%s/v5/asset/withdraw/withdrawable-amount/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	type queryParams struct {
		Coin string `url:"coin"`
	}

	params := queryParams{
		Coin: coin,
	}

	sortedParams, err := query.Values(params)
	if err != nil {
		return decimal.Zero, fmt.Errorf("query.Values: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, sortedParams.Encode(), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetContext(ctx).
		Get(fmt.Sprintf("%s?%s", url, sortedParams.Encode()))

	if err != nil {
		return decimal.Zero, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return decimal.Zero, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getWithdrawalAmountResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return decimal.Zero, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getWithdrawalAmount: error: 10002")
			return b.GetWithdrawalAmount(ctx, coin)
		}
		return decimal.Zero, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	res, ok := result.Result.WithdrawableAmount[model.Fund]
	if !ok {
		return decimal.Zero, fmt.Errorf("no fund")
	}

	return decimal.Decimal(res.WithdrawableAmount), nil
}

type withdrawResponse struct {
	RetCode    returnCode     `json:"retCode"`
	RetMsg     string         `json:"retMsg"`
	Result     withdrawResult `json:"result"`
	RetExtInfo map[string]any `json:"retExtInfo"`
	Time       int64          `json:"time"`
}

type withdrawResult struct {
	ID intWrapper `json:"id"`
}

func (b *ByBit) Withdraw(ctx context.Context, w *model.Withdraw) (int64, error) {
	url := fmt.Sprintf("%s/v5/asset/withdraw/create/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	var feeType string
	if w.IncludeFeeInAmount {
		feeType = "0"
	} else {
		feeType = "1"
	}

	params := map[string]string{
		"coin":        w.Coin,
		"chain":       w.Network,
		"address":     w.Address,
		"amount":      w.Amount.String(),
		"accountType": model.Fund,
		"timestamp":   strconv.FormatInt(timeStamp, 10),
		"feeType":     feeType,
		"requestId":   w.InternalID.String(), // idempotent key
	}

	if w.Tag != "" {
		params["tag"] = w.Tag
	}

	body, err := json.Marshal(params)
	if err != nil {
		return 0, fmt.Errorf("marshal: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, string(body), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetBody(body).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return 0, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return 0, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result withdrawResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return 0, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "withdraw: error: 10002")
			return b.Withdraw(ctx, w)
		}
		return 0, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	return int64(result.Result.ID), nil
}

type getWithdrawResponse struct {
	RetCode    returnCode        `json:"retCode"`
	RetMsg     string            `json:"retMsg"`
	Result     getWithdrawResult `json:"result"`
	RetExtInfo interface{}       `json:"retExtInfo"`
	Time       int64             `json:"time"`
}

type getWithdrawResult struct {
	Rows           []getWithdrawRow `json:"rows"`
	NextPageCursor string           `json:"nextPageCursor"`
}

type getWithdrawRow struct {
	Coin         string         `json:"coin"`
	Chain        string         `json:"chain"`
	Amount       decimalWrapper `json:"amount"`
	TxID         string         `json:"txID"`
	Status       string         `json:"status"`
	ToAddress    string         `json:"toAddress"`
	Tag          string         `json:"tag"`
	WithdrawFee  string         `json:"withdrawFee"`
	CreateTime   string         `json:"createTime"`
	UpdateTime   string         `json:"updateTime"`
	WithdrawID   intWrapper     `json:"withdrawId"`
	WithdrawType intWrapper     `json:"withdrawType"`
}

func (b *ByBit) GetWithdraw(ctx context.Context, withdrawID int) (*model.Withdraw, error) {
	url := fmt.Sprintf("%s/v5/asset/withdraw/query-record", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	type queryParams struct {
		WithdrawID string `url:"withdrawID"`
	}

	params := queryParams{
		WithdrawID: strconv.Itoa(withdrawID),
	}

	sortedParams, err := query.Values(params)
	if err != nil {
		return nil, fmt.Errorf("query.Values: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, sortedParams.Encode(), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetContext(ctx).
		Get(fmt.Sprintf("%s?%s", url, sortedParams.Encode()))

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getWithdrawResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getWithdraw: error: 10002")
			// не делаем повторную попытку, риск вывести монеты несколько раз
			//return b.GetWithdraw(ctx, withdrawID)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	for _, w := range result.Result.Rows {
		if int(w.WithdrawID) != withdrawID {
			continue
		}

		log.Debug(ctx, "withdraw status in stock: %s", w.Status)

		var status model.WithdrawStatus

		// https://bybit-exchange.github.io/docs/v5/enum#withdrawstatus
		switch w.Status {
		case "success":
			status = model.WsSuccess
		case "BlockchainConfirmed":
			status = model.WsBlockchainConfirmed
		case "Reject", "Fail":
			status = model.WsFailed
		default:
			status = model.WsUnknown
		}

		return &model.Withdraw{
			ID:      int64(withdrawID),
			HashID:  w.TxID,
			Coin:    w.Coin,
			Network: w.Chain,
			Address: w.ToAddress,
			Tag:     w.Tag,
			Amount:  decimal.Decimal(w.Amount),
			Status:  status,
		}, nil
	}

	return nil, fmt.Errorf("withdear not found")
}

type transferFromSubToMasterResponse struct {
	RetCode    returnCode                    `json:"retCode"`
	RetMsg     string                        `json:"retMsg"`
	Result     transferFromSubToMasterResult `json:"result"`
	RetExtInfo map[string]any                `json:"retExtInfo"`
	Time       int64                         `json:"time"`
}

type transferFromSubToMasterResult struct {
	TransferID string `json:"transferId"`
	Status     string `json:"status"`
}

// TransferFromSubToMaster return isExecuted, isSuccess, error
func (b *ByBit) TransferFromSubToMaster(ctx context.Context, transfer *model.InternalTransfer) (model.InternalTransferStatus, error) {
	url := fmt.Sprintf("%s/v5/asset/transfer/universal-transfer/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	params := map[string]string{
		"transferId":      transfer.ID.String(),
		"coin":            transfer.Coin,
		"amount":          transfer.Amount.String(),
		"fromMemberId":    strconv.Itoa(int(transfer.FromAccountID)),
		"toMemberId":      strconv.Itoa(b.masterUid),
		"fromAccountType": model.UNIFIED,
		"toAccountType":   model.Fund,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return model.ItsUnknown, fmt.Errorf("marshal: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, string(body), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetBody(body).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return 0, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return 0, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result transferFromSubToMasterResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return model.ItsUnknown, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "transferFromSubToMaster: error: 10002")
			return b.TransferFromSubToMaster(ctx, transfer)
		}
		return model.ItsUnknown, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	switch result.Result.Status {
	case "SUCCESS":
		return model.ItsSuccess, nil
	case "FAILED":
		return model.ItsFailed, nil
	case "PENDING":
		return model.ItsPending, nil
	case "STATUS_UNKNOWN":
		return model.ItsUnknown, fmt.Errorf("status_unknown: %s", result.Result.Status)
	default:
		return model.ItsUnknown, fmt.Errorf("unexpected response: transafer_id: %s, status: %s",
			result.Result.TransferID, result.Result.Status)
	}
}

func (b *ByBit) TransferFromMasterToSub(ctx context.Context, toAccountID int, coin string, amount decimal.Decimal) (string, error) {
	url := fmt.Sprintf("%s/v5/asset/transfer/universal-transfer/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	params := map[string]string{
		"transferId":      uuid.New().String(),
		"coin":            coin,
		"amount":          amount.String(),
		"fromMemberId":    strconv.Itoa(b.masterUid),
		"toMemberId":      strconv.Itoa(toAccountID),
		"fromAccountType": model.Fund,
		"toAccountType":   model.Spot,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, string(body), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetBody(body).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return "", fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result transferFromSubToMasterResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "transferFromMasterToSub: error: 10002")
			return b.TransferFromMasterToSub(ctx, toAccountID, coin, amount)

		}
		return "", wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	if result.Result.Status != "SUCCESS" {
		return "", fmt.Errorf("status: %s", result.Result.Status)
	}

	return result.Result.TransferID, nil
}

type createInternalTransferResponse struct {
	RetCode returnCode `json:"retCode"`
	RetMsg  string     `json:"retMsg"`
	Result  struct {
		TransferID string `json:"transferId"`
		Status     string `json:"status"`
	} `json:"result"`
	Time int64 `json:"time"`
}

func (b *ByBit) CreateInternalTransfer(ctx context.Context, trs *model.InternalTransfer, apiKey, apiSecret string) (*model.InternalTransfer, error) {
	url := fmt.Sprintf("%s/v5/asset/transfer/inter-transfer", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	params := map[string]string{
		"transferId":      trs.ID.String(),
		"coin":            trs.Coin,
		"fromAccountType": trs.FromAccountType,
		"toAccountType":   trs.ToAccountType,
		"amount":          trs.Amount.String(),
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     apiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(apiKey, apiSecret, string(body), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetBody(body).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result createInternalTransferResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "createInternalTransfer: error: 10002")
			return b.CreateInternalTransfer(ctx, trs, apiKey, apiSecret)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	switch result.Result.Status {
	case "SUCCESS":
		trs.Status = model.ItsSuccess
		return trs, nil
	case "FAILED":
		trs.Status = model.ItsFailed
		return trs, nil
	case "PENDING":
		trs.Status = model.ItsPending
		return trs, nil
	case "STATUS_UNKNOWN":
		trs.Status = model.ItsUnknown
		return trs, fmt.Errorf("status_unknown: %s", result.Result.Status)
	default:
		trs.Status = model.ItsUnknown
		return trs, fmt.Errorf("unexpected response: transafer_id: %s, status: %s",
			result.Result.TransferID, result.Result.Status)
	}
}

type getTransferResponse struct {
	RetCode returnCode        `json:"retCode"`
	RetMsg  string            `json:"retMsg"`
	Result  getTransferResult `json:"result"`
	Time    int64             `json:"time"`
}

type getTransferResult struct {
	List           []transferItem `json:"list"`
	NextPageCursor string         `json:"nextPageCursor"`
}

type transferItem struct {
	TransferId      string         `json:"transferId"`
	Coin            string         `json:"coin"`
	Amount          decimalWrapper `json:"amount"`
	Timestamp       intWrapper     `json:"timestamp"`
	Status          string         `json:"status"`
	FromAccountType string         `json:"fromAccountType"`
	ToAccountType   string         `json:"toAccountType"`
	FromMemberId    intWrapper     `json:"fromMemberId"`
	ToMemberId      intWrapper     `json:"toMemberId"`
}

func (b *ByBit) GetTransfer(ctx context.Context, transferID uuid.UUID) (*model.InternalTransfer, error) {
	url := fmt.Sprintf("%s/v5/asset/transfer/query-universal-transfer-list", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	type queryParams struct {
		TransferID string `url:"transferId"`
	}

	params := queryParams{
		TransferID: transferID.String(),
	}

	sortedParams, err := query.Values(params)
	if err != nil {
		return nil, fmt.Errorf("query.Values: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, sortedParams.Encode(), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetContext(ctx).
		Get(fmt.Sprintf("%s?%s", url, sortedParams.Encode()))

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getTransferResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getTransfer: error: 10002")
			return b.GetTransfer(ctx, transferID)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	for _, transfer := range result.Result.List {
		if transfer.TransferId != transferID.String() {
			continue
		}
		tm := time.Unix(int64(transfer.Timestamp), 0)

		var status model.InternalTransferStatus

		switch transfer.Status {
		case "SUCCESS":
			status = model.ItsSuccess
		case "FAILED":
			status = model.ItsFailed
		case "PENDING":
			status = model.ItsPending
		case "STATUS_UNKNOWN":
			status = model.ItsUnknown
		default:
			status = model.ItsUnknown
		}

		return &model.InternalTransfer{
			ID:              transferID,
			Coin:            transfer.Coin,
			Amount:          decimal.Decimal(transfer.Amount),
			FromAccountID:   int64(transfer.FromMemberId),
			ToAccountID:     int64(transfer.ToMemberId),
			FromAccountType: transfer.FromAccountType,
			ToAccountType:   transfer.ToAccountType,
			Status:          status,
			UpdatedAt:       tm,
		}, nil
	}

	return nil, fmt.Errorf("transfer not found")
}

type placeOrderResponse struct {
	RetCode    returnCode            `json:"retCode"`
	RetMsg     string                `json:"retMsg"`
	Result     placeOrderOrderResult `json:"result"`
	RetExtInfo interface{}           `json:"retExtInfo"`
	Time       int64                 `json:"time"`
}

type placeOrderOrderResult struct {
	OrderID     string `json:"orderId"`
	OrderLinkID string `json:"orderLinkId"`
}

func (b *ByBit) PlaceOrder(ctx context.Context, order *model.Order, secrets *model.Secrets) error {
	url := fmt.Sprintf("%s/v5/order/create/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	params := map[string]string{
		"orderLinkId": order.ID.String(),
		"category":    strings.ToLower(order.Category),
		"symbol":      order.Symbol,
		"side":        order.Direction.String(),
		"orderType":   "Market",
		"qty":         order.AmountFrom.String(),
	}

	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     secrets.ApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(secrets.ApiKey, secrets.ApiSecret, string(body), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetBody(body).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result placeOrderResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "placeOrder: error: 10002")
			return b.PlaceOrder(ctx, order, secrets)
		}
		return wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	return nil
}

type getInstrumentsInfoLotSizeFilter struct {
	BasePrecision  decimalWrapper `json:"basePrecision"`
	QuotePrecision decimalWrapper `json:"quotePrecision"`
	MinOrderQty    decimalWrapper `json:"minOrderQty"`
	MaxOrderQty    decimalWrapper `json:"maxOrderQty"`
	MinOrderAmt    decimalWrapper `json:"minOrderAmt"`
	MaxOrderAmt    decimalWrapper `json:"maxOrderAmt"`
}

type getInstrumentsInfoPriceFilter struct {
	TickSize string `json:"tickSize"`
}

type getInstrumentsInfoRiskParameters struct {
	LimitParameter  string `json:"limitParameter"`
	MarketParameter string `json:"marketParameter"`
}

type getInstrumentsInfoListItem struct {
	Symbol         string                           `json:"symbol"`
	BaseCoin       string                           `json:"baseCoin"`
	QuoteCoin      string                           `json:"quoteCoin"`
	Innovation     string                           `json:"innovation"`
	Status         string                           `json:"status"`
	MarginTrading  string                           `json:"marginTrading"`
	LotSizeFilter  getInstrumentsInfoLotSizeFilter  `json:"lotSizeFilter"`
	PriceFilter    getInstrumentsInfoPriceFilter    `json:"priceFilter"`
	RiskParameters getInstrumentsInfoRiskParameters `json:"riskParameters"`
}

type getInstrumentsInfoResult struct {
	Category string                       `json:"category"`
	List     []getInstrumentsInfoListItem `json:"list"`
}

type getInstrumentsInfoResponse struct {
	RetCode    returnCode               `json:"retCode"`
	RetMsg     string                   `json:"retMsg"`
	Result     getInstrumentsInfoResult `json:"result"`
	RetExtInfo interface{}              `json:"retExtInfo"`
	Time       int64                    `json:"time"`
}

func (b *ByBit) GetInstrumentsInfo(ctx context.Context) (map[string]*model.Symbol, error) {
	url := fmt.Sprintf("%s/v5/market/instruments-info", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	type queryParams struct {
		Category string `url:"category"`
	}

	params := queryParams{
		Category: strings.ToLower(model.Spot),
	}

	sortedParams, err := query.Values(params)
	if err != nil {
		return nil, fmt.Errorf("query.Values: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, sortedParams.Encode(), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetContext(ctx).
		Get(fmt.Sprintf("%s?%s", url, sortedParams.Encode()))

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getInstrumentsInfoResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getInstrumentsInfo: error: 10002")
			return b.GetInstrumentsInfo(ctx)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	if result.Result.Category != strings.ToLower(model.Spot) {
		return nil, fmt.Errorf("result.Result.Category: %s", result.Result.Category)
	}

	symbols := make(map[string]*model.Symbol, len(result.Result.List))

	for i := range result.Result.List {
		symbols[result.Result.List[i].Symbol] = &model.Symbol{
			Title:          result.Result.List[i].Symbol,
			BaseCoin:       result.Result.List[i].BaseCoin,
			QuoteCoin:      result.Result.List[i].QuoteCoin,
			BasePrecision:  decimal.Decimal(result.Result.List[i].LotSizeFilter.BasePrecision),
			QuotePrecision: decimal.Decimal(result.Result.List[i].LotSizeFilter.QuotePrecision),
			MinOrderQty:    decimal.Decimal(result.Result.List[i].LotSizeFilter.MinOrderQty),
			MaxOrderQty:    decimal.Decimal(result.Result.List[i].LotSizeFilter.MaxOrderQty),
			MinOrderAmt:    decimal.Decimal(result.Result.List[i].LotSizeFilter.MinOrderAmt),
			MaxOrderAmt:    decimal.Decimal(result.Result.List[i].LotSizeFilter.MaxOrderAmt),
		}
	}

	return symbols, nil
}

type getDepositRecordsResponse struct {
	RetCode    returnCode              `json:"retCode"`
	RetMsg     string                  `json:"retMsg"`
	Result     getDepositRecordsResult `json:"result"`
	RetExtInfo interface{}             `json:"retExtInfo"`
	Time       int64                   `json:"time"`
}

type getDepositRecordsResult struct {
	Rows           []getDepositRecordsRow `json:"rows"`
	NextPageCursor string                 `json:"nextPageCursor"`
}

type getDepositRecordsRow struct {
	Coin              string         `json:"coin"`
	Chain             string         `json:"chain"`
	Amount            decimalWrapper `json:"amount"`
	TxID              string         `json:"txID"`
	Status            intWrapper     `json:"status"`
	ToAddress         string         `json:"toAddress"`
	Tag               string         `json:"tag"`
	DepositFee        decimalWrapper `json:"depositFee"`
	SuccessAt         string         `json:"successAt"`
	Confirmations     string         `json:"confirmations"`
	TxIndex           string         `json:"txIndex"`
	BlockHash         string         `json:"blockHash"`
	BatchReleaseLimit string         `json:"batchReleaseLimit"`
	DepositType       string         `json:"depositType"`
}

func (b *ByBit) GetDepositRecords(ctx context.Context, coin string, startTime time.Time, subApikey, subApiSecret string) ([]*model.Deposit, error) {
	url := fmt.Sprintf("%s/v5/asset/deposit/query-record/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	type queryParams struct {
		Coin      string `url:"coin"`
		StartTime string `url:"startTime"`
		EndTime   string `url:"endTime"`
	}

	params := queryParams{
		Coin:      coin,
		StartTime: strconv.FormatInt(startTime.UTC().UnixMilli(), 10),
		EndTime:   strconv.FormatInt(time.Now().UTC().UnixMilli(), 10),
	}

	sortedParams, err := query.Values(params)
	if err != nil {
		return nil, fmt.Errorf("query.Values: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     subApikey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(subApikey, subApiSecret, sortedParams.Encode(), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetContext(ctx).
		Get(fmt.Sprintf("%s?%s", url, sortedParams.Encode()))

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getDepositRecordsResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getDepositRecords: error: 10002")
			return b.GetDepositRecords(ctx, coin, startTime, subApikey, subApiSecret)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	deposits := make([]*model.Deposit, 0)

	for _, row := range result.Result.Rows {
		var status model.DepositStatus
		switch int(row.Status) {
		case 3:
			status = model.DepositSuccess
		case 4:
			status = model.DepositFailed
		default:
			status = model.DepositPending
		}

		var depositType model.DepositType
		switch row.DepositType {
		case "0":
			depositType = model.DepositNormal
		default:
			depositType = model.DepositAbnormal
		}

		deposits = append(deposits, &model.Deposit{
			TxID:        row.TxID,
			Coin:        row.Coin,
			Amount:      decimal.Decimal(row.Amount),
			Fee:         decimal.Decimal(row.DepositFee),
			Status:      status,
			UpdatedAt:   time.Now().UTC(),
			DepositType: depositType,
		})
	}

	return deposits, nil
}

type getOrderHistoryResponse struct {
	RetCode    returnCode            `json:"retCode"`
	RetMsg     string                `json:"retMsg"`
	Result     getOrderHistoryResult `json:"result"`
	RetExtInfo interface{}           `json:"retExtInfo"`
	Time       int64                 `json:"time"`
}

type getOrderHistoryResult struct {
	List           []getOrderHistoryOrder `json:"list"`
	NextPageCursor string                 `json:"nextPageCursor"`
	Category       string                 `json:"category"`
}

type getOrderHistoryOrder struct {
	OrderID            string     `json:"orderId"`
	OrderLinkID        string     `json:"orderLinkId"`
	BlockTradeID       string     `json:"blockTradeId"`
	Symbol             string     `json:"symbol"`
	Price              string     `json:"price"`
	Qty                string     `json:"qty"`
	Side               string     `json:"side"`
	IsLeverage         string     `json:"isLeverage"`
	PositionIdx        intWrapper `json:"positionIdx"`
	OrderStatus        string     `json:"orderStatus"`
	CancelType         string     `json:"cancelType"`
	RejectReason       string     `json:"rejectReason"`
	AvgPrice           string     `json:"avgPrice"`
	LeavesQty          string     `json:"leavesQty"`
	LeavesValue        string     `json:"leavesValue"`
	CumExecQty         string     `json:"cumExecQty"`
	CumExecValue       string     `json:"cumExecValue"`
	CumExecFee         string     `json:"cumExecFee"`
	TimeInForce        string     `json:"timeInForce"`
	OrderType          string     `json:"orderType"`
	StopOrderType      string     `json:"stopOrderType"`
	OrderIv            string     `json:"orderIv"`
	TriggerPrice       string     `json:"triggerPrice"`
	TakeProfit         string     `json:"takeProfit"`
	StopLoss           string     `json:"stopLoss"`
	TpTriggerBy        string     `json:"tpTriggerBy"`
	SlTriggerBy        string     `json:"slTriggerBy"`
	TriggerDirection   intWrapper `json:"triggerDirection"`
	TriggerBy          string     `json:"triggerBy"`
	LastPriceOnCreated string     `json:"lastPriceOnCreated"`
	ReduceOnly         bool       `json:"reduceOnly"`
	CloseOnTrigger     bool       `json:"closeOnTrigger"`
	SmpType            string     `json:"smpType"`
	SmpGroup           intWrapper `json:"smpGroup"`
	SmpOrderID         string     `json:"smpOrderId"`
	TpslMode           string     `json:"tpslMode"`
	TpLimitPrice       string     `json:"tpLimitPrice"`
	SlLimitPrice       string     `json:"slLimitPrice"`
	PlaceType          string     `json:"placeType"`
	CreatedTime        string     `json:"createdTime"`
	UpdatedTime        string     `json:"updatedTime"`
}

func (b *ByBit) GetOrderStatus(ctx context.Context, orderID uuid.UUID, subApiKey, subApiSecret string) (model.OrderStatus, error) {
	url := fmt.Sprintf("%s/v5/order/history", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	type queryParams struct {
		Category    string `url:"category"`
		OrderLinkID string `json:"orderLinkId"`
	}

	params := queryParams{
		Category:    strings.ToLower(model.Spot),
		OrderLinkID: orderID.String(),
	}

	sortedParams, err := query.Values(params)
	if err != nil {
		return 0, fmt.Errorf("query.Values: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     subApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(subApiKey, subApiSecret, sortedParams.Encode(), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetContext(ctx).
		Get(fmt.Sprintf("%s?%s", url, sortedParams.Encode()))

	if err != nil {
		return 0, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return 0, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getOrderHistoryResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return 0, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getOrderStatus: error: 10002")
			return b.GetOrderStatus(ctx, orderID, subApiKey, subApiSecret)
		}
		return 0, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	var status model.OrderStatus
	for _, ord := range result.Result.List {
		if ord.OrderLinkID != orderID.String() {
			log.Error(ctx, "getOrderStatus: id unmatched: id: %s, exp: %s", ord.OrderLinkID, orderID.String())
			continue
		}
		switch ord.OrderStatus {
		case "Rejected", "Cancelled":
			status = model.OrderFailed
		case "Filled":
			status = model.OrderFilled
		case "PartiallyFilledCanceled":
			status = model.OrderPartiallyFilled
		default:
			status = model.OrderUnknown
		}
	}
	return status, nil
}

type requestAQuoteResponse struct {
	RetCode    returnCode          `json:"retCode"`
	RetMsg     string              `json:"retMsg"`
	Result     requestAQuoteResult `json:"result"`
	RetExtInfo struct{}            `json:"retExtInfo"`
	Time       int64               `json:"time"`
}

type requestAQuoteResult struct {
	QuoteTxId    string         `json:"quoteTxId"`
	ExchangeRate decimalWrapper `json:"exchangeRate"`
	FromCoin     string         `json:"fromCoin"`
	FromCoinType string         `json:"fromCoinType"`
	ToCoin       string         `json:"toCoin"`
	ToCoinType   string         `json:"toCoinType"`
	FromAmount   decimalWrapper `json:"fromAmount"`
	ToAmount     decimalWrapper `json:"toAmount"`
	ExpiredTime  intWrapper     `json:"expiredTime"`
}

func (b *ByBit) RequestAQuote(ctx context.Context, from, to, accountType string, amount decimal.Decimal) (*model.Quote, error) {
	url := fmt.Sprintf("%s/v5/asset/exchange/quote-apply/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	params := map[string]string{
		"fromCoin":      from,
		"toCoin":        to,
		"requestCoin":   from,
		"requestAmount": amount.String(),
		"accountType":   accountType,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, string(body), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetBody(body).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result requestAQuoteResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "requestAQuote: error: 10002")
			return b.RequestAQuote(ctx, from, to, accountType, amount)
		}
		return nil, wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	expiredTime := time.UnixMilli(int64(result.Result.ExpiredTime))

	return &model.Quote{
		ID:           result.Result.QuoteTxId,
		Rate:         decimal.Decimal(result.Result.ExchangeRate),
		FromCoin:     result.Result.FromCoin,
		FromCoinType: result.Result.FromCoinType,
		ToCoin:       result.Result.ToCoin,
		ToCoinType:   result.Result.ToCoinType,
		FromAmount:   decimal.Decimal(result.Result.FromAmount),
		ToAmount:     decimal.Decimal(result.Result.ToAmount),
		ExpiredTime:  expiredTime,
	}, nil
}

type confirmAQuoteResponse struct {
	RetCode    returnCode          `json:"retCode"`
	RetMsg     string              `json:"retMsg"`
	Result     confirmAQuoteResult `json:"result"`
	RetExtInfo interface{}         `json:"retExtInfo"`
	Time       int64               `json:"time"`
}

type confirmAQuoteResult struct {
	ExchangeStatus string `json:"exchangeStatus"`
}

func (b *ByBit) ConfirmAQuote(ctx context.Context, id string) (string, error) {
	url := fmt.Sprintf("%s/v5/asset/exchange/convert-execute/", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	params := map[string]string{
		"quoteTxId": id,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, string(body), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetBody(body).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return "", fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result confirmAQuoteResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "confirmAQuote: error: 10002")
			return b.ConfirmAQuote(ctx, id)
		}
		return "", wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	return result.Result.ExchangeStatus, nil
}

type getConvertStatusResponse struct {
	RetCode    returnCode                  `json:"retCode"`
	RetMsg     string                      `json:"retMsg"`
	Result     getConvertStatusOuterResult `json:"result"`
	RetExtInfo interface{}                 `json:"retExtInfo"`
	Time       int64                       `json:"time"`
}

type getConvertStatusOuterResult struct {
	Result getConvertStatusInnerResult `json:"result"`
}

type getConvertStatusInnerResult struct {
	ExchangeStatus string `json:"exchangeStatus"`
}

func (b *ByBit) GetConvertStatus(ctx context.Context, id, accountType string) (string, error) {
	url := fmt.Sprintf("%s/v5/asset/exchange/convert-result-query", b.apiUrl)
	timeStamp := time.Now().UTC().UnixNano() / 1000000

	type queryParams struct {
		QuoteTxId   string `url:"quoteTxId"`
		AccountType string `url:"accountType"`
	}

	params := queryParams{
		QuoteTxId:   id,
		AccountType: accountType,
	}

	sortedParams, err := query.Values(params)
	if err != nil {
		return "", fmt.Errorf("query.Values: %w", err)
	}

	headers := map[string]string{
		"X-BAPI-API-KEY":     b.masterApiKey,
		"X-BAPI-TIMESTAMP":   strconv.FormatInt(timeStamp, 10),
		"X-BAPI-RECV-WINDOW": recvWindow,
		"X-BAPI-SIGN":        b.getSignature(b.masterApiKey, b.masterApiSecret, sortedParams.Encode(), timeStamp),
	}

	resp, err := b.client.R().
		SetHeaders(headers).
		SetContext(ctx).
		Get(fmt.Sprintf("%s?%s", url, sortedParams.Encode()))

	if err != nil {
		return "", fmt.Errorf("doRequest: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result getConvertStatusResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	if result.RetCode != returnCodeSuccess {
		if result.RetCode == returnCodeInvalidTimestamp {
			log.Warn(ctx, "getOrderStatus: error: 10002")
			return b.GetConvertStatus(ctx, id, accountType)
		}
		return "", wrapError(ctx, resp, result.RetCode, result.RetMsg)
	}

	return result.Result.Result.ExchangeStatus, nil
}

func (b *ByBit) getSignature(apiKey, apiSecret, params string, timestamp int64) string {
	hmac256 := hmac.New(sha256.New, []byte(apiSecret))
	hmac256.Write([]byte(strconv.FormatInt(timestamp, 10) + apiKey + recvWindow + params))
	return hex.EncodeToString(hmac256.Sum(nil))
}

func generateRandomString(timestamp int64) (string, error) {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabetLen := len(alphabet)
	r := rand.New(rand.NewSource(timestamp))
	var randomString strings.Builder
	for i := 0; i < 16; i++ {
		char := alphabet[r.Intn(alphabetLen)]
		err := randomString.WriteByte(char)
		if err != nil {
			return "", fmt.Errorf("generateRandomString: %w", err)
		}
	}
	return randomString.String(), nil
}

func wrapError(ctx context.Context, resp *resty.Response, code returnCode, msg string) error {
	var traceID, url, body string

	if resp != nil {
		traceID = resp.Header().Get("traceId")

		url = resp.Request.URL

		switch v := resp.Request.Body.(type) {
		case string:
			body = v
		case []byte:
			body = string(v)
		default:
			body = fmt.Sprintf("unsupported body type: %T", v)
		}
	}

	log.SDebug(ctx, "byBit Error", map[string]any{
		"byBit_trace_id": traceID,
		"url":            url,
		"body":           body,
	})

	return fmt.Errorf("result.RetCode: %d, trace_id: %s, msg: %s", code, traceID, msg)
}
