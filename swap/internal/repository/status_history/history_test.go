package status_history

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/b2b/swap/mocks/internal_/repository"
	"code.emcdtech.com/b2b/swap/model"
	"code.emcdtech.com/b2b/swap/package/gokit"
)

func Test_Add(t *testing.T) {
	ctx := context.Background()
	swapId := uuid.New()
	validSwap := &model.Swap{ID: swapId, Status: model.Unknown}

	tests := []struct {
		name                        string
		swapRepositoryMock          func() *repository.MockSwap
		statusHistoryRepositoryMock func() *repository.MockSwapStatusHistory
		swap                        *model.Swap
		expectError                 bool
	}{
		{
			name: "successful add",
			swapRepositoryMock: func() *repository.MockSwap {
				mockSwap := &repository.MockSwap{}
				mockSwap.On("Add", ctx, validSwap).Return(nil)
				return mockSwap
			},
			statusHistoryRepositoryMock: func() *repository.MockSwapStatusHistory {
				mockStatusHistory := &repository.MockSwapStatusHistory{}
				mockStatusHistory.On("Add", ctx, validSwap.ID, mock.Anything).Return(nil)
				return mockStatusHistory
			},
			swap:        validSwap,
			expectError: false,
		},
		{
			name: "error from swap repository add",
			swapRepositoryMock: func() *repository.MockSwap {
				mockSwap := &repository.MockSwap{}
				mockSwap.On("Add", ctx, validSwap).Return(fmt.Errorf("swap repository add error"))
				return mockSwap
			},
			statusHistoryRepositoryMock: func() *repository.MockSwapStatusHistory {
				mockStatusHistory := &repository.MockSwapStatusHistory{}
				return mockStatusHistory
			},
			swap:        validSwap,
			expectError: true,
		},
		{
			name: "error adding to status history",
			swapRepositoryMock: func() *repository.MockSwap {
				mockSwap := &repository.MockSwap{}
				mockSwap.On("Add", ctx, validSwap).Return(nil)
				return mockSwap
			},
			statusHistoryRepositoryMock: func() *repository.MockSwapStatusHistory {
				mockStatusHistory := &repository.MockSwapStatusHistory{}
				mockStatusHistory.On("Add", ctx, validSwap.ID, mock.Anything).Return(fmt.Errorf("status history add error"))
				return mockStatusHistory
			},
			swap:        validSwap,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			swapRepo := test.swapRepositoryMock()
			statusHistoryRepo := test.statusHistoryRepositoryMock()

			history := NewHistory(swapRepo, statusHistoryRepo)

			err := history.Add(ctx, test.swap)
			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			swapRepo.AssertExpectations(t)
			statusHistoryRepo.AssertExpectations(t)
		})
	}
}

func Test_Update(t *testing.T) {
	ctx := context.Background()

	swapId := uuid.New()
	validSwap := &model.Swap{ID: swapId, Status: model.Unknown}
	validPartial := &model.SwapPartial{Status: gokit.Ptr(model.WaitDeposit)}
	emptyPartial := &model.SwapPartial{}
	validFilter := &model.SwapFilter{}

	tests := []struct {
		name                        string
		swapRepositoryMock          func() *repository.MockSwap
		statusHistoryRepositoryMock func() *repository.MockSwapStatusHistory
		swap                        *model.Swap
		filter                      *model.SwapFilter
		partial                     *model.SwapPartial
		expectError                 bool
	}{
		{
			name: "successful update with non-nil status",
			swapRepositoryMock: func() *repository.MockSwap {
				mockSwap := &repository.MockSwap{}
				mockSwap.On("Update", ctx, validSwap, validFilter, validPartial).Return(nil)
				return mockSwap
			},
			statusHistoryRepositoryMock: func() *repository.MockSwapStatusHistory {
				mockStatusHistory := &repository.MockSwapStatusHistory{}
				mockStatusHistory.On("Add", ctx, validSwap.ID, mock.Anything).Return(nil)
				return mockStatusHistory
			},
			swap:        validSwap,
			filter:      validFilter,
			partial:     validPartial,
			expectError: false,
		},
		{
			name: "successful update with nil status",
			swapRepositoryMock: func() *repository.MockSwap {
				mockSwap := &repository.MockSwap{}
				mockSwap.On("Update", ctx, validSwap, validFilter, emptyPartial).Return(nil)
				return mockSwap
			},
			statusHistoryRepositoryMock: func() *repository.MockSwapStatusHistory {
				mockStatusHistory := &repository.MockSwapStatusHistory{}
				return mockStatusHistory
			},
			swap:        validSwap,
			filter:      validFilter,
			partial:     emptyPartial,
			expectError: false,
		},
		{
			name: "error from swap repository update",
			swapRepositoryMock: func() *repository.MockSwap {
				mockSwap := &repository.MockSwap{}
				mockSwap.On("Update", ctx, validSwap, validFilter, validPartial).Return(fmt.Errorf("update error"))
				return mockSwap
			},
			statusHistoryRepositoryMock: func() *repository.MockSwapStatusHistory {
				mockStatusHistory := &repository.MockSwapStatusHistory{}
				return mockStatusHistory
			},
			swap:        validSwap,
			filter:      validFilter,
			partial:     validPartial,
			expectError: true,
		},
		{
			name: "error adding to status history",
			swapRepositoryMock: func() *repository.MockSwap {
				mockSwap := &repository.MockSwap{}
				mockSwap.On("Update", ctx, validSwap, validFilter, validPartial).Return(nil)
				return mockSwap
			},
			statusHistoryRepositoryMock: func() *repository.MockSwapStatusHistory {
				mockStatusHistory := &repository.MockSwapStatusHistory{}
				mockStatusHistory.On("Add", ctx, validSwap.ID, mock.Anything).Return(fmt.Errorf("status history error"))
				return mockStatusHistory
			},
			swap:        validSwap,
			filter:      validFilter,
			partial:     validPartial,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			swapRepo := test.swapRepositoryMock()
			statusHistoryRepo := test.statusHistoryRepositoryMock()

			history := NewHistory(swapRepo, statusHistoryRepo)

			err := history.Update(ctx, test.swap, test.filter, test.partial)
			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			swapRepo.AssertExpectations(t)
			statusHistoryRepo.AssertExpectations(t)
		})
	}
}
