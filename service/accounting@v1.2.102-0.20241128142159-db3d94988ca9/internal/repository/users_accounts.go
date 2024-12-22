package repository

import (
	"context"
	"fmt"
	"time"

	sdkLog "code.emcdtech.com/emcd/sdk/log"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
)

const repoName = "user_account"

type UserAccountRepo interface {
	AddUserAccounts(ctx context.Context, userId int32, _ uuid.UUID, userAccounts model.UserAccounts) error
	AddUserAccount(ctx context.Context, userAccount *model.UserAccount) error
	FindUserAccountByFilter(ctx context.Context, filter *model.UserAccountFilter) (*uint64, model.UserAccounts, error)
	FindUserAccountByFilterMigrateOnly(ctx context.Context, filter *model.UserAccountFilter) (*uint64, model.UserAccounts, error)
	UpdateUserAccountByPartial(ctx context.Context, userAccount *model.UserAccount, partial *model.UserAccountPartial) error
	UpdateUserAccountForMigrateUserIdNew(ctx context.Context, userAccounts model.UserAccounts) error
	UpdateUserAccountForMigrateCoinNew(ctx context.Context, userAccounts model.UserAccounts) error

	GetUserAccountByIdLegacy(ctx context.Context, id int32) (*model.UserAccount, error)
	FindUserAccountByUserIdLegacy(ctx context.Context, userIdNew uuid.UUID) (model.UserAccounts, error)

	transactor.PgxTransactor
}

type userAccountsImpl struct {
	pool *pgxpool.Pool
	transactor.PgxTransactor
}

func NewUserAccountRepo(pool *pgxpool.Pool) UserAccountRepo {
	pgxTransactor := transactor.NewPgxTransactor(pool)

	return &userAccountsImpl{
		pool:          pool,
		PgxTransactor: pgxTransactor,
	}
}

func (a *userAccountsImpl) AddUserAccount(ctx context.Context, userAccount *model.UserAccount) error {

	return a.AddUserAccounts(ctx, userAccount.UserID, userAccount.UserIDNew.UUID, model.UserAccounts{userAccount})
}

