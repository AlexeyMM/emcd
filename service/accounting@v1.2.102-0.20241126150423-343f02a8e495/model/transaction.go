package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType int64

const (
	// IncomeBillTrTypeID Входящие поступления на пользовательские адреса
	IncomeBillTrTypeID TransactionType = 1
	// CnhldInitialTrTypeID Пользователь переводит получившийся баланс на один из продуктов — Coinhold или Hedge
	CnhldInitialTrTypeID TransactionType = 2
	// PoolIncomeFromOuterPoolsMainTrTypeID Пул получает приход денег от внешних пулов за майнинг основной монеты
	PoolIncomeFromOuterPoolsMainTrTypeID TransactionType = 3
	// PoolIncomeFromOuterPoolsMergeTrTypeID Пул получает приход денег от внешних пулов за майнинг мердж монет
	PoolIncomeFromOuterPoolsMergeTrTypeID TransactionType = 4
	// MainCoinMiningPayoutTrTypeID Пользователь получает начисление за майнинг основной монеты  на счет пользователя
	MainCoinMiningPayoutTrTypeID TransactionType = 5
	// MergeCoinMiningPayoutTrTypeID Пользователь получает начисление за майнинг мердж монет на счет пользователя
	MergeCoinMiningPayoutTrTypeID TransactionType = 6
	// UserPaysPoolComsTrTypeID Пользователь платит комиссию пулу
	UserPaysPoolComsTrTypeID TransactionType = 7
	// UserPaysPoolDonationsTrTypeID Пользователь платит donation пулу
	UserPaysPoolDonationsTrTypeID TransactionType = 8
	// PoolPaysUserCompensationTrTypeID Пул платит пользователю компенсацию или бонус (Pool)
	PoolPaysUserCompensationTrTypeID TransactionType = 9
	// PayUserReferralBonusTrTypeID Пользователю выплачивается его рефферальное вознаграждение
	PayUserReferralBonusTrTypeID TransactionType = 10
	// FromMiningReplenishmentTrTypeID Пользователь пополняет свой депозитный счет со счета для майнинга
	FromMiningReplenishmentTrTypeID TransactionType = 11
	// FromPlatformReplenishmentTrTypeID Пользователь пополняет свой депозитный счет заводя деньги через платформу
	FromPlatformReplenishmentTrTypeID TransactionType = 12
	// PayoutTrTypeID Пользователь получает начисление процентов
	PayoutTrTypeID TransactionType = 13
	// UserPercentCapTypeID Пользователь получает капитализацию процентов
	UserPercentCapTypeID TransactionType = 14
	// CnhldCloseTrTypeID Пользователь выводит депозит (депозит закрывается - сумма депозита и начисленные проценты переводятся на транзитный счет)
	CnhldCloseTrTypeID TransactionType = 15
	// CnhldWithdrawalTrTypeID Пул выплачивает депозит
	CnhldWithdrawalTrTypeID TransactionType = 16
	// PoolTransfersCoinsExternalTrTypeID Пул перечисляет монеты на депозит во внешнем источнике (Coinloan)
	PoolTransfersCoinsExternalTrTypeID TransactionType = 17
	// PoolGetsPercentsExternalCapTrTypeID Пул получает начисление процентов во внешнем источнике (Coinloan) - капитализирует
	PoolGetsPercentsExternalCapTrTypeID TransactionType = 18
	// PoolGetsPercentsExternalTrTypeID Пул получает начисление процентов во внешнем источнике (Coinloan)
	PoolGetsPercentsExternalTrTypeID TransactionType = 19
	// PoolGetsWithdrawalExternalTrTypeID Пул получает выплату процентов во внешнем источнике (Coinloan)
	PoolGetsWithdrawalExternalTrTypeID TransactionType = 20
	// PoolPaysUsersBalanceTrTypeID Пул выплачивает пользователю его баланс
	PoolPaysUsersBalanceTrTypeID TransactionType = 21
	// PoolPaysUsersReferralsTrTypeID Пул начисляет пользователю его реферальное вознаграждение
	PoolPaysUsersReferralsTrTypeID TransactionType = 22
	// PoolPaysBenefitOtherUserTrTypeID Пользователь переводит отчисление процента от дохода другому пользователю
	PoolPaysBenefitOtherUserTrTypeID TransactionType = 23
	// DBTransferAccType1TrTypeID Транзакция для корректировки данных во время переноса бд (accType 1)
	DBTransferAccType1TrTypeID TransactionType = 24
	// DBTransferAccType2TrTypeID Транзакция для корректировки данных во время переноса бд (accType 2)
	DBTransferAccType2TrTypeID TransactionType = 25
	// DBTransferAccType3TrTypeID Транзакция для корректировки данных во время переноса бд (accType 3)
	DBTransferAccType3TrTypeID TransactionType = 26
	// DBTransferAccType4TrTypeID Транзакция для корректировки данных во время переноса бд (accType 4)
	DBTransferAccType4TrTypeID TransactionType = 27
	// EarlyPaymentTrTypeID Досрочная выплата
	EarlyPaymentTrTypeID TransactionType = 28
	// TransferBetweenAccountsTrTypeID Перевод между аккаунтами
	TransferBetweenAccountsTrTypeID TransactionType = 29
	// WithdrawalWithCommissionTrTypeID Платный вывод средств с кошелька платформы
	WithdrawalWithCommissionTrTypeID TransactionType = 30
	// WalletMiningTransferTrTypeID Перевод средств с кошелька (платформы) на пул
	WalletMiningTransferTrTypeID TransactionType = 31
	// PayoutCommissionBillTrTypeID Комиссия за вывод с комиссией
	PayoutCommissionBillTrTypeID TransactionType = 32
	// AccrualsChineseTechAccountTrTypeID Начисления для китайского тех аккаунта
	AccrualsChineseTechAccountTrTypeID TransactionType = 33
	// CnhldDiffBalanceTrTypeID Перемещения средств между основным счётом и счётом коинхолда
	CnhldDiffBalanceTrTypeID TransactionType = 34
	// InterestAccrualTrTypeID Начисление процента интереса на счёт коинхолда
	InterestAccrualTrTypeID TransactionType = 35
	// CnhldEarlyCloseTrTypeID Перевод пользовательских средств при досрочном закрытии фиксированного коинхолда
	CnhldEarlyCloseTrTypeID TransactionType = 36
	// ReturnInterestsTrTypeID Возврат начисленного интереса на счёт компании при досрочном закрытии фиксированного коинхолда
	ReturnInterestsTrTypeID TransactionType = 37
	// InterestCapitalizationTrTypeID Возврат выплаченного интереса с основного счёта на счёт коинхолда при включенной капитализации
	InterestCapitalizationTrTypeID TransactionType = 38
	// BaseCoinFeesTokenWithdrawalsTrTypeID Траты на комиссии базовой монеты для пользовательских выводов токенов
	BaseCoinFeesTokenWithdrawalsTrTypeID TransactionType = 39
	// CommissionFeesTrTypeID Траты на комиссии за сбор и перевод на финансовый адрес токенов
	CommissionFeesTrTypeID TransactionType = 40
	// ToFinanceTransferBillTrTypeID Перевод на финансовый адрес токенов
	ToFinanceTransferBillTrTypeID TransactionType = 41
	// FiatWithdrawTrTypeID Блокировка средств при создании заявки на вывод в фиат
	FiatWithdrawTrTypeID TransactionType = 42
	// AdminUnblockWithdrawalTrTypeID Разблокировка средств в пользу админа при успешном выводе
	AdminUnblockWithdrawalTrTypeID TransactionType = 43
	// UnblockBackToUserFailWithdrawalTrTypeID Разблокировка средств обратно пользователю при неуспешном выводе
	UnblockBackToUserFailWithdrawalTrTypeID TransactionType = 44
	// HedgeBuyBlockTrTypeID Блокировка средств при покупке продукта hedge
	HedgeBuyBlockTrTypeID TransactionType = 45
	// UserUnblockTransferCommonWalletTrTypeID Разблокировка средств пользователя и перевод на общий кошелек
	UserUnblockTransferCommonWalletTrTypeID TransactionType = 46
	// FundsTransferHedgeAccTrTypeID Перевод средств на счет hedge
	FundsTransferHedgeAccTrTypeID TransactionType = 47
	// ComRemoveHedgeAccTrTypeID Снятие комиссии за продукт hedge в пользу хедж аккаунта компании
	ComRemoveHedgeAccTrTypeID TransactionType = 48
	// UserHedgeIncomeTrTypeID Начисление пользователю доходности хеджа с системного счёта хеджа
	UserHedgeIncomeTrTypeID TransactionType = 49
	// UserUnblockRefundHedgeTrTypeID Разблокировка и возврат средств пользователя на общий кошелек при неудачной покупке хеджа
	UserUnblockRefundHedgeTrTypeID TransactionType = 50
	// MiningWalletPayoutTrType Перевод средств с майнинга на основной счет
	MiningWalletPayoutTrType TransactionType = 51
	// PaymentReferralAccrualTrType Выплата реферальных начислений на кошелек
	PaymentReferralAccrualTrType TransactionType = 52
	// BalanceTransferToAdminTrType Перевод баланса неактивных аккаунтов на счет админа
	BalanceTransferToAdminTrType TransactionType = 53
	// CnhldRefBonusTrTypeID Выплата бонуса рефералу коинхолда на коинхолд или кошелек
	CnhldRefBonusTrTypeID TransactionType = 54
	// PromoAuthorBonusTrTypeID Выплата бонуса автору промокода коинхолда на кошелек
	PromoAuthorBonusTrTypeID TransactionType = 55
	// UserCompensationOrBonusTrTypeID Пул платит пользователю компенсацию или бонус (Wallets)
	UserCompensationOrBonusTrTypeID TransactionType = 56
	// ExchBlockTrTypeID Crypto Exchange блокировка перед обменом
	ExchBlockTrTypeID TransactionType = 57
	// ExchRollbackTrTypeID Crypto Exchange разблокировка при отмене обмена
	ExchRollbackTrTypeID TransactionType = 58
	// ExchUnblockTrTypeID Crypto Exchange разблокировка при удачном обмене
	ExchUnblockTrTypeID TransactionType = 59
	// ExchPayoutTrTypeID Crypto Exchange выплата пользователю обменяной монеты
	ExchPayoutTrTypeID TransactionType = 60
	// ToKucoinTradeTransferTrTypeID Перевод с расчётных адресов на trade аккаунт Kucoin
	ToKucoinTradeTransferTrTypeID TransactionType = 61
	// FromKucoinTradeTransferTrTypeID Перевод с trade аккаунта Kucoin на расчётные адреса
	FromKucoinTradeTransferTrTypeID TransactionType = 62
	// KucoinTransferFeeSpendTrTypeID Траты комиссии за переводы с/на trade аккаунт Kucoin
	KucoinTransferFeeSpendTrTypeID TransactionType = 63
	// WalletToWalletTrTypeID Перевод средств между пользователями с основного на основной
	WalletToWalletTrTypeID TransactionType = 64
	// MiningToWalletTrTypeID Перевод средств между пользователями с майнинга на основной
	MiningToWalletTrTypeID TransactionType = 65
	// P2PSellTrType P2P блокировка средств при создании заявки на продажу криптовалюты
	P2PSellTrType TransactionType = 66
	// P2PBuyTrType P2P разблокировка средств в пользу админа при успешной продаже криптовалюты
	P2PBuyTrType TransactionType = 67
	// P2PBuyFailTrType P2P разблокировка средств обратно пользователю при неуспешной продаже криптовалюты
	P2PBuyFailTrType TransactionType = 68
	// P2PSellCommissionTrType P2P блокировка комиссии при создании заявки на продажу криптовалюты
	P2PSellCommissionTrType TransactionType = 69
	// RejectedWithdrawalTrTypeID Возврат отклонённого вывода
	RejectedWithdrawalTrTypeID TransactionType = 70
	// RejectedWithdrawalFeeTrTypeID Возврат комиссии отклонённого вывода
	RejectedWithdrawalFeeTrTypeID TransactionType = 71
	// P2PUnblockToAdminTrTypeID P2P разблокировка комиссии в пользу админа при успешной продаже криптовалюты
	P2PUnblockToAdminTrTypeID TransactionType = 72
	// P2PUnblockFailToSellerTrTypeID P2P разблокировка комиссии обратно продавцу при неуспешной продаже криптовалюты
	P2PUnblockFailToSellerTrTypeID TransactionType = 73
	// TransferFromWalletToP2PAccount Перевод с кошелька на p2p-аккаунт
	TransferFromWalletToP2PAccount TransactionType = 74
	// TransferFromP2PAccountToWallet Перевод с p2p-аккаунта на кошелек
	TransferFromP2PAccountToWallet TransactionType = 75
	// CreditingFreeBonusToAccount Зачисление пользователю безвозмездного бонуса
	CreditingFreeBonusToAccount TransactionType = 76
	// CompensationForProblemsInProduct Зачисление пользователю компенсации за проблемы в продукте
	CompensationForProblemsInProduct TransactionType = 77
	// CorrectionErrorsInBilling Ручная корректировка балансов в связи с ошибками системы
	CorrectionErrorsInBilling TransactionType = 78
	// UnblockPayoutMiningToNode Выплата с майниг счета на ноду
	UnblockPayoutMiningToNode TransactionType = 79
	// RefundCanceledMiningPayouts Возврат выплаты с майниг счета
	RefundCanceledMiningPayouts TransactionType = 80
	// UnblockPayoutFreeWalletToNode Выплата бесплатного вывода с кошелька на ноду
	UnblockPayoutFreeWalletToNode TransactionType = 81
	// RefundCanceledFreePayoutFromWallet Возврат бесплатного вывода с кошелька
	RefundCanceledFreePayoutFromWallet TransactionType = 82
	// BlockTransferFromDepositToWallet Блокировка средств при выводе с депозита на кошелек
	BlockTransferFromDepositToWallet TransactionType = 83
	// UnblockTransferFromDepositToWallet Разблокировка средств при выводе с депозита на кошелек
	UnblockTransferFromDepositToWallet TransactionType = 84
	// IncomeForSharedMining Начисление за Shared Mining
	IncomeForSharedMining TransactionType = 85
	// FiatCardWithdraw Вывод фиата на карту
	FiatCardWithdraw TransactionType = 86
	// RefundFiatCardWithdrawFailed Возврат при неуспешном выводе фиата на карту
	RefundFiatCardWithdrawFailed TransactionType = 87
	// UserPaysWlComsTrTypeID Пользователь платит комиссию вайтлэйблу
	UserPaysWlComsTrTypeID TransactionType = 88
	//  WlPaysUserReferralsTrTypeID Вайтлэйбл выплачивает реферальную комиссию рефоводу
	WlPaysUserReferralsTrTypeID TransactionType = 89
	// OffchainWalletWithdraw Оффчейн кошельковый вывод
	OffchainWalletWithdraw TransactionType = 90
	// OffchainWalletReceipt Оффчейн кошельковое поступление
	OffchainWalletReceipt TransactionType = 91
	// OffchainMiningPayout Оффчейн майнерская выплата
	OffchainMiningPayout TransactionType = 92
	// P2PTransfer Перевод средств от одного пользователя к другому через p2p
	P2PTransfer TransactionType = 93
	// FiatProviderExternalAddrTransfer Вывод с расчетного счета на внешний адрес фиатного провайдера
	FiatProviderExternalAddrTransfer TransactionType = 94
	// P2PCryptoHold заморозка крипты при создании ордера (wallet -> p2p)
	P2PCryptoHold TransactionType = 95
	// P2PCryptoRelease разморозка крипты при отмене ордера (p2p -> wallet)
	P2PCryptoRelease TransactionType = 96
	// P2PTransferFee Комиссия в пользу emcd за успешную p2p-сделку
	P2PTransferFee TransactionType = 97
	// P2PTransferFee Комиссия в пользу рефовода за успешную p2p-сделку
	P2PTransferFeeRef TransactionType = 98

	// LastSupportedTransactionNumber Номер последней поддерживаемой транзакции, !!!обязательно указывать при добавлении новых констант!!!
	LastSupportedTransactionNumber = 98
)

