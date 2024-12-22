package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"

	"code.emcdtech.com/emcd/sdk/log"
	pgTx "code.emcdtech.com/emcd/sdk/pg"

	"code.emcdtech.com/emcd/service/profile/internal/model"
)

const (
	searchPatternStartsWith = `%s%%`
)

type Profile interface {
	Create(ctx context.Context, pr *model.Profile) error
	SaveV3(ctx context.Context, pr *model.Profile) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error)
	UpdatePassword(ctx context.Context, email, password string) error
	GetUserByEmailAndWl(ctx context.Context, email string, whiteLabelID uuid.UUID) (*model.User, error)
	GetAllUsers(ctx context.Context, skip, take int, sortField string, asc bool, searchEmail string) ([]*model.User, int, error)
	GetAllUsersByWlID(
		ctx context.Context,
		skip, take int,
		sortField string,
		asc bool,
		searchEmail string,
		whiteLabelID uuid.UUID,
	) ([]*model.User, int, error)
	GetUserByApiKey(ctx context.Context, apiKey string) (*model.User, error)
	UpdateUserIsActive(ctx context.Context, email string, active bool) error
	GetUserIsActive(ctx context.Context, userID uuid.UUID) (bool, error)
	ExistAppleID(ctx context.Context, appleID string) (bool, string, error)
	UpdateAppleID(ctx context.Context, email, appleID string) (int, error)
	RelatedUsers(ctx context.Context, firstID, secondID uuid.UUID) (bool, error)
	GetAllSubUsers(ctx context.Context, userID uuid.UUID) ([]*model.User, error)
	GetAllUserIDsByUsername(ctx context.Context) (map[string]uuid.UUID, error)
	GetReferrals(ctx context.Context, userID uuid.UUID, skip, take int, sortField string, asc bool) ([]*model.User, int, error)
	GetUsernamesByIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]string, error)
	GetSuspended(ctx context.Context, userID uuid.UUID) (bool, error)
	SetSuspended(ctx context.Context, userIDs []uuid.UUID, suspended bool) error
	GetEmailsByIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]string, error)
	SafeDeleteByID(ctx context.Context, userID uuid.UUID) error
	SaveUser(ctx context.Context, u *model.User) error
	GetByUsernames(ctx context.Context, usernames []string) ([]*model.Profile, error)
	UpdateRefID(ctx context.Context, userID uuid.UUID, refID int32, newRefID uuid.UUID) error
	GetFlagReferralLinkGenerated(ctx context.Context, userID uuid.UUID) (bool, error)
	SetFlagReferralLinkGenerated(ctx context.Context, userID uuid.UUID, value bool) error
	GetAllReferrals(
		ctx context.Context,
		referrerID uuid.UUID,
		retrieveAfterUsername *string,
		sortField string,
		asc bool,
		take int,
	) (map[string]*model.Referral, []string, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, mutate func(user *model.User) error) error
	SetLanguage(ctx context.Context, userID uuid.UUID, language string) error
	SoftDeleteSubUser(ctx context.Context, subUserID uuid.UUID, newParentID uuid.UUID, newEmail model.SubUserEmail) error

	GetCountUserWithWL(ctx context.Context, wlUUID uuid.UUID, offset, limit int32) ([]model.UserShortInfo, int64, error)

	GetUsersByUUIDs(ctx context.Context, userUUIDs []uuid.UUID) ([]model.UserShortInfo, error)

	Begin() pgTx.PgxTransactor

	/* API key's section */
	// use it only then we migrate all API keys to profile.user table
	// GetAPIKeyByUserIDAndParentID get API key by userID and parentID
	// GetAPIKeyByUserIDAndParentID(ctx context.Context, userID, parentID uuid.UUID) (bool, string, error)
	// GetAPIKeyByUserID get API key by userID
	// GetAPIKeyByUserID(ctx context.Context, userID uuid.UUID) (bool, string, error)
	// SetAPIKeyForUserIDAndParentID set API key for pair userID and parentID
	// SetAPIKeyForUserIDAndParentID(ctx context.Context, apiKey string, userID, parentID uuid.UUID) error
	// SetAPIKeyForUserID set API key for userID
	// SetAPIKeyForUserID(ctx context.Context, apiKey string, userID uuid.UUID) error

	GetByUsernamesForReferrals(ctx context.Context, usernames []string) ([]*model.Profile, error)
	GetLastSubUserEmail(ctx context.Context, parentID uuid.UUID) (string, error)
	IsUniqueUsername(ctx context.Context, username string) (bool, error)
	GetUserByTg(ctx context.Context, tgID string) (*model.User, error)

	SetSecretKeyForUserID(ctx context.Context, secretKey string, userID uuid.UUID) error
}

type profile struct {
	trx    pgTx.PgxTransactor
	oldTrx pgTx.PgxTransactor
}

func NewProfile(trx pgTx.PgxTransactor, oldTrx pgTx.PgxTransactor) *profile {
	return &profile{
		trx:    trx,
		oldTrx: oldTrx,
	}
}

func (p *profile) Create(ctx context.Context, pr *model.Profile) error {
	err := p.createUser(ctx, pr.User)
	if err != nil {
		return fmt.Errorf("profile: create: %w", err)
	}
	return nil
}

