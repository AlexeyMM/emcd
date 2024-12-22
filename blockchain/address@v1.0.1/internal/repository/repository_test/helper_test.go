package repository_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/repository"
	"code.emcdtech.com/emcd/blockchain/address/model"
	addressPb "code.emcdtech.com/emcd/blockchain/address/protocol/address"
)

const addressCommonType = addressPb.AddressType_ADDRESS_TYPE_BASED_ID

func newMockDerivedAddress(userUuid uuid.UUID, networkGroup nodeCommon.NetworkGroupEnum) *model.Address {

	return model.NewAddress(
		uuid.New(),
		"",
		userUuid,
		addressPb.AddressType_ADDRESS_TYPE_DERIVED,
		networkGroup,
	)
}

func newMockDerivedProcessingAddress(userUuid, processingUuid uuid.UUID, networkGroup nodeCommon.NetworkGroupEnum) *model.Address {

	return model.NewProcessingAddress(
		uuid.New(),
		"",
		userUuid,
		processingUuid,
		addressPb.AddressType_ADDRESS_TYPE_DERIVED,
		networkGroup,
	)
}

func newMockCommonAddress(userUuid uuid.UUID, networkGroup nodeCommon.NetworkGroupEnum) *model.Address {

	return model.NewAddress(
		uuid.New(),
		uuid.New().String(),
		userUuid,
		addressCommonType,
		networkGroup,
	)
}

func newMockOldAddress(userUuid uuid.UUID, network nodeCommon.NetworkEnum, userAccountId int32, coin string) *model.AddressOld {

	return model.NewAddressOld(
		uuid.New(),
		uuid.New().String(),
		userUuid,
		addressCommonType,
		network,
		userAccountId,
		coin,
	)
}

func newMockDirtyAddress(address string, network nodeCommon.NetworkEnum, isDirty bool, currentTime time.Time) *model.AddressDirty {

	return model.NewAddressDirty(
		address,
		network,
		isDirty,
		currentTime,
	)
}

func newMockPersonalAddress(address string, userUuid uuid.UUID, network nodeCommon.NetworkEnum) *model.AddressPersonal {

	return model.NewAddressPersonal(
		uuid.New(),
		address,
		userUuid,
		network,
		0.0,
	)
}

func addDerivedAddress(ctx context.Context, wg *sync.WaitGroup, t *testing.T, addressRepo repository.AddressRepository, address *model.Address, masterKeyId uint32, derivedFunc repository.DerivedFunc) {
	defer wg.Done()

	if err := addressRepo.AddNewDerivedAddress(ctx, address, masterKeyId, derivedFunc); err != nil {
		t.Error(err)
		t.Fail()

	}
}

func addCommonAddress(ctx context.Context, wg *sync.WaitGroup, t *testing.T, addressRepo repository.AddressRepository, address *model.Address) {
	defer wg.Done()

	if err := addressRepo.AddNewCommonAddress(ctx, address); err != nil {
		t.Errorf(err.Error())
		t.Fail()

	}
}

func addOldAddress(ctx context.Context, wg *sync.WaitGroup, t *testing.T, addressRepo repository.AddressRepository, addressOld *model.AddressOld) {
	defer wg.Done()

	if err := addressRepo.AddOldAddress(ctx, addressOld); err != nil {
		t.Errorf(err.Error())
		t.Fail()

	}
}

func addDirtyAddress(ctx context.Context, t *testing.T, addressRepo repository.AddressRepository, addressDirty *model.AddressDirty) {
	if err := addressRepo.AddOrUpdateDirtyAddress(ctx, addressDirty); err != nil {
		t.Errorf(err.Error())
		t.Fail()

	}
}

func addPersonalAddress(ctx context.Context, t *testing.T, addressRepo repository.AddressRepository, addressPersonal *model.AddressPersonal) {
	if err := addressRepo.AddPersonalAddress(ctx, addressPersonal); err != nil {
		t.Errorf(err.Error())
		t.Fail()

	}
}

func truncateTables(ctx context.Context, t *testing.T, dbPool *pgxpool.Pool, tables ...string) {
	for _, table := range tables {
		_, err := dbPool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s cascade", table))
		require.NoError(t, err)
	}
}
