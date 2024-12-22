package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/b2b/swap/internal/client"
	"code.emcdtech.com/b2b/swap/model"
)

type ByBit struct {
	cliAccount     client.ExchangeAccount
	cliTransaction client.ExchangeTransaction
	cliMarket      client.Market
	histogram      *prometheus.HistogramVec
}

func NewByBit(
	cliAccount client.ExchangeAccount,
	cliTransaction client.ExchangeTransaction,
	cliMarket client.Market,
	histogram *prometheus.HistogramVec,
) *ByBit {
	return &ByBit{
		cliAccount:     cliAccount,
		cliTransaction: cliTransaction,
		cliMarket:      cliMarket,
		histogram:      histogram,
	}
}

func (m *ByBit) CreateSubAccount(ctx context.Context) (*model.Account, error) {
	start := time.Now()
	acc, err := m.cliAccount.CreateSubAccount(ctx)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.CreateSubAccount", strconv.FormatBool(err == nil)).Observe(duration)
	return acc, err
}

func (m *ByBit) CreateSubAPIKey(ctx context.Context, accountID int) (*model.Secrets, error) {
	start := time.Now()
	sec, err := m.cliAccount.CreateSubAPIKey(ctx, accountID)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.CreateSubAPIKey", strconv.FormatBool(err == nil)).Observe(duration)
	return sec, err
}

func (m *ByBit) DeleteSubAccount(ctx context.Context, id int) error {
	start := time.Now()
	err := m.cliAccount.DeleteSubAccount(ctx, id)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.DeleteSubAccount", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (m *ByBit) DeleteSubAPIKey(ctx context.Context, apiKey string) error {
	start := time.Now()
	err := m.cliAccount.DeleteSubAPIKey(ctx, apiKey)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.DeleteSubAPIKey", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (m *ByBit) GetSubDepositAddress(ctx context.Context, accountID int, coin, network string) (*model.AddressData, error) {
	start := time.Now()
	addr, err := m.cliAccount.GetSubDepositAddress(ctx, accountID, coin, network)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetSubDepositAddress", strconv.FormatBool(err == nil)).Observe(duration)
	return addr, err
}

func (m *ByBit) GetBalanceByCoin(ctx context.Context, accountID int, coin, accountType string) (*model.Balance, error) {
	start := time.Now()
	balance, err := m.cliAccount.GetBalanceByCoin(ctx, accountID, coin, accountType)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetBalanceByCoin", strconv.FormatBool(err == nil)).Observe(duration)
	return balance, err
}

func (m *ByBit) GetDepositRecords(ctx context.Context, coin string, startTime time.Time, subAccApikey, subAccApiSecret string) ([]*model.Deposit, error) {
	start := time.Now()
	rec, err := m.cliAccount.GetDepositRecords(ctx, coin, startTime, subAccApikey, subAccApiSecret)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetDepositRecords", strconv.FormatBool(err == nil)).Observe(duration)
	return rec, err
}

func (m *ByBit) GetWithdrawalAmount(ctx context.Context, coin string) (decimal.Decimal, error) {
	start := time.Now()
	amount, err := m.cliAccount.GetWithdrawalAmount(ctx, coin)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetWithdrawalAmount", strconv.FormatBool(err == nil)).Observe(duration)
	return amount, err
}

func (m *ByBit) GetSubAccounts(ctx context.Context) ([]*model.Account, error) {
	start := time.Now()
	accs, err := m.cliAccount.GetSubAccounts(ctx)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetSubAccounts", strconv.FormatBool(err == nil)).Observe(duration)
	return accs, err
}

func (m *ByBit) TransferFromSubToMaster(ctx context.Context, transfer *model.InternalTransfer) (model.InternalTransferStatus, error) {
	start := time.Now()
	status, err := m.cliTransaction.TransferFromSubToMaster(ctx, transfer)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.TransferFromSubToMaster", strconv.FormatBool(err == nil)).Observe(duration)
	return status, err
}

