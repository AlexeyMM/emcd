package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	userAccountModelEnum "code.emcdtech.com/emcd/service/accounting/model/enum"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/emcd/sdk/log"
	pgTx "code.emcdtech.com/emcd/sdk/pg"

	"code.emcdtech.com/emcd/service/profile/internal/model"
)

type OldUsers interface {
	Create(ctx context.Context, u *model.User) (int32, uuid.UUID, error)
	// CreateUser only for create
	CreateUser(ctx context.Context, u *model.User) (int32, error)
	SaveV2(ctx context.Context, u *model.User) (int32, error)
	UpdatePassword(ctx context.Context, email, password string) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetIDs(ctx context.Context, ids []uuid.UUID) ([]*model.ID, error)
	GetMiningAccountID(ctx context.Context, userID, coinID int) (int, error)
	GetDefaultAddress(ctx context.Context, miningAccountID int) (string, error)
	GetRefIDAndEmail(ctx context.Context, userID int) (int32, string, error)
	UpdateUserMinPay(ctx context.Context, miningAccountID int, minpay float32) error
	GetSegmentID(ctx context.Context, userID int32) (int, error)
	Begin() pgTx.PgxTransactor
	CreateAddress(ctx context.Context, userID int32, coinID int, address string) error
	DeleteAddress(ctx context.Context, userID int32, coinID int) error
	UpdateUserAccountChangedAt(ctx context.Context, userID int32, coinID int) error
	GetEmailWithParentID(ctx context.Context, userID int32) (string, error)
	SaveID(ctx context.Context, id uuid.UUID, email string) (*model.ID, error)
	UpdateRefID(ctx context.Context, userID int32, refID int32) error
	UpdateUser(ctx context.Context, u *model.User) error
	SetTimezone(ctx context.Context, userID int, timezone string) error
	SetLanguage(ctx context.Context, userID int, language string) error
	GetAddresses(ctx context.Context, userID uuid.UUID, coinCodesMining []string) ([]*model.Address, error)

	/* API key's section */

	// GetAPIKeyByUserIDAndParentID get API key by userID and parentID
	GetAPIKeyByUserIDAndParentID(ctx context.Context, userID, parentID int32) (bool, string, error)
	// SetAPIKeyForUserIDAndParentID set API key for pair userID and parentID
	SetAPIKeyForUserIDAndParentID(ctx context.Context, apiKey string, userID, parentID int32) error
	// SetAPIKeyForUserID set API key for userID
	SetAPIKeyForUserID(ctx context.Context, apiKey string, userID int32) error

	/* No Pay */

	// UpdateNoPay update no pay field in old user table
	UpdateNoPay(ctx context.Context, userID int32, val bool) error
	GetNoPay(ctx context.Context, userID uuid.UUID) (bool, error)
	SoftDeleteSubUser(ctx context.Context, subUserID int32, newParentID int32, newEmail model.SubUserEmail) error
	GetUserByOldID(ctx context.Context, userID int) (*model.User, error)
	GetUserByUUID(ctx context.Context, userID uuid.UUID) (*model.User, error)

	GetSpecRefFee(ctx context.Context, userID int32) (float64, error)
	ApplyParentPromoCodes(ctx context.Context, parentID int, subID int32) error
	CreateReferralAccounts(ctx context.Context, tx pgx.Tx, userAccountID, coinID int32) error
	CreateAutoPayAddress(ctx context.Context, tx pgx.Tx, userAccountID int32, coinID int32, address string) error
	CreateAccountsPool(ctx context.Context, tx pgx.Tx, userAccountID int32) error

	GetUserIDBySegmentID(ctx context.Context, segmentID int32) (int32, error)
}

type oldUsers struct {
	trx pgTx.PgxTransactor
}

func NewOldUsers(trx pgTx.PgxTransactor) *oldUsers {
	return &oldUsers{
		trx: trx,
	}
}

func (s *oldUsers) GetUserIDBySegmentID(ctx context.Context, segmentID int32) (int32, error) {
	query := `SELECT user_id FROM emcd.histories.segment_userids WHERE segment_id=$1`
	var userID int32
	err := s.trx.Runner(ctx).QueryRow(ctx, query, segmentID).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("query row: scan: %w", err)
	}
	return userID, nil
}

