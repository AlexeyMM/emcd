package rabbitmqkit

import (
	"maps"
	"slices"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/propagation"
)

var _ propagation.TextMapCarrier = &RabbitmqHeadersCarrier{}

type RabbitmqHeadersCarrier amqp.Table

func (c RabbitmqHeadersCarrier) Get(key string) string {
	v, ok := c[key]
	if !ok {
		return ""
	}

	return v.(string) //nolint:forcetypeassert
}

func (c RabbitmqHeadersCarrier) Set(key string, value string) {
	c[key] = value
}

func (c RabbitmqHeadersCarrier) Keys() []string {
	return slices.Collect(maps.Keys(c))
}
