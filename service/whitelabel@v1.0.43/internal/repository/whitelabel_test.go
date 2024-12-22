package repository

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/service/whitelabel/internal/model"
)

func TestWhiteLabel_Create(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)
	expected := &model.WhiteLabel{
		ID:     uuid.New(),
		Domain: "dewiohsorgi 324 qreho",
	}
	err := wlRepo.Create(ctx, expected)
	require.NoError(t, err)
	actual, total, err := wlRepo.GetAll(ctx, 0, 1, "domain", false)
	require.NoError(t, err)
	require.Equal(t, 1, len(actual))
	require.Equal(t, *expected, *actual[0])
	require.Equal(t, 1, total)
}

func TestWhiteLabel_GetAll(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)

	expected := make([]*model.WhiteLabel, 26)
	for i := 0; i < 26; i++ {
		prefix := byte('a' + i)
		sb := strings.Builder{}
		sb.WriteByte(prefix)
		expected[i] = &model.WhiteLabel{
			ID:     uuid.New(),
			Domain: sb.String(),
		}
		err := wlRepo.Create(ctx, expected[i])
		require.NoError(t, err)
		_, err = dbPool.Exec(ctx, "UPDATE whitelabel SET user_id=$1 WHERE id=$2", i, expected[i].ID)
		expected[i].UserID = int32(i)
		require.NoError(t, err)
	}
	actual, total, err := wlRepo.GetAll(ctx, 15, 5, "domain", true)
	require.NoError(t, err)
	require.Equal(t, 5, len(actual))
	require.Equal(t, 26, total)
	for i := range actual {
		require.Equal(t, expected[i+15], actual[i])
	}
}

func TestWhiteLabel_GetAllSortCommissions(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)
	expected := &model.WhiteLabel{
		ID:     uuid.New(),
		Domain: "dewiohsorgi 324 qreho",
	}
	err := wlRepo.Create(ctx, expected)
	require.NoError(t, err)
	actual, total, err := wlRepo.GetAll(ctx, 0, 1, "domain", false)
	require.NoError(t, err)
	require.Equal(t, 1, len(actual))
	require.Equal(t, 1, total)
	require.Equal(t, expected, actual[0])
}

func TestWhiteLabel_Update(t *testing.T) {
	fmt.Println(time.Now().Add(-24 * time.Hour))
	ctx := context.Background()
	defer truncateAllWl(t, ctx)
	wl := &model.WhiteLabel{
		ID:     uuid.New(),
		Domain: "dewiohsorgi 324 qreho",
	}
	expected := &model.WhiteLabel{
		ID:     wl.ID,
		Domain: "dewiohsorgi 324 qreho",
	}
	err := wlRepo.Create(ctx, wl)
	require.NoError(t, err)
	err = wlRepo.Update(ctx, expected)
	require.NoError(t, err)
	actual, total, err := wlRepo.GetAll(ctx, 0, 1, "domain", false)
	require.NoError(t, err)
	require.Equal(t, 1, len(actual))
	require.Equal(t, expected, actual[0])
	require.Equal(t, 1, total)
}

func TestWhiteLabel_Delete(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)
	wl := &model.WhiteLabel{
		ID:     uuid.New(),
		Domain: "dewiohsorgi 324 qreho",
	}
	err := wlRepo.Create(ctx, wl)
	require.NoError(t, err)
	err = wlRepo.Delete(ctx, wl.ID)
	require.NoError(t, err)
	actual, total, err := wlRepo.GetAll(ctx, 0, 1, "domain", false)
	require.NoError(t, err)
	require.Equal(t, 0, len(actual))
	require.Equal(t, total, 0)
}

func TestWhiteLabel_GetByID(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)
	expected := &model.WhiteLabel{
		ID:     uuid.New(),
		Domain: "dewiohsorgi 324 qreho",
	}
	err := wlRepo.Create(ctx, expected)
	require.NoError(t, err)
	actual, err := wlRepo.GetByID(ctx, expected.ID)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestWhiteLabel_GetBySegmentID(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)

	var segmentID int32 = 1234

	expected := &model.WhiteLabel{
		ID:          uuid.New(),
		UserID:      1111,
		SegmentID:   segmentID,
		Origin:      "test_get_by_segment_id",
		Prefix:      "t",
		SenderEmail: "test_get_by_segment_id@emcd.io",
		Domain:      "test_get_by_segment_id",
		APIKey:      "test_get_by_segment_id",
		URL:         "test_get_by_segment_id",
		Version:     1,
	}
	err := wlRepo.Create(ctx, expected)
	require.NoError(t, err)
	actual, err := wlRepo.GetBySegmentID(ctx, int(segmentID))
	require.NoError(t, err)
	//verifyEqualWhiteLabels(t, expected, actual)
	verifyEqualWhiteLabelsWithoutCommissions(t, expected, actual)
}