func (s *oldUsers) SaveID(ctx context.Context, id uuid.UUID, email string) (*model.ID, error) {
	log.Error(ctx, "fix user %s %s", id, email)
	email = strings.ToLower(email)
	_, err := s.trx.Runner(ctx).Exec(ctx, "UPDATE emcd.users SET new_id=$1, updated_at = NOW() WHERE email=$2", id, email)
	if err != nil {
		return nil, fmt.Errorf("old users: update password: %w", err)
	}
	u, err := s.GetUserByUUID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.ID{Old: u.OldID, New: u.ID}, nil
}

func (s *oldUsers) Create(ctx context.Context, u *model.User) (int32, uuid.UUID, error) {
	var (
		id    uuid.NullUUID
		oldID int32
	)
	err := s.trx.Runner(ctx).
		QueryRow(ctx, `INSERT INTO emcd.users (new_id,username,ref_id,email,password,created_at,api_key) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (LOWER(email))
DO UPDATE SET email=EXCLUDED.email RETURNING new_id, id`,
			u.ID, u.Username, u.RefID, strings.ToLower(u.Email), u.Password, u.CreatedAt, u.ApiKey).
		Scan(&id, &oldID)
	if err != nil {
		return 0, uuid.Nil, fmt.Errorf("old users: create: %w", err)
	}
	if !id.Valid {
		return 0, uuid.Nil, fmt.Errorf("old users: create: email: %s: %w", u.Email, err)
	}

	return oldID, id.UUID, nil
}

func (s *oldUsers) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var (
		u           model.User
		apiKey      sql.NullString
		suspended   sql.NullTime
		oldParentID sql.NullInt32
	)
	err := s.trx.Runner(ctx).
		QueryRow(ctx, "SELECT id,new_id,username,ref_id,LOWER(email),password,created_at,api_key,suspended,parent_id FROM emcd.users WHERE LOWER(email)=LOWER($1)", email).
		Scan(&u.OldID, &u.ID, &u.Username, &u.RefID, &u.Email, &u.Password, &u.CreatedAt, &apiKey, &suspended, &oldParentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan: %w", err)
	}
	if apiKey.Valid {
		u.ApiKey = apiKey.String
	}
	if suspended.Valid {
		u.Suspended = true
	}
	segmentID, err := s.getSegmentIDByUserID(ctx, u.OldID)
	if err != nil {
		return nil, fmt.Errorf("getSegmentIDByUserID: %w", err)
	}

	u.SegmentID = int(segmentID)
	vip, err := s.isVipUser(ctx, u.OldID)
	if err != nil {
		return nil, fmt.Errorf("isVipUser: %w", err)
	}
	u.Vip = vip
	if oldParentID.Valid {
		err = s.trx.Runner(ctx).QueryRow(ctx, `SELECT new_id FROM emcd.users WHERE id = $1`, oldParentID.Int32).Scan(&u.ParentID)
		if err != nil {
			return nil, fmt.Errorf("queryRow oldParentID: %w", err)
		}

	}

	return &u, nil
}

func (s *oldUsers) getSegmentIDByUserID(ctx context.Context, userID int32) (int32, error) {
	var segmentID int32
	if err := s.trx.Runner(ctx).QueryRow(ctx, `SELECT segment_id FROM histories.segment_userids WHERE user_id=$1`, userID).Scan(&segmentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error(ctx, "not found segment_id by user_id: %d", userID)
			segmentID, err = s.addSegment(ctx, userID)
			if err != nil {
				return 0, fmt.Errorf("addSegment: %w", err)
			}
		} else {
			return 0, fmt.Errorf("get segment id by user id: %w", err)
		}
	}
	return segmentID, nil
}

