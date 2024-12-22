package model

import (
	"errors"
	"fmt"
)

type ErrorCode string

var (
	ErrorCodeNoSuchMerchant           ErrorCode = "NO_SUCH_MERCHANT"
	ErrorCodeNoAvailableAddress       ErrorCode = "NO_AVAILABLE_ADDRESS"
	ErrorCodeInternal                 ErrorCode = "INTERNAL_ERROR"
	ErrorCodeNoSuchInvoiceForm        ErrorCode = "NO_SUCH_INVOICE_FORM"
	ErrorCodeInvalidArgument          ErrorCode = "INVALID_ARGUMENT"
	ErrorCodeNoSuchCoin               ErrorCode = "NO_SUCH_COIN"
	ErrorCodeNoSuchNetwork            ErrorCode = "NO_SUCH_NETWORK"
	ErrorCodeNoSuchInvoice            ErrorCode = "NO_SUCH_INVOICE"
	ErrorCodeTransactionAlreadyExists ErrorCode = "TRANSACTION_ALREADY_EXISTS"
	ErrorCodeInvoiceFormExpired       ErrorCode = "INVOICE_FORM_EXPIRED"
)

type Error struct {
	Inner   error
	Code    ErrorCode
	Message string
}

func (e *Error) Error() string {
	message := string(e.Code)

	if e.Message != "" {
		message = fmt.Sprintf("%s: %s", message, e.Message)
	}

	if e.Inner != nil {
		message = fmt.Sprintf("%s: %s", message, e.Inner)
	}

	return message
}

func (e *Error) Is(target error) bool {
	var targetErr *Error
	if !errors.As(target, &targetErr) {
		return errors.Is(e.Inner, target)
	}

	if e.Code != targetErr.Code {
		return false
	}

	if e.Inner != nil && targetErr.Inner != nil {
		return errors.Is(e.Inner, targetErr.Inner)
	}

	return true
}

var codeToMessage = map[ErrorCode]string{
	ErrorCodeInternal: "We encountered an internal error. Please try again.",
}

func ErrorWithDefaultMessage(err *Error) *Error {
	err.Message = codeToMessage[err.Code]

	return err
}

func ErrCode(err error) ErrorCode {
	var targetErr *Error
	if errors.As(err, &targetErr) {
		return targetErr.Code
	}

	return ""
}
