package pg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"code.emcdtech.com/b2b/endpoint/internal/encryptor"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
)

var (
	db      *pgxpool.Pool
	encrypt encryptor.Encryptor
)

func TestMain(m *testing.M) {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal(context.Background(), "Could not connect to docker: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresResource := initializePostgres(ctx, dockerPool, newPostgresConfig())

	encrypt = encryptor.NewEncryptor("12345678901234567890123456789012")

	code := m.Run()

	purgeResources(dockerPool, postgresResource)

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

		db, err = pgxpool.New(ctx, cfg.getConnectionString(dbHostAndPort))
		if err != nil {
			return err
		}

		return db.Ping(ctx)
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
		"-locations=filesystem:../../../migrations",
		fmt.Sprintf("-url=jdbc:postgresql://%v/postgres", dbHostAndPort),
		"migrate",
	}
}

func truncateAll(ctx context.Context) error {
	tables := []string{
		"endpoint.clients",
		"endpoint.secrets",
		"endpoint.whitelist_ips",
		"endpoint.request_logs",
	}

	for _, table := range tables {
		_, err := db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			return err
		}
	}
	return nil
}
