package repository_test

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/repository"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

const derivedLen = 100
const otherLen = 100
const masterKeyId = 0

func createDerivedFunc(networkGroup nodeCommon.NetworkGroupEnum) repository.DerivedFunc {

	return func(offset uint32) (string, error) {

		return fmt.Sprintf("d:%s:%d", networkGroup.ToString(), offset), nil
	}
}

func TestAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer truncateTables(ctx, t, dbPool, "address", "address_derived")

	repo := repository.NewAddressRepository(dbPool)

	t.Run("add new addresses concurrently success", func(t *testing.T) {
		networkGroup := nodeCommon.EthNetworkGroupId
		var addressList model.Addresses
		for i := 0; i < derivedLen; i++ {
			addressList = append(addressList, newMockDerivedAddress(uuid.New(), networkGroup))

		}

		for i := 0; i < otherLen; i++ {
			addressList = append(addressList, newMockCommonAddress(uuid.New(), networkGroup))

		}

		wg := new(sync.WaitGroup)
		for _, i := range rand.Perm(derivedLen + otherLen) {
			address := addressList[i]
			if address.AddressType.Number() == addressPb.AddressType_ADDRESS_TYPE_DERIVED.Number() {
				wg.Add(1)
				go addDerivedAddress(ctx, wg, t, repo, address, masterKeyId, createDerivedFunc(networkGroup))

			} else {
				wg.Add(1)
				go addCommonAddress(ctx, wg, t, repo, address)

			}
		}

		wg.Wait()

		time.Sleep(100 * time.Millisecond)
		currentTime := time.Now().UTC()
		time.Sleep(2 * time.Millisecond)

		userUuid1 := uuid.New()
		userUuid2 := uuid.New()
		processingUuid := uuid.New()

		address1 := newMockDerivedAddress(userUuid1, networkGroup)
		address2 := newMockDerivedProcessingAddress(userUuid2, processingUuid, networkGroup)

		err1 := repo.AddNewDerivedAddress(ctx, address1, masterKeyId, createDerivedFunc(networkGroup))
		err2 := repo.AddNewDerivedAddress(ctx, address2, masterKeyId, createDerivedFunc(networkGroup))

		t.Run("add derived count success", func(t *testing.T) {
			require.NoError(t, err1)
			require.Equal(t, address1.GetAddressDerived().AddressUuid, address1.Id)
			require.Equal(t, address1.GetAddressDerived().DerivedOffset, uint32(derivedLen+0))

			require.NoError(t, err2)
			require.Equal(t, address2.GetAddressDerived().AddressUuid, address2.Id)
			require.Equal(t, address2.GetAddressDerived().DerivedOffset, uint32(derivedLen+1))
		})

		t.Run("get total success", func(t *testing.T) {
			addressFilter := &model.AddressFilter{
				Id:       nil,
				Address:  nil,
				UserUuid: nil,
			}

			_, addressRespList, err := repo.GetNewAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, derivedLen+otherLen+2)

		})

		t.Run("get id success", func(t *testing.T) {
			addressFilter := &model.AddressFilter{
				Id: &address1.Id,
			}

			_, addressRespList, err := repo.GetNewAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)

		})

		t.Run("get address success", func(t *testing.T) {
			addressFilter := &model.AddressFilter{
				Address: &address1.Address,
			}

			_, addressRespList, err := repo.GetNewAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)

		})

		t.Run("get user_uuid success", func(t *testing.T) {
			addressFilter := &model.AddressFilter{
				UserUuid: &address1.UserUuid,
			}

			_, addressRespList, err := repo.GetNewAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)

		})

		t.Run("get non processing success", func(t *testing.T) {
			addressFilter := &model.AddressFilter{
				UserUuid:     &address1.UserUuid,
				IsProcessing: utils.BoolToPtr(false),
			}

			_, addressRespList, err := repo.GetNewAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)
			require.Equal(t, addressRespList[0].Address, address1.Address)

		})

		t.Run("get processing success", func(t *testing.T) {
			addressFilter := &model.AddressFilter{
				IsProcessing: utils.BoolToPtr(true),
			}

			_, addressRespList, err := repo.GetNewAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)
			require.Equal(t, addressRespList[0].Address, address2.Address)

		})

		t.Run("get created_at_gt success", func(t *testing.T) {
			addressFilter := &model.AddressFilter{
				CreatedAtGt: &currentTime,
			}

			_, addressRespList, err := repo.GetNewAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 2)

		})

		t.Run("get pagination success", func(t *testing.T) {
			addressFilter := &model.AddressFilter{
				Pagination: &model.Pagination{
					Limit:  1,
					Offset: 0,
				},
			}

			totalCount, addressRespList, err := repo.GetNewAddresses(ctx, addressFilter)
			require.NotNil(t, totalCount)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)
			require.Equal(t, *totalCount, uint64(derivedLen+otherLen+2))

		})
	})
}
