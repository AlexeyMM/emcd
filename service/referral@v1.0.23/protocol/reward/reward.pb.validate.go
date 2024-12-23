// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: protocol/reward/reward.proto

package reward

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Transaction with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Transaction) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Transaction with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in TransactionMultiError, or
// nil if none found.
func (m *Transaction) ValidateAll() error {
	return m.validate(true)
}

func (m *Transaction) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	// no validation rules for Type

	// no validation rules for Amount

	if len(errors) > 0 {
		return TransactionMultiError(errors)
	}

	return nil
}

// TransactionMultiError is an error wrapping multiple validation errors
// returned by Transaction.ValidateAll() if the designated constraints aren't met.
type TransactionMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TransactionMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TransactionMultiError) AllErrors() []error { return m }

// TransactionValidationError is the validation error returned by
// Transaction.Validate if the designated constraints aren't met.
type TransactionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TransactionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TransactionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TransactionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TransactionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TransactionValidationError) ErrorName() string { return "TransactionValidationError" }

// Error satisfies the builtin error interface
func (e TransactionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTransaction.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TransactionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TransactionValidationError{}

// Validate checks the field values on CalculateRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *CalculateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CalculateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CalculateRequestMultiError, or nil if none found.
func (m *CalculateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CalculateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	// no validation rules for Product

	// no validation rules for Coin

	// no validation rules for Amount

	if len(errors) > 0 {
		return CalculateRequestMultiError(errors)
	}

	return nil
}

// CalculateRequestMultiError is an error wrapping multiple validation errors
// returned by CalculateRequest.ValidateAll() if the designated constraints
// aren't met.
type CalculateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CalculateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CalculateRequestMultiError) AllErrors() []error { return m }

// CalculateRequestValidationError is the validation error returned by
// CalculateRequest.Validate if the designated constraints aren't met.
type CalculateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CalculateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CalculateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CalculateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CalculateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CalculateRequestValidationError) ErrorName() string { return "CalculateRequestValidationError" }

// Error satisfies the builtin error interface
func (e CalculateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCalculateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CalculateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CalculateRequestValidationError{}

// Validate checks the field values on CalculateResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *CalculateResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CalculateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CalculateResponseMultiError, or nil if none found.
func (m *CalculateResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CalculateResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetTxs() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CalculateResponseValidationError{
						field:  fmt.Sprintf("Txs[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CalculateResponseValidationError{
						field:  fmt.Sprintf("Txs[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CalculateResponseValidationError{
					field:  fmt.Sprintf("Txs[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return CalculateResponseMultiError(errors)
	}

	return nil
}

// CalculateResponseMultiError is an error wrapping multiple validation errors
// returned by CalculateResponse.ValidateAll() if the designated constraints
// aren't met.
type CalculateResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CalculateResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CalculateResponseMultiError) AllErrors() []error { return m }

// CalculateResponseValidationError is the validation error returned by
// CalculateResponse.Validate if the designated constraints aren't met.
type CalculateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CalculateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CalculateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CalculateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CalculateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CalculateResponseValidationError) ErrorName() string {
	return "CalculateResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CalculateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCalculateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CalculateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CalculateResponseValidationError{}

// Validate checks the field values on UpdateWithMultiplierRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateWithMultiplierRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateWithMultiplierRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateWithMultiplierRequestMultiError, or nil if none found.
func (m *UpdateWithMultiplierRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateWithMultiplierRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Multiplier

	// no validation rules for UserId

	// no validation rules for Product

	if len(errors) > 0 {
		return UpdateWithMultiplierRequestMultiError(errors)
	}

	return nil
}

// UpdateWithMultiplierRequestMultiError is an error wrapping multiple
// validation errors returned by UpdateWithMultiplierRequest.ValidateAll() if
// the designated constraints aren't met.
type UpdateWithMultiplierRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateWithMultiplierRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateWithMultiplierRequestMultiError) AllErrors() []error { return m }

// UpdateWithMultiplierRequestValidationError is the validation error returned
// by UpdateWithMultiplierRequest.Validate if the designated constraints
// aren't met.
type UpdateWithMultiplierRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateWithMultiplierRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateWithMultiplierRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateWithMultiplierRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateWithMultiplierRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateWithMultiplierRequestValidationError) ErrorName() string {
	return "UpdateWithMultiplierRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateWithMultiplierRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateWithMultiplierRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateWithMultiplierRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateWithMultiplierRequestValidationError{}

// Validate checks the field values on UpdateWithMultiplierResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateWithMultiplierResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateWithMultiplierResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateWithMultiplierResponseMultiError, or nil if none found.
func (m *UpdateWithMultiplierResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateWithMultiplierResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return UpdateWithMultiplierResponseMultiError(errors)
	}

	return nil
}

// UpdateWithMultiplierResponseMultiError is an error wrapping multiple
// validation errors returned by UpdateWithMultiplierResponse.ValidateAll() if
// the designated constraints aren't met.
type UpdateWithMultiplierResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateWithMultiplierResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateWithMultiplierResponseMultiError) AllErrors() []error { return m }

// UpdateWithMultiplierResponseValidationError is the validation error returned
// by UpdateWithMultiplierResponse.Validate if the designated constraints
// aren't met.
type UpdateWithMultiplierResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateWithMultiplierResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateWithMultiplierResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateWithMultiplierResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateWithMultiplierResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateWithMultiplierResponseValidationError) ErrorName() string {
	return "UpdateWithMultiplierResponseValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateWithMultiplierResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateWithMultiplierResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateWithMultiplierResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateWithMultiplierResponseValidationError{}
