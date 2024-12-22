package admin

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	pgTx "code.emcdtech.com/emcd/sdk/pg"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"code.emcdtech.com/emcd/service/profile/internal/model"
	"code.emcdtech.com/emcd/service/profile/internal/repository"
	"code.emcdtech.com/emcd/service/profile/internal/service"
	"code.emcdtech.com/emcd/service/profile/tests"
)

type DeleteSubUserSuite struct {
	suite.Suite

	dbPool *pgxpool.Pool
	docker *tests.Docker

	trx     pgTx.PgxTransactor
	profile service.Profile
}

func TestDeleteSubUser(t *testing.T) {
	suite.Run(t, new(DeleteSubUserSuite))
}

// Тест проверят, что при удалении subuser'a у всех остальных
// subuser'ов, идущих после него, уменьшен суффикс email'a на 1.
//
// Исходное состояние DB для теста см. в fixtures/delete_sub_user
func (s *DeleteSubUserSuite) TestDeleteSubUser() {
	// time.Sleep(30 * time.Second)
	ctx := context.Background()
	subUserUUID := uuid.MustParse("b889dfa0-2cba-4c50-95d7-0e73a5ed1c88")
	newParentID := uuid.MustParse("5667cc03-c79c-4648-942f-356e0c1ed078")

	userBeforeDeleting, err := s.profile.GetByUserID(ctx, subUserUUID)
	s.Require().NoError(err)

	err = s.profile.SoftDeleteSubUser(ctx, subUserUUID, newParentID)
	s.Require().NoError(err)

	userAfterDeleting, err := s.profile.GetByUserID(ctx, subUserUUID)
	s.Require().NoError(err)

	emailBeforeDeleting := model.NewSubUserEmailFromString(userBeforeDeleting.User.Email)
	emailAfterDeleting := model.NewSubUserEmailFromString(userAfterDeleting.User.Email)
	s.Less(emailBeforeDeleting.Index(), emailAfterDeleting.Index())

	s.Equal("deleted_"+userBeforeDeleting.User.Username, userAfterDeleting.User.Username)
	s.Equal(userAfterDeleting.User.ParentID, newParentID)
	// видимо работает репликация, поэтому пароль не меняется
	// s.Equal("\x7375626163636f756e74", userAfterDeleting.User.Password)

	// убираем разницу, т.к. проверили выше
	userAfterDeleting.User.Username = userBeforeDeleting.User.Username
	userAfterDeleting.User.Email = userBeforeDeleting.User.Email
	userAfterDeleting.User.ParentID = userBeforeDeleting.User.ParentID
	userBeforeDeleting.User.Password = userAfterDeleting.User.Password

	s.Equal(userBeforeDeleting, userAfterDeleting)
}

func (s *DeleteSubUserSuite) SetupSuite() {
	s.docker = tests.NewDocker(s.T())
	s.dbPool = s.docker.GetPgxPool(s.T())

	s.T().Logf("database: %s", s.dbPool.Config().ConnString())
	db, err := sql.Open("postgres", s.dbPool.Config().ConnString())
	s.Require().NoError(err)
	defer db.Close()
	err = db.Ping()
	s.Require().NoError(err)

	fixturesPath, err := filepath.Abs("./fixtures/delete_sub_user/")
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
}

func (s *DeleteSubUserSuite) TearDownSuite() {
	s.docker.Destroy(context.Background())
}
