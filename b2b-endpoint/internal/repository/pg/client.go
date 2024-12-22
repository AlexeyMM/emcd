package pg

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/endpoint/internal/model"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool *pgxpool.Pool
}

func NewClient(pool *pgxpool.Pool) *Client {
	return &Client{
		pool: pool,
	}
}

func (c *Client) Add(ctx context.Context, client *model.Client) error {
	sql, args, err := squirrel.
		Insert("endpoint.clients").
		Columns(
			"id",
			"name").
		Values(
			client.ID,
			client.Name).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = c.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
