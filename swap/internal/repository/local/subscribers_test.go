package local

import (
	"testing"

	"code.emcdtech.com/b2b/swap/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSubscribers_AddAndGetBySwapID(t *testing.T) {
	subs := NewSubscribers()

	swapID := uuid.New()
	clientID := uuid.New()
	ch := make(chan model.PublicStatus, 1)

	subs.Add(swapID, clientID, ch)

	result := subs.GetBySwapID(swapID)
	require.Len(t, result, 1)
	require.Equal(t, clientID, result[0].ClientID)
	require.Equal(t, ch, result[0].Ch)
}

func TestSubscribers_AddMultipleSubscribers(t *testing.T) {
	subs := NewSubscribers()

	swapID := uuid.New()
	clientID1 := uuid.New()
	clientID2 := uuid.New()

	ch1 := make(chan model.PublicStatus, 1)
	ch2 := make(chan model.PublicStatus, 1)

	subs.Add(swapID, clientID1, ch1)
	subs.Add(swapID, clientID2, ch2)

	result := subs.GetBySwapID(swapID)
	require.Len(t, result, 2)
	require.ElementsMatch(t, []*model.Subscriber{
		{ClientID: clientID1, Ch: ch1},
		{ClientID: clientID2, Ch: ch2},
	}, result)
}

func TestSubscribers_DeleteSubscriber(t *testing.T) {
	subs := NewSubscribers()

	swapID := uuid.New()
	clientID1 := uuid.New()
	clientID2 := uuid.New()

	ch1 := make(chan model.PublicStatus, 1)
	ch2 := make(chan model.PublicStatus, 1)

	subs.Add(swapID, clientID1, ch1)
	subs.Add(swapID, clientID2, ch2)

	subs.Delete(swapID, clientID1)

	result := subs.GetBySwapID(swapID)
	require.Len(t, result, 1)
	require.Equal(t, clientID2, result[0].ClientID)
}

func TestSubscribers_DeleteOnlySubscriber(t *testing.T) {
	subs := NewSubscribers()

	swapID := uuid.New()
	clientID := uuid.New()

	ch := make(chan model.PublicStatus, 1)

	subs.Add(swapID, clientID, ch)

	subs.Delete(swapID, clientID)

	result := subs.GetBySwapID(swapID)
	require.Len(t, result, 0)
}

func TestSubscribers_DeleteNonExistentSubscriber(t *testing.T) {
	subs := NewSubscribers()

	swapID := uuid.New()
	clientID1 := uuid.New()
	clientID2 := uuid.New()

	ch1 := make(chan model.PublicStatus, 1)
	ch2 := make(chan model.PublicStatus, 1)

	subs.Add(swapID, clientID1, ch1)
	subs.Add(swapID, clientID2, ch2)

	subs.Delete(swapID, uuid.New())

	result := subs.GetBySwapID(swapID)
	require.Len(t, result, 2)
}

func TestSubscribers_GetByNonExistentSwapID(t *testing.T) {
	subs := NewSubscribers()

	result := subs.GetBySwapID(uuid.New())
	require.Len(t, result, 0)
}
