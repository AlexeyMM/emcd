package profile

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	pgTx "code.emcdtech.com/emcd/sdk/pg"
	"code.emcdtech.com/emcd/service/profile/internal/repository"
	"code.emcdtech.com/emcd/service/profile/internal/server"
	"code.emcdtech.com/emcd/service/profile/protocol/profile"
	"code.emcdtech.com/emcd/service/profile/tests"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type OldNewIDsSuite struct {
	suite.Suite

	dbPool *pgxpool.Pool
	docker *tests.Docker

	handler     *server.Profile
	oldUserRepo repository.OldUsers
}

func TestNewOldIDs(t *testing.T) {
	suite.Run(t, new(OldNewIDsSuite))
}

func (s *OldNewIDsSuite) TestGetOldIDs() {
	id1 := "10cbfbc2-b1ab-45ad-af79-951860973223"

	oldNewIDs, err := s.handler.GetOldIDByID(context.Background(), &profile.GetOldIDByIDRequest{
		Id: id1,
	})
	s.Require().NoError(err)

	s.Equal(int32(1), oldNewIDs.OldId)
}

func (s *OldNewIDsSuite) SetupSuite() {
	s.docker = tests.NewDocker(s.T())
	s.dbPool = s.docker.GetPgxPool(s.T())

	s.T().Logf("database: %s", s.dbPool.Config().ConnString())
	db, err := sql.Open("postgres", s.dbPool.Config().ConnString())
	s.Require().NoError(err)
	defer db.Close()
	err = db.Ping()
	s.Require().NoError(err)

	fixturesPath, err := filepath.Abs("./fixtures/old_new_ids/")
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

	trx := pgTx.NewPgxTransactor(s.dbPool)

	s.oldUserRepo = repository.NewOldUsers(trx)
	s.handler = server.NewProfile(nil, s.oldUserRepo)
}

func (s *OldNewIDsSuite) TearDownSuite() {
	s.docker.Destroy(context.Background())
}
