package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Status uint8

// При добавлении нового статуса, а так же при изменении названия текущего, необходимо обновить БД
// swap.swap_statuses
const (
	Unknown Status = iota

	WaitDeposit
	CheckDeposit
	DepositError // Deprecated
	TransferToUnified
	CreateOrder
	PlaceOrder
	CheckOrder
	PlaceAdditionalOrder
	CheckAdditionalOrder
	TransferFromSubToMaster
	CheckTransferFromSubToMaster
	PrepareWithdraw
	WithdrawSwapStatus
	WaitWithdraw

	Completed
	Cancel
	Error
	ManualCompleted
)

func (s Status) String() string {
	switch s {
	case Unknown:
		return "Unknown"
	case WaitDeposit:
		return "WaitDeposit"
	case CheckDeposit:
		return "CheckDeposit"
	case DepositError:
		return "DepositError"
	case TransferToUnified:
		return "TransferToUnified"
	case CreateOrder:
		return "CreateOrder"
	case PlaceOrder:
		return "PlaceOrder"
	case CheckOrder:
		return "CheckOrder"
	case PlaceAdditionalOrder:
		return "PlaceAdditionalOrder"
	case CheckAdditionalOrder:
		return "CheckAdditionalOrder"
	case TransferFromSubToMaster:
		return "TransferFromSubToMaster"
	case CheckTransferFromSubToMaster:
		return "CheckTransferFromSubToMaster"
	case PrepareWithdraw:
		return "PrepareWithdraw"
	case WithdrawSwapStatus:
		return "WithdrawSwapStatus"
	case WaitWithdraw:
		return "WaitWithdraw"
	case Completed:
		return "Completed"
	case Cancel:
		return "Cancel"
	case Error:
		return "Error"
	case ManualCompleted:
		return "ManualCompleted"
	default:
		return "Unknown"
	}
}

// PublicStatus статус, который видит пользователь
type PublicStatus uint8

const (
	PSUnknown      PublicStatus = iota
	PSWaitDeposit               // Ожидаем депозит
	PSCheckDeposit              // Подтверждаем депозит
	PSSwap                      // Обмениваем
	PSWithdraw                  // Выводим
	PSCompleted                 // Вывели
	PSError                     // Ошибка
	PSCancel                    // Swap отменён
)

func ConvertInternalToPublicStatus(internalStatus Status) PublicStatus {
	switch internalStatus {
	case WaitDeposit:
		return PSWaitDeposit
	case CheckDeposit:
		return PSCheckDeposit
	case TransferToUnified, CreateOrder, PlaceOrder, CheckOrder, PlaceAdditionalOrder, CheckAdditionalOrder,
		TransferFromSubToMaster, CheckTransferFromSubToMaster, PrepareWithdraw:
		return PSSwap
	case WithdrawSwapStatus, WaitWithdraw:
		return PSWithdraw
	case Completed:
		return PSCompleted
	case Cancel:
		return PSCancel
	case DepositError, Error:
		return PSError
	default:
		return PSUnknown
	}
}

type Swap struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	AccountFrom int64
	CoinFrom    string
	AddressFrom string
	TagFrom     string
	NetworkFrom string
	CoinTo      string
	AddressTo   string
	NetworkTo   string
	TagTo       string
	AmountFrom  decimal.Decimal
	AmountTo    decimal.Decimal
	Status      Status
	StartTime   time.Time
	EndTime     time.Time
	PartnerID   string
}

type SwapFilter struct {
	ID            *uuid.UUID
	NotEqStatus   []Status
	TxID          *string
	UserID        *uuid.UUID
	Email         *string
	AddressFrom   *string
	StartTimeFrom *time.Time
	StartTimeTo   *time.Time

	Offset *int
	Limit  *int
}

type SwapPartial struct {
	UserID     *uuid.UUID
	Status     *Status
	AmountFrom *decimal.Decimal
	AmountTo   *decimal.Decimal
	StartTime  *time.Time
	EndTime    *time.Time
	AddressTo  *string
	TagTo      *string
}

func (s *Swap) Update(partial *SwapPartial) {
	if partial.UserID != nil {
		s.UserID = *partial.UserID
	}
	if partial.Status != nil {
		s.Status = *partial.Status
	}
	if partial.AmountFrom != nil {
		s.AmountFrom = *partial.AmountFrom
	}
	if partial.AmountTo != nil {
		s.AmountTo = *partial.AmountTo
	}
	if partial.StartTime != nil {
		s.StartTime = *partial.StartTime
	}
	if partial.EndTime != nil {
		s.EndTime = *partial.EndTime
	}
	if partial.AddressTo != nil {
		s.AddressTo = *partial.AddressTo
	}
	if partial.TagTo != nil {
		s.TagTo = *partial.TagTo
	}
}

type Swaps []*Swap

type Estimate struct {
	AmountFrom decimal.Decimal
	AmountTo   decimal.Decimal
	Rate       decimal.Decimal
	Limits     *Limits
}

type Limits struct {
	Min decimal.Decimal
	Max decimal.Decimal
}

type EstimateRequest struct {
	CoinFrom    string
	CoinTo      string
	NetworkFrom string
	NetworkTo   string
	AmountFrom  decimal.Decimal
	AmountTo    decimal.Decimal
}

type SwapRequest struct {
	CoinFrom    string
	CoinTo      string
	NetworkFrom string
	NetworkTo   string
	AmountFrom  decimal.Decimal
	AmountTo    decimal.Decimal
	AddressTo   *AddressData
	PartnerID   string
}

type SwapByIDResponse struct {
	CoinFrom     string
	CoinTo       string
	NetworkFrom  string
	NetworkTo    string
	Rate         decimal.Decimal
	AmountFrom   decimal.Decimal
	AmountTo     decimal.Decimal
	AddressTo    *AddressData
	AddressFrom  *AddressData
	StartTime    time.Time
	SwapDuration time.Duration
	Status       Status
}

type Balance struct {
	WalletBalance   decimal.Decimal
	TransferBalance decimal.Decimal
}

type SwapStatusHistoryItem struct {
	Status Status
	SetAt  time.Time
}

type SwapStatusHistoryFilter struct {
	SwapID *uuid.UUID
}