func (a *userAccountsImpl) AddUserAccounts(ctx context.Context, _ int32, userIdNew uuid.UUID, userAccounts model.UserAccounts) error {
	const queryLockUserId = `SELECT EXISTS(SELECT id FROM emcd.users_accounts WHERE user_id_new = @user_id_new FOR UPDATE)`

	//	const queryUniq = `
	//	INSERT INTO emcd.users_accounts (user_id, coin_id, account_type_id, minpay, address, changed_at, img1, img2, is_active, created_at, updated_at, fee, user_id_new, coin_new)
	//	SELECT @user_id, @coin_id, @account_type_id, @minpay, @address, @changed_at, @img1, @img2, @is_active, @created_at, @updated_at, @fee, @user_id_new, @coin_new
	//	WHERE NOT EXISTS (SELECT id FROM emcd.users_accounts WHERE user_id_new = @user_id_new AND coin_id = @coin_id AND account_type_id = @account_type_id)
	//	`
	//
	//	const queryNonUniq = `
	//	INSERT INTO emcd.users_accounts (user_id, coin_id, account_type_id, minpay, address, changed_at, img1, img2, is_active, created_at, updated_at, fee, user_id_new, coin_new)
	//	VALUES (@user_id, @coin_id, @account_type_id, @minpay, @address, @changed_at, @img1, @img2, @is_active, @created_at, @updated_at, @fee, @user_id_new, @coin_new)
	// `

	const queryUniq = `
	WITH new_row AS (
	INSERT INTO emcd.users_accounts (user_id, coin_id, account_type_id, minpay, address, changed_at, img1, img2, is_active, created_at, updated_at, fee, user_id_new, coin_new)
	SELECT @user_id, @coin_id, @account_type_id, @minpay, @address, @changed_at, @img1, @img2, @is_active, @created_at, @updated_at, @fee, @user_id_new, @coin_new
	WHERE NOT EXISTS (SELECT id FROM emcd.users_accounts WHERE user_id = @user_id AND coin_id = @coin_id AND account_type_id = @account_type_id)
	RETURNING id, user_id, coin_id, account_type_id, minpay, address, changed_at, img1, img2, is_active, created_at, updated_at, fee, user_id_new, coin_new
	)
	SELECT id, user_id, coin_id, account_type_id, minpay, address, changed_at, img1, img2, is_active, created_at, updated_at, fee, user_id_new, coin_new
	FROM new_row
	UNION
	SELECT id, user_id, coin_id, account_type_id, minpay, address, changed_at, img1, img2, is_active, created_at, updated_at, fee, user_id_new, coin_new
	FROM emcd.users_accounts
	WHERE user_id = @user_id AND coin_id = @coin_id AND account_type_id = @account_type_id
	`

	const queryNonUniq = `
	INSERT INTO emcd.users_accounts (user_id, coin_id, account_type_id, minpay, address, changed_at, img1, img2, is_active, created_at, updated_at, fee, user_id_new, coin_new)
	VALUES (@user_id, @coin_id, @account_type_id, @minpay, @address, @changed_at, @img1, @img2, @is_active, @created_at, @updated_at, @fee, @user_id_new, @coin_new)
	RETURNING id, user_id, coin_id, account_type_id, minpay, address, changed_at, img1, img2, is_active, created_at, updated_at, fee, user_id_new, coin_new
	`

	batch := pgx.Batch{}
	batch.Queue(queryLockUserId, pgx.NamedArgs{"user_id_new": userIdNew})

	for _, userAccount := range userAccounts {
		var query string
		if userAccount.AccountTypeID.AccountTypeId != enum.CoinholdAccountTypeID {
			query = queryUniq

		} else {
			query = queryNonUniq

		}
		dt := time.Now().UTC()

		if !userAccount.ChangedAt.Valid {
			userAccount.ChangedAt.Time = dt
			userAccount.ChangedAt.Valid = true

		}

		if !userAccount.IsActive.Valid {
			userAccount.IsActive.Bool = true
			userAccount.IsActive.Valid = true

		}

		if !userAccount.UpdatedAt.Valid {
			userAccount.UpdatedAt.Time = dt
			userAccount.UpdatedAt.Valid = true

		}

		if !userAccount.CreatedAt.Valid {
			userAccount.CreatedAt.Time = dt
			userAccount.CreatedAt.Valid = true

		}

		userSql := userAccountSql{UserAccount: userAccount}

		batch.Queue(query, userSql.ToNamedArgs())

	}

	batchResult := a.PgxTransactor.Runner(ctx).SendBatch(ctx, &batch)
	defer func(br pgx.BatchResults) {
		if err := br.Close(); err != nil {
			sdkLog.Error(ctx, "failed to close batch: %s", err.Error())

		}
	}(batchResult)

	var scanVar bool
	if err := batchResult.QueryRow().Scan(&scanVar); err != nil {

		return fmt.Errorf("failed scan %s: %w", repoName, err)
	}

	for _, userAccount := range userAccounts {
		row := batchResult.QueryRow()

		if err := row.Scan(&userAccount.ID,
			&userAccount.UserID,
			&userAccount.CoinID,
			&userAccount.AccountTypeID,
			&userAccount.Minpay,
			&userAccount.Address,
			&userAccount.ChangedAt,
			&userAccount.Img1,
			&userAccount.Img2,
			&userAccount.IsActive,
			&userAccount.CreatedAt,
			&userAccount.UpdatedAt,
			&userAccount.Fee,
			&userAccount.UserIDNew,
			&userAccount.CoinNew,
		); err != nil {

			return fmt.Errorf("failed scan %s: %w", repoName, err)
		} else if err := userAccount.SqlNoEmptyValidate(); err != nil {

			return fmt.Errorf("failed validate %s: %w", repoName, err)
		}
	}

	return nil
}

func (a *userAccountsImpl) FindUserAccountByFilter(ctx context.Context, filter *model.UserAccountFilter) (*uint64, model.UserAccounts, error) {
	if totalCount, usersAccount, err := a.findUserAccountByFilterWithValidate(ctx, filter); err != nil {

		return nil, nil, err
	} else {
		for _, userAccount := range usersAccount {
			if err := userAccount.SqlNoEmptyValidate(); err != nil {

				return nil, nil, fmt.Errorf("failed validate %s: %w", repoName, err)
			}
		}

		return totalCount, usersAccount, nil
	}
}

