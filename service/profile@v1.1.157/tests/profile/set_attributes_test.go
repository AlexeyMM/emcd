package profile

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	pgTx "code.emcdtech.com/emcd/sdk/pg"
	"github.com/brianvoe/gofakeit"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"go.openly.dev/pointy"

	"code.emcdtech.com/emcd/service/profile/internal/repository"
	"code.emcdtech.com/emcd/service/profile/internal/server"
	"code.emcdtech.com/emcd/service/profile/internal/service"
	"code.emcdtech.com/emcd/service/profile/protocol/profile"
	"code.emcdtech.com/emcd/service/profile/tests"
)

type SetAttributesSuite struct {
	suite.Suite

	dbPool *pgxpool.Pool
	docker *tests.Docker

	trx     pgTx.PgxTransactor
	profile service.Profile
	handler *server.Profile
}

func TestSetAttributes(t *testing.T) {
	suite.Run(t, new(SetAttributesSuite))
}

func (s *SetAttributesSuite) TestSetAttributes() {
	userID, err := uuid.Parse("10cbfbc2-b1ab-45ad-af79-951860973223")
	s.Require().NoError(err)

	ctx := context.Background()
	changes := &profile.SetUserAttributesRequest{
		UserId:                   userID.String(),
		Username:                 pointy.Pointer(gofakeit.Username()),
		Language:                 pointy.String("en"),
		ParentId:                 pointy.String(gofakeit.UUID()),
		WhiteLabelID:             pointy.String(gofakeit.UUID()),
		WasReferralLinkGenerated: pointy.Bool(true),
		IsAmbassador:             pointy.Bool(true),
		PoolType:                 pointy.String("pool"),
	}
	_, err = s.handler.SetUserAttributes(ctx, changes)
	s.Require().NoError(err)

	user, err := s.profile.GetByUserID(ctx, userID)
	s.Require().NoError(err)

	s.Equal(*changes.Username, user.User.Username)
	s.Equal(*changes.Language, user.User.Language)
	s.Equal(*changes.ParentId, user.User.ParentID.String())
	s.Equal(*changes.WhiteLabelID, user.User.WhiteLabelID.String())
	s.Equal(*changes.WasReferralLinkGenerated, user.User.WasReferralLinkGenerated)
	s.Equal(*changes.IsAmbassador, user.User.IsAmbassador)
	s.Equal(*changes.PoolType, user.User.PoolType)
}

func (s *SetAttributesSuite) TestGetAttributes() {
	userID, err := uuid.Parse("10cbfbc2-b1ab-45ad-af79-951860973223")
	s.Require().NoError(err)

	ctx := context.Background()
	changes := &profile.SetUserAttributesRequest{
		UserId:       userID.String(),
		IsAmbassador: pointy.Bool(true),
	}
	_, err = s.handler.SetUserAttributes(ctx, changes)
	s.Require().NoError(err)

	user, err := s.handler.GetByUserID(ctx, &profile.GetByUserIDRequest{
		UserID: userID.String(),
	})
	s.Require().NoError(err)

	s.Equal(*changes.IsAmbassador, user.Profile.User.IsAmbassador)
}

func (s *SetAttributesSuite) SetupSuite() {
	s.docker = tests.NewDocker(s.T())
	s.dbPool = s.docker.GetPgxPool(s.T())

	s.T().Logf("database: %s", s.dbPool.Config().ConnString())
	db, err := sql.Open("postgres", s.dbPool.Config().ConnString())
	s.Require().NoError(err)
	defer db.Close()
	err = db.Ping()
	s.Require().NoError(err)

	fixturesPath, err := filepath.Abs("./fixtures/set_attributes/")
	s.Require().NoError(err)
	s.T().Logf("Test data is loading from %+v", fixturesPath)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(fixturesPath),
	)
	s.Require().NoError(err)
	s.T().Log("test data is loaded")

	err = fixtures.Load()
	s.Require().NoError(err)

	s.trx = pgTx.NewPgxTransactor(s.dbPool)
	profileRepo := repository.NewProfile(s.trx, s.trx)
	oldUserRepo := repository.NewOldUsers(s.trx)

	// can initialize only using dependencies
	s.profile = service.NewProfile(
		profileRepo,
		oldUserRepo,
		nil,
		nil,
		nil,
		"",
		nil,
		0,
		nil,
		nil,
		0,
		nil,
		nil,
		nil,
		nil,
		nil)
	s.handler = server.NewProfile(s.profile)
}

func (s *SetAttributesSuite) TearDownSuite() {
	s.docker.Destroy(context.Background())
}
