// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: protocol/proto/admin_profile.proto

package profile

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

// Validate checks the field values on DeleteSubUserRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteSubUserRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteSubUserRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteSubUserRequestMultiError, or nil if none found.
func (m *DeleteSubUserRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteSubUserRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for SubuserId

	// no validation rules for NewParentId

	if len(errors) > 0 {
		return DeleteSubUserRequestMultiError(errors)
	}

	return nil
}

// DeleteSubUserRequestMultiError is an error wrapping multiple validation
// errors returned by DeleteSubUserRequest.ValidateAll() if the designated
// constraints aren't met.
type DeleteSubUserRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteSubUserRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteSubUserRequestMultiError) AllErrors() []error { return m }

// DeleteSubUserRequestValidationError is the validation error returned by
// DeleteSubUserRequest.Validate if the designated constraints aren't met.
type DeleteSubUserRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteSubUserRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteSubUserRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteSubUserRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteSubUserRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteSubUserRequestValidationError) ErrorName() string {
	return "DeleteSubUserRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteSubUserRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteSubUserRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteSubUserRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteSubUserRequestValidationError{}
