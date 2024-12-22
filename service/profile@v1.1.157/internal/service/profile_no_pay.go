package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/sdk/log"
)

const (
	opGetNoPayStatus     = "server.profile.GetNoPayStatus"
	opUpdateNoPayToFalse = "server.profile.UpdateNoPayToFalse"
	opCancelOffNoPay     = "server.profile.CancelOffNoPay"
)

func (p *profile) GetNoPayStatus(ctx context.Context, userUUID uuid.UUID) (bool, time.Time, error) {
	var nextTime time.Time
	noPayStatus, err := p.oldUsers.GetNoPay(ctx, userUUID)
	if err != nil {
		log.Error(ctx, "%s: get no pay status: %s, user uuid: %s", opGetNoPayStatus, err.Error(), userUUID.String())
		return false, nextTime, fmt.Errorf("%w: get no pay status", ErrInternal)
	}

	if !noPayStatus {
		return false, nextTime, nil
	}

	user, err := p.oldUsers.GetUserByUUID(ctx, userUUID)
	if err != nil {
		log.Error(ctx, "%s: get old user id: %s, user uuid: %s", opGetNoPayStatus, err.Error(), userUUID.String())
		return noPayStatus, nextTime, fmt.Errorf("%w: get old user id: %s", ErrNotFound, err.Error())
	}

	resp, err := p.jobsClient.GetActiveJob("settings", fmt.Sprintf("%d_enableNopay", user.OldID))
	if err != nil {
		log.Error(ctx, "%s: get job status: %s, user uuid: %s", opGetNoPayStatus, err.Error(), userUUID.String())
		return noPayStatus, nextTime, fmt.Errorf("%w: get job status: %s", ErrInternal, err.Error())
	}

	if len(resp.Current) > 0 {
		nextTime = resp.NextTime
	}

	return noPayStatus, nextTime, nil
}

func (p *profile) UpdateNoPayToFalse(ctx context.Context, userUUID uuid.UUID) error {
	user, err := p.oldUsers.GetUserByUUID(ctx, userUUID)
	if err != nil {
		log.Error(ctx, "%s: get old user id: %s, user uuid: %s", opUpdateNoPayToFalse, err.Error(), userUUID.String())
		return fmt.Errorf("%w: get old user id: %s", ErrNotFound, err.Error())
	}

	err = p.oldUsers.UpdateNoPay(ctx, user.OldID, false)
	if err != nil {
		log.Error(ctx, "%s: update nopay by user id: %s, user uuid: %s", opUpdateNoPayToFalse, err.Error(), userUUID.String())
		return fmt.Errorf("%w: update nopay by user id: %s", ErrInternal, err.Error())
	}

	return nil
}

func (p *profile) CancelJobOffNoPay(ctx context.Context, userUUID uuid.UUID) error {
	user, err := p.oldUsers.GetUserByUUID(ctx, userUUID)
	if err != nil {
		log.Error(ctx, "%s: get old user id: %s, user uuid: %s", opCancelOffNoPay, err.Error(), userUUID.String())
		return fmt.Errorf("%w: get old user id: %s", ErrNotFound, err.Error())
	}

	err = p.jobsClient.DeleteJobSet("settings", fmt.Sprintf("%d_enableNopay", user.OldID))
	if err != nil {
		log.Error(ctx, "%s: get job status: %s, user uuid: %s", opCancelOffNoPay, err.Error(), userUUID.String())
		return fmt.Errorf("%w: get job status: %s", ErrInternal, err.Error())
	}
	return nil
}