func (p *profile) createUser(ctx context.Context, u *model.User) error {
	createUserQuery := `
INSERT INTO users (id,
                   username,
                   ref_id,
                   email,
                   password,
                   created_at,
                   whitelabel_id,
                   api_key,
                   new_ref_id,
                   pool_type,
                   parent_id,
                   language) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := p.trx.Runner(ctx).Exec(
		ctx,
		createUserQuery,
		u.ID,
		u.Username,
		u.RefID,
		strings.ToLower(u.Email),
		u.Password,
		u.CreatedAt,
		u.WhiteLabelID,
		u.ApiKey,
		u.NewRefID,
		u.PoolType,
		u.ParentID,
		u.Language,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (p *profile) SaveV3(ctx context.Context, pr *model.Profile) error {
	const op = "repo.profile.SaveV3"

	err := p.saveUserV2(ctx, pr.User)
	if err != nil {
		return fmt.Errorf("%s: save user: %w", op, err)
	}
	return nil
}

func (p *profile) saveUserV2(ctx context.Context, u *model.User) error {
	var appleID sql.NullString
	if u.AppleID != "" {
		appleID.Valid = true
		appleID.String = u.AppleID
	}

	upsertUserQuery := "INSERT INTO users (id,username,ref_id,email,password,created_at,whitelabel_id,api_key,is_active,apple_id,parent_id,new_ref_id,pool_type,language) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) ON CONFLICT (email,whitelabel_id) DO " +
		"UPDATE SET username=$2,ref_id=$3,password=$5,created_at=$6,api_key=$8,is_active=$9,apple_id=$10,parent_id=$11,new_ref_id=$12,pool_type=$13,language=$14"
	_, err := p.trx.Runner(ctx).
		Exec(ctx, upsertUserQuery, u.ID, u.Username, u.RefID, u.Email, u.Password, u.CreatedAt, u.WhiteLabelID, u.ApiKey, u.IsActive, appleID, u.ParentID, u.NewRefID, u.PoolType, u.Language)
	if err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	return nil
}

func (p *profile) UpdatePassword(ctx context.Context, email, password string) error {
	log.Info(ctx, "update password %s", email)

	email = strings.ToLower(email)
	updatePasswordQuery := "UPDATE users SET password=$1 WHERE LOWER(email)=$2"
	commandTag, err := p.trx.Runner(ctx).Exec(ctx, updatePasswordQuery, password, email)
	if err != nil {
		return fmt.Errorf("profile: update password: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("profile: update password: unknown email: %s", email)
	}
	return nil
}

func (p *profile) GetUserByEmailAndWl(ctx context.Context, email string, whiteLabelID uuid.UUID) (*model.User, error) {
	getUserByEmailQuery := `SELECT id,username,ref_id,email,password,created_at,whitelabel_id,api_key,parent_id,is_active,new_ref_id,suspended,telegram_id,telegram_username
		FROM users WHERE LOWER(email)=LOWER($1) AND whitelabel_id=$2`
	var (
		u          model.User
		pwd        []byte
		isActive   sql.NullBool
		suspended  sql.NullTime
		tgID       sql.NullString
		tgUsername sql.NullString
	)
	err := p.trx.Runner(ctx).QueryRow(ctx, getUserByEmailQuery, email, whiteLabelID).Scan(&u.ID, &u.Username, &u.RefID, &u.Email,
		&pwd, &u.CreatedAt, &u.WhiteLabelID, &u.ApiKey, &u.ParentID, &isActive, &u.NewRefID, &suspended, &tgID, &tgUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("queryRow1: %w", err)
	}
	u.Password = string(pwd)
	if isActive.Valid {
		u.IsActive = isActive.Bool
	}
	if suspended.Valid {
		u.Suspended = true
	}
	if tgID.Valid {
		u.TgID = tgID.String
	}
	if tgUsername.Valid {
		u.TgUsername = tgUsername.String
	}

	query := `--anti slave
		SELECT
		CASE
			WHEN v.user_id IS NOT NULL THEN true
			ELSE false
		END AS vip_user,
			su.segment_id
		FROM emcd.users u
		LEFT JOIN emcd.vip_users v ON u.id = v.user_id
		LEFT JOIN histories.segment_userids su ON u.id = su.user_id
		WHERE u.new_id = $1;`
	var (
		vip       sql.NullBool
		segmentID sql.NullInt32
	)
	err = p.oldTrx.Runner(ctx).QueryRow(ctx, query, u.ID).Scan(&vip, &segmentID)
	if err != nil {
		return nil, fmt.Errorf("queryRow2: %w", err)
	}
	if vip.Valid {
		u.Vip = vip.Bool
	}
	if segmentID.Valid {
		u.SegmentID = int(segmentID.Int32)
	}
	if segmentID.Int32 == 0 {
		u.SegmentID, err = p.addSegment(ctx, u.ID)
		if err != nil {
			return nil, fmt.Errorf("addSegment: %w", err)
		}
	}
	return &u, nil
}

func (p *profile) GetCountUserWithWL(
	ctx context.Context,
	wlUUID uuid.UUID,
	offset, limit int32,
) ([]model.UserShortInfo, int64, error) {
	b := pgx.Batch{}
	const getUsersByWlUUID = `
SELECT id,username,email,created_at
FROM users WHERE whitelabel_id = $1
           ORDER BY created_at LIMIT $2 OFFSET $3`
	b.Queue(getUsersByWlUUID, wlUUID, limit, offset)

	const countQuery = `SELECT COUNT(*) FROM users WHERE whitelabel_id = $1`
	b.Queue(countQuery, wlUUID)

	res := p.trx.Runner(ctx).SendBatch(ctx, &b)
	defer func() {
		err := res.Close()
		if err != nil {
			log.Error(ctx, "GetCountUserWithWL: %v", err)
		}
	}()

	rows, err := res.Query()
	if err != nil {
		return nil, 0, fmt.Errorf("GetCountUserWithWL: query getUsersByWlUUID: %w", err)
	}
	defer rows.Close()

	userNames := make([]model.UserShortInfo, 0, 10)
	for rows.Next() {
		var userUUID sql.NullString
		var username sql.NullString
		var userEmail sql.NullString
		var createdAt sql.NullTime
		err = rows.Scan(&userUUID, &username, &userEmail, &createdAt)
		if err != nil {
			return nil, 0, fmt.Errorf("GetCountUserWithWL: scan username: %w", err)
		}
		userNames = append(userNames, model.UserShortInfo{
			UserUUID:  userUUID.String,
			UserName:  username.String,
			Email:     userEmail.String,
			CreatedAt: createdAt.Time,
		})
	}

	var count int64
	err = res.QueryRow().Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("GetCountUserWithWL: query countQuery: %w", err)
	}
	return userNames, count, nil
}

func (p *profile) GetAllUsers(
	ctx context.Context,
	skip, take int,
	sortField string,
	asc bool,
	searchEmail string,
) ([]*model.User, int, error) {
	searchEmail = strings.ToLower(searchEmail)
	orderDirection := "DESC"
	if asc {
		orderDirection = "ASC"
	}

	oldUsersMap, err := p.getAllOldUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("getAllOldUsers: %w", err)
	}

	b := pgx.Batch{}
	query := `SELECT  u.id,username,email,password,created_at,u.whitelabel_id,api_key,ref_id,parent_id,is_active,new_ref_id,pool_type
		FROM users u
		WHERE LOWER(email) LIKE $1 ORDER BY %s %s OFFSET $2 LIMIT $3`
	b.Queue(fmt.Sprintf(query, sortField, orderDirection),
		fmt.Sprintf(searchPatternStartsWith, searchEmail), skip, take)
	b.Queue("SELECT COUNT(*) FROM users WHERE LOWER(email) LIKE $1", fmt.Sprintf(searchPatternStartsWith, searchEmail))
	res := p.trx.Runner(ctx).SendBatch(ctx, &b)
	defer func() {
		err := res.Close()
		if err != nil {
			log.Error(ctx, "users get all: close batch: %v", err)
		}
	}()
	rows, err := res.Query()
	if err != nil {
		return nil, 0, fmt.Errorf("users get all: query users: %w", err)
	}
	defer rows.Close()
	users := make([]*model.User, 0)
	for rows.Next() {
		var u model.User
		var pwd []byte
		var (
			isActive sql.NullBool
		)
		err = rows.Scan(&u.ID, &u.Username, &u.Email, &pwd, &u.CreatedAt, &u.WhiteLabelID, &u.ApiKey,
			&u.RefID, &u.ParentID, &isActive, &u.NewRefID, &u.PoolType)
		if err != nil {
			return nil, 0, fmt.Errorf("users get all: scan users: %w", err)
		}
		u.Password = string(pwd)
		if isActive.Valid {
			u.IsActive = isActive.Bool
		}
		oldUser, ok := oldUsersMap[u.ID]
		if ok {
			u.SegmentID = oldUser.SegmentID
			u.Vip = oldUser.Vip
		}
		users = append(users, &u)
	}
	var count int
	err = res.QueryRow().Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("users get all: scan total count: %w", err)
	}
	return users, count, nil
}

func (p *profile) getAllOldUsers(ctx context.Context) (map[uuid.UUID]*model.User, error) {
	query := `
		SELECT
		    u.new_id,
		CASE
			WHEN v.user_id IS NOT NULL THEN true
			ELSE false
		END AS vip_user,
			su.segment_id
		FROM emcd.users u
		LEFT JOIN emcd.vip_users v ON u.id = v.user_id
		LEFT JOIN histories.segment_userids su ON u.id = su.user_id;`
	rows, err := p.oldTrx.Runner(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	oldUsersMap := make(map[uuid.UUID]*model.User)
	for rows.Next() {
		var (
			u         model.User
			vip       sql.NullBool
			segmentID sql.NullInt32
		)
		err = rows.Scan(&u.ID, &vip, &segmentID)
		if err != nil {
			return nil, fmt.Errorf("scan old users: %w", err)
		}
		if vip.Valid {
			u.Vip = vip.Bool
		}
		if segmentID.Valid {
			u.SegmentID = int(segmentID.Int32)
		}

		oldUsersMap[u.ID] = &u
	}
	return oldUsersMap, nil
}

func (p *profile) GetAllUsersByWlID(
	ctx context.Context,
	skip, take int,
	sortField string,
	asc bool,
	searchEmail string,
	whiteLabelID uuid.UUID,
) ([]*model.User, int, error) {
	searchEmail = strings.ToLower(searchEmail)
	orderDirection := "DESC"
	if asc {
		orderDirection = "ASC"
	}

	oldUsersMap, err := p.getAllOldUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("getAllOldUsers: %w", err)
	}

	b := pgx.Batch{}
	query := `SELECT u.id,username,email,password,created_at,u.whitelabel_id,api_key,ref_id,parent_id,is_active,new_ref_id,pool_type
		FROM users u
		WHERE u.whitelabel_id=$1 AND LOWER(email) LIKE $2 ORDER BY %s %s OFFSET $3 LIMIT $4`
	b.Queue(fmt.Sprintf(query, sortField, orderDirection),
		whiteLabelID, fmt.Sprintf(searchPatternStartsWith, searchEmail), skip, take)
	b.Queue(
		"SELECT COUNT(*) FROM users WHERE whitelabel_id=$1 AND LOWER(email) LIKE $2",
		whiteLabelID,
		fmt.Sprintf(searchPatternStartsWith, searchEmail),
	)
	res := p.trx.Runner(ctx).SendBatch(ctx, &b)
	defer func() {
		err := res.Close()
		if err != nil {
			log.Error(ctx, "users get all by wl id: close batch: %v", err)
		}
	}()
	rows, err := res.Query()
	if err != nil {
		return nil, 0, fmt.Errorf("users get all by wl id: query users: %w", err)
	}
	defer rows.Close()
	users := make([]*model.User, 0)
	for rows.Next() {
		var u model.User
		var pwd []byte
		var (
			isActive sql.NullBool
		)
		err = rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&pwd,
			&u.CreatedAt,
			&u.WhiteLabelID,
			&u.ApiKey,
			&u.RefID,
			&u.ParentID,
			&isActive,
			&u.NewRefID,
			&u.PoolType,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("users get all by wl id: scan users: %w", err)
		}
		u.Password = string(pwd)
		if isActive.Valid {
			u.IsActive = isActive.Bool
		}
		oldUser, ok := oldUsersMap[u.ID]
		if ok {
			u.SegmentID = oldUser.SegmentID
			u.Vip = oldUser.Vip
		}
		users = append(users, &u)
	}
	var count int
	err = res.QueryRow().Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("users get all by wl id: scan total count: %w", err)
	}
	return users, count, nil
}

func (p *profile) GetUserByApiKey(ctx context.Context, apiKey string) (*model.User, error) {
	getUserByEmailQuery := `
SELECT u.id,username,ref_id,LOWER(email),password,created_at,u.whitelabel_id,api_key,parent_id,is_active,new_ref_id
  FROM users u
 WHERE api_key=$1
 `
	var u model.User
	var pwd []byte
	var (
		isActive sql.NullBool
	)
	err := p.trx.Runner(ctx).QueryRow(ctx, getUserByEmailQuery, apiKey).Scan(&u.ID, &u.Username, &u.RefID, &u.Email,
		&pwd, &u.CreatedAt, &u.WhiteLabelID, &u.ApiKey, &u.ParentID, &isActive, &u.NewRefID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("profile: get user by api key: %w", err)
	}
	if isActive.Valid {
		u.IsActive = isActive.Bool
	}
	u.Password = string(pwd)
	query := `
		SELECT
		CASE
			WHEN v.user_id IS NOT NULL THEN true
			ELSE false
		END AS vip_user,
			su.segment_id
		FROM emcd.users u
		LEFT JOIN emcd.vip_users v ON u.id = v.user_id
		LEFT JOIN histories.segment_userids su ON u.id = su.user_id
		WHERE u.new_id = $1;`
	var (
		vip       sql.NullBool
		segmentID sql.NullInt32
	)
	err = p.oldTrx.Runner(ctx).QueryRow(ctx, query, u.ID).Scan(&vip, &segmentID)
	if err != nil {
		return nil, fmt.Errorf("queryRow: %w", err)
	}
	if vip.Valid {
		u.Vip = vip.Bool
	}
	if segmentID.Valid {
		u.SegmentID = int(segmentID.Int32)
	}
	if segmentID.Int32 == 0 {
		u.SegmentID, err = p.addSegment(ctx, u.ID)
		if err != nil {
			return nil, fmt.Errorf("addSegment: %w", err)
		}
	}
	return &u, nil
}

func (p *profile) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	u, err := p.getUserByID(ctx, userID, false)
	if err != nil {
		return nil, fmt.Errorf("profile: get by user id: %w", err)
	}
	if u == nil {
		return nil, nil
	}
	return &model.Profile{
		User: u,
	}, nil
}

func (p *profile) getUserByID(ctx context.Context, userID uuid.UUID, forUpdate bool) (*model.User, error) {
	var u model.User
	var pwd []byte
	var (
		isActive   sql.NullBool
		lang       sql.NullString
		tgID       sql.NullString
		tgUsername sql.NullString
	)
	getUserByIDQuery := `SELECT
		u.id,
		username,
		ref_id,
		LOWER(email),
		password,
		created_at,
		u.whitelabel_id,
		api_key,
		secret_key,
		parent_id,
		is_active,
		new_ref_id,
		language,
		pool_type,
		telegram_id,
		telegram_username,
		was_referral_link_generated,
		is_ambassador
	FROM users u
	WHERE u.id=$1`
	if forUpdate {
		getUserByIDQuery += " FOR UPDATE of u"
	}
	err := p.trx.Runner(ctx).QueryRow(ctx, getUserByIDQuery, userID).Scan(
		&u.ID,
		&u.Username,
		&u.RefID,
		&u.Email,
		&pwd,
		&u.CreatedAt,
		&u.WhiteLabelID,
		&u.ApiKey,
		&u.SecretKey,
		&u.ParentID,
		&isActive,
		&u.NewRefID,
		&lang,
		&u.PoolType,
		&tgID,
		&tgUsername,
		&u.WasReferralLinkGenerated,
		&u.IsAmbassador,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	if isActive.Bool {
		u.IsActive = isActive.Bool
	}
	if lang.Valid {
		u.Language = lang.String
	}
	if tgID.Valid {
		u.TgID = tgID.String
	}
	if tgUsername.Valid {
		u.TgUsername = tgUsername.String
	}
	u.Password = string(pwd)
	query := `
		SELECT
		CASE
			WHEN v.user_id IS NOT NULL THEN true
			ELSE false
		END AS vip_user,
			su.segment_id
		FROM emcd.users u
		LEFT JOIN emcd.vip_users v ON u.id = v.user_id
		LEFT JOIN histories.segment_userids su ON u.id = su.user_id
		WHERE u.new_id = $1;`
	var (
		vip       sql.NullBool
		segmentID sql.NullInt32
	)
	err = p.oldTrx.Runner(ctx).QueryRow(ctx, query, u.ID).Scan(&vip, &segmentID)
	if err != nil {
		return nil, fmt.Errorf("queryRow: %w", err)
	}
	if vip.Valid {
		u.Vip = vip.Bool
	}
	if segmentID.Valid {
		u.SegmentID = int(segmentID.Int32)
	}
	if segmentID.Int32 == 0 {
		u.SegmentID, err = p.addSegment(ctx, u.ID)
		if err != nil {
			return nil, fmt.Errorf("addSegment: %w", err)
		}
	}
	return &u, nil
}

func (p *profile) GetByUsernames(ctx context.Context, usernames []string) ([]*model.Profile, error) {
	u, err := p.getUsersByUsername(ctx, usernames)
	if err != nil {
		return nil, fmt.Errorf("profile: get by usernames: %w", err)
	}
	if u == nil {
		return nil, nil
	}

	profiles := make([]*model.Profile, 0)
	for _, user := range u {
		profiles = append(profiles, &model.Profile{
			User: user,
		})
	}
	return profiles, nil
}

func (p *profile) GetByUsernamesForReferrals(ctx context.Context, usernames []string) ([]*model.Profile, error) {
	u, err := p.getUsersByUsernameForReferrals(ctx, usernames)
	if err != nil {
		return nil, fmt.Errorf("profile: get by usernames for referrals: %w", err)
	}
	if u == nil {
		return nil, nil
	}

	profiles := make([]*model.Profile, 0)
	for _, user := range u {
		profiles = append(profiles, &model.Profile{
			User: user,
		})
	}
	return profiles, nil
}

func (p *profile) getUsersByUsernameForReferrals(ctx context.Context, usernames []string) ([]*model.User, error) {
	getUsersByUsernameForReferrals := `
SELECT
    id,
    username,
    ref_id,
    LOWER(email),
    created_at,
    whitelabel_id,
    parent_id,
    is_active,
    new_ref_id,
    language,
    pool_type
FROM users WHERE LOWER(username)=ANY($1)`

	rows, err := p.trx.Runner(ctx).Query(ctx, getUsersByUsernameForReferrals, usernames)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("getUsersByUsernameForReferrals: %w", err)
	}
	defer rows.Close()
	users := make([]*model.User, 0)
	for rows.Next() {
		var u model.User
		var (
			isActive sql.NullBool
			lang     sql.NullString
		)
		err = rows.Scan(
			&u.ID,
			&u.Username,
			&u.RefID,
			&u.Email,
			&u.CreatedAt,
			&u.WhiteLabelID,
			&u.ParentID,
			&isActive,
			&u.NewRefID,
			&lang,
			&u.PoolType)
		if err != nil {
			return nil, fmt.Errorf("scan old users: %w", err)
		}
		if isActive.Bool {
			u.IsActive = isActive.Bool
		}
		if lang.Valid {
			u.Language = lang.String
		}
		users = append(users, &u)
	}
	return users, nil
}

func (p *profile) getUsersByUsername(ctx context.Context, usernames []string) ([]*model.User, error) {
	getUserByUsernamesQuery := `
SELECT u.id,username,ref_id,LOWER(email),password,created_at,u.whitelabel_id,api_key,parent_id,is_active,new_ref_id,language,pool_type
  FROM users u
 WHERE LOWER(u.username)=ANY($1)
`

	rows, err := p.trx.Runner(ctx).Query(ctx, getUserByUsernamesQuery, usernames)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("getUsersByUsername: %w", err)
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		var u model.User
		var pwd []byte
		var (
			isActive sql.NullBool
			lang     sql.NullString
		)
		err = rows.Scan(&u.ID, &u.Username, &u.RefID, &u.Email,
			&pwd, &u.CreatedAt, &u.WhiteLabelID, &u.ApiKey, &u.ParentID, &isActive, &u.NewRefID, &lang, &u.PoolType)
		if err != nil {
			return nil, fmt.Errorf("scan old users: %w", err)
		}
		if isActive.Bool {
			u.IsActive = isActive.Bool
		}
		if lang.Valid {
			u.Language = lang.String
		}
		u.Password = string(pwd)

		query := `
		SELECT
		CASE
			WHEN v.user_id IS NOT NULL THEN true
			ELSE false
		END AS vip_user,
			su.segment_id
		FROM emcd.users u
		LEFT JOIN emcd.vip_users v ON u.id = v.user_id
		LEFT JOIN histories.segment_userids su ON u.id = su.user_id
		WHERE u.new_id = $1;`
		var (
			vip       sql.NullBool
			segmentID sql.NullInt32
		)
		err = p.oldTrx.Runner(ctx).QueryRow(ctx, query, u.ID).Scan(&vip, &segmentID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Info(ctx, "queryRow by segments not found %s, ignoring...", u.ID.String())
				u.Vip = false
			} else {
				return nil, fmt.Errorf("queryRow: %w, %s", err, u.ID.String())
			}
		} else {
			if vip.Valid {
				u.Vip = vip.Bool
			}
			if segmentID.Valid {
				u.SegmentID = int(segmentID.Int32)
			}
			if segmentID.Int32 == 0 {
				u.SegmentID, err = p.addSegment(ctx, u.ID)
				if err != nil {
					return nil, fmt.Errorf("addSegment: %w", err)
				}
			}
		}

		users = append(users, &u)
	}

	return users, nil
}

func (p *profile) UpdateUserIsActive(ctx context.Context, email string, active bool) error {
	email = strings.ToLower(email)
	query := `UPDATE users SET is_active = $1 WHERE email = $2`
	_, err := p.trx.Runner(ctx).Exec(ctx, query, active, email)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (p *profile) GetUserIsActive(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `
	SELECT COALESCE(u2.is_active, u1.is_active)
	FROM users u1
	LEFT JOIN users u2 ON u1.parent_id = u2.id
	WHERE u1.id = $1;`
	var isActive bool
	err := p.trx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&isActive)
	if err != nil {
		return false, fmt.Errorf("queryRow: %w", err)
	}
	return isActive, nil
}

func (p *profile) ExistAppleID(ctx context.Context, appleID string) (bool, string, error) {
	query := `SELECT LOWER(email) FROM users WHERE apple_id = $1`
	var em string
	err := p.trx.Runner(ctx).QueryRow(ctx, query, appleID).Scan(&em)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, "", fmt.Errorf("queryRow: %w", err)
	} else if errors.Is(err, pgx.ErrNoRows) {
		return false, "", nil
	}
	return true, em, nil
}

func (p *profile) UpdateAppleID(ctx context.Context, email, appleID string) (int, error) {
	email = strings.ToLower(email)
	query := `UPDATE users SET apple_id = $1 WHERE LOWER(email) = $2`
	commandTag, err := p.trx.Runner(ctx).Exec(ctx, query, appleID, email)
	if err != nil {
		return 0, fmt.Errorf("exec: %w", err)
	}
	return int(commandTag.RowsAffected()), nil
}

func (p *profile) RelatedUsers(ctx context.Context, firstID, secondID uuid.UUID) (bool, error) {
	query := `SELECT id, parent_id FROM users WHERE id = $1 or id = $2`
	rows, err := p.trx.Runner(ctx).Query(ctx, query, firstID, secondID)
	if err != nil {
		return false, fmt.Errorf("query")
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.ParentID)
		if err != nil {
			return false, fmt.Errorf("scan: %w", err)
		}
		users = append(users, &user)
	}
	if len(users) != 2 {
		return false, fmt.Errorf("unexpected count of users: %d by id: %s and id: %s", len(users), firstID, secondID)
	}

	u1 := users[0]
	u2 := users[1]

	// Parent and sub
	if u1.ParentID == uuid.Nil && u2.ParentID == u1.ID {
		return true, nil
	}

	// Sub and sub
	if u1.ParentID != uuid.Nil && u2.ParentID != uuid.Nil && u1.ParentID == u2.ParentID {
		return true, nil
	}

	// Sub and parent
	if u2.ParentID == uuid.Nil && u1.ParentID == u2.ID {
		return true, nil
	}

	return false, nil
}

func (p *profile) GetAllSubUsers(ctx context.Context, userID uuid.UUID) ([]*model.User, error) {
	query := `
SELECT users_child.new_id, users_child.username, users_child.id
  FROM emcd.users users_child join 
       emcd.users users_parent on users_parent.id = users_child.parent_id
 WHERE users_parent.new_id = $1;`
	rows, err := p.oldTrx.Runner(ctx).Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Username, &user.OldID)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		users = append(users, &user)
	}
	return users, nil
}

func (p *profile) GetAllUserIDsByUsername(ctx context.Context) (map[string]uuid.UUID, error) {
	query := `SELECT id, username FROM users`
	rows, err := p.trx.Runner(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	m := make(map[string]uuid.UUID)
	for rows.Next() {
		var (
			id       uuid.UUID
			username string
		)
		err = rows.Scan(&id, &username)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		m[username] = id
	}
	return m, nil
}

func (p *profile) GetReferrals(
	ctx context.Context,
	userID uuid.UUID,
	skip, take int,
	sortField string,
	asc bool,
) ([]*model.User, int, error) {
	orderDirection := "DESC"
	if asc {
		orderDirection = "ASC"
	}

	selectAllUsers := `
	SELECT id, username, LOWER(email), created_at
	FROM users
	WHERE new_ref_id = $1 ORDER BY %s %s OFFSET $2 LIMIT $3`

	selectCount := `SELECT count(1) FROM users WHERE new_ref_id = $1`

	b := pgx.Batch{}
	b.Queue(fmt.Sprintf(selectAllUsers, sortField, orderDirection), userID, skip, take)
	b.Queue(selectCount, userID)
	res := p.trx.Runner(ctx).SendBatch(ctx, &b)
	defer func() {
		err := res.Close()
		if err != nil {
			log.Error(ctx, "getReferrals res.Close: %s", err.Error())
		}
	}()

	rows, err := res.Query()
	if err != nil {
		return nil, 0, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	refs := make([]*model.User, 0)
	for rows.Next() {
		var u model.User
		err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan selectAllUsers: %w", err)
		}
		refs = append(refs, &u)
	}

	var count int
	err = res.QueryRow().Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("scan selectCount: %w", err)
	}

	return refs, count, nil
}

func (p *profile) GetSuspended(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `SELECT suspended FROM users WHERE id = $1`
	var suspended sql.NullTime
	err := p.trx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&suspended)
	if err != nil {
		return false, fmt.Errorf("select: %w", err)
	}
	if suspended.Valid {
		return true, nil
	} else {
		return false, nil
	}
}

func (p *profile) SetSuspended(ctx context.Context, userIDs []uuid.UUID, suspended bool) error {
	var query string
	switch suspended {
	case true:
		query = `UPDATE users SET suspended = now() WHERE id = ANY($1)`
	case false:
		query = `UPDATE users SET suspended = null WHERE id = ANY($1)`
	}

	_, err := p.trx.Runner(ctx).Exec(ctx, query, userIDs)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (p *profile) GetUsernamesByIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]string, error) {
	query := `SELECT id, username FROM users WHERE id=ANY($1)`
	rows, err := p.trx.Runner(ctx).Query(ctx, query, userIDs)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	usernames := make(map[uuid.UUID]string)
	for rows.Next() {
		var (
			username string
			id       uuid.UUID
		)
		err = rows.Scan(&id, &username)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		usernames[id] = username
	}
	return usernames, nil
}

func (p *profile) GetUsersByUUIDs(ctx context.Context, userUUIDs []uuid.UUID) ([]model.UserShortInfo, error) {
	query := `SELECT id,username,email,created_at FROM users WHERE id=ANY($1) ORDER BY created_at DESC`
	rows, err := p.trx.Runner(ctx).Query(ctx, query, userUUIDs)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	userNames := make([]model.UserShortInfo, 0, 10)
	for rows.Next() {
		var userUUID sql.NullString
		var username sql.NullString
		var userEmail sql.NullString
		var createdAt sql.NullTime
		err = rows.Scan(&userUUID, &username, &userEmail, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("GetUsersByUUIDs: scan username: %w", err)
		}
		userNames = append(userNames, model.UserShortInfo{
			UserUUID:  userUUID.String,
			UserName:  username.String,
			Email:     userEmail.String,
			CreatedAt: createdAt.Time,
		})
	}
	return userNames, nil
}

func (p *profile) GetEmailsByIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]string, error) {
	query := `SELECT id, LOWER(email) FROM users WHERE id=ANY($1)`
	rows, err := p.trx.Runner(ctx).Query(ctx, query, userIDs)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	emails := make(map[uuid.UUID]string)
	for rows.Next() {
		var (
			email string
			id    uuid.UUID
		)
		err = rows.Scan(&id, &email)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		emails[id] = email
	}
	return emails, nil
}

func (p *profile) addSegment(ctx context.Context, userID uuid.UUID) (int, error) {
	var oldID int
	query := `SELECT id FROM emcd.users WHERE new_id = $1`
	err := p.oldTrx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&oldID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("queryRow1: %w", err)
	}

	query = `INSERT INTO histories.segment_userids (user_id) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err = p.oldTrx.Runner(ctx).Exec(ctx, query, oldID)
	if err != nil {
		return 0, fmt.Errorf("exec: %w", err)
	}

	var segmentID int
	query = `--antislave
		SELECT segment_id FROM histories.segment_userids WHERE user_id = $1`
	err = p.oldTrx.Runner(ctx).QueryRow(ctx, query, oldID).Scan(&segmentID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("queryRow2: %w", err)
	}

	return segmentID, nil
}

