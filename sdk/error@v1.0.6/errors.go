// Package errors contains errors processing logic
package errors

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "code.emcdtech.com/emcd/sdk/error/proto"
)

// Error provides information about a business error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code %s: %s", e.Code, e.Message)
}

// NewError builds a new Error.
func NewError(code string, msg string) error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

// NewGRPCError builds new business gRPC status error with details.
func NewGRPCError(code string, msg string) error {
	pbStatus, err := status.New(codes.Unknown, msg).WithDetails(&pb.ErrorDetails{
		Code:    code,
		Message: msg,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal error details: %s", err)
	}
	return pbStatus.Err()
}

func IsCode(err error, code string) bool {
	var e *Error
	if errors.As(err, &e) && e.Code == code {
		return true
	}

	return false
}
