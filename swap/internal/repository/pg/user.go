package pg

import (
	"context"
	"errors"
	"fmt"

	"code.emcdtech.com/b2b/swap/model"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	errUserNotFound = errors.New("user not found")
)

type User struct {
	db *pgxpool.Pool
	transactor.PgxTransactor
}

func NewUser(db *pgxpool.Pool) *User {
	return &User{
		db:            db,
		PgxTransactor: transactor.NewPgxTransactor(db),
	}
}

func (u *User) Add(ctx context.Context, user *model.User) error {
	query := squirrel.Insert("swap.users").
		Columns(
			"id",
			"email",
			"language").
		Values(
			user.ID,
			user.Email,
			user.Language,
		).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = u.Runner(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (u *User) FindOne(ctx context.Context, filter *model.UserFilter) (*model.User, error) {
	users, err := u.find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}
	if len(users) != 1 {
		if len(users) == 0 {
			return nil, errUserNotFound
		}

		return nil, fmt.Errorf("unexpected number of users: %d", len(users))
	}

	return users[0], nil
}

func (u *User) find(ctx context.Context, filter *model.UserFilter) ([]*model.User, error) {
	query := squirrel.Select(
		"id",
		"email",
		"language").
		From("swap.users").
		PlaceholderFormat(squirrel.Dollar)

	query = newUserFilterSql(filter).applyToSelectQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	rows, err := u.Runner(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	if rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Email, &user.Language)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return users, nil
}
