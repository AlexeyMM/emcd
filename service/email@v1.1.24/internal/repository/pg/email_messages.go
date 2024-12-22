package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

type EmailMessagesStore struct {
	pool *pgxpool.Pool
}

func NewEmailMessages(pool *pgxpool.Pool) *EmailMessagesStore {
	return &EmailMessagesStore{
		pool: pool,
	}
}

func (e *EmailMessagesStore) Create(ctx context.Context, em *model.EmailMessageEvent) error {
	createEventQuery := `INSERT INTO sent_email_messages (id,email,type,created_at) VALUES ($1,$2,$3,$4)`
	_, err := e.pool.Exec(ctx, createEventQuery, em.ID, em.Email, em.Type, em.CreatedAt)
	if err != nil {
		return fmt.Errorf("email messages: create: %w", err)
	}
	return nil
}

func (e *EmailMessagesStore) ListMessages(
	ctx context.Context,
	email *string,
	eventType *string,
	skip, take int32,
) ([]*model.EmailMessageEvent, int, error) {
	batch := pgx.Batch{}

	query := `
				SELECT id,
					   email,
					   type,
					   created_at
				FROM sent_email_messages
				WHERE (email = $1 OR Coalesce($1,'') = '') and (type = $2 OR Coalesce($2,'') = '')
				LIMIT $3 OFFSET $4
				`

	batch.Queue(query, email, eventType, take, skip)

	query = `
				SELECT COUNT(*)
				FROM sent_email_messages
				WHERE (email = $1 OR Coalesce($1,'') = '') and (type = $2 OR Coalesce($2,'') = '')
				`

	batch.Queue(query, email, eventType)

	res := e.pool.SendBatch(ctx, &batch)
	defer func() {
		if err := res.Close(); err != nil {
			log.Error(ctx, "sent_email_messages.List: close batch: %v", err)
		}
	}()

	rows, err := res.Query()
	if err != nil {
		return nil, 0, fmt.Errorf("query %w", err)
	}

	defer rows.Close()

	list := make([]*model.EmailMessageEvent, 0)

	for rows.Next() {
		o := new(model.EmailMessageEvent)

		err = rows.Scan(
			&o.ID,
			&o.Email,
			&o.Type,
			&o.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan %w", err)
		}

		list = append(list, o)
	}

	var totalCount int

	err = res.QueryRow().Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("queryRow %w", err)
	}

	return list, totalCount, nil
}