func (a *userAccountsImpl) FindUserAccountByFilterMigrateOnly(ctx context.Context, filter *model.UserAccountFilter) (*uint64, model.UserAccounts, error) {

	return a.findUserAccountByFilterWithValidate(ctx, filter)
}

func (a *userAccountsImpl) findUserAccountByFilterWithValidate(ctx context.Context, filter *model.UserAccountFilter) (*uint64, model.UserAccounts, error) {
	var totalCount *uint64
	query := squirrel.
		Select("id",
			"user_id",
			"coin_id",
			"account_type_id",
			"minpay",
			"address",
			"changed_at",
			"img1",
			"img2",
			"is_active",
			"created_at",
			"updated_at",
			"fee",
			"user_id_new",
			"coin_new",
		).From("emcd.users_accounts as ua").
		PlaceholderFormat(squirrel.Dollar)

	if filter.Pagination != nil {
		var totalCountScan uint64
		queryCount := squirrel.Select("count(*)").
			From("emcd.users_accounts as ua").
			PlaceholderFormat(squirrel.Dollar)

		queryCount = newUserAccountFilterSql(filter).ApplyToQuery(queryCount)

		if querySql, args, err := queryCount.ToSql(); err != nil {

			return nil, nil, fmt.Errorf("failed to sql %s: %w", repoName, err)
		} else if err := a.PgxTransactor.Runner(ctx).QueryRow(ctx, querySql, args...).Scan(&totalCountScan); err != nil {

			return nil, nil, fmt.Errorf("failed scan count %s: %w", repoName, err)
		} else {
			query = newPaginationSql(filter.Pagination).ApplyToQuery(query)
			query = query.OrderBy("ua.id ASC")
			totalCount = &totalCountScan

		}
	}

	query = newUserAccountFilterSql(filter).ApplyToQuery(query)

	if querySql, args, err := query.ToSql(); err != nil {

		return nil, nil, fmt.Errorf("failed to sql %s: %w", repoName, err)
	} else if rows, err := a.PgxTransactor.Runner(ctx).Query(ctx, querySql, args...); err != nil {

		return nil, nil, fmt.Errorf("failed query rows %s: %w", repoName, err)
	} else {
		defer rows.Close()

		var userAccountList model.UserAccounts
		for rows.Next() {
			var userAccount model.UserAccount
			if err := rows.Scan(&userAccount.ID,
				&userAccount.UserID,
				&userAccount.CoinID,
				&userAccount.AccountTypeID,
				&userAccount.Minpay,
				&userAccount.Address,
				&userAccount.ChangedAt,
				&userAccount.Img1,
				&userAccount.Img2,
				&userAccount.IsActive,
				&userAccount.CreatedAt,
				&userAccount.UpdatedAt,
				&userAccount.Fee,
				&userAccount.UserIDNew,
				&userAccount.CoinNew,
			); err != nil {

				return nil, nil, fmt.Errorf("failed scan rows %s: %w", repoName, err)
			}

			userAccountList = append(userAccountList, &userAccount)

		}

		return totalCount, userAccountList, nil
	}
}

// GetUserAccountByIdLegacy - deprecated, use FindUserAccountByFilter
func (a *userAccountsImpl) GetUserAccountByIdLegacy(ctx context.Context, id int32) (*model.UserAccount, error) {
	filter := &model.UserAccountFilter{
		ID:              &id,
		UserID:          nil,
		AccountTypeID:   nil,
		UserIDNew:       nil,
		CoinNew:         nil,
		IsActive:        nil,
		Pagination:      nil,
		UserIDNewIsNull: nil,
		CoinNewIsNull:   nil,
	}

	if _, userAccounts, err := a.FindUserAccountByFilter(ctx, filter); err != nil {

		return nil, err
	} else if len(userAccounts) == 0 {

		return nil, fmt.Errorf("user account is not found")
	} else if len(userAccounts) > 1 {

		return nil, fmt.Errorf("user account is not uniq id")
	} else {

		return userAccounts[0], nil
	}
}

