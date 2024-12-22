package service

import (
	"errors"
)

var (
	ErrInternal     = errors.New("internal error")
	ErrAlreadyExist = errors.New("already exist")
	ErrNotFound     = errors.New("not found")
)
