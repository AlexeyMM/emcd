package tests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"code.emcdtech.com/emcd/sdk/config"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

type Docker struct {
	*dockertest.Pool

	Network     *docker.Network
	dbContainer *dockertest.Resource
	resources   []*dockertest.Resource
}

func (d *Docker) Destroy(ctx context.Context) {
	for i := range d.resources {
		if err := d.Purge(d.resources[i]); err != nil {
			log.Error(ctx, "could not purge resource: %s", err.Error())
		}
		err := d.resources[i].Expire(1)
		if err != nil {
			log.Error(ctx, err.Error())
		}
	}
	err := d.Pool.Client.RemoveNetwork(d.Network.ID)
	if err != nil {
		log.Error(ctx, err.Error())
	}
}

func (d *Docker) dockerHostOpts(cfg *docker.HostConfig) {
	cfg.AutoRemove = false
	cfg.RestartPolicy = docker.RestartPolicy{
		Name: "no",
	}
}

func NewDocker(t *testing.T) *Docker {
	const networkName = "test-network"
	dockerPool, err := dockertest.NewPool("")
	require.NoError(t, err)

	network, err := dockerPool.Client.NetworkInfo(networkName)
	noErr := &docker.NoSuchNetwork{}
	if errors.As(err, &noErr) {
		network, err = dockerPool.Client.CreateNetwork(docker.CreateNetworkOptions{
			Name: networkName,
		})
	}
	require.NoError(t, err)

	err = dockerPool.Client.Ping()
	require.NoError(t, err)

	return &Docker{
		Pool:    dockerPool,
		Network: network,
	}
}

func (d *Docker) GetPgxPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	var pgPool *pgxpool.Pool
	sync.OnceFunc(func() {
		getConnectionString := func(port string) string {
			gitlabCIHost := os.Getenv("DATABASE_HOST")
			if gitlabCIHost == "" {
				gitlabCIHost = "localhost"
			}
			return fmt.Sprintf("postgres://user:secret@%s:%s/profile_test?sslmode=disable", gitlabCIHost, port)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		runOptions := &dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "latest",
			Env: []string{
				"POSTGRES_USER=user",
				"POSTGRES_PASSWORD=secret",
				"POSTGRES_DB=profile_test", // should has suffix _test for testfixtures pkg
				"listen_addresses = '*'",
			},
			NetworkID: d.Network.ID,
		}
		resource, err := d.RunWithOptions(runOptions, d.dockerHostOpts)
		require.NoError(t, err)
		err = resource.Expire(300)
		require.NoError(t, err)

		d.Pool.MaxWait = 60 * time.Second

		var port string

		log.Info(ctx, "connecting string: %s", getConnectionString(port))
		makePgxPool := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			port = resource.GetPort("5432/tcp")
			pgPool, err = config.NewPGXPoolFromDSN(ctx, noop.NewTracerProvider(), getConnectionString(port))
			if err != nil {
				return err
			}

			return pgPool.Ping(ctx)
		}

		err = d.Pool.Retry(makePgxPool)
		require.NoError(t, err)

		d.dbContainer = resource
		d.resources = append(d.resources, resource)
		d.applyMigration(t, "public", "../../migrations")
		d.applyMigration(t, "emcd", "../migrations")
	})()
	return pgPool
}

func (d *Docker) applyMigration(t *testing.T, scheme, relativePath string) {
	pathMigration, err := filepath.Abs(relativePath)
	require.NoError(t, err)
	t.Logf("migrations path %q: %s", scheme, pathMigration)

	host := d.dbContainer.Container.Name
	runOptions := &dockertest.RunOptions{
		Repository: "flyway/flyway",
		Tag:        "latest",
		Cmd: []string{
			fmt.Sprintf("-url=jdbc:postgresql:/%s:%s/profile_test ", host, "5432"),
			"-user=user",
			"-password=secret",
			fmt.Sprintf("-schemas=%s", scheme),
			"-validateMigrationNaming=true",
			"-connectRetries=3",
			"migrate",
		},
		NetworkID: d.Network.ID,
		Mounts: []string{
			pathMigration + ":/flyway/sql",
		},
	}

	resource, err := d.RunWithOptions(runOptions, d.dockerHostOpts)
	require.NoError(t, err)

	err = resource.Expire(60)
	require.NoError(t, err)

	running := true
	t.Logf("waiting for migration %q to finish", scheme)
	for running {
		ctn, err := d.Client.InspectContainer(resource.Container.ID)
		if err != nil {
			t.Logf("up migration: %s", err.Error())
			break
		}
		running = ctn.State.Running
		time.Sleep(1 * time.Second)
	}
	t.Logf("migration %q finished", scheme)
	d.resources = append(d.resources, resource)
}