// FindUserAccountByUserIdLegacy - deprecated, use FindUserAccountByFilter
func (a *userAccountsImpl) FindUserAccountByUserIdLegacy(ctx context.Context, userIdNew uuid.UUID) (model.UserAccounts, error) {
	filter := &model.UserAccountFilter{
		ID:              nil,
		UserID:          nil,
		AccountTypeID:   nil,
		UserIDNew:       &userIdNew,
		CoinNew:         nil,
		IsActive:        nil,
		Pagination:      nil,
		UserIDNewIsNull: nil,
		CoinNewIsNull:   nil,
	}

	_, userAccounts, err := a.FindUserAccountByFilter(ctx, filter)

	return userAccounts, err
}

func (a *userAccountsImpl) UpdateUserAccountByPartial(ctx context.Context, userAccount *model.UserAccount, partial *model.UserAccountPartial) error {
	query := squirrel.
		Update("emcd.users_accounts").
		Where(squirrel.Eq{"id": userAccount.ID})

	query = newUserAccountPartialSql(partial).ApplyToQuery(query)
	query = query.PlaceholderFormat(squirrel.Dollar)

	if querySql, args, err := query.ToSql(); err != nil {
		return fmt.Errorf("failed to sql %s: %w", repoName, err)

	} else if _, err := a.Runner(ctx).Exec(ctx, querySql, args...); err != nil {

		return fmt.Errorf("failed exec %s: %w", repoName, err)
	} else {
		userAccount.Update(partial)

		return nil
	}
}

func (a *userAccountsImpl) UpdateUserAccountForMigrateUserIdNew(ctx context.Context, userAccounts model.UserAccounts) error {
	const queryUpdate = `update emcd.users_accounts as c set user_id_new = j.new_id from emcd.users AS j where c.user_id = j.id and c.id = @id`

	return a.updateUserAccountForMigrate(ctx, userAccounts, queryUpdate)
}

func (a *userAccountsImpl) UpdateUserAccountForMigrateCoinNew(ctx context.Context, userAccounts model.UserAccounts) error {
	const queryUpdate = `update emcd.users_accounts as c set coin_new = j.code from emcd.coins AS j where c.coin_id = j.id and c.id = @id`

	return a.updateUserAccountForMigrate(ctx, userAccounts, queryUpdate)
}

func (a *userAccountsImpl) updateUserAccountForMigrate(ctx context.Context, userAccounts model.UserAccounts, queryUpdate string) error {
	batch := pgx.Batch{}

	for _, userAccount := range userAccounts {
		batch.Queue(queryUpdate, newUserAccountSql(userAccount).ToNamedArgs())

	}

	batchResult := a.PgxTransactor.Runner(ctx).SendBatch(ctx, &batch)
	defer func(br pgx.BatchResults) {
		if err := br.Close(); err != nil {
			sdkLog.Error(ctx, "failed to close batch: %s", err.Error())

		}
	}(batchResult)

	if _, err := batchResult.Exec(); err != nil {

		return fmt.Errorf("failed batch query %s: %w", repoName, err)
	}

	return nil
}

type userAccountSql struct {
	*model.UserAccount
}

func newUserAccountSql(userAccount *model.UserAccount) *userAccountSql {

	return &userAccountSql{
		UserAccount: userAccount,
	}
}

func (u *userAccountSql) Validate(userId int32, userIdNew uuid.UUID) error {
	switch {
	case u.UserID != userId:
		return fmt.Errorf("mismatch user id %s", repoName)
	case u.UserIDNew.UUID != userIdNew:
		return fmt.Errorf("mismatch user id new %s", repoName)
	default:
		return nil
	}
}

func (u *userAccountSql) ToNamedArgs() *pgx.NamedArgs {

	r := &pgx.NamedArgs{
		"id":              u.ID,
		"user_id":         u.UserID,
		"coin_id":         u.CoinID,
		"account_type_id": u.AccountTypeID.AccountTypeId,
		"minpay":          u.Minpay,
		"address":         u.Address,
		"changed_at":      u.ChangedAt,
		"img1":            u.Img1,
		"img2":            u.Img2,
		"is_active":       u.IsActive,
		"created_at":      u.CreatedAt,
		"updated_at":      u.UpdatedAt,
		"user_id_new":     u.UserIDNew.UUID,
		"coin_new":        u.CoinNew,
	}

	return r
}
