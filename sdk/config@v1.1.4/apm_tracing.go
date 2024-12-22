package config

import (
	"fmt"

	"github.com/pkg/errors"
	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/transport"
)

type APMTracing transport.HTTPTransportOptions

func (cfg APMTracing) New(serviceName, serviceVersion, environment string) (*apm.Tracer, error) {
	transport, err := transport.NewHTTPTransport(transport.HTTPTransportOptions(cfg))
	if err != nil {
		return nil, fmt.Errorf("create transport for apm: %w", errors.WithStack(err))
	}
	tracer, err := apm.NewTracerOptions(apm.TracerOptions{
		ServiceName:        serviceName,
		ServiceVersion:     serviceVersion,
		ServiceEnvironment: environment,
		Transport:          transport,
	})
	if err != nil {
		return nil, fmt.Errorf("create tracer: %w", errors.WithStack(err))
	}
	return tracer, nil
}
