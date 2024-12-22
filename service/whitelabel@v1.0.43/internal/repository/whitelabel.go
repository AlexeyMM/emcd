package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"code.emcdtech.com/emcd/service/whitelabel/internal/model"
)

const (
	numberOfCoinTypes       = 13
	searchPatternStartsWith = `%s%%`
)

type WhiteLabel interface {
	Create(ctx context.Context, wl *model.WhiteLabel) error
	GetAll(ctx context.Context, skip, take int, orderBy string, asc bool) ([]*model.WhiteLabel, int, error)
	Update(ctx context.Context, wl *model.WhiteLabel) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.WhiteLabel, error)
	GetBySegmentID(ctx context.Context, segmentID int) (*model.WhiteLabel, error)
	GetByUserID(ctx context.Context, userID int) (*model.WhiteLabel, error)
	GetByOrigin(ctx context.Context, origin string) (*model.WhiteLabel, error)
	CheckByUserID(ctx context.Context, userID int) (bool, error)
	CheckByUserIDAndOrigin(ctx context.Context, userID int, origin string) (bool, error)
	GetV2WLs(ctx context.Context) ([]*model.WhiteLabel, error)
	GetConfigByOrigin(ctx context.Context, origin string) (*model.WlConfig, error)
	SetConfigByRefID(ctx context.Context, conf *model.WlConfig) error
	SetAllowOrigin(ctx context.Context, req *model.AllowOrigin) error
	GetAllowOrigins(ctx context.Context) ([]*model.AllowOrigin, error)
	SetStratum(ctx context.Context, req *model.Stratum) error
	GetFullByUserID(ctx context.Context, userID int) (*model.WhiteLabel, error)
	GetCoins(ctx context.Context, wlId uuid.UUID) ([]*model.WLCoins, error)
	AddCoin(ctx context.Context, wlID uuid.UUID, coinID string) error
	DeleteCoin(ctx context.Context, wlID uuid.UUID, coinID string) error
	GetWLStratumList(ctx context.Context, refID string) ([]model.Stratum, error)
	GetWLStratumListV2(ctx context.Context, refID int32) ([]model.Stratum, error)
}

type whiteLabel struct {
	pool *pgxpool.Pool
}

func NewWhiteLabel(pool *pgxpool.Pool) *whiteLabel {
	return &whiteLabel{
		pool: pool,
	}
}

func (w *whiteLabel) Create(ctx context.Context, wl *model.WhiteLabel) error {
	tx, err := w.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("whitelabel create: %w", err)
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error().Msgf("whitelabel create: rollback: %v, whitelabel id: %d", err, wl.ID)
		}
	}()

	if err = w.create(ctx, tx, wl); err != nil {
		return fmt.Errorf("whitelabel create: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("whitelabel create: commit: %w", err)
	}
	return nil
}

