package repository

import (
	"context"
	"fmt"
	"time"

	nodeCommon "code.emcdtech.com/emcd/blockchain/node/common"
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

func (r *addressRepositoryImpl) addSeriesAddressDerivedRepeat(ctx context.Context, addressUuid uuid.UUID, masterKeyId uint32, networkGroup nodeCommon.NetworkGroupEnum, repeat int) (uint32, error) {
	var derivedOffset uint32

	if err := r.WithinTransactionWithOptions(ctx, func(ctx context.Context) error {
		if derivedOffsetInside, err := r.addSeriesAddressDerived(ctx, addressUuid, masterKeyId, networkGroup); err != nil {

			return fmt.Errorf("within transaction: %w", err)
		} else {
			derivedOffset = derivedOffsetInside

			return nil
		}
	}, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.Deferrable,
		BeginQuery:     "",
	}); err != nil {
		if repeat > 0 {
			time.Sleep(repeatTimeout)

			return r.addSeriesAddressDerivedRepeat(ctx, addressUuid, masterKeyId, networkGroup, repeat-1)
		} else {

			return 0, fmt.Errorf("transaction: %w", err)
		}
	} else {

		return derivedOffset, nil
	}
}

func (r *addressRepositoryImpl) addSeriesAddressDerived(ctx context.Context, addressUuid uuid.UUID, masterKeyId uint32, networkGroup nodeCommon.NetworkGroupEnum) (uint32, error) {
	const queryLockNetworkGroup = `SELECT EXISTS(SELECT * FROM address_derived WHERE network_group = @network_group AND derived_offset >= 0 FOR UPDATE)`

	const query = `
		insert into address_derived (address_uuid, network_group, master_key_id, derived_offset)
		select @address_uuid, @network_group, @master_key_id, coalesce(max(derived_offset), -1) + 1 from address_derived 
	   	where network_group = @network_group and master_key_id = @master_key_id
	   	returning derived_offset
`

	batch := pgx.Batch{}
	batch.Queue(queryLockNetworkGroup, pgx.NamedArgs{"network_group": networkGroup})

	addressDerivedSql := newAddressDerivedSql(model.NewAddressDerived(addressUuid, networkGroup, masterKeyId))

	batch.Queue(query, addressDerivedSql.toNamedArgs())

	batchResult := r.Runner(ctx).SendBatch(ctx, &batch)
	defer func() {
		if batchResult != nil {
			if err := batchResult.Close(); err != nil {
				sdkLog.Warn(ctx, "can not to close batch: %v", err)

			}
		}
	}()

	var derivedOffset uint32

	if _, err := batchResult.Exec(); err != nil {

		return 0, fmt.Errorf("exec: %w", err)
	} else if err := batchResult.QueryRow().Scan(&derivedOffset); err != nil {

		return 0, fmt.Errorf("scan offset: %w", err)
	} else if err := batchResult.Close(); err != nil {

		return 0, fmt.Errorf("batch close: %w", err)
	} else {
		batchResult = nil

		return derivedOffset, nil
	}
}
