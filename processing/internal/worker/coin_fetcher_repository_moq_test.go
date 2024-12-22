// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package worker

import (
	"context"
	"sync"
)

// Ensure, that RepositoryMock does implement Repository.
// If this is not the case, regenerate this file with moq.
var _ Repository = &RepositoryMock{}

// RepositoryMock is a mock implementation of Repository.
//
//	func TestSomethingThatUsesRepository(t *testing.T) {
//
//		// make and configure a mocked Repository
//		mockedRepository := &RepositoryMock{
//			FetchCoinsFunc: func(ctx context.Context) error {
//				panic("mock out the FetchCoins method")
//			},
//		}
//
//		// use mockedRepository in code that requires Repository
//		// and then make assertions.
//
//	}
type RepositoryMock struct {
	// FetchCoinsFunc mocks the FetchCoins method.
	FetchCoinsFunc func(ctx context.Context) error

	// calls tracks calls to the methods.
	calls struct {
		// FetchCoins holds details about calls to the FetchCoins method.
		FetchCoins []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
	}
	lockFetchCoins sync.RWMutex
}

// FetchCoins calls FetchCoinsFunc.
func (mock *RepositoryMock) FetchCoins(ctx context.Context) error {
	if mock.FetchCoinsFunc == nil {
		panic("RepositoryMock.FetchCoinsFunc: method is nil but Repository.FetchCoins was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockFetchCoins.Lock()
	mock.calls.FetchCoins = append(mock.calls.FetchCoins, callInfo)
	mock.lockFetchCoins.Unlock()
	return mock.FetchCoinsFunc(ctx)
}

// FetchCoinsCalls gets all the calls that were made to FetchCoins.
// Check the length with:
//
//	len(mockedRepository.FetchCoinsCalls())
func (mock *RepositoryMock) FetchCoinsCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockFetchCoins.RLock()
	calls = mock.calls.FetchCoins
	mock.lockFetchCoins.RUnlock()
	return calls
}