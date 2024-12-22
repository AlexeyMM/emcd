package businessError

import sdkError "code.emcdtech.com/emcd/sdk/error"

var (
	SameCoinsErr                  = sdkError.NewError("swap-001", "the same coins error")
	NotSupportedWithdrawByNetwork = sdkError.NewError("swap-002", "withdraw of the selected network is not supported")
	MaxLimitExceededErr           = sdkError.NewError("swap-003", "max limit exceeded")
	BelowMinLimitErr              = sdkError.NewError("swap-004", "sum is below than min limit")
	CreateSubAccountErr           = sdkError.NewError("swap-005", "create sub account error")
	CalculateSwapErr              = sdkError.NewError("swap-006", "calculate swap error")
	CoinNotFoundErr               = sdkError.NewError("swap-007", "coin not found")
	NetworkNotFountErr            = sdkError.NewError("swap-008", "network not found")
	FeeNotFoundErr                = sdkError.NewError("swap-009", "fee not found")
	SymbolNotFoundErr             = sdkError.NewError("swap-010", "symbol not found")
	OrderBookNotFoundErr          = sdkError.NewError("swap-011", "order book not found")
	MarketDepthExceededErr        = sdkError.NewError("swap-012", "market depth exceeded")
	NoPathToSwapErr               = sdkError.NewError("swap-013", "there's no way to swap")
	TransactionNotFoundErr        = sdkError.NewError("swap-014", "transaction not found")
	SwapActiveLimitExceededErr    = sdkError.NewError("swap-015", "the limit on the number of swaps at the same time has been reached")
	SwapAlreadyExistsErr          = sdkError.NewError("swap-016", "swap already exists")
)

var (
	SymbolNotFoundAdminErr                          = sdkError.NewError("swap-admin-001", "symbol not found")
	WithdrawImpossibleBecauseOfSwapStatus           = sdkError.NewError("swap-admin-002", "withdraw impossible because because it doesn't match the swap id")
	WithdrawImpossibleBecauseWithdrawHasAlreadyBeen = sdkError.NewError("swap-admin-003", "withdraw impossible because has already been withdrawn")
)