func (s *oldUsers) addSegment(ctx context.Context, userID int32) (int32, error) {
	query := `INSERT INTO histories.segment_userids (user_id) VALUES ($1) ON CONFLICT DO NOTHING `
	_, err := s.trx.Runner(ctx).Exec(ctx, query, userID)
	if err != nil {
		return 0, fmt.Errorf("exec: %w", err)
	}

	var segmentID int32
	query = `--antislave
		SELECT segment_id FROM histories.segment_userids WHERE user_id = $1`
	err = s.trx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&segmentID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("queryRow: %w", err)
	}

	return segmentID, nil
}

func (s *oldUsers) isVipUser(ctx context.Context, userID int32) (bool, error) {
	var count int
	err := s.trx.Runner(ctx).QueryRow(ctx, "SELECT COUNT(*) FROM emcd.vip_users WHERE user_id=$1", userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("is vip user: %w", err)
	}
	return count != 0, nil
}

func (s *oldUsers) UpdatePassword(ctx context.Context, email, password string) error {
	email = strings.ToLower(email)
	_, err := s.trx.Runner(ctx).
		Exec(ctx, "UPDATE emcd.users SET password=$1, updated_at = NOW() WHERE LOWER(email)=$2", password, email)
	if err != nil {
		return fmt.Errorf("old users: update password: %w", err)
	}
	return nil
}

// CreateUser only for create
func (s *oldUsers) CreateUser(ctx context.Context, u *model.User) (int32, error) {
	var id int32

	queryCreateUser := `
INSERT 
	INTO emcd.users (new_id,
	                 username,
	                 ref_id,
	                 email,
	                 password,
	                 created_at,
	                 api_key,
	                 is_active,
	                 parent_id,
	                 language) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`

	err := s.trx.Runner(ctx).
		QueryRow(ctx, queryCreateUser,
			u.ID,
			u.Username,
			u.RefID,
			strings.ToLower(u.Email),
			u.Password,
			u.CreatedAt,
			u.ApiKey,
			u.IsActive,
			u.OldParentID,
			u.Language).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("old users: save v2: %w", err)
	}
	return id, nil
}

func (s *oldUsers) SaveV2(ctx context.Context, u *model.User) (int32, error) {
	var id int32
	u.Email = strings.ToLower(u.Email)

	err := s.trx.Runner(ctx).
		QueryRow(ctx, `INSERT INTO emcd.users (new_id,username,ref_id,email,password,created_at,api_key,is_active,parent_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT (LOWER(email))
			DO UPDATE SET new_id=$1,username=$2,ref_id=$3,password=$5,created_at=$6,api_key=$7,is_active=$8 RETURNING id`,
			u.ID, u.Username, u.RefID, u.Email, u.Password, u.CreatedAt, u.ApiKey, u.IsActive, u.OldParentID).
		Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("old users: save v2: %w", err)
	}
	return id, nil
}

func (s *oldUsers) GetIDs(ctx context.Context, ids []uuid.UUID) ([]*model.ID, error) {
	getByIDsQuery := `SELECT id,new_id FROM emcd.users WHERE new_id=ANY($1::uuid[])`
	rows, err := s.trx.Runner(ctx).Query(ctx, getByIDsQuery, ids)
	if err != nil {
		return nil, fmt.Errorf("old users: get by ids: %w", err)
	}
	defer rows.Close()
	res := make([]*model.ID, 0)
	for rows.Next() {
		var id model.ID
		err = rows.Scan(&id.Old, &id.New)
		if err != nil {
			return nil, fmt.Errorf("old users: get by ids: %w", err)
		}
		res = append(res, &id)
	}
	return res, nil
}

func (s *oldUsers) GetMiningAccountID(ctx context.Context, userID, coinID int) (int, error) {
	query := `SELECT id FROM emcd.users_accounts WHERE account_type_id = 2 AND user_id = $1 AND coin_id = $2`

	var miningAccountID int
	if err := s.trx.Runner(ctx).QueryRow(ctx, query, userID, coinID).Scan(&miningAccountID); err != nil {
		return 0, fmt.Errorf("s.pool.QueryRow: %w", err)
	}

	return miningAccountID, nil
}