func (t TransactionType) Int64() int64 {
	return int64(t)
}

// Transaction
// TODO::Сопоставить все поля из базы
type Transaction struct {
	ID                int64           `json:"id"`
	Type              TransactionType `json:"type"`
	CreatedAt         time.Time       `json:"created_at"`
	SenderAccountID   int64           `json:"sender_account_id"`
	ReceiverAccountID int64           `json:"receiver_account_id"`
	CoinID            int64           `json:"coin_id"`
	TokenID           int64           `json:"token_id"`
	Amount            decimal.Decimal `json:"amount"`
	Comment           string          `json:"comment"`
	Hash              string          `json:"hash"`
	ReceiverAddress   string          `json:"receiver_address"`
	Hashrate          int64           `json:"hashrate"`
	FromReferralId    int64           `json:"from_refferal_id"`
	ActionID          string          `json:"action_id"`
	UnblockAccountId  int64
	BlockedTill       time.Time
}

type TransactionSelectionWithBlock struct {
	SenderAccountID      int64
	ReceiverAccountID    int64
	CoinID               int64
	Type                 TransactionType
	Amount               decimal.Decimal
	BlockID              int64
	UnblockToAccountID   int64
	UnblockTransactionID int64
	ActionID             string
}

// TransactionCollectorFilter - TODO: refactor Transaction repo
type TransactionCollectorFilter struct {
	Types        []int32
	CoinId       *int32
	CreatedAtGt  *time.Time
	CreatedAtLte *time.Time
	Pagination   *Pagination
}
