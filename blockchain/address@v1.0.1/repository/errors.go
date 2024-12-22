package repository

import sdkError "code.emcdtech.com/emcd/sdk/error"

var (
	ErrAddr1011 = sdkError.NewError("adr-1011", "failed mapping proto request")
	ErrAddr1012 = sdkError.NewError("adr-1012", "failed determinate address type by network group by new way")
	ErrAddr1013 = sdkError.NewError("adr-1013", "failed get address by constraint by new way")
	ErrAddr1014 = sdkError.NewError("adr-1014", "failed create address by new way")

	ErrAddr1015 = sdkError.NewError("adr-1015", "failed determinate address type by network group by new way")
	ErrAddr1016 = sdkError.NewError("adr-1016", "failed get address by constraint by old way")
	ErrAddr1017 = sdkError.NewError("adr-1017", "failed create address by old way")

	ErrAddr1021 = sdkError.NewError("adr-1021", "failed mapping proto processing request")
	ErrAddr1022 = sdkError.NewError("adr-1022", "failed determinate address type by network group for processing address")
	ErrAddr1023 = sdkError.NewError("adr-1023", "failed create processing address")

	ErrAddr1031 = sdkError.NewError("adr-1031", "failed parse uuid")
	ErrAddr1032 = sdkError.NewError("adr-1032", "failed get new way addresses")
	ErrAddr1033 = sdkError.NewError("adr-1033", "failed get old way addresses")
	ErrAddr1034 = sdkError.NewError("adr-1034", "get multiple old addresses")
	ErrAddr1035 = sdkError.NewError("adr-1035", "get multiple new addresses")
	ErrAddr1036 = sdkError.NewError("adr-1036", "address not found by uuid")

	ErrAddr1041 = sdkError.NewError("adr-1041", "failed get new way addresses")
	ErrAddr1042 = sdkError.NewError("adr-1042", "failed get old way addresses")
	ErrAddr1043 = sdkError.NewError("adr-1043", "get multiple old addresses")
	ErrAddr1044 = sdkError.NewError("adr-1044", "get multiple new addresses")
	ErrAddr1045 = sdkError.NewError("adr-1045", "address not found by string")

	ErrAddr1051 = sdkError.NewError("adr-1051", "failed parse uuid")
	ErrAddr1052 = sdkError.NewError("adr-1052", "failed get new way addresses")
	ErrAddr1053 = sdkError.NewError("adr-1053", "failed get old way addresses")

	ErrAddr1061 = sdkError.NewError("adr-1061", "failed parse old filter request")
	ErrAddr1062 = sdkError.NewError("adr-1062", "failed get old addresses")

	ErrAddr1071 = sdkError.NewError("adr-1071", "failed parse new filter request")
	ErrAddr1072 = sdkError.NewError("adr-1072", "failed get new addresses")

	ErrAddr1081 = sdkError.NewError("adr-1081", "failed parse personal filter request")
	ErrAddr1082 = sdkError.NewError("adr-1082", "failed get personal address")
	ErrAddr1083 = sdkError.NewError("adr-1083", "minimal payout is required for new personal address")
	ErrAddr1084 = sdkError.NewError("adr-1084", "failed get default minimal payout")
	ErrAddr1085 = sdkError.NewError("adr-1085", "minimal payout is less then default minimal payout")
	ErrAddr1086 = sdkError.NewError("adr-1086", "failed create personal address")
	ErrAddr1087 = sdkError.NewError("adr-1087", "failed update personal address")

	ErrAddr1091 = sdkError.NewError("adr-1091", "failed parse personal delete request")
	ErrAddr1092 = sdkError.NewError("adr-1092", "failed delete personal address")

	ErrAddr1101 = sdkError.NewError("adr-1101", "failed parse personal filter request")
	ErrAddr1102 = sdkError.NewError("adr-1102", "failed get personal addresses")

	ErrAddr1201 = sdkError.NewError("adr-1201", "failed parse personal filter request")
	ErrAddr1202 = sdkError.NewError("adr-1202", "failed get personal addresses")
	ErrAddr1203 = sdkError.NewError("adr-1203", "get multiple personal addresses")

	ErrAddr1301 = sdkError.NewError("adr-1301", "failed parse dirty form request")
	ErrAddr1302 = sdkError.NewError("adr-1302", "failed create or update dirty address")

	ErrAddr1401 = sdkError.NewError("adr-1401", "failed parse dirty filter request")
	ErrAddr1402 = sdkError.NewError("adr-1402", "failed get dirty addresses")
)