func (s *oldUsers) GetDefaultAddress(ctx context.Context, miningAccountID int) (string, error) {
	query := `
		SELECT COALESCE(address, '')
		FROM emcd.autopay_addresses
		WHERE user_account_id = $1
		  AND percent = 100`

	var address string
	err := s.trx.Runner(ctx).QueryRow(ctx, query, miningAccountID).Scan(&address)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("queryRow: %w", err)
	}

	return address, nil
}

func (s *oldUsers) GetRefIDAndEmail(ctx context.Context, userID int) (int32, string, error) {
	query := `SELECT ref_id, LOWER(email) FROM emcd.users WHERE id = $1`
	var (
		refID sql.NullInt32
		email sql.NullString
	)
	err := s.trx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&refID, &email)
	if err != nil {
		return 0, "", fmt.Errorf("queryRow: %w", err)
	}
	if !refID.Valid || !email.Valid {
		return 0, "", fmt.Errorf("refID or email are invalid")
	}
	return refID.Int32, email.String, nil
}

func (s *oldUsers) UpdateUserMinPay(ctx context.Context, miningAccountID int, minpay float32) error {
	mp := decimal.NewFromFloat32(minpay)
	query := `UPDATE emcd.users_accounts SET minpay = $1 WHERE id = $2`
	if _, err := s.trx.Runner(ctx).Exec(ctx, query, mp, miningAccountID); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (s *oldUsers) GetSegmentID(ctx context.Context, userID int32) (int, error) {
	query := `SELECT segment_id FROM histories.segment_userids WHERE user_id=$1`
	var segmentID int32
	err := s.trx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&segmentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error(ctx, "not found segment_id by user_id: %d", userID)
			segmentID, err = s.addSegment(ctx, userID)
			if err != nil {
				return 0, fmt.Errorf("addSegment: %w", err)
			}
		} else {
			return 0, fmt.Errorf("queryRow: %w", err)
		}
	}
	return int(segmentID), nil
}

func (s *oldUsers) Begin() pgTx.PgxTransactor {
	return s.trx
}

func (s *oldUsers) DeleteAddress(ctx context.Context, userID int32, coinID int) error {
	query := `DELETE FROM emcd.autopay_addresses WHERE user_account_id =
    (SELECT id FROM emcd.users_accounts WHERE user_id = $1 AND coin_id = $2 AND account_type_id = 2)`
	_, err := s.trx.Runner(ctx).Exec(ctx, query, userID, coinID)
	if err != nil {
		return fmt.Errorf("old users: delete address: %w", err)
	}
	return nil
}

func (s *oldUsers) CreateAddress(ctx context.Context, userID int32, coinID int, address string) error {
	query := `INSERT INTO emcd.autopay_addresses (user_account_id, address, percent)
VALUES ((SELECT id FROM emcd.users_accounts WHERE user_id = $1 AND coin_id = $2 AND account_type_id = 2), $3, 100);`
	_, err := s.trx.Runner(ctx).Exec(ctx, query, userID, coinID, address)
	if err != nil {
		return fmt.Errorf("old users: create address: %w", err)
	}
	return nil
}

func (s *oldUsers) UpdateUserAccountChangedAt(ctx context.Context, userID int32, coinID int) error {
	query := `UPDATE emcd.users_accounts SET changed_at=NOW() WHERE user_id = $1 AND coin_id = $2 AND account_type_id = 2`
	_, err := s.trx.Runner(ctx).Exec(ctx, query, userID, coinID)
	if err != nil {
		return fmt.Errorf("update user account changed at: %w", err)
	}
	return nil
}

func (s *oldUsers) GetEmailWithParentID(ctx context.Context, userID int32) (string, error) {
	query := `
	SELECT
		CASE
			WHEN us.parent_id IS NULL THEN LOWER(us.email)
			ELSE LOWER(up.email)
		END
	FROM
		emcd.users us
	LEFT JOIN
		emcd.users up ON up.id = us.parent_id
	WHERE
		us.id =$1`
	var e sql.NullString
	err := s.trx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&e)
	if err != nil {
		return "", fmt.Errorf("old users: get email: to notify: %w", err)
	}
	var res string
	if e.Valid {
		res = e.String
	}
	return res, nil
}

