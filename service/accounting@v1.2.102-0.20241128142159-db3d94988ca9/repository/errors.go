package repository

import sdkError "code.emcdtech.com/emcd/sdk/error"

var (
	ErrAcc1011 = sdkError.NewError("acc-1011", "failed mapping proto request")
	ErrAcc1012 = sdkError.NewError("acc-1012", "failed create user accounts")

	ErrAcc1021 = sdkError.NewError("acc-1021", "failed mapping proto request")
	ErrAcc1022 = sdkError.NewError("acc-1022", "failed find user account")

	ErrAcc1031 = sdkError.NewError("acc-1031", "failed mapping proto filter")
	ErrAcc1032 = sdkError.NewError("acc-1032", "failed find user accounts")

	ErrAcc1041 = sdkError.NewError("acc-1041", "failed find user account")

	ErrAcc1051 = sdkError.NewError("acc-1051", "failed parse constraint request")
	ErrAcc1052 = sdkError.NewError("acc-1052", "failed find user account")

	ErrAcc1061 = sdkError.NewError("acc-1061", "failed parse user uuid")
	ErrAcc1062 = sdkError.NewError("acc-1062", "failed find user accounts")
)
