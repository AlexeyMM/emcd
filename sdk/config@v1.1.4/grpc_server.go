package config

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type GRPCServer struct {
	ListenAddr                     string        `env:"GRPC_LISTEN_ADDR"                        envDefault:":9090"`
	UseReflection                  bool          `env:"GRPC_USE_REFLECTION"                     envDefault:"True"`
	KeepaliveMaxConnectionIdle     time.Duration `env:"GRPC_KEEPALIVE_MAX_CONNECTION_IDLE"      envDefault:"0s"`
	KeepaliveMaxConnectionAge      time.Duration `env:"GRPC_KEEPALIVE_MAX_CONNECTION_AGE"       envDefault:"0s"`
	KeepaliveMaxConnectionAgeGrace time.Duration `env:"GRPC_KEEPALIVE_MAX_CONNECTION_AGE_GRACE" envDefault:"0s"`
	KeepaliveTime                  time.Duration `env:"GRPC_KEEPALIVE_TIME"                     envDefault:"0s"`
	KeepaliveTimeout               time.Duration `env:"GRPC_KEEPALIVE_TIMEOUT"                  envDefault:"0s"`
}

func (cfg GRPCServer) New(opts ...grpc.ServerOption) *grpc.Server {
	opts = append(opts, grpc.KeepaliveParams(
		keepalive.ServerParameters{
			MaxConnectionIdle:     cfg.KeepaliveMaxConnectionIdle,
			MaxConnectionAge:      cfg.KeepaliveMaxConnectionAge,
			MaxConnectionAgeGrace: cfg.KeepaliveMaxConnectionAgeGrace,
			Time:                  cfg.KeepaliveTime,
			Timeout:               cfg.KeepaliveTimeout,
		},
	))
	return grpc.NewServer(opts...)
}
