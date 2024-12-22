package pg_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace/noop"

	"code.emcdtech.com/emcd/sdk/config"

	"code.emcdtech.com/emcd/service/email/internal/model"
	"code.emcdtech.com/emcd/service/email/internal/repository/pg"
)

const (
	DotEnvFilename = "../../../.local.env"
)

func TestPg(t *testing.T) {
	suite.Run(t, &testPg{})
}

type testPg struct {
	suite.Suite
	pool          *pgxpool.Pool
	templateStore *pg.TemplateStore
	// sentEmailStore       *pg.EmailSendingAuditLogStore
	provideSettingsStore *pg.ProvideSettingsStore
}

func (s *testPg) SetupSuite() {
	err := godotenv.Load(DotEnvFilename)
	s.Require().NoError(err)
	var cfg config.PGXPool
	err = env.Parse(&cfg)
	s.Require().NoError(err)
	s.pool, err = cfg.New(context.Background(), noop.NewTracerProvider())
	s.Require().NoError(err)
	s.Require().NotNil(s.pool)

	s.templateStore = pg.NewTemplateStore(s.pool)
	// s.sentEmailStore = pg.NewEmailSendingAuditLogStore(s.pool)
	s.provideSettingsStore = pg.NewProvideSettingsStore(s.pool)
}

func (s *testPg) SetupTest() {
	s.pool.Exec(context.Background(), "TRUNCATE TABLE email_templates")
	// s.pool.Exec(context.Background(), "TRUNCATE TABLE sent_email_messages")
	s.pool.Exec(context.Background(), "TRUNCATE TABLE provider_settings")
}

func (s *testPg) randomString() string {
	return random.String(32, random.Alphanumeric)
}

func (s *testPg) newRandomTemplate() model.Template {
	return model.Template{
		WhiteLabelID: uuid.New(),
		Language:     s.randomString(),
		Template:     model.NewTextTemplate(s.randomString()),
		Type:         model.NewCodeTemplate(s.randomString()),
		Subject:      model.NewTextTemplate(s.randomString()),
		Footer:       s.randomString(),
	}
}

func (s *testPg) newRandomSetting() model.Setting {
	smtpSetting, err := model.GetProvider(model.SmtpSetting{
		Username:               s.randomString(),
		Password:               s.randomString(),
		ServerAddress:          s.randomString(),
		ServerPort:             rand.Int(),
		FromAddress:            s.randomString(),
		FromAddressDisplayedAs: s.randomString(),
	})
	s.Require().NoError(err)
	//nolint
	mailgun, err := model.GetProvider(model.MailgunSetting{
		Domain:                 s.randomString(),
		ApiKey:                 s.randomString(),
		FromAddress:            s.randomString(),
		ApiBase:                s.randomString(),
		FromAddressDisplayedAs: s.randomString(),
	})
	return model.Setting{
		WhiteLabelID: uuid.New(),
		Providers:    []model.Provider{smtpSetting, mailgun},
		CreatedAt:    time.Now().Add(time.Duration(rand.Int()) * time.Second).UTC().Truncate(time.Second),
		UpdatedAt:    time.Now().Add(time.Duration(rand.Int()) * time.Second).UTC().Truncate(time.Second),
	}
}
