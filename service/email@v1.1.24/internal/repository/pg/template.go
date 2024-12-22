package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"code.emcdtech.com/emcd/service/email/internal/model"
	"code.emcdtech.com/emcd/service/email/internal/repository"
)

type TemplateStore struct {
	pool *pgxpool.Pool
}

func (s *TemplateStore) Create(ctx context.Context, t model.Template) error {
	const createTemplateSQL = `
INSERT INTO email_templates(whitelabel_id, template, language, type, subject, footer) 
VALUES ($1, $2, $3, $4, $5, $6)
`
	_, err := s.pool.Exec(ctx, createTemplateSQL, t.WhiteLabelID, t.Template, t.Language, t.Type, t.Subject, t.Footer)
	if err != nil {
		return fmt.Errorf("exect createTemplateSQL: %w", err)
	}
	return nil
}

func (s *TemplateStore) Get(
	ctx context.Context,
	whitelabelID uuid.UUID,
	language string,
	_type model.CodeTemplate,
) (model.Template, error) {
	const getTemplateSQL = `
SELECT whitelabel_id, template, language, type, subject, footer 
  FROM email_templates 
 WHERE whitelabel_id = $1 AND language = $2 AND type = $3
`
	var r model.Template
	err := s.pool.
		QueryRow(ctx, getTemplateSQL, whitelabelID, language, _type).
		Scan(&r.WhiteLabelID, &r.Template, &r.Language, &r.Type, &r.Subject, &r.Footer)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		if errors.Is(err, pgx.ErrNoRows) {
			return r, repository.ErrNotFound
		}
		return r, fmt.Errorf("exect getTemplateSQL: %w", err)
	}
	return r, nil
}

func (s *TemplateStore) Update(ctx context.Context, t model.Template) error {
	const updateTemplateSQL = `
UPDATE email_templates
   SET template = $1,
       subject = $2,
       footer = $3
 WHERE whitelabel_id = $4 AND language = $5 AND type = $6
`
	_, err := s.pool.Exec(ctx, updateTemplateSQL, t.Template, t.Subject, t.Footer, t.WhiteLabelID, t.Language, t.Type)
	if err != nil {
		return fmt.Errorf("exect updateTemplateSQL: %w", err)
	}
	return nil
}

func (s *TemplateStore) List(ctx context.Context, p repository.Pagination) ([]model.Template, error) {
	const listTemplatesSQL = `
SELECT whitelabel_id, template, language, type, subject, footer 
  FROM email_templates
 ORDER BY whitelabel_id, language, type  
 LIMIT $1 OFFSET $2
`
	rows, err := s.pool.Query(ctx, listTemplatesSQL, p.Size, p.Offset())
	if err != nil {
		return nil, fmt.Errorf("exect listTemplatesSQL: %w", err)
	}
	defer rows.Close()
	result := make([]model.Template, 0, p.Size)
	var template model.Template
	for rows.Next() {
		err = rows.Scan(&template.WhiteLabelID, &template.Template, &template.Language, &template.Type,
			&template.Subject, &template.Footer)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		result = append(result, template)
	}
	return result, nil
}

func NewTemplateStore(pool *pgxpool.Pool) *TemplateStore {
	return &TemplateStore{
		pool: pool,
	}
}