const (
	adminUsername = "admin"
	deletePrefix  = "deleted"
	deleteDeter   = "_"
)

func (p *profile) SafeDeleteByID(ctx context.Context, userID uuid.UUID) error {
	user, err := p.getUserByID(ctx, userID, false)
	if err != nil {
		return fmt.Errorf("profile: get user by id: %w", err)
	}

	userOldId, err := p.getOldIDbyID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("profile: getOldIDbyID: %w", err)
	}

	admin, err := p.getUserByUsername(ctx, adminUsername)
	if err != nil {
		return fmt.Errorf("profile: get user by username: %w", err)
	}

	adminOldId, err := p.getOldIDbyID(ctx, admin.ID)
	if err != nil {
		return fmt.Errorf("profile: getOldIDbyID: %w", err)
	}

	uniq := false
	var delUsername string
	var delEmail string
	for i := 1; i <= 50; i++ {
		deter := strings.Repeat(deleteDeter, i)
		delUsername = deletePrefix + deter + user.Username
		delEmail = deletePrefix + deter + user.Email

		check, err := p.findUniqUsernameOrEmail(ctx, delUsername, delEmail)
		if err != nil {
			return fmt.Errorf("profile: checkUniqUsernameOrEmail: %w", err)
		}
		if check == false {
			uniq = true
			break
		}
	}
	if uniq == false {
		return fmt.Errorf("profile: safeDelete: count of duplicate is too more. user id: %v", user.ID)
	}

	err = p.trx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// txCtx is a context provided by pool db Runner
		// while ctx is a context provided by emcd db Runner
		// need use txCtx for quering to pool db
		// and ctx - to emcd db
		inTxErr := p.safeDeleteUser(txCtx, userID, delUsername, delEmail, admin.ID, adminOldId)
		if inTxErr != nil {
			return fmt.Errorf("profile: safeDelete: %w", inTxErr)
		}
		inTxErr = p.updateReferrals(txCtx, user.ID, admin.ID, adminOldId)
		if inTxErr != nil {
			return fmt.Errorf("profile: safeDelete: updateReferrals: %w", inTxErr)
		}

		inTxErr = p.safeDeleteUserOld(ctx, userID, delUsername, delEmail, adminOldId)
		if inTxErr != nil {
			return fmt.Errorf("profile: safeDeleteOld: %w", inTxErr)
		}
		inTxErr = p.updateReferralsOld(ctx, userOldId, adminOldId)
		if inTxErr != nil {
			return fmt.Errorf("profile: safeDelete: updateReferralsOld: %w", inTxErr)
		}
		return nil
	})

	return err
}