func (m *ByBit) CreateInternalTransfer(ctx context.Context, trs *model.InternalTransfer, apiKey, apiSecret string) (*model.InternalTransfer, error) {
	start := time.Now()
	status, err := m.cliTransaction.CreateInternalTransfer(ctx, trs, apiKey, apiSecret)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.CreateInternalTransfer", strconv.FormatBool(err == nil)).Observe(duration)
	return status, err
}

func (m *ByBit) GetTransfer(ctx context.Context, transferID uuid.UUID) (*model.InternalTransfer, error) {
	start := time.Now()
	tr, err := m.cliTransaction.GetTransfer(ctx, transferID)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetTransfer", strconv.FormatBool(err == nil)).Observe(duration)
	return tr, err
}

func (m *ByBit) Withdraw(ctx context.Context, withdraw *model.Withdraw) (int64, error) {
	start := time.Now()
	status, err := m.cliTransaction.Withdraw(ctx, withdraw)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.Withdraw", strconv.FormatBool(err == nil)).Observe(duration)
	return status, err
}

func (m *ByBit) GetWithdraw(ctx context.Context, withdrawID int) (*model.Withdraw, error) {
	start := time.Now()
	wt, err := m.cliTransaction.GetWithdraw(ctx, withdrawID)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetWithdraw", strconv.FormatBool(err == nil)).Observe(duration)
	return wt, err
}

func (m *ByBit) PlaceOrder(ctx context.Context, order *model.Order, secrets *model.Secrets) error {
	start := time.Now()
	err := m.cliMarket.PlaceOrder(ctx, order, secrets)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.PlaceOrder", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (m *ByBit) GetOrderStatus(ctx context.Context, orderID uuid.UUID, subApiKey, subApiSecret string) (model.OrderStatus, error) {
	start := time.Now()
	status, err := m.cliMarket.GetOrderStatus(ctx, orderID, subApiKey, subApiSecret)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetOrderStatus", strconv.FormatBool(err == nil)).Observe(duration)
	return status, err
}

func (m *ByBit) GetAllFeeRate(ctx context.Context) (map[string]*model.Fee, error) {
	start := time.Now()
	rates, err := m.cliMarket.GetAllFeeRate(ctx)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetAllFeeRate", strconv.FormatBool(err == nil)).Observe(duration)
	return rates, err
}

func (m *ByBit) GetCoinInfo(ctx context.Context) ([]*model.Coin, error) {
	start := time.Now()
	info, err := m.cliMarket.GetCoinInfo(ctx)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetCoinInfo", strconv.FormatBool(err == nil)).Observe(duration)
	return info, err
}

func (m *ByBit) GetConvertCoinList(ctx context.Context) (map[string]int, error) {
	start := time.Now()
	info, err := m.cliMarket.GetConvertCoinList(ctx)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetConvertCoinList", strconv.FormatBool(err == nil)).Observe(duration)
	return info, err
}

func (m *ByBit) GetInstrumentsInfo(ctx context.Context) (map[string]*model.Symbol, error) {
	start := time.Now()
	info, err := m.cliMarket.GetInstrumentsInfo(ctx)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetInstrumentsInfo", strconv.FormatBool(err == nil)).Observe(duration)
	return info, err
}

func (m *ByBit) RequestAQuote(ctx context.Context, from, to, accountType string, amount decimal.Decimal) (*model.Quote, error) {
	start := time.Now()
	quote, err := m.cliMarket.RequestAQuote(ctx, from, to, accountType, amount)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.RequestAQuote", strconv.FormatBool(err == nil)).Observe(duration)
	return quote, err
}

func (m *ByBit) ConfirmAQuote(ctx context.Context, id string) (string, error) {
	start := time.Now()
	quote, err := m.cliMarket.ConfirmAQuote(ctx, id)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.ConfirmAQuote", strconv.FormatBool(err == nil)).Observe(duration)
	return quote, err
}

func (m *ByBit) GetConvertStatus(ctx context.Context, id, accountType string) (string, error) {
	start := time.Now()
	status, err := m.cliMarket.GetConvertStatus(ctx, id, accountType)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("exchangeAccount.GetConvertStatus", strconv.FormatBool(err == nil)).Observe(duration)
	return status, err
}
