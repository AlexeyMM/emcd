package client

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/b2b/swap/model"
)

type ExchangeAccount interface {
	CreateSubAccount(ctx context.Context) (*model.Account, error)
	CreateSubAPIKey(ctx context.Context, accountID int) (*model.Secrets, error)
	DeleteSubAccount(ctx context.Context, id int) error
	DeleteSubAPIKey(ctx context.Context, apiKey string) error
	GetSubDepositAddress(ctx context.Context, accountID int, coin, network string) (*model.AddressData, error)
	GetBalanceByCoin(ctx context.Context, accountID int, coin, accountType string) (*model.Balance, error)
	GetDepositRecords(ctx context.Context, coin string, startTime time.Time, subAccApikey, subAccApiSecret string) ([]*model.Deposit, error)
	GetWithdrawalAmount(ctx context.Context, coin string) (decimal.Decimal, error)
	GetSubAccounts(ctx context.Context) ([]*model.Account, error)
}

type ExchangeTransaction interface {
	TransferFromSubToMaster(ctx context.Context, transfer *model.InternalTransfer) (model.InternalTransferStatus, error)
	CreateInternalTransfer(ctx context.Context, trs *model.InternalTransfer, apiKey, apiSecret string) (*model.InternalTransfer, error)
	GetTransfer(ctx context.Context, transferID uuid.UUID) (*model.InternalTransfer, error)
	Withdraw(ctx context.Context, withdraw *model.Withdraw) (int64, error)
	GetWithdraw(ctx context.Context, withdrawID int) (*model.Withdraw, error)
}

//go:generate mockery --name=Market
type Market interface {
	PlaceOrder(ctx context.Context, order *model.Order, secrets *model.Secrets) error
	GetOrderStatus(ctx context.Context, orderID uuid.UUID, subApiKey, subApiSecret string) (model.OrderStatus, error)
	GetAllFeeRate(ctx context.Context) (map[string]*model.Fee, error)
	GetCoinInfo(ctx context.Context) ([]*model.Coin, error)
	GetConvertCoinList(ctx context.Context) (map[string]int, error)
	GetInstrumentsInfo(ctx context.Context) (map[string]*model.Symbol, error)
	RequestAQuote(ctx context.Context, from, to, accountType string, amount decimal.Decimal) (*model.Quote, error)
	ConfirmAQuote(ctx context.Context, id string) (string, error)
	GetConvertStatus(ctx context.Context, id, accountType string) (string, error)
}

type Subscriber interface {
	SubscribeOnWallet(ctx context.Context, swapID uuid.UUID, swapCreated time.Time, account *model.Account, coin string, amount decimal.Decimal) error
	SubscribeOnOrders(ctx context.Context, account *model.Account, orders []*model.Order, receivedFirstOrder bool) error
	SubscribeOnOrderbooks(ctx context.Context, symbols []*model.Symbol) error
}
