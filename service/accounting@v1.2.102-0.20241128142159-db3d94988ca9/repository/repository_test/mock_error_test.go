package repository_test

import "errors"

type mockError error

type mockErrorImp struct{}

func newMockError() mockError {

	return mockErrorImp{}
}

func (e mockErrorImp) Error() string {

	return errors.New("mock error").Error()
}
