package pg_test

import (
	"context"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/email/internal/model"
	"code.emcdtech.com/emcd/service/email/internal/repository"
)

func (s *testPg) TestProviderSettings_CRU() {
	setting := s.newRandomSetting()
	s.Run("create", func() {
		err := s.provideSettingsStore.Create(context.Background(), setting)
		s.Require().NoError(err)
	})
	s.Run("create (empty provider)", func() {
		emptySettings := s.newRandomSetting()
		emptySettings.Providers = emptySettings.Providers[0:0]
		err := s.provideSettingsStore.Create(context.Background(), emptySettings)
		s.Require().NoError(err)
		v, err := s.provideSettingsStore.Get(context.Background(), emptySettings.WhiteLabelID)
		s.Require().NoError(err)
		emptySettings.Providers = nil
		s.Equal(emptySettings, v)
	})
	s.Run("create (nil provider)", func() {
		emptySettings := s.newRandomSetting()
		emptySettings.Providers = nil
		err := s.provideSettingsStore.Create(context.Background(), emptySettings)
		s.Require().NoError(err)
		v, err := s.provideSettingsStore.Get(context.Background(), emptySettings.WhiteLabelID)
		s.Require().NoError(err)
		s.Equal(emptySettings, v)
	})
	s.Run("get", func() {
		v, err := s.provideSettingsStore.Get(context.Background(), setting.WhiteLabelID)
		s.Require().NoError(err)
		s.Equal(setting, v)
	})
	s.Run("update", func() {
		upd := s.newRandomSetting()
		upd.WhiteLabelID = setting.WhiteLabelID
		err := s.provideSettingsStore.Update(context.Background(), upd)
		s.Require().NoError(err)
		getSetting, err := s.provideSettingsStore.Get(context.Background(), setting.WhiteLabelID)
		s.Require().NoError(err)
		s.Equal(upd, getSetting)
		s.NotEqual(setting, getSetting)
	})
	s.Run("get (not found)", func() {
		_, err := s.provideSettingsStore.Get(context.Background(), uuid.New())
		s.Require().ErrorIs(err, repository.ErrNotFound)
	})
}

func (s *testPg) TestProviderSettings_List() {
	const count = 7
	testData := make(map[uuid.UUID]model.Setting, count)
	for range count {
		e := s.newRandomSetting()
		testData[e.WhiteLabelID] = e
		err := s.provideSettingsStore.Create(context.Background(), e)
		s.Require().NoError(err)
		err = s.provideSettingsStore.Create(context.Background(), e)
		s.Require().Error(err)
	}
	pg := repository.NewPagination(count / 2)
	for pg.Offset() < count {
		list, err := s.provideSettingsStore.List(context.Background(), pg)
		s.Require().NoError(err)
		for i := range list {
			_, ok := testData[list[i].WhiteLabelID]
			s.Require().True(ok)
			// TODO добавить сравнение
			delete(testData, list[i].WhiteLabelID)
		}
		pg = pg.Next()
	}
	s.Require().True(len(testData) == 0)
}
