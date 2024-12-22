package repository_test

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/common/utils"
	"code.emcdtech.com/emcd/blockchain/address/internal/repository"
	"code.emcdtech.com/emcd/blockchain/address/model"
)

const personalLen = 10
const epsilon10 = 1e-10

func TestAddressPersonal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer truncateTables(ctx, t, dbPool, "address_personal")

	repo := repository.NewAddressRepository(dbPool)

	network := nodeCommon.EthNetworkId

	t.Run("add personal addresses success", func(t *testing.T) {
		var addressList model.AddressesPersonal
		for i := 0; i < personalLen; i++ {
			addressList = append(addressList, newMockPersonalAddress(uuid.NewString(), uuid.New(), network))

		}

		for _, i := range rand.Perm(personalLen) {
			address := addressList[i]
			go addPersonalAddress(ctx, t, repo, address)

		}

		time.Sleep(10 * time.Millisecond)
		currentTime := time.Now().UTC()
		time.Sleep(2 * time.Millisecond)

		userUuid := uuid.New()
		address := newMockPersonalAddress(uuid.NewString(), userUuid, network)

		err := repo.AddPersonalAddress(ctx, address)

		t.Run("add personal address", func(t *testing.T) {
			require.NoError(t, err)
		})

		addressStr := uuid.NewString()
		addressPartial1 := &model.AddressPersonalPartial{
			Address:   &addressStr,
			MinPayout: nil,
			DeletedAt: nil,
			UpdatedAt: &currentTime,
		}

		err = repo.UpdatePersonalAddress(ctx, address, addressPartial1)

		t.Run("update personal address", func(t *testing.T) {
			require.NoError(t, err)
			require.Equal(t, address.Address, *addressPartial1.Address)
		})

		addressPartial2 := &model.AddressPersonalPartial{
			Address:   nil,
			MinPayout: utils.Float64ToPtr(1.0),
			DeletedAt: nil,
			UpdatedAt: &currentTime,
		}

		err = repo.UpdatePersonalAddress(ctx, address, addressPartial2)

		t.Run("update minimum payout", func(t *testing.T) {
			require.NoError(t, err)
			require.InEpsilon(t, address.MinPayout, *addressPartial2.MinPayout, epsilon10)
			require.False(t, address.DeletedAt.Valid)
		})

		addressPartial3 := &model.AddressPersonalPartial{
			Address:   utils.StringToPtr(""),
			MinPayout: nil,
			DeletedAt: &sql.NullTime{Time: currentTime, Valid: true},
			UpdatedAt: &currentTime,
		}

		err = repo.UpdatePersonalAddress(ctx, address, addressPartial3)

		t.Run("delete personal address", func(t *testing.T) {
			require.NoError(t, err)
			require.Equal(t, address.Address, *addressPartial3.Address)
			require.True(t, address.DeletedAt.Valid)
		})

		t.Run("get total success", func(t *testing.T) {
			addressFilter := &model.AddressPersonalFilter{
				Id:         nil,
				Address:    nil,
				UserUuid:   nil,
				Network:    nil,
				IsDeleted:  nil,
				Pagination: nil,
			}

			_, addressRespList, err := repo.GetPersonalAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, personalLen+1)

		})

		t.Run("get id success", func(t *testing.T) {
			addressFilter := &model.AddressPersonalFilter{
				Id: &address.Id,
			}

			_, addressRespList, err := repo.GetPersonalAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)

		})

		t.Run("get address success", func(t *testing.T) {
			addressFilter := &model.AddressPersonalFilter{
				Address: &address.Address,
			}

			_, addressRespList, err := repo.GetPersonalAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)

		})

		t.Run("get user_uuid success", func(t *testing.T) {
			addressFilter := &model.AddressPersonalFilter{
				UserUuid: &address.UserUuid,
			}

			_, addressRespList, err := repo.GetPersonalAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)

		})

		t.Run("get deleted success", func(t *testing.T) {
			addressFilter := &model.AddressPersonalFilter{
				IsDeleted: utils.BoolToPtr(true),
			}

			_, addressRespList, err := repo.GetPersonalAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)

		})

		t.Run("get non-deleted success", func(t *testing.T) {
			addressFilter := &model.AddressPersonalFilter{
				IsDeleted: utils.BoolToPtr(false),
			}

			_, addressRespList, err := repo.GetPersonalAddresses(ctx, addressFilter)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, personalLen)

		})

		t.Run("get pagination success", func(t *testing.T) {
			addressFilter := &model.AddressPersonalFilter{
				Pagination: &model.Pagination{
					Limit:  1,
					Offset: 0,
				},
			}

			totalCount, addressRespList, err := repo.GetPersonalAddresses(ctx, addressFilter)
			require.NotNil(t, totalCount)
			require.NotEmpty(t, addressRespList)
			require.NoError(t, err)

			require.Len(t, addressRespList, 1)
			require.Equal(t, *totalCount, uint64(personalLen+1))

		})
	})
}