const (
	createNewWL = `
		INSERT INTO whitelabel
		    (id, user_id, segment_id, origin, prefix, sender_email, domain, api_key, url, version, master_slave, master_fee,
		     is_two_fa_enabled, is_captcha_enabled, is_email_confirm_enabled)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
)

func (w *whiteLabel) create(ctx context.Context, tx pgx.Tx, wl *model.WhiteLabel) error {
	_, err := tx.Exec(ctx, createNewWL,
		wl.ID, wl.UserID, wl.SegmentID, wl.Origin, wl.Prefix, wl.SenderEmail, wl.Domain,
		wl.APIKey, wl.URL, wl.Version, wl.MasterSlave, wl.MasterFee, wl.IsTwoFAEnabled, wl.IsCaptchaEnabled,
		wl.IsEmailConfirmEnabled)
	if err != nil {
		return fmt.Errorf("create to whitelabel table: %w", err)
	}
	return nil
}

func (w *whiteLabel) GetAll(ctx context.Context, skip, take int, orderBy string, asc bool) ([]*model.WhiteLabel, int, error) {
	orderDirection := "DESC"
	if asc {
		orderDirection = "ASC"
	}
	b := pgx.Batch{}
	query := `
	  	SELECT id, user_id, segment_id, origin, prefix, sender_email, domain, api_key, url, version, master_slave,
	  	       master_fee, is_two_fa_enabled, is_captcha_enabled, is_email_confirm_enabled
		FROM whitelabel
		ORDER BY %s %s
		OFFSET $1 LIMIT $2
	`
	b.Queue(fmt.Sprintf(query, orderBy, orderDirection), skip, take)
	b.Queue("SELECT COUNT(*) FROM whitelabel")
	res := w.pool.SendBatch(ctx, &b)
	defer func() {
		err := res.Close()
		if err != nil {
			log.Error().Msgf("whitelabel get all: close batch: %v", err)
		}
	}()
	rows, err := res.Query()
	if err != nil {
		return nil, 0, fmt.Errorf("whitelabel get all: %w", err)
	}
	defer rows.Close()
	whiteLabels := make([]*model.WhiteLabel, 0)
	for rows.Next() {
		var wl model.WhiteLabel

		err = rows.Scan(
			&wl.ID,
			&wl.UserID,
			&wl.SegmentID,
			&wl.Origin,
			&wl.Prefix,
			&wl.SenderEmail,
			&wl.Domain,
			&wl.APIKey,
			&wl.URL,
			&wl.Version,
			&wl.MasterSlave,
			&wl.MasterFee,
			&wl.IsTwoFAEnabled,
			&wl.IsCaptchaEnabled,
			&wl.IsEmailConfirmEnabled,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("whitelabel get all: scan: %w", err)
		}
		whiteLabels = append(whiteLabels, &wl)
	}
	var count int
	err = res.QueryRow().Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("whitelabel get all: scan total count: %w", err)
	}
	return whiteLabels, count, nil
}

func (w *whiteLabel) Update(ctx context.Context, wl *model.WhiteLabel) error {
	tx, err := w.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("whitelabel update: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error().Msgf("whitelabel update: rollback: %v, whitelabel id: %d", err, wl.ID)
		}
	}()
	err = w.update(ctx, tx, wl)
	if err != nil {
		return fmt.Errorf("whitelabel update: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("whitelabel update: commit: %w", err)
	}
	return nil
}

func (w *whiteLabel) update(ctx context.Context, tx pgx.Tx, wl *model.WhiteLabel) error {
	_, err := tx.Exec(ctx, "UPDATE whitelabel SET domain=$1 WHERE id=$2", wl.Domain, wl.ID)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (w *whiteLabel) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := w.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("whitelabel delete: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error().Msgf("whitelabel delete: rollback: %v, whitelabel id: %d", err, id)
		}
	}()
	err = w.delete(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("whitelabel delete: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("whitelabel delete: commit: %w", err)
	}
	return nil
}

func (w *whiteLabel) delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	_, err := tx.Exec(ctx, "DELETE FROM whitelabel WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (w *whiteLabel) GetByID(ctx context.Context, id uuid.UUID) (*model.WhiteLabel, error) {
	const selectByIDQuery = `
		SELECT id, user_id, segment_id, origin, prefix, sender_email, domain, api_key, url, version, master_slave,
		       master_fee, is_two_fa_enabled, is_captcha_enabled, is_email_confirm_enabled
		FROM whitelabel 
		WHERE id=$1`
	var (
		wl model.WhiteLabel
	)

	err := w.pool.QueryRow(ctx, selectByIDQuery, id).Scan(
		&wl.ID, &wl.UserID, &wl.SegmentID, &wl.Origin, &wl.Prefix, &wl.SenderEmail, &wl.Domain, &wl.APIKey, &wl.URL,
		&wl.Version, &wl.MasterSlave, &wl.MasterFee, &wl.IsTwoFAEnabled, &wl.IsCaptchaEnabled, &wl.IsEmailConfirmEnabled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("whitelabel get by id: scan whitelabel: %w", err)
	}
	return &wl, nil
}

const (
	selectBySegmentID = `
		SELECT
			id, user_id, segment_id, origin, prefix, sender_email, domain, api_key, url, version, master_slave,
			master_fee, is_two_fa_enabled, is_captcha_enabled, is_email_confirm_enabled
		FROM
			whitelabel
		WHERE
		    segment_id = $1
		    AND is_active = true`
)

func (w *whiteLabel) GetBySegmentID(ctx context.Context, segmentID int) (*model.WhiteLabel, error) {
	b := pgx.Batch{}
	b.Queue(selectBySegmentID, segmentID)
	res := w.pool.SendBatch(ctx, &b)
	defer func() {
		if err := res.Close(); err != nil {
			log.Error().Msgf("whitelabel get by segment_id: close batch: %v", err)
		}
	}()

	var wl model.WhiteLabel

	if err := res.QueryRow().Scan(
		&wl.ID,
		&wl.UserID,
		&wl.SegmentID,
		&wl.Origin,
		&wl.Prefix,
		&wl.SenderEmail,
		&wl.Domain,
		&wl.APIKey,
		&wl.URL,
		&wl.Version,
		&wl.MasterSlave,
		&wl.MasterFee,
		&wl.IsTwoFAEnabled,
		&wl.IsCaptchaEnabled,
		&wl.IsEmailConfirmEnabled,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("whitelabel get by segment_id: scan whitelabel: %w", err)
	}

	return &wl, nil
}

const (
	selectByUserID = `
		SELECT
			user_id, segment_id, origin, prefix, sender_email, domain, api_key, url, version, master_slave, master_fee,
			id, is_two_fa_enabled, is_captcha_enabled, is_email_confirm_enabled
		FROM
			whitelabel
		WHERE
		    user_id = $1
		    AND is_active = true`
)

func (w *whiteLabel) GetByUserID(ctx context.Context, userID int) (*model.WhiteLabel, error) {
	b := pgx.Batch{}
	b.Queue(selectByUserID, userID)
	res := w.pool.SendBatch(ctx, &b)
	defer func() {
		if err := res.Close(); err != nil {
			log.Error().Msgf("whitelabel get by user_id: close batch: %v", err)
		}
	}()

	var wl model.WhiteLabel

	if err := res.QueryRow().Scan(&wl.UserID, &wl.SegmentID, &wl.Origin, &wl.Prefix, &wl.SenderEmail, &wl.Domain,
		&wl.APIKey, &wl.URL, &wl.Version, &wl.MasterSlave, &wl.MasterFee, &wl.ID, &wl.IsTwoFAEnabled,
		&wl.IsCaptchaEnabled, &wl.IsEmailConfirmEnabled); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("whitelabel get by user_id: scan whitelabel: %w", err)
	}
	return &wl, nil
}

const (
	selectByOrigin = `
		SELECT
			id, user_id, segment_id, origin, prefix, sender_email, domain, api_key, url, version, master_slave,
			master_fee, is_two_fa_enabled, is_captcha_enabled, is_email_confirm_enabled
		FROM
			whitelabel
		WHERE
		  $1 ilike '%' || origin || '%'
		  AND is_active = true`
)

func (w *whiteLabel) GetByOrigin(ctx context.Context, origin string) (*model.WhiteLabel, error) {
	b := pgx.Batch{}
	b.Queue(selectByOrigin, origin)
	res := w.pool.SendBatch(ctx, &b)
	defer func() {
		if err := res.Close(); err != nil {
			log.Error().Msgf("whitelabel get by origin: close batch: %v", err)
		}
	}()

	var wl model.WhiteLabel

	if err := res.QueryRow().Scan(
		&wl.ID, &wl.UserID, &wl.SegmentID, &wl.Origin, &wl.Prefix, &wl.SenderEmail, &wl.Domain, &wl.APIKey, &wl.URL,
		&wl.Version, &wl.MasterSlave, &wl.MasterFee, &wl.IsTwoFAEnabled, &wl.IsCaptchaEnabled,
		&wl.IsEmailConfirmEnabled); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("whitelabel get by origin: scan whitelabel: %w", err)
	}

	return &wl, nil
}

const (
	checkByUserID = `
		SELECT
			count(*)
		FROM
			whitelabel
		WHERE
		    user_id = $1
		    AND is_active = true`
)

func (w *whiteLabel) CheckByUserID(ctx context.Context, userID int) (bool, error) {
	var count int

	if err := w.pool.QueryRow(ctx, checkByUserID, userID).Scan(&count); err != nil {
		return false, err
	}

	return count >= 1, nil
}

const (
	checkByUserIDAndOrigin = `
		SELECT
			count(*)
		FROM
			whitelabel
		WHERE
		    user_id = $1
		  	AND  $2 ilike '%' || origin || '%'
		    AND is_active = true`
)

func (w *whiteLabel) CheckByUserIDAndOrigin(ctx context.Context, userID int, origin string) (bool, error) {
	var count int

	if err := w.pool.QueryRow(ctx, checkByUserIDAndOrigin, userID, origin).Scan(&count); err != nil {
		return false, err
	}

	return count >= 1, nil
}

const (
	selectV2WLs = `
		SELECT
			id, user_id, segment_id, api_key, url
		FROM
			whitelabel
		WHERE
		    version = 2
		    AND is_active = true`
)

func (w *whiteLabel) GetV2WLs(ctx context.Context) ([]*model.WhiteLabel, error) {
	rows, err := w.pool.Query(ctx, selectV2WLs)
	if err != nil {
		return nil, err
	}

	var opts []*model.WhiteLabel

	for rows.Next() {
		var opt model.WhiteLabel

		if err = rows.Scan(&opt.ID, &opt.UserID, &opt.SegmentID, &opt.APIKey, &opt.URL); err != nil {
			return nil, err
		}

		opts = append(opts, &opt)
	}

	return opts, nil
}

const (
	selectWlConfigByOrigin = `
		SELECT
			f.ref_id::varchar, f.origin, f.title, f.commission, f.colors, f.logo, f.favicon, f.firmware_instruction,
			f.lang, f.possible_languages
		FROM
			frontend_configs f
		LEFT JOIN
			whitelabel w ON w.segment_id = f.ref_id
		WHERE
			$1 ilike '%' || w.origin || '%'`
)

func (w *whiteLabel) GetConfigByOrigin(ctx context.Context, origin string) (*model.WlConfig, error) {
	b := pgx.Batch{}
	b.Queue(selectWlConfigByOrigin, origin)
	res := w.pool.SendBatch(ctx, &b)
	defer func() {
		if err := res.Close(); err != nil {
			log.Error().Msgf("whitelabel get by origin: close batch: %v", err)
		}
	}()

	var conf model.WlConfig
	var possibleLang string
	if err := res.QueryRow().Scan(&conf.RefID, &conf.Origin, &conf.Title, &conf.Commission,
		&conf.ColorsJB, &conf.Logo, &conf.Favicon, &conf.FirmwareInstruction, &conf.Lang, &possibleLang); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &conf, nil
		}
		return nil, fmt.Errorf("whitelabel get by origin: scan whitelabel: %w", err)
	}

	if conf.ColorsJB != nil {
		if err := json.Unmarshal(conf.ColorsJB, &conf.Colors); err != nil {
			log.Error().Msg(err.Error())
		}
	}

	var err error

	conf.StratumLists, err = w.GetWLStratumList(ctx, conf.RefID)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	conf.PossibleLang = strings.Split(possibleLang, ",")

	return &conf, nil
}

const (
	configExists = `SELECT COUNT(*) FROM frontend_configs WHERE ref_id = $1`
	insertConfig = `INSERT INTO frontend_configs (
                              ref_id, origin, title, commission, colors, logo, favicon, created_at, updated_at,
                              firmware_instruction, lang) VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now(),$8, $9)`
	updateConfig = `UPDATE frontend_configs
    				SET ref_id = $1, origin = $2, title = $3, commission = $4, colors = $5, logo = $6,
    				    favicon = $7, updated_at = now(),firmware_instruction=$8, lang=$9 WHERE ref_id = $1`
)

func (w *whiteLabel) SetConfigByRefID(ctx context.Context, conf *model.WlConfig) error {
	tx, err := w.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("tx create err: %w", err)
	}

	defer func(tx pgx.Tx) {
		if err = tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Error().Msg(err.Error())
		}
	}(tx)

	var (
		count int
		query = insertConfig
	)

	if err = tx.QueryRow(ctx, configExists, conf.RefID).Scan(&count); err != nil {
		return fmt.Errorf("get count err: %w", err)
	}

	if count > 0 {
		query = updateConfig
	}

	if conf.ColorsJB, err = json.Marshal(conf.Colors); err != nil {
		return fmt.Errorf("json marshal err: %w", err)
	}

	if _, err = tx.Exec(ctx, query, conf.RefID, conf.Origin, conf.Title,
		conf.Commission, conf.ColorsJB, conf.Logo, conf.Favicon, conf.FirmwareInstruction, conf.Lang); err != nil {
		return fmt.Errorf("insert/update err: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit err: %w", err)
	}

	return nil
}

const (
	insertAllowOrigin = `INSERT INTO
    							allow_origins (user_id, origin, created_at)
						 VALUES
						     	($1, $2, now())`
)

func (w *whiteLabel) SetAllowOrigin(ctx context.Context, req *model.AllowOrigin) error {
	if _, err := w.pool.Exec(ctx, insertAllowOrigin, req.UserID, req.Origin); err != nil {
		return fmt.Errorf("insert err: %w", err)
	}

	return nil
}

const (
	selectAllowOrigins = `
		SELECT
			user_id, origin
		FROM
			allow_origins
		WHERE
		    is_active = true`
)

func (w *whiteLabel) GetAllowOrigins(ctx context.Context) ([]*model.AllowOrigin, error) {
	rows, err := w.pool.Query(ctx, selectAllowOrigins)
	if err != nil {
		return nil, err
	}

	var opts []*model.AllowOrigin

	for rows.Next() {
		var opt model.AllowOrigin

		if err = rows.Scan(&opt.UserID, &opt.Origin); err != nil {
			return nil, err
		}

		opts = append(opts, &opt)
	}

	return opts, nil
}

const (
	insertWLStratum = `
			INSERT INTO
				stratum (ref_id, coin, region, number, url, created_at, updated_at)
			VALUES
				($1, $2, $3, $4, $5, now(), now())`
)

func (w *whiteLabel) SetStratum(ctx context.Context, req *model.Stratum) error {
	if _, err := w.pool.Exec(ctx, insertWLStratum, req.RefID, req.Coin, req.Region, req.Number, req.Url); err != nil {
		return fmt.Errorf("insert err: %w", err)
	}

	return nil
}

const (
	selectWLStratumList = `
		SELECT
			coin, region, number, url
		FROM
			stratum
		WHERE
		    ref_id =
		    (CASE
				WHEN (SELECT count(*) FROM stratum WHERE ref_id = $1) > 0 THEN $1
				ELSE 1
			END)`

	selectWLStratumListV2 = `
		SELECT
			coin, region, number, url
		FROM
			stratum
		WHERE
		    ref_id = $1`
)

func (w *whiteLabel) GetWLStratumList(ctx context.Context, refID string) ([]model.Stratum, error) {
	rows, err := w.pool.Query(ctx, selectWLStratumList, refID)
	if err != nil {
		return nil, err
	}

	var list []model.Stratum

	for rows.Next() {
		var stratum model.Stratum

		if err = rows.Scan(&stratum.Coin, &stratum.Region, &stratum.Number, &stratum.Url); err != nil {
			return nil, err
		}

		stratum.RefID = refID

		list = append(list, stratum)
	}

	return list, nil
}

func (w *whiteLabel) GetWLStratumListV2(ctx context.Context, refID int32) ([]model.Stratum, error) {
	rows, err := w.pool.Query(ctx, selectWLStratumListV2, refID)
	if err != nil {
		return nil, err
	}

	var list []model.Stratum

	for rows.Next() {
		var stratum model.Stratum

		if err = rows.Scan(&stratum.Coin, &stratum.Region, &stratum.Number, &stratum.Url); err != nil {
			return nil, err
		}

		stratum.RefID = string(refID)

		list = append(list, stratum)
	}

	return list, nil
}

func (w *whiteLabel) GetFullByUserID(ctx context.Context, userID int) (*model.WhiteLabel, error) {
	var wl model.WhiteLabel

	err := w.pool.QueryRow(ctx, `SELECT id,domain,user_id, segment_id, origin, prefix, sender_email, domain,
       api_key, url, version, master_slave, master_fee, is_two_fa_enabled, is_captcha_enabled, is_email_confirm_enabled
	FROM whitelabel WHERE user_id=$1`, userID).Scan(
		&wl.ID, &wl.Domain, &wl.UserID, &wl.SegmentID, &wl.Origin, &wl.Prefix, &wl.SenderEmail, &wl.Domain,
		&wl.APIKey, &wl.URL, &wl.Version, &wl.MasterSlave, &wl.MasterFee, &wl.IsTwoFAEnabled, &wl.IsCaptchaEnabled,
		&wl.IsEmailConfirmEnabled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get full by user id: %w", err)
	}
	return &wl, nil
}

const selectCoins = `
	SELECT
		wl_id,
		coin_id
	FROM whitelabel_coins where wl_id=$1`

func (w *whiteLabel) GetCoins(ctx context.Context, wlId uuid.UUID) ([]*model.WLCoins, error) {
	rows, err := w.pool.Query(ctx, selectCoins, wlId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	coins := make([]*model.WLCoins, 0)
	for rows.Next() {
		var item model.WLCoins
		if err = rows.Scan(
			&item.WlID,
			&item.CoinID,
		); err != nil {
			return nil, errors.WithStack(err)
		}
		coins = append(coins, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return coins, nil
}

func (w *whiteLabel) AddCoin(ctx context.Context, wlID uuid.UUID, coinID string) error {
	_, err := w.pool.Exec(ctx, `
		INSERT INTO whitelabel_coins (wl_id, coin_id)
		VALUES ($1, $2)`,
		wlID, strings.ToLower(coinID))
	if err != nil {
		return fmt.Errorf("failed to add coin: %w", err)
	}
	return nil
}

func (w *whiteLabel) DeleteCoin(ctx context.Context, wlID uuid.UUID, coinID string) error {
	_, err := w.pool.Exec(ctx, `
		DELETE FROM whitelabel_coins
		WHERE wl_id = $1 AND coin_id = $2`,
		wlID, strings.ToLower(coinID))
	if err != nil {
		return fmt.Errorf("failed to delete coin: %w", err)
	}
	return nil
}
