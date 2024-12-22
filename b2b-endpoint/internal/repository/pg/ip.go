package pg

import (
	"context"
	"fmt"

	"code.emcdtech.com/b2b/endpoint/internal/model"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IP struct {
	pool *pgxpool.Pool
}

func NewIP(pool *pgxpool.Pool) *IP {
	return &IP{
		pool: pool,
	}
}

func (i *IP) Add(ctx context.Context, ips []*model.IP) error {
	query := squirrel.
		Insert("endpoint.whitelist_ips").
		Columns(
			"api_key",
			"ip_address",
			"created_at").
		PlaceholderFormat(squirrel.Dollar)

	for _, ip := range ips {
		query = query.Values(
			ip.ApiKey,
			ip.Address,
			ip.CreatedAt.UTC())
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = i.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (i *IP) Find(ctx context.Context, filter *model.IPFilter) ([]*model.IP, error) {
	query := squirrel.Select(
		"api_key",
		"ip_address",
		"created_at").
		From("endpoint.whitelist_ips").
		PlaceholderFormat(squirrel.Dollar)

	query = newIPFilterSql(filter).applyToSelectQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	rows, err := i.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var ips []*model.IP
	for rows.Next() {
		var (
			ip model.IP
		)
		err = rows.Scan(&ip.ApiKey, &ip.Address, &ip.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		ips = append(ips, &ip)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return ips, nil
}

func (i *IP) Delete(ctx context.Context, filter *model.IPFilter) error {
	query := squirrel.Delete("endpoint.whitelist_ips").
		PlaceholderFormat(squirrel.Dollar)

	query = newIPFilterSql(filter).applyToDeleteQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = i.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
