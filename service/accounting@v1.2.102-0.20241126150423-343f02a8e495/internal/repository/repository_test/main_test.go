package repository_test

import (
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		sdkLog.Fatal(ctx, "can't connect to docker: %s", err)
	}

	postgresCfg := newPostgresConfig()
	postgresResource := initializePostgres(ctx, dockerPool, postgresCfg)
	code := m.Run()

	purgeResources(ctx, dockerPool, postgresResource)

	os.Exit(code)
}

func initializePostgres(ctx context.Context, dockerPool *dockertest.Pool, cfg *postgresConfig) *dockertest.Resource {
	resource, err := dockerPool.Run(cfg.dbName, cfg.version, cfg.envVariables)
	if err != nil {
		sdkLog.Fatal(ctx, "can't start resource: %s", err)
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
		dbPool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}

		return dbPool.Ping(ctx)
	})
	if err != nil {
		sdkLog.Fatal(ctx, "can't connect to database: %s", err)
	}

	migrationArgs := cfg.getFlywayMigrationArgs(dbHostAndPort)
	migrationArgsString := strings.Join(migrationArgs, " ")
	sdkLog.Info(ctx, "flyway args: %s", migrationArgsString)
	cmd := exec.Command("flyway", migrationArgs...)

	err = cmd.Run()
	if err != nil {
		sdkLog.Fatal(ctx, "can't run migrations: %s", err)
	}

	// remove extra flyway report files
	cmd = exec.Command("rm", "report.html", "report.json")
	err = cmd.Run()
	if err != nil {
		sdkLog.Error(ctx, "can't delete report files: %v", err)
	}

	return resource
}

func purgeResources(ctx context.Context, dockerPool *dockertest.Pool, resources ...*dockertest.Resource) {
	for i := range resources {
		if err := dockerPool.Purge(resources[i]); err != nil {
			sdkLog.Fatal(ctx, "can't purge resource: %s", err)
		}

		err := resources[i].Expire(1)
		if err != nil {
			sdkLog.Fatal(ctx, err.Error())
		}
	}

}

func (p *postgresConfig) getFlywayMigrationArgs(dbHostAndPort string) []string {
	return []string{
		"-user=postgres",
		"-password=password123",
		"-locations=filesystem:../migrations",
		fmt.Sprintf("-url=jdbc:postgresql://%s/postgres", dbHostAndPort),
		"migrate",
	}
}