func TestWhiteLabel_GetByUserID(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)

	var userID int32 = 1234

	expected := &model.WhiteLabel{
		ID:          uuid.New(),
		UserID:      userID,
		SegmentID:   1111,
		Origin:      "Origin",
		Prefix:      "t",
		SenderEmail: "SenderEmail@emcd.io",
		Domain:      "Domain",
		APIKey:      "APIKey",
		URL:         "URL",
		Version:     1,
	}
	err := wlRepo.Create(ctx, expected)
	require.NoError(t, err)
	actual, err := wlRepo.GetByUserID(ctx, int(userID))
	require.NoError(t, err)
	//verifyEqualWhiteLabels(t, expected, actual)
	verifyEqualWhiteLabelsWithoutCommissions(t, expected, actual)
}

func TestWhiteLabel_GetByOrigin(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)

	var origin = "test_get_by_origin"

	expected := &model.WhiteLabel{
		ID:          uuid.New(),
		UserID:      1234,
		SegmentID:   1111,
		Origin:      origin,
		Prefix:      "t",
		SenderEmail: "test_get_by_origin@emcd.io",
		Domain:      "test_get_by_origin",
		APIKey:      "test_get_by_origin",
		URL:         "test_get_by_origin",
		Version:     1,
	}
	err := wlRepo.Create(ctx, expected)
	require.NoError(t, err)
	actual, err := wlRepo.GetByOrigin(ctx, origin)
	require.NoError(t, err)
	//verifyEqualWhiteLabels(t, expected, actual)
	verifyEqualWhiteLabelsWithoutCommissions(t, expected, actual)
}

func TestWhiteLabel_SetConfig(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)

	expected := &model.WlConfig{
		RefID:               "1234",
		MediaID:             "test_set_config",
		Origin:              "test_set_config",
		Colors:              model.Colors{},
		FirmwareInstruction: "instruction",
		Lang:                "ru",
	}
	err := wlRepo.Create(ctx, &model.WhiteLabel{
		ID:        uuid.New(),
		SegmentID: 1234,
		Origin:    expected.Origin,
	})
	require.NoError(t, err)
	err = wlRepo.SetConfigByRefID(ctx, expected)
	require.NoError(t, err)
	actual, err := wlRepo.GetConfigByOrigin(ctx, expected.Origin)
	require.NoError(t, err)
	require.Equal(t, expected.RefID, actual.RefID)
	require.Equal(t, expected.Origin, actual.Origin)
}

func TestWhiteLabel_SetConfigUpdate(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)

	expected := &model.WlConfig{
		RefID:               "1234",
		MediaID:             "test_set_config",
		Origin:              "test_set_config",
		Colors:              model.Colors{},
		FirmwareInstruction: "instruction",
		Lang:                "en",
	}
	err := wlRepo.Create(ctx, &model.WhiteLabel{
		ID:        uuid.New(),
		SegmentID: 1234,
		Origin:    expected.Origin,
	})
	require.NoError(t, err)
	err = wlRepo.SetConfigByRefID(ctx, expected)
	require.NoError(t, err)
	expected.FirmwareInstruction = "updated"
	expected.MediaID = "updated media id"
	err = wlRepo.SetConfigByRefID(ctx, expected)
	require.NoError(t, err)
	actual, err := wlRepo.GetConfigByOrigin(ctx, expected.Origin)
	require.NoError(t, err)
	require.Equal(t, expected.RefID, actual.RefID)
	require.Equal(t, expected.Origin, actual.Origin)
}

