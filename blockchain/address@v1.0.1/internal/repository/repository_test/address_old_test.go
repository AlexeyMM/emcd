package repository_test

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/repository"
	"code.emcdtech.com/emcd/blockchain/address/model"
)

const oldLen = 100

func TestAddressOld(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer truncateTables(ctx, t, dbPool, "address_old")

	repo := repository.NewAddressRepository(dbPool)

	network := nodeCommon.EthNetworkId
	userAccountId := int32(1)
	coin := "eth"

	t.Run("add old addresses concurrently success", func(t *testing.T) {
		var addressList model.AddressesOld
		for i := 0; i < oldLen; i++ {
			userAccountId++
			addressList = append(addressList, newMockOldAddress(uuid.New(), network, userAccountId, coin))

		}

		wg := new(sync.WaitGroup)
		for _, i := range rand.Perm(oldLen) {
			address := addressList[i]
			wg.Add(1)
			go addOldAddress(ctx, wg, t, repo, address)

		}

		wg.Wait()

		time.Sleep(10 * time.Millisecond)
		currentTime := time.Now().UTC()
		time.Sleep(2 * time.Millisecond)

		userUuid := uuid.New()
		userAccountId++
		address := newMockOldAddress(userUuid, network, userAccountId, coin)

		err := repo.AddOldAddress(ctx, address)

		t.Run("add old address", func(t *testing.T) {
			require.NoError(t, err)
		})

		t.Run("get total success", func(t *testing.T) {
			addressFilter := &model.AddressOldFilter{
				Id:       nil,
				Address:  nil,
				UserUuid: nil,
			}

			_, addressRespList, err := repo.GetOldAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, oldLen+1)

		})

		t.Run("get id success", func(t *testing.T) {
			addressFilter := &model.AddressOldFilter{
				Id: &address.Id,
			}

			_, addressRespList, err := repo.GetOldAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)

		})

		t.Run("get address success", func(t *testing.T) {
			addressFilter := &model.AddressOldFilter{
				Address: &address.Address,
			}

			_, addressRespList, err := repo.GetOldAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)

		})

		t.Run("get user_uuid success", func(t *testing.T) {
			addressFilter := &model.AddressOldFilter{
				UserUuid: &address.UserUuid,
			}

			_, addressRespList, err := repo.GetOldAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)

		})

		t.Run("get created_at_gt success", func(t *testing.T) {
			addressFilter := &model.AddressOldFilter{
				CreatedAtGt: &currentTime,
			}

			_, addressRespList, err := repo.GetOldAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)

		})

		t.Run("get pagination success", func(t *testing.T) {
			addressFilter := &model.AddressOldFilter{
				Pagination: &model.Pagination{
					Limit:  1,
					Offset: 0,
				},
			}

			totalCount, addressRespList, err := repo.GetOldAddresses(ctx, addressFilter)
			require.NotNil(t, totalCount)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)
			require.Equal(t, *totalCount, uint64(oldLen+1))

		})
	})
}
