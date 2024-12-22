package pg_test

import (
	"context"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/email/internal/model"
	"code.emcdtech.com/emcd/service/email/internal/repository"
)

func (s *testPg) TestTemplateStore_CRU() {
	template := s.newRandomTemplate()
	s.Run("create", func() {
		err := s.templateStore.Create(context.Background(), template)
		s.Require().NoError(err)
	})
	s.Run("update", func() {
		template.Template = model.NewTextTemplate(s.randomString())
		template.Subject = model.NewTextTemplate(s.randomString())
		template.Footer = s.randomString()
		err := s.templateStore.Update(context.Background(), template)
		s.Require().NoError(err)
	})
	s.Run("get", func() {
		v, err := s.templateStore.Get(context.Background(), template.WhiteLabelID, template.Language, template.Type)
		s.Require().NoError(err)
		s.Require().EqualValues(template, v)
	})
	s.Run("get (not found)", func() {
		var empty model.Template
		v, err := s.templateStore.Get(
			context.Background(),
			template.WhiteLabelID,
			s.randomString(),
			model.NewCodeTemplate(s.randomString()),
		)
		s.Require().ErrorIs(err, repository.ErrNotFound)
		s.Require().EqualValues(empty, v)
	})
}

func (s *testPg) TestTemplateStore_List() {
	const count = 7
	templates := make(map[uuid.UUID]model.Template)
	for range count {
		template := s.newRandomTemplate()
		templates[template.WhiteLabelID] = template
	}
	for _, template := range templates {
		err := s.templateStore.Create(context.Background(), template)
		s.Require().NoError(err)
	}

	pg := repository.NewPagination(3)
	for pg.Offset() < count {
		list, err := s.templateStore.List(context.Background(), pg)
		s.Require().NoError(err)
		s.Require().True(len(list) <= pg.Size)
		for i := range list {
			want, ok := templates[list[i].WhiteLabelID]
			s.Require().True(ok)
			s.Require().EqualValues(want, list[i])
			delete(templates, list[i].WhiteLabelID)
		}
		pg = pg.Next()
	}
	s.Require().True(len(templates) == 0)
}