func (p *profile) safeDeleteUser(
	ctx context.Context,
	uuid uuid.UUID,
	newUsername string,
	newEmail string,
	newRefID uuid.UUID,
	newOldRefID int32,
) error {
	newEmail = strings.ToLower(newEmail)

	query := `UPDATE users SET email=$1, password=$2, username=$3, is_active = false, new_ref_id = $4, ref_id = $5, apple_id = null, api_key = ''
                    where id=$6;`
	_, err := p.trx.Runner(ctx).Exec(ctx, query, newEmail, deletePrefix, newUsername, newRefID, newOldRefID, uuid)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (p *profile) safeDeleteUserOld(
	ctx context.Context,
	uuid uuid.UUID,
	newUsername string,
	newEmail string,
	newRefID int32,
) error {
	newEmail = strings.ToLower(newEmail)

	query := `UPDATE users SET email=$1, password=$2, username=$3, is_active = false,
			  nopay=true, is_email_notifications_on = false, is_tg_notifications_on = false, api_key = null,
			  tg_id = NULL, ref_id=$4 WHERE new_id=$5;`
	_, err := p.oldTrx.Runner(ctx).Exec(ctx, query, newEmail, deletePrefix, newUsername, newRefID, uuid)
	if err != nil {
		return fmt.Errorf("update (emcd): %w", err)
	}

	return nil
}

func (p *profile) updateReferrals(ctx context.Context, oldNewRefID uuid.UUID, newRefID uuid.UUID, refID int32) error {
	query := `UPDATE users SET ref_id=$1, new_ref_id=$2
                    where new_ref_id=$3;`
	_, err := p.trx.Runner(ctx).Exec(ctx, query, refID, newRefID, oldNewRefID)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (p *profile) updateReferralsOld(ctx context.Context, oldRefID int32, newRefID int32) error {
	query := `UPDATE users SET ref_id=$1
                    where ref_id=$2;`
	_, err := p.oldTrx.Runner(ctx).Exec(ctx, query, newRefID, oldRefID)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (p *profile) findUniqUsernameOrEmail(ctx context.Context, username string, email string) (bool, error) {
	var exists bool
	email = strings.ToLower(email)
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1 OR LOWER(email)=$2);`
	err := p.trx.Runner(ctx).QueryRow(ctx, query, username, email).Scan(&exists)
	return exists, err
}

func (p *profile) getOldIDbyID(ctx context.Context, uuid uuid.UUID) (int32, error) {
	var oldID int32
	query := `SELECT id from users where new_id = $1;`
	err := p.oldTrx.Runner(ctx).QueryRow(ctx, query, uuid).Scan(&oldID)
	return oldID, err
}

func (p *profile) getUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	var pwd []byte
	var (
		isActive sql.NullBool
		lang     sql.NullString
	)
	getUserByIDQuery := `
SELECT u.id,username,ref_id,LOWER(email),password,created_at,u.whitelabel_id,api_key,parent_id,is_active,new_ref_id,language
  FROM users u
 WHERE u.username=$1 LIMIT 1`
	err := p.trx.Runner(ctx).QueryRow(ctx, getUserByIDQuery, username).Scan(&u.ID, &u.Username, &u.RefID, &u.Email,
		&pwd, &u.CreatedAt, &u.WhiteLabelID, &u.ApiKey, &u.ParentID, &isActive, &u.NewRefID, &lang)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	if isActive.Bool {
		u.IsActive = isActive.Bool
	}
	if lang.Valid {
		u.Language = lang.String
	}
	u.Password = string(pwd)
	query := `
		SELECT
		CASE
			WHEN v.user_id IS NOT NULL THEN true
			ELSE false
		END AS vip_user,
			su.segment_id
		FROM emcd.users u
		LEFT JOIN emcd.vip_users v ON u.id = v.user_id
		LEFT JOIN histories.segment_userids su ON u.id = su.user_id
		WHERE u.username = $1 LIMIT 1;`
	var (
		vip       sql.NullBool
		segmentID sql.NullInt32
	)
	err = p.oldTrx.Runner(ctx).QueryRow(ctx, query, u.Username).Scan(&vip, &segmentID)
	if err != nil {
		return nil, fmt.Errorf("queryRow: %w", err)
	}
	if vip.Valid {
		u.Vip = vip.Bool
	}
	if segmentID.Valid {
		u.SegmentID = int(segmentID.Int32)
	}
	return &u, nil
}

// update
func (p *profile) UpdateUser(
	ctx context.Context,
	userID uuid.UUID,
	mutate func(user *model.User) error,
) error {
	tx, err := p.trx.Runner(ctx).Begin(ctx)
	if err != nil {
		return fmt.Errorf("UpdateUser: update user: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	user, err := p.getUserByID(ctx, userID, true)
	if err != nil {
		return fmt.Errorf("UpdateUser: get user: %w", err)
	}

	err = mutate(user)
	if err != nil {
		return fmt.Errorf("UpdateUser: mutate user: %w", err)
	}

	err = p.SaveUser(ctx, user)
	if err != nil {
		return fmt.Errorf("UpdateUser: save user: %w", err)
	}

	return tx.Commit(ctx)
}

// Deprecated: use UpdateUser instead. Should became private
func (p *profile) SaveUser(ctx context.Context, u *model.User) error {
	const op = "repo.profile.SaveUser"

	u.Email = strings.ToLower(u.Email)
	var appleID sql.NullString
	if u.AppleID != "" {
		appleID.Valid = true
		appleID.String = u.AppleID
	}

	upsertUserQuery := `
INSERT INTO users (
	id,
	username,
	ref_id,
	email,
	password,
	created_at,
	whitelabel_id,
	api_key,
	is_active,
	apple_id,
	parent_id,
	new_ref_id,
	pool_type,
	language,
	telegram_id,
	telegram_username,
	was_referral_link_generated,
	is_ambassador
)
VALUES (
	@id,
	@username,
	@ref_id,
	@email,
	@password,
	@created_at,
	@whitelabel_id,
	@api_key,
	@is_active,
	@apple_id,
	@parent_id,
	@new_ref_id,
	@pool_type,
	@language,
	@telegram_id,
	@telegram_username,
	@was_referral_link_generated,
	@is_ambassador
)
ON CONFLICT (id) DO
	UPDATE SET
		username=@username,
		ref_id=@ref_id,
		password=@password,
		created_at=@created_at,
		api_key=@api_key,
		is_active=@is_active,
		apple_id=@apple_id,
		parent_id=@parent_id,
		new_ref_id=@new_ref_id,
		pool_type=@pool_type,
		whitelabel_id=@whitelabel_id,
		language=@language,
		telegram_id=@telegram_id,
		telegram_username=@telegram_username,
		was_referral_link_generated=@was_referral_link_generated,
		is_ambassador=@is_ambassador
`
	_, err := p.trx.Runner(ctx).
		Exec(
			ctx,
			upsertUserQuery,
			pgx.NamedArgs{
				"id":                          u.ID,
				"username":                    u.Username,
				"ref_id":                      u.RefID,
				"email":                       u.Email,
				"password":                    u.Password,
				"created_at":                  u.CreatedAt,
				"whitelabel_id":               u.WhiteLabelID,
				"api_key":                     u.ApiKey,
				"is_active":                   u.IsActive,
				"apple_id":                    appleID,
				"parent_id":                   u.ParentID,
				"new_ref_id":                  u.NewRefID,
				"pool_type":                   u.PoolType,
				"language":                    u.Language,
				"telegram_id":                 u.TgID,
				"telegram_username":           u.TgUsername,
				"was_referral_link_generated": u.WasReferralLinkGenerated,
				"is_ambassador":               u.IsAmbassador,
			},
		)
	if err != nil {
		return fmt.Errorf("%s: failed to save user: %w", op, err)
	}
	return nil
}

func (p *profile) UpdateRefID(ctx context.Context, userID uuid.UUID, refID int32, newRefID uuid.UUID) error {
	log.Info(ctx, "user: %s update ref_id: %d", userID, refID)
	query := `UPDATE users SET ref_id=$1, new_ref_id=$2 WHERE id=$3 OR parent_id=$3`

	commandTag, err := p.trx.Runner(ctx).Exec(ctx, query, refID, newRefID, userID)
	if err != nil {
		return fmt.Errorf("profile: update ref_id: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("profile: update ref_id: unknown user_id: %s", userID)
	}
	return nil
}

func (p *profile) GetFlagReferralLinkGenerated(ctx context.Context, userID uuid.UUID) (bool, error) {
	var flag bool
	log.Info(ctx, "user: %s get was_referral_link_generated: %s", userID, userID)
	query := `SELECT was_referral_link_generated from users where id = $1`

	err := p.trx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&flag)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return flag, fmt.Errorf("profile: user %s not found", userID)
	} else if err != nil {
		return false, fmt.Errorf("profile: get user %w", err)
	}
	return flag, nil
}

// Deprecated: use UpdateUser instead
func (p *profile) SetFlagReferralLinkGenerated(ctx context.Context, userID uuid.UUID, value bool) error {
	log.Info(ctx, "user: %s update was_referral_link_generated: %s", userID, userID)
	query := `UPDATE users SET was_referral_link_generated=$1 WHERE id=$2`
	commandTag, err := p.trx.Runner(ctx).Exec(ctx, query, value, userID)
	if err != nil {
		return fmt.Errorf("profile: update was_referral_link_generated: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("profile: update was_referral_link_generated: unknown user_id: %s", userID)
	}
	return nil
}

func (p *profile) GetAllReferrals(
	ctx context.Context,
	referrerID uuid.UUID,
	retrieveAfterUsername *string,
	sortField string,
	asc bool,
	take int,
) (map[string]*model.Referral, []string, error) {
	order := buildOrderClause(sortField, asc)
	var query string
	var args []interface{}

	if retrieveAfterUsername == nil {
		query = `SELECT id, username, email, created_at FROM users WHERE new_ref_id = $1 ORDER BY %s LIMIT $2`
		query = fmt.Sprintf(query, order)
		args = append(args, referrerID, take)
	} else {
		if asc {
			query = `SELECT id, username, email, created_at FROM users WHERE new_ref_id = $1 AND %s > (SELECT %s FROM users WHERE username = $2) ORDER BY %s LIMIT $3`
		} else {
			query = `SELECT id, username, email, created_at FROM users WHERE new_ref_id = $1 AND %s < (SELECT %s FROM users WHERE username = $2) ORDER BY %s LIMIT $3`
		}
		query = fmt.Sprintf(query, sortField, sortField, order)
		args = append(args, referrerID, *retrieveAfterUsername, take)
	}

	rows, err := p.trx.Runner(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	referrals := make(map[string]*model.Referral)
	orderedUsernames := make([]string, 0)
	for rows.Next() {
		var r model.Referral
		err = rows.Scan(&r.UserID, &r.Username, &r.Email, &r.CreatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("scan: %w", err)
		}
		orderedUsernames = append(orderedUsernames, r.Username)
		referrals[r.Username] = &r
	}
	return referrals, orderedUsernames, nil
}

func (p *profile) SetLanguage(ctx context.Context, userID uuid.UUID, language string) error {
	query := `UPDATE users SET language = $1 WHERE id = $2`
	_, err := p.trx.Runner(ctx).Exec(ctx, query, language, userID)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func buildOrderClause(sortField string, asc bool) string {
	order := sortField
	if asc {
		order += " ASC"
	} else {
		order += " DESC"
	}
	return order
}

func (p *profile) GetAPIKeyByUserIDAndParentID(ctx context.Context, userID, parentID uuid.UUID) (bool, string, error) {
	query := `SELECT api_key FROM users WHERE id = $1 AND parent_id = $2`
	var apiKey string

	err := p.trx.Runner(ctx).QueryRow(ctx, query, userID, parentID).Scan(&apiKey)

	switch {
	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		return false, "", fmt.Errorf("profile: get api key by pair user_id and parent_id: %w", err)
	case errors.Is(err, pgx.ErrNoRows):
		return false, "", nil
	}

	return apiKey != "", apiKey, nil
}

func (p *profile) SetAPIKeyForUserIDAndParentID(ctx context.Context, apiKey string, userID, parentID uuid.UUID) error {
	query := `UPDATE users SET api_key=$1 WHERE id=$2 AND parent_id=$3`
	commandTag, err := p.trx.Runner(ctx).Exec(ctx, query, apiKey, userID, parentID)
	if err != nil {
		return fmt.Errorf("profile: update api_key by user_id and parent_id: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("profile: update api_key: unknown pair user_id and parent_id: %s, %s", userID, parentID)
	}
	return nil
}

func (p *profile) GetAPIKeyByUserID(ctx context.Context, userID uuid.UUID) (bool, string, error) {
	query := `SELECT api_key FROM users WHERE id = $1`
	var apiKey string

	err := p.trx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&apiKey)

	switch {
	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		return false, "", fmt.Errorf("profile: get api key by user_id: %w", err)
	case errors.Is(err, pgx.ErrNoRows):
		return false, "", nil
	}

	return apiKey != "", apiKey, nil
}

func (p *profile) SetAPIKeyForUserID(ctx context.Context, apiKey string, userID uuid.UUID) error {
	query := `UPDATE users SET api_key=$1 WHERE id=$2`
	commandTag, err := p.trx.Runner(ctx).Exec(ctx, query, apiKey, userID)
	if err != nil {
		return fmt.Errorf("profile: update api_key by user_id: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("profile: update api_key: unknown user_id: %s", userID)
	}
	return nil
}

func (p *profile) SetSecretKeyForUserID(ctx context.Context, secretKey string, userID uuid.UUID) error {
	query := `UPDATE users SET secret_key=$1 WHERE id=$2`
	commandTag, err := p.trx.Runner(ctx).Exec(ctx, query, secretKey, userID)
	if err != nil {
		return fmt.Errorf("profile: update secret_key by user_id: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("profile: update secret_key: unknown user_id: %s", userID)
	}
	return nil
}

func (p *profile) SoftDeleteSubUser(ctx context.Context, subUserID uuid.UUID, newParentID uuid.UUID, newEmail model.SubUserEmail) error {
	const op = "profile.SoftDeleteSubUser"
	err := p.markSubUserAsDeleted(ctx, newParentID, newEmail, subUserID)
	if err != nil {
		return fmt.Errorf("%s: mark user as deleted : %w", op, err)
	}

	return nil
}

func (p *profile) markSubUserAsDeleted(ctx context.Context, newParentID uuid.UUID, newEmail model.SubUserEmail, subUserID uuid.UUID) error {
	const query = `
		UPDATE users
		SET
			parent_id = $1, 
			username = CONCAT('deleted_', username), 
			password = '\x7375626163636f756e74', 
			email = $2
		WHERE
			id = $3
	`
	_, err := p.trx.Runner(ctx).Exec(ctx, query, newParentID, newEmail, subUserID)
	if err != nil {
		return fmt.Errorf("exec query: %v", err)
	}
	return nil
}

func (p *profile) Begin() pgTx.PgxTransactor {
	return p.trx
}

func (p *profile) GetLastSubUserEmail(ctx context.Context, parentID uuid.UUID) (string, error) {
	selectLastSubUserEmail := `SELECT email FROM users WHERE parent_id = $1 ORDER BY email DESC LIMIT 1`

	var lastSubUserEmail string
	if err := p.trx.Runner(ctx).QueryRow(ctx, selectLastSubUserEmail, parentID).Scan(&lastSubUserEmail); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("queryRow: %w", err)
	}
	return lastSubUserEmail, nil
}

func (p *profile) IsUniqueUsername(ctx context.Context, username string) (bool, error) {
	query := `--antislave 
SELECT COUNT(*) FROM users WHERE lower(username)=lower($1)`
	var cnt int
	err := p.trx.Runner(ctx).QueryRow(ctx, query, username).Scan(&cnt)
	if err != nil {
		return false, fmt.Errorf("query row: scan: %w", err)
	}
	return cnt == 0, nil
}

func (p *profile) GetUserByTg(ctx context.Context, tgID string) (*model.User, error) {
	if tgID == "" {
		return nil, nil
	}
	getUserByEmailQuery := `SELECT id,username,ref_id,email,password,created_at,whitelabel_id,api_key,parent_id,is_active,new_ref_id,suspended,telegram_id,telegram_username
		FROM users WHERE LOWER(telegram_id)=LOWER($1) AND telegram_id IS NOT NULL`
	var (
		u          model.User
		pwd        []byte
		isActive   sql.NullBool
		suspended  sql.NullTime
		tgId       sql.NullString
		tgUsername sql.NullString
	)
	err := p.trx.Runner(ctx).QueryRow(ctx, getUserByEmailQuery, tgID).Scan(&u.ID, &u.Username, &u.RefID, &u.Email,
		&pwd, &u.CreatedAt, &u.WhiteLabelID, &u.ApiKey, &u.ParentID, &isActive, &u.NewRefID, &suspended, &tgId, &tgUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("queryRow1: %w", err)
	}
	u.Password = string(pwd)
	if isActive.Valid {
		u.IsActive = isActive.Bool
	}
	if suspended.Valid {
		u.Suspended = true
	}
	if tgId.Valid {
		u.TgID = tgId.String
	}
	if tgUsername.Valid {
		u.TgUsername = tgUsername.String
	}

	query := `--anti slave
		SELECT
		CASE
			WHEN v.user_id IS NOT NULL THEN true
			ELSE false
		END AS vip_user,
			su.segment_id
		FROM emcd.users u
		LEFT JOIN emcd.vip_users v ON u.id = v.user_id
		LEFT JOIN histories.segment_userids su ON u.id = su.user_id
		WHERE u.new_id = $1;`
	var (
		vip       sql.NullBool
		segmentID sql.NullInt32
	)
	err = p.oldTrx.Runner(ctx).QueryRow(ctx, query, u.ID).Scan(&vip, &segmentID)
	if err != nil {
		return nil, fmt.Errorf("queryRow2: %w", err)
	}
	if vip.Valid {
		u.Vip = vip.Bool
	}
	if segmentID.Valid {
		u.SegmentID = int(segmentID.Int32)
	}
	if segmentID.Int32 == 0 {
		u.SegmentID, err = p.addSegment(ctx, u.ID)
		if err != nil {
			return nil, fmt.Errorf("addSegment: %w", err)
		}
	}
	return &u, nil
}