func (s *oldUsers) GetNoPay(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `SELECT nopay FROM emcd.users WHERE new_id = $1`
	nopay := true
	err := s.trx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&nopay)
	if err != nil {
		return false, fmt.Errorf("queryRow: %w", err)
	}
	return nopay, nil
}

func (s *oldUsers) UpdateNoPay(ctx context.Context, userID int32, val bool) error {
	query := `UPDATE emcd.users SET nopay = $1 WHERE id = $2 OR parent_id = $2;`
	commandTag, err := s.trx.Runner(ctx).Exec(ctx, query, val, userID)
	switch {
	case err != nil:
		return fmt.Errorf("old users: update nopay: %w", err)
	case commandTag.RowsAffected() == 0:
		return fmt.Errorf("old users: update nopay: unknown user_id: %d", userID)
	}
	return nil
}

func (p *oldUsers) UpdateRefID(ctx context.Context, userID int32, refID int32) error {
	query := `UPDATE emcd.users SET ref_id=$1 WHERE id=$2 OR parent_id=$2`

	commandTag, err := p.trx.Runner(ctx).Exec(ctx, query, refID, userID)
	if err != nil {
		return fmt.Errorf("old users: update ref_id: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("old users: update ref_id: unknown user_id: %d", userID)
	}
	return nil
}

func (p *oldUsers) UpdateUser(ctx context.Context, u *model.User) error {
	updateUserQuery := "UPDATE emcd.users SET username=$1,ref_id=$2,email=$3,api_key=$4,is_active=$5 WHERE new_id=$6"
	_, err := p.trx.Runner(ctx).Exec(ctx, updateUserQuery, u.Username, u.RefID, u.Email, u.ApiKey, u.IsActive, u.ID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (p *oldUsers) SetTimezone(ctx context.Context, userID int, timezone string) error {
	query := `UPDATE emcd.users SET timezone = $1 WHERE id = $2`
	_, err := p.trx.Runner(ctx).Exec(ctx, query, timezone, userID)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (p *oldUsers) SetLanguage(ctx context.Context, userID int, language string) error {
	query := `UPDATE emcd.users SET language = $1 WHERE id = $2`
	_, err := p.trx.Runner(ctx).Exec(ctx, query, language, userID)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (p *oldUsers) GetAddresses(ctx context.Context, userUuid uuid.UUID, coinCodesMining []string) ([]*model.Address, error) {
	query := `SELECT
			c.code AS coin,
			ua.account_type_id,
			ua.minpay,
			COALESCE(a.address, COALESCE(ua.address, '')) AS wallet_address,
			COALESCE(aa.address, '') AS mining_address
		FROM
			emcd.coins c
		LEFT JOIN
			emcd.users_accounts ua ON ua.coin_new = c.code AND ua.account_type_id in (1,2)
		LEFT JOIN
			emcd.users u ON u.id = ua.user_id /* bug#1 after fix uuid.New() in users_accounts, start use u.new_id = ua.user_id_new; bug#2 (some users with same new_id has different id) start use emcd.users u ON u.new_id = ua.user_id_new and u.id = ua.user_id */
		LEFT JOIN
			emcd.autopay_addresses aa ON ua.id = aa.user_account_id AND aa.percent = 100
		LEFT JOIN
			emcd.addresses a ON a.user_account_id = ua.id and token_id is null /* token_id is null means is native network (deprecated soon) */
		WHERE
			u.new_id = @user_uuid AND ua.account_type_id in (1,2) and c.code = ANY(@coin_codes)`

	namedArgs := pgx.NamedArgs{
		"user_uuid":  userUuid,
		"coin_codes": pq.Array(coinCodesMining),
	}

	rows, err := p.trx.Runner(ctx).Query(ctx, query, namedArgs)
	if err != nil {
		return nil, fmt.Errorf("failed query oldUsers.GetAddresses: %w", err)
	}
	defer rows.Close()

	addrMap := make(map[string]*model.Address)
	for rows.Next() {
		var addr model.Address

		if err = rows.Scan(
			&addr.Coin,
			&addr.AccountTypeId,
			&addr.MinPay,
			&addr.WalletAddress,
			&addr.MiningAddress,
		); err != nil {
			return nil, fmt.Errorf("failed scan oldUsers.GetAddresses: %w", err)
		}

		if _, ok := addrMap[addr.Coin]; !ok {
			addrMap[addr.Coin] = &addr
		} else {
			addrMap[addr.Coin].MinPay = decimal.Max(addrMap[addr.Coin].MinPay, addr.MinPay) // MinPay can be mess

			if addr.AccountTypeId.AccountTypeId == userAccountModelEnum.WalletAccountTypeID {
				addrMap[addr.Coin].WalletAddress = addr.WalletAddress
			} else if addr.AccountTypeId.AccountTypeId == userAccountModelEnum.MiningAccountTypeID {
				addrMap[addr.Coin].MiningAddress = addr.MiningAddress
			} else {
				// pass
			}
		}
	}

	var addrList []*model.Address
	for _, addr := range addrMap {
		if addr.MiningAddress == "" {
			addr.MiningAddress = addr.WalletAddress
		}

		addrList = append(addrList, addr)

	}

	return addrList, nil
}

func (p *oldUsers) GetAPIKeyByUserIDAndParentID(ctx context.Context, userID, parentID int32) (bool, string, error) {
	query := `SELECT api_key FROM emcd.users WHERE id = $1 AND parent_id = $2`
	var apiKey sql.NullString

	err := p.trx.Runner(ctx).QueryRow(ctx, query, userID, parentID).Scan(&apiKey)

	switch {
	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		return false, "", fmt.Errorf("profile: get api key by pair user_id and parent_id: %w", err)
	case errors.Is(err, pgx.ErrNoRows):
		return false, "", nil
	}

	return apiKey.String != "", apiKey.String, nil
}

func (p *oldUsers) SetAPIKeyForUserIDAndParentID(ctx context.Context, apiKey string, userID, parentID int32) error {
	query := `UPDATE emcd.users SET api_key=$1 WHERE id=$2 AND parent_id=$3`
	commandTag, err := p.trx.Runner(ctx).Exec(ctx, query, apiKey, userID, parentID)
	if err != nil {
		return fmt.Errorf("profile: update api_key by user_id and parent_id: %d, %d: %w", userID, parentID, err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("profile: update api_key: unknown pair user_id and parent_id: %d, %d", userID, parentID)
	}
	return nil
}

func (p *oldUsers) SetAPIKeyForUserID(ctx context.Context, apiKey string, userID int32) error {
	query := `UPDATE emcd.users SET api_key=$1 WHERE id=$2`
	commandTag, err := p.trx.Runner(ctx).Exec(ctx, query, apiKey, userID)
	if err != nil {
		return fmt.Errorf("profile: update api_key by user_id: %d: %w", userID, err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("profile: update api_key: unknown user_id: %d", userID)
	}
	return nil
}

func (p *oldUsers) SoftDeleteSubUser(ctx context.Context, subUserID int32, newParentID int32, newEmail model.SubUserEmail) error {
	const op = "repository.old_users.SoftDeleteSubUser"

	err := p.markSubUserAsDeleted(ctx, newParentID, newEmail, subUserID)
	if err != nil {
		return fmt.Errorf("%s: delete sub user: %w", op, err)
	}

	return nil
}

func (s *oldUsers) GetUserByOldID(ctx context.Context, userID int) (*model.User, error) {
	query := `--antislave
		SELECT id,new_id,username,ref_id,LOWER(email),password,created_at,api_key,language FROM emcd.users WHERE id=$1`
	var u model.User
	var apiKey sql.NullString
	var language sql.NullString
	err := s.trx.Runner(ctx).QueryRow(ctx, query, userID).
		Scan(&u.OldID, &u.ID, &u.Username, &u.RefID, &u.Email, &u.Password, &u.CreatedAt, &apiKey, &language)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("old users: get by old id: %w", err)
	}
	if apiKey.Valid {
		u.ApiKey = apiKey.String
	}
	u.Language = "en"
	if language.Valid {
		u.Language = language.String
	}
	return &u, nil
}

func (p *oldUsers) markSubUserAsDeleted(ctx context.Context, newParentID int32, newEmail model.SubUserEmail, subUserID int32) error {
	const op = "repository.old_users.deleteSubUser"
	softDeleteQuery := `
	UPDATE emcd.users SET
		parent_id=@newParentID,
		username=CONCAT('deleted_',username),
		password='subaccount',
		email=@newEmail
	WHERE id = @subUserID`
	if _, err := p.trx.Runner(ctx).Exec(ctx, softDeleteQuery, pgx.NamedArgs{
		"newParentID": newParentID,
		"newEmail":    newEmail,
		"subUserID":   subUserID,
	}); err != nil {
		return fmt.Errorf("%s: soft_delete: %w", op, err)
	}
	return nil
}

func (s *oldUsers) GetUserByUUID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	query := `SELECT id,new_id,username,ref_id,LOWER(email),password,created_at,api_key,language FROM emcd.users WHERE new_id=$1`
	var u model.User
	var apiKey sql.NullString
	var language sql.NullString
	err := s.trx.Runner(ctx).QueryRow(ctx, query, userID).
		Scan(&u.OldID, &u.ID, &u.Username, &u.RefID, &u.Email, &u.Password, &u.CreatedAt, &apiKey, &language)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("old users: get by uuid: %w", err)
	}
	if apiKey.Valid {
		u.ApiKey = apiKey.String
	}
	u.Language = "en"
	if language.Valid {
		u.Language = language.String
	}
	return &u, nil
}

func (s *oldUsers) GetSpecRefFee(ctx context.Context, userID int32) (float64, error) {
	var specRefFee float64
	query := `SELECT spec_ref_fee FROM emcd.vip_users WHERE user_id = $1 AND is_blocked_ref_tier=true`
	if err := s.trx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&specRefFee); err != nil {
		return 0, fmt.Errorf("queryRow: %w", err)
	}

	return specRefFee, nil
}

func (s *oldUsers) ApplyParentPromoCodes(ctx context.Context, parentID int, subID int32) error {
	selectParensPromoCodes := `
		SELECT
			promocode_id, created_at, expires_at
		FROM
			emcd.users_promocodes
		WHERE
			user_id = $1`

	insertUserPromoCode := `
		INSERT INTO emcd.users_promocodes
			(user_id, promocode_id, created_at, expires_at)
		VALUES
			($1, $2, $3, $4)`

	rows, err := s.trx.Runner(ctx).Query(ctx, selectParensPromoCodes, parentID)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var codes []*model.UserPromoCode
	for rows.Next() {
		var p model.UserPromoCode
		if err = rows.Scan(&p.ID, &p.CreatedAt, &p.ExpiresAt); err != nil {
			return fmt.Errorf("scan: %w", err)
		}
		codes = append(codes, &p)
	}

	for _, p := range codes {
		if _, err = s.trx.Runner(ctx).Exec(ctx, insertUserPromoCode, subID, p.ID, p.CreatedAt, p.ExpiresAt); err != nil {
			return fmt.Errorf("exec: %w", err)
		}
	}

	return nil
}

func (s *oldUsers) CreateReferralAccounts(ctx context.Context, tx pgx.Tx, userAccountID, coinID int32) error {
	fee := 0.0005
	query := `INSERT INTO emcd.accounts_referral (account_id, tier, referral_fee, active_referrals, coin_id)
	VALUES ($1, 1, $2, 0, $3)`
	_, err := tx.Exec(ctx, query, userAccountID, fee, coinID)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (s *oldUsers) CreateAutoPayAddress(ctx context.Context, tx pgx.Tx, userAccountID int32, coinID int32, address string) error {
	query := `INSERT INTO emcd.autopay_addresses (user_account_id, address, percent) VALUES ($1, $2, 100);`
	_, err := tx.Exec(ctx, query, userAccountID, address)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (s *oldUsers) CreateAccountsPool(ctx context.Context, tx pgx.Tx, userAccountID int32) error {
	query := `INSERT INTO emcd.accounts_pool (account_id, emcd_address_autopay) VALUES ($1, TRUE)`
	_, err := tx.Exec(ctx, query, userAccountID)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
