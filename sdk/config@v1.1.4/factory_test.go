package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/sdk/config"
)

func TestNew(t *testing.T) {
	type Config struct {
		Environment config.Environment
	}
	randString := random.String(32, random.Alphanumeric)
	err := os.Setenv("ENVIRONMENT", randString)
	require.NoError(t, err)
	if err != nil {
		require.NoError(t, err)
	}

	cfg, err := config.New[Config]()
	require.NoError(t, err)
	require.Equal(t, cfg.Environment.Name, randString)
}

func TestNew_WithPrefix(t *testing.T) {
	type Config struct {
		Environment config.Environment
	}
	randString := random.String(32, random.Alphanumeric)
	err := os.Setenv("MY_APP_ENVIRONMENT", randString)
	require.NoError(t, err)
	if err != nil {
		require.NoError(t, err)
	}

	cfg, err := config.New[Config](config.WithPrefix("MY_APP_"))
	require.NoError(t, err)
	require.Equal(t, cfg.Environment.Name, randString)
}

func ExampleNew() {
	type Config struct {
		Environment config.Environment
		Log         config.Log
		GRPC        config.GRPCServer
		Tracing     config.APMTracing
		PGXPool     config.PGXPool
		HTTP        config.HTTPServer
		Redis       config.RedisConfig
		Minio       config.MinioConfig
	}
	cfg, err := config.New[Config]()
	if err != nil {
		panic(err)
	}
	fmt.Sprintf(cfg.HTTP.ListenAddr)
}
