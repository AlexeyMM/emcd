package repository

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/ory/dockertest"

	"code.emcdtech.com/emcd/sdk/log"
	pgTx "code.emcdtech.com/emcd/sdk/pg"
)

const (
	DotEnvFilename = "../../.env"
)

var (
	transactor               pgTx.PgxTransactor
	dbPool                   *pgxpool.Pool
	profileRepo              *profile
	oldUsersRepo             *oldUsers
	notificationSettingsRepo *notificationSettings
)

func TestMain(m *testing.M) {
	godotenv.Load(DotEnvFilename)
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal(context.Background(), "Could not connect to docker: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var postgresResource *dockertest.Resource
	if os.Getenv("LOCAL_TEST") == "" {
		postgresResource = initializePostgres(ctx, dockerPool, newPostgresConfig())
	} else {
		dbPool, err = pgxpool.New(ctx, os.Getenv("DATABASE_HOST"))
		if err != nil {
			log.Fatal(context.Background(), "new pgxpool: %s", err)
		}
		if err = dbPool.Ping(ctx); err != nil {
			log.Fatal(context.Background(), "new pgxpool: %s", err)
		}
	}
	createEmcdTables(ctx)
	transactor = pgTx.NewPgxTransactor(dbPool)
	profileRepo = NewProfile(transactor, transactor)
	oldUsersRepo = NewOldUsers(transactor)
	notificationSettingsRepo = NewNotificationSettings(transactor)

	code := m.Run()
	if os.Getenv("LOCAL_TEST") == "" {
		purgeResources(dockerPool, postgresResource)
	}
	os.Exit(code)
}

func initializePostgres(ctx context.Context, dockerPool *dockertest.Pool, cfg *postgresConfig) *dockertest.Resource {
	resource, err := dockerPool.Run(cfg.Repository, cfg.Version, cfg.EnvVariables)
	if err != nil {
		log.Fatal(ctx, "Could not start resource: %s", err)
	}

	var dbHostAndPort string

	err = dockerPool.Retry(func() error {
		var dbHost string
		gitlabCIHost := os.Getenv("DATABASE_HOST")
		if gitlabCIHost != "" {
			dbHost = gitlabCIHost
		} else {
			dbHost = "localhost"
		}
		port := resource.GetPort("5432/tcp")
		dbHostAndPort = fmt.Sprintf("%s:%s", dbHost, port)

		dbPool, err = pgxpool.New(ctx, cfg.getConnectionString(dbHostAndPort))
		if err != nil {
			return err
		}

		return dbPool.Ping(ctx)
	})
	if err != nil {
		log.Fatal(ctx, "Could not connect to database: %s", err)
	}
	log.Info(ctx, strings.Join(cfg.getFlywayMigrationArgs(dbHostAndPort), " "))
	cmd := exec.Command("flyway", cfg.getFlywayMigrationArgs(dbHostAndPort)...)

	err = cmd.Run()
	if err != nil {
		log.Fatal(ctx, "There are errors in migrations: %v", err)
	}
	return resource
}

func purgeResources(dockerPool *dockertest.Pool, resources ...*dockertest.Resource) {
	bgCtx := context.Background()
	for i := range resources {
		if err := dockerPool.Purge(resources[i]); err != nil {
			log.Fatal(bgCtx, "Could not purge resource: %s", err)
		}

		err := resources[i].Expire(1)
		if err != nil {
			log.Fatal(bgCtx, err.Error())
		}
	}
}

type postgresConfig struct {
	Repository   string
	Version      string
	EnvVariables []string
	PortID       string
}

func newPostgresConfig() *postgresConfig {
	return &postgresConfig{
		Repository:   "postgres",
		Version:      "14.1-alpine",
		EnvVariables: []string{"POSTGRES_PASSWORD=password123"},
		PortID:       "5432/tcp",
	}
}

func (p *postgresConfig) getConnectionString(dbHostAndPort string) string {
	return fmt.Sprintf("postgresql://postgres:password123@%v/%s", dbHostAndPort, p.Repository)
}

func (p *postgresConfig) getFlywayMigrationArgs(dbHostAndPort string) []string {
	return []string{
		"-user=postgres",
		"-password=password123",
		"-locations=filesystem:../../migrations",
		fmt.Sprintf("-url=jdbc:postgresql://%v/postgres", dbHostAndPort),
		"migrate",
	}
}

func createEmcdTables(ctx context.Context) {
	_, err := dbPool.Exec(ctx, "CREATE SCHEMA emcd")
	if err != nil {
		log.Fatal(ctx, err.Error())
	}
	_, err = dbPool.Exec(ctx, `
		CREATE TABLE emcd.users
		(
			new_id   uuid,
			username varchar(100) NOT NULL,
			ref_id int NOT NULL DEFAULT 1,
			password text,
			created_at timestamp DEFAULT NOW(),
			api_key       varchar(100) NOT NULL,
			id SERIAL,
			email text UNIQUE,
			updated_at timestamp DEFAULT NOW(),
			suspended timestamp,
			language varchar(10),
			parent_id int,
			is_active bool,
			nopay bool,
			timezone varchar(20)
		);`)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}
	_, err = dbPool.Exec(ctx, `CREATE UNIQUE INDEX users_email_uindex ON emcd.users USING btree (LOWER(email));`)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}
	_, err = dbPool.Exec(ctx, "CREATE SCHEMA histories")
	if err != nil {
		log.Fatal(ctx, err.Error())
	}
	_, err = dbPool.Exec(ctx, `CREATE TABLE histories.segment_userids (
    user_id int PRIMARY KEY,
    segment_id serial
)`)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}
	_, err = dbPool.Exec(ctx, `CREATE TABLE emcd.vip_users (
    	user_id int
)`)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}
	_, err = dbPool.Exec(ctx, `CREATE TABLE emcd.coins (
    	id SERIAL PRIMARY KEY,
    	code text,
    	name text
)`)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}
	_, err = dbPool.Exec(ctx, `CREATE TABLE emcd.accounts_referral (
    	id SERIAL PRIMARY KEY,
    	account_id integer NOT NULL,
    tier integer,
    referral_fee numeric NOT NULL,
    active_referrals integer,
    coin_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
)`)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}
	_, err = dbPool.Exec(ctx, `CREATE TABLE emcd.users_accounts (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    coin_id integer NOT NULL,
    account_type_id integer NOT NULL,
    minpay numeric,
    address character varying COLLATE pg_catalog."default",
    changed_at timestamp without time zone,
    img1 numeric,
    img2 numeric,
    is_active boolean,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    fee numeric DEFAULT 0.015
    )`)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}

	_, err = dbPool.Exec(ctx, `CREATE TABLE emcd.accounts_pool(
    id int,
    emcd_address_autopay bool,
    account_id int
)`)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}

	_, err = dbPool.Exec(ctx, `CREATE TABLE emcd.autopay_addresses (
    id SERIAL,
    user_account_id int,
    address varchar (200),
    percent int
)`)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}

	_, err = dbPool.Exec(ctx, `CREATE TABLE emcd.histories (
    user_id int,
    segment_id int
)`)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}

	_, err = dbPool.Exec(ctx, `CREATE TABLE emcd.user_logs (
		id int,
		user_id int,
		change_type varchar(30),
		ip varchar(50),
		token varchar(200),
		old_value varchar(30),
		value varchar(30),
		active bool,
		used bool,
		is_segment_sended bool
)`)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}

	insertCoins(ctx)
}
