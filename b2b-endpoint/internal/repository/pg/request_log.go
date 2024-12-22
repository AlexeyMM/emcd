package pg

import (
	"context"
	"errors"
	"fmt"

	"code.emcdtech.com/b2b/endpoint/internal/model"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var LogsNotFoundError = errors.New("logs not found")

type RequestLog struct {
	pool *pgxpool.Pool
}

func NewRequestLog(pool *pgxpool.Pool) *RequestLog {
	return &RequestLog{
		pool: pool,
	}
}

func (r *RequestLog) Add(ctx context.Context, log *model.RequestLog) error {
	sql, args, err := squirrel.
		Insert("endpoint.request_logs").
		Columns(
			"api_key",
			"request_hash",
			"created_at").
		Values(
			log.ApiKey,
			log.RequestHash,
			log.CreatedAt.UTC()).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("toSql: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *RequestLog) Find(ctx context.Context, filter *model.RequestLogFilter) ([]*model.RequestLog, error) {
	return r.find(ctx, filter)
}

func (r *RequestLog) FindOne(ctx context.Context, filter *model.RequestLogFilter) (*model.RequestLog, error) {
	logs, err := r.find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(logs) == 1 {
		return logs[0], nil
	} else if len(logs) > 1 {
		return nil, fmt.Errorf("unexpected number of logs: %d", len(logs))
	}
	return nil, LogsNotFoundError
}

func (r *RequestLog) find(ctx context.Context, filter *model.RequestLogFilter) ([]*model.RequestLog, error) {
	query := squirrel.Select(
		"api_key",
		"request_hash",
		"created_at").
		From("endpoint.request_logs").
		PlaceholderFormat(squirrel.Dollar)

	query = newRequestLogFilterSql(filter).applyToSelectQuery(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var logs []*model.RequestLog
	for rows.Next() {
		var log model.RequestLog
		err = rows.Scan(
			&log.ApiKey,
			&log.RequestHash,
			&log.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		logs = append(logs, &log)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return logs, nil
}
