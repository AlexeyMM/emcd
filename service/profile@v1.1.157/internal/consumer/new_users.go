package consumer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"

	"code.emcdtech.com/emcd/sdk/log"
	pgTx "code.emcdtech.com/emcd/sdk/pg"

	"code.emcdtech.com/emcd/service/profile/internal/model"
	"code.emcdtech.com/emcd/service/profile/internal/service"
)

type FinderNewUsers struct {
	emcdTrx pgTx.PgxTransactor
	wl      whitelabel.WhitelabelServiceClient
	serv    service.Profile
}

func NewFinderNewUsers(emcdTrx pgTx.PgxTransactor, wl whitelabel.WhitelabelServiceClient, serv service.Profile) *FinderNewUsers {
	return &FinderNewUsers{
		emcdTrx: emcdTrx,
		wl:      wl,
		serv:    serv,
	}
}

func (f *FinderNewUsers) Consume(ctx context.Context) {
	log.Info(ctx, "finderNewUsers consumer started")
	defer func() {
		log.Info(ctx, "end FinderNewUsers consumer")
	}()

	t := time.NewTicker(time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			f.find(ctx)
		}
	}
}

func (f *FinderNewUsers) find(ctx context.Context) {
	query := `SELECT id, email, is_active, ref_id FROM emcd.users WHERE new_id is NULL`
	rows, err := f.emcdTrx.Runner(ctx).Query(ctx, query)
	if err != nil {
		log.Error(ctx, "finderNewUsers query: %s", err.Error())
		return
	}
	for rows.Next() {
		var (
			id       int32
			email    string
			isActive bool
			refID    int
		)
		err = rows.Scan(&id, &email, &isActive, &refID)
		if err != nil {
			log.Error(ctx, "finderNewUsers couldn't scan id, email: %s", err.Error())
			continue
		}

		log.Info(ctx, "finderNewUsers execute id: %d, email: %s, refID: %d", id, email, refID)
		newRefID, err := f.getNewIDByOldID(ctx, int32(refID))
		if err != nil {
			log.Error(ctx, "finderNewUsers getNewIDByOldID: %s, email: %s", err.Error(), email)
			continue
		}

		var whitelabelID uuid.UUID
		resp, err := f.wl.GetByUserID(ctx, &whitelabel.GetByUserIDRequest{UserId: id})
		if err != nil {
			log.Error(ctx, "finderNewUsers wl.GetByUserID: %s, email: %s", err.Error(), email)
			continue
		}
		if resp.WhiteLabel != nil {
			log.Info(
				ctx,
				"finderNewUsers detected wl: %s by user_id: %d, email: %s",
				resp.WhiteLabel.Id,
				resp.WhiteLabel.UserId,
				email,
			)
			whitelabelID, err = uuid.Parse(resp.WhiteLabel.Id)
			if err != nil {
				log.Error(ctx, "finderNewUsers: parse wl: %s, email: %s", err.Error(), email)
				continue
			}
		}

		user, err := f.serv.GetUserByEmailAndWl(ctx, email, whitelabelID)
		if err != nil {
			log.Error(ctx, "finderNewUsers getUserByEmailAndWl: %s, email: %s", err.Error(), email)
			continue
		}
		if user != nil {
			log.Info(ctx, "finderNewUsers userAlreadyExistInProfile: %s", email)
			continue
		}

		pr, err := f.getOldUserFromEmcd(ctx, email)
		if err != nil {
			log.Error(ctx, "finderNewUsers getOldByEmailAndWlV2: %s, email: %s", err.Error(), email)
			continue
		}
		pr.User.ID = uuid.New()
		pr.User.IsActive = isActive
		pr.User.NewRefID = newRefID
		pr.User.WhiteLabelID = whitelabelID

		// if user.ref_id=whitelabel.user_id then user.whitelabel_id=whitelabel.id
		if whitelabelID == uuid.Nil {
			wl, err := f.wl.GetByUserID(ctx, &whitelabel.GetByUserIDRequest{
				UserId: int32(refID),
			})
			if err != nil {
				log.Error(ctx, "finderNewUsers: wl.getByUserID: %s, email: %s", err.Error(), email)
				continue
			}
			if wl.WhiteLabel != nil {
				wlID, err := uuid.Parse(wl.WhiteLabel.Id)
				if err != nil {
					log.Error(
						ctx,
						"finderNewUsers: parse whitelabel_id: %s, value: %s, email: %s",
						err.Error(),
						wl.WhiteLabel.Id,
						email,
					)
				}
				pr.User.WhiteLabelID = wlID
			}
		}

		_, err = f.serv.SaveV3(ctx, pr)
		if err != nil {
			log.Error(ctx, "finderNewUsers saveV3: %s", err.Error())
			continue
		}
	}
}

func (f *FinderNewUsers) getOldUserFromEmcd(ctx context.Context, email string) (*model.Profile, error) {
	prof, err := f.serv.GetOldByEmailAndWl(ctx, email, uuid.Nil)
	if err != nil {
		return nil, fmt.Errorf("getOldByEmailAndWlV2: %w", err)
	}
	if prof.User == nil {
		return nil, fmt.Errorf("getOldByEmailAndWlV2: returned nil by email: %s", email)
	}
	return &model.Profile{
		User: prof.User,
	}, nil
}

func (f *FinderNewUsers) getNewIDByOldID(ctx context.Context, userID int32) (uuid.UUID, error) {
	query := `SELECT new_id FROM emcd.users WHERE id = $1`
	var newUserID uuid.UUID
	err := f.emcdTrx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&newUserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("queryRow: %w", err)
	}
	return newUserID, nil
}
