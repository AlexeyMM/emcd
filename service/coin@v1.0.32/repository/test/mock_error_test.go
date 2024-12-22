package test

import "errors"

type mockError error

type mockImpError struct{}

func newMockError() mockError {
	return mockImpError{}
}

func (e mockImpError) Error() string {
	return errors.New("mock error").Error()
}
