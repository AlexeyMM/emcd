package businessError

import sdkError "code.emcdtech.com/emcd/sdk/error"

var (
	Internal       = sdkError.NewError("0001", "internal Error")
	InvalidRequest = sdkError.NewError("0002", "invalid request")

	// Auth
	InvalidApiKey    = sdkError.NewError("0003", "invalid api key")
	InvalidTimestamp = sdkError.NewError("0004", "invalid timestamp")
	InvalidSignature = sdkError.NewError("0005", "invalid signature")
	DoubleRequest    = sdkError.NewError("0006", "request previously executed")
	InvalidIP        = sdkError.NewError("0007", "invalid ip")
)
