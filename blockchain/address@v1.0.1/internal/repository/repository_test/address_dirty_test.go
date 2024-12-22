package repository_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/repository"
	"code.emcdtech.com/emcd/blockchain/address/model"
)

const dirtyLen = 10

func TestAddressDirty(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer truncateTables(ctx, t, dbPool, "address_dirty")

	repo := repository.NewAddressRepository(dbPool)

	network := nodeCommon.EthNetworkId

	t.Run("add dirty addresses success", func(t *testing.T) {
		currentTime := time.Now().UTC()
		var addressList model.AddressesDirty
		for i := 0; i < dirtyLen; i++ {
			addressList = append(addressList, newMockDirtyAddress(uuid.NewString(), network, true, currentTime))

		}

		for _, i := range rand.Perm(personalLen) {
			address := addressList[i]
			addDirtyAddress(ctx, t, repo, address)

		}

		time.Sleep(10 * time.Millisecond)

		addressList[0].IsDirty = false
		errUpdate := repo.AddOrUpdateDirtyAddress(ctx, addressList[0])

		t.Run("update dirty address", func(t *testing.T) {
			require.NoError(t, errUpdate)
		})

		t.Run("get total success", func(t *testing.T) {
			addressFilter := &model.AddressDirtyFilter{
				Address: nil,
				Network: nil,
			}

			addressRespList, err := repo.GetDirtyAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, dirtyLen)

		})

		t.Run("get new success", func(t *testing.T) {
			addressFilter := &model.AddressDirtyFilter{
				Address: &addressList[1].Address,
				Network: &network,
			}

			addressRespList, err := repo.GetDirtyAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)
			require.True(t, addressRespList[0].IsDirty)

		})

		t.Run("get updated success", func(t *testing.T) {
			addressFilter := &model.AddressDirtyFilter{
				Address: &addressList[0].Address,
				Network: &network,
			}

			addressRespList, err := repo.GetDirtyAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)
			require.False(t, addressRespList[0].IsDirty)

		})
	})
}
