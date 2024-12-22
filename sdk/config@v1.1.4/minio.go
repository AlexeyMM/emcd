package config

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type MinioConfig struct {
	Endpoint  string `env:"MINIO_ENDPOINT"   envDefault:"localhost:9001"`
	AccessKey string `env:"MINIO_ACCESS_KEY" envDefault:"minio123"`
	Secret    string `env:"MINIO_SECRET"     envDefault:"minio123"`
	Secure    bool   `env:"MINIO_SECURE"     envDefault:"true"`
}

func (c MinioConfig) New(provider trace.TracerProvider) (*minio.Client, error) {
	t, err := minio.DefaultTransport(c.Secure)
	if err != nil {
		return nil, fmt.Errorf("creating minio transport: %w", err)
	}

	return minio.New(c.Endpoint, &minio.Options{
		Creds:           credentials.NewStaticV4(c.AccessKey, c.Secret, ""),
		Secure:          c.Secure,
		Transport:       otelhttp.NewTransport(t, otelhttp.WithTracerProvider(provider)),
		Region:          "",
		BucketLookup:    0,
		TrailingHeaders: false,
		CustomMD5:       nil,
		CustomSHA256:    nil,
	})
}
