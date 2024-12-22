package handler_test

type mockError error

type mockErrorImp struct{}

func newMockError() mockError {
	return mockErrorImp{}
}

func (e mockErrorImp) Error() string {

	return "mock error"
}
