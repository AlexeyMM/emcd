package repository_migration_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
	"github.com/rs/zerolog/log"
)

type postgresConfig struct {
	dbName       string
	version      string
	envVariables []string
	portID       string
}

var (
	dbPool *pgxpool.Pool
)

func newPostgresConfig() *postgresConfig {
	return &postgresConfig{
		dbName:       "postgres",
		version:      "14.1-alpine",
		envVariables: []string{"POSTGRES_PASSWORD=password123"},
		portID:       "5432/tcp",
	}
}

func TestMain(m *testing.M) {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal().Msgf("can't connect to docker: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresCfg := newPostgresConfig()
	postgresResource := initializePostgres(ctx, dockerPool, postgresCfg)
	code := m.Run()

	purgeResources(dockerPool, postgresResource)

	os.Exit(code)
}

func initializePostgres(ctx context.Context, dockerPool *dockertest.Pool, cfg *postgresConfig) *dockertest.Resource {
	resource, err := dockerPool.Run(cfg.dbName, cfg.version, cfg.envVariables)
	if err != nil {
		log.Fatal().Msgf("can't start resource: %s", err)
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
		dsn := fmt.Sprintf("postgresql://postgres:password123@%s/%s", dbHostAndPort, cfg.dbName)

		pgxConfig, err := pgxpool.ParseConfig(dsn)
		pgxConfig.MaxConns = 100

		if err != nil {
			return err
		}
		dbPool, err = pgxpool.NewWithConfig(ctx, pgxConfig)
		if err != nil {
			return err
		}

		return dbPool.Ping(ctx)
	})
	if err != nil {
		log.Fatal().Msgf("can't connect to database: %s", err)
	}

	migrationArgs := cfg.getFlywayMigrationArgs(dbHostAndPort)
	migrationArgsString := strings.Join(migrationArgs, " ")
	log.Info().Msgf("flyway args: %s", migrationArgsString)
	cmd := exec.Command("flyway", migrationArgs...)

	err = cmd.Run()
	if err != nil {
		log.Fatal().Msgf("can't run migrations: %s", err)
	}

	// remove extra flyway report files
	cmd = exec.Command("rm", "report.html", "report.json")
	err = cmd.Run()
	if err != nil {
		log.Error().Msgf("can't delete report files: %s", err)
	}

	return resource
}

func purgeResources(dockerPool *dockertest.Pool, resources ...*dockertest.Resource) {
	for i := range resources {
		if err := dockerPool.Purge(resources[i]); err != nil {
			log.Fatal().Msgf("can't purge resource: %s", err)
		}

		err := resources[i].Expire(1)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}

}

func (p *postgresConfig) getFlywayMigrationArgs(dbHostAndPort string) []string {
	return []string{
		"-user=postgres",
		"-password=password123",
		"-locations=filesystem:./migrations_test",
		fmt.Sprintf("-url=jdbc:postgresql://%s/postgres", dbHostAndPort),
		"migrate",
	}
}