func TestWhiteLabel_GetConfig(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)

	expected := &model.WlConfig{
		RefID:   "1234",
		MediaID: "test_get_config",
		Origin:  "test_get_config",
		Colors:  model.Colors{},
		Lang:    "ru",
	}
	err := wlRepo.Create(ctx, &model.WhiteLabel{
		ID:        uuid.New(),
		SegmentID: 1234,
		Origin:    expected.Origin,
	})
	require.NoError(t, err)
	err = wlRepo.SetConfigByRefID(ctx, expected)
	require.NoError(t, err)
	actual, err := wlRepo.GetConfigByOrigin(ctx, expected.Origin)
	require.NoError(t, err)
	require.Equal(t, expected.RefID, actual.RefID)
	require.Equal(t, expected.Origin, actual.Origin)
	require.Len(t, actual.PossibleLang, 1)
	require.Equal(t, actual.PossibleLang[0], "ru")
}

func verifyEqualWhiteLabelsWithoutCommissions(t *testing.T, a, b *model.WhiteLabel) {
	require.Equal(t, a.Domain, b.Domain)
	require.Equal(t, a.UserID, b.UserID)
	require.Equal(t, a.SegmentID, b.SegmentID)
	require.Equal(t, a.ID, b.ID, a.ID.String()+"   "+b.ID.String())
}

func truncateAllWl(t *testing.T, ctx context.Context) {
	_, err := dbPool.Exec(ctx, "TRUNCATE whitelabel CASCADE")
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, "TRUNCATE frontend_configs CASCADE")
	require.NoError(t, err)
}

func TestWhiteLabel_GetFullByUserID(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)
	expected := model.WhiteLabel{
		ID:          uuid.New(),
		UserID:      123,
		SegmentID:   4324,
		Origin:      "origin",
		Prefix:      "prefix",
		SenderEmail: "sender",
		Domain:      "domain",
		APIKey:      "apiKey",
		URL:         "url",
		Version:     2,
		MasterSlave: true,
		MasterFee:   1454,
	}
	err := wlRepo.Create(ctx, &expected)
	require.NoError(t, err)
	actual, err := wlRepo.GetFullByUserID(ctx, int(expected.UserID))
	require.NoError(t, err)
	require.Equal(t, &expected, actual)
}

func TestGetCoins(t *testing.T) {
	validUUID := uuid.New()
	testData := []*model.WLCoins{{WlID: validUUID, CoinID: "btc"}, {WlID: validUUID, CoinID: "eth"}}

	cases := []struct {
		name        string
		wlId        uuid.UUID
		expectCoins []*model.WLCoins
	}{
		{
			name:        "ValidUUID",
			wlId:        validUUID,
			expectCoins: testData,
		},
		{
			name:        "NonExistentUUID",
			wlId:        uuid.UUID{},
			expectCoins: []*model.WLCoins{},
		},
	}

	ctx := context.Background()

	for _, coin := range testData {
		err := wlRepo.AddCoin(ctx, coin.WlID, coin.CoinID)
		require.NoError(t, err, "failed to insert test data")
	}

	defer func(UUID uuid.UUID, coinID string) {
		err := wlRepo.DeleteCoin(ctx, UUID, coinID)
		require.NoError(t, err, "failed to clean up test data")
	}(validUUID, "btc")

	defer func(UUID uuid.UUID, coinID string) {
		err := wlRepo.DeleteCoin(ctx, UUID, coinID)
		require.NoError(t, err, "failed to clean up test data")
	}(validUUID, "eth")

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			coins, err := wlRepo.GetCoins(ctx, tc.wlId)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.expectCoins, coins, "The expected and actual results do not match")
		})
	}
}

func TestWhiteLabel_GetAllEmptyCommissions(t *testing.T) {
	ctx := context.Background()
	defer truncateAllWl(t, ctx)

	expected := make([]*model.WhiteLabel, 26)
	for i := 0; i < 26; i++ {
		prefix := byte('a' + i)
		sb := strings.Builder{}
		sb.WriteByte(prefix)
		expected[i] = &model.WhiteLabel{
			ID:     uuid.New(),
			Domain: sb.String(),
		}
		err := wlRepo.Create(ctx, expected[i])
		require.NoError(t, err)
	}
	actual, total, err := wlRepo.GetAll(ctx, 15, 5, "domain", true)
	require.NoError(t, err)
	require.Equal(t, 5, len(actual))
	require.Equal(t, 26, total)

	for i := range actual {
		require.Equal(t, expected[i+15], actual[i])
	}
}
